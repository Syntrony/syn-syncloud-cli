package docker

type ContainerList []Container

type Container struct {
	ID      string `json:"ID"`
	Names   string `json:"Names"`
	Image   string `json:"Image"`
	Command string `json:"Command"`
	Created string `json:"Created"`
	Status  string `json:"Status"`
	Ports   string `json:"Ports"`
	State   string `json:"State"`
	Labels  string `json:"Labels"`
}
