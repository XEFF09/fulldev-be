package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title  			string `json:"title"`
	Author 			string `json:"author"`
	PublisherID 	uint `json:"publisher_id"`
	Publisher	 	Publisher
	Authors 		[]Author `gorm:"many2many:author_books;"`
}

// Handler functions
// getBooks godoc
// @Summary Get all books
// @Description Get details of all books
// @Tags books
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} Book
// @Router /books [get]
func GetBooks(db *gorm.DB) ([]Book, error) {
	books := []Book{}
	result := db.Preload("Authors").Preload("Publisher").Find(&books)

	if result.Error != nil {
		return nil, result.Error	
	}
	return books, nil
}

func GetBook(db *gorm.DB, id uint) (*Book, error) {
	book := Book{}
	result := db.Preload("Authors").Preload("Publisher").First(&book, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func CreateBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateBook(db *gorm.DB, book *Book) error {
	result := db.Save(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteBook(db *gorm.DB, id uint) error {
	book := Book{}
	result := db.Delete(&book, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}