package dto

type ResourceDetailDto struct {
	Id        string
	Name      string
	Kind      string
	Runtime   string
	CreatedAt string
	UpdatedAt string
	Status    map[string]interface{}
	Spec      map[string]interface{}
	LiveSpec  map[string]interface{}
	Events    []string
}
