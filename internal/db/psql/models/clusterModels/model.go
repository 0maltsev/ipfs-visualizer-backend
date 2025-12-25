package clustermodels

import "time"

type ClusterSqlModel struct {
	ClusterID          string  `db:"cluster_id" json:"clusterId"`
	ClusterName        *string `db:"cluster_name" json:"clusterName,omitempty"`
	Replicas           int     `db:"replicas" json:"replicas"`
	ServiceType        *string `db:"service_type" json:"serviceType,omitempty"`
	StorageClass       *string `db:"storage_class" json:"storageClass,omitempty"`
	ClusterStorageSize *string `db:"cluster_storage_size" json:"clusterStorageSize,omitempty"`
	IPFSStorageSize    *string `db:"ipfs_storage_size" json:"ipfsStorageSize,omitempty"`

	IPFSImage        *string `db:"ipfs_image" json:"ipfs,omitempty"`
	IPFSClusterImage *string `db:"ipfs_cluster_image" json:"ipfsCluster,omitempty"`

	EnvConfig     *string `db:"env_config" json:"envConfig,omitempty"`
	ScriptsConfig *string `db:"scripts_config" json:"scriptsConfig,omitempty"`

	ClusterSecret    *string `db:"cluster_secret" json:"clusterSecret,omitempty"`
	BootstrapPrivKey *string `db:"bootstrap_priv_key" json:"bootstrapPrivKey,omitempty"`
	BootstrapPeerID  *string `db:"bootstrap_peer_id" json:"bootstrapPeerId,omitempty"`

	NodeIDs     []string  `db:"nodes" json:"nodes"` // тут только ids
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
