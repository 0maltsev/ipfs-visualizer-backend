package nodehandlers

import (
	"encoding/json"
	"ipfs-visualizer/internal/services/nodes"
	"net/http"
)

func NewNodeHandler() *NodeHandler {
	return &NodeHandler{

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
