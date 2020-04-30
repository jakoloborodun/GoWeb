package models

import (
	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	Title       string
	Description string
}

func (cat *Category) Create(db *gorm.DB) {
	db.Create(&cat)
}

func (cat *Category) Update(db *gorm.DB) {
	db.Save(&cat)
}

func (cat *Category) Delete(db *gorm.DB) {
	db.Delete(&cat)
}

func GetCategory(id int64, db *gorm.DB) (cat Category) {
	db.First(&cat, id)

	return
}

func GetAllCategories(db *gorm.DB) (categories []Category) {
	db.Find(&categories)

	return
}

func NewCategory(title, description string) *Category {
	return &Category{
		Title:       title,
		Description: description,
	}
}
