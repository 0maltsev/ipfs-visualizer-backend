package clusterhandlers

import (
	"database/sql"
	"ipfs-visualizer/config"
	"ipfs-visualizer/internal/services/clusters"
	"ipfs-visualizer/internal/services/nodes"

	"k8s.io/client-go/kubernetes"
)

type ClusterHandler struct {
	sqlDbPool  *sql.DB
	clusterCfg *config.ClusterConfig
	kubeClientSet *kubernetes.Clientset
	nodeCfg *config.NodeConfig
}

type GetAllClustersResponseBody struct {
	ClusterList []clusters.ClusterSpec `json:"clusterList"`
}

type GetClusterByIDResponseBody struct {
	Cluster clusters.ClusterSpec `json:"cluster"`
}

type CreateClusterRequestBody struct {
	ClusterName        string           `json:"clusterName" validate:"required"`
	Replicas           int              `json:"replicas" validate:"required,min=1"`
	ClusterStorageSize string           `json:"clusterStorageSize,omitempty"`
	IPFSStorageSize    string           `json:"ipfsStorageSize,omitempty"`
	Nodes              []nodes.NodeSpec `json:"nodes,omitempty"`
}

type UpdateClusterRequestBody struct {
	Replicas           *int                 `json:"replicas,omitempty"`
	ServiceType        *string              `json:"serviceType,omitempty"`
	StorageClass       *string              `json:"storageClass,omitempty"`
	ClusterStorageSize *string              `json:"clusterStorageSize,omitempty"`
	IPFSStorageSize    *string              `json:"ipfsStorageSize,omitempty"`
	Nodes              []nodes.NodeSpec     `json:"nodes,omitempty"`
	Images             *clusters.ImagesSpec `json:"images,omitempty"`
}

type GetClusterStatusByIDResponseBody struct {
	Status clusters.ClusterStatus `json:"status"` // creating, running, degraded, error
}

type GetClusterNodesByIDResponseBody struct {
	NodeList []nodes.NodeSpec `json:"nodeList"`
}

type AddNodeToClusterRequestBody struct {
	NodeRole string `json:"role"`
}

type DeleteNodeByIDResponseBody struct {
	Cluster clusters.ClusterSpec
}