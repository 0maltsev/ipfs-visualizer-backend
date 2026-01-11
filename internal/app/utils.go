package app

import (
	"database/sql"
	"ipfs-visualizer/config"
	psql_connection "ipfs-visualizer/internal/db/psql/connection"
	clustermodels "ipfs-visualizer/internal/db/psql/models/clusterModels"
	kubemodels "ipfs-visualizer/internal/db/psql/models/kubeModels"
	nodemodels "ipfs-visualizer/internal/db/psql/models/nodeModels"
	"ipfs-visualizer/internal/kube"
	"log/slog"

	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

func (a *App) loadGeneralCfg(cfg *config.Config) {
	a.serverCfg = &cfg.ServerCfg
	a.clusterCfg = &cfg.ClusterCfg
	a.nodeCfg = &cfg.NodeCfg
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
	clustermodels.CreateClustersTableIfNotExist(sqlPool)
	nodemodels.CreateNodesTableIfNotExist(sqlPool)
	kubemodels.CreateClusterKubeResourcesTableIfNotExist(sqlPool)
	kubemodels.CreateNodeKubeResourcesTableIfNotExist(sqlPool)
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
