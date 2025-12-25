package app

import (
	clusterhandlers "ipfs-visualizer/internal/handlers/clusterHandlers"
	nodehandlers "ipfs-visualizer/internal/handlers/nodeHandlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		r.Route("/clusters", a.loadClusterRoutes)
		r.Route("/nodes", a.loadNodeRoutes)
	})

	a.router = router
}

func (a *App) loadClusterRoutes(router chi.Router) {
	clusterHandler := clusterhandlers.NewClusterHandler()

	router.Get("/", clusterHandler.GetAllClusters)
	router.Post("/", clusterHandler.CreateCluster)
	router.Get("/{clusterID}", clusterHandler.GetClusterByID)
	router.Delete("/{clusterID}", clusterHandler.DeleteClusterByID)
	router.Put("/{clusterID}", clusterHandler.UpdateClusterByID)
	router.Get("/{clusterID}/status", clusterHandler.GetClusterStatusByID)
	router.Get("/{clusterID}/nodes", clusterHandler.GetClusterNodesByID)
	router.Post("/{clusterID}/nodes", clusterHandler.AddNodeToClusterByID)
}

func (a *App) loadNodeRoutes(router chi.Router) {
	nodeHandler := nodehandlers.NewNodeHandler()

	router.Get("/{nodeID}", nodeHandler.GetNodeByID)
	router.Delete("/{nodeID}", nodeHandler.DeleteNodeByID)
	router.Put("/{nodeID}", nodeHandler.UpdateNodeByID)
	router.Get("/{nodeID}/logs", nodeHandler.GetNodeLogsByID)
}
