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
	// Development is the string for the server running in development
	Development string = "dev"
	// Production is the string for the server running in production
	Production = "prod"
)

// Config is the servers dynamic configuration
type Config struct {
	AppVersion *semver.Version
	Env        string
	Port       int

	LogLevel         string
	LogStatsInterval int

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
	fmt.Printf("Loaded config \n%+v\n", serverConfig)
}

func new() *Config {
	appVersionStr := getEnv("APP_VERSION", "1.0.0", false)
	appVersion := semver.MustParse(appVersionStr)
	supportedClient := GetSupportedClient(appVersion)

	return &Config{
		// Non-required configs
		AppVersion:       appVersion,
		Env:              getEnv("ENVIRONMENT", Development, false),
		Port:             getEnvAsInt("PORT", 8000, false),
		LogLevel:         getEnv("LOG_LEVEL", "debug", false),
		LogStatsInterval: getEnvAsInt("LOG_STATS_INTERVAL", 300, false),
		CertFile:         getEnv("CERT_FILE", "", false),
		KeyFile:          getEnv("KEY_FILE", "", false),

		// Required configs
		SupportedClient: supportedClient,
	}
}

// Get returns the singleton Config instance for the server
func Get() *Config {
	return serverConfig
}

// RunTLS checks if the configuration was loaded with a CertFile and KeyFile for
// running the server in TLS mode. This is useful in local dev mode
func (c *Config) RunTLS() bool {
	return c.CertFile != "" && c.KeyFile != ""
}

// IsValidClient validates whether clientVerStr is a client version
// that this server accepts
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
