package main

import (
	_"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Book struct {
	ID		uuid.UUID `json:"id"`
	Title	string `json:"title"` 
	Author	string `json:"author"`
}
func NewBook(title string, author string) Book {
	return Book{
		ID: uuid.New(),
		Title: title,
		Author: author,
	}
}

type BookRequest struct {
	Title	string `json:"title"`
	Author	string `json:"author"`
}

var books []Book

func main() {
	app := fiber.New()

	books = append(
		books,
		NewBook("The Hitchhiker's Guide to the Galaxy", "Douglas Adams"),
		NewBook("Nineteen Eighty-Four", "George Orwell"),
		NewBook("Brave New World", "Aldous Huxley"),
	)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete(("/books/:id"), deleteBook)

	app.Listen(":8088")
}

func getBooks(c *fiber.Ctx) error {
	return c.JSON(books)
}

func getBook(c *fiber.Ctx) error {
	bookId := c.Params("id")
	for _, book := range books {
		if book.ID.String() == bookId {
			return c.JSON(book)
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}

func createBook(c *fiber.Ctx) error {
	req := new(BookRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	newBook := NewBook(req.Title, req.Author)
	books = append(books, newBook)
	return c.JSON(newBook)	
}

func updateBook(c *fiber.Ctx) error {
	bookId := c.Params("id")
	req := new(BookRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books {
		if book.ID.String() == bookId {
			book.Title = req.Title
			book.Author = req.Author
			books[i] = book
			return c.JSON(book)
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}

func deleteBook(c *fiber.Ctx) error {
	bookId := c.Params("id")
	for i, book := range books {
		if book.ID.String() == bookId {
			books = append(books[:i], books[i+1:]...)
			return c.SendString("Book deleted")
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}