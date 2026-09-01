package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Username  string
	Content   string
	CreatedAt time.Time
}

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(postID, userID int, content string) (int, error) {
	stmt := `
		INSERT INTO comments (post_id, user_id, content, created_at)
		VALUES (?, ?, ?, DATETIME('now'))`

	result, err := m.DB.Exec(stmt, postID, userID, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *CommentModel) ByPostID(postID int) ([]*Comment, error) {
	stmt := `
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at
		FROM comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC`

	rows, err := m.DB.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		c := &Comment{}
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
