package models

import (
	"database/sql"
	"errors"
)

var ErrInvalidReaction = errors.New("models: invalid reaction type")

type ReactionModel struct {
	DB *sql.DB
}

func validReaction(reactionType string) bool {
	return reactionType == "like" || reactionType == "dislike"
}

// ToggleForPost creates a reaction, changes it, or removes it when the same
// reaction is pressed for a second time.
func (m *ReactionModel) ToggleForPost(userID, postID int, reactionType string) error {
	if !validReaction(reactionType) {
		return ErrInvalidReaction
	}
	return m.toggle(userID, postID, reactionType, true)
}

func (m *ReactionModel) ToggleForComment(userID, commentID int, reactionType string) error {
	if !validReaction(reactionType) {
		return ErrInvalidReaction
	}
	return m.toggle(userID, commentID, reactionType, false)
}

func (m *ReactionModel) toggle(userID, targetID int, reactionType string, isPost bool) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	targetColumn := "comment_id"
	otherColumn := "post_id"
	if isPost {
		targetColumn = "post_id"
		otherColumn = "comment_id"
	}

	var current string
	err = tx.QueryRow(`SELECT reaction_type FROM reactions WHERE user_id = ? AND `+targetColumn+` = ?`, userID, targetID).Scan(&current)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = tx.Exec(`INSERT INTO reactions (user_id, `+targetColumn+`, `+otherColumn+`, reaction_type, created_at)
			VALUES (?, ?, NULL, ?, DATETIME('now'))`, userID, targetID, reactionType)
	case err != nil:
		return err
	case current == reactionType:
		_, err = tx.Exec(`DELETE FROM reactions WHERE user_id = ? AND `+targetColumn+` = ?`, userID, targetID)
	default:
		_, err = tx.Exec(`UPDATE reactions SET reaction_type = ?, created_at = DATETIME('now')
			WHERE user_id = ? AND `+targetColumn+` = ?`, reactionType, userID, targetID)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}
