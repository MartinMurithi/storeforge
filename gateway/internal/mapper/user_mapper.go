package mapper

import (
	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/user"
)

// MapUserProtoToDTO transforms a gRPC User message into a public UserResponse DTO.
// This is used to present user data to the frontend in a clean JSON format.
func MapUserProtoToDTO(u *userv1.User) *user.UserResponse {
    if u == nil {
        return nil
    }

    dto := &user.UserResponse{
        ID:         u.Id,
        Email:      u.Email,
        IsVerified: u.IsVerified,
        CreatedAt:  u.CreatedAt.AsTime(),
        UpdatedAt:  u.UpdatedAt.AsTime(),
    }

    if u.Profile != nil {
        dto.Profile = user.UserProfileDTO{
            FullName:     u.Profile.FullName,
            Phone:        u.Profile.Phone,
            BusinessName: u.Profile.BusinessName,
            BusinessType: u.Profile.BusinessType,
        }
    }

    // Handle optional DeletedAt from Proto if it exists
    if u.DeletedAt != nil {
        t := u.DeletedAt.AsTime()
        dto.DeletedAt = &t
    }

    return dto
}

// MapUserProtosToDTOs converts a slice of gRPC User messages into a slice of UserResponse DTOs.
// Useful for list endpoints like GetAllUsers.
func MapUserProtosToDTOs(users []*userv1.User) []user.UserResponse {
    if users == nil {
        return []user.UserResponse{}
    }
    
    out := make([]user.UserResponse, 0, len(users))
	
    for _, u := range users {
        if dtoUser := MapUserProtoToDTO(u); dtoUser != nil {
            out = append(out, *dtoUser)
        }
    }
    return out
}