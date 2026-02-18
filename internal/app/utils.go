package app

import (
	"database/sql"
	"ipfs-visualizer/config"
	psql_connection "ipfs-visualizer/internal/db/psql/connection"
	topologymodels "ipfs-visualizer/internal/db/psql/models/topologyModels"
	"ipfs-visualizer/internal/kube"
	"log/slog"

	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

func (a *App) loadGeneralCfg(cfg *config.Config) {
	a.serverCfg = &cfg.ServerCfg
}

func (a *App) createStorageConnections(cfg *config.Config) error {
	a.sqlDBCfg = &cfg.PostgreSqlCfg

	sqlPool, err := psql_connection.NewSqlDBPool(&cfg.PostgreSqlCfg)
	if err != nil {
		return NewStorageError("CreateStorageConnections", "failed to create postgresql pool", err)
	}

	a.sqlDBPool = sqlPool
	CreatePSQLTablesIfNotExist(sqlPool)
	return nil
}

func CreatePSQLTablesIfNotExist(sqlPool *sql.DB) {
	topologymodels.CreateTopologyTablesIfNotExist(sqlPool)
}

func (a *App) CloseStorageConnections() error {
	if err := a.sqlDBPool.Close(); err != nil {
		return NewStorageError("CloseStorageConnections", "failed to close sql db pool", err)
	}
	return nil
}

func (a *App) createKubeEnv(cfg *config.Config) error {
	a.kubernetesCfg = &cfg.KubeCfg
	err := a.CreateKubeClientSet()
	return err
}

func (a *App) CreateKubeClientSet() error {
	kubeCfg, err := kube.CreateKubeconfig(*a.kubernetesCfg)
	if err != nil {
		return NewClientError("CreateKubeClientSet", "failed to create kubeconfig", err)
	}

	clientSet, err := kube.CreateKubeClientSet(kubeCfg)
	if err != nil {
		return NewClientError("CreateKubeClientSet", "failed to create kube clientset", err)
	}

	slog.Info("Creating APIClientSet")
	APIClientSet, err := apiextension.NewForConfig(kubeCfg)
	if err != nil {
		return NewClientError("CreateKubeClientSet", "failed to create api extension clientset", err)
	}

	a.kubernetesClient = clientSet
	a.kubernetesAPIClient = APIClientSet

	return nil
}
