package clusterhandlers

import (
	"encoding/json"
	"ipfs-visualizer/internal/handlers"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *ClusterHandler) GetAllClusters(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to get all clusters")
	clusterList, err := GetAllClusters()
	if err != nil {
		slog.Error("Error getting all agents", "error", err)
		http.Error(w, handlers.NewClusterError("GetAllClusters", "failed to get agents", err).Error(), http.StatusInternalServerError)
		return
	}

	if err := WriteGetAllClustersResponse(w, clusterList); err != nil {
		slog.Error("Error encoding response", "error", err)
		http.Error(w, handlers.NewResponseError("WriteGetAllClustersResponse", "failed to encode response", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ClusterHandler) CreateCluster(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to create cluster")

	var clusterReqBody CreateClusterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&clusterReqBody); err != nil {
		http.Error(w, handlers.NewRequestError("CreateCluster", "invalid request body", err).Error(), http.StatusBadRequest)
		return
	}

	clusterBody := BuildCreateClusterReqBody(clusterReqBody)
	cluster, err := CreateCluster(clusterBody)
	if err != nil {
		slog.Error("Error creating cluster", "error", err)
		http.Error(w, handlers.NewClusterError("CreateCluster", "failed to create cluster", err).Error(), http.StatusInternalServerError)
		return
	}

	WriteCreateClusterResponse(w, cluster)
}

func (h *ClusterHandler) GetClusterByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to get cluster by ID")

	clusterID := chi.URLParam(r, "clusterID")
	intClusterID, err := strconv.Atoi(clusterID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("GetClusterByID", "invalid cluster ID", err).Error(), http.StatusBadRequest)
		return
	}

	cluster, err := GetClusterByID(intClusterID)
	if err != nil {
		slog.Error("Error getting cluster by ID", "error", err)
		http.Error(w, handlers.NewClusterError("GetClusterByID", "failed to get cluster", err).Error(), http.StatusInternalServerError)
		return
	}

	if err := WriteGetClusterByIDResponse(w, cluster); err != nil {
		slog.Error("Error encoding response", "error", err)
		http.Error(w, handlers.NewResponseError("WriteGetClusterByIDResponse", "failed to encode response", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ClusterHandler) DeleteClusterByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to delete cluster")

	clusterID := chi.URLParam(r, "clusterID")
	intClusterID, err := strconv.Atoi(clusterID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("DeleteClusterByID", "invalid cluster ID", err).Error(), http.StatusBadRequest)
		return
	}

	if err := DeleteClusterByID(intClusterID); err != nil {
		slog.Error("Error deleting cluster by ID", "error", err)
		http.Error(w, handlers.NewClusterError("DeleteClusterByID", "failed to delete cluster", err).Error(), http.StatusInternalServerError)
		return
	}

	WriteDeleteClusterResponse(w)
}

func (h *ClusterHandler) UpdateClusterByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to update cluster")

	var clusterReqBody UpdateClusterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&clusterReqBody); err != nil {
		http.Error(w, handlers.NewRequestError("UpdateClusterByID", "invalid request body", err).Error(), http.StatusBadRequest)
		return
	}

	clusterBody := BuildUpdateClusterReqBody(clusterReqBody)
	clusterID := chi.URLParam(r, "clusterID")
	intClusterID, err := strconv.Atoi(clusterID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("UpdateClusterByID", "invalid cluster ID", err).Error(), http.StatusBadRequest)
		return
	}

	cluster, err := UpdateClusterByID(intClusterID, clusterBody)
	if err != nil {
		slog.Error("Error updating cluster by ID", "error", err)
		http.Error(w, handlers.NewClusterError("UpdateClusterByID", "failed to update cluster", err).Error(), http.StatusInternalServerError)
		return
	}

	WriteUpdateClusterResponse(w, cluster)
}

func (h *ClusterHandler) GetClusterStatusByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to get cluster status by ID")

	clusterID := chi.URLParam(r, "clusterID")
	intClusterID, err := strconv.Atoi(clusterID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("GetClusterStatusByID", "invalid cluster ID", err).Error(), http.StatusBadRequest)
		return
	}

	status, err := GetClusterStatusByID(intClusterID)
	if err != nil {
		slog.Error("Error getting cluster status by ID", "error", err)
		http.Error(w, handlers.NewClusterError("GetClusterStatusByID", "failed to get cluster status", err).Error(), http.StatusInternalServerError)
		return
	}

	if err := WriteGetClusterStatusResponse(w, status); err != nil {
		slog.Error("Error encoding response", "error", err)
		http.Error(w, handlers.NewResponseError("WriteGetClusterStatusResponse", "failed to encode response", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ClusterHandler) GetClusterNodesByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to get cluster nodes by ID")

	clusterID := chi.URLParam(r, "clusterID")
	intClusterID, err := strconv.Atoi(clusterID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("GetClusterNodesByID", "invalid cluster ID", err).Error(), http.StatusBadRequest)
		return
	}

	nodeList, err := GetClusterNodesByID(intClusterID)
	if err != nil {
		slog.Error("Error getting cluster nodes by ID", "error", err)
		http.Error(w, handlers.NewClusterError("GetClusterNodesByID", "failed to get cluster nodes", err).Error(), http.StatusInternalServerError)
		return
	}

	if err := WriteGetClusterNodesResponse(w, nodeList); err != nil {
		slog.Error("Error encoding response", "error", err)
		http.Error(w, handlers.NewResponseError("WriteGetClusterNodesResponse", "failed to encode response", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ClusterHandler) AddNodeToClusterByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to add node to cluster")

	var nodeToClusterReqBody AddNodeToClusterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&nodeToClusterReqBody); err != nil {
		http.Error(w, handlers.NewRequestError("AddNodeToClusterByID", "invalid request body", err).Error(), http.StatusBadRequest)
		return
	}

	clusterBody := BuildAddNodeToClusterReqBody(nodeToClusterReqBody)
	clusterID := chi.URLParam(r, "clusterID")
	intClusterID, err := strconv.Atoi(clusterID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("AddNodeToClusterByID", "invalid cluster ID", err).Error(), http.StatusBadRequest)
		return
	}

	cluster, err := AddNodeToClusterByID(intClusterID, clusterBody)
	if err != nil {
		slog.Error("Error adding node to cluster by ID", "error", err)
		http.Error(w, handlers.NewClusterError("AddNodeToClusterByID", "failed to add node to cluster", err).Error(), http.StatusInternalServerError)
		return
	}

	WriteAddNodeToClusterResponse(w, cluster)
}
