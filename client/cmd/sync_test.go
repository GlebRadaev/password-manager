package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/client/models"
)

func TestSyncCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSyncService := NewMockSyncServiceInterface(ctrl)
	originalSyncService := syncService
	syncService = mockSyncService
	defer func() { syncService = originalSyncService }()

	t.Run("sync without conflicts", func(t *testing.T) {
		mockSyncService.EXPECT().
			Sync().
			Return(&models.SyncResponse{}, nil)

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cmd := &cobra.Command{Use: "pm"}
		cmd.AddCommand(syncCmd)
		cmd.SetArgs([]string{"sync"})
		err := cmd.Execute()

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = oldStdout

		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Sync completed successfully")
	})

	t.Run("sync with conflicts", func(t *testing.T) {
		conflicts := []models.Conflict{
			{ConflictID: "conflict1", DataID: "data1"},
			{ConflictID: "conflict2", DataID: "data2"},
		}
		mockSyncService.EXPECT().
			Sync().
			Return(&models.SyncResponse{Conflicts: conflicts}, nil)

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cmd := &cobra.Command{Use: "pm"}
		cmd.AddCommand(syncCmd)
		cmd.SetArgs([]string{"sync"})
		err := cmd.Execute()

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = oldStdout

		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Found 2 conflicts:")
		assert.Contains(t, buf.String(), "data1")
		assert.Contains(t, buf.String(), "data2")
	})
}

func TestResolveCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSyncService := NewMockSyncServiceInterface(ctrl)
	originalSyncService := syncService
	syncService = mockSyncService
	defer func() { syncService = originalSyncService }()

	testConflictID := "conflict123"
	testStrategy := "server"
	testResponse := &models.ResolutionResponse{Message: "Conflict resolved successfully"}

	mockSyncService.EXPECT().
		Resolve(testConflictID, testStrategy).
		Return(testResponse, nil)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "pm"}
	cmd.AddCommand(resolveCmd)
	cmd.SetArgs([]string{"resolve", testConflictID, "--strategy", testStrategy})
	err := cmd.Execute()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Conflict resolved: "+testResponse.Message)
}
