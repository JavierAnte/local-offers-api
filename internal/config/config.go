package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	JWTSecret string

	PublicBaseURL string
	UploadDir     string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	publicBaseURL, err := validatePublicBaseURL(os.Getenv("PUBLIC_BASE_URL"))
	if err != nil {
		return nil, err
	}

	uploadDir := strings.TrimSpace(os.Getenv("UPLOAD_DIR"))
	if uploadDir == "" {
		return nil, fmt.Errorf("UPLOAD_DIR is not set")
	}

	return &Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		JWTSecret:     jwtSecret,
		PublicBaseURL: publicBaseURL,
		UploadDir:     uploadDir,
	}, nil
}

func validatePublicBaseURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", fmt.Errorf("PUBLIC_BASE_URL is not set")
	}

	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", fmt.Errorf("PUBLIC_BASE_URL must be an absolute HTTP(S) URL")
	}
	if parsed.Path != "" && parsed.Path != "/" {
		return "", fmt.Errorf("PUBLIC_BASE_URL must not contain a path")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("PUBLIC_BASE_URL must not contain a query or fragment")
	}

	return strings.TrimRight(value, "/"), nil
}
