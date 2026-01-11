package clusters

import (
	"context"
	"database/sql"
	"ipfs-visualizer/config"
	clustermodels "ipfs-visualizer/internal/db/psql/models/clusterModels"
	kubemodels "ipfs-visualizer/internal/db/psql/models/kubeModels"
	"ipfs-visualizer/internal/kube"
	"ipfs-visualizer/internal/services/nodes"

	"k8s.io/client-go/kubernetes"
)

// TODO
const namespace = "default"

func GetAllClusters(sqlDBPool *sql.DB) ([]ClusterSpec, error) {
	var clusterList []ClusterSpec
	ctx := context.Background()

	clusterModelList, err := clustermodels.GetAllClusters(ctx, sqlDBPool)
	if err != nil {
		return clusterList, &GetAllClustersErrors{Inner: err}
	}

	for _, clusterModel := range clusterModelList {
		clusterList = append(clusterList, ConvertClusterSqlModelToClusterSpec(clusterModel))
	}

	return clusterList, nil
}

func CreateCluster(cluster ClusterSpec, sqlDBPool *sql.DB, kubeClientSet *kubernetes.Clientset) (ClusterSpec, error) {
	ctx := context.Background()

	clusterModel := ConvertClusterSpecToClusterSqlModel(cluster)
	if err := clustermodels.InsertCluster(ctx, sqlDBPool, &clusterModel); err != nil {
		return cluster, &InsertClusterError{ClusterID: cluster.ClusterID, Inner: err}
	}

	emptyClusterKubeRes := ConvertClusterSpecToClusterKubeResources(cluster, namespace)

	clusterKubeResources, err := kube.CreateCluster(ctx, kubeClientSet, namespace, *emptyClusterKubeRes)
	if err != nil {
		return cluster, &CreateClusterKubeResourcesError{Inner: err}
	}

	clusterKubeResourcesModel := ConvertClusterKubeResourcesToModel(&clusterKubeResources)
	if err := kubemodels.InsertClusterKubeResources(ctx, sqlDBPool, &clusterKubeResourcesModel); err != nil {
		return cluster, &InsertClusterKubeResourceError{ClusterID: cluster.ClusterID, Inner: err}
	}

	return cluster, nil
}

func GetClusterByID(clusterID string, sqlDBPool *sql.DB) (ClusterSpec, error) {
	var cluster ClusterSpec
	ctx := context.Background()
	clusterModel, err := clustermodels.GetClusterByID(ctx, sqlDBPool, clusterID)
	if err != nil {
		return cluster, &GetClusterByIDError{ClusterID: clusterID, Inner: err}
	}
	cluster = ConvertClusterSqlModelToClusterSpec(*clusterModel)
	return cluster, nil
}

func DeleteClusterByID(clusterID string, sqlDBPool *sql.DB, kubeClientSet *kubernetes.Clientset) (error) {
	ctx := context.Background()
	clusterModel, err := clustermodels.GetClusterByID(ctx, sqlDBPool, clusterID)
	if err != nil {
		return &GetClusterByIDError{ClusterID: clusterID, Inner: err}
	}

	clusterKubeResourceModel, err := kubemodels.GetClusterKubeResourcesByClusterID(ctx, sqlDBPool, clusterModel.ClusterID)
	if err != nil {
		return &GetClusterKubeResourceError{ClusterID: clusterID, Inner: err}
	}
	clusterKubeResources := ConvertClusterKubeResourcesModelToClusterKubeResources(*clusterKubeResourceModel)

	if err := kube.DeleteCluster(ctx, kubeClientSet, namespace, clusterKubeResources.ClusterName); err != nil {
		return &DeleteClusterByIDError{ClusterID: clusterID, Inner: err}
	}

	if err := kubemodels.DeleteClusterKubeResourcesByClusterID(ctx, sqlDBPool, clusterID); err != nil {
		return &DeleteClusterByIDError{ClusterID: clusterID, Inner: err}
	}

	return nil
}

// TODO
func UpdateClusterByID(clusterID string, clusterBody ClusterSpec) (ClusterSpec, error) {
	return  clusterBody, nil
}

// TODO
func GetClusterStatusByID(clusterID string) (ClusterStatus, error) {
	return "status", nil
}

func GetClusterNodesByID(clusterID string, sqlDBPool *sql.DB) ([]nodes.NodeSpec, error) {
	ctx := context.Background()
	var nodeList []nodes.NodeSpec

	clusterModel, err := clustermodels.GetClusterByID(ctx, sqlDBPool, clusterID)
	if err != nil {
		return nodeList, &GetClusterNodesByIDError{ClusterID: clusterID, Inner: err}
	}

	for _, nodeID := range clusterModel.NodeIDs {
		nodeSpec, err := nodes.GetNodeByID(nodeID, sqlDBPool)
		if err != nil {
			return nodeList, &GetClusterNodesByIDError{ClusterID: clusterID, Inner: err}
		}
		nodeList = append(nodeList, nodeSpec)
	}
	return nodeList, nil
}

func AddNodeToClusterByID(clusterID string, nodeRole string, sqlDBPool *sql.DB, kubeClientSet *kubernetes.Clientset, nodeCfg *config.NodeConfig) (ClusterSpec, error) {
	ctx := context.Background()
	var cluster ClusterSpec

	clusterModel, err := clustermodels.GetClusterByID(ctx, sqlDBPool, clusterID)
	if err != nil {
		return cluster, &AddNodeToClusterByIDError{ClusterID: clusterID, Inner: err}
	}

	node, err := nodes.CreateNode(nodeRole, clusterID, cluster.ClusterName, sqlDBPool, kubeClientSet, nodeCfg)
	if err != nil {
		return cluster, &AddNodeToClusterByIDError{ClusterID: clusterID, Inner: err}
	}

	clusterModel.NodeIDs = append(clusterModel.NodeIDs, node.NodeID)
	
	if err = clustermodels.UpdateCluster(ctx, sqlDBPool, clusterModel); err != nil {
		return  cluster, &AddNodeToClusterByIDError{ClusterID: clusterID, Inner: err}
	}

	cluster = ConvertClusterSqlModelToClusterSpec(*clusterModel)
	return cluster, nil
}


func RemoveNodeFromClusterByID(clusterID string, nodeID string, sqlDBPool *sql.DB, kubeClientSet *kubernetes.Clientset) (ClusterSpec, error) {
	ctx := context.Background()
	var cluster ClusterSpec

	clusterModel, err := clustermodels.GetClusterByID(ctx, sqlDBPool, clusterID)
	if err != nil {
		return cluster, &DeleteNodeFromClusterByIDError{NodeID:nodeID, ClusterID: clusterID, Inner: err}
	}

	err = nodes.DeleteNodeByID(nodeID, sqlDBPool, kubeClientSet)
	if err != nil {
		return cluster, &DeleteNodeFromClusterByIDError{NodeID:nodeID, ClusterID: clusterID, Inner: err}
	}

	removeByValue(clusterModel.NodeIDs, nodeID)
	
	if err = clustermodels.UpdateCluster(ctx, sqlDBPool, clusterModel); err != nil {
		return  cluster, &DeleteNodeFromClusterByIDError{NodeID:nodeID, ClusterID: clusterID, Inner: err}
	}

	cluster = ConvertClusterSqlModelToClusterSpec(*clusterModel)
	return cluster, nil
}