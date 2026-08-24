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

func (m *PostModel) Insert(userID int, title, content string) (int, error) {
	stmt := `INSERT INTO posts (user_id, title, content, created_at)
					VALUES(?, ?, ?, DATETIME('now'))`

	result, err := m.DB.Exec(stmt, userID, title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
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
