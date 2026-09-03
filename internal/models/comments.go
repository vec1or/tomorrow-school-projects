package models

import (
	"database/sql"
	"errors"
	"time"
)

type Comment struct {
	ID           int
	PostID       int
	UserID       int
	Username     string
	Content      string
	CreatedAt    time.Time
	LikeCount    int
	DislikeCount int
}

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(postID, userID int, content string) (int, error) {
	result, err := m.DB.Exec(`
		INSERT INTO comments (post_id, user_id, content, created_at)
		VALUES (?, ?, ?, DATETIME('now'))`, postID, userID, content)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *CommentModel) GetByID(id int) (*Comment, error) {
	c := &Comment{}
	err := m.DB.QueryRow(`
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at,
			(SELECT COUNT(*) FROM reactions r WHERE r.comment_id = c.id AND r.reaction_type = 'like'),
			(SELECT COUNT(*) FROM reactions r WHERE r.comment_id = c.id AND r.reaction_type = 'dislike')
		FROM comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.id = ?`, id).Scan(
		&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt,
		&c.LikeCount, &c.DislikeCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return c, nil
}

func (m *CommentModel) ByPostID(postID int) ([]*Comment, error) {
	rows, err := m.DB.Query(`
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at,
			(SELECT COUNT(*) FROM reactions r WHERE r.comment_id = c.id AND r.reaction_type = 'like'),
			(SELECT COUNT(*) FROM reactions r WHERE r.comment_id = c.id AND r.reaction_type = 'dislike')
		FROM comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC, c.id ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		c := &Comment{}
		if err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt,
			&c.LikeCount, &c.DislikeCount,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
