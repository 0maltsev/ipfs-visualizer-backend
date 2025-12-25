package config

type ServerConfig struct {
	ServerAddressPort string `env:"SERVER_ADDRESS_PORT" envDefault:"3001"`
	DestributionType  string `env:"DESTRIBUTION_TYPE"`
}

type PostgreSqlConfig struct {
	PostrgesqlUsername         string `env:"POSTGRE_SQL_USERNAME"`
	PostrgesqlPassword         string `env:"POSTGRE_SQL_PASSWORD"`
	PostrgesqlHost             string `env:"POSTGRE_SQL_HOST"`
	PostrgesqlPort             string `env:"POSTGRE_SQL_PORT"`
	PostrgesqlName             string `env:"POSTGRE_SQL_NAME"`
	PostgresqlSSLMode          string `env:"POSTGRE_SQL_SSL_MODE"`
	PostgresqlConnectionScheme string `env:"POSTGRE_SQL_CONNECTION_SCHEME"`
}
