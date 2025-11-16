package models

type Tenant struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OwnerId   string `gorm:"foreignKey:UserId"`
	StoreName string `gorm:"unique;not null" json:"storeName"`
	Slug      string `gorm:"not null" json:"slug"`
	SubDomain string `gorm:"unique;not null" json:"subDomain"`
	Status    string `gorm:"not null; default:'provisioning'" json:"status"`  //provisioning, active, suspended, pending deletion, deleted

	Users []User `gorm:"many2many:user_tenants;joinForeignKey:TenantID;JoinReferences:UserID, omitempty"` // all users in this tenant
}
