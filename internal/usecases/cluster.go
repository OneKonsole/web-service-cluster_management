package usecases

import (
	"context"
	"errors"

	"github.com/oneKonsole/web-service-cluster_management/internal/models"
	iRepository "github.com/oneKonsole/web-service-cluster_management/internal/repositories/interfaces"
	eUseCase "github.com/oneKonsole/web-service-cluster_management/internal/usecases/errors"
	iUseCase "github.com/oneKonsole/web-service-cluster_management/internal/usecases/interfaces"
)

type clusterUseCase struct {
	clusterRepository iRepository.ClusterRepository
}

func NewClusterUseCase(clusterRepository iRepository.ClusterRepository) iUseCase.ClusterUsecase {
	return &clusterUseCase{
		clusterRepository: clusterRepository,
	}
}

func (c *clusterUseCase) FindAll(ctx context.Context) ([]models.Cluster, error) {
	result, err := c.clusterRepository.FindAll(ctx)
	if err != nil || len(result) == 0 {
		return result, eUseCase.NewErrorNotFound(errors.New("no clusters found"))
	}
	return result, nil
}

func (c *clusterUseCase) FindByClientName(ctx context.Context, clientName string) ([]models.Cluster, error) {
	result, err := c.clusterRepository.FindByClientName(ctx, clientName)
	if err != nil || len(result) == 0 {
		return result, eUseCase.NewErrorNotFound(errors.New("no clusters found"))
	}
	return result, nil
}

func (c *clusterUseCase) FindByClientNameAndClusterName(ctx context.Context, clientName string, clusterName string) (models.Cluster, error) {
	result, err := c.clusterRepository.FindByClientNameAndClusterName(ctx, clientName, clusterName)
	if err != nil || len(result.Name) == 0 {
		return result, eUseCase.NewErrorNotFound(errors.New("no clusters found"))
	}
	return result, nil
}
func (c *clusterUseCase) GetKubeConfig(ctx context.Context, clientName string, clusterName string) (string, error) {
	result, err := c.clusterRepository.GetKubeConfig(ctx, clientName, clusterName)
	if err != nil || len(result) == 0 {
		return result, eUseCase.NewErrorNotFound(errors.New("no clusters found"))
	}
	return result, nil
}

func (c *clusterUseCase) Delete(ctx context.Context, clientName string, clusterName string) error {
	err := c.clusterRepository.Delete(ctx, clientName, clusterName)
	if err != nil {
		return eUseCase.NewErrorBusinessException(errors.New("error deleting cluster"))
	}
	return nil
}
