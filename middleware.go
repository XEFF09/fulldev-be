package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func authentication(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtZcret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(
		cookie,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
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
	fmt.Println("URL: "+c.OriginalURL(), "METHOD: "+c.Method(), "TIME: "+start.String())

	return c.Next()
}