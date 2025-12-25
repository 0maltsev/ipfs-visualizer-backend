package psql_connection

import (
	"ipfs-visualizer/config"
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"
)

func NewSqlDBPool(cfg *config.PostgreSqlConfig) (*sql.DB, error) {
	dsn := url.URL{
		Scheme: cfg.PostgresqlConnectionScheme,
		User:   url.UserPassword(cfg.PostrgesqlUsername, cfg.PostrgesqlPassword),
		Host:   fmt.Sprintf("%s:%s", cfg.PostrgesqlHost, cfg.PostrgesqlPort),
		Path:   cfg.PostrgesqlName,
	}

	q := dsn.Query()
	q.Add("sslmode", cfg.PostgresqlSSLMode)
	dsn.RawQuery = q.Encode()

	pool, err := sql.Open("postgres", dsn.String())
	if err != nil {
		return nil, newPostgresConnectionError("NewSqlDBPool", "failed to open PostgreSQL connection", err)
	}

	if err := pool.Ping(); err != nil {
		return nil, newPostgresConnectionError("NewSqlDBPool", "failed to ping PostgreSQL", err)
	}

	return pool, nil
}
