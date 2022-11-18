package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Server ServerCfg `yaml:"server"`
	Db     DBCfg     `yaml:"db"`
}

type DBCfg struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

type ServerCfg struct {
	Addr          string        `yaml:"addr"`
	TokenDuration time.Duration `yaml:"tokenDuration"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
