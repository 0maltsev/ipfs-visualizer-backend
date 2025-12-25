package nodes

type NodeSpec struct {
	NodeID        string            `json:"nodeId"`
	NodeName      string            `json:"nodeName"`
	Role          string            `json:"role"` // bootstrap | worker
	Ports         PortsSpec         `json:"ports"`
	Env           map[string]string `json:"env"`
	InitContainer InitContainerSpec `json:"initContainer"`
	Containers    []ContainerSpec   `json:"containers"`
	Volumes       VolumesSpec       `json:"volumes"`
}

type PortsSpec struct {
	SwarmTCP     int `json:"swarmTCP"`
	SwarmUDP     int `json:"swarmUDP"`
	API          int `json:"api"`
	HTTPGateway  int `json:"httpGateway"`
	WS           int `json:"ws"`
	ClusterAPI   int `json:"clusterAPI"`
	ClusterProxy int `json:"clusterProxy"`
	ClusterSwarm int `json:"clusterSwarm"`
}

type InitContainerSpec struct {
	Name    string   `json:"name"`
	Command []string `json:"command"`
}

type ContainerSpec struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command,omitempty"`
	Ports   []string `json:"ports,omitempty"` // имена портов из PortsSpec
}

type VolumesSpec struct {
	IPFSStorage    string `json:"ipfsStorage"`
	ClusterStorage string `json:"clusterStorage"`
	ScriptsConfig  string `json:"scriptsConfig"`
}
