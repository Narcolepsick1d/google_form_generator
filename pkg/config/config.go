package config

import (
	"github.com/spf13/viper"
	"log"
	"strconv"
)

type Config struct {
	TelegramToken string
	Host          string
	Port          int
	Username_DB   string
	DBName        string
	SSLMode       string
	Password      string
}

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := fromEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	return nil
}

func fromEnv(cfg *Config) error {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("нету значений")
	}

	cfg.TelegramToken = viper.GetString("TelegramToken")
	cfg.Host = viper.GetString("HOST")
	cfg.Port, err = strconv.Atoi(viper.GetString("PORT"))
	cfg.Username_DB = viper.GetString("USERNAMEDB")
	cfg.DBName = viper.GetString("DBNAME")
	cfg.SSLMode = viper.GetString("SSLMODE")
	cfg.Password = viper.GetString("PASSWORD")

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("app")

	return viper.ReadInConfig()
}
