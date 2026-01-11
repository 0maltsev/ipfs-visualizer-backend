package clusters

import "fmt"

type GetAllClustersErrors struct{ Inner error }

func (e GetAllClustersErrors) Error() string {
	return fmt.Sprintf("failed to get all clusters: %v", e.Inner)
}
func (e GetAllClustersErrors) Unwrap() error { return e.Inner }

type InsertClusterError struct {
	ClusterID string
	Inner     error
}

func (e InsertClusterError) Error() string {
	return fmt.Sprintf("failed to insert cluster with ID %s: %v", e.ClusterID, e.Inner)
}
func (e InsertClusterError) Unwrap() error { return e.Inner }

type CreateClusterKubeResourcesError struct {
	Inner  error
}

func (e CreateClusterKubeResourcesError) Error() string {
	return fmt.Sprintf("failed to create cluster Kubernetes resources: %v", e.Inner)
}
func (e CreateClusterKubeResourcesError) Unwrap() error { return e.Inner }

type InsertClusterKubeResourceError struct {
	ClusterID string
	Inner  error
}

func (e InsertClusterKubeResourceError) Error() string {
	return fmt.Sprintf("failed to insert kube resource for cluster ID %s: %v", e.ClusterID, e.Inner)
}
func (e InsertClusterKubeResourceError) Unwrap() error { return e.Inner }

type GetClusterByIDError struct {
	ClusterID string
	Inner  error
}

func (e GetClusterByIDError) Error() string {
	return fmt.Sprintf("failed to get cluster with ID %s: %v", e.ClusterID, e.Inner)
}
func (e GetClusterByIDError) Unwrap() error { return e.Inner }

type GetClusterKubeResourceError struct {
	ClusterID string
	Inner  error
}

func (e GetClusterKubeResourceError) Error() string {
	return fmt.Sprintf("failed to get kube resource for cluster ID %s: %v", e.ClusterID, e.Inner)
}
func (e GetClusterKubeResourceError) Unwrap() error { return e.Inner }

type DeleteClusterByIDError struct {
	ClusterID string
	Inner  error
}

func (e DeleteClusterByIDError) Error() string {
	return fmt.Sprintf("failed to delete cluster with ID %s: %v", e.ClusterID, e.Inner)
}
func (e DeleteClusterByIDError) Unwrap() error { return e.Inner }

type GetClusterNodesByIDError struct {
	ClusterID string
	Inner  error
}

func (e GetClusterNodesByIDError) Error() string {
	return fmt.Sprintf("failed to get cluster nodes with ID %s: %v", e.ClusterID, e.Inner)
}
func (e GetClusterNodesByIDError) Unwrap() error { return e.Inner }

type AddNodeToClusterByIDError struct {
	ClusterID string
	Inner  error
}

func (e AddNodeToClusterByIDError) Error() string {
	return fmt.Sprintf("failed to add node to cluster with ID %s: %v", e.ClusterID, e.Inner)
}
func (e AddNodeToClusterByIDError) Unwrap() error { return e.Inner }

type DeleteNodeFromClusterByIDError struct {
	ClusterID string
	NodeID string
	Inner  error
}

func (e DeleteNodeFromClusterByIDError) Error() string {
	return fmt.Sprintf("failed to add node %s to cluster with ID %s: %v", e.NodeID, e.ClusterID, e.Inner)
}
func (e DeleteNodeFromClusterByIDError) Unwrap() error { return e.Inner }
