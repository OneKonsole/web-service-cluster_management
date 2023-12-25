package handlers

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/gofiber/fiber/v2/log"
	iUseCase "github.com/oneKonsole/web-service-cluster_management/internal/usecases/interfaces"
)

type ClusterHandler struct {
	clusterUseCase iUseCase.ClusterUsecase
}

func NewClusterHandler(clusterUseCase iUseCase.ClusterUsecase) *ClusterHandler {
	return &ClusterHandler{
		clusterUseCase: clusterUseCase,
	}
}

func (h *ClusterHandler) FindByClientName(ctx *fiber.Ctx) error {
	log.Info("GetClusters called for client: ", ctx.Params("client"))

	client := ctx.Params("client")

	clusterList, err := h.clusterUseCase.FindByClientName(ctx.Context(), client)

	if err != nil {
		log.Error("Error getting cluster list: ", err)
		return err
	}

	return ctx.JSON(clusterList)
}

func (h *ClusterHandler) FindAll(ctx *fiber.Ctx) error {

	clusterList, err := h.clusterUseCase.FindAll(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(clusterList)
}

func (h *ClusterHandler) GetKubeconfig(ctx *fiber.Ctx) error {
	log.Info("GetKubeConfig called for client: ", ctx.Params("client"), " and cluster: ", ctx.Params("cluster"))

	client := ctx.Params("client")
	cluster := ctx.Params("cluster")

	kubeConfig, err := h.clusterUseCase.GetKubeConfig(ctx.Context(), client, cluster)
	if err != nil {
		log.Error("Error getting kubeconfig file for cluster: ", ctx.Params("cluster"), " and client: ", ctx.Params("client"), " - ", err)
		return err
	}
	// TODO: return kubeconfig file as attachment instead of json
	return ctx.JSON(kubeConfig)
}

func (h *ClusterHandler) FindByClientNameAndClusterName(ctx *fiber.Ctx) error {
	log.Info("GetClusterDetail called for client: ", ctx.Params("client"), " and cluster: ", ctx.Params("cluster"))

	client := ctx.Params("client")
	cluster := ctx.Params("cluster")

	result, err := h.clusterUseCase.FindByClientNameAndClusterName(ctx.Context(), client, cluster)
	if err != nil {
		log.Error("Error getting cluster detail for cluster: ", ctx.Params("cluster"), " and client: ", ctx.Params("client"), " - ", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(result)
}

func (h *ClusterHandler) Delete(ctx *fiber.Ctx) error {
	log.Info("DeleteCluster called for client: ", ctx.Params("client"), " and cluster: ", ctx.Params("cluster"))

	client := ctx.Params("client")
	cluster := ctx.Params("cluster")

	err := h.clusterUseCase.Delete(ctx.Context(), client, cluster)
	if err != nil {
		log.Error("Error deleting cluster: ", ctx.Params("cluster"), " and client: ", ctx.Params("client"), " - ", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
