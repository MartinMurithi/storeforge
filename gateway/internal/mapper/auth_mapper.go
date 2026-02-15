package mapper

import (
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/user"
	authv1 "github.com/MartinMurithi/storeforge/api/protos/auth/v1"
)

func ProtoLoginResponseToDTO(res *authv1.LoginResponse) *user.LoginResponse {
	if res == nil {
		return nil
	}

	return &user.LoginResponse{
		User:  ProtoUserToDTO(res.User),
		Token: ProtoTokenToDTO(res.Token),
	}
}

func ProtoRegisterResponseToDTO(res *authv1.RegisterResponse) *dto.RegisterResponse {
	if res == nil {
		return nil
	}

	return &dto.RegisterResponse{
		User:    ProtoUserToDTO(res.User),
		Message: res.Message,
	}
}


func ProtoTokenToDTO(t *authv1.Token) *dto.TokenResponse {
	if t == nil {
		return nil
	}

	return &dto.TokenResponse{
		AccessToken: t.AccessToken,
		TokenType:   t.TokenType,
		ExpiresIn:   t.ExpiresIn,
		ExpiresAt:   t.ExpiresAt.AsTime(),
		IssuedAt:    t.IssuedAt.AsTime(),
	}
}