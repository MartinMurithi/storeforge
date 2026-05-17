package entity

import "time"

type StoreSchema struct {
    SchemaVersion string    `json:"schema_version"` // controls migrations + compatibility
    TenantID string `json:"tenant_id"`
    Tenant Tenant `json:"tenant"`
    Settings Settings `json:"settings"`
    Theme ThemeConfig `json:"theme"`
    Storefront Storefront `json:"storefront"` // (Step 2 system)
    Actions Actions `json:"actions"` // behavior system
    Permissions Permissions `json:"permissions"` // access control

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Storefront represents the UI structure of a tenant's store.
// It defines pages, layout trees, and section composition.
type Storefront struct {
    Pages []Page `json:"pages"`
}

// Page represents a single route in the storefront (e.g. "/", "/products").
type Page struct {
    ID    string `json:"id"`
    Path  string `json:"path"`

    Type  string `json:"type"` // static | product | collection | cart

    Layout []Section `json:"layout"`
}

// Section represents a UI building block in the storefront.
//
// Sections form a tree structure that the frontend renderer converts into UI components.
type Section struct {
    ID   string `json:"id"`
    Type string `json:"type"` // MUST map to a registered component

    Props map[string]any `json:"props"`

    Children []Section `json:"children,omitempty"`
}

// Actions represent executable behaviors triggered from the UI.
// The frontend does NOT implement logic — it only references action IDs.
// Actions is a collection of all executable behavior definitions for a tenant.
type Actions struct {
    Definitions []Action `json:"definitions"`
}

// Action defines a single executable operation in the system.
//
// Actions are referenced by UI components but executed by backend services.
type Action struct {
    ID string `json:"id"`

    Type string `json:"type"`

    InputSchema map[string]string `json:"input_schema"`

    AuthRequired bool `json:"auth_required"`
}

// Permissions control:
//   - who can modify schema
//   - who can execute actions
//   - who can access admin features
// Permissions defines role-based access control configuration for a tenant.
type Permissions struct {
    Roles []Role `json:"roles"`
}

// Role defines a user role and its associated permissions within a tenant.
type Role struct {
    Name string `json:"name"`

    Permissions []string `json:"permissions"`
}