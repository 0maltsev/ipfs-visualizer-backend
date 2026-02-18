package topology

import (
	"context"
	"fmt"
	"strings"

	"ipfs-visualizer/internal/kube"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type NodeInfo struct {
	NodeID   string
	Label    string
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	Role string
}

type EdgeInfo struct {
	SourceNodeID string
	TargetNodeID string
}

type DeployConfig struct {
	TopologyID  string
	Name        string
	Namespace   string
	Nodes       []NodeInfo
	Edges       []EdgeInfo
	BootstrapID string
}

func Deploy(ctx context.Context, client kubernetes.Interface, cfg DeployConfig) error {
	keyPair, err := kube.GenerateBootstrapPrivateKey()
	if err != nil {
		return fmt.Errorf("generate bootstrap key: %w", err)
	}
	clusterSecret, err := kube.GenerateClusterSecret()
	if err != nil {
		return fmt.Errorf("generate cluster secret: %w", err)
	}

	svcName := "ipfs-" + strings.ReplaceAll(cfg.TopologyID, "-", "")[:12]

	entrypointScript := getEntrypointScript(svcName, keyPair.PeerID)
	configureIPFSScript := getConfigureIPFSScript()

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: svcName + "-scripts", Namespace: cfg.Namespace},
		Data: map[string]string{
			"entrypoint.sh":     entrypointScript,
			"configure-ipfs.sh": configureIPFSScript,
		},
	}
	if _, err := client.CoreV1().ConfigMaps(cfg.Namespace).Create(ctx, cm, metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("create configmap: %w", err)
		}
	}

	envCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: svcName + "-env", Namespace: cfg.Namespace},
		Data: map[string]string{
			"bootstrap-peer-id": keyPair.PeerID,
		},
	}
	if _, err := client.CoreV1().ConfigMaps(cfg.Namespace).Create(ctx, envCM, metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("create env configmap: %w", err)
		}
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: svcName + "-secrets", Namespace: cfg.Namespace},
		Type:       corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"cluster-secret":          []byte(clusterSecret),
			"bootstrap-peer-priv-key": []byte(keyPair.PrivateKey),
		},
	}
	if _, err := client.CoreV1().Secrets(cfg.Namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("create secret: %w", err)
		}
	}

	headlessSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: svcName, Namespace: cfg.Namespace},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector:  map[string]string{"app": svcName},
			Ports: []corev1.ServicePort{
				{Name: "swarm", Port: 4001},
				{Name: "swarm-udp", Port: 4002, Protocol: corev1.ProtocolUDP},
				{Name: "api", Port: 5001},
				{Name: "ws", Port: 8081},
				{Name: "http", Port: 8080},
				{Name: "cluster-swarm", Port: 9096},
				{Name: "cluster-api", Port: 9094},
				{Name: "cluster-proxy", Port: 9095},
			},
		},
	}
	if _, err := client.CoreV1().Services(cfg.Namespace).Create(ctx, headlessSvc, metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("create headless service: %w", err)
		}
	}

	externalSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: svcName + "-external", Namespace: cfg.Namespace},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeLoadBalancer,
			Selector: map[string]string{"app": svcName},
			Ports: []corev1.ServicePort{
				{Name: "swarm", Port: 4001, TargetPort: intstr.FromInt(4001)},
				{Name: "swarm-udp", Port: 4002, TargetPort: intstr.FromInt(4002), Protocol: corev1.ProtocolUDP},
				{Name: "api", Port: 5001, TargetPort: intstr.FromInt(5001)},
				{Name: "http", Port: 8080, TargetPort: intstr.FromInt(8080)},
				{Name: "ws", Port: 8081, TargetPort: intstr.FromInt(8081)},
				{Name: "cluster-api", Port: 9094, TargetPort: intstr.FromInt(9094)},
				{Name: "cluster-proxy", Port: 9095, TargetPort: intstr.FromInt(9095)},
				{Name: "cluster-swarm", Port: 9096, TargetPort: intstr.FromInt(9096)},
			},
		},
	}
	if _, err := client.CoreV1().Services(cfg.Namespace).Create(ctx, externalSvc, metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("create external service: %w", err)
		}
	}

	replicas := int32(len(cfg.Nodes))
	if replicas < 1 {
		replicas = 1
	}

	sts := buildStatefulSet(svcName, cfg.Namespace, replicas, keyPair.PeerID)
	if _, err := client.AppsV1().StatefulSets(cfg.Namespace).Create(ctx, sts, metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("create statefulset: %w", err)
		}
	}

	return nil
}

