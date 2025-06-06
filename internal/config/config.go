package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AppPort      string
	AppSecretKey string
	Postgres     Postgres
	S3           S3
}

type Postgres struct {
	DBHost    string
	DBPort    string
	DBUser    string
	DBName    string
	DBPass    string
	DBSSLMode string
}

type S3 struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	UseSSL    bool
}

// Вспомогательная функция для env с дефолтом и логом
func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Printf("%s environment variable is not set. Using default value: %s\n", key, def)
		return def
	}
	return val
}

// Для булевых значений
func getEnvBool(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "" {
		fmt.Printf("%s environment variable is not set. Using default value: %v\n", key, def)
		return def
	}
	val = strings.ToLower(val)
	return val == "true" || val == "1"
}

func GetConfig() Config {
	return Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		AppSecretKey: getEnv("APP_SECRET_KEY", "secret"),
		Postgres: Postgres{
			DBHost:    getEnv("DB_HOST", "localhost"),
			DBPort:    getEnv("DB_PORT", "5432"),
			DBUser:    getEnv("DB_USER", "postgres"),
			DBName:    getEnv("DB_NAME", "postgres"),
			DBPass:    getEnv("DB_PASS", "password"),
			DBSSLMode: getEnv("DBSSL_MODE", "disable"),
		},
		S3: S3{
			Endpoint:  getEnv("S3_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("S3_BUCKET", "mybucket"),
			Region:    getEnv("S3_REGION", "us-east-1"),
			UseSSL:    getEnvBool("S3_USE_SSL", false),
		},
	}
}
