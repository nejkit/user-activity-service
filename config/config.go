package config

import (
	"fmt"
	"time"
)

type Config struct {
	BackgroundWorkerConfig `envPrefix:"BACKGROUND_WORKER_"`
	DatabaseConfig         `envPrefix:"DATABASE_"`
	ApplicationPort        int    `env:"APP_PORT" envDefault:"80"`
	LoggerLevel            string `env:"LOG_LEVEL" envDefault:"INFO"`
}

type BackgroundWorkerConfig struct {
	Interval             time.Duration `env:"INTERVAL" envDefault:"10s"`
	EventsPeriodDuration time.Duration `env:"PERIOD_DURATION" envDefault:"1m"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"5432"`
	Username string `env:"USERNAME" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	DbName   string `env:"DB_NAME" envDefault:"postgres"`
}

func (d *DatabaseConfig) ToConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.Username, d.Password, d.DbName)
}
