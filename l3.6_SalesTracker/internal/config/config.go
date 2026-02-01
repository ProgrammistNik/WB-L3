package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server Server `mapstructure:"server"`
	DB     DB     `mapstructure:"db"`
}

type Server struct {
	Address         string        `mapstructure:"address"`
	Timeout         time.Duration `mapstructure:"timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
}

type ConfigLoader struct {
	v *viper.Viper
}

func New() *ConfigLoader {
	v := viper.New()
	return &ConfigLoader{v: v}
}

func (c *ConfigLoader) Load(path string) error {
	c.v.SetConfigFile(path)
	return c.v.ReadInConfig()
}

func (c *ConfigLoader) Unmarshal(rawVal any) error {
	return c.v.Unmarshal(rawVal)
}
