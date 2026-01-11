package kubemodels

import (
	"context"
	"database/sql"

	sqlmodelerrors "ipfs-visualizer/internal/db/psql/models"
)

func CreateClusterKubeResourcesTableIfNotExist(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS cluster_kube_resources (
		cluster_id VARCHAR(255) PRIMARY KEY,

		namespace VARCHAR(255) NOT NULL,

		statefulset VARCHAR(255),
		service VARCHAR(255),

		env_configmap VARCHAR(255),
		scripts_configmap VARCHAR(255),

		cluster_secret VARCHAR(255),
		bootstrap_secret VARCHAR(255),

		ipfs_pvc VARCHAR(255),
		cluster_pvc VARCHAR(255),

		headless_service VARCHAR(255),

		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`
	if _, err := db.Exec(query); err != nil {
		return sqlmodelerrors.NewPostgresModelError(
			"CreateClusterKubeResourcesTableIfNotExist",
			"failed to create cluster_kube_resources table",
			err,
		)
	}
	return nil
}

func GetAllClusterKubeResources(
	ctx context.Context,
	db *sql.DB,
) ([]ClusterKubeResourcesModel, error) {

	var resources []ClusterKubeResourcesModel

	rows, err := db.QueryContext(ctx, getAllClusterKubeResourcesQuery)
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError(
			"GetAllClusterKubeResources",
			"query failed",
			err,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var r ClusterKubeResourcesModel
		if err := rows.Scan(
			&r.ClusterID,
			&r.Namespace,
			&r.StatefulSet,
			&r.Service,
			&r.EnvConfigMap,
			&r.ScriptsConfigMap,
			&r.ClusterSecret,
			&r.BootstrapSecret,
			&r.IPFSPVC,
			&r.ClusterPVC,
			&r.HeadlessService,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, sqlmodelerrors.NewPostgresModelError(
				"GetAllClusterKubeResources",
				"scan failed",
				err,
			)
		}
		resources = append(resources, r)
	}

	return resources, nil
}

func GetClusterKubeResourcesByClusterID(
	ctx context.Context,
	db *sql.DB,
	clusterID string,
) (*ClusterKubeResourcesModel, error) {

	var r ClusterKubeResourcesModel

	err := db.QueryRowContext(
		ctx,
		getClusterKubeResourcesByClusterIDQuery,
		clusterID,
	).Scan(
		&r.ClusterID,
		&r.Namespace,
		&r.StatefulSet,
		&r.Service,
		&r.EnvConfigMap,
		&r.ScriptsConfigMap,
		&r.ClusterSecret,
		&r.BootstrapSecret,
		&r.IPFSPVC,
		&r.ClusterPVC,
		&r.HeadlessService,
		&r.CreatedAt,
		&r.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError(
			"GetClusterKubeResourcesByClusterID",
			"query failed",
			err,
		)
	}

	return &r, nil
}

func InsertClusterKubeResources(
	ctx context.Context,
	db *sql.DB,
	r *ClusterKubeResourcesModel,
) error {

	if err := db.QueryRowContext(
		ctx,
		insertClusterKubeResourcesQuery,
		r.ClusterID,
		r.Namespace,
		r.StatefulSet,
		r.Service,
		r.EnvConfigMap,
		r.ScriptsConfigMap,
		r.ClusterSecret,
		r.BootstrapSecret,
		r.IPFSPVC,
		r.ClusterPVC,
		r.HeadlessService,
	).Scan(&r.CreatedAt, &r.UpdatedAt); err != nil {

		return sqlmodelerrors.NewPostgresModelError(
			"InsertClusterKubeResources",
			"insert failed",
			err,
		)
	}

	return nil
}

func UpdateClusterKubeResources(
	ctx context.Context,
	db *sql.DB,
	r *ClusterKubeResourcesModel,
) error {

	if err := db.QueryRowContext(
		ctx,
		updateClusterKubeResourcesQuery,
		r.Namespace,
		r.StatefulSet,
		r.Service,
		r.EnvConfigMap,
		r.ScriptsConfigMap,
		r.ClusterSecret,
		r.BootstrapSecret,
		r.IPFSPVC,
		r.ClusterPVC,
		r.HeadlessService,
		r.ClusterID,
	).Scan(&r.UpdatedAt); err != nil {

		return sqlmodelerrors.NewPostgresModelError(
			"UpdateClusterKubeResources",
			"update failed",
			err,
		)
	}

	return nil
}

func DeleteClusterKubeResourcesByClusterID(
	ctx context.Context,
	db *sql.DB,
	clusterID string,
) error {

	if _, err := db.ExecContext(
		ctx,
		deleteClusterKubeResourcesQuery,
		clusterID,
	); err != nil {

		return sqlmodelerrors.NewPostgresModelError(
			"DeleteClusterKubeResourcesByClusterID",
			"delete failed",
			err,
		)
	}

	return nil
}

func CreateNodeKubeResourcesTableIfNotExist(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS node_kube_resources (
		node_id VARCHAR(255) PRIMARY KEY,
		node_name VARCHAR(255) NOT NULL,
		cluster_id VARCHAR(255) NOT NULL,
		namespace VARCHAR(255) NOT NULL,
		pod_name VARCHAR(255),
		containers TEXT,
		service TEXT,
		configmap TEXT,
		secret TEXT,
		pvcs TEXT,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`
	if _, err := db.Exec(query); err != nil {
		return sqlmodelerrors.NewPostgresModelError(
			"CreateNodeKubeResourcesTableIfNotExist",
			"failed to create node_kube_resources table",
			err,
		)
	}
	return nil
}

func GetAllNodeKubeResources(ctx context.Context, db *sql.DB) ([]NodeKubeResourcesModel, error) {
	var resources []NodeKubeResourcesModel

	rows, err := db.QueryContext(ctx, getAllNodeKubeResourcesQuery)
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetAllNodeKubeResources", "query failed", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r NodeKubeResourcesModel
		if err := rows.Scan(
			&r.NodeID,
			&r.NodeName,
			&r.ClusterID,
			&r.Namespace,
			&r.PodName,
			&r.Containers,
			&r.Service,
			&r.ConfigMap,
			&r.Secret,
			&r.PVCs,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, sqlmodelerrors.NewPostgresModelError("GetAllNodeKubeResources", "scan failed", err)
		}
		resources = append(resources, r)
	}

	return resources, nil
}

