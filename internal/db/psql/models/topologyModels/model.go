package topologymodels

import "time"

type TopologyModel struct {
	TopologyID   string     `db:"topology_id" json:"topologyId"`
	Name         string     `db:"name" json:"name"`
	DeployStatus string     `db:"deploy_status" json:"deployStatus"`
	K8sNamespace *string    `db:"k8s_namespace" json:"k8sNamespace,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updatedAt"`
}

type TopologyNodeModel struct {
	TopologyID string   `db:"topology_id" json:"topologyId"`
	NodeID     string   `db:"node_id" json:"nodeId"`
	Label      string   `db:"label" json:"label"`
	PosX       float64  `db:"pos_x" json:"posX"`
	PosY       float64  `db:"pos_y" json:"posY"`
	Role       string   `db:"role" json:"role"`
}

type TopologyEdgeModel struct {
	TopologyID    string `db:"topology_id" json:"topologyId"`
	EdgeID        string `db:"edge_id" json:"edgeId"`
	SourceNodeID  string `db:"source_node_id" json:"sourceNodeId"`
	TargetNodeID  string `db:"target_node_id" json:"targetNodeId"`
}

type TopologySummaryRow struct {
	TopologyID   string     `db:"topology_id"`
	Name         string     `db:"name"`
	DeployStatus string     `db:"deploy_status"`
	K8sNamespace *string    `db:"k8s_namespace"`
	CreatedAt    time.Time  `db:"created_at"`
	NodeCount    int        `db:"node_count"`
	EdgeCount    int        `db:"edge_count"`
}
