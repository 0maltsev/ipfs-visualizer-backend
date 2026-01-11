package kubemodels

const (
	getAllClusterKubeResourcesQuery = `
		SELECT
			cluster_id,
			namespace,
			statefulset,
			service,
			env_configmap,
			scripts_configmap,
			cluster_secret,
			bootstrap_secret,
			ipfs_pvc,
			cluster_pvc,
			headless_service,
			created_at,
			updated_at
		FROM cluster_kube_resources;
	`

	getClusterKubeResourcesByClusterIDQuery = `
		SELECT
			cluster_id,
			namespace,
			statefulset,
			service,
			env_configmap,
			scripts_configmap,
			cluster_secret,
			bootstrap_secret,
			ipfs_pvc,
			cluster_pvc,
			headless_service,
			created_at,
			updated_at
		FROM cluster_kube_resources
		WHERE cluster_id=$1;
	`

	insertClusterKubeResourcesQuery = `
		INSERT INTO cluster_kube_resources (
			cluster_id,
			namespace,
			statefulset,
			service,
			env_configmap,
			scripts_configmap,
			cluster_secret,
			bootstrap_secret,
			ipfs_pvc,
			cluster_pvc,
			headless_service
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING created_at, updated_at;
	`

	updateClusterKubeResourcesQuery = `
		UPDATE cluster_kube_resources
		SET
			namespace=$1,
			statefulset=$2,
			service=$3,
			env_configmap=$4,
			scripts_configmap=$5,
			cluster_secret=$6,
			bootstrap_secret=$7,
			ipfs_pvc=$8,
			cluster_pvc=$9,
			headless_service=$10,
			updated_at=NOW()
		WHERE cluster_id=$11
		RETURNING updated_at;
	`

	deleteClusterKubeResourcesQuery = `
		DELETE FROM cluster_kube_resources
		WHERE cluster_id=$1;
	`
)

const (
	getAllNodeKubeResourcesQuery = `
		SELECT
			node_id,
			node_name,
			cluster_id,
			namespace,
			pod_name,
			containers,
			service,
			configmap,
			secret,
			pvcs,
			created_at,
			updated_at
		FROM node_kube_resources;
	`

	getNodeKubeResourcesByNodeIDQuery = `
		SELECT
			node_id,
			node_name,
			cluster_id,
			namespace,
			pod_name,
			containers,
			service,
			configmap,
			secret,
			pvcs,
			created_at,
			updated_at
		FROM node_kube_resources
		WHERE node_id=$1;
	`

	insertNodeKubeResourcesQuery = `
		INSERT INTO node_kube_resources (
			node_id,
			node_name,
			cluster_id,
			namespace,
			pod_name,
			containers,
			service,
			configmap,
			secret,
			pvcs
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING created_at, updated_at;
	`

	updateNodeKubeResourcesQuery = `
		UPDATE node_kube_resources
		SET
			node_name=$1,
			cluster_id=$2,
			namespace=$3,
			pod_name=$4,
			containers=$5,
			service=$6,
			configmap=$7,
			secret=$8,
			pvcs=$9,
			updated_at=NOW()
		WHERE node_id=$10
		RETURNING updated_at;
	`

	deleteNodeKubeResourcesQuery = `
		DELETE FROM node_kube_resources
		WHERE node_id=$1;
	`
)
