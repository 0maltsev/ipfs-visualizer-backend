package clusters

import "ipfs-visualizer/internal/services/nodes"

type ClusterSpec struct {
	ClusterID          string           `json:"clusterId"`
	ClusterName        string           `json:"clusterName"`
	Replicas           int              `json:"replicas"`
	ServiceType        string           `json:"serviceType"`
	StorageClass       string           `json:"storageClass"`
	ClusterStorageSize string           `json:"clusterStorageSize"`
	IPFSStorageSize    string           `json:"ipfsStorageSize"`
	Images             ImagesSpec       `json:"images"`
	ConfigMaps         ConfigMapsSpec   `json:"configMaps"`
	Secrets            SecretsSpec      `json:"secrets"`
	BootstrapPeerID    string           `json:"bootstrapPeerId"`
	ClusterSecret      string           `json:"clusterSecret"`
	Nodes              []nodes.NodeSpec `json:"nodes"`
}

type ImagesSpec struct {
	IPFS        string `json:"ipfs"`
	IPFSCluster string `json:"ipfsCluster"`
}

type ConfigMapsSpec struct {
	EnvConfig     string `json:"envConfig"`
	ScriptsConfig string `json:"scriptsConfig"`
}

type SecretsSpec struct {
	ClusterSecret    string `json:"clusterSecret"`
	BootstrapPrivKey string `json:"bootstrapPrivKey"`
}
