package entity

// FeatureFlags controls optional system capabilities per tenant.
//
// These flags enable or disable features dynamically without schema changes.
type FeatureFlags struct {
    ReviewsEnabled bool `json:"reviews_enabled"`
    CouponsEnabled bool `json:"coupons_enabled"`
    BetaCheckoutV2 bool `json:"beta_checkout_v2"`
    MultiWarehouse bool `json:"multi_warehouse"`
}