package clustermodels

const (
	getAllClustersQuery = `
		SELECT
			cluster_id, cluster_name, replicas, service_type, storage_class,
			cluster_storage_size, ipfs_storage_size,
			ipfs_image, ipfs_cluster_image,
			env_config, scripts_config,
			cluster_secret, bootstrap_priv_key, bootstrap_peer_id,
			nodes, created_at, updated_at
		FROM clusters;
	`

	getClusterByIDQuery = `
		SELECT
			cluster_id, cluster_name, replicas, service_type, storage_class,
			cluster_storage_size, ipfs_storage_size,
			ipfs_image, ipfs_cluster_image,
			env_config, scripts_config,
			cluster_secret, bootstrap_priv_key, bootstrap_peer_id,
			nodes, created_at, updated_at
		FROM clusters
		WHERE cluster_id = $1;
	`

	insertClusterQuery = `
		INSERT INTO clusters (
			cluster_id, cluster_name, replicas, service_type, storage_class,
			cluster_storage_size, ipfs_storage_size,
			ipfs_image, ipfs_cluster_image,
			env_config, scripts_config,
			cluster_secret, bootstrap_priv_key, bootstrap_peer_id,
			nodes
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING created_at, updated_at;
	`

	updateClusterQuery = `
		UPDATE clusters SET
			cluster_name=$1,
			replicas=$2,
			service_type=$3,
			storage_class=$4,
			cluster_storage_size=$5,
			ipfs_storage_size=$6,
			ipfs_image=$7,
			ipfs_cluster_image=$8,
			env_config=$9,
			scripts_config=$10,
			cluster_secret=$11,
			bootstrap_priv_key=$12,
			bootstrap_peer_id=$13,
			nodes=$14,
			updated_at=NOW()
		WHERE cluster_id=$15
		RETURNING updated_at;
	`

	deleteClusterQuery = `
		DELETE FROM clusters WHERE cluster_id=$1;
	`
)
