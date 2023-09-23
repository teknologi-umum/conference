package main

import (
	"os"

	"dario.cat/mergo"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host" envconfig:"DB_HOST" default:"localhost"`
		Port     uint16 `yaml:"port" envconfig:"DB_PORT" default:"5432"`
		User     string `yaml:"user" envconfig:"DB_USER" default:"postgres"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD" default:"postgres"`
		Name     string `yaml:"database" envconfig:"DB_NAME" default:"postgres"`
	} `yaml:"database"`
	Environment string `yaml:"environment" envconfig:"ENVIRONMENT" default:"production"`
	Port        string `yaml:"port" envconfig:"PORT" default:"8080"`
	Mailer      struct {
		Hostname string `yaml:"hostname" envconfig:"SMTP_HOSTNAME"`
		Port     string `yaml:"port" envconfig:"SMTP_PORT"`
		From     string `yaml:"from" envconfig:"SMTP_FROM"`
		Password string `yaml:"password" envconfig:"SMTP_PASSWORD"`
	} `yaml:"mailer"`
	BlobUrl   string `yaml:"blob_url" envconfig:"BLOB_URL" default:"file:///data/"`
	Signature struct {
		PublicKey  string `yaml:"public_key" envconfig:"SIGNATURE_PUBLIC_KEY"`
		PrivateKey string `yaml:"private_key" envconfig:"SIGNATURE_PRIVATE_KEY"`
	} `yaml:"signature"`
}

func GetConfig(configurationFile string) (Config, error) {
	var configurationFromEnvironment Config
	err := envconfig.Process("", &configurationFromEnvironment)
	if err != nil {
		return Config{}, err
	}

	var configurationFromYaml Config
	if configurationFile != "" {
		f, err := os.Open(configurationFile)
		if err != nil {
			return Config{}, err
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Error().Err(err).Msg("closing configuration file")
			}
		}()
		err = yaml.NewDecoder(f).Decode(&configurationFromYaml)
		if err != nil {
			return Config{}, err
		}
	}

	// Environment variables set the precedence
	err = mergo.Merge(&configurationFromYaml, configurationFromEnvironment)
	if err != nil {
		return Config{}, err
	}

	return configurationFromYaml, nil
}
