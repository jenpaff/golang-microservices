package config

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"github.com/imdario/mergo"
	"github.com/jenpaff/golang-microservices/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Build creates and initializes the configuration
func BuildConfig(envConfigFilePath string, secretsDirectoryPath string, injectedSecrets string) (Config, error) {
	config := defaultConfig

	envConfig, err := readEnvSpecificConfig(envConfigFilePath)
	if err != nil {
		return config, err
	}

	err = mergo.Merge(&config, envConfig, mergo.WithOverride)
	if err != nil {
		return config, fmt.Errorf("could not merge configs: %s: %w", err.Error(), errors.InternalServerError)
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return config, fmt.Errorf("could not marshal config: %s: %w", err.Error(), errors.InternalServerError)
	}

	secretsTemplate, err := template.New("config").Parse(string(configJSON))
	if err != nil {
		return config, fmt.Errorf("could not create secrets template: %s: %w", err.Error(), errors.InternalServerError)
	}

	secrets := make(map[string]string)

	secrets, err = injectSecretsFromSecretsDir(secretsDirectoryPath, secrets)
	if err != nil {
		return config, err
	}

	secrets, err = injectSecretsFromEnvironment(injectedSecrets, secrets)
	if err != nil {
		return config, err
	}

	var tpl bytes.Buffer
	err = secretsTemplate.Execute(&tpl, secrets)
	if err != nil {
		return config, fmt.Errorf("could not insert secrets into template: %s: %w", err.Error(), errors.InternalServerError)
	}

	err = json.Unmarshal(tpl.Bytes(), &config)
	if err != nil {
		return config, fmt.Errorf("could not unmarshal config: %s: %w", err.Error(), errors.InternalServerError)
	}

	return config, nil
}

func injectSecretsFromEnvironment(injectedSecrets string, secrets map[string]string) (map[string]string, error) {
	if injectedSecrets != "" {
		for _, splitSecret := range strings.Split(injectedSecrets, ";") {
			secret := strings.Split(splitSecret, ":")

			secretName := strings.TrimSpace(secret[0])

			unencodedSecret, err := base64.StdEncoding.DecodeString(secret[1])
			if err != nil {
				return nil, fmt.Errorf("could not decode injected secret %s: %s: %w", secretName, err.Error(), errors.InternalServerError)
			}

			secretValue := strings.TrimSpace(string(unencodedSecret))

			secrets[secretName] = secretValue
		}
	}
	return secrets, nil
}

func injectSecretsFromSecretsDir(secretsDirectoryPath string, secrets map[string]string) (map[string]string, error) {
	secretsDir, err := ioutil.ReadDir(secretsDirectoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Infof("Secrets directory %s not found. Skipping.", secretsDirectoryPath)
		} else {
			return nil, fmt.Errorf("could not read secrets dir %s: %s: %w", secretsDirectoryPath, err.Error(), errors.InternalServerError)
		}
	}

	for _, fileInfo := range secretsDir {
		if fileInfo.Name()[0:1] == "." || fileInfo.IsDir() {
			// ignore "hidden" files and dirs
			continue
		}

		secretPath := filepath.Join(secretsDirectoryPath, fileInfo.Name())
		fileData, err := ioutil.ReadFile(secretPath)
		if err != nil {
			return nil, fmt.Errorf("could not read secret fileInfo %s: %s: %w", secretPath, err.Error(), errors.InternalServerError)
		}

		secrets[fileInfo.Name()] = strings.TrimSpace(string(fileData))
	}
	return secrets, nil
}

func readEnvSpecificConfig(envConfigFilePath string) (Config, error) {
	fileBytes, err := ioutil.ReadFile(envConfigFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("could not read config file %s: %s: %w", envConfigFilePath, err.Error(), errors.InternalServerError)
	}

	envConfig := Config{}
	err = json.Unmarshal(fileBytes, &envConfig)
	if err != nil {
		return Config{}, fmt.Errorf("could not unmarshal secret file %s: %s: %w", envConfigFilePath, err.Error(), errors.InternalServerError)
	}
	return envConfig, nil
}
