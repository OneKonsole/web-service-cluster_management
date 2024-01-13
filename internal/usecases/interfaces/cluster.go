package interfaces

import (
	"context"

	"github.com/oneKonsole/web-service-cluster_management/internal/models"
)

type ClusterUsecase interface {
	FindAll(ctx context.Context) ([]models.Cluster, error)
	FindByClientName(ctx context.Context, clientName string) ([]models.Cluster, error)
	FindByClientNameAndClusterName(ctx context.Context, clientName string, clusterName string) (models.Cluster, error)
	GetKubeConfig(ctx context.Context, clientName string, clusterName string) (string, error)
	Delete(ctx context.Context, clientName string, clusterName string) error
}
