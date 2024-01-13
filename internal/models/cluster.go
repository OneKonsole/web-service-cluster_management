package models

type Cluster struct {
	Name              string       `json:"name"`
	Status            string       `json:"status"`
	ControlPlane      ControlPlane `json:"controlPlane"`
	NodeList          []Node       `json:"nodes"`
	KubernetesVersion string       `json:"kubernetesVersion"`
	OrderID           string       `json:"orderId"`
}

type Element struct {
	Name                   string `json:"name"`
	ReadyNumber            int    `json:"readyNumber"`
	DesiredNumberScheduled int    `json:"desiredNumberScheduled"`
}

type ControlPlane struct {
	KonnectivityServer    Element `json:"konnectivity-server"`
	KubeApiserver         Element `json:"kube-apiserver"`
	KubeControllerManager Element `json:"kube-controller-manager"`
	KubeScheduler         Element `json:"kube-scheduler"`
}

type Node struct {
	Name  string   `json:"name"`
	Ready bool     `json:"ready"`
	Roles []string `json:"roles"`
}
