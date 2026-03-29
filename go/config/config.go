package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"database"`
	Redis  RedisConfig  `mapstructure:"redis"`
}

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`  // seconds
	WriteTimeout int    `mapstructure:"write_timeout"` // seconds
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
