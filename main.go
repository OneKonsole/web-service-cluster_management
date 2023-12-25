package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	flags "github.com/jessevdk/go-flags"
	handler "github.com/oneKonsole/web-service-cluster_management/internal/controllers/handlers"
	repository "github.com/oneKonsole/web-service-cluster_management/internal/repositories"
	service "github.com/oneKonsole/web-service-cluster_management/internal/services"
	usecase "github.com/oneKonsole/web-service-cluster_management/internal/usecases"
	"github.com/oneKonsole/web-service-cluster_management/internal/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Arguments struct {
	TypeKubernetesConnection string `short:"t" long:"type" description:"Type of Kubernetes connection" choice:"inCluster" choice:"kubeConfig" required:"true"`
	KubeConfigPath           string `short:"k" long:"kubeConfig" description:"Path to kubeconfig file"`
}

var arguments = Arguments{
	TypeKubernetesConnection: "kubeConfig",
	KubeConfigPath:           os.Getenv("HOME") + "/.kube/config",
}

var clientSet *kubernetes.Clientset

func main() {
	_, err := flags.Parse(&arguments)
	if err != nil {
		fmt.Println("Error parsing flags: ", err)
		os.Exit(1)
	}

	// Change log level to debug
	app := fiber.New(
		fiber.Config{
			Prefork: true, // When set to true, this will spawn multiple Go processes listening on the same port.
		},
	)

	// Connect to Kubernetes cluster and get clientSet
	switch arguments.TypeKubernetesConnection {
	case "inCluster":
		config, err := rest.InClusterConfig()
		if err != nil {
			fmt.Println("Error creating Kubernetes config: ", err)
			os.Exit(1)
		}
		clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println("Error creating Kubernetes client: ", err)
			os.Exit(1)
		}
	case "kubeConfig":
		clientSet, err = utils.GetKubernetesClientsetFromFilePath(arguments.KubeConfigPath)
		if err != nil {
			fmt.Println("Error creating Kubernetes client: ", err)
			os.Exit(1)
		}
	}

	// TODO; Externalize this to a config file or environment variable (see dependency injection as wire and viper for the config file)
	// Initialize config
	clusterManager := service.NewKubernetesClusterManager(clientSet)
	clusterRepository := repository.NewClusterDatabase(clientSet, clusterManager)
	clusterUseCase := usecase.NewClusterUseCase(clusterRepository)
	clusterHandler := handler.NewClusterHandler(clusterUseCase)

	app.Get("/:client", clusterHandler.FindByClientName)
	//TODO: Restrict access to this endpoint for admin users only
	app.Get("/", clusterHandler.FindAll)
	app.Get("/:client/:cluster/kubeconfig", clusterHandler.GetKubeconfig)
	app.Get("/:client/:cluster", clusterHandler.FindByClientNameAndClusterName)
	app.Delete("/:client/:cluster", clusterHandler.Delete)

	app.Listen(":3000")
}
