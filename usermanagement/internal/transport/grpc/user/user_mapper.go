package user

import (
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toProtoUser maps a domain User entity to its gRPC protobuf representation.
//
// This function performs a pure transformation:
//
// It exists strictly at the transport boundary to adapt the domain model
// to the AuthService gRPC contract.
//
// Returns nil if the input user is nil.
func ToProtoUser(u *entity.User) *userv1.User {
	if u == nil {
		return nil
	}

	return &userv1.User{
		Id:         u.ID.String(),
		Email:      u.Email,
		IsVerified: u.IsVerified,
		Profile: &userv1.UserProfile{
			FullName:     u.FullName,
			Phone:        u.Phone,
			BusinessName: u.BusinessName,
			BusinessType: u.BusinessType,
		},
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: toProtoTimestamp(u.UpdatedAt),
		DeletedAt: toProtoTimestamp(u.DeletedAt),
	}
}

// toProtoTimestamp converts an optional time.Time pointer into a protobuf Timestamp.
//
// A nil time value is preserved as nil, allowing optional timestamps
// (e.g. updated_at, deleted_at) to remain unset in the wire representation.
//
// This helper avoids leaking protobuf concerns into the domain layer.
func toProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

// ToProtoPaginationMeta converts internal pagination metadata to gRPC proto.
func ToProtoPaginationMeta(meta *dto.PaginationMeta) *userv1.PaginationMeta {
	if meta == nil {
		return nil
	}

	return &userv1.PaginationMeta{
		Page:       int32(meta.Page),
		Limit:      int32(meta.Limit),
		Total:      int32(meta.Total),
		TotalPages: int32(meta.TotalPages),
		HasNext:    meta.HasNext,
		HasPrev:    meta.HasPrev,
	}
}

// ToProtoFetchAllUsersResponse maps service-layer users + pagination → gRPC response
func ToProtoFetchAllUsersResponse(users []*entity.User, meta *dto.PaginationMeta) *userv1.GetAllUsersResponse {
	protoUsers := make([]*userv1.User, len(users))

	for i, u := range users {
		protoUsers[i] = ToProtoUser(u)
	}

	return &userv1.GetAllUsersResponse{
		Users: protoUsers,
		Meta:  ToProtoPaginationMeta(meta),
	}
}

func ToProtoUpdateUserResponse(user *entity.User) *userv1.UpdateUserResponse {
	if user == nil {
		return nil
	}

	return &userv1.UpdateUserResponse{
		User: ToProtoUser(user),
	}
}

func ToProtoSoftDeleteUserResponse() *userv1.DeleteUserResponse {

	return &userv1.DeleteUserResponse{
		Message: "User deleted successfully",
	}
}
