package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/gofiber/swagger"
    _"github.com/XEFF09/fulldev-be/docs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host 	 = "postgres"
	port   	 = 5432
	user   	 = "myuser"
	password = "mypassword"
	dbname 	 = "fulldev-postgres-db"
)

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

type User struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

var books []Book

var memberUser = User{
	Email: "user@gmail.com",
	Password: "password",
}

func checkMiddleware(c *fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized 
	}
	
	start := time.Now()
	fmt.Println("URL: " + c.OriginalURL(), "METHOD: " + c.Method(), "TIME: " + start.String())

	return c.Next()
}

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8088
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	fmt.Println(db, "Database connected25!")
	db.AutoMigrate(&Book{})

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	fmt.Println("hello woroldjkdajwlkdja")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, HEAD, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	books = append(
		books,
		NewBook("The Hitchhiker's Guide to the Galaxy", "Douglas Adams"),
		NewBook("Nineteen Eighty-Four", "George Orwell"),
		NewBook("Brave New World", "Aldous Huxley"),
	)

	app.Post("/login", login)
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	app.Use(checkMiddleware)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete(("/books/:id"), deleteBook)
	app.Post("/upload", uploadFile)
	app.Get("/env", getEnv)

	app.Listen(":8088")
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

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = c.SaveFile(file, "./uploads/" + file.Filename)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File uploaded")
}

func getEnv(c *fiber.Ctx) error {

	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "default-secret"
	}

	return c.JSON(fiber.Map{
		"env": secret,
	})
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email != memberUser.Email || user.Password != memberUser.Password {
		return fiber.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Login success",
		"token": t,
	})
}