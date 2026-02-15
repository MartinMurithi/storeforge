package mapper

import (
	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/user"
)

func ProtoUserToDTO(u *userv1.User) *user.UserResponse {
	if u == nil {
		return nil
	}
	return &user.UserResponse{
		ID:         u.Id,
		Email:      u.Email,
		IsVerified: u.IsVerified,
		Profile: user.UserProfileDTO{
			FullName:     u.Profile.FullName,
			Phone:        u.Profile.Phone,
			BusinessName: u.Profile.BusinessName,
			BusinessType: u.Profile.BusinessType,
		},
		CreatedAt: u.CreatedAt.AsTime(),
		UpdatedAt: u.UpdatedAt.AsTime(),
	}
}

func ProtoUsersToDTO(users []*userv1.User) []user.UserResponse {
	out := make([]user.UserResponse, 0, len(users))

	for _, u := range users {
		dtoUser := ProtoUserToDTO(u)

		if dtoUser != nil {
			out = append(out, *dtoUser)
		}
	}

	return out
}
