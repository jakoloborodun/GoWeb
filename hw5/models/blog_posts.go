package models

import (
	"github.com/jinzhu/gorm"
)

type BlogPost struct {
	gorm.Model
	Title    string
	Text     string
	Status   bool
	Category *Category
	Content  string
}

func (post *BlogPost) Create(db *gorm.DB) {
	db.Save(&post)
}

func (post *BlogPost) Update(db *gorm.DB) {
	db.Update(&post)
}

func (post *BlogPost) Delete(db *gorm.DB) {
	db.Delete(&post)
}

func GetPost(id int64, db *gorm.DB) (post BlogPost) {
	db.First(&post, id)

	return
}

func GetAllPosts(db *gorm.DB) (posts []BlogPost) {
	db.Find(&posts)

	return
}

func NewBlogPost(title string, text string, status bool, category *Category, content string) *BlogPost {
	return &BlogPost{
		Title:    title,
		Text:     text,
		Status:   status,
		Category: category,
		Content:  content,
	}
}
