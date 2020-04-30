package models

import (
	"database/sql"
)

type Category struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"desc"`
}

type Categories []Category

func (cat *Category) Create(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO Category (title, description) VALUES (?, ?)",
		cat.Title, cat.Description)

	return err
}

func (cat *Category) Update(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE Category SET title = ?, description = ? WHERE id = ?",
		cat.Title, cat.Description, cat.ID,
	)

	return err
}

func (cat *Category) Delete(db *sql.DB) error {
	_, err := db.Exec(
		"DELETE FROM Category WHERE id = ?",
		cat.ID,
	)

	return err
}

func GetCategory(id int64, db *sql.DB) (Category, error) {
	cat := Category{}

	row := db.QueryRow("SELECT * FROM Category WHERE id = ?", id)
	err := row.Scan(&cat.ID, &cat.Title, &cat.Description)
	if err != nil {
		return cat, err
	}

	return cat, nil
}

func GetAllCategories(db *sql.DB) (Categories, error) {
	rows, err := db.Query("SELECT * FROM Category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make(Categories, 0, 10)
	for rows.Next() {
		c := Category{}
		if err = rows.Scan(&c.ID, &c.Title, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, err
}

func NewCategory(id int64, title, description string) *Category {
	return &Category{
		ID:          id,
		Title:       title,
		Description: description,
	}
}
