package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	Id      string    `gorm:"primaryKey;column:id"`
	UserId  string    `gorm:"column:user_id;not null"`
	User    User      `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
	Expires time.Time `gorm:"column:expires"`
}

func (o *Session) AfterFind(tx *gorm.DB) error {
	return tx.First(&o.User, &User{Id: o.UserId}).Error
}

func (o *Session) BeforeCreate(tx *gorm.DB) (err error) {
	if o.Id == "" {
		o.Id = uuid.NewString()
	}

	return nil
}
