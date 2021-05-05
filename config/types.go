package config

type Config struct {
	FeatureToggles FeatureToggles `json:"featuretoggles"`
}

type FeatureToggles map[string]bool
