package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type HTTPServer struct {
	ServeIP           string        `envconfig:"SERVE_IP" default:"127.0.0.1"`
	ServePort         int           `envconfig:"SERVE_PORT" default:"8000"`
	ServerWiteTimeout time.Duration `envconfig:"SERVE_WRITE_TIMEOUT" default:"10s"`
	ServerReadTimeout time.Duration `envconfig:"SERVE_READ_TIMEOUT" default:"10s"`
}

func ServerConfig() (HTTPServer, error) {
	var config HTTPServer

	err := envconfig.Process("", &config)
	return config, err
}

type Database struct {
	DBUser           string        `envconfig:"PGUSER" required:"true"`
	DBPassword       string        `envconfig:"PGPASSWORD" required:"true"`
	DBHost           string        `envconfig:"PGHOST" required:"true"`
	DBPort           int           `envconfig:"PGPORT" default:"5432"`
	DBName           string        `envconfig:"PGDATABASE" required:"true"`
	DBConnectTimeout time.Duration `envconfig:"DB_RETRY_TIMEOUT" default:"30s"`
}

func DBConfig() (Database, error) {
	var config Database

	err := envconfig.Process("", &config)
	return config, err
}
