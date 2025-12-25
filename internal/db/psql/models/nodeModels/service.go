package nodemodels

import (
	"context"
	"database/sql"
	sqlmodelerrors "ipfs-visualizer/internal/db/psql/models"
)

func CreateNodesTableIfNotExist(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS nodes (
		node_id VARCHAR(255) PRIMARY KEY,
		node_name VARCHAR(255),
		role VARCHAR(50) NOT NULL,

		swarm_tcp INT NOT NULL,
		swarm_udp INT NOT NULL,
		api INT NOT NULL,
		http_gateway INT NOT NULL,
		ws INT NOT NULL,
		cluster_api INT NOT NULL,
		cluster_proxy INT NOT NULL,
		cluster_swarm INT NOT NULL,

		ipfs_storage VARCHAR(255),
		cluster_storage VARCHAR(255),
		scripts_config VARCHAR(255),

		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`
	if _, err := db.Exec(query); err != nil {
		return sqlmodelerrors.NewPostgresModelError(
			"CreateNodesTableIfNotExist",
			"failed to create nodes table",
			err,
		)
	}
	return nil
}

func GetAllNodes(ctx context.Context, db *sql.DB) ([]NodeSqlModel, error) {
	var nodes []NodeSqlModel

	rows, err := db.QueryContext(ctx, getAllNodesQuery)
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetAllNodes", "query failed", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n NodeSqlModel
		if err := rows.Scan(
			&n.NodeID,
			&n.NodeName,
			&n.Role,
			&n.SwarmTCP,
			&n.SwarmUDP,
			&n.API,
			&n.HTTPGateway,
			&n.WS,
			&n.ClusterAPI,
			&n.ClusterProxy,
			&n.ClusterSwarm,
			&n.IPFSStorage,
			&n.ClusterStorage,
			&n.ScriptsConfig,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, sqlmodelerrors.NewPostgresModelError("GetAllNodes", "scan failed", err)
		}
		nodes = append(nodes, n)
	}

	return nodes, nil
}

func GetNodeByID(ctx context.Context, db *sql.DB, id string) (*NodeSqlModel, error) {
	var n NodeSqlModel
	err := db.QueryRowContext(ctx, getNodeByIDQuery, id).Scan(
		&n.NodeID,
		&n.NodeName,
		&n.Role,
		&n.SwarmTCP,
		&n.SwarmUDP,
		&n.API,
		&n.HTTPGateway,
		&n.WS,
		&n.ClusterAPI,
		&n.ClusterProxy,
		&n.ClusterSwarm,
		&n.IPFSStorage,
		&n.ClusterStorage,
		&n.ScriptsConfig,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetNodeByID", "query failed", err)
	}
	return &n, nil
}

func InsertNode(ctx context.Context, db *sql.DB, n *NodeSqlModel) error {
	if err := db.QueryRowContext(
		ctx,
		insertNodeQuery,
		n.NodeID,
		n.NodeName,
		n.Role,
		n.SwarmTCP,
		n.SwarmUDP,
		n.API,
		n.HTTPGateway,
		n.WS,
		n.ClusterAPI,
		n.ClusterProxy,
		n.ClusterSwarm,
		n.IPFSStorage,
		n.ClusterStorage,
		n.ScriptsConfig,
	).Scan(&n.CreatedAt, &n.UpdatedAt); err != nil {
		return sqlmodelerrors.NewPostgresModelError("InsertNode", "insert failed", err)
	}
	return nil
}

func UpdateNode(ctx context.Context, db *sql.DB, n *NodeSqlModel) error {
	if err := db.QueryRowContext(
		ctx,
		updateNodeQuery,
		n.NodeName,
		n.Role,
		n.SwarmTCP,
		n.SwarmUDP,
		n.API,
		n.HTTPGateway,
		n.WS,
		n.ClusterAPI,
		n.ClusterProxy,
		n.ClusterSwarm,
		n.IPFSStorage,
		n.ClusterStorage,
		n.ScriptsConfig,
		n.NodeID,
	).Scan(&n.UpdatedAt); err != nil {
		return sqlmodelerrors.NewPostgresModelError("UpdateNode", "update failed", err)
	}
	return nil
}

func DeleteNodeByID(ctx context.Context, db *sql.DB, id string) error {
	if _, err := db.ExecContext(ctx, deleteNodeQuery, id); err != nil {
		return sqlmodelerrors.NewPostgresModelError("DeleteNodeByID", "delete failed", err)
	}
	return nil
}

