package config

import (
	"log/slog"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerCfg     ServerConfig
	PostgreSqlCfg PostgreSqlConfig
	ClusterCfg    ClusterConfig
	NodeCfg       NodeConfig
	KubeCfg       KubeConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("'.env' is not found. Loading config from envs...")
	}

	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
