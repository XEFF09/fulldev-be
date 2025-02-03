package models

import (
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name		string	`json:"name"`
	Books		[]Book	`gorm:"many2many:author_books;"`
}

func CreateAuthor(db *gorm.DB, author *Author, book *Book) error {
	author.Books = append(author.Books, *book)
	_ = db.Create(author)
	return nil
}