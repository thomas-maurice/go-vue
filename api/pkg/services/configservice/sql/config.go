package configservice

import (
	"crypto/ecdsa"
	"errors"
	"strings"

	"github.com/thomas-maurice/api/go-vue/pkg/services/configservice"
	"github.com/thomas-maurice/api/go-vue/pkg/services/configservice/sql/models"
	"gorm.io/gorm"
)

type ConfigService struct {
	DB         *gorm.DB
	SigningKey *ecdsa.PrivateKey
}

func oidcProviderFromModel(input *models.OIDCProvider) *configservice.OIDCProvider {
	return &configservice.OIDCProvider{
		Name:        input.Name,
		Active:      input.Active,
		DisplayName: input.DisplayName,

		Issuer:       input.Issuer,
		ClientID:     input.ClientID,
		ClientSecret: input.ClientSecret,
		Scopes:       strings.Split(input.Scopes, ","),

		Created: input.Created,
	}
}

func oidcProviderToModel(input *configservice.OIDCProvider) *models.OIDCProvider {
	return &models.OIDCProvider{
		Name:        input.Name,
		Active:      input.Active,
		DisplayName: input.DisplayName,

		Issuer:       input.Issuer,
		ClientID:     input.ClientID,
		ClientSecret: input.ClientSecret,
		Scopes:       strings.Join(input.Scopes, ","),

		Created: input.Created,
	}
}

func NewConfigService(db *gorm.DB) (configservice.ConfigService, error) {
	if err := db.AutoMigrate(models.OIDCProvider{}); err != nil {
		return nil, err
	}

	return &ConfigService{
		DB: db,
	}, nil
}

func (s *ConfigService) GetOIDCProvider(name string) (*configservice.OIDCProvider, error) {
	var prov models.OIDCProvider
	if err := s.DB.Where(&models.OIDCProvider{Name: name}).First(&prov).Error; err != nil {
		return nil, err
	}

	return oidcProviderFromModel(&prov), nil
}

func (s *ConfigService) UpsertOIDCProvider(prov *configservice.OIDCProvider) (*configservice.OIDCProvider, error) {
	var prv models.OIDCProvider
	if err := s.DB.Where(&models.OIDCProvider{Name: prov.Name}).First(&prv).Error; err == nil {
		if err := s.DB.Save(oidcProviderToModel(prov)).Error; err != nil {
			return nil, err
		}
		return prov, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := s.DB.Create(oidcProviderToModel(prov)).Error; err != nil {
			return prov, err
		}
	}

	return oidcProviderFromModel(&prv), nil
}

func (s *ConfigService) GetOIDCProviders() ([]configservice.OIDCProvider, error) {
	var providers []models.OIDCProvider
	if err := s.DB.Find(&providers).Error; err != nil {
		return nil, err
	}

	plist := make([]configservice.OIDCProvider, 0)
	for _, p := range providers {
		plist = append(plist, *oidcProviderFromModel(&p))
	}

	return plist, nil
}

func (s *ConfigService) CreateOIDCProvider(prov *configservice.OIDCProvider) (*configservice.OIDCProvider, error) {
	if err := s.DB.Create(oidcProviderToModel(prov)).Error; err != nil {
		return nil, err
	}

	return prov, nil
}
