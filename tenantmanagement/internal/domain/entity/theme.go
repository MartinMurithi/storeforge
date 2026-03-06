package entity

import (
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
)

type Theme struct {
    ID          value_object.ThemeID
    Name        string
    Description string
    DefaultConfig *Settings 
    IsActive    bool
    CreatedAt   time.Time
}