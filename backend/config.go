package main

import (
	"os"

	"dario.cat/mergo"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	FeatureFlags struct {
		RegistrationClosed bool `yaml:"registration_closed" envconfig:"FEATURE_REGISTRATION_CLOSED" default:"false"`
	} `yaml:"feature_flags"`
	Database struct {
		Host     string `yaml:"host" envconfig:"DB_HOST" default:"localhost"`
		Port     uint16 `yaml:"port" envconfig:"DB_PORT" default:"5432"`
		User     string `yaml:"user" envconfig:"DB_USER" default:"conference"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD" default:"VeryStrongPassword"`
		Name     string `yaml:"database" envconfig:"DB_NAME" default:"conference"`
	} `yaml:"database"`
	Environment string `yaml:"environment" envconfig:"ENVIRONMENT" default:"local"`
	Port        string `yaml:"port" envconfig:"PORT" default:"8080"`
	Mailer      struct {
		Hostname string `yaml:"hostname" envconfig:"SMTP_HOSTNAME" default:"localhost"`
		Port     string `yaml:"port" envconfig:"SMTP_PORT" default:"1025"`
		From     string `yaml:"from" envconfig:"SMTP_FROM"`
		Password string `yaml:"password" envconfig:"SMTP_PASSWORD"`
	} `yaml:"mailer"`
	BlobUrl string `yaml:"blob_url" envconfig:"BLOB_URL" default:"file:///tmp/"`
	// The default value for these is safe to use for local environment.
	// The code to generate the keys is available here: https://go.dev/play/p/FNe2KGmgc1_f
	Signature struct {
		PublicKey  string `yaml:"public_key" envconfig:"SIGNATURE_PUBLIC_KEY" default:"b0598b81d98ada39a2d2d2d79a855ef9b56444954bdf59edf5979c6ef5a3eca0"`
		PrivateKey string `yaml:"private_key" envconfig:"SIGNATURE_PRIVATE_KEY" default:"82538826d574ba6d85a4c00ba1fc1a202e58397e8f102ff1931d699b6aca1aa3b0598b81d98ada39a2d2d2d79a855ef9b56444954bdf59edf5979c6ef5a3eca0"`
	} `yaml:"signature"`
	EmailTemplate struct {
		TicketPrice                         string `yaml:"ticket_price" envconfig:"EMAIL_TEMPLATE_TICKET_PRICE"`
		TicketStudentCollegePrice           string `yaml:"ticket_student_college_price" envconfig:"EMAIL_TEMPLATE_TICKET_STUDENT_COLLEGE_PRICE"`
		TicketStudentHighSchoolPrice        string `yaml:"ticket_student_high_school_price" envconfig:"EMAIL_TEMPLATE_TICKET_STUDENT_HIGH_SCHOOL_PRICE"`
		TicketStudentCollegeDiscount        string `yaml:"ticket_student_college_discount" envconfig:"EMAIL_TEMPLATE_TICKET_STUDENT_COLLEGE_DISCOUNT"`
		TicketStudentHighSchoolDiscount     string `yaml:"ticket_student_high_school_discount" envconfig:"EMAIL_TEMPLATE_TICKET_STUDENT_HIGH_SCHOOL_DISCOUNT"`
		PercentageStudentCollegeDiscount    string `yaml:"percentage_student_college_discount" envconfig:"EMAIL_TEMPLATE_PERCENTAGE_STUDENT_COLLEGE_DISCOUNT"`
		PercentageStudentHighSchoolDiscount string `yaml:"percentage_student_high_school_discount" envconfig:"EMAIL_TEMPLATE_PERCENTAGE_STUDENT_HIGH_SCHOOL_DISCOUNT"`
		ConferenceEmail                     string `yaml:"conference_email" envconfig:"EMAIL_TEMPLATE_CONFERENCE_EMAIL"`
		BankAccounts                        string `yaml:"bank_accounts" envconfig:"EMAIL_TEMPLATE_BANK_ACCOUNTS"` // List of bank accounts for payments in HTML format
	} `yaml:"email_template"`
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
