package nodes

import (
	"context"
	"database/sql"
	"ipfs-visualizer/config"
	kubemodels "ipfs-visualizer/internal/db/psql/models/kubeModels"
	nodemodels "ipfs-visualizer/internal/db/psql/models/nodeModels"
	"ipfs-visualizer/internal/kube"

	"k8s.io/client-go/kubernetes"
)

// TODO
const namespace = "default"

func GetNodeByID(nodeID string, sqlDBPool *sql.DB) (NodeSpec, error) {
	var node NodeSpec
	ctx := context.Background()
	nodeModel, err := nodemodels.GetNodeByID(ctx, sqlDBPool, nodeID)
	if err != nil {
		return node, &GetNodeByIDError{NodeID: nodeID, Inner: err}
	}
	node = ConvertNodeSqlModelToNodeSpec(*nodeModel)
	return node, nil
}

func CreateNode(nodeRole string, clusterID string, clusterName string, sqlDBPool *sql.DB, kubeClientSet *kubernetes.Clientset, nodeCfg *config.NodeConfig) (NodeSpec, error) {
	ctx := context.Background()

	node := BuildNodeSpecFromRole(nodeRole, nodeCfg)
	nodeModel := ConvertNodeSpecToNodeSqlModel(node)
	if err := nodemodels.InsertNode(ctx, sqlDBPool, &nodeModel); err != nil {
		return node, &InsertNodeError{NodeID: node.NodeID, Inner: err}
	}
	nodeKubeResources := ConvertNodeSpecToNodeKubeResources(node, clusterID, namespace)

	_, err := kube.AddNodeToCluster(ctx, kubeClientSet, namespace, clusterName, *nodeKubeResources)
	if err != nil {
		return node, &CreateNodeKubeResourcesError{Inner: err}
	}

	nodeKubeResourcesModel, err := ConvertNodeKubeResourcesToModel(nodeKubeResources)
	if err != nil {
		return node, &InsertNodeKubeResourceError{NodeID: node.NodeID, Inner: err}
	}

	if err := kubemodels.InsertNodeKubeResources(ctx, sqlDBPool, nodeKubeResourcesModel); err != nil {
		return node, &InsertNodeKubeResourceError{NodeID: node.NodeID, Inner: err}
	}

	return node, nil
}

func DeleteNodeByID(nodeID string, sqlDBPool *sql.DB, kubeClientSet *kubernetes.Clientset) (error) {
	ctx := context.Background()
	nodeModel, err := nodemodels.GetNodeByID(ctx, sqlDBPool, nodeID)
	if err != nil {
		return &GetNodeByIDError{NodeID: nodeID, Inner: err}
	}
	
	if err := kube.RemoveNodeFromCluster(ctx, kubeClientSet, namespace, nodeModel.ClusterName); err != nil {
		return &DeleteNodeByIDError{NodeID: nodeID, Inner: err}
	}

	if err := kubemodels.DeleteNodeKubeResourcesByNodeID(ctx, sqlDBPool, nodeID); err != nil {
		return &DeleteNodeByIDError{NodeID: nodeID, Inner: err}
	}

	return nil
}

// TODO
func UpdateNodeByID(nodeID string, nodeBody NodeSpec) (NodeSpec, error) {
	return nodeBody, nil
}

// TODO
func GetNodeLogsByID(nodeID string) (string, error) {
	return "logs", nil
}