package kube

import (
	"ipfs-visualizer/config"

	"golang.org/x/exp/slog"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateKubeconfig(cfg config.KubeConfig) (*rest.Config, error) {
	slog.Info("Creating Kubeconfig")
	if cfg.ManualKubeConfigFlag {
		kubeconfig := cfg.KubeConfigPath
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, &ManualKubeConfigCreationError{Path: kubeconfig, Inner: err}
		}
		return config, nil
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, &AutoKubeConfigCreationError{Inner: err}
	}
	return config, nil
}

func CreateKubeClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	slog.Info("Creating KubeClientSet")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, &ClientCreationError{Config: config, Inner: err}
	}
	return clientset, nil
}
