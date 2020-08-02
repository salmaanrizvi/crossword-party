package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	Development string = "dev"
	Production         = "prod"
)

type Config struct {
	AppVersion string
	Env        string
	LogLevel   string
	Port       int
	CertFile   string
	KeyFile    string
}

func init() {
	env, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		env = Development
	}

	filename := fmt.Sprintf(".%s.env", env)
	godotenv.Load(filename)
}

func New() *Config {
	return &Config{
		AppVersion: getEnv("APP_VERSION", "1.0.0", false),
		Env:        getEnv("ENVIRONMENT", Development, false),
		Port:       getEnvAsInt("PORT", 8000, false),
		CertFile:   getEnv("CERT_FILE", "", true),
		KeyFile:    getEnv("KEY_FILE", "", true),
	}
}

func getEnv(key string, defaultVal string, required bool) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if required {
		log.Fatalf("Couldn't find environment variable %s. Are you sure its set?", key)
	}

	return defaultVal
}

func getEnvAsInt(key string, defaultVal int, required bool) int {
	valueStr := getEnv(key, "", required)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
