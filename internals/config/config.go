package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lenarsaitov/go-task/pkg/logging"
	"sync"
)

type Config struct {
	Listen struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	Postgres struct {
		Host     string `yaml:"host" env-default:"localhost"`
		Port     string `yaml:"port" env-default:"8080"`
		Username string `yaml:"username" env-default:"docker"`
		Password string `yaml:"password" env-default:"docker"`
		Name     string `yaml:"name" env-default:"docker"`
		SSL      string `yaml:"ssl" env-default:"disable"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, err := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
