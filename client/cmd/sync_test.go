package cmd

import (
	"bytes"
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

		syncCmd := &cobra.Command{
			Use: "sync",
			RunE: func(cmd *cobra.Command, args []string) error {
				resp, err := syncService.Sync()
				if err != nil {
					return err
				}
				if len(resp.Conflicts) == 0 {
					cmd.Println("Sync completed successfully")
				}
				return nil
			},
		}

		buf := new(bytes.Buffer)
		syncCmd.SetOut(buf)
		syncCmd.SetErr(buf)

		err := syncCmd.Execute()
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

		syncCmd := &cobra.Command{
			Use: "sync",
			RunE: func(cmd *cobra.Command, args []string) error {
				resp, err := syncService.Sync()
				if err != nil {
					return err
				}
				if len(resp.Conflicts) > 0 {
					cmd.Printf("Found %d conflicts:\n", len(resp.Conflicts))
					for _, c := range resp.Conflicts {
						cmd.Printf("- %s (ID: %s)\n", c.DataID, c.ConflictID)
					}
				}
				return nil
			},
		}

		buf := new(bytes.Buffer)
		syncCmd.SetOut(buf)
		syncCmd.SetErr(buf)

		err := syncCmd.Execute()
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

	resolveCmd := &cobra.Command{
		Use: "resolve",
		RunE: func(cmd *cobra.Command, args []string) error {
			strategy, _ := cmd.Flags().GetString("strategy")
			resp, err := syncService.Resolve(args[0], strategy)
			if err != nil {
				return err
			}
			cmd.Println("Conflict resolved:", resp.Message)
			return nil
		},
	}
	resolveCmd.Flags().StringP("strategy", "s", "", "Resolution strategy")

	buf := new(bytes.Buffer)
	resolveCmd.SetOut(buf)
	resolveCmd.SetErr(buf)

	resolveCmd.SetArgs([]string{testConflictID, "--strategy", testStrategy})
	err := resolveCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Conflict resolved: "+testResponse.Message)
}
