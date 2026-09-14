package models

import (
	"database/sql"
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
		JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?`

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

func (m *CategoryModel) Exists(id int) (bool, error) {
	var exists bool
	err := m.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)`, id).Scan(&exists)
	return exists, err
}

func (m *CategoryModel) AllExist(ids []int) (bool, error) {
	for _, id := range ids {
		exists, err := m.Exists(id)
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
	}
	return true, nil
}

func (m *CategoryModel) Insert(name string) (int, error) {
	stmt := `INSERT INTO categories (name) VALUES (?)`
	result, err := m.DB.Exec(stmt, name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
