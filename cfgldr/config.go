package cfgldr

import (
	"github.com/spf13/viper"
)

type Database struct {
	Url string `mapstructure:"URL"`
}

type Bucket struct {
	Endpoint        string `mapstructure:"ENDPOINT"`
	AccessKeyID     string `mapstructure:"ACCESS_KEY"`
	SecretAccessKey string `mapstructure:"SECRET_KEY"`
	UseSSL          bool   `mapstructure:"USE_SSL"`
	BucketName      string `mapstructure:"NAME"`
}

type App struct {
	Port int    `mapstructure:"PORT"`
	Env  string `mapstructure:"ENV"`
}

type Config struct {
	App      App
	Database Database
	Bucket   Bucket
}

func LoadConfig() (*Config, error) {
	dbCfgLdr := viper.New()
	dbCfgLdr.SetEnvPrefix("DB")
	dbCfgLdr.AutomaticEnv()
	dbCfgLdr.AllowEmptyEnv(false)
	dbConfig := Database{}
	if err := dbCfgLdr.Unmarshal(&dbConfig); err != nil {
		return nil, err
	}

	appCfgLdr := viper.New()
	appCfgLdr.SetEnvPrefix("APP")
	appCfgLdr.AutomaticEnv()
	dbCfgLdr.AllowEmptyEnv(false)
	appConfig := App{}
	if err := appCfgLdr.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	bucketCfgLdr := viper.New()
	bucketCfgLdr.SetEnvPrefix("BUCKET")
	bucketCfgLdr.AutomaticEnv()
	bucketCfgLdr.AllowEmptyEnv(false)
	bucketConfig := Bucket{}
	if err := bucketCfgLdr.Unmarshal(&bucketConfig); err != nil {
		return nil, err
	}

	return &Config{
		Database: dbConfig,
		App:      appConfig,
		Bucket:   bucketConfig,
	}, nil
}

func (ac *App) IsDevelopment() bool {
	return ac.Env == "development"
}
