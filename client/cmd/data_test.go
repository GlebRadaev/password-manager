package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/client/models"
)

func TestAddCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataService := NewMockDataServiceInterface(ctrl)
	mockDataService.EXPECT().
		Add(gomock.Any()).
		Return(nil)

	originalDataService := dataService
	dataService = mockDataService
	defer func() { dataService = originalDataService }()

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add new data entry",
		Run: func(cmd *cobra.Command, args []string) {
			dataType, _ := cmd.Flags().GetString("type")
			content, _ := cmd.Flags().GetString("data")

			entry := &models.DataEntry{
				ID:        uuid.New().String(),
				Type:      models.DataTypeFromString(dataType),
				Data:      []byte(content),
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}

			if err := dataService.Add(entry); err != nil {
				cmd.PrintErrln("Add failed:", err)
				return
			}
			cmd.Printf("Added entry with ID: %s\n", entry.ID)
		},
	}
	addCmd.Flags().StringP("type", "t", "", "Entry type (login|note|card|binary)")
	addCmd.Flags().StringP("data", "d", "", "Entry content (JSON format for structured types)")

	buf := new(bytes.Buffer)
	addCmd.SetOut(buf)
	addCmd.SetErr(buf)

	addCmd.SetArgs([]string{"--type", "login", "--data", `{"username":"user","password":"pass"}`})
	err := addCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Added entry with ID:")
}

func TestListCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataService := NewMockDataServiceInterface(ctrl)
	testEntries := []*models.DataEntry{
		{
			ID:        uuid.NewString(),
			Type:      models.Login,
			UpdatedAt: time.Now().Unix(),
		},
		{
			ID:        uuid.NewString(),
			Type:      models.Note,
			UpdatedAt: time.Now().Unix(),
		},
	}
	mockDataService.EXPECT().
		List().
		Return(testEntries, nil)

	originalDataService := dataService
	dataService = mockDataService
	defer func() { dataService = originalDataService }()

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all entries",
		Run: func(cmd *cobra.Command, args []string) {
			entries, err := dataService.List()
			if err != nil {
				cmd.PrintErrln("List failed:", err)
				return
			}

			for i, e := range entries {
				cmd.Printf("%d. %s [%s] %s\n", i+1, e.ID, e.Type.String(),
					time.Unix(e.UpdatedAt, 0).Format("2006-01-02"))
			}
		},
	}

	buf := new(bytes.Buffer)
	listCmd.SetOut(buf)
	listCmd.SetErr(buf)

	listCmd.SetArgs([]string{})
	err := listCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), testEntries[0].ID)
	assert.Contains(t, buf.String(), testEntries[1].ID)
}

func TestViewCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataService := NewMockDataServiceInterface(ctrl)
	testID := uuid.NewString()
	testEntry := &models.DataEntry{
		ID:        testID,
		Type:      models.Login,
		Data:      []byte(`{"username":"user","password":"pass"}`),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	mockDataService.EXPECT().
		Get(testID).
		Return(testEntry, nil)

	originalDataService := dataService
	dataService = mockDataService
	defer func() { dataService = originalDataService }()

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View entry details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			entry, err := dataService.Get(args[0])
			if err != nil {
				cmd.PrintErrln("View failed:", err)
				return
			}

			cmd.Printf("ID: %s\n", entry.ID)
			cmd.Printf("Type: %s\n", entry.Type.String())
			cmd.Printf("Created: %s\n", time.Unix(entry.CreatedAt, 0).Format(time.RFC822))
			cmd.Printf("Updated: %s\n", time.Unix(entry.UpdatedAt, 0).Format(time.RFC822))
			cmd.Printf("Data: %s\n", string(entry.Data))
		},
	}

	buf := new(bytes.Buffer)
	viewCmd.SetOut(buf)
	viewCmd.SetErr(buf)

	viewCmd.SetArgs([]string{testID})
	err := viewCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "ID: "+testID)
	assert.Contains(t, buf.String(), "Type: login")
	assert.Contains(t, buf.String(), "Data: "+string(testEntry.Data))
}

func TestDeleteCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataService := NewMockDataServiceInterface(ctrl)
	testID := uuid.NewString()
	mockDataService.EXPECT().
		Delete(testID).
		Return(nil)

	originalDataService := dataService
	dataService = mockDataService
	defer func() { dataService = originalDataService }()

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete entry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := dataService.Delete(args[0]); err != nil {
				cmd.PrintErrln("Delete failed:", err)
				return
			}
			cmd.Println("Entry deleted")
		},
	}

	buf := new(bytes.Buffer)
	deleteCmd.SetOut(buf)
	deleteCmd.SetErr(buf)

	deleteCmd.SetArgs([]string{testID})
	err := deleteCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Entry deleted")
}
