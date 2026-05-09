package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBUser     string
	DBPassword string
	DBPort     string
	DBName     string
	DBHost     string
}

var AppConfig Config

func init() {
	mode := flag.String("env", "dev", "dev | prod")
	flag.Parse()

	envFile := ".env." + *mode

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("can not env file %s for mode %s", envFile, *mode)
	}

	AppConfig = Config{
		AppEnv:     *mode,
		AppPort:    getEnv("APP_PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBName:     getEnv("DB_NAME"),
	}
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("mongodb://%s:%s", c.DBHost, c.DBPort)
}

func getEnv(key string) string {
	val := os.Getenv(key)

	// if val == "" {
	// 	log.Fatalf("missing env var %s in env file", key)
	// }

	return val
}
