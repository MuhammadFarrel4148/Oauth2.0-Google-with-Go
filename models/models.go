package models

import "gorm.io/gorm"

type Input_User struct {
	gorm.Model
	Email	string		`json:"unique;not null"`
}