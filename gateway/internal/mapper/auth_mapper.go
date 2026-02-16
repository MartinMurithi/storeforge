package mapper

import (
	authv1 "github.com/MartinMurithi/storeforge/api/protos/auth/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/user"
)

// MapTokenProtoToDTO converts the gRPC Token message into a JSON-friendly TokenDTO.
func MapTokenProtoToDTO(t *authv1.Token) *user.TokenDTO {
	if t == nil {
		return &user.TokenDTO{}
	}

	return &user.TokenDTO{
		AccessToken: t.AccessToken,
		TokenType:   t.TokenType,
		ExpiresIn:   t.ExpiresIn,
		IssuedAt:    t.IssuedAt.AsTime(),
		ExpiresAt:   t.ExpiresAt.AsTime(),
	}
}

// MapRegisterResponseProtoToDTO adapts the Auth service registration result for the frontend.
func MapRegisterResponseProtoToDTO(pb *authv1.RegisterResponse) *user.RegisterResponse {
	if pb == nil {
		return nil
	}

	return &user.RegisterResponse{
		User:    *MapUserProtoToDTO(pb.User),
		Message: pb.Message,
	}
}

// MapLoginResponseProtoToDTO adapts the Auth service login result for the frontend,
// combining the identity (User) and the credentials (Token).
func MapLoginResponseProtoToDTO(pb *authv1.LoginResponse) *user.LoginResponse {
	if pb == nil {
		return nil
	}

	return &user.LoginResponse{
		User:  *MapUserProtoToDTO(pb.User),
		Token: *MapTokenProtoToDTO(pb.Token),
	}
}
