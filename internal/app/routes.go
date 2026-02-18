package app

import (
	topologyhandlers "ipfs-visualizer/internal/handlers/topologyHandlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		th := topologyhandlers.NewHandler(a.sqlDBPool, a.kubernetesClient)
		r.Route("/topologies", func(r chi.Router) {
			r.Get("/", th.GetAll)
			r.Post("/", th.Create)
			r.Get("/{topologyId}", th.GetByID)
			r.Put("/{topologyId}", th.Update)
			r.Delete("/{topologyId}", th.Delete)
			r.Post("/{topologyId}/deploy", th.Deploy)
			r.Post("/{topologyId}/undeploy", th.Undeploy)
			r.Get("/{topologyId}/status", th.GetStatus)
		})
	})

	a.router = router
}
