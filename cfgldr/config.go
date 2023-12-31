package cfgldr

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSL      string `mapstructure:"ssl"`
}

type S3 struct {
	BucketName string `mapstructure:"bucket_name"`
	Region     string `mapstructure:"region"`
}

type App struct {
	Port  int  `mapstructure:"port"`
	Debug bool `mapstructure:"debug"`
}

type Config struct {
	App      App      `mapstructure:"app"`
	Database Database `mapstructure:"database"`
	S3       S3       `mapstructure:"s3"`
}

func LoadConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error occurs while reading the config")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.Wrap(err, "error occurs while unmarshal the config")
	}

	return
}
