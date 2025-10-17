package config

import "github.com/caarlos0/env/v11"

func GetConfig() *Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
