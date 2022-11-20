package util

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	Server ServerCfg `yaml:"server"`
	Db     DBCfg     `yaml:"db"`
}

type DBCfg struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Source   string `yaml:"source"`
}

type ServerCfg struct {
	Addr          string        `yaml:"addr"`
	TokenDuration time.Duration `yaml:"tokenDuration"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	//for _, key := range viper.AllKeys() {
	//	val := viper.Get(key)
	//	if val != "" {
	//		viper.Set(key, val)
	//	}
	//}

	config.Db.Source = fmt.Sprintf(
		config.Db.Source,
		config.Db.Username,
		config.Db.Password,
		config.Db.Host,
		config.Db.Port,
	)
	return
}

func replaceEnvInConfig(body []byte) []byte {
	search := regexp.MustCompile(`\$\{([^{}]+)\}`)
	replacedBody := search.ReplaceAllFunc(body, func(b []byte) []byte {
		group1 := search.ReplaceAllString(string(b), `$1`)

		envValue := os.Getenv(group1)
		if len(envValue) > 0 {
			return []byte(envValue)
		}
		return []byte("")
	})
	return replacedBody
}
