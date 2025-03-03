package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nickname   string `gorm:"size:255;not null;unique" json:"nickname"`
	First_name string `gorm:"size:255" json:"first_name"`
	Last_name  string `gorm:"size:255" json:"last_name"`
	Email      string `gorm:"size:100;unique" json:"email"`
	Password   string `gorm:"size:100;not null;" json:"password"`
	Phone      string `gorm:"size:255;" json:"phone"`
	Pic        string `gorm:"size:255" json:"pic"`
	Otoritas   uint32 `gorm:"size:11;" json:"otoritas"`
	Status     string `json:"status" gorm:"column:status"`
}
func (User) TableName() string {
	return "user"
}