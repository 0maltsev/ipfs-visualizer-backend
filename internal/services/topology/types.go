package topology

type Topology struct {
	TopologyID   string        `json:"topologyId"`
	Name         string        `json:"name"`
	Nodes        []TopologyNode `json:"nodes"`
	Edges        []TopologyEdge `json:"edges"`
	DeployStatus string        `json:"deployStatus"`
	K8sNamespace *string       `json:"k8sNamespace,omitempty"`
	CreatedAt    string        `json:"createdAt"`
	UpdatedAt    string        `json:"updatedAt"`
}

type TopologySummary struct {
	TopologyID   string `json:"topologyId"`
	Name         string `json:"name"`
	NodeCount    int    `json:"nodeCount"`
	EdgeCount    int    `json:"edgeCount"`
	DeployStatus string `json:"deployStatus"`
	CreatedAt    string `json:"createdAt"`
}

type TopologyNode struct {
	NodeID   string    `json:"nodeId"`
	Label    string    `json:"label"`
	Position Position  `json:"position"`
	Role     string    `json:"role"` // bootstrap | worker
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type TopologyEdge struct {
	EdgeID       string `json:"edgeId"`
	SourceNodeID string `json:"sourceNodeId"` // worker bootstraps to target
	TargetNodeID string `json:"targetNodeId"` // bootstrap
}

type TopologyCreate struct {
	Name  string         `json:"name"`
	Nodes []TopologyNode `json:"nodes,omitempty"`
	Edges []TopologyEdge `json:"edges,omitempty"`
}

type TopologyUpdate struct {
	Name  *string        `json:"name,omitempty"`
	Nodes []TopologyNode `json:"nodes,omitempty"`
	Edges []TopologyEdge `json:"edges,omitempty"`
}

type DeployResult struct {
	TopologyID string `json:"topologyId"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
}

type DeployStatus struct {
	TopologyID   string      `json:"topologyId"`
	Status       string      `json:"status"`
	Message      string      `json:"message,omitempty"`
	Pods         []PodStatus `json:"pods,omitempty"`
}

type PodStatus struct {
	NodeID  string `json:"nodeId"`
	PodName string `json:"podName"`
	Phase   string `json:"phase"`
	Ready   bool   `json:"ready"`
}
