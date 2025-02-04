package repositories

import (
	"github.com/XEFF09/fulldev-be/cogs/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(user *entities.User) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (g *GormUserRepository) CreateUser(user *entities.User) error {
	result := g.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (g *GormUserRepository) GetUserByEmail(user *entities.User) error {
	result := g.db.Where("email = ?", user.Email).First(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
