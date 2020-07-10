package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Code ...
type Code struct {
	gorm.Model
	SendID    uint
	Code      string    `gorm:"type:varchar(50);not null"`
	IsUse     int32     `gorm:"type:tinyint(1);default:0;not null"`
	ExpiresAt time.Time `gorm:"type:datetime;not null"`
	Send      Send
}
