package models

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"os"
)

type User struct {
	gorm.Model
	Email string `gorm:"unique"`
	Password string 
}

func CreateUser(db *gorm.DB, user *User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)

	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return  nil
}

func LoginUser(db *gorm.DB, req *User) (string, error) {
	user := User{}
	result := db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return "", result.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	jwtZcret := os.Getenv("JWT_SECRET")

	t, err := token.SignedString([]byte(jwtZcret))
	if err != nil {
		return "", err
	}

	return t, nil
}