package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort    string `mapstructure:"APP_PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
	SOURCE1    string `mapstructure:"SOURCE1"`
	SOURCE2    string `mapstructure:"SOURCE2"`
	SOURCE3    string `mapstructure:"SOURCE3"`
	SOURCE4    string `mapstructure:"SOURCE4"`
}

func LoadConfig() (Config, error) {

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := Config{
		AppPort:    viper.GetString("APP_PORT"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBSSLMode:  viper.GetString("DB_SSLMODE"),
		SOURCE1:    viper.GetString("SOURCE1"),
		SOURCE2:    viper.GetString("SOURCE2"),
		SOURCE3:    viper.GetString("SOURCE3"),
		SOURCE4:    viper.GetString("SOURCE4"),
	}

	log.Printf("%+v\n", cfg)

	return cfg, nil
}
