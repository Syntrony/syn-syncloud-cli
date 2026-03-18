package domain

type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionNoop   ActionType = "noop"
)

type Action struct {
	Type     ActionType
	Resource *Resource
}

type ReconcileContext struct {
	Resources []*Resource
	State     *State
	Actions   []*Action
}
