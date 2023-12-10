package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	flags "github.com/jessevdk/go-flags"
	"github.com/oneKonsole/web-service-cluster_management/internal/controllers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
		clientSet, err = GetKubernetesClientset()
		if err != nil {
			fmt.Println("Error creating Kubernetes client: ", err)
			os.Exit(1)
		}
	}

	// Create a new KubeController
	kubeController := controllers.KubeController{
		ClientSet: clientSet,
	}

	app.Get("/:client", kubeController.GetClusters)
	//TODO: Restrict access to this endpoint for admin users only
	app.Get("/", kubeController.GetAllClusters)

	app.Listen(":3000")
}

// GetKubernetesClientset returns a Kubernetes clientset using the kubeconfig file at the default location.
func GetKubernetesClientset() (*kubernetes.Clientset, error) {
	// Get the kubeconfig file path from the default location
	kubeconfig := arguments.KubeConfigPath

	// Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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
