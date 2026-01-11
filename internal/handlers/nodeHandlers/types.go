package nodehandlers

import (
	"database/sql"
	"ipfs-visualizer/config"
	"ipfs-visualizer/internal/services/nodes"

	"k8s.io/client-go/kubernetes"
)

type NodeHandler struct {
	sqlDbPool     *sql.DB
	nodeCfg       *config.NodeConfig
	kubeClientSet *kubernetes.Clientset
}

type GetNodeByIDResponseBody struct {
	Node nodes.NodeSpec `json:"node"`
}

type GetNodeLogsResponseBody struct {
	NodeLogs string `json:"nodeLogs"`
}

type CreateNodeRequestBody struct {
	Node nodes.NodeSpec `json:"node"`
}

type UpdateNodeRequestBody struct {
	Ports *nodes.PortsSpec `json:"ports,omitempty"`
}
