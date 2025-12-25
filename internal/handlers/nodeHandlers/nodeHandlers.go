package nodehandlers

import (
	"encoding/json"
	"ipfs-visualizer/internal/handlers"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *NodeHandler) GetNodeByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to get node by ID")

	nodeID := chi.URLParam(r, "nodeID")
	intNodeID, err := strconv.Atoi(nodeID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("GetNodeByID", "invalid node ID", err).Error(), http.StatusBadRequest)
		return
	}

	node, err := GetNodeByID(intNodeID)
	if err != nil {
		slog.Error("Error getting node by ID", "error", err)
		http.Error(w, handlers.NewNodeError("GetNodeByID", "failed to get node", err).Error(), http.StatusInternalServerError)
		return
	}

	if err := WriteGetNodeByIDResponse(w, node); err != nil {
		slog.Error("Error encoding response", "error", err)
		http.Error(w, handlers.NewResponseError("WriteGetNodeByIDResponse", "failed to encode response", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NodeHandler) DeleteNodeByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to delete node")

	nodeID := chi.URLParam(r, "nodeID")
	intNodeID, err := strconv.Atoi(nodeID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("DeleteNodeByID", "invalid node ID", err).Error(), http.StatusBadRequest)
		return
	}

	if err := DeleteNodeByID(intNodeID); err != nil {
		slog.Error("Error deleting node by ID", "error", err)
		http.Error(w, handlers.NewNodeError("DeleteNodeByID", "failed to delete node", err).Error(), http.StatusInternalServerError)
		return
	}

	WriteDeleteNodeResponse(w)
}

func (h *NodeHandler) UpdateNodeByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to update node")

	var nodeReqBody UpdateNodeRequestBody
	if err := json.NewDecoder(r.Body).Decode(&nodeReqBody); err != nil {
		http.Error(w, handlers.NewRequestError("UpdateNodeByID", "invalid request body", err).Error(), http.StatusBadRequest)
		return
	}

	nodeBody := BuildUpdateNodeReqBody(nodeReqBody)
	nodeID := chi.URLParam(r, "nodeID")
	intNodeID, err := strconv.Atoi(nodeID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("UpdateNodeByID", "invalid node ID", err).Error(), http.StatusBadRequest)
		return
	}

	node, err := UpdateNodeByID(intNodeID, nodeBody)
	if err != nil {
		slog.Error("Error updating node by ID", "error", err)
		http.Error(w, handlers.NewNodeError("UpdateNodeByID", "failed to update node", err).Error(), http.StatusInternalServerError)
		return
	}

	WriteUpdateNodeResponse(w, node)
}

func (h *NodeHandler) GetNodeLogsByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("Server got request to get node logs by ID")

	nodeID := chi.URLParam(r, "nodeID")
	intNodeID, err := strconv.Atoi(nodeID)
	if err != nil {
		slog.Error("Can't parse ID from URL", "error", err)
		http.Error(w, handlers.NewRequestError("GetNodeLogsByID", "invalid node ID", err).Error(), http.StatusBadRequest)
		return
	}

	nodeLogs, err := GetNodeLogsByID(intNodeID)
	if err != nil {
		slog.Error("Error getting node by ID", "error", err)
		http.Error(w, handlers.NewNodeError("GetNodeLogsByID", "failed to get node logs", err).Error(), http.StatusInternalServerError)
		return
	}

	if err := WriteGetNodeLogsByIDResponse(w, nodeLogs); err != nil {
		slog.Error("Error encoding response", "error", err)
		http.Error(w, handlers.NewResponseError("WriteGetNodeLogsByIDResponse", "failed to encode response", err).Error(), http.StatusInternalServerError)
		return
	}
}
