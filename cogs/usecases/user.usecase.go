package usecases

import (
	"github.com/XEFF09/fulldev-be/cogs/adapters/repositories"
	"github.com/XEFF09/fulldev-be/cogs/entities"
	_"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"os"
)

type UserUseCase interface {
	RegisterUser(uesr *entities.User) error
	LoginUser(user *entities.User) (string, error)
}

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserUseCase {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(user *entities.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)

	err = s.repo.CreateUser(user)
	if err != nil {
		return err
	}

	return  nil
}
func (s *UserService) LoginUser(req *entities.User) (string, error) {
	user := entities.User{Email: req.Email}
	err := s.repo.GetUserByEmail(&user)
	if err != nil {
		return "", err
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