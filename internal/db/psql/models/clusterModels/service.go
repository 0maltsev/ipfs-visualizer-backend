package clustermodels

import (
	"context"
	"database/sql"
	sqlmodelerrors "ipfs-visualizer/internal/db/psql/models"
)

func CreateClustersTableIfNotExist(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS clusters (
		cluster_id VARCHAR(255) PRIMARY KEY,
		cluster_name VARCHAR(255),
		replicas INT NOT NULL,
		service_type VARCHAR(255),
		storage_class VARCHAR(255),
		cluster_storage_size VARCHAR(255),
		ipfs_storage_size VARCHAR(255),

		ipfs_image VARCHAR(255),
		ipfs_cluster_image VARCHAR(255),

		env_config VARCHAR(255),
		scripts_config VARCHAR(255),

		cluster_secret VARCHAR(255),
		bootstrap_priv_key VARCHAR(255),
		bootstrap_peer_id VARCHAR(255),

		nodes JSONB NOT NULL,

		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`
	if _, err := db.Exec(query); err != nil {
		return sqlmodelerrors.NewPostgresModelError(
			"CreateClustersTableIfNotExist",
			"failed to create clusters table",
			err,
		)
	}
	return nil
}

func GetAllClusters(ctx context.Context, db *sql.DB) ([]ClusterSqlModel, error) {
	var clusters []ClusterSqlModel

	rows, err := db.QueryContext(ctx, getAllClustersQuery)
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetAllClusters", "query failed", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c ClusterSqlModel
		if err := rows.Scan(
			&c.ClusterID,
			&c.ClusterName,
			&c.Replicas,
			&c.ServiceType,
			&c.StorageClass,
			&c.ClusterStorageSize,
			&c.IPFSStorageSize,
			&c.IPFSImage,
			&c.IPFSClusterImage,
			&c.EnvConfig,
			&c.ScriptsConfig,
			&c.ClusterSecret,
			&c.BootstrapPrivKey,
			&c.BootstrapPeerID,
			&c.NodeIDs,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, sqlmodelerrors.NewPostgresModelError("GetAllClusters", "scan failed", err)
		}
		clusters = append(clusters, c)
	}

	return clusters, nil
}

func GetClusterByID(ctx context.Context, db *sql.DB, id string) (*ClusterSqlModel, error) {
	var c ClusterSqlModel
	err := db.QueryRowContext(ctx, getClusterByIDQuery, id).Scan(
		&c.ClusterID,
		&c.ClusterName,
		&c.Replicas,
		&c.ServiceType,
		&c.StorageClass,
		&c.ClusterStorageSize,
		&c.IPFSStorageSize,
		&c.IPFSImage,
		&c.IPFSClusterImage,
		&c.EnvConfig,
		&c.ScriptsConfig,
		&c.ClusterSecret,
		&c.BootstrapPrivKey,
		&c.BootstrapPeerID,
		&c.NodeIDs,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetClusterByID", "query failed", err)
	}
	return &c, nil
}

func InsertCluster(ctx context.Context, db *sql.DB, c *ClusterSqlModel) error {
	if err := db.QueryRowContext(
		ctx,
		insertClusterQuery,
		c.ClusterID,
		c.ClusterName,
		c.Replicas,
		c.ServiceType,
		c.StorageClass,
		c.ClusterStorageSize,
		c.IPFSStorageSize,
		c.IPFSImage,
		c.IPFSClusterImage,
		c.EnvConfig,
		c.ScriptsConfig,
		c.ClusterSecret,
		c.BootstrapPrivKey,
		c.BootstrapPeerID,
		c.NodeIDs,
	).Scan(&c.CreatedAt, &c.UpdatedAt); err != nil {
		return sqlmodelerrors.NewPostgresModelError("InsertCluster", "insert failed", err)
	}
	return nil
}

func UpdateCluster(ctx context.Context, db *sql.DB, c *ClusterSqlModel) error {
	if err := db.QueryRowContext(
		ctx,
		updateClusterQuery,
		c.ClusterName,
		c.Replicas,
		c.ServiceType,
		c.StorageClass,
		c.ClusterStorageSize,
		c.IPFSStorageSize,
		c.IPFSImage,
		c.IPFSClusterImage,
		c.EnvConfig,
		c.ScriptsConfig,
		c.ClusterSecret,
		c.BootstrapPrivKey,
		c.BootstrapPeerID,
		c.NodeIDs,
		c.ClusterID,
	).Scan(&c.UpdatedAt); err != nil {
		return sqlmodelerrors.NewPostgresModelError("UpdateCluster", "update failed", err)
	}
	return nil
}

func DeleteClusterByID(ctx context.Context, db *sql.DB, id string) error {
	if _, err := db.ExecContext(ctx, deleteClusterQuery, id); err != nil {
		return sqlmodelerrors.NewPostgresModelError("DeleteClusterByID", "delete failed", err)
	}
	return nil
}

