package models

import(
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"
)

type Session struct {
	Token string
	UserID int 
	Expiry time.Time
}

type SessionModel struct {
	DB *sql.DB
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (m *SessionModel) Insert(userID int, duration time.Duration) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	expiry := time.Now().Add(duration)
	stmt := `INSERT INTO sessions (token, user_id, expiry) VALUES (?, ?, ?)`

	_, err = m.DB.Exec(stmt, token, userID, expiry)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *SessionModel) GetUserID(token string) (int, error) {
	var userID int
	var expiry time.Time

	stmt := `SELECT user_id, expiry FROM sessions WHERE token = ?`
	err := m.DB.QueryRow(stmt, token).Scan(&userID, &expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoRecord
		}
		return 0, err
	}

	if time.Now().After(expiry) {
		_ = m.Delete(token)
		return 0, ErrNoRecord
	}
	return userID, nil
}

func (m *SessionModel) Delete(token string) error {
	stmt := `DELETE FROM sessions WHERE token = ?`
	_, err := m.DB.Exec(stmt, token)
	return err
}