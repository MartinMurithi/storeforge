package models

type User struct {
	ID           string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	FullNames    string `gorm:"not null" json:"fullNames"`
	Email        string `gorm:"unique;not null" json:"email"`
	Phone        string `gorm:"unique;not null" json:"phone"`
	PasswordHash string `gorm:"not null" json:"password"`
	BusinessType string `gorm:"not null" json:"businessType"`        //help select default theme
	BusinessName string `gorm:"unique;not null" json:"businessName"` //generates slug for domain
	IsVerified   bool   `gorm:"default:false" json:"isVerified"`
	Roles        []Role `gorm:"many2many:user_roles" json:"-"`

	Tenants []Tenant `gorm:"many2many:user_tenants;joinForeignKey:UserID;JoinReferences:TenantID, omitempty"` // all tenants user belongs to
}
