package main

import "github.com/spf13/viper"

type Config struct {
	Logger LoggerConf `yaml:"logger"`
	HTTP   HTTPConf   `yaml:"http"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type HTTPConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func NewConfig(path string) (Config, error) {
	var cfg Config

	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
