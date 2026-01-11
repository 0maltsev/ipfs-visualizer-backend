package kube

import (
	"context"
	"strconv"

	"k8s.io/apimachinery/pkg/api/resource"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateCluster(
    ctx context.Context,
    client kubernetes.Interface,
    namespace string,
    clusterKubeRes ClusterKubeResources,
) (ClusterKubeResources,error) {

    configMap, err := createClusterConfigMap(ctx, client, namespace, clusterKubeRes)
	if err != nil {
        return clusterKubeRes, err
    }
	clusterKubeRes.ScriptsConfigMap = configMap

    secrets, err := createClusterSecrets(ctx, client, namespace, clusterKubeRes)
    if err != nil {
        return clusterKubeRes, err
    }
	clusterKubeRes.ClusterSecret = secrets.ClusterKubeSecret

    headlessService, err := createHeadlessService(ctx, client, namespace, clusterKubeRes)
	if err != nil {
        return clusterKubeRes, err
    }
	clusterKubeRes.HeadlessService = headlessService

    externalService, err := createExternalService(ctx, client, namespace, clusterKubeRes)
	if err != nil {
        return clusterKubeRes, err
    }
	clusterKubeRes.ExternalService = externalService

    statefulSet, err := createStatefulSet(ctx, client, namespace, clusterKubeRes)
	if err != nil {
        return clusterKubeRes, err
    }
	clusterKubeRes.StatefulSet = statefulSet

    return clusterKubeRes, nil
}

func createClusterConfigMap(
    ctx context.Context,
    client kubernetes.Interface,
    ns string,
    spec ClusterKubeResources,
) (*corev1.ConfigMap, error) {

    cm := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: spec.ScriptsConfig,
        },
		// TODO
        Data: map[string]string{
            "entrypoint.sh":      "...",
            "configure-ipfs.sh":  "...",
        },
    }

    configMap, err := client.CoreV1().
        ConfigMaps(ns).
        Create(ctx, cm, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

    return configMap, nil
}

type ClusterSecrets struct {
    ClusterSecret    string
    BootstrapPrivKey string
	ClusterKubeSecret   *corev1.Secret
}

func createClusterSecrets(
    ctx context.Context,
    client kubernetes.Interface,
    ns string,
    spec ClusterKubeResources,
) (*ClusterSecrets, error) {

	clusterSecret, err := generateClusterSecret()
	if err != nil {
		return  nil, err
	}

	bootstrapPrivKey, err := generateBootstrapPrivateKey()
	if err != nil {
		return  nil, err
	}

    secrets := &ClusterSecrets{
        ClusterSecret:    clusterSecret,
        BootstrapPrivKey: bootstrapPrivKey.PrivateKey,
    }

    secret := &corev1.Secret{
        ObjectMeta: metav1.ObjectMeta{
            Name: spec.ClusterName + "-secrets",
        },
        Type: corev1.SecretTypeOpaque,
        Data: map[string][]byte{
            "cluster-secret":        []byte(secrets.ClusterSecret),
            "bootstrap-peer-priv":   []byte(secrets.BootstrapPrivKey),
        },
    }

    clusterKubeSecret, err := client.CoreV1().
        Secrets(ns).
        Create(ctx, secret, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	secrets.ClusterKubeSecret = clusterKubeSecret

    return secrets, err
}

func createHeadlessService(
    ctx context.Context,
    client kubernetes.Interface,
    ns string,
    spec ClusterKubeResources,
) (*corev1.Service, error) {

    svc := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name: spec.ClusterName,
        },
        Spec: corev1.ServiceSpec{
            ClusterIP: "None",
            Selector: map[string]string{
                "app": spec.ClusterName,
            },
            Ports: []corev1.ServicePort{
                { Name: "cluster-swarm", Port: 9096 },
            },
        },
    }

    headlessService, err := client.CoreV1().
        Services(ns).
        Create(ctx, svc, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

    return headlessService, nil
}

func createExternalService(
    ctx context.Context,
    client kubernetes.Interface,
    ns string,
    spec ClusterKubeResources,
) (*corev1.Service, error) {

    svc := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name: spec.ClusterName + "-external",
        },
        Spec: corev1.ServiceSpec{
            Type: corev1.ServiceType(spec.ServiceType),
            Selector: map[string]string{
                "app": spec.ClusterName,
            },
			// TODO
            Ports: []corev1.ServicePort{
                { Name: "swarm", Port: 4001 },
                { Name: "api", Port: 5001 },
                { Name: "gateway", Port: 8080 },
            },
        },
    }

    externalService, err := client.CoreV1().
        Services(ns).
        Create(ctx, svc, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

    return externalService, nil
}

func createStatefulSet(
	ctx context.Context,
	client kubernetes.Interface,
	ns string,
	clusterSpec ClusterKubeResources,
) (*appsv1.StatefulSet, error) {

	replicas := int32(1) // обычно 1 для bootstrap

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterSpec.ClusterName,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: clusterSpec.ClusterName,
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": clusterSpec.ClusterName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": clusterSpec.ClusterName,
					},
				},
				// TODO
				Spec: corev1.PodSpec{
					// Можно оставить пустым или с базовыми контейнерами bootstrap
					// Контейнеры для дополнительных нод будут добавляться через CreateNode
				},
			},
			VolumeClaimTemplates: buildPVCs(clusterSpec),
		},
	}

	statefulSet, err := client.AppsV1().
		StatefulSets(ns).
		Create(ctx, sts, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return statefulSet, nil
}


