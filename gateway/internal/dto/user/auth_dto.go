package user

type RegisterUserRequestDTO struct {
    FullName     string `json:"fullName"`
    Email        string `json:"email"`
    Phone        string `json:"phone"`
    Password     string `json:"password"`
    BusinessType string `json:"businessType"`
    BusinessName string `json:"businessName"`
}

type RegisterUserResponseDTO struct {
    User    *UserResponseDTO `json:"user"`
    Message string           `json:"message"`
}

type LoginUserRequestDTO struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginUserResponseDTO struct {
    User  *UserResponseDTO `json:"user"`
    Token string           `json:"token"` // string for JSON-friendly token
}

type PatchUserRequestDTO struct {
    BusinessName *string `json:"businessName,omitempty"`
    BusinessType *string `json:"businessType,omitempty"`
}