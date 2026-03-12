package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/config"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database/postgres"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var repo repository.IAuthRepository

func TestMain(m *testing.M) {
	ctx := context.Background()

	cfg, err := config.Load()

	// Setup real local DB connection
	// Ensure config points to local postgres
	pool, err := postgres.Connect(ctx, &cfg.DB)
	if err != nil {
		panic(err)
	}

	dbAdapter := postgres.NewAdapter(pool.Pool)

	repo = repository.NewUAuthRepository(dbAdapter)

	os.Exit(m.Run())
}

func TestUpdateSctiveSessionContext(t *testing.T) {

	targetUserId := pgtype.UUID{Bytes: uuid.MustParse("567c50eb-1e56-4a6b-aeb3-705801cb0162"), Valid: true}
	targetTenatId := pgtype.UUID{Bytes: uuid.New(), Valid: false}
	targetRole := "owner"

	testInput := struct {
		UserId   pgtype.UUID
		TenantId pgtype.UUID
		Role     string
	}{
		UserId:   targetUserId,
		TenantId: targetTenatId,
		Role:     "",
	}

	// Create a Global session first (Tenant Id and role is Nil)
	initialToken := &entity.RefreshToken{
		UserId:    testInput.UserId,
		LastRole:  testInput.Role,
		ExpiresAt: time.Now().Add(time.Hour),
		TokenHash: "token-hash162",
	}

	err := repo.CreateRefreshToken(context.Background(), initialToken)
	assert.NoError(t, err)

	testInput.TenantId = pgtype.UUID{Bytes: uuid.MustParse("c6731dc2-89a2-4c01-8f2f-ceb4f9d8ac6f"), Valid: true}

	// Update the active session
	updatedSession, err := repo.UpdateActiveSessionContext(context.Background(), testInput.UserId, testInput.TenantId, targetRole)

	// assertions
	require.NoError(t, err)
	require.NotNil(t, updatedSession)
	assert.Equal(t, targetRole, updatedSession.LastRole)
	assert.True(t, updatedSession.LastTenantId.Valid)
	assert.Equal(t, testInput.TenantId.Bytes, updatedSession.LastTenantId.Bytes)

}
