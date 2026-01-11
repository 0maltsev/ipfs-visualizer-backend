package nodes

import "fmt"

type GetNodeByIDError struct {
	NodeID string
	Inner  error
}

func (e GetNodeByIDError) Error() string {
	return fmt.Sprintf("failed to get node with ID %s: %v", e.NodeID, e.Inner)
}
func (e GetNodeByIDError) Unwrap() error { return e.Inner }

type GetNodeKubeResourceError struct {
	NodeID string
	Inner  error
}

func (e GetNodeKubeResourceError) Error() string {
	return fmt.Sprintf("failed to get kube resource for node ID %s: %v", e.NodeID, e.Inner)
}
func (e GetNodeKubeResourceError) Unwrap() error { return e.Inner }

type InsertNodeError struct {
	NodeID string
	Inner     error
}

func (e InsertNodeError) Error() string {
	return fmt.Sprintf("failed to insert node with ID %s: %v", e.NodeID, e.Inner)
}
func (e InsertNodeError) Unwrap() error { return e.Inner }

type CreateNodeKubeResourcesError struct {
	Inner  error
}

func (e CreateNodeKubeResourcesError) Error() string {
	return fmt.Sprintf("failed to create node Kubernetes resources: %v", e.Inner)
}
func (e CreateNodeKubeResourcesError) Unwrap() error { return e.Inner }

type InsertNodeKubeResourceError struct {
	NodeID string
	Inner  error
}

func (e InsertNodeKubeResourceError) Error() string {
	return fmt.Sprintf("failed to insert kube resource for node ID %s: %v", e.NodeID, e.Inner)
}
func (e InsertNodeKubeResourceError) Unwrap() error { return e.Inner }

type DeleteNodeByIDError struct {
	NodeID string
	Inner  error
}

func (e DeleteNodeByIDError) Error() string {
	return fmt.Sprintf("failed to delete node with ID %s: %v", e.NodeID, e.Inner)
}
func (e DeleteNodeByIDError) Unwrap() error { return e.Inner }
