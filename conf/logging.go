package conf

import (
	"github.com/sirupsen/logrus"
)

// LoggingConf is the configuration model of the logger
type loggingConfig struct {
	Level            string                 `mapstructure:"log_level" json:"log_level"`
	File             string                 `mapstructure:"log_file" json:"log_file"`
	DisableColors    bool                   `mapstructure:"disable_colors" split_words:"true" json:"disable_colors"`
	QuoteEmptyFields bool                   `mapstructure:"quote_empty_fields" split_words:"true" json:"quote_empty_fields"`
	TSFormat         string                 `mapstructure:"ts_format" json:"ts_format"`
	Fields           map[string]interface{} `mapstructure:"fields" json:"fields"`
}

// ConfigureLogging sets the global logging configuration accepts LoggingConfig
func ConfigureLogging(config *loggingConfig) (*logrus.Entry, error) {

}
