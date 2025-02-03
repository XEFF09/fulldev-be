package models

import (
	"gorm.io/gorm"
)

type Publisher struct {
	gorm.Model
	Details		string `json:"details"`
	Name		string `json:"name"`
	Books		[]Book
}

func CreatePublisher(db *gorm.DB, pub *Publisher) error {
	_ = db.Create(pub)
	return nil
}