package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	TokenSynmetricKey string `mapstructure:"TOKEN_SYNMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	MailtrapURL              string        `mapstructure:"MAILTRAP_URL"`
	MailtrapAuthToken        string        `mapstructure:"MAILTRAP_AUTH_TOKEN"`
	SendGridURL              string        `mapstructure:"SENDGRID_URL"`
	SendGridAuthToken        string        `mapstructure:"SENDGRID_AUTH_TOKEN"`
	DefaultFromEmail         string        `mapstructure:"DEFAULT_FROM_EMAIL"`
	EmailSubjectPrefit       string        `mapstructure:"EMAIL_SUBJECT_PREFIX"`
	Provider                 string        `mapstructure:"EMAIL_PROVIDER"`
	VerificationCodeDuration time.Duration `mapstructure:"VERIFICATION_CODE_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
 