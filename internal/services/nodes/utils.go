package nodes

import (
	"encoding/json"
	"ipfs-visualizer/config"
	kubemodels "ipfs-visualizer/internal/db/psql/models/kubeModels"
	nodemodels "ipfs-visualizer/internal/db/psql/models/nodeModels"
	"ipfs-visualizer/internal/kube"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/uuid"
)

func ConvertNodeSqlModelToNodeSpec(sqlModel nodemodels.NodeSqlModel) NodeSpec {
	node := NodeSpec{
		NodeID: sqlModel.NodeID,
		Role:   sqlModel.Role,
		ClusterName: sqlModel.ClusterName,
		Ports: PortsSpec{
			SwarmTCP:     sqlModel.SwarmTCP,
			SwarmUDP:     sqlModel.SwarmUDP,
			API:          sqlModel.API,
			HTTPGateway:  sqlModel.HTTPGateway,
			WS:           sqlModel.WS,
			ClusterAPI:   sqlModel.ClusterAPI,
			ClusterProxy: sqlModel.ClusterProxy,
			ClusterSwarm: sqlModel.ClusterSwarm,
		},
		Env:        make(map[string]string),
		Containers: make([]ContainerSpec, 0),
	}

	if sqlModel.NodeName != nil {
		node.NodeName = *sqlModel.NodeName
	} else {
		node.NodeName = ""
	}

	if sqlModel.IPFSStorage != nil {
		node.Volumes.IPFSStorage = *sqlModel.IPFSStorage
	} else {
		node.Volumes.IPFSStorage = ""
	}

	if sqlModel.ClusterStorage != nil {
		node.Volumes.ClusterStorage = *sqlModel.ClusterStorage
	} else {
		node.Volumes.ClusterStorage = ""
	}

	if sqlModel.ScriptsConfig != nil {
		node.Volumes.ScriptsConfig = *sqlModel.ScriptsConfig
	} else {
		node.Volumes.ScriptsConfig = ""
	}

	return node
}

func ConvertNodeSpecToNodeSqlModel(node NodeSpec) nodemodels.NodeSqlModel {

	if node.NodeID == "" {
		node.NodeID = uuid.NewString()
	}

	sqlModel := nodemodels.NodeSqlModel{
		NodeID: node.NodeID,
		Role:   node.Role,
		ClusterName: node.ClusterName,

		SwarmTCP:     node.Ports.SwarmTCP,
		SwarmUDP:     node.Ports.SwarmUDP,
		API:          node.Ports.API,
		HTTPGateway:  node.Ports.HTTPGateway,
		WS:           node.Ports.WS,
		ClusterAPI:   node.Ports.ClusterAPI,
		ClusterProxy: node.Ports.ClusterProxy,
		ClusterSwarm: node.Ports.ClusterSwarm,
	}

	if node.NodeName != "" {
		sqlModel.NodeName = &node.NodeName
	}

	if node.Volumes.IPFSStorage != "" {
		sqlModel.IPFSStorage = &node.Volumes.IPFSStorage
	}

	if node.Volumes.ClusterStorage != "" {
		sqlModel.ClusterStorage = &node.Volumes.ClusterStorage
	}

	if node.Volumes.ScriptsConfig != "" {
		sqlModel.ScriptsConfig = &node.Volumes.ScriptsConfig
	}

	return sqlModel
}

// TODO не все поля заполняются
func BuildNodeSpecFromRole(nodeRole string, nodeCfg *config.NodeConfig) NodeSpec {
	return NodeSpec{
		NodeName: "NodeIPFS",
		Role: nodeRole,
		Ports: PortsSpec{
			SwarmTCP: nodeCfg.SwarmTCP,
			SwarmUDP: nodeCfg.SwarmUDP,
			API: nodeCfg.API,
			HTTPGateway: nodeCfg.HTTPGateway,
			WS: nodeCfg.WS,
			ClusterAPI: nodeCfg.ClusterAPI,
			ClusterProxy: nodeCfg.ClusterProxy,
			ClusterSwarm: nodeCfg.ClusterSwarm,
		},
		InitContainer: InitContainerSpec{
			Name: nodeCfg.InitContName,
			Image: nodeCfg.InitContImage,
			Command: nodeCfg.InitContCommand,
		},
	}
}

