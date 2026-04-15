package domain

type Node struct {
	Id       string `json:"id"`
	Hostname string `json:"hostname"`
	Role     string `json:"role"`
	Ip       string `json:"ip"`
	Os       string `json:"os,omitempty"`
	Arch     string `json:"arch,omitempty"`
}
