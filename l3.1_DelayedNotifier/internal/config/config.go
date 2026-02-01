package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server         `mapstructure:"server"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type Server struct {
	Address         string        `mapstructure:"address"`
	Timeout         time.Duration `mapstructure:"timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type RabbitMQConfig struct {
	URL   string            `mapstructure:"url"`
	Queue string            `mapstructure:"queue"`
	Retry RabbitRetryConfig `mapstructure:"retry"`
}

type RabbitRetryConfig struct {
	Attempts int     `mapstructure:"attempts"`
	DelayMS  int     `mapstructure:"delay_ms"`
	Backoff  float64 `mapstructure:"backoff"`
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