func buildStatefulSet(svcName, namespace string, replicas int32, bootstrapPeerID string) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: svcName, Namespace: namespace},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: svcName,
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": svcName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": svcName}},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "configure-ipfs",
							Image: "ipfs/kubo:release",
							Command: []string{"sh", "/custom/configure-ipfs.sh"},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "ipfs-storage", MountPath: "/data/ipfs"},
								{Name: "configure-script", MountPath: "/custom"},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "ipfs",
							Image: "ipfs/kubo:release",
							Env:   []corev1.EnvVar{{Name: "IPFS_FD_MAX", Value: "4096"}},
							Ports: []corev1.ContainerPort{
								{Name: "swarm", ContainerPort: 4001, Protocol: corev1.ProtocolTCP},
								{Name: "swarm-udp", ContainerPort: 4002, Protocol: corev1.ProtocolUDP},
								{Name: "api", ContainerPort: 5001},
								{Name: "ws", ContainerPort: 8081},
								{Name: "http", ContainerPort: 8080},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "ipfs-storage", MountPath: "/data/ipfs"},
								{Name: "configure-script", MountPath: "/custom"},
							},
						},
						{
							Name:  "ipfs-cluster",
							Image: "ipfs/ipfs-cluster:latest",
							Command: []string{"sh", "/custom/entrypoint.sh"},
							Env: []corev1.EnvVar{
								{Name: "BOOTSTRAP_PEER_ID", ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{Name: svcName + "-env"},
										Key:                  "bootstrap-peer-id",
									},
								}},
								{Name: "BOOTSTRAP_PEER_PRIV_KEY", ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{Name: svcName + "-secrets"},
										Key:                  "bootstrap-peer-priv-key",
									},
								}},
								{Name: "CLUSTER_SECRET", ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{Name: svcName + "-secrets"},
										Key:                  "cluster-secret",
									},
								}},
								{Name: "CLUSTER_MONITOR_PING_INTERVAL", Value: "3m"},
								{Name: "SVC_NAME", Value: svcName},
							},
							Ports: []corev1.ContainerPort{
								{Name: "api-http", ContainerPort: 9094},
								{Name: "proxy-http", ContainerPort: 9095},
								{Name: "cluster-swarm", ContainerPort: 9096},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "cluster-storage", MountPath: "/data/ipfs-cluster"},
								{Name: "configure-script", MountPath: "/custom"},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "configure-script",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: svcName + "-scripts"},
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "cluster-storage"},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
						StorageClassName: strPtr("standard"),
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("30Gi"),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "ipfs-storage"},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
						StorageClassName: strPtr("standard"),
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("30Gi"),
							},
						},
					},
				},
			},
		},
	}
}

func strPtr(s string) *string { return &s }

func getEntrypointScript(svcName, bootstrapPeerID string) string {
	return `#!/bin/sh
user=ipfs
sleep 10
if [ ! -f /data/ipfs-cluster/service.json ]; then
  ipfs-cluster-service init
fi
sed -i 's~/ip4/127.0.0.1/tcp/9095~/ip4/0.0.0.0/tcp/9095~g' /data/ipfs-cluster/service.json
sed -i 's~/ip4/127.0.0.1/tcp/9094~/ip4/0.0.0.0/tcp/9094~g' /data/ipfs-cluster/service.json

if echo $(cat /proc/sys/kernel/hostname) | grep -q "` + svcName + `-0"; then
  CLUSTER_ID=${BOOTSTRAP_PEER_ID} \
  CLUSTER_PRIVATEKEY=${BOOTSTRAP_PEER_PRIV_KEY} \
  exec ipfs-cluster-service daemon --upgrade
else
  BOOTSTRAP_ADDR=/dns4/` + svcName + `-0.` + svcName + `/tcp/9096/ipfs/` + bootstrapPeerID + `
  exec ipfs-cluster-service daemon --upgrade --bootstrap $BOOTSTRAP_ADDR --leave
fi
`
}

func getConfigureIPFSScript() string {
	return `#!/bin/sh
set -e
set -x
user=root
mkdir -p /data/ipfs && chown -R ipfs /data/ipfs
user=ipfs
if [ -f /data/ipfs/config ]; then
  if [ -f /data/ipfs/repo.lock ]; then
    rm /data/ipfs/repo.lock
  fi
  exit 0
fi
ipfs init --profile=badgerds,server
ipfs config --json Addresses.API /ip4/0.0.0.0/tcp/5001
ipfs config --json Addresses.Gateway /ip4/0.0.0.0/tcp/8080
ipfs config --json Swarm.ConnMgr.HighWater 2000
ipfs config --json Datastore.BloomFilterSize 1048576
ipfs config Datastore.StorageMax 100GB
`
}

func Undeploy(ctx context.Context, client kubernetes.Interface, topologyID, namespace string) error {
	svcName := "ipfs-" + strings.ReplaceAll(topologyID, "-", "")[:12]
	propagation := metav1.DeletePropagationForeground
	_ = client.AppsV1().StatefulSets(namespace).Delete(ctx, svcName, metav1.DeleteOptions{PropagationPolicy: &propagation})
	_ = client.CoreV1().Services(namespace).Delete(ctx, svcName, metav1.DeleteOptions{})
	_ = client.CoreV1().Services(namespace).Delete(ctx, svcName+"-external", metav1.DeleteOptions{})
	_ = client.CoreV1().ConfigMaps(namespace).Delete(ctx, svcName+"-scripts", metav1.DeleteOptions{})
	_ = client.CoreV1().ConfigMaps(namespace).Delete(ctx, svcName+"-env", metav1.DeleteOptions{})
	_ = client.CoreV1().Secrets(namespace).Delete(ctx, svcName+"-secrets", metav1.DeleteOptions{})
	return nil
}

type PodStatusResult struct {
	NodeID  string
	PodName string
	Phase   string
	Ready   bool
}

func GetPodsStatus(ctx context.Context, client kubernetes.Interface, topologyID, namespace string) ([]PodStatusResult, error) {
	svcName := "ipfs-" + strings.ReplaceAll(topologyID, "-", "")[:12]
	pods, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=" + svcName,
	})
	if err != nil {
		return nil, err
	}
	result := make([]PodStatusResult, 0, len(pods.Items))
	for _, p := range pods.Items {
		ready := false
		for _, c := range p.Status.Conditions {
			if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
				ready = true
				break
			}
		}
		result = append(result, PodStatusResult{
			NodeID:  p.Labels["statefulset.kubernetes.io/pod-name"],
			PodName: p.Name,
			Phase:   string(p.Status.Phase),
			Ready:   ready,
		})
	}
	return result, nil
}
