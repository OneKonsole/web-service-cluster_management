package models

import (
	"context"
	"encoding/json"
	"fmt"

	kamaji "github.com/clastix/kamaji/api/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type Cluster struct {
	Name                    string                  `json:"name"`
	Status                  string                  `json:"status"`
	ControlPlaneElementList ControlPlaneElementList `json:"controlPlaneElements"`
	NodeList                NodeList                `json:"nodes"`
	KubernetesVersion       string                  `json:"kubernetesVersion"`
	OrderID                 string                  `json:"orderId"`
}

type ClusterList struct {
	Clusters []Cluster `json:"clusters"`
}

type ControlPlaneElementList struct {
	ControlPlaneElements []ControlPlaneElement `json:"controlPlaneElements"`
}

type ControlPlaneElement struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Replicas int    `json:"replicas"`
	Memory   string `json:"memory"`
	Cpu      string `json:"cpu"`
}

type NodeList struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Roles  string `json:"roles"`
}

// GetClusterList returns a list of clusters for a specific client
func GetClusterList(
	ctx context.Context,
	clientSet *kubernetes.Clientset,
	clientName string,
) (ClusterList, error) {
	var clusterList ClusterList

	// Get all tenantcontrolplanes.kamaji.clastix.io CRs in the cluster filtered by clientName label
	path := fmt.Sprintf("/apis/kamaji.clastix.io/v1alpha1/namespaces/%s/tenantcontrolplanes", clientName)
	response := clientSet.CoreV1().RESTClient().Get().
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

	// Iterate over all tenantcontrolplanes.kamaji.clastix.io CRs
	for _, tenantControlPlane := range kamaji.Items {
		var cluster Cluster
		cluster.Name = tenantControlPlane.GetName()
		cluster.Status = tenantControlPlane.Status.ControlPlaneEndpoint
		cluster.KubernetesVersion = tenantControlPlane.Spec.Kubernetes.Version
		cluster.OrderID = tenantControlPlane.GetLabels()["order"]

		clusterList.Clusters = append(clusterList.Clusters, cluster)
	}

	// json clusterList and print it
	jsonClusterList, err := json.Marshal(clusterList)
	if err != nil {
		return clusterList, err
	}
	fmt.Println(string(jsonClusterList))

	return clusterList, nil
}

// GetAllClusterList returns a list of all clusters in the cluster
func GetAllClusterList(
	ctx context.Context,
	clientSet *kubernetes.Clientset,
) (ClusterList, error) {
	var clusterList ClusterList

	// Get all tenantcontrolplanes.kamaji.clastix.io CRs in the cluster filtered by clientName label
	path := fmt.Sprintf("/apis/kamaji.clastix.io/v1alpha1/tenantcontrolplanes")
	response := clientSet.CoreV1().RESTClient().Get().
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

	// Iterate over all tenantcontrolplanes.kamaji.clastix.io CRs
	for _, tenantControlPlane := range kamaji.Items {
		var cluster Cluster
		cluster.Name = tenantControlPlane.GetName()
		cluster.Status = tenantControlPlane.Status.ControlPlaneEndpoint
		cluster.KubernetesVersion = tenantControlPlane.Spec.Kubernetes.Version
		cluster.OrderID = tenantControlPlane.GetLabels()["order"]

		clusterList.Clusters = append(clusterList.Clusters, cluster)
	}

	// json clusterList and print it
	jsonClusterList, err := json.Marshal(clusterList)
	if err != nil {
		return clusterList, err
	}
	fmt.Println(string(jsonClusterList))

	return clusterList, nil
}
