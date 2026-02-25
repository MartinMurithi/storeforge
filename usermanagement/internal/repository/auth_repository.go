package repository

import (
	"context"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/database"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database/postgres"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
)

type AuthRepository struct {
	DB database.DB
}

type IAuthRepository interface {
	CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id string) (*entity.RefreshToken, error)
}

func NewUAuthRepository(db database.DB) IAuthRepository {
	return &AuthRepository{DB: db}
}

func (repo *AuthRepository) CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	const op = "auth.createRefreshToken"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked)
        VALUES ($1,$2,$3,$4,$5)
		RETURNING (id, token_hash)
    `
	err := repo.DB.QueryRow(ctx, query,
		token.Id,
		token.UserId,
		token.TokenHash,
		token.ExpiresAt,
		token.Revoked,
	).Scan(&token.Id, &token.TokenHash)

	if err != nil {
		log.Printf("[%s]: error creating refresh token: %v", op, err)

		// Map DB/driver errors → infra errors
		infraErr := postgres.MapPostgresError(err)

		// Translate infra → domain errors
		return TranslateUserRepoError(infraErr)
	}
	return nil
}

// GetRefreshTokenByHash implements [IAuthRepository].
func (repo *AuthRepository) GetRefreshTokenByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	const op = "authRepo.GetRefreshTokenByHash"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
	SELECT id, user_id, token_hash, expires_at, revoked FROM refresh_tokens 
	WHERE token_hash = $1
	`

	refreshToken := &entity.RefreshToken{}

	err := repo.DB.QueryRow(ctx, query).Scan(
		&refreshToken.Id,
		&refreshToken.UserId,
		&refreshToken.TokenHash,
		&refreshToken.Revoked,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		log.Printf("[%s]: error getting refresh token: %v", op, err)

		// Map DB/driver errors → infra errors
		infraErr := postgres.MapPostgresError(err)

		// Translate infra → domain errors
		return nil, TranslateUserRepoError(infraErr)
	}

	return refreshToken, nil

}

// RevokeRefreshToken implements [IAuthRepository].
func (repo *AuthRepository) RevokeRefreshToken(ctx context.Context, id string) (*entity.RefreshToken, error) {
	const op = "auth.revokeRefreshToken"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `	UPDATE refresh_tokens SET revoked = TRUE WHERE id = $1`

	err := repo.DB.QueryRow(ctx, query, id).Scan()

	if err != nil {
		log.Printf("[%s]: error getting refresh token: %v", op, err)

		// Map DB/driver errors → infra errors
		infraErr := postgres.MapPostgresError(err)

		// Translate infra → domain errors
		return nil, TranslateUserRepoError(infraErr)
	}
	return nil, nil
}
