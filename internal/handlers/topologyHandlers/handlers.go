package topologyhandlers

import (
	"database/sql"
	"encoding/json"
	"ipfs-visualizer/internal/services/topology"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"k8s.io/client-go/kubernetes"
)

type Handler struct {
	db  *sql.DB
	k8s *kubernetes.Clientset
}

func NewHandler(db *sql.DB, k8s *kubernetes.Clientset) *Handler {
	return &Handler{db: db, k8s: k8s}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	list, err := topology.GetAllTopologies(ctx, h.db)
	if err != nil {
		slog.Error("GetAllTopologies", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "topologyId")
	t, err := topology.GetTopologyByID(ctx, h.db, id)
	if err != nil {
		slog.Error("GetTopologyByID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if t == nil {
		http.Error(w, "topology not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req topology.TopologyCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	t, err := topology.CreateTopology(ctx, h.db, req)
	if err != nil {
		slog.Error("CreateTopology", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(t)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "topologyId")
	var req topology.TopologyUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	t, err := topology.UpdateTopology(ctx, h.db, id, req)
	if err != nil {
		slog.Error("UpdateTopology", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if t == nil {
		http.Error(w, "topology not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "topologyId")
	t, err := topology.GetTopologyByID(ctx, h.db, id)
	if err != nil {
		slog.Error("GetTopologyByID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if t == nil {
		http.Error(w, "topology not found", http.StatusNotFound)
		return
	}
	if err := topology.DeleteTopology(ctx, h.db, id); err != nil {
		slog.Error("DeleteTopology", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Deploy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "topologyId")
	namespace := "default"
	var body struct {
		Namespace string `json:"namespace"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Namespace != "" {
		namespace = body.Namespace
	}
	result, err := topology.DeployTopology(ctx, h.db, h.k8s, id, namespace)
	if err != nil {
		slog.Error("DeployTopology", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(result)
}

func (h *Handler) Undeploy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "topologyId")
	if err := topology.UndeployTopology(ctx, h.db, h.k8s, id); err != nil {
		slog.Error("UndeployTopology", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "topologyId")
	status, err := topology.GetDeployStatus(ctx, h.db, h.k8s, id)
	if err != nil {
		slog.Error("GetDeployStatus", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
