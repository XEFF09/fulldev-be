package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_"github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"strconv"

	"github.com/gofiber/swagger"
    _"github.com/XEFF09/fulldev-be/docs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

const (
	host 	 = "postgres"
	port   	 = 5432
	user   	 = "myuser"
	password = "mypassword"
	dbname 	 = "fulldev-postgres-db"
)

var memberUser = User{
	Email: "user@gmail.com",
	Password: "password",
}

func checkMiddleware(c *fiber.Ctx) error {

	cookie := c.Cookies("jwt")
	jwtZcret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(
		cookie,
		jwt.MapClaims{},
		func (token *jwt.Token) (interface{}, error) {
			return []byte(jwtZcret), nil
		},
	)
	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claim := token.Claims.(jwt.MapClaims)
	if claim["role"] != "admin" {
		return c.SendStatus(fiber.StatusUnauthorized)
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

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s "+
  		"password=%s dbname=%s sslmode=disable", 
		host, port, user, password, dbname,
	)
	
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:				time.Second,
			LogLevel:					logger.Info,
			Colorful:					true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to database!")
	}
	fmt.Println("Database connected!")

	err = db.AutoMigrate(&Book{}, &User{})
	if err != nil {
		panic("Failed to migrate database!")
	}
	fmt.Println("Database migrated!")

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, HEAD, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Post("/register", func (c *fiber.Ctx) error {
		user := User{}
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		err = createUser(db, &user)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.JSON(user)
	})
	app.Post("/login", func (c *fiber.Ctx) error {
		req := User{}
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		token, err := loginUser(db, &req)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}

		c.Cookie(&fiber.Cookie{
			Name: "jwt",
			Value: token,
			Expires: time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
		})

		return c.JSON(fiber.Map{
			"message": "Login Successful!",
			"token": token,
		})
	})

	app.Get("/swagger/*", swagger.HandlerDefault)


	app.Use("/books", checkMiddleware)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/books", func (c *fiber.Ctx) error {
		books, err := getBooks(db)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return c.JSON(books)
	})
	app.Get("/books/:id", func (c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		book, err := getBook(db, uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return c.JSON(book)
	})
	app.Post("/books", func (c *fiber.Ctx) error {
		book := Book{}
		err := c.BodyParser(&book)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		err = createBook(db, &book)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(book)
	}) 
	app.Put("/books/:id", func (c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))	
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		req := Book{}
		err = c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		book, err := getBook(db, uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}

		book.Title = req.Title
		book.Author = req.Author

		err = updateBook(db, book)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(book)
	})
	app.Delete(("/books/:id"), func (c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))	
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		err = deleteBook(db, uint(id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(fiber.Map{
			"message": "Deleted Successful!",
		})
	})
	app.Post("/upload", uploadFile)

	app.Listen(":8088")
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