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

type ClusterConfig struct {
	ServiceType      string `env:"CLUSTER_SERVICE_TYPE"`
	StorageClass     string `env:"CLUSTER_STORAGE_CLASS"`
	IpfsImage        string `env:"CLUSTER_IPFS_IMAGE"`
	IpfsClusterImage string `env:"CLUSTER_IPFS_CLUSTER_IMAGE"`
	ConfigMapEnv     string `env:"CLUSTER_CONFIG_MAP_ENV"`
	ConfigMapScript  string `env:"CLUSTER_CONFIG_MAP_SCRIPT"`
	ClusterSecret    string `env:"CLUSTER_SECRET"`
	BootstrapPrivKey string `env:"CLUSTER_BOOTSTRAP_PRIV_KEY"`
	BootstrapPeerID  string `env:"CLUSTER_BOOTSTRAP_PEER_ID"`
}

type NodeConfig struct {
	SwarmTCP     int `env:"NODE_PORT_SWARM_TCP"`
	SwarmUDP     int `env:"NODE_PORT_SWARM_UDP"`
	API          int `env:"NODE_PORT_API"`
	HTTPGateway  int `env:"NODE_PORT_HTTP_GATEWAY"`
	WS           int `env:"NODE_PORT_WS"`
	ClusterAPI   int `env:"NODE_PORT_CLUSTER_API"`
	ClusterProxy int `env:"NODE_PORT_CLUSTER_PROXY"`
	ClusterSwarm int `env:"NODE_PORT_CLUSTER_SWARM"`

	InitContName	string `env:"INIT_CONTAINER_NAME"`
	InitContImage	string `env:"INIT_CONTAINER_IMAGE"`
	InitContCommand	[]string `env:"INIT_CONTAINER_COMMAND"`
}

type KubeConfig struct {
	KubeConfigPath       string `env:"KUBE_CONFIG_PATH"`
	ManualKubeConfigFlag bool   `env:"MANUAL_KUBE_CONFIG_FLAG"`
}
