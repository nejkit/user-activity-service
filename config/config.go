package config

import (
	"fmt"
	"time"
)

type Config struct {
	BackgroundWorkerConfig
	DatabaseConfig
	ApplicationPort int
}

type BackgroundWorkerConfig struct {
	Interval             time.Duration
	EventsPeriodDuration time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

func (d *DatabaseConfig) ToConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.Username, d.Password, d.DbName)
}
