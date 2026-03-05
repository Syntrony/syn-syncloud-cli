package resource

type Ownership struct {
	OwnerId   string `json:"ownerId,omitempty"`
	OwnerKind string `json:"ownerKind,omitempty"`
	OwnerName string `json:"ownerName,omitempty"`
}
