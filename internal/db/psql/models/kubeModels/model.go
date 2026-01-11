package kubemodels

import "time"

type ClusterKubeResourcesModel struct {
	ClusterID        string    `db:"cluster_id" json:"clusterId"`
	Namespace        string    `db:"namespace" json:"namespace"`
	StatefulSet      string    `db:"statefulset,omitempty" json:"statefulSet,omitempty"`
	Service          string    `db:"service,omitempty" json:"service,omitempty"`
	EnvConfigMap     string    `db:"env_configmap,omitempty" json:"envConfigMap,omitempty"`
	ScriptsConfigMap string    `db:"scripts_configmap,omitempty" json:"scriptsConfigMap,omitempty"`
	ClusterSecret    string    `db:"cluster_secret,omitempty" json:"clusterSecret,omitempty"`
	BootstrapSecret  string    `db:"bootstrap_secret,omitempty" json:"bootstrapSecret,omitempty"`
	IPFSPVC          string    `db:"ipfs_pvc,omitempty" json:"ipfsPvc,omitempty"`
	ClusterPVC       string    `db:"cluster_pvc,omitempty" json:"clusterPvc,omitempty"`
	HeadlessService  string    `db:"headless_service,omitempty" json:"headlessService,omitempty"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

type NodeKubeResourcesModel struct {
	NodeID     string    `db:"node_id" json:"nodeId"`
	NodeName   string    `db:"node_name" json:"nodeName"`
	ClusterID  string    `db:"cluster_id" json:"clusterId"`
	Namespace  string    `db:"namespace" json:"namespace"`
	PodName    string    `db:"pod_name,omitempty" json:"podName,omitempty"`

	Containers   string `db:"containers,omitempty" json:"containers,omitempty"` // JSON сериализация контейнеров
	Service      string `db:"service,omitempty" json:"service,omitempty"`
	ConfigMap    string `db:"configmap,omitempty" json:"configMap,omitempty"`
	Secret       string `db:"secret,omitempty" json:"secret,omitempty"`
	PVCs         string `db:"pvcs,omitempty" json:"pvcs,omitempty"` // JSON сериализация PVC

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
