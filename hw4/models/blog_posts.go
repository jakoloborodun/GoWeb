package models

import (
	"database/sql"
	"time"
)

type BlogPost struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	Created  time.Time `json:"created"`
	Status   bool      `json:"status"`
	Category *Category `json:"category"`
	Content  string    `json:"content"`
}

type BlogPosts []BlogPost

func (post *BlogPost) Create(db *sql.DB) error {
	_, err := db.Query(
		"INSERT INTO BlogPost (title, text, created, status, category_id, content) VALUES (?, ?, ?, ?, ?, ?)",
		post.Title, post.Text, post.Created.Unix(), post.Status, post.Category.ID, post.Content)

	return err
}

func (post *BlogPost) Update(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE BlogPost SET title = ?, text = ?, status = ?, category_id = ?, content = ? WHERE id = ?",
		post.Title, post.Text, post.Status, post.Category.ID, post.Content, post.ID,
	)

	return err
}

func (post *BlogPost) Delete(db *sql.DB) error {
	_, err := db.Exec(
		"DELETE FROM BlogPost WHERE id = ?",
		post.ID,
	)

	return err
}

func GetPost(id int64, db *sql.DB) (BlogPost, error) {
	post := BlogPost{}
	var categoryId int64
	var timestamp int64

	row := db.QueryRow("SELECT * FROM BlogPost WHERE id = ?", id)
	err := row.Scan(&post.ID, &post.Title, &post.Text, &timestamp, &post.Status, &categoryId, &post.Content)
	if err != nil {
		return post, err
	}

	post.Created = time.Unix(timestamp, 0)
	category, err := GetCategory(categoryId, db)
	if err != nil {
		return post, err
	}
	post.Category = &category

	return post, nil
}

func GetAllPosts(db *sql.DB) (BlogPosts, error) {
	rows, err := db.Query("SELECT id, title, text, created, status, category_id, content FROM BlogPost")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make(BlogPosts, 0, 10)
	var categoryId int64
	var timestamp int64

	for rows.Next() {
		p := BlogPost{}
		if err = rows.Scan(&p.ID, &p.Title, &p.Text, &timestamp, &p.Status, &categoryId, &p.Content); err != nil {
			return nil, err
		}
		p.Created = time.Unix(timestamp, 0)
		category, err := GetCategory(categoryId, db)
		if err != nil {
			return posts, err
		}
		p.Category = &category

		posts = append(posts, p)
	}

	return posts, err
}

func GetPostsByCategory(cid int64, db *sql.DB) (BlogPosts, error) {
	rows, err := db.Query("SELECT * FROM BlogPost WHERE category_id = ?", cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := BlogPosts{}
	var categoryId int64
	var timestamp int64

	for rows.Next() {
		p := BlogPost{}
		if err = rows.Scan(&p.ID, &p.Title, &p.Text, &timestamp, &p.Status, &categoryId, &p.Content); err != nil {
			return nil, err
		}
		p.Created = time.Unix(timestamp, 0)
		category, err := GetCategory(categoryId, db)
		if err != nil {
			return posts, err
		}
		p.Category = &category

		posts = append(posts, p)
	}

	return posts, err
}

func NewBlogPost(id int64, title string, text string, status bool, category *Category, content string) *BlogPost {
	return &BlogPost{
		ID:       id,
		Title:    title,
		Text:     text,
		Created:  time.Now(),
		Status:   status,
		Category: category,
		Content:  content,
	}
}