func ConvertNodeSpecToNodeKubeResources(
	node NodeSpec,
	clusterID string,
	namespace string,
) *kube.NodeKubeResources {

	res := &kube.NodeKubeResources{
		NodeID:    node.NodeID,
		NodeName:  node.NodeName,
		ClusterID: clusterID,
		Namespace: namespace,

		Env:   node.Env,
		Ports: kube.PortsSpec(node.Ports),

		Labels: map[string]string{
			"node-id":    node.NodeID,
			"node-name":  node.NodeName,
			"cluster-id": clusterID,
			"role":       node.Role,
		},
		Annotations: map[string]string{},
	}

	// Containers
	for _, c := range node.Containers {
		container := corev1.Container{
			Name:    c.Name,
			Image:   c.Image,
			Command: c.Command,
		}
		res.Containers = append(res.Containers, container)
	}

	return res
}

func ConvertNodeKubeResourcesToNodeSpec(
	res *kube.NodeKubeResources,
) NodeSpec {

	node := NodeSpec{
		NodeID:   res.NodeID,
		NodeName: res.NodeName,
		Role:     res.Labels["role"],

		Env:   res.Env,
		Ports: PortsSpec(res.Ports),
	}

	for _, c := range res.Containers {
		node.Containers = append(node.Containers, ContainerSpec{
			Name:    c.Name,
			Image:   c.Image,
			Command: c.Command,
		})
	}

	return node
}

func ConvertNodeKubeResourcesToModel(
	res *kube.NodeKubeResources,
) (*kubemodels.NodeKubeResourcesModel, error) {

	containersJSON, err := json.Marshal(res.Containers)
	if err != nil {
		return nil, err
	}

	pvcsJSON, err := json.Marshal(res.PVCs)
	if err != nil {
		return nil, err
	}

	model := &kubemodels.NodeKubeResourcesModel{
		NodeID:    res.NodeID,
		NodeName:  res.NodeName,
		ClusterID: res.ClusterID,
		Namespace: res.Namespace,
		PodName:   res.PodName,

		Containers: string(containersJSON),
		PVCs:       string(pvcsJSON),
	}

	if res.Service != nil {
		model.Service = res.Service.Name
	}

	if res.ConfigMap != nil {
		model.ConfigMap = res.ConfigMap.Name
	}

	if res.Secret != nil {
		model.Secret = res.Secret.Name
	}

	return model, nil
}

func ConvertNodeKubeResourcesModelToNodeKubeResources(
	model kubemodels.NodeKubeResourcesModel,
) (*kube.NodeKubeResources, error) {

	var containers []corev1.Container
	if model.Containers != "" {
		if err := json.Unmarshal([]byte(model.Containers), &containers); err != nil {
			return nil, err
		}
	}

	var pvcs []*corev1.PersistentVolumeClaim
	if model.PVCs != "" {
		if err := json.Unmarshal([]byte(model.PVCs), &pvcs); err != nil {
			return nil, err
		}
	}

	res := &kube.NodeKubeResources{
		NodeID:     model.NodeID,
		NodeName:   model.NodeName,
		ClusterID:  model.ClusterID,
		Namespace:  model.Namespace,
		PodName:    model.PodName,
		Containers: containers,
		PVCs:       pvcs,
	}

	if model.Service != "" {
		res.Service = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.Service,
				Namespace: model.Namespace,
			},
		}
	}

	if model.ConfigMap != "" {
		res.ConfigMap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.ConfigMap,
				Namespace: model.Namespace,
			},
		}
	}

	if model.Secret != "" {
		res.Secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.Secret,
				Namespace: model.Namespace,
			},
		}
	}

	return res, nil
}

