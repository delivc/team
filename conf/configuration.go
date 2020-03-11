package conf

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// EmailProviderConfiguration holds email related configs
type EmailProviderConfiguration struct {
	Disabled bool `json:"disabled"`
}

// DBConfiguration holds all the database related configuration.
type DBConfiguration struct {
	Driver         string `json:"driver" required:"true"`
	URL            string `json:"url" envconfig:"DELIVC_DATABASE_URL" required:"true"`
	Namespace      string `json:"namespace"`
	MigrationsPath string `json:"migrations_path" split_words:"true" default:"./migrations"`
}

// SMTPConfiguration is the SMTP config for the Mailer
type SMTPConfiguration struct {
	MaxFrequency time.Duration `json:"max_frequency" split_words:"true"`
	Host         string        `json:"host"`
	Port         int           `json:"port,omitempty" default:"587"`
	User         string        `json:"user"`
	Pass         string        `json:"pass,omitempty"`
	AdminEmail   string        `json:"admin_email" split_words:"true"`
}

// GlobalConfiguration holds all the configuration that applies to all instances.
type GlobalConfiguration struct {
	API struct {
		Host            string
		Port            int `envconfig:"PORT" default:"8083"`
		Endpoint        string
		RequestIDHeader string `envconfig:"REQUEST_ID_HEADER"`
	}
	IdentityEndpoint string        `envconfig:"DELIVC_IDENTITY_ENDPOINT" required:"true"`
	Logging          loggingConfig `envconfig:"LOG"`
	OperatorToken    string        `split_words:"true" required:"true"`
	DB               DBConfiguration
	SMTP             SMTPConfiguration
}

func loadEnvironment(filename string) error {
	var err error

	if filename != "" {
		err = godotenv.Load(filename)
	} else {
		err = godotenv.Load()
		// handle if .env file does not exist, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

// LoadGlobal loads configuration from file and environment variables.
func LoadGlobal(filename string) (*GlobalConfiguration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(GlobalConfiguration)
	if err := envconfig.Process("delivc", config); err != nil {
		return nil, err
	}
	if _, err := ConfigureLogging(&config.Logging); err != nil {
		return nil, err
	}

	if config.SMTP.MaxFrequency == 0 {
		config.SMTP.MaxFrequency = 15 * time.Minute
	}
	return config, nil
}

// EmailContentConfiguration holds the configuration for emails, both subjects and template URLs.
type EmailContentConfiguration struct {
	Invite       string `json:"invite"`
	Confirmation string `json:"confirmation"`
}

// Configuration holds all the per-instance configuration.
type Configuration struct {
	SiteURL string            `json:"site_url" split_words:"true" required:"true"`
	SMTP    SMTPConfiguration `json:"smtp"`
	Mailer  struct {
		Subjects  EmailContentConfiguration `json:"subjects"`
		Templates EmailContentConfiguration `json:"templates"`
		URLPaths  EmailContentConfiguration `json:"url_paths"`
	} `json:"mailer"`
}

// LoadConfig loads per-instance configuration.
func LoadConfig(filename string) (*Configuration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(Configuration)
	if err := envconfig.Process("delivc", config); err != nil {
		return nil, err
	}
	config.ApplyDefaults()
	return config, nil
}

// ApplyDefaults sets defaults for a Configuration
func (config *Configuration) ApplyDefaults() {

	if config.Mailer.URLPaths.Invite == "" {
		config.Mailer.URLPaths.Invite = "/"
	}
	if config.Mailer.URLPaths.Confirmation == "" {
		config.Mailer.URLPaths.Confirmation = "/"
	}

	if config.SMTP.MaxFrequency == 0 {
		config.SMTP.MaxFrequency = 15 * time.Minute
	}
}

// Value returns the configuration as string
func (config *Configuration) Value() (driver.Value, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return driver.Value(""), err
	}
	return driver.Value(string(data)), nil
}

// Scan unmarshals given interface
func (config *Configuration) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	default:
		return errors.New("Invalid data type for Configuration")
	}

	if len(source) == 0 {
		source = []byte("{}")
	}
	return json.Unmarshal(source, &config)
}
