package util

import "github.com/spf13/viper"

type Config struct {
	Server Server `yaml:"server"`
	Db     DB     `yaml:"db"`
}

type DB struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

type Server struct {
	Addr string `yaml:"addr"`
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
