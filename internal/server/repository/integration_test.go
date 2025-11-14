//go:build integration

package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/etoneja/go-keeper/internal/server/stypes"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
	testDB   *pgxpool.Pool
	migrator *migrate.Migrate
)

func TestMain(m *testing.M) {
	dbURL := os.Getenv("GOKEEPER_TEST_DB_URL")
	dbURL = "postgres://gokeeper:gokeeper@127.0.0.1:5432/gokeepertest?sslmode=disable"
	if dbURL == "" {
		log.Fatal("GOKEEPER_TEST_DB_URL environment variable is required for integration tests")
	}

	var err error
	migrator, err = migrate.New("file://../../..//migrations", dbURL)
	if err != nil {
		log.Fatal("Failed to create migrator:", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration up failed:", err)
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}
	testDB = db

	code := m.Run()

	if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Warning: failed to rollback migrations: %v", err)
	}

	db.Close()
	migrator.Close()
	os.Exit(code)
}

func getTestDB(t *testing.T) *pgxpool.Pool {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}
	return testDB
}

type testUser struct {
	ID       string
	Login    string
	Password string
}

type testSecret struct {
	ID     string
	UserID string
	Data   []byte
	Hash   string
}

func createTestUser(t *testing.T, userRepo *UserRepository, login, password string) *testUser {
	ctx := context.Background()
	db := getTestDB(t)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user, err := userRepo.CreateUser(ctx, db, login, string(passwordHash))
	require.NoError(t, err)

	return &testUser{
		ID:       user.ID,
		Login:    user.Login,
		Password: password,
	}
}

func generateTestID(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), os.Getpid())
}

