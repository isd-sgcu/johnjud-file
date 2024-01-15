package cfgldr

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Database struct {
	Url string `mapstructure:"db_url"`
}

type Bucket struct {
	Endpoint        string `mapstructure:"bucket_endpoint"`
	AccessKeyID     string `mapstructure:"bucket_access_key"`
	SecretAccessKey string `mapstructure:"bucket_secret_key"`
	UseSSL          bool   `mapstructure:"bucket_use_ssl"`
	BucketName      string `mapstructure:"bucket_name"`
}

type App struct {
	Port int    `mapstructure:"app_port"`
	Env  string `mapstructure:"app_env"`
}

type Config struct {
	App      App
	Database Database
	Bucket   Bucket
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "production" {
		viper.AutomaticEnv()
	}
	if env == "development" {
		viper.SetConfigFile(".env")
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal().Err(err).
				Str("service", "file").
				Msg("Failed to load .env file")
		}
	}

	var dbConfig Database
	if err := viper.Unmarshal(&dbConfig); err != nil {
		return nil, err
	}

	var appConfig App
	if err := viper.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	var bucketConfig Bucket
	if err := viper.Unmarshal(&bucketConfig); err != nil {
		return nil, err
	}

	config := &Config{
		Database: dbConfig,
		App:      appConfig,
		Bucket:   bucketConfig,
	}

	log.Info().Interface("config", config).Msg("Loaded config")

	return config, nil
}

func (ac *App) IsDevelopment() bool {
	return ac.Env == "development"
}
