package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Http     HttpServer `json:"http" yaml:"http"`
	Database Database   `json:"database" yaml:"database"`
	Storage  Storage    `json:"storage" yaml:"storage"`
	Mail     MailConfig `json:"mail" yaml:"mail"`
}

type MailConfig struct {
	From     string `json:"from" yaml:"from"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
}

type Storage struct {
	CloudName string `json:"cloudName" yaml:"cloudName"`
	APIKey    string `json:"apiKey" yaml:"apiKey"`
	APISecret string `json:"apiSecret" yaml:"apiSecret"`
}

type HttpServer struct {
	Path string `json:"path" yaml:"path"`
	Port string `json:"port" yaml:"port"`
}

type Database struct {
	Driver   string `json:"driver" yaml:"driver"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	DBName   string `json:"dbname" yaml:"dbname"`
	SSLMode  string `json:"sslmode" yaml:"sslmode"`
	// Legacy field for backward compatibility
	Source string `json:"source" yaml:"source"`
}

// GetDSN trả về connection string cho PostgreSQL
func (d *Database) GetDSN() string {
	if d.Source != "" {
		return d.Source
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

// LoadConfig đọc file config từ đường dẫn được chỉ định
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("không thể đọc file config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("không thể parse config: %w", err)
	}

	return &config, nil
}
