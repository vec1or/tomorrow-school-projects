package models

import (
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
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

	stmt1 := `INSERT INTO posts (user_id, title, content, created_at)
					VALUES(?, ?, ?, DATETIME('now'))`

	result, err := tx.Exec(stmt1, userID, title, content)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	stmt2 := `INSERT INTO post_categories (post_id, category_id) VALUES(?, ?)`
	for _, catID := range categoryIDs {
		_, err := tx.Exec(stmt2, postID, catID)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (m *PostModel) GetByID(id int) (*Post, error) {
	stmt := `SELECT id, user_id, title, content, created_at FROM posts WHERE id = ?`

	p := &Post{}
	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("models: no matching record found")
		}
		return nil, err
	}
	return p, nil
}

func (m *PostModel) Latest() ([]*Post, error) {
	stmt := `SELECT id, user_id, title, content, created_at 
					FROM posts
					ORDER BY created_at DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		p := &Post{}
		err = rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
