package nodemodels

const (
	getAllNodesQuery = `
		SELECT
			node_id, node_name, role,
			swarm_tcp, swarm_udp, api, http_gateway, ws,
			cluster_api, cluster_proxy, cluster_swarm,
			ipfs_storage, cluster_storage, scripts_config,
			created_at, updated_at
		FROM nodes;
	`

	getNodeByIDQuery = `
		SELECT
			node_id, node_name, role,
			swarm_tcp, swarm_udp, api, http_gateway, ws,
			cluster_api, cluster_proxy, cluster_swarm,
			ipfs_storage, cluster_storage, scripts_config,
			created_at, updated_at
		FROM nodes
		WHERE node_id = $1;
	`

	insertNodeQuery = `
		INSERT INTO nodes (
			node_id, node_name, role,
			swarm_tcp, swarm_udp, api, http_gateway, ws,
			cluster_api, cluster_proxy, cluster_swarm,
			ipfs_storage, cluster_storage, scripts_config
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		RETURNING created_at, updated_at;
	`

	updateNodeQuery = `
		UPDATE nodes SET
			node_name=$1,
			role=$2,
			swarm_tcp=$3,
			swarm_udp=$4,
			api=$5,
			http_gateway=$6,
			ws=$7,
			cluster_api=$8,
			cluster_proxy=$9,
			cluster_swarm=$10,
			ipfs_storage=$11,
			cluster_storage=$12,
			scripts_config=$13,
			updated_at=NOW()
		WHERE node_id=$14
		RETURNING updated_at;
	`

	deleteNodeQuery = `
		DELETE FROM nodes WHERE node_id=$1;
	`
)