func GetNodeKubeResourcesByNodeID(ctx context.Context, db *sql.DB, nodeID string) (*NodeKubeResourcesModel, error) {
	var r NodeKubeResourcesModel

	err := db.QueryRowContext(ctx, getNodeKubeResourcesByNodeIDQuery, nodeID).Scan(
		&r.NodeID,
		&r.NodeName,
		&r.ClusterID,
		&r.Namespace,
		&r.PodName,
		&r.Containers,
		&r.Service,
		&r.ConfigMap,
		&r.Secret,
		&r.PVCs,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetNodeKubeResourcesByNodeID", "query failed", err)
	}

	return &r, nil
}

func InsertNodeKubeResources(ctx context.Context, db *sql.DB, r *NodeKubeResourcesModel) error {
	if err := db.QueryRowContext(
		ctx,
		insertNodeKubeResourcesQuery,
		r.NodeID,
		r.NodeName,
		r.ClusterID,
		r.Namespace,
		r.PodName,
		r.Containers,
		r.Service,
		r.ConfigMap,
		r.Secret,
		r.PVCs,
	).Scan(&r.CreatedAt, &r.UpdatedAt); err != nil {
		return sqlmodelerrors.NewPostgresModelError("InsertNodeKubeResources", "insert failed", err)
	}
	return nil
}

func UpdateNodeKubeResources(ctx context.Context, db *sql.DB, r *NodeKubeResourcesModel) error {
	if err := db.QueryRowContext(
		ctx,
		updateNodeKubeResourcesQuery,
		r.NodeName,
		r.ClusterID,
		r.Namespace,
		r.PodName,
		r.Containers,
		r.Service,
		r.ConfigMap,
		r.Secret,
		r.PVCs,
		r.NodeID,
	).Scan(&r.UpdatedAt); err != nil {
		return sqlmodelerrors.NewPostgresModelError("UpdateNodeKubeResources", "update failed", err)
	}
	return nil
}

func DeleteNodeKubeResourcesByNodeID(ctx context.Context, db *sql.DB, nodeID string) error {
	if _, err := db.ExecContext(ctx, deleteNodeKubeResourcesQuery, nodeID); err != nil {
		return sqlmodelerrors.NewPostgresModelError("DeleteNodeKubeResourcesByNodeID", "delete failed", err)
	}
	return nil
}
