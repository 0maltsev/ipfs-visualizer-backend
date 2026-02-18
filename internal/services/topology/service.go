package topology

import (
	"context"
	"database/sql"
	"fmt"

	topologymodels "ipfs-visualizer/internal/db/psql/models/topologyModels"
	kubetopo "ipfs-visualizer/internal/kube/topology"

	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"
)

func GetAllTopologies(ctx context.Context, db *sql.DB) ([]TopologySummary, error) {
	rows, err := topologymodels.GetAllTopologies(ctx, db)
	if err != nil {
		return nil, err
	}
	result := make([]TopologySummary, 0, len(rows))
	for _, r := range rows {
		result = append(result, TopologySummary{
			TopologyID:   r.TopologyID,
			Name:         r.Name,
			NodeCount:    r.NodeCount,
			EdgeCount:    r.EdgeCount,
			DeployStatus: r.DeployStatus,
			CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result, nil
}

func GetTopologyByID(ctx context.Context, db *sql.DB, id string) (*Topology, error) {
	m, err := topologymodels.GetTopologyByID(ctx, db, id)
	if err != nil || m == nil {
		return nil, err
	}
	nodes, err := topologymodels.GetNodesByTopology(ctx, db, id)
	if err != nil {
		return nil, err
	}
	edges, err := topologymodels.GetEdgesByTopology(ctx, db, id)
	if err != nil {
		return nil, err
	}

	t := &Topology{
		TopologyID:   m.TopologyID,
		Name:         m.Name,
		DeployStatus: m.DeployStatus,
		K8sNamespace: m.K8sNamespace,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	for _, n := range nodes {
		t.Nodes = append(t.Nodes, TopologyNode{
			NodeID:   n.NodeID,
			Label:    n.Label,
			Position: Position{X: n.PosX, Y: n.PosY},
			Role:     n.Role,
		})
	}
	for _, e := range edges {
		t.Edges = append(t.Edges, TopologyEdge{
			EdgeID:       e.EdgeID,
			SourceNodeID: e.SourceNodeID,
			TargetNodeID: e.TargetNodeID,
		})
	}
	return t, nil
}

func CreateTopology(ctx context.Context, db *sql.DB, req TopologyCreate) (*Topology, error) {
	id := uuid.NewString()
	m := &topologymodels.TopologyModel{
		TopologyID:   id,
		Name:         req.Name,
		DeployStatus: "none",
	}
	if err := topologymodels.InsertTopology(ctx, db, m); err != nil {
		return nil, err
	}
	nodeModels := modelsFromNodes(id, req.Nodes)
	if len(nodeModels) > 0 {
		if err := topologymodels.ReplaceTopologyNodes(ctx, db, id, nodeModels); err != nil {
			return nil, err
		}
	}
	edgeModels := modelsFromEdges(id, req.Edges)
	if len(edgeModels) > 0 {
		if err := topologymodels.ReplaceTopologyEdges(ctx, db, id, edgeModels); err != nil {
			return nil, err
		}
	}
	return GetTopologyByID(ctx, db, id)
}

func modelsFromNodes(topologyID string, nodes []TopologyNode) []topologymodels.TopologyNodeModel {
	out := make([]topologymodels.TopologyNodeModel, 0, len(nodes))
	for _, n := range nodes {
		nodeID := n.NodeID
		if nodeID == "" {
			nodeID = uuid.NewString()
		}
		role := n.Role
		if role == "" {
			role = "worker"
		}
		out = append(out, topologymodels.TopologyNodeModel{
			TopologyID: topologyID,
			NodeID:     nodeID,
			Label:      n.Label,
			PosX:       n.Position.X,
			PosY:       n.Position.Y,
			Role:       role,
		})
	}
	return out
}

func modelsFromEdges(topologyID string, edges []TopologyEdge) []topologymodels.TopologyEdgeModel {
	out := make([]topologymodels.TopologyEdgeModel, 0, len(edges))
	for _, e := range edges {
		edgeID := e.EdgeID
		if edgeID == "" {
			edgeID = uuid.NewString()
		}
		out = append(out, topologymodels.TopologyEdgeModel{
			TopologyID:   topologyID,
			EdgeID:       edgeID,
			SourceNodeID: e.SourceNodeID,
			TargetNodeID: e.TargetNodeID,
		})
	}
	return out
}

func UpdateTopology(ctx context.Context, db *sql.DB, id string, req TopologyUpdate) (*Topology, error) {
	m, err := topologymodels.GetTopologyByID(ctx, db, id)
	if err != nil || m == nil {
		return nil, err
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Nodes != nil {
		if err := topologymodels.ReplaceTopologyNodes(ctx, db, id, modelsFromNodes(id, req.Nodes)); err != nil {
			return nil, err
		}
	}
	if req.Edges != nil {
		if err := topologymodels.ReplaceTopologyEdges(ctx, db, id, modelsFromEdges(id, req.Edges)); err != nil {
			return nil, err
		}
	}
	if err := topologymodels.UpdateTopology(ctx, db, m); err != nil {
		return nil, err
	}
	return GetTopologyByID(ctx, db, id)
}

func DeleteTopology(ctx context.Context, db *sql.DB, id string) error {
	_, err := topologymodels.GetTopologyByID(ctx, db, id)
	if err != nil {
		return err
	}
	return topologymodels.DeleteTopology(ctx, db, id)
}

func DeployTopology(ctx context.Context, db *sql.DB, k8s *kubernetes.Clientset, id string, namespace string, private bool) (*DeployResult, error) {
	t, err := GetTopologyByID(ctx, db, id)
	if err != nil || t == nil {
		return nil, fmt.Errorf("topology not found: %s", id)
	}
	if len(t.Nodes) == 0 {
		return nil, fmt.Errorf("topology has no nodes")
	}

	bootstrapID := resolveBootstrapNode(t)
	if bootstrapID == "" {
		return nil, fmt.Errorf("topology must have exactly one bootstrap node (node that others connect to)")
	}

	if err := topologymodels.UpdateTopologyDeployStatus(ctx, db, id, "deploying", &namespace); err != nil {
		return nil, err
	}

	cfg := kubetopo.DeployConfig{
		TopologyID:  id,
		Name:        t.Name,
		Namespace:   namespace,
		BootstrapID: bootstrapID,
		Private:     private,
	}
	for _, n := range t.Nodes {
		cfg.Nodes = append(cfg.Nodes, kubetopo.NodeInfo{
			NodeID:   n.NodeID,
			Label:    n.Label,
			Position: struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			}{X: n.Position.X, Y: n.Position.Y},
			Role: n.Role,
		})
	}
	for _, e := range t.Edges {
		cfg.Edges = append(cfg.Edges, kubetopo.EdgeInfo{SourceNodeID: e.SourceNodeID, TargetNodeID: e.TargetNodeID})
	}
	if err := kubetopo.Deploy(ctx, k8s, cfg); err != nil {
		_ = topologymodels.UpdateTopologyDeployStatus(ctx, db, id, "error", &namespace)
		return nil, err
	}
	_ = topologymodels.UpdateTopologyDeployStatus(ctx, db, id, "running", &namespace)

	return &DeployResult{TopologyID: id, Status: "deploying", Message: "Deployment started"}, nil
}

func resolveBootstrapNode(t *Topology) string {
	targets := make(map[string]bool)
	for _, e := range t.Edges {
		targets[e.TargetNodeID] = true
	}
	sources := make(map[string]bool)
	for _, e := range t.Edges {
		sources[e.SourceNodeID] = true
	}
	for _, n := range t.Nodes {
		if n.Role == "bootstrap" {
			return n.NodeID
		}
	}
	for _, n := range t.Nodes {
		if targets[n.NodeID] && !sources[n.NodeID] {
			return n.NodeID
		}
	}
	if len(t.Nodes) == 1 {
		return t.Nodes[0].NodeID
	}
	return ""
}

func UndeployTopology(ctx context.Context, db *sql.DB, k8s *kubernetes.Clientset, id string) error {
	t, err := GetTopologyByID(ctx, db, id)
	if err != nil || t == nil {
		return fmt.Errorf("topology not found: %s", id)
	}
	ns := "default"
	if t.K8sNamespace != nil {
		ns = *t.K8sNamespace
	}
	if err := kubetopo.Undeploy(ctx, k8s, id, ns); err != nil {
		return err
	}
	return topologymodels.UpdateTopologyDeployStatus(ctx, db, id, "none", nil)
}

func GetDeployStatus(ctx context.Context, db *sql.DB, k8s *kubernetes.Clientset, id string) (*DeployStatus, error) {
	t, err := GetTopologyByID(ctx, db, id)
	if err != nil || t == nil {
		return nil, fmt.Errorf("topology not found: %s", id)
	}
	ns := "default"
	if t.K8sNamespace != nil {
		ns = *t.K8sNamespace
	}
	rawPods, err := kubetopo.GetPodsStatus(ctx, k8s, id, ns)
	if err != nil {
		return &DeployStatus{TopologyID: id, Status: t.DeployStatus, Message: err.Error()}, nil
	}
	pods := make([]PodStatus, 0, len(rawPods))
	for _, p := range rawPods {
		pods = append(pods, PodStatus{NodeID: p.NodeID, PodName: p.PodName, Phase: p.Phase, Ready: p.Ready})
	}
	return &DeployStatus{
		TopologyID: id,
		Status:     t.DeployStatus,
		Pods:       pods,
	}, nil
}

func GetPodLogs(ctx context.Context, db *sql.DB, k8s *kubernetes.Clientset, topologyID, podName, container string) (string, error) {
	t, err := GetTopologyByID(ctx, db, topologyID)
	if err != nil || t == nil {
		return "", fmt.Errorf("topology not found: %s", topologyID)
	}
	ns := "default"
	if t.K8sNamespace != nil {
		ns = *t.K8sNamespace
	}
	return kubetopo.GetPodLogs(ctx, k8s, ns, podName, container)
}
