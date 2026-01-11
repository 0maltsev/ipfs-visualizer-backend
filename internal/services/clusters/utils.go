package clusters

import (
	clustermodels "ipfs-visualizer/internal/db/psql/models/clusterModels"
	kubemodels "ipfs-visualizer/internal/db/psql/models/kubeModels"
	"ipfs-visualizer/internal/kube"
	"ipfs-visualizer/internal/services/nodes"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/uuid"
)

func ConvertClusterSqlModelToClusterSpec(
	sqlModel clustermodels.ClusterSqlModel,
) ClusterSpec {
	cluster := ClusterSpec{
		ClusterID: sqlModel.ClusterID,
		Replicas:  sqlModel.Replicas,
		Nodes:     make([]nodes.NodeSpec, 0),
	}

	if sqlModel.ClusterName != nil {
		cluster.ClusterName = *sqlModel.ClusterName
	}

	if sqlModel.ServiceType != nil {
		cluster.ServiceType = *sqlModel.ServiceType
	}

	if sqlModel.StorageClass != nil {
		cluster.StorageClass = *sqlModel.StorageClass
	}

	if sqlModel.ClusterStorageSize != nil {
		cluster.ClusterStorageSize = *sqlModel.ClusterStorageSize
	}

	if sqlModel.IPFSStorageSize != nil {
		cluster.IPFSStorageSize = *sqlModel.IPFSStorageSize
	}

	// Images
	if sqlModel.IPFSImage != nil {
		cluster.Images.IPFS = *sqlModel.IPFSImage
	}

	if sqlModel.IPFSClusterImage != nil {
		cluster.Images.IPFSCluster = *sqlModel.IPFSClusterImage
	}

	// ConfigMaps
	if sqlModel.EnvConfig != nil {
		cluster.ConfigMaps.EnvConfig = *sqlModel.EnvConfig
	}

	if sqlModel.ScriptsConfig != nil {
		cluster.ConfigMaps.ScriptsConfig = *sqlModel.ScriptsConfig
	}

	// Secrets
	if sqlModel.ClusterSecret != nil {
		cluster.Secrets.ClusterSecret = *sqlModel.ClusterSecret
	}

	if sqlModel.BootstrapPrivKey != nil {
		cluster.Secrets.BootstrapPrivKey = *sqlModel.BootstrapPrivKey
	}

	if sqlModel.BootstrapPeerID != nil {
		cluster.BootstrapPeerID = *sqlModel.BootstrapPeerID
	}

	return cluster
}

func ConvertClusterSpecToClusterSqlModel(
	cluster ClusterSpec,
) clustermodels.ClusterSqlModel {

	if cluster.ClusterID == "" {
		cluster.ClusterID = uuid.NewString()
	}

	sqlModel := clustermodels.ClusterSqlModel{
		ClusterID: cluster.ClusterID,
		Replicas:  cluster.Replicas,
		NodeIDs:   make([]string, 0),
	}

	if cluster.ClusterName != "" {
		sqlModel.ClusterName = &cluster.ClusterName
	}

	if cluster.ServiceType != "" {
		sqlModel.ServiceType = &cluster.ServiceType
	}

	if cluster.StorageClass != "" {
		sqlModel.StorageClass = &cluster.StorageClass
	}

	if cluster.ClusterStorageSize != "" {
		sqlModel.ClusterStorageSize = &cluster.ClusterStorageSize
	}

	if cluster.IPFSStorageSize != "" {
		sqlModel.IPFSStorageSize = &cluster.IPFSStorageSize
	}

	// Images
	if cluster.Images.IPFS != "" {
		sqlModel.IPFSImage = &cluster.Images.IPFS
	}

	if cluster.Images.IPFSCluster != "" {
		sqlModel.IPFSClusterImage = &cluster.Images.IPFSCluster
	}

	// ConfigMaps
	if cluster.ConfigMaps.EnvConfig != "" {
		sqlModel.EnvConfig = &cluster.ConfigMaps.EnvConfig
	}

	if cluster.ConfigMaps.ScriptsConfig != "" {
		sqlModel.ScriptsConfig = &cluster.ConfigMaps.ScriptsConfig
	}

	// Secrets
	if cluster.Secrets.ClusterSecret != "" {
		sqlModel.ClusterSecret = &cluster.Secrets.ClusterSecret
	}

	if cluster.Secrets.BootstrapPrivKey != "" {
		sqlModel.BootstrapPrivKey = &cluster.Secrets.BootstrapPrivKey
	}

	if cluster.BootstrapPeerID != "" {
		sqlModel.BootstrapPeerID = &cluster.BootstrapPeerID
	}

	// Nodes → сохраняем только IDs
	for _, node := range cluster.Nodes {
		sqlModel.NodeIDs = append(sqlModel.NodeIDs, node.NodeID)
	}

	return sqlModel
}


