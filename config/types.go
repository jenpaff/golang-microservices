package config

type Config struct {
	FeatureToggles FeatureToggles `json:"featuretoggles"`
}

type FeatureToggles map[string]bool

type Postgres struct {
	Host       string
	Port       int
	UserName   string
	Password   string
	DBName     string
	SSLEnabled bool
}
