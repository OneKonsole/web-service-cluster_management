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

	if err != nil && err.Error() != "no clusters found" {
		log.Error("Error getting cluster list: ", err, " for client: ", ctx.Params("client"))
		return err
	} else if err != nil && err.Error() == "no clusters found" {
		log.Info("No clusters found for client: ", ctx.Params("client"))
		return ctx.Status(fiber.StatusNotFound).Send([]byte("No clusters found for client: " + client))
	}

	return ctx.JSON(clusterList)
}

func (h *ClusterHandler) FindAll(ctx *fiber.Ctx) error {

	clusterList, err := h.clusterUseCase.FindAll(ctx.Context())
	if err != nil && err.Error() != "no clusters found" {
		return err
	} else if err != nil && err.Error() == "no clusters found" {
		return ctx.Status(fiber.StatusNotFound).Send([]byte("No clusters found"))
	}

	return ctx.JSON(clusterList)
}

func (h *ClusterHandler) GetKubeconfig(ctx *fiber.Ctx) error {
	log.Info("GetKubeConfig called for client: ", ctx.Params("client"), " and cluster: ", ctx.Params("cluster"))

	client := ctx.Params("client")
	cluster := ctx.Params("cluster")

	format := ctx.Query("format")

	kubeConfig, err := h.clusterUseCase.GetKubeConfig(ctx.Context(), client, cluster)
	if err != nil && err.Error() != "no clusters found" {
		log.Error("Error getting kubeconfig file for cluster: ", ctx.Params("cluster"), " and client: ", ctx.Params("client"), " - ", err)
		return err
	} else if err != nil && err.Error() == "no clusters found" {
		log.Info("Cluster not found for client: ", ctx.Params("client"), " and cluster: ", ctx.Params("cluster"))
		return ctx.Status(fiber.StatusNotFound).Send([]byte("Cluster not found for client: " + client + " and cluster: " + cluster))
	}

	if format == "file" {
		log.Info("Sending kubeconfig file for cluster: ", ctx.Params("cluster"), " and client: ", ctx.Params("client"))
		ctx.Set("Content-Type", "application/octet-stream")
		ctx.Set("Content-Disposition", "attachment; filename="+cluster+"-kubeconfig")
		return ctx.SendString(kubeConfig)
	}

	// TODO: return kubeconfig file as attachment instead of json
	return ctx.JSON(kubeConfig)
}

func (h *ClusterHandler) FindByClientNameAndClusterName(ctx *fiber.Ctx) error {
	log.Info("GetClusterDetail called for client: ", ctx.Params("client"), " and cluster: ", ctx.Params("cluster"))

	client := ctx.Params("client")
	cluster := ctx.Params("cluster")

	result, err := h.clusterUseCase.FindByClientNameAndClusterName(ctx.Context(), client, cluster)
	if err != nil && err.Error() != "no clusters found" {
		log.Error("Error getting cluster detail for cluster: ", ctx.Params("cluster"), " and client: ", ctx.Params("client"), " - ", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if err != nil && err.Error() == "no clusters found" {
		log.Info("Cluster not found for client: ", ctx.Params("client"), " and cluster: ", ctx.Params("cluster"))
		return ctx.Status(fiber.StatusNotFound).Send([]byte("Cluster not found for client: " + client + " and cluster: " + cluster))
	}

	return ctx.JSON(result)
}

func (h *ClusterHandler) Delete(ctx *fiber.Ctx) error {
	log.Info("DeleteCluster called for client: ", ctx.Params("client"), " and cluster: ", ctx.Params("cluster"))

	client := ctx.Params("client")
	cluster := ctx.Params("cluster")

	err := h.clusterUseCase.Delete(ctx.Context(), client, cluster)
	if err != nil && err.Error() != "no clusters found" {
		log.Error("Error deleting cluster: ", ctx.Params("cluster"), " and client: ", ctx.Params("client"), " - ", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if err != nil && err.Error() == "no clusters found" {
		log.Info("Cluster not found for client: ", ctx.Params("client"), " and cluster: ", ctx.Params("cluster"))
		return ctx.Status(fiber.StatusNotFound).Send([]byte("Cluster not found for client: " + client + " and cluster: " + cluster))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
