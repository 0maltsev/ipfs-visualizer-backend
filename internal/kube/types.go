package kube

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type ClusterKubeResources struct {
	ClusterID          string
	ClusterName        string
	Namespace          string
	ServiceType        string
	StorageClass       string
	ClusterStorageSize string
	IPFSStorageSize    string
	IPFSClusterImage   string
	IPFSImage          string
	EnvConfig          string
	ScriptsConfig      string

	StatefulSet *appsv1.StatefulSet

	HeadlessService *corev1.Service
	ExternalService *corev1.Service

	EnvConfigMap     *corev1.ConfigMap
	ScriptsConfigMap *corev1.ConfigMap

	ClusterSecret *corev1.Secret

	IPFSPVC    *corev1.PersistentVolumeClaim
	ClusterPVC *corev1.PersistentVolumeClaim

	Labels      map[string]string
	Annotations map[string]string
}

type NodeKubeResources struct {
	NodeID    string            `json:"nodeId"`
	NodeName  string            `json:"nodeName"`
	ClusterID string            `json:"clusterId"`
	Namespace string            `json:"namespace"`
	Env       map[string]string `json:"env"`
	Ports     PortsSpec         `json:"ports"`

	PodName    string             `json:"podName,omitempty"`
	Containers []corev1.Container `json:"containers,omitempty"`

	Service   *corev1.Service                 `json:"service,omitempty"`
	ConfigMap *corev1.ConfigMap               `json:"configMap,omitempty"`
	Secret    *corev1.Secret                  `json:"secret,omitempty"`
	PVCs      []*corev1.PersistentVolumeClaim `json:"pvcs,omitempty"`

	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type PortsSpec struct {
	SwarmTCP     int `json:"swarmTCP"`
	SwarmUDP     int `json:"swarmUDP"`
	API          int `json:"api"`
	HTTPGateway  int `json:"httpGateway"`
	WS           int `json:"ws"`
	ClusterAPI   int `json:"clusterAPI"`
	ClusterProxy int `json:"clusterProxy"`
	ClusterSwarm int `json:"clusterSwarm"`
}
