package rest

import (
	"github.com/XEFF09/fulldev-be/cogs/usecases"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/XEFF09/fulldev-be/cogs/entities"
)

type UserRest struct {
	userUsecase usecases.UserUseCase
}

func NewUserRest(userUsecase usecases.UserUseCase) *UserRest {
	return &UserRest{userUsecase: userUsecase}
}

func (ur *UserRest) Register(c *fiber.Ctx) error {
	user := entities.User{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := ur.userUsecase.RegisterUser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(user)
}

func (ur *UserRest) Login(c *fiber.Ctx) error {
	user := entities.User{}
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	token, err := ur.userUsecase.LoginUser(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"message": "Login Successful!",
		"token":   token,
	})
}

