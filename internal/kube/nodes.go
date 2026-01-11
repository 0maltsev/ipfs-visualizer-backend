package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
)

func AddNodeToCluster(
	ctx context.Context,
	client kubernetes.Interface,
	namespace string,
	clusterName string,
	nodeSpec NodeKubeResources,
) (*appsv1.StatefulSet, error) {

	sts, err := client.AppsV1().
		StatefulSets(namespace).
		Get(ctx, clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	newReplicas := *sts.Spec.Replicas + 1
	sts.Spec.Replicas = &newReplicas

	applyNodeSpecToPodTemplate(&sts.Spec.Template, nodeSpec)

	statefulSet, err := client.AppsV1().
		StatefulSets(namespace).
		Update(ctx, sts, metav1.UpdateOptions{})

	if err != nil {
		return nil, err
	}

	return statefulSet, nil
}

func applyNodeSpecToPodTemplate(
	podTemplate *corev1.PodTemplateSpec,
	nodeSpec NodeKubeResources,
) {

	// ENV
	for i := range podTemplate.Spec.Containers {
		for k, v := range nodeSpec.Env {
			podTemplate.Spec.Containers[i].Env = append(
				podTemplate.Spec.Containers[i].Env,
				corev1.EnvVar{
					Name:  k,
					Value: v,
				},
			)
		}
	}

	applyPorts(podTemplate, nodeSpec.Ports)
}

func applyPorts(
	podTemplate *corev1.PodTemplateSpec,
	ports PortsSpec,
) {
	for i := range podTemplate.Spec.Containers {
		c := &podTemplate.Spec.Containers[i]

		c.Ports = []corev1.ContainerPort{
			{ Name: "swarm", ContainerPort: int32(ports.SwarmTCP) },
			{ Name: "api", ContainerPort: int32(ports.API) },
			{ Name: "gateway", ContainerPort: int32(ports.HTTPGateway) },
			{ Name: "cluster-api", ContainerPort: int32(ports.ClusterAPI) },
		}
	}
}

func RemoveNodeFromCluster(
	ctx context.Context,
	client kubernetes.Interface,
	namespace string,
	clusterName string,
) error {

	sts, err := client.AppsV1().
		StatefulSets(namespace).
		Get(ctx, clusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if *sts.Spec.Replicas == 0 {
		return nil
	}

	newReplicas := *sts.Spec.Replicas - 1
	sts.Spec.Replicas = &newReplicas

	_, err = client.AppsV1().
		StatefulSets(namespace).
		Update(ctx, sts, metav1.UpdateOptions{})

	return err
}
