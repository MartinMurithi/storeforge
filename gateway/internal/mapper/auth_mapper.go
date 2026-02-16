package mapper

import (
	authv1 "github.com/MartinMurithi/storeforge/api/protos/auth/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
)

// MapTokenProtoToDTO converts the gRPC Token message into a JSON-friendly TokenDTO.
func MapTokenProtoToDTO(t *authv1.Token) *dto.TokenDTO {
	if t == nil {
		return &dto.TokenDTO{}
	}

	return &dto.TokenDTO{
		AccessToken: t.AccessToken,
		TokenType:   t.TokenType,
		ExpiresIn:   t.ExpiresIn,
		IssuedAt:    t.IssuedAt.AsTime(),
		ExpiresAt:   t.ExpiresAt.AsTime(),
	}
}

// MapRegisterResponseProtoToDTO adapts the Auth service registration result for the frontend.
func MapRegisterResponseProtoToDTO(pb *authv1.RegisterResponse) *dto.RegisterResponseDTO {
	if pb == nil {
		return nil
	}

	return &dto.RegisterResponseDTO{
		User:    *MapUserProtoToDTO(pb.User),
		Message: pb.Message,
	}
}

// MapLoginResponseProtoToDTO adapts the Auth service login result for the frontend,
// combining the identity (User) and the credentials (Token).
func MapLoginResponseProtoToDTO(pb *authv1.LoginResponse) *dto.LoginResponseDTO {
	if pb == nil {
		return nil
	}

	return &dto.LoginResponseDTO{
		User:  *MapUserProtoToDTO(pb.User),
		Token: *MapTokenProtoToDTO(pb.Token),
	}
}
