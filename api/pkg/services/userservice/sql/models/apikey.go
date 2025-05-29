package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APIKey struct {
	Id      string    `gorm:"primaryKey;column:id"`
	Hash    string    `gorm:"column:hash;not null"`
	Name    string    `gorm:"column:string;not null"`
	Active  bool      `gorm:"column:active;not null;default:true"`
	UserId  string    `gorm:"column:user_id;not null"`
	User    User      `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
	Expires time.Time `gorm:"column:expires;default:null"`
}

func (o *APIKey) TableName() string {
	return "api_keys"
}

func (o *APIKey) AfterFind(tx *gorm.DB) error {
	return tx.First(&o.User, &User{Id: o.UserId}).Error
}

func (o *APIKey) BeforeCreate(tx *gorm.DB) (err error) {
	if o.Id == "" {
		o.Id = uuid.NewString()
	}

	return nil
}
