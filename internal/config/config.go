package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

// Config includes all params of application
type Config struct {
	PrometheusPort string     `yaml:"prometheus_port"`
	HTTPServer     HTTPServer `yaml:"http_server"`
	PostgresDB     PostgresDB `yaml:"postgres_db"`
}

type HTTPServer struct {
	Addr        string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type PostgresDB struct {
	Host              string        `yaml:"host" env-default:"localhost"`
	Port              string        `yaml:"port" env-default:"5432"`
	DBName            string        `yaml:"db_name"`
	Username          string        `yaml:"username"`
	Password          string        `env:"DB_PASSWORD"`
	MaxConns          int           `yaml:"max_conns" env-default:"10"`
	MinConns          int           `yaml:"min_conns" env-default:"0"`
	MaxConnIdleTime   time.Duration `yaml:"max_conn_idle_time" env-default:"5m"`
	HealthCheckPeriod time.Duration `yaml:"health_check_period" env-default:"10s"`
}

// MustLoad loads configuration or stopping application
func MustLoad() *Config {
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		log.Println("CONFIG_PATH environment variable not set")
		configPath = "config/local.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist")
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
