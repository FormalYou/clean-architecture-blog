package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"database"`
	JWT struct {
		Secret           string `mapstructure:"secret"`
		ExpiresInMinutes int    `mapstructure:"expires_in_minutes"`
	} `mapstructure:"jwt"`
	Logger struct {
		Level    string `mapstructure:"level"`
		Encoding string `mapstructure:"encoding"`
		File     struct {
			Filename   string `mapstructure:"filename"`
			MaxSize    int    `mapstructure:"max_size"`
			MaxBackups int    `mapstructure:"max_backups"`
			MaxAge     int    `mapstructure:"max_age"`
			Compress   bool   `mapstructure:"compress"`
		} `mapstructure:"file"`
	} `mapstructure:"logger"`
	AuditLog struct {
		File string `mapstructure:"file"`
	} `mapstructure:"audit_log"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}