func TestSecretRepository_Integration(t *testing.T) {
	db := getTestDB(t)
	userRepo := NewUserRepository()
	secretRepo := NewSecretRepository()

	t.Run("SetSecret and GetSecret", func(t *testing.T) {
		ctx := context.Background()

		user := createTestUser(t, userRepo, generateTestID("user"), "password")
		secretID := generateTestID("secret")

		secret := &stypes.Secret{
			ID:           secretID,
			UserID:       user.ID,
			Data:         []byte("secret data"),
			Hash:         generateTestID("hash"),
			LastModified: time.Now(),
		}

		err := secretRepo.SetSecret(ctx, db, secret)
		require.NoError(t, err)

		retrieved, err := secretRepo.GetSecret(ctx, db, user.ID, secretID)
		require.NoError(t, err)
		assert.Equal(t, secret.ID, retrieved.ID)
		assert.Equal(t, secret.UserID, retrieved.UserID)
		assert.Equal(t, secret.Data, retrieved.Data)
		assert.Equal(t, secret.Hash, retrieved.Hash)
	})

	t.Run("GetSecret not found", func(t *testing.T) {
		ctx := context.Background()

		user := createTestUser(t, userRepo, generateTestID("user"), "password")

		secret, err := secretRepo.GetSecret(ctx, db, user.ID, generateTestID("nonexistent"))
		require.Error(t, err)
		assert.Nil(t, secret)
		assert.ErrorIs(t, err, ErrSecretNotFound)
	})

	t.Run("UpdateSecret", func(t *testing.T) {
		ctx := context.Background()

		user := createTestUser(t, userRepo, generateTestID("user"), "password")
		secretID := generateTestID("secret")

		secret1 := &stypes.Secret{
			ID:           secretID,
			UserID:       user.ID,
			Data:         []byte("data v1"),
			Hash:         "hash1",
			LastModified: time.Now(),
		}

		err := secretRepo.SetSecret(ctx, db, secret1)
		require.NoError(t, err)

		secret2 := &stypes.Secret{
			ID:           secretID,
			UserID:       user.ID,
			Data:         []byte("data v2"),
			Hash:         "hash2",
			LastModified: time.Now().Add(time.Hour),
		}

		err = secretRepo.SetSecret(ctx, db, secret2)
		require.NoError(t, err)

		retrieved, err := secretRepo.GetSecret(ctx, db, user.ID, secretID)
		require.NoError(t, err)
		assert.Equal(t, []byte("data v2"), retrieved.Data)
		assert.Equal(t, "hash2", retrieved.Hash)
	})

	t.Run("DeleteSecret", func(t *testing.T) {
		ctx := context.Background()

		user := createTestUser(t, userRepo, generateTestID("user"), "password")
		secretID := generateTestID("secret")

		secret := &stypes.Secret{
			ID:           secretID,
			UserID:       user.ID,
			Data:         []byte("data"),
			Hash:         generateTestID("hash"),
			LastModified: time.Now(),
		}

		err := secretRepo.SetSecret(ctx, db, secret)
		require.NoError(t, err)

		err = secretRepo.DeleteSecret(ctx, db, user.ID, secretID)
		require.NoError(t, err)

		_, err = secretRepo.GetSecret(ctx, db, user.ID, secretID)
		assert.ErrorIs(t, err, ErrSecretNotFound)
	})

	t.Run("DeleteSecret not found", func(t *testing.T) {
		ctx := context.Background()

		user := createTestUser(t, userRepo, generateTestID("user"), "password")

		err := secretRepo.DeleteSecret(ctx, db, user.ID, generateTestID("nonexistent"))
		assert.ErrorIs(t, err, ErrSecretNotFound)
	})

	t.Run("ListSecrets user isolation", func(t *testing.T) {
		ctx := context.Background()

		user1 := createTestUser(t, userRepo, generateTestID("user1"), "password")
		user2 := createTestUser(t, userRepo, generateTestID("user2"), "password")

		secret1 := &stypes.Secret{
			ID:           generateTestID("secret"),
			UserID:       user1.ID,
			Data:         []byte("data1"),
			Hash:         generateTestID("hash"),
			LastModified: time.Now(),
		}
		secret2 := &stypes.Secret{
			ID:           generateTestID("secret"),
			UserID:       user1.ID,
			Data:         []byte("data2"),
			Hash:         generateTestID("hash"),
			LastModified: time.Now().Add(-time.Hour),
		}
		secret3 := &stypes.Secret{
			ID:           generateTestID("secret"),
			UserID:       user2.ID,
			Data:         []byte("data3"),
			Hash:         generateTestID("hash"),
			LastModified: time.Now(),
		}

		secretRepo.SetSecret(ctx, db, secret1)
		secretRepo.SetSecret(ctx, db, secret2)
		secretRepo.SetSecret(ctx, db, secret3)

		user1Secrets, err := secretRepo.ListSecrets(ctx, db, user1.ID)
		require.NoError(t, err)
		assert.Len(t, user1Secrets, 2)

		user2Secrets, err := secretRepo.ListSecrets(ctx, db, user2.ID)
		require.NoError(t, err)
		assert.Len(t, user2Secrets, 1)
	})

	t.Run("ListSecrets empty", func(t *testing.T) {
		ctx := context.Background()

		user := createTestUser(t, userRepo, generateTestID("user"), "password")

		secrets, err := secretRepo.ListSecrets(ctx, db, user.ID)
		require.NoError(t, err)
		assert.Empty(t, secrets)
	})

	t.Run("ListSecrets order by last_modified desc", func(t *testing.T) {
		ctx := context.Background()

		user := createTestUser(t, userRepo, generateTestID("user"), "password")

		secret1 := &stypes.Secret{
			ID:           generateTestID("secret"),
			UserID:       user.ID,
			Data:         []byte("data1"),
			Hash:         "hash1",
			LastModified: time.Now().Add(-2 * time.Hour),
		}
		secret2 := &stypes.Secret{
			ID:           generateTestID("secret"),
			UserID:       user.ID,
			Data:         []byte("data2"),
			Hash:         "hash2",
			LastModified: time.Now().Add(-1 * time.Hour),
		}
		secret3 := &stypes.Secret{
			ID:           generateTestID("secret"),
			UserID:       user.ID,
			Data:         []byte("data3"),
			Hash:         "hash3",
			LastModified: time.Now(),
		}

		secretRepo.SetSecret(ctx, db, secret1)
		secretRepo.SetSecret(ctx, db, secret2)
		secretRepo.SetSecret(ctx, db, secret3)

		secrets, err := secretRepo.ListSecrets(ctx, db, user.ID)
		require.NoError(t, err)
		require.Len(t, secrets, 3)
		assert.Equal(t, "hash3", secrets[0].Hash)
		assert.Equal(t, "hash2", secrets[1].Hash)
		assert.Equal(t, "hash1", secrets[2].Hash)
	})
}
