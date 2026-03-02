package mapper

import (
	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
)

// MapUserProtoToDTO transforms a gRPC User message into a public UserResponse DTO.
// This is used to present user data to the frontend in a clean JSON format.
func MapUserProtoToDTO(u *userv1.User) *dto.UserResponseDTO {
	if u == nil {
		return nil
	}

	res := &dto.UserResponseDTO{
		ID:         u.Id,
		Email:      u.Email,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt.AsTime(),
		UpdatedAt:  u.UpdatedAt.AsTime(),
	}

	if u.Profile != nil {
		res.Profile = dto.UserProfileDTO{
			FullName:     u.Profile.FullName,
			Phone:        u.Profile.Phone,
		}
	}

	// Handle optional DeletedAt from Proto if it exists
	if u.DeletedAt != nil {
		t := u.DeletedAt.AsTime()
		res.DeletedAt = &t
	}

	return res
}

// MapUserProtosToDTOs converts a slice of gRPC User messages into a slice of UserResponse DTOs.
// Useful for list endpoints like GetAllUsers.
func MapUserProtosToDTOs(users []*userv1.User) []dto.UserResponseDTO {
	if users == nil {
		return []dto.UserResponseDTO{} // Returns empty slice of users
	}

	// Ensure 'out' matches the return signature type exactly
	out := make([]dto.UserResponseDTO, 0, len(users))

	for _, u := range users {
		if dtoUser := MapUserProtoToDTO(u); dtoUser != nil {
			out = append(out, *dtoUser)
		}
	}
	return out
}
