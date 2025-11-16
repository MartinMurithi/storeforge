package models

type UserTenant struct {
	UserID   string `gorm:"primaryKey" json:"userId"`
	TenantId string `gorm:"primaryKey" json:"tenantId"`
}
