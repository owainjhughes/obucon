package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	DBName         string
	JWTSecret      string
	DBSSLMode      string
	AllowedOrigins []string
	GinMode        string
	TrustedProxies []string
	CookieSecure   bool
}

func Load() *Config {
	godotenv.Load()

	if secretName := os.Getenv("APP_SECRET_NAME"); secretName != "" {
		if err := loadSecretsManagerEnv(secretName); err != nil {
			log.Printf("Warning: failed to load secrets from Secrets Manager (%s): %v", secretName, err)
		}
	}

	allowedOrigins := getEnv("ALLOWED_ORIGINS", "")
	if allowedOrigins == "" {
		allowedOrigins = getEnv("ALLOWED_ORIGIN", "http://localhost:3000,http://127.0.0.1:3000")
	}

	ginMode := getEnv("GIN_MODE", "")
	if ginMode == "" {
		if strings.EqualFold(getEnv("APP_ENV", ""), "production") {
			ginMode = "release"
		} else {
			ginMode = "debug"
		}
	}

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPass:         getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "obucon"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-key"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		AllowedOrigins: splitCSV(allowedOrigins, []string{"http://localhost:3000", "http://127.0.0.1:3000"}),
		GinMode:        ginMode,
		TrustedProxies: splitCSV(getEnv("TRUSTED_PROXIES", "127.0.0.1,::1"), []string{"127.0.0.1", "::1"}),
		CookieSecure:   getEnvBool("COOKIE_SECURE", false),
	}
}

// loadSecretsManagerEnv fetches the named secret from AWS Secrets Manager and
// injects each key into the process environment. Existing env vars are not
// overwritten, so docker-compose / local overrides still take precedence.
func loadSecretsManagerEnv(secretName string) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	client := secretsmanager.NewFromConfig(cfg)
	result, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return err
	}

	var secrets map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secrets); err != nil {
		return err
	}

	for k, v := range secrets {
		if _, exists := os.LookupEnv(k); !exists {
			os.Setenv(k, v)
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func splitCSV(value string, defaults []string) []string {
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		origins = append(origins, trimmed)
	}

	if len(origins) == 0 {
		return defaults
	}

	return origins
}
