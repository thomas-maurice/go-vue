package configservice

type ConfigService interface {
	GetOIDCProvider(name string) (*OIDCProvider, error)
	UpsertOIDCProvider(prov *OIDCProvider) (*OIDCProvider, error)
	GetOIDCProviders() ([]OIDCProvider, error)
	CreateOIDCProvider(prov *OIDCProvider) (*OIDCProvider, error)
}
