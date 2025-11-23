package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr  string
	DatabaseURL string
	EthNodeURL  string
	StartBlock  uint64
	JWTSecret   string
}

func Load() (*Config, error) {
	// 尝试加载 .env 文件（如果存在）
	_ = godotenv.Load()

	cfg := &Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/wallet?sslmode=disable"),
		EthNodeURL:  getEnv("ETH_NODE_URL", "https://mainnet.infura.io/v3/YOUR-PROJECT-ID"),
		StartBlock:  0, // 可以从环境变量读取
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
