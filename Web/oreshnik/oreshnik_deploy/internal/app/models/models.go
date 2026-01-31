package models

import "gorm.io/gorm"

type RevokedToken struct {
	ID    uint   `gorm:"primaryKey"`
	Token string `gorm:"unique"`
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
	IsAdmin  bool `gorm:"default:false"`
}

type Product struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
}

type Purchase struct {
	gorm.Model
	UserID    uint
	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID"`
}

var _ = gorm.Model{}
