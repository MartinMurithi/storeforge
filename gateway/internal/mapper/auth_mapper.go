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

// MapRefreshTokenResponseProtoToSafeDTO adapts the gRPC refresh result for the frontend.
// It explicitly excludes the RefreshToken string since that will be handled via HttpOnly cookies.
func MapRefreshTokenResponseProtoToSafeDTO(pb *authv1.RefreshTokenResponse) *dto.RefreshTokenResponseDTO {
    if pb == nil || pb.Token == nil {
        return nil
    }

    return &dto.RefreshTokenResponseDTO{
        Token: dto.TokenDTO{
            AccessToken: pb.Token.AccessToken,
            ExpiresIn:   pb.Token.ExpiresIn,
			ExpiresAt: pb.Token.ExpiresAt.AsTime(),
			IssuedAt: pb.Token.IssuedAt.AsTime(),
            TokenType:   "Bearer",
        },
    }
}