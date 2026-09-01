package models

import (
	"database/sql"
	"errors"
	"time"
)

type Reaction struct {
	ID           int
	UserID       int
	PostID       int
	CommentID    int
	ReactionType string
	CreatedAt    time.Time
}

type ReactionModel struct {
	DB *sql.DB
}

func (m *ReactionModel) InsertForPost(userID, postID int, reactionType string) error {
	stmt := `
		INSERT INTO reactions (user_id, post_id, reaction_type, created_at)
		VALUES (?, ?, ?, DATETIME('now'))
		ON CONFLICT(user_id, post_id, comment_id) DO UPDATE SET reaction_type = excluded.reaction_type, created_at = DATETIME('now')`

	_, err := m.DB.Exec(stmt, userID, postID, reactionType)
	return err
}

func (m *ReactionModel) InsertForComment(userID, commentID int, reactionType string) error {
	stmt := `
		INSERT INTO reactions (user_id, comment_id, reaction_type, created_at)
		VALUES (?, ?, ?, DATETIME('now'))
		ON CONFLICT(user_id, post_id, comment_id) DO UPDATE SET reaction_type = excluded.reaction_type, created_at = DATETIME('now')`

	_, err := m.DB.Exec(stmt, userID, commentID, reactionType)
	return err
}

func (m *ReactionModel) CountForPost(postID int) (int, int, error) {
	likesStmt := `SELECT COUNT(*) FROM reactions WHERE post_id = ? AND reaction_type = 'like'`
	dislikesStmt := `SELECT COUNT(*) FROM reactions WHERE post_id = ? AND reaction_type = 'dislike'`

	var likes, dislikes int
	if err := m.DB.QueryRow(likesStmt, postID).Scan(&likes); err != nil {
		return 0, 0, err
	}
	if err := m.DB.QueryRow(dislikesStmt, postID).Scan(&dislikes); err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

func (m *ReactionModel) CountForComment(commentID int) (int, int, error) {
	likesStmt := `SELECT COUNT(*) FROM reactions WHERE comment_id = ? AND reaction_type = 'like'`
	dislikesStmt := `SELECT COUNT(*) FROM reactions WHERE comment_id = ? AND reaction_type = 'dislike'`

	var likes, dislikes int
	if err := m.DB.QueryRow(likesStmt, commentID).Scan(&likes); err != nil {
		return 0, 0, err
	}
	if err := m.DB.QueryRow(dislikesStmt, commentID).Scan(&dislikes); err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

var ErrInvalidReaction = errors.New("models: invalid reaction type")
