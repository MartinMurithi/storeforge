package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/database"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database/postgres"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/jackc/pgx/v5"
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
        INSERT INTO refresh_tokens (user_id, token_hash, expires_at, revoked)
        VALUES ($1,$2,$3,$4)
		RETURNING id, token_hash
    `
	err := repo.DB.QueryRow(ctx, query,
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

func (repo *AuthRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	const op = "authRepo.GetRefreshTokenByHash"

	if tokenHash == "" {
		return nil, fmt.Errorf("[%s] token hash is empty", op)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, token_hash, expires_at, revoked, created_at, revoked_at
		FROM refresh_tokens 
		WHERE token_hash = $1
	`

	refreshToken := &entity.RefreshToken{}

	log.Printf("Executing query with token_hash=%s", tokenHash)

	err := repo.DB.QueryRow(ctx, query, tokenHash).Scan(
		&refreshToken.Id,
		&refreshToken.UserId,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.Revoked,
		&refreshToken.CreatedAt,
		&refreshToken.RevokedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Not found
			return nil, fmt.Errorf("[%s] refresh token not found", op)
		}
		log.Printf("[%s] error getting refresh token: %v", op, err)
		infraErr := postgres.MapPostgresError(err)
		return nil, TranslateUserRepoError(infraErr)
	}

	return refreshToken, nil
}

func (repo *AuthRepository) RevokeRefreshToken(ctx context.Context, id string) (*entity.RefreshToken, error) {
    const op = "auth.revokeRefreshToken"

    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        UPDATE refresh_tokens 
        SET 
            revoked = TRUE,
            revoked_at = NOW()
        WHERE id = $1
        RETURNING id, user_id, token_hash, expires_at, revoked, created_at, revoked_at
    `

    refreshToken := &entity.RefreshToken{}

    err := repo.DB.QueryRow(ctx, query, id).Scan(
        &refreshToken.Id,
        &refreshToken.UserId,
        &refreshToken.TokenHash,
        &refreshToken.ExpiresAt,
        &refreshToken.Revoked,
        &refreshToken.CreatedAt,
        &refreshToken.RevokedAt,
    )

    if err != nil {
        log.Printf("[%s]: error revoking refresh token: %v", op, err)

        infraErr := postgres.MapPostgresError(err)
        return nil, TranslateUserRepoError(infraErr)
    }

    return refreshToken, nil
}
