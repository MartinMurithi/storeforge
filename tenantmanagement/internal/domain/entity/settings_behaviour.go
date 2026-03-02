package entity

import (
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
)

// ChangeTheme swaps the layout engine and resets the config to a baseline.
func (s *Settings) ChangeTheme(themeID value_object.ThemeID, baseConfig ThemeConfig) {
	s.ThemeID = themeID
	s.Config = baseConfig
	s.Version++
	s.UpdatedAt = time.Now()
}

// UpdateOverrides performs a partial merge of new visual settings.
func (s *Settings) UpdateOverrides(overrides ThemeConfig) {
	if s.Config == nil {
		s.Config = make(ThemeConfig)
	}

	for key, value := range overrides {
		s.Config[key] = value
	}

	s.Version++
	s.UpdatedAt = time.Now()
}