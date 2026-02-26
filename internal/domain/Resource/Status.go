package resource

type Status struct {
	Phase            string `json:"phase,omitempty"`
	RuntimeId        string `json:"runtimeId,omitempty"`
	LastReconciledAt string `json:"lastReconciledAt,omitempty"`
}
