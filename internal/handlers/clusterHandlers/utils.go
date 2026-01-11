package clusterhandlers

import (
	"encoding/json"
	"ipfs-visualizer/internal/services/clusters"
	"ipfs-visualizer/internal/services/nodes"
	"net/http"
	"database/sql"
	"ipfs-visualizer/config"
	"k8s.io/client-go/kubernetes"
)

func NewClusterHandler(sqlDbPool *sql.DB,  clusterCfg *config.ClusterConfig, kubeClientSet *kubernetes.Clientset, nodeCfg *config.NodeConfig) *ClusterHandler {
	return &ClusterHandler{
		sqlDbPool: sqlDbPool,
		clusterCfg: clusterCfg,
		kubeClientSet: kubeClientSet,
		nodeCfg: nodeCfg,
	}
}

func WriteGetAllClustersResponse(w http.ResponseWriter, clusterList []clusters.ClusterSpec) error {
	response := GetAllClustersResponseBody{
		ClusterList: clusterList,
	}
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response)
}

func WriteGetClusterByIDResponse(w http.ResponseWriter, cluster clusters.ClusterSpec) error {
	response := GetClusterByIDResponseBody{
		Cluster: cluster,
	}
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response)
}

func WriteCreateClusterResponse(w http.ResponseWriter, cluster clusters.ClusterSpec) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func WriteUpdateClusterResponse(w http.ResponseWriter, cluster clusters.ClusterSpec) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func WriteDeleteClusterResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func WriteGetClusterStatusResponse(w http.ResponseWriter, status clusters.ClusterStatus) error {
	response := GetClusterStatusByIDResponseBody{
		Status: status,
	}
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response)
}

func WriteGetClusterNodesResponse(w http.ResponseWriter, nodeList []nodes.NodeSpec) error {
	response := GetClusterNodesByIDResponseBody{
		NodeList: nodeList,
	}
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response)
}

func WriteAddNodeToClusterResponse(w http.ResponseWriter, cluster clusters.ClusterSpec) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func BuildCreateClusterReqBody(clusterReqBody CreateClusterRequestBody) clusters.ClusterSpec {
	return clusters.ClusterSpec{

	}
}

func BuildUpdateClusterReqBody(clusterReqBody UpdateClusterRequestBody) clusters.ClusterSpec {
	return clusters.ClusterSpec{

	}
}

func BuildAddNodeToClusterReqBody(nodeToClusterReqBody AddNodeToClusterRequestBody) clusters.ClusterSpec {
	return clusters.ClusterSpec{

	}
}

func WriteRemoveNodeFromClusterResponse(w http.ResponseWriter, cluster clusters.ClusterSpec) error {
	response := GetClusterByIDResponseBody{
			Cluster: cluster,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

		return json.NewEncoder(w).Encode(response)
	}