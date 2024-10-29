package configs

import (
	"github.com/spf13/viper"
	"log"
)

var EnvConfigs *envConfig

type envConfig struct {
	ClientID      string `mapstructure:"CLIENT_ID"`
	DbName        string `mapstructure:"DB_NAME"`
	DbPassword    string `mapstructure:"DB_PASSWORD"`
	DbPort        string `mapstructure:"DB_PORT"`
	DbUser        string `mapstructure:"DB_USER"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	SessionSecret string `mapstructure:"SESSION_SECRET"`
}

func InitEnvConfig() {
	EnvConfigs = loadEnv()
}

func loadEnv() (config *envConfig) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
	}

	return
}
