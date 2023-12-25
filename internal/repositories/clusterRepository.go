package repositories

import (
	"context"
	"fmt"
	"strings"

	kamaji "github.com/clastix/kamaji/api/v1alpha1"
	"github.com/oneKonsole/web-service-cluster_management/internal/models"
	iRepository "github.com/oneKonsole/web-service-cluster_management/internal/repositories/interfaces"
	iService "github.com/oneKonsole/web-service-cluster_management/internal/services/interfaces"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type clusterDatabase struct {
	clientsetClusterMaster *kubernetes.Clientset
	clusterManager         iService.ClusterManager
}

func NewClusterDatabase(clientsetClusterMaster *kubernetes.Clientset, clusterManager iService.ClusterManager) iRepository.ClusterRepository {
	return &clusterDatabase{
		clientsetClusterMaster: clientsetClusterMaster,
		clusterManager:         clusterManager,
	}
}

func (c *clusterDatabase) FindAll(ctx context.Context) ([]models.Cluster, error) {
	var clusterList []models.Cluster

	// Get all tenantcontrolplanes.kamaji.clastix.io CRs in the cluster filtered by clientName label
	path := fmt.Sprintf("/apis/kamaji.clastix.io/v1alpha1/tenantcontrolplanes")
	response := c.clientsetClusterMaster.CoreV1().RESTClient().Get().
		AbsPath(path).
		Do(ctx)

	// Check if there is an error
	err := response.Error()
	if err != nil {
		return clusterList, err
	}

	// Get the response json body
	var kamaji kamaji.TenantControlPlaneList
	err = response.Into(&kamaji)
	if err != nil {
		return clusterList, err
	}

	for _, tenantControlPlane := range kamaji.Items {
		var cluster models.Cluster
		cluster.Name = tenantControlPlane.GetName()
		cluster.Status = string(*tenantControlPlane.Status.Kubernetes.Version.Status)
		cluster.KubernetesVersion = tenantControlPlane.Spec.Kubernetes.Version
		cluster.OrderID = tenantControlPlane.GetLabels()["order"]

		if cluster.Status != "Ready" {
			clusterList = append(clusterList, cluster)
			continue
		}

		client := tenantControlPlane.GetLabels()["client"]

		clientsetTenant, err := c.clusterManager.GetKubernetesClientsetFromSecret(client, cluster.Name, ctx)
		if err != nil {
			fmt.Println("Failed to get clientset for cluster " + cluster.Name)
			clusterList = append(clusterList, cluster)
			continue
		}

		nodes, err := clientsetTenant.CoreV1().Nodes().List(ctx, v1.ListOptions{})
		if err != nil {
			fmt.Println("Failed to get nodes for cluster " + cluster.Name)
			clusterList = append(clusterList, cluster)
			continue
		}

		for _, node := range nodes.Items {
			labels := node.GetLabels()
			roles := make([]string, 0, len(labels))
			// Get a label key starting with node-role.kubernetes.io/
			for key := range labels {
				if strings.HasPrefix(key, "node-role.kubernetes.io/") {
					roles = append(roles, strings.TrimPrefix(key, "node-role.kubernetes.io/"))
				}
			}

			conditions := node.Status.Conditions
			isReady := false
			for _, condition := range conditions {
				if condition.Type == "Ready" && condition.Status == "True" {
					isReady = true
				}
			}

			cluster.NodeList = append(cluster.NodeList, models.Node{
				Name:  node.GetName(),
				Ready: isReady,
				Roles: roles,
			})
		}

		if err != nil {
			fmt.Println("Failed to get nodes for cluster " + cluster.Name)
			clusterList = append(clusterList, cluster)
			continue
		}

		clusterList = append(clusterList, cluster)
	}

	return clusterList, nil
}

func (c *clusterDatabase) FindByClientName(ctx context.Context, clientName string) ([]models.Cluster, error) {
	var clusterList []models.Cluster

	// Get all tenantcontrolplanes.kamaji.clastix.io CRs in the cluster filtered by clientName label
	path := fmt.Sprintf("/apis/kamaji.clastix.io/v1alpha1/namespaces/%s/tenantcontrolplanes", clientName)
	response := c.clientsetClusterMaster.CoreV1().RESTClient().Get().
		AbsPath(path).
		Param("labelSelector", "client="+clientName).
		Do(ctx)

	// Check if there is an error
	err := response.Error()
	if err != nil {
		return clusterList, err
	}

	// Get the response json body
	var kamaji kamaji.TenantControlPlaneList
	err = response.Into(&kamaji)
	if err != nil {
		return clusterList, err
	}

	for _, tenantControlPlane := range kamaji.Items {
		var cluster models.Cluster
		cluster.Name = tenantControlPlane.GetName()
		cluster.Status = string(*tenantControlPlane.Status.Kubernetes.Version.Status)
		cluster.KubernetesVersion = tenantControlPlane.Spec.Kubernetes.Version
		cluster.OrderID = tenantControlPlane.GetLabels()["order"]

		if cluster.Status != "Ready" {
			clusterList = append(clusterList, cluster)
			continue
		}

		clientsetTenant, err := c.clusterManager.GetKubernetesClientsetFromSecret(clientName, cluster.Name, ctx)
		if err != nil {
			fmt.Println("Failed to get clientset for cluster " + cluster.Name)
			clusterList = append(clusterList, cluster)
			continue
		}

		nodes, err := clientsetTenant.CoreV1().Nodes().List(ctx, v1.ListOptions{})
		if err != nil {
			fmt.Println("Failed to get nodes for cluster " + cluster.Name)
			clusterList = append(clusterList, cluster)
			continue
		}

		for _, node := range nodes.Items {
			labels := node.GetLabels()
			roles := make([]string, 0, len(labels))
			// Get a label key starting with node-role.kubernetes.io/
			for key := range labels {
				if strings.HasPrefix(key, "node-role.kubernetes.io/") {
					roles = append(roles, strings.TrimPrefix(key, "node-role.kubernetes.io/"))
				}
			}

			conditions := node.Status.Conditions
			isReady := false
			for _, condition := range conditions {
				if condition.Type == "Ready" && condition.Status == "True" {
					isReady = true
				}
			}

			cluster.NodeList = append(cluster.NodeList, models.Node{
				Name:  node.GetName(),
				Ready: isReady,
				Roles: roles,
			})
		}

		if err != nil {
			fmt.Println("Failed to get nodes for cluster " + cluster.Name)
			clusterList = append(clusterList, cluster)
			continue
		}

		clusterList = append(clusterList, cluster)
	}

	return clusterList, nil
}

func (c *clusterDatabase) FindByClientNameAndClusterName(ctx context.Context, clientName string, clusterName string) (models.Cluster, error) {
	var cluster models.Cluster

	// Get all tenantcontrolplanes.kamaji.clastix.io CRs in the cluster filtered by clientName label
	path := fmt.Sprintf("/apis/kamaji.clastix.io/v1alpha1/namespaces/%s/tenantcontrolplanes/", clientName)
	response := c.clientsetClusterMaster.CoreV1().RESTClient().Get().
		AbsPath(path).
		Param("labelSelector", "client="+clientName).
		Do(ctx)

	// Check if there is an error
	err := response.Error()
	if err != nil {
		return cluster, err
	}

	// Get the response json body
	var kamaji kamaji.TenantControlPlaneList
	err = response.Into(&kamaji)
	if err != nil {
		return cluster, err
	}

	for _, tenantControlPlane := range kamaji.Items {
		if tenantControlPlane.GetName() != clusterName {
			continue
		}

		cluster.Name = tenantControlPlane.GetName()
		cluster.Status = string(*tenantControlPlane.Status.Kubernetes.Version.Status)
		cluster.KubernetesVersion = tenantControlPlane.Spec.Kubernetes.Version
		cluster.OrderID = tenantControlPlane.GetLabels()["order"]

		if cluster.Status != "Ready" {
			return cluster, nil
		}

		clientsetTenant, err := c.clusterManager.GetKubernetesClientsetFromSecret(clientName, cluster.Name, ctx)
		if err != nil {
			fmt.Println("Failed to get clientset for cluster " + cluster.Name)
			continue
		}

		nodes, err := clientsetTenant.CoreV1().Nodes().List(ctx, v1.ListOptions{})
		if err != nil {
			fmt.Println("Failed to get nodes for cluster " + cluster.Name)
			continue
		}

		for _, node := range nodes.Items {
			labels := node.GetLabels()
			roles := make([]string, 0, len(labels))
			// Get a label key starting with node-role.kubernetes.io/
			for key := range labels {
				if strings.HasPrefix(key, "node-role.kubernetes.io/") {
					roles = append(roles, strings.TrimPrefix(key, "node-role.kubernetes.io/"))
				}
			}

			conditions := node.Status.Conditions
			isReady := false
			for _, condition := range conditions {
				if condition.Type == "Ready" && condition.Status == "True" {
					isReady = true
				}
			}

			cluster.NodeList = append(cluster.NodeList, models.Node{
				Name:  node.GetName(),
				Ready: isReady,
				Roles: roles,
			})
		}

		if err != nil {
			fmt.Println("Failed to get nodes for cluster " + cluster.Name)
			continue
		}

	}

	return cluster, nil
}

func (c *clusterDatabase) GetKubeConfig(ctx context.Context, clientName string, clusterName string) (string, error) {
	// Get the kubeconfig secret for the cluster
	kubeConfigSecret, err := c.clientsetClusterMaster.CoreV1().Secrets(clientName).Get(ctx, clusterName+"-admin-kubeconfig", v1.GetOptions{})
	if err != nil {
		fmt.Println("Failed to get kubeconfig for cluster " + clusterName)
		return "", err
	}

	// fmt.Println("Kubeconfig secret: ", string(kubeConfigSecret.Data["admin.conf"]))

	kubeConfig := kubeConfigSecret.Data["admin.conf"]

	return string(kubeConfig), nil
}

func (c *clusterDatabase) Delete(ctx context.Context, clientName string, clusterName string) error {
	// Delete the tenantcontrolplane CR
	path := fmt.Sprintf("/apis/kamaji.clastix.io/v1alpha1/namespaces/%s/tenantcontrolplanes/%s", clientName, clusterName)
	response := c.clientsetClusterMaster.CoreV1().RESTClient().Delete().
		AbsPath(path).
		Do(ctx)

	// Check if there is an error
	err := response.Error()
	if err != nil {
		return err
	}

	return nil
}
