package controllers

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/gofiber/fiber/v2/log"
	"github.com/oneKonsole/web-service-cluster_management/internal/models"
	"k8s.io/client-go/kubernetes"
)

type KubeController struct {
	ClientSet *kubernetes.Clientset
}

func (kc KubeController) GetClusters(ctx *fiber.Ctx) error {
	log.Info("GetClusters called for client: ", ctx.Params("client"))

	clusterList, err := models.GetClusterList(ctx.UserContext(), kc.ClientSet, ctx.Params("client"))
	if err != nil {
		log.Error("Error getting cluster list: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(clusterList.Clusters) == 0 {
		log.Error("No clusters found for client: ", ctx.Params("client"))
		return fiber.NewError(fiber.StatusNotFound, "No clusters found")
	}

	return ctx.JSON(clusterList)
}

func (kc KubeController) GetAllClusters(ctx *fiber.Ctx) error {

	clusterList, err := models.GetAllClusterList(ctx.UserContext(), kc.ClientSet)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(clusterList)
}
