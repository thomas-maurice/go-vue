package models

import "time"

type OIDCProvider struct {
	Name   string `gorm:"primaryKey;column:name"`
	Active bool   `gorm:"column:active;not null;default:true"`

	DisplayName string `gorm:"column:display_name"`

	Issuer       string `gorm:"column:issuer;not null"`
	ClientID     string `gorm:"column:client_id;not null"`
	ClientSecret string `gorm:"column:client_secret;not null"`
	Scopes       string `gorm:"column:scopes;not null;default:openid,profile,email,groups"`

	Created time.Time `gorm:"created,default:null"`
}

func (o *OIDCProvider) TableName() string {
	return "oidc_providers"
}
