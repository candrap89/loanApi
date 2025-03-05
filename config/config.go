package config

import (
	"encoding/json"
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
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
