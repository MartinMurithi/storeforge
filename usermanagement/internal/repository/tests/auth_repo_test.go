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
var permRepo repository.IPermissionRepository

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
	permRepo = repository.NewPermissionRepository(dbAdapter)

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

func TestGetPermissions(t *testing.T) {

	permissions, err := permRepo.GetPermissions(context.Background())
	assert.NoError(t, err)

	// assertions
	assert.NoError(t, err)
	assert.NotNil(t, permissions)
	assert.Equal(t, len(permissions), 11, "should have exactly 11 seeded permissions")

	var found = false

	for _, perm := range permissions {
		// Check if at least one expected permission exists in the slice
		assert.NotEqual(t, uuid.Nil, perm.Id, "each permission should have a generated uuid")

		if perm.Slug == "products:write" {
			found = true
			assert.NotEmpty(t, perm.Description)
		}
	}
	assert.True(t, found, "The 'products:write' permission should be present in the results")
}

func TestGetPermissionsById(t *testing.T) {
    ctx := context.Background()

    // get all permissions so we have valid IDs to work with
    allPerms, err := permRepo.GetPermissions(ctx)
    require.NoError(t, err)
    require.True(t, len(allPerms) >= 2, "Need at least 2 permissions in DB to run this test")

    // Pick the first two IDs
    permIds := []pgtype.UUID{
        {Bytes: allPerms[0].Id.Bytes, Valid: true},
        {Bytes: allPerms[1].Id.Bytes, Valid: true},
    }

    // test the GetByIDs function
    permissions, err := permRepo.GetPermissionsById(ctx, permIds)
    
    assert.NoError(t, err)
    require.NotNil(t, permissions)
    assert.Equal(t, 2, len(permissions), "Should return exactly the number of unique IDs requested")

    // 4. Verify the IDs returned match the ones we sent
    assert.Equal(t, allPerms[0].Id, permissions[0].Id)
    assert.Equal(t, allPerms[1].Id, permissions[1].Id)
}
