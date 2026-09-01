package models

import (
	"database/sql"
	"errors"
)

type Category struct {
	ID   int
	Name string
}

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) All() ([]*Category, error) {
	stmt := `SELECT id, name FROM categories ORDER BY name ASC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		c := &Category{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *CategoryModel) ByPostID(postID int) ([]*Category, error) {
	stmt := `
		SELECT c.id, c.name
		FROM categories c
		INNER JOIN post_categories pc ON pc.category_id = c.id
		WHERE pc.post_id = ?
		ORDER BY c.name ASC`

	rows, err := m.DB.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		c := &Category{}
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *CategoryModel) Insert(name string) (int, error) {
	stmt := `INSERT INTO categories (name) VALUES (?)`
	result, err := m.DB.Exec(stmt, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
