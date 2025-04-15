package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/GlebRadaev/password-manager/client/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalStorage(t *testing.T) {
	tempDir := t.TempDir()

	storage := &LocalStorage{
		path:         filepath.Join(tempDir, "data"),
		syncFilePath: filepath.Join(tempDir, "data", ".sync_status"),
	}

	testData := []byte("encrypted-data")
	now := time.Now().Unix()

	t.Run("Add and Get", func(t *testing.T) {
		entry := &models.DataEntry{
			ID:        "test1",
			Type:      models.Login,
			Data:      testData,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := storage.Add(entry)
		require.NoError(t, err)

		retrieved, err := storage.Get("test1")
		require.NoError(t, err)
		assert.Equal(t, entry, retrieved)
	})

	t.Run("Get non-existent", func(t *testing.T) {
		_, err := storage.Get("nonexistent")
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("GetAll", func(t *testing.T) {
		entries := []*models.DataEntry{
			{
				ID:        "test2",
				Type:      models.Card,
				Data:      testData,
				CreatedAt: 123,
				UpdatedAt: 123,
			},
			{
				ID:        "test3",
				Type:      models.Note,
				Data:      testData,
				CreatedAt: 456,
				UpdatedAt: 456,
			},
		}

		for _, e := range entries {
			require.NoError(t, storage.Add(e))
		}

		all, err := storage.GetAll()
		require.NoError(t, err)
		assert.Len(t, all, 3)

		ids := make(map[string]bool)
		for _, e := range all {
			ids[e.ID] = true
		}
		assert.True(t, ids["test1"])
		assert.True(t, ids["test2"])
		assert.True(t, ids["test3"])
	})

	t.Run("Delete", func(t *testing.T) {
		err := storage.Delete("test1")
		require.NoError(t, err)

		_, err = storage.Get("test1")
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("Sync status", func(t *testing.T) {
		syncStorage := &LocalStorage{
			path:         filepath.Join(tempDir, "sync_data"),
			syncFilePath: filepath.Join(tempDir, "sync_data", ".sync_status"),
		}

		entry := &models.DataEntry{
			ID:        "sync_test",
			Type:      models.Login,
			Data:      testData,
			CreatedAt: 100,
			UpdatedAt: 100,
		}

		require.NoError(t, syncStorage.Add(entry))

		pending, err := syncStorage.GetPendingSyncEntries()
		require.NoError(t, err)
		require.Len(t, pending, 1, "Должна быть одна запись, ожидающая синхронизации")
		assert.Equal(t, "sync_test", pending[0].ID)

		require.NoError(t, syncStorage.UpdateSyncStatus([]*models.DataEntry{entry}))

		pending, err = syncStorage.GetPendingSyncEntries()
		require.NoError(t, err)
		assert.Empty(t, pending, "После синхронизации не должно быть pending записей")

		require.NoError(t, syncStorage.ClearPendingSync())
	})

	t.Run("File permissions", func(t *testing.T) {
		entry := &models.DataEntry{
			ID:        "perm_test",
			Type:      models.Login,
			Data:      testData,
			CreatedAt: now,
			UpdatedAt: now,
		}

		require.NoError(t, storage.Add(entry))

		filePath := filepath.Join(storage.path, "perm_test")
		stat, err := os.Stat(filePath)
		require.NoError(t, err)

		assert.Equal(t, os.FileMode(0600), stat.Mode().Perm())
	})

	t.Run("Corrupted data", func(t *testing.T) {
		badFile := filepath.Join(storage.path, "corrupted")
		require.NoError(t, os.WriteFile(badFile, []byte("not a json"), 0600))

		all, err := storage.GetAll()
		require.NoError(t, err)

		for _, e := range all {
			assert.NotEqual(t, "corrupted", e.ID)
		}
	})
}

func TestGetAuthToken(t *testing.T) {
	tempDir := t.TempDir()
	tokenPath := filepath.Join(tempDir, ".pm_token")

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	storage := &LocalStorage{}

	t.Run("No token file", func(t *testing.T) {
		token, err := storage.GetAuthToken()
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("With token file", func(t *testing.T) {
		testToken := "test-token-123"
		require.NoError(t, os.WriteFile(tokenPath, []byte(testToken), 0600))

		token, err := storage.GetAuthToken()
		require.NoError(t, err)
		assert.Equal(t, testToken, token)
	})

	t.Run("File permissions", func(t *testing.T) {
		testToken := "test-token-456"
		require.NoError(t, os.WriteFile(tokenPath, []byte(testToken), 0600))

		stat, err := os.Stat(tokenPath)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0600), stat.Mode().Perm())
	})
}
