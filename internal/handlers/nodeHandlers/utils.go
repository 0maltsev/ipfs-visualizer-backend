package nodehandlers

import (
	"encoding/json"
	"ipfs-visualizer/internal/services/nodes"
	"net/http"
	"database/sql"
	"ipfs-visualizer/config"
	"k8s.io/client-go/kubernetes"
)

func NewNodeHandler(sqlDbPool *sql.DB, nodeCfg *config.NodeConfig, kubeClientSet *kubernetes.Clientset) *NodeHandler {
	return &NodeHandler{
		sqlDbPool: sqlDbPool,
		nodeCfg: nodeCfg,
		kubeClientSet: kubeClientSet,
	}
}

func WriteGetNodeByIDResponse(w http.ResponseWriter, node nodes.NodeSpec) error {
	response := GetNodeByIDResponseBody{
		Node: node,
	}
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response)
}

func WriteGetNodeLogsByIDResponse(w http.ResponseWriter, nodeLogs string) error {
	response := GetNodeLogsResponseBody{
		NodeLogs: nodeLogs,
	}
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response)
}

func WriteUpdateNodeResponse(w http.ResponseWriter, node nodes.NodeSpec) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func WriteDeleteNodeResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func BuildUpdateNodeReqBody(nodeReqBody UpdateNodeRequestBody) nodes.NodeSpec {
	return nodes.NodeSpec{

	}
}
