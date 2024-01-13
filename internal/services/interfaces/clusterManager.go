package interfaces

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

type ClusterManager interface {
	GetKubernetesClientsetFromSecret(clientName, clusterName string, ctx context.Context) (*kubernetes.Clientset, error)
}
