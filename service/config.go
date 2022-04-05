package service

import (
	"github.com/mabaro3009/example-architecture-go/data"
	"time"
)

type Config struct {
	Environment string `envconfig:"environment" default:"local-dev"`

	ListenAddress string `envconfig:"listen_address" default:":8081"`

	DatabaseType data.Type `envconfig:"database_type" default:"0"`

	DatabaseHost           string        `envconfig:"db_host" default:"localhost"`
	DatabasePort           int           `envconfig:"db_port" default:"5432"`
	DatabaseName           string        `envconfig:"db_name" default:"postgres"`
	DatabaseUser           string        `envconfig:"db_user" default:"postgres"`
	DatabasePassword       string        `envconfig:"db_password" default:"postgres"`
	DatabaseConnectTimeout time.Duration `envconfig:"db_connect_timeout" default:"15s"`
}
