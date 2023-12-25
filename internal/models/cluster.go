package models

type Cluster struct {
	Name                    string                `json:"name"`
	Status                  string                `json:"status"`
	ControlPlaneElementList []ControlPlaneElement `json:"controlPlaneElements"`
	NodeList                []Node                `json:"nodes"`
	KubernetesVersion       string                `json:"kubernetesVersion"`
	OrderID                 string                `json:"orderId"`
}

type ControlPlaneElement struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Replicas int    `json:"replicas"`
	Memory   string `json:"memory"`
	Cpu      string `json:"cpu"`
}

type Node struct {
	Name  string   `json:"name"`
	Ready bool     `json:"ready"`
	Roles []string `json:"roles"`
}
