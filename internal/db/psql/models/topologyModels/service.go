package topologymodels

import (
	"context"
	"database/sql"

	sqlmodelerrors "ipfs-visualizer/internal/db/psql/models"
)

func CreateTopologyTablesIfNotExist(db *sql.DB) error {
	if _, err := db.Exec(createTopologiesTable); err != nil {
		return sqlmodelerrors.NewPostgresModelError("CreateTopologyTablesIfNotExist", "failed to create topologies table", err)
	}
	if _, err := db.Exec(createTopologyNodesTable); err != nil {
		return sqlmodelerrors.NewPostgresModelError("CreateTopologyTablesIfNotExist", "failed to create topology_nodes table", err)
	}
	if _, err := db.Exec(createTopologyEdgesTable); err != nil {
		return sqlmodelerrors.NewPostgresModelError("CreateTopologyTablesIfNotExist", "failed to create topology_edges table", err)
	}
	return nil
}

func GetAllTopologies(ctx context.Context, db *sql.DB) ([]TopologySummaryRow, error) {
	rows, err := db.QueryContext(ctx, getAllTopologiesQuery)
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetAllTopologies", "query failed", err)
	}
	defer rows.Close()

	var list []TopologySummaryRow
	for rows.Next() {
		var r TopologySummaryRow
		if err := rows.Scan(&r.TopologyID, &r.Name, &r.DeployStatus, &r.K8sNamespace, &r.CreatedAt, &r.NodeCount, &r.EdgeCount); err != nil {
			return nil, sqlmodelerrors.NewPostgresModelError("GetAllTopologies", "scan failed", err)
		}
		list = append(list, r)
	}
	return list, nil
}

func GetTopologyByID(ctx context.Context, db *sql.DB, id string) (*TopologyModel, error) {
	var m TopologyModel
	err := db.QueryRowContext(ctx, getTopologyByIDQuery, id).Scan(
		&m.TopologyID, &m.Name, &m.DeployStatus, &m.K8sNamespace, &m.CreatedAt, &m.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetTopologyByID", "query failed", err)
	}
	return &m, nil
}

func InsertTopology(ctx context.Context, db *sql.DB, m *TopologyModel) error {
	if err := db.QueryRowContext(ctx, insertTopologyQuery,
		m.TopologyID, m.Name, m.DeployStatus, m.K8sNamespace,
	).Scan(&m.CreatedAt, &m.UpdatedAt); err != nil {
		return sqlmodelerrors.NewPostgresModelError("InsertTopology", "insert failed", err)
	}
	return nil
}

func UpdateTopology(ctx context.Context, db *sql.DB, m *TopologyModel) error {
	_, err := db.ExecContext(ctx, updateTopologyQuery, m.Name, m.DeployStatus, m.K8sNamespace, m.TopologyID)
	return err
}

func UpdateTopologyDeployStatus(ctx context.Context, db *sql.DB, topologyID, status string, namespace *string) error {
	_, err := db.ExecContext(ctx, `UPDATE topologies SET deploy_status = $1, k8s_namespace = $2, updated_at = NOW() WHERE topology_id = $3`,
		status, namespace, topologyID)
	return err
}

func DeleteTopology(ctx context.Context, db *sql.DB, id string) error {
	_, err := db.ExecContext(ctx, deleteTopologyQuery, id)
	return err
}

func GetNodesByTopology(ctx context.Context, db *sql.DB, topologyID string) ([]TopologyNodeModel, error) {
	rows, err := db.QueryContext(ctx, getNodesByTopologyQuery, topologyID)
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetNodesByTopology", "query failed", err)
	}
	defer rows.Close()

	var list []TopologyNodeModel
	for rows.Next() {
		var n TopologyNodeModel
		if err := rows.Scan(&n.TopologyID, &n.NodeID, &n.Label, &n.PosX, &n.PosY, &n.Role); err != nil {
			return nil, sqlmodelerrors.NewPostgresModelError("GetNodesByTopology", "scan failed", err)
		}
		list = append(list, n)
	}
	return list, nil
}

func GetEdgesByTopology(ctx context.Context, db *sql.DB, topologyID string) ([]TopologyEdgeModel, error) {
	rows, err := db.QueryContext(ctx, getEdgesByTopologyQuery, topologyID)
	if err != nil {
		return nil, sqlmodelerrors.NewPostgresModelError("GetEdgesByTopology", "query failed", err)
	}
	defer rows.Close()

	var list []TopologyEdgeModel
	for rows.Next() {
		var e TopologyEdgeModel
		if err := rows.Scan(&e.TopologyID, &e.EdgeID, &e.SourceNodeID, &e.TargetNodeID); err != nil {
			return nil, sqlmodelerrors.NewPostgresModelError("GetEdgesByTopology", "scan failed", err)
		}
		list = append(list, e)
	}
	return list, nil
}

func ReplaceTopologyNodes(ctx context.Context, db *sql.DB, topologyID string, nodes []TopologyNodeModel) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, deleteNodesByTopologyQuery, topologyID); err != nil {
		return err
	}
	for _, n := range nodes {
		if _, err := tx.ExecContext(ctx, insertNodeQuery,
			n.TopologyID, n.NodeID, n.Label, n.PosX, n.PosY, n.Role); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func ReplaceTopologyEdges(ctx context.Context, db *sql.DB, topologyID string, edges []TopologyEdgeModel) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, deleteEdgesByTopologyQuery, topologyID); err != nil {
		return err
	}
	for _, e := range edges {
		if _, err := tx.ExecContext(ctx, insertEdgeQuery,
			e.TopologyID, e.EdgeID, e.SourceNodeID, e.TargetNodeID); err != nil {
			return err
		}
	}
	return tx.Commit()
}
