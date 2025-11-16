package models
type UserTenant struct {
    UserID   string `gorm:"type:uuid;primaryKey"`
    TenantID string `gorm:"type:uuid;primaryKey"`
    RoleID   string `gorm:"type:uuid"`

    User   User   `gorm:"foreignKey:UserID"`
    Tenant Tenant `gorm:"foreignKey:TenantID"`
    Role   Role   `gorm:"foreignKey:RoleID"`
}
