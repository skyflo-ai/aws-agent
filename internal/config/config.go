package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration settings.
type Config struct {
	AWSAccessKey    string
	AWSSecretKey    string
	AWSRegion       string
	BackendEndpoint string
}

// LoadConfig reads the required environment variables.
func LoadConfig() (*Config, error) {
	awsAccessKey := strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY"))
	awsSecretKey := strings.TrimSpace(os.Getenv("AWS_SECRET_KEY"))
	awsRegion := strings.TrimSpace(os.Getenv("AWS_REGION"))
	backendEndpoint := strings.TrimSpace(os.Getenv("BACKEND_ENDPOINT"))

	if awsAccessKey == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY is not set")
	}
	if awsSecretKey == "" {
		return nil, fmt.Errorf("AWS_SECRET_KEY is not set")
	}
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}
	if backendEndpoint == "" {
		return nil, fmt.Errorf("BACKEND_ENDPOINT is not set")
	}

	return &Config{
		AWSAccessKey:    awsAccessKey,
		AWSSecretKey:    awsSecretKey,
		AWSRegion:       awsRegion,
		BackendEndpoint: backendEndpoint,
	}, nil
}