func removeByValue(s []string, val string) []string {
	for i, v := range s {
		if v == val {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func ConvertClusterSpecToClusterKubeResources(
	spec ClusterSpec,
	namespace string,
) *kube.ClusterKubeResources {

	res := &kube.ClusterKubeResources{
		ClusterID:   spec.ClusterID,
		ClusterName: spec.ClusterName,
		Namespace:   namespace,

		ServiceType:        spec.ServiceType,
		StorageClass:       spec.StorageClass,
		ClusterStorageSize: spec.ClusterStorageSize,
		IPFSStorageSize:    spec.IPFSStorageSize,

		IPFSImage:        spec.Images.IPFS,
		IPFSClusterImage: spec.Images.IPFSCluster,

		EnvConfig:     spec.ConfigMaps.EnvConfig,
		ScriptsConfig: spec.ConfigMaps.ScriptsConfig,

		Labels: map[string]string{
			"app":        spec.ClusterName,
			"cluster-id": spec.ClusterID,
		},
		Annotations: map[string]string{},
	}

	// secrets (если уже сгенерированы и лежат в spec)
	if spec.Secrets.ClusterSecret != "" {
		res.ClusterSecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: spec.ClusterName + "-cluster-secret",
			},
		}
	}


	return res
}

func ClusterKubeResourcesToConvertClusterSpec(
	res *kube.ClusterKubeResources,
) ClusterSpec {

	spec := ClusterSpec{
		ClusterID:          res.ClusterID,
		ClusterName:        res.ClusterName,
		ServiceType:        res.ServiceType,
		StorageClass:       res.StorageClass,
		ClusterStorageSize: res.ClusterStorageSize,
		IPFSStorageSize:    res.IPFSStorageSize,

		Images: ImagesSpec{
			IPFS:        res.IPFSImage,
			IPFSCluster: res.IPFSClusterImage,
		},

		ConfigMaps: ConfigMapsSpec{
			EnvConfig:     res.EnvConfig,
			ScriptsConfig: res.ScriptsConfig,
		},

		Secrets: SecretsSpec{},
		Nodes:   []nodes.NodeSpec{},
	}

	// Secrets (если присутствуют)
	if res.ClusterSecret != nil {
		spec.Secrets.ClusterSecret = res.ClusterSecret.Name
	}


	// Реплики — только если StatefulSet уже существует
	if res.StatefulSet != nil && res.StatefulSet.Spec.Replicas != nil {
		spec.Replicas = int(*res.StatefulSet.Spec.Replicas)
	}

	return spec
}

func ConvertClusterKubeResourcesToModel(
	res *kube.ClusterKubeResources,
) kubemodels.ClusterKubeResourcesModel {

	model := kubemodels.ClusterKubeResourcesModel{
		ClusterID: res.ClusterID,
		Namespace: res.Namespace,
	}

	if res.StatefulSet != nil {
		model.StatefulSet = res.StatefulSet.Name
	}

	if res.ExternalService != nil {
		model.Service = res.ExternalService.Name
	}

	if res.HeadlessService != nil {
		model.HeadlessService = res.HeadlessService.Name
	}

	if res.EnvConfigMap != nil {
		model.EnvConfigMap = res.EnvConfigMap.Name
	}

	if res.ScriptsConfigMap != nil {
		model.ScriptsConfigMap = res.ScriptsConfigMap.Name
	}

	if res.ClusterSecret != nil {
		model.ClusterSecret = res.ClusterSecret.Name
	}


	if res.IPFSPVC != nil {
		model.IPFSPVC = res.IPFSPVC.Name
	}

	if res.ClusterPVC != nil {
		model.ClusterPVC = res.ClusterPVC.Name
	}

	return model
}

func ConvertClusterKubeResourcesModelToClusterKubeResources(
	model kubemodels.ClusterKubeResourcesModel,
) *kube.ClusterKubeResources {

	res := &kube.ClusterKubeResources{
		ClusterID: model.ClusterID,
		Namespace: model.Namespace,
	}

	if model.StatefulSet != "" {
		res.StatefulSet = &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.StatefulSet,
				Namespace: model.Namespace,
			},
		}
	}

	if model.Service != "" {
		res.ExternalService = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.Service,
				Namespace: model.Namespace,
			},
		}
	}

	if model.HeadlessService != "" {
		res.HeadlessService = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.HeadlessService,
				Namespace: model.Namespace,
			},
		}
	}

	if model.EnvConfigMap != "" {
		res.EnvConfigMap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.EnvConfigMap,
				Namespace: model.Namespace,
			},
		}
	}

	if model.ScriptsConfigMap != "" {
		res.ScriptsConfigMap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.ScriptsConfigMap,
				Namespace: model.Namespace,
			},
		}
	}

	if model.ClusterSecret != "" {
		res.ClusterSecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.ClusterSecret,
				Namespace: model.Namespace,
			},
		}
	}


	if model.IPFSPVC != "" {
		res.IPFSPVC = &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.IPFSPVC,
				Namespace: model.Namespace,
			},
		}
	}

	if model.ClusterPVC != "" {
		res.ClusterPVC = &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.ClusterPVC,
				Namespace: model.Namespace,
			},
		}
	}

	return res
}

