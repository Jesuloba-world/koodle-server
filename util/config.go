package util

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"

)

type Config struct {
	// PGUser      string `mapstructure:"PGUSER"`
	// PGPassword  string `mapstructure:"PGPASSWORD"`
	// PGHost      string `mapstructure:"PGHOST"`
	// PGPort      string `mapstructure:"PGPORT"`
	// PGDB        string `mapstructure:"PGDATABASE"`
	SecretKey   string `mapstructure:"SECRET_KEY"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`

	// SMTP Configuration
	SMTPHost    string `mapstructure:"SMTP_HOST"`
	SMTPPort    int    `mapstructure:"SMTP_PORT"`
	SMTPUser    string `mapstructure:"SMTP_USER"`
	SMTPPass    string `mapstructure:"SMTP_PASS"`
	EmailSender string `mapstructure:"EMAIL_SENDER"`
}

var (
	config Config
	once   sync.Once
	loaded bool
)

func loadConfig() (Config, error) {
	var err error

	once.Do(func() {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		if err = viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("Config file not found, relying on environment variables")
			} else {
				return
			}
		}

		// viper.BindEnv("PGUSER")
		// viper.BindEnv("PGPASSWORD")
		// viper.BindEnv("PGHOST")
		// viper.BindEnv("PGPORT")
		// viper.BindEnv("PGDATABASE")
		viper.BindEnv("SECRET_KEY")
		viper.BindEnv("BREVO_API_KEY")
		viper.BindEnv("SMTP_HOST")
		viper.BindEnv("SMTP_PORT")
		viper.BindEnv("SMTP_USER")
		viper.BindEnv("SMTP_PASS")
		viper.BindEnv("EMAIL_SENDER")

		err = viper.Unmarshal(&config)
		if err == nil {
			loaded = true
		}
	})

	return config, err
}

func GetConfig() (Config, error) {
	if !loaded {
		return loadConfig()
	}
	return config, nil
}
