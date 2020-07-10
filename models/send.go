package models

import (
	"github.com/jinzhu/gorm"
)

// Send ...
type Send struct {
	gorm.Model
	TemplateID uint   `gorm:"index"`
	UserID     uint   `gorm:"index"`
	SID        string `gorm:"column:sid;type:varchar(255)"`
	Provider   string `gorm:"type:varchar(100);not null"`
	BizType    int32  `gorm:"type:tinyint(1);default:0;index;not null"`
	Mobile     string `gorm:"type:varchar(50);not null"`
	Content    string `gorm:"type:text;not null"`
	Success    int32  `gorm:"type:tinyint(1);default:0;not null"`
	Message    string `gorm:"type:varchar(500)"`
	Code       *Code
}
