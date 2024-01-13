package services

import (
	"context"
	"fmt"

	"github.com/oneKonsole/web-service-cluster_management/internal/services/interfaces"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClusterManager struct {
	clientset *kubernetes.Clientset
}

func NewKubernetesClusterManager(clientset *kubernetes.Clientset) interfaces.ClusterManager {
	return &KubernetesClusterManager{
		clientset: clientset,
	}
}

func (kcm *KubernetesClusterManager) GetKubernetesClientsetFromSecret(clientName, clusterName string, ctx context.Context) (*kubernetes.Clientset, error) {
	// Get the kubeconfig secret for the cluster
	kubeConfigSecret, err := kcm.clientset.CoreV1().Secrets(clientName).Get(ctx, clusterName+"-admin-kubeconfig", v1.GetOptions{})
	if err != nil {
		fmt.Println("Failed to get kubeconfig for cluster " + clusterName)
		return nil, err
	}

	// fmt.Println("Kubeconfig secret: ", string(kubeConfigSecret.Data["admin.conf"]))

	kubeConfigSecretContent := kubeConfigSecret.Data["admin.conf"]
	clientset, err := getKubernetesClientsetFromKubeConfig(kubeConfigSecretContent)
	if err != nil {
		fmt.Println("Failed to connect to cluster " + clusterName + " using kubeconfig secret")
		return nil, err
	}

	return clientset, nil
}

func getKubernetesClientsetFromKubeConfig(kubeConfig []byte) (*kubernetes.Clientset, error) {
	// Use the current context in kubeconfig
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig: %v", err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	return clientset, nil
}
