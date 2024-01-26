package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	env_file string = "../server.env"

	env_user_field   string = "AWP_DB_USER"
	env_pass_field   string = "AWP_DB_PASSWORD"
	env_host_field   string = "AWP_DB_HOST"
	env_dbname_field string = "AWP_DB_NAME"
)

type Config struct {
	User     string
	Password string
	Host     string
	DBName   string
}

func (cfg Config) GetPostgresConnString() string {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.DBName,
	)

	return connString
}

func checkConfigFields(cfg *Config) {
	if len(cfg.User) == 0 {
		log.Fatal("please set AWP_DB_USER environment variable")
	}

	if len(cfg.DBName) == 0 {
		log.Fatal("please set AWP_DB_NAME environment variable")
	}

	if len(cfg.Host) == 0 {
		log.Fatal("please set AWP_DB_HOST environment variable")
	}
}

func ParseConfigurationFile() *Config {
	err := godotenv.Load(env_file)
	if err != nil {
		log.Fatalf("rrror loading configuration file '%s': %v", env_file, err.Error())
	}

	config := &Config{
		User:     os.Getenv(env_user_field),
		Password: os.Getenv(env_pass_field),
		Host:     os.Getenv(env_host_field),
		DBName:   os.Getenv(env_dbname_field),
	}

	checkConfigFields(config)

	return config
}
