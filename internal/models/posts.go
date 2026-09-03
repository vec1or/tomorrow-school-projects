package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Post struct {
	ID           int
	UserID       int
	Username     string
	Title        string
	Content      string
	CreatedAt    time.Time
	Categories   []*Category
	LikeCount    int
	DislikeCount int
	CommentCount int
}

type PostFilters struct {
	CategoryID      int
	CreatedByUserID int
	LikedByUserID   int
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(userID int, title, content string, categoryIDs []int) (int, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`INSERT INTO posts (user_id, title, content, created_at)
		VALUES (?, ?, ?, DATETIME('now'))`, userID, title, content)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, categoryID := range categoryIDs {
		_, err = tx.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, categoryID)
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return int(postID), nil
}

const postSelect = `
	SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
		(SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type = 'like'),
		(SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type = 'dislike'),
		(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id)
	FROM posts p
	INNER JOIN users u ON u.id = p.user_id`

func (m *PostModel) GetByID(id int) (*Post, error) {
	p := &Post{}
	err := m.DB.QueryRow(postSelect+` WHERE p.id = ?`, id).Scan(
		&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt,
		&p.LikeCount, &p.DislikeCount, &p.CommentCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	p.Categories, err = categoriesForPost(m.DB, p.ID)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (m *PostModel) Latest() ([]*Post, error) {
	return m.Filter(PostFilters{})
}

func (m *PostModel) Filter(filters PostFilters) ([]*Post, error) {
	query := postSelect
	var conditions []string
	var args []any

	if filters.CategoryID > 0 {
		conditions = append(conditions, `EXISTS (
			SELECT 1 FROM post_categories pc
			WHERE pc.post_id = p.id AND pc.category_id = ?
		)`)
		args = append(args, filters.CategoryID)
	}
	if filters.CreatedByUserID > 0 {
		conditions = append(conditions, `p.user_id = ?`)
		args = append(args, filters.CreatedByUserID)
	}
	if filters.LikedByUserID > 0 {
		conditions = append(conditions, `EXISTS (
			SELECT 1 FROM reactions ur
			WHERE ur.post_id = p.id AND ur.user_id = ? AND ur.reaction_type = 'like'
		)`)
		args = append(args, filters.LikedByUserID)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += ` ORDER BY p.created_at DESC, p.id DESC`

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var posts []*Post
	for rows.Next() {
		p := &Post{}
		err = rows.Scan(
			&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt,
			&p.LikeCount, &p.DislikeCount, &p.CommentCount,
		)
		if err != nil {
			rows.Close()
			return nil, err
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}

	// The database intentionally uses one connection, so category lookups happen
	// only after the post rows are closed.
	for _, p := range posts {
		p.Categories, err = categoriesForPost(m.DB, p.ID)
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func categoriesForPost(db *sql.DB, postID int) ([]*Category, error) {
	rows, err := db.Query(`
		SELECT c.id, c.name
		FROM categories c
		INNER JOIN post_categories pc ON pc.category_id = c.id
		WHERE pc.post_id = ?
		ORDER BY c.name`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category := &Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}
