package service

import (
	pflag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

var (
	OptionRedisAddress = "redis.address"
)

type Config struct {
	Redis RedisSecretConfig
}

func (this *Config) Init() (err error) {
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("API")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return
}

func (this *Config) BindPFlag(option string, flag *pflag.Flag) (err error) {
	err = viper.BindPFlag(option, flag)
	return err
}

func (this *Config) Unmarshal(configName string) error {
	viper.SetConfigName(configName)

	err := viper.ReadInConfig()
	if err != nil {
		if err, notFound := err.(viper.ConfigFileNotFoundError); !notFound {
			return err
		}
	}

	err = viper.Unmarshal(this)
	if err != nil {
		return err
	}

	return err
}
