package nodehandlers

import "ipfs-visualizer/internal/services/nodes"

type NodeHandler struct{}

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
	Ports *nodes.PortsSpec  `json:"ports,omitempty"`
}
