package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/thomas-maurice/api/go-vue/pkg/services/userservice"
	"gorm.io/gorm"
)

type User struct {
	Id          string               `gorm:"column:id;primary_key"`
	Username    string               `gorm:"column:username;unique;not null"`
	Email       string               `gorm:"column:email;unique;not null"`
	Active      bool                 `gorm:"column:active;not null;default:true"`
	DisplayName string               `gorm:"column:display_name"`
	Kind        userservice.UserKind `gorm:"column:kind;not null"`
	Admin       bool                 `gorm:"column:admin;default:false"`
	Password    string               `gorm:"column:password"`
	Created     time.Time            `gorm:"column:created;default:null"`
	LastLogin   time.Time            `gorm:"column:last_login;default:null"`
}

func (o *User) TableName() string {
	return "users"
}

func (o *User) BeforeCreate(tx *gorm.DB) (err error) {
	if o.Id == "" {
		o.Id = uuid.NewString()
	}

	return nil
}
