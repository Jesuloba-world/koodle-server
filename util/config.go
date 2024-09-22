package util

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	PGUser      string `mapstructure:"PGUSER"`
	PGPassword  string `mapstructure:"PGPASSWORD"`
	PGHost      string `mapstructure:"PGHOST"`
	PGPort      string `mapstructure:"PGPORT"`
	PGDB        string `mapstructure:"PGDATABASE"`
	SecretKey   string `mapstructure:"SECRET_KEY"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
	BrevoKey    string `mapstructure:"BREVO_API_KEY"`
	MsKey       string `mapstructure:"MAILERSEND_API_KEY"`
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

		viper.BindEnv("PGUSER")
		viper.BindEnv("PGPASSWORD")
		viper.BindEnv("PGHOST")
		viper.BindEnv("PGPORT")
		viper.BindEnv("PGDATABASE")
		viper.BindEnv("SECRET_KEY")
		viper.BindEnv("BREVO_API_KEY")
		viper.BindEnv("MAILERSEND_API_KEY")

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
