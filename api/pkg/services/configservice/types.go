package configservice

import "time"

type OIDCProvider struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`

	DisplayName string `json:"display_name"`

	Issuer       string   `json:"issuer"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`

	Created time.Time `json:"created"`
}
