package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type SchedulerConfig struct {
	Interval string `json:"interval"` // "minute", "hour", "day", "week"
	Time     string `json:"time"`     // Specific time for daily/weekly schedules (e.g., "15:04")
}

type Config struct {
	Database  DatabaseConfig  `json:"database"`
	Scheduler SchedulerConfig `json:"scheduler"`
	APIKey    APIKeyConfig    `json:"api_keys"`
	RateLimit RateLimitConfig `json:"rate_limits"`
	Redis     Redis           `json:"redis"`
}

type RateLimitConfig struct {
	Delinguents int `json:"delinguents"`
	Outstanding int `json:"outstanding"`
	Payment     int `json:"payment"`
}

type APIKeyConfig struct {
	DelinguentsKey string `json:"delinguents_key"`
	OutstandingKey string `json:"outstanding_key"`
	PaymentKey     string `json:"payment_key"`
}

type Redis struct {
	Host     []string `json:"host"`
	Password string   `json:"Password"`
	DB       int      `json:"Db"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Debug: Print the file content
	fmt.Println("File content:", string(file))

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &config, nil
}
