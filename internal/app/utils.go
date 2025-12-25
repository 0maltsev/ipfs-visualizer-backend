package app

import (
	"ipfs-visualizer/config"
	"database/sql"
	psql_connection "ipfs-visualizer/internal/db/psql/connection"
)


func (a *App) loadServerCfg(cfg *config.Config) {
	serverCfg := ComposeServerCfgFromGeneralCfg(cfg)
	a.serverCfg = serverCfg
}

func ComposeServerCfgFromGeneralCfg(generalCfg *config.Config) *config.ServerConfig {
	return &config.ServerConfig{
		ServerAddressPort: generalCfg.ServerCfg.ServerAddressPort,
		DestributionType:  generalCfg.ServerCfg.DestributionType,
	}
}

func ComposePostgresqlCfgFromGeneralCfg(generalCfg *config.Config) *config.PostgreSqlConfig {
	return &config.PostgreSqlConfig{
		PostrgesqlUsername:         generalCfg.PostgreSqlCfg.PostrgesqlUsername,
		PostrgesqlPassword:         generalCfg.PostgreSqlCfg.PostrgesqlPassword,
		PostrgesqlHost:             generalCfg.PostgreSqlCfg.PostrgesqlHost,
		PostrgesqlPort:             generalCfg.PostgreSqlCfg.PostrgesqlPort,
		PostrgesqlName:             generalCfg.PostgreSqlCfg.PostrgesqlName,
		PostgresqlSSLMode:          generalCfg.PostgreSqlCfg.PostgresqlSSLMode,
		PostgresqlConnectionScheme: generalCfg.PostgreSqlCfg.PostgresqlConnectionScheme,
	}
}

func (a *App) createStorageConnections(cfg *config.Config) (error) {
	sqlDBCfg := ComposePostgresqlCfgFromGeneralCfg(cfg)
	a.sqlDBCfg = sqlDBCfg
	
	sqlPool, err := psql_connection.NewSqlDBPool(sqlDBCfg)
	if err != nil {
		return NewStorageError("CreateStorageConnections", "failed to create postgresql pool", err)
	}

	a.sqlDBPool = sqlPool
	return nil
}

func CreatePSQLTablesIfNotExist(sqlPool *sql.DB) {
	// TODO
}

func (a *App) CloseStorageConnections() error {
	if err := a.sqlDBPool.Close(); err != nil {
		return NewStorageError("CloseStorageConnections", "failed to close sql db pool", err)
	}
	return nil
}
