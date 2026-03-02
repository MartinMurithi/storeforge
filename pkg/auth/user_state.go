package auth

// UserState represents the current progress of a user within the StoreForge ecosystem.
// This state dictates the level of access granted by the API Gateway.
type UserState string

const (

    // StatePendingOnboarding: User is authenticated but has no associated tenant.
    // Access: Restricted to the Onboarding Service to provide Business Name/Theme.
    // Token: Identity-only (TenantID is nil).
    StatePendingOnboarding UserState = "PENDING_ONBOARDING"

    // StateActive: User has successfully provisioned at least one tenant.
    // Access: Full access to the Dashboard and Store resources based on PBAC.
    // Token: Full Context (contains TenantID and Role).
    StateActive UserState = "ACTIVE"

    // StateSuspended: User or Tenant has been flagged/disabled.
    // Access: Denied across all services until administrative resolution.
    StateSuspended UserState = "SUSPENDED"
)