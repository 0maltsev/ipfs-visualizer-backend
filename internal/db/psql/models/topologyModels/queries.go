package topologymodels

const (
	createTopologiesTable = `
		CREATE TABLE IF NOT EXISTS topologies (
			topology_id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			deploy_status VARCHAR(50) DEFAULT 'none',
			k8s_namespace VARCHAR(255),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`

	createTopologyNodesTable = `
		CREATE TABLE IF NOT EXISTS topology_nodes (
			topology_id VARCHAR(255) NOT NULL REFERENCES topologies(topology_id) ON DELETE CASCADE,
			node_id VARCHAR(255) NOT NULL,
			label VARCHAR(255) NOT NULL,
			pos_x DOUBLE PRECISION NOT NULL,
			pos_y DOUBLE PRECISION NOT NULL,
			role VARCHAR(50) DEFAULT 'worker',
			PRIMARY KEY (topology_id, node_id)
		);`

	createTopologyEdgesTable = `
		CREATE TABLE IF NOT EXISTS topology_edges (
			topology_id VARCHAR(255) NOT NULL REFERENCES topologies(topology_id) ON DELETE CASCADE,
			edge_id VARCHAR(255) NOT NULL,
			source_node_id VARCHAR(255) NOT NULL,
			target_node_id VARCHAR(255) NOT NULL,
			PRIMARY KEY (topology_id, edge_id)
		);`

	getAllTopologiesQuery = `
		SELECT t.topology_id, t.name, t.deploy_status, t.k8s_namespace, t.created_at,
		       COALESCE((SELECT COUNT(*)::int FROM topology_nodes WHERE topology_id = t.topology_id), 0) AS node_count,
		       COALESCE((SELECT COUNT(*)::int FROM topology_edges WHERE topology_id = t.topology_id), 0) AS edge_count
		FROM topologies t
		ORDER BY t.updated_at DESC;`

	getTopologyByIDQuery = `
		SELECT topology_id, name, deploy_status, k8s_namespace, created_at, updated_at
		FROM topologies WHERE topology_id = $1;`

	insertTopologyQuery = `
		INSERT INTO topologies (topology_id, name, deploy_status, k8s_namespace)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at;`

	updateTopologyQuery = `
		UPDATE topologies SET name = $1, deploy_status = $2, k8s_namespace = $3, updated_at = NOW()
		WHERE topology_id = $4
		RETURNING updated_at;`

	deleteTopologyQuery = `DELETE FROM topologies WHERE topology_id = $1;`

	getNodesByTopologyQuery = `
		SELECT topology_id, node_id, label, pos_x, pos_y, role
		FROM topology_nodes WHERE topology_id = $1;`

	getEdgesByTopologyQuery = `
		SELECT topology_id, edge_id, source_node_id, target_node_id
		FROM topology_edges WHERE topology_id = $1;`

	insertNodeQuery = `
		INSERT INTO topology_nodes (topology_id, node_id, label, pos_x, pos_y, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (topology_id, node_id) DO UPDATE SET label = $3, pos_x = $4, pos_y = $5, role = $6;`

	insertEdgeQuery = `
		INSERT INTO topology_edges (topology_id, edge_id, source_node_id, target_node_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (topology_id, edge_id) DO UPDATE SET source_node_id = $3, target_node_id = $4;`

	deleteNodesByTopologyQuery = `DELETE FROM topology_nodes WHERE topology_id = $1;`
	deleteEdgesByTopologyQuery = `DELETE FROM topology_edges WHERE topology_id = $1;`
)
