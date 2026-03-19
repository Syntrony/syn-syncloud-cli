package domain

type RuntimeResourceState struct {
	Exists bool
	Spec   map[string]interface{}
}
