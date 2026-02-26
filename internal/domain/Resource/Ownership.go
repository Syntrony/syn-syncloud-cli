package resource

type Ownership struct {
	OwnerType string `json:"ownerType,omitempty"`
	OwnerId   string `json:"ownerId,omitempty"`
}
