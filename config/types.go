package config

type Config struct {
	Persistence    PersistenceConfig `json:"persistence"`
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

type PersistenceConfig struct {
	DbName     string `json:"dbName"`
	DbHost     string `json:"dbHost"`
	DbPort     int    `json:"dbPort"`
	DbUsername string `json:"dbUsername"`
	DbPassword string `json:"dbPassword"`
	SslEnabled bool   `json:"sslEnabled"`
}

