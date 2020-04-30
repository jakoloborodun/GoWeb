package models

import (
	"github.com/jinzhu/gorm"
)

type BlogPost struct {
	gorm.Model
	Title      string
	Text       string
	Status     bool
	Category   *Category `gorm:"foreignkey:CategoryID"`
	CategoryID int
	Content    string
}

func (post *BlogPost) Create(db *gorm.DB) {
	db.Save(&post)
}

func (post *BlogPost) Update(db *gorm.DB) {
	db.Save(&post)
}

func (post *BlogPost) Delete(db *gorm.DB) {
	db.Delete(&post)
}

func GetPost(id int64, db *gorm.DB) (post BlogPost) {
	db.Preload("Category").First(&post, id)
	return
}

func GetAllPosts(db *gorm.DB) (posts []BlogPost) {
	db.Preload("Category").Find(&posts)
	return
}

func GetPostsByCategory(cid int64, db *gorm.DB) (posts []BlogPost) {
	db.Preload("Category").Find(&posts, "category_id = ?", cid)
	return
}

func NewBlogPost(title string, text string, status bool, category int, content string) *BlogPost {
	return &BlogPost{
		Title:      title,
		Text:       text,
		Status:     status,
		CategoryID: category,
		Content:    content,
	}
}