func DeleteCluster(
	ctx context.Context,
	client kubernetes.Interface,
	namespace string,
	clusterName string,
) error {

	if err := deleteStatefulSet(ctx, client, namespace, clusterName); err != nil {
		return err
	}

	if err := deleteServices(ctx, client, namespace, clusterName); err != nil {
		return err
	}

	if err := deleteConfigMaps(ctx, client, namespace, clusterName); err != nil {
		return err
	}

	if err := deleteSecrets(ctx, client, namespace, clusterName); err != nil {
		return err
	}

	if err := deletePVCs(ctx, client, namespace, clusterName); err != nil {
		return err
	}

	return nil
}

func deleteStatefulSet(
	ctx context.Context,
	client kubernetes.Interface,
	ns string,
	clusterName string,
) error {

	propagation := metav1.DeletePropagationForeground

	return client.AppsV1().
		StatefulSets(ns).
		Delete(ctx, clusterName, metav1.DeleteOptions{
			PropagationPolicy: &propagation,
		})
}

func deleteServices(
	ctx context.Context,
	client kubernetes.Interface,
	ns string,
	clusterName string,
) error {

	services := []string{
		clusterName,                 // headless
		clusterName + "-external",   // external
	}

	for _, name := range services {
		_ = client.CoreV1().
			Services(ns).
			Delete(ctx, name, metav1.DeleteOptions{})
	}

	return nil
}

func deleteConfigMaps(
	ctx context.Context,
	client kubernetes.Interface,
	ns string,
	clusterName string,
) error {

	cmName := clusterName + "-scripts"

	_ = client.CoreV1().
		ConfigMaps(ns).
		Delete(ctx, cmName, metav1.DeleteOptions{})

	return nil
}

func deleteSecrets(
	ctx context.Context,
	client kubernetes.Interface,
	ns string,
	clusterName string,
) error {

	secretName := clusterName + "-secrets"

	_ = client.CoreV1().
		Secrets(ns).
		Delete(ctx, secretName, metav1.DeleteOptions{})

	return nil
}

func deletePVCs(
	ctx context.Context,
	client kubernetes.Interface,
	ns string,
	clusterName string,
) error {

	pvcs, err := client.CoreV1().
		PersistentVolumeClaims(ns).
		List(ctx, metav1.ListOptions{
			LabelSelector: "app=" + clusterName,
		})
	if err != nil {
		return err
	}

	for _, pvc := range pvcs.Items {
		_ = client.CoreV1().
			PersistentVolumeClaims(ns).
			Delete(ctx, pvc.Name, metav1.DeleteOptions{})
	}

	return nil
}

func buildIPFSContainer(
	cluster ClusterKubeResources,
	node NodeKubeResources,
) corev1.Container {

	return corev1.Container{
		Name:  "ipfs",
		Image: cluster.IPFSImage,

		Ports: []corev1.ContainerPort{
			{Name: "swarm-tcp", ContainerPort: int32(node.Ports.SwarmTCP)},
			{Name: "swarm-udp", ContainerPort: int32(node.Ports.SwarmUDP), Protocol: corev1.ProtocolUDP},
			{Name: "api", ContainerPort: int32(node.Ports.API)},
			{Name: "gateway", ContainerPort: int32(node.Ports.HTTPGateway)},
			{Name: "ws", ContainerPort: int32(node.Ports.WS)},
		},

		EnvFrom: []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cluster.EnvConfig,
					},
				},
			},
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cluster.ClusterName + "-secrets",
					},
				},
			},
		},

		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "ipfs-storage",
				MountPath: "/data/ipfs",
			},
			{
				Name:      "scripts",
				MountPath: "/scripts",
			},
		},

		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/api/v0/id",
					Port: intstr.FromInt(node.Ports.API),
				},
			},
			InitialDelaySeconds: 10,
			PeriodSeconds:       10,
		},
	}
}

func buildIPFSClusterContainer(
	cluster ClusterKubeResources,
	node NodeKubeResources,
) corev1.Container {

	env := []corev1.EnvVar{
		{
			Name:  "CLUSTER_PEERNAME",
			Value: node.NodeName,
		},
		{
			Name:  "CLUSTER_BOOTSTRAP",
			Value: strconv.FormatBool(node.Labels["role"] == "bootstrap"),
		},
	}

	return corev1.Container{
		Name:  "ipfs-cluster",
		Image: cluster.IPFSClusterImage,

		Ports: []corev1.ContainerPort{
			{Name: "cluster-api", ContainerPort: int32(node.Ports.ClusterAPI)},
			{Name: "cluster-proxy", ContainerPort: int32(node.Ports.ClusterProxy)},
			{Name: "cluster-swarm", ContainerPort: int32(node.Ports.ClusterSwarm)},
		},

		Env: env,

		EnvFrom: []corev1.EnvFromSource{
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cluster.ClusterName + "-secrets",
					},
				},
			},
		},

		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "cluster-storage",
				MountPath: "/data/ipfs-cluster",
			},
		},
	}
}

func buildPVCs(cluster ClusterKubeResources) []corev1.PersistentVolumeClaim {

	return []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ipfs-storage",
				Labels: map[string]string{
					"app":  cluster.ClusterName,
					"role": "ipfs",
				},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				StorageClassName: &cluster.StorageClass,
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(cluster.IPFSStorageSize),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster-storage",
				Labels: map[string]string{
					"app":  cluster.ClusterName,
					"role": "cluster",
				},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				StorageClassName: &cluster.StorageClass,
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(cluster.ClusterStorageSize),
					},
				},
			},
		},
	}
}
