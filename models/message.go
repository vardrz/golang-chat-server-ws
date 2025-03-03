package models

import (
	"gorm.io/gorm"
)

type Message struct {
    gorm.Model
    FromID    uint   `json:"from_id" gorm:"not null"`  // Add not null constraint
    ToID      uint   `json:"to_id" gorm:"not null"`    // Add not null constraint
    Content   string `json:"content" gorm:"type:text"`
    IsRead    bool   `json:"is_read" gorm:"default:false"`
    From      User   `gorm:"foreignKey:FromID"`
    To        User   `gorm:"foreignKey:ToID"`
}

func (Message) TableName() string {
	return "messages"
}