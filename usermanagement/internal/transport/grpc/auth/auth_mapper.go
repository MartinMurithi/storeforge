package auth

import (
	"time"

	authv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/auth/v1"
	userv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/user/v1"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
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

// toProtoToken maps a domain Token entity into its protobuf representation.
//
// The Token message is a transport-level construct used by AuthService responses.
// The token type is fixed to "Bearer" as part of the authentication contract.
//
// Assumes the token has already been created and validated by the application layer.
func toProtoToken(t *entity.Token) *authv1.Token {
	if t == nil {
		return nil
	}

	return &authv1.Token{
		AccessToken: t.AccessToken,
		RefreshToken: t.RefreshToken,
		TokenType:   "Bearer",
		ExpiresIn:   t.ExpiresIn,
		IssuedAt:    timestamppb.New(t.IssuedAt),
		ExpiresAt:   timestamppb.New(t.ExpiresAt),
	}
}

// ToProtoRegisterResponse builds the gRPC RegisterResponse from a domain User.
//
// This mapper defines the exact wire response for a successful registration.
// The presence of the user implies that registration has already succeeded
// at the application layer.
//
// The success message is part of the API contract, not domain logic.
func ToProtoRegisterResponse(user *entity.User) *authv1.RegisterResponse {
	if user == nil {
		return nil
	}

	return &authv1.RegisterResponse{
		User:    ToProtoUser(user),
		Message: "Registration successful. Please verify your email.",
	}
}

// ToProtoLoginResponse builds the gRPC LoginResponse from domain User and Token entities.
//
// This function is called only after authentication has succeeded.
// It combines identity (User) and authorization (Token) into a single
// transport response as defined by the AuthService proto contract.
func ToProtoLoginResponse(user *entity.User, token *entity.Token) *authv1.LoginResponse {
	if user == nil || token == nil {
		return nil
	}

	return &authv1.LoginResponse{
		User:  ToProtoUser(user),
		Token: toProtoToken(token),
	}
}
