package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Masterminds/semver"
	"github.com/joho/godotenv"
)

const (
	Development string = "dev"
	Production         = "prod"
)

type Config struct {
	AppVersion      *semver.Version
	Env             string
	LogLevel        string
	Port            int
	CertFile        string
	KeyFile         string
	SupportedClient *SupportedClient
}

var serverConfig *Config

func init() {
	env, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		env = Development
	}

	filename := fmt.Sprintf(".%s.env", env)
	godotenv.Load(filename)
	serverConfig = new()
}

func new() *Config {
	appVersionStr := getEnv("APP_VERSION", "1.0.0", false)
	appVersion := semver.MustParse(appVersionStr)
	supportedClients := GetSupportedClients(appVersion)

	return &Config{
		// Non-required configs
		AppVersion: appVersion,
		Env:        getEnv("ENVIRONMENT", Development, false),
		Port:       getEnvAsInt("PORT", 8000, false),
		LogLevel:   getEnv("LOG_LEVEL", "debug", false),
		CertFile:   getEnv("CERT_FILE", "", false),
		KeyFile:    getEnv("KEY_FILE", "", false),

		// Required configs
		SupportedClient: supportedClients,
	}
}

func Get() *Config {
	return serverConfig
}

func (c *Config) RunTLS() bool {
	retirn c.CertFile != "" && c.KeyFile != ""
}

func (c *Config) IsValidClient(clientVerStr string) bool {
	clientVersion, err := semver.NewVersion(clientVerStr)
	if err != nil {
		return false
	}

	return c.SupportedClient.Constraints.Check(clientVersion)
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

