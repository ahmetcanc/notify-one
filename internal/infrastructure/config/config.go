package config

import (
	"os"
	"strconv"

	"github.com/ahmetcanc/notify-one/internal/infrastructure/cache"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/database"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	AppPort       string
}

// LoadConfig loads configuration from environment
func LoadConfig() *Config {
	// Load .env file into the system's environment
	_ = godotenv.Load()

	// Convert string DB to int safely
	dbNum, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		dbNum = 0 // default to 0 if conversion fails
	}

	return &Config{
		// Database
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		// Redis
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       dbNum,

		// Application
		AppPort: os.Getenv("APP_PORT"),
	}
}

// ToDatabaseConfig maps general config to database specific config
func (c *Config) ToDatabaseConfig() database.Config {
	return database.Config{
		Host:     c.DBHost,
		Port:     c.DBPort,
		User:     c.DBUser,
		Password: c.DBPassword,
		DBName:   c.DBName,
		SSLMode:  c.DBSSLMode,
	}
}

// ToRedisConfig maps general config to redis specific config
func (c *Config) ToRedisConfig() cache.Config {
	return cache.Config{
		Host:     c.RedisHost,
		Port:     c.RedisPort,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	}
}
