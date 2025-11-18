package models

type Tenant struct {
	ID      string `json:"id"`
	OwnerID string `json:"ownerId"`
	Owner   User

	StoreName string `json:"storeName"`
	Slug      string `json:"slug"`
	SubDomain string `json:"subDomain"`
	Status    string `json:"status"` //provisioning, active, suspended, pending deletion, deleted

	Users []User // all users in this tenant
}
