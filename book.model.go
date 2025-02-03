package main

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
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
func getBooks(db *gorm.DB) ([]Book, error) {
	books := []Book{}
	result := db.Find(&books)

	if result.Error != nil {
		return nil, result.Error	
	}
	return books, nil
}

func getBook(db *gorm.DB, id uint) (*Book, error) {
	book := Book{}
	result := db.First(&book, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func createBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func updateBook(db *gorm.DB, book *Book) error {
	result := db.Save(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func deleteBook(db *gorm.DB, id uint) error {
	book := Book{}
	result := db.Delete(&book, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}