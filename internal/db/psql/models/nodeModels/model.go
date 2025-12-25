package nodemodels

import "time"

type NodeSqlModel struct {
	NodeID   string  `db:"node_id" json:"nodeId"`
	NodeName *string `db:"node_name" json:"nodeName,omitempty"`
	Role     string  `db:"role" json:"role"`

	SwarmTCP     int `db:"swarm_tcp" json:"swarmTCP"`
	SwarmUDP     int `db:"swarm_udp" json:"swarmUDP"`
	API          int `db:"api" json:"api"`
	HTTPGateway  int `db:"http_gateway" json:"httpGateway"`
	WS           int `db:"ws" json:"ws"`
	ClusterAPI   int `db:"cluster_api" json:"clusterAPI"`
	ClusterProxy int `db:"cluster_proxy" json:"clusterProxy"`
	ClusterSwarm int `db:"cluster_swarm" json:"clusterSwarm"`

	IPFSStorage    *string `db:"ipfs_storage" json:"ipfsStorage,omitempty"`
	ClusterStorage *string `db:"cluster_storage" json:"clusterStorage,omitempty"`
	ScriptsConfig  *string `db:"scripts_config" json:"scriptsConfig,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
