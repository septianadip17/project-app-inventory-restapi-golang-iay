package config

import "github.com/spf13/viper"

type Config struct {
	DBUrl      string `mapstructure:"DB_URL"`
	ServerPort string `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	err := viper.Unmarshal(&cfg)
	return &cfg, err
}
