package cmd

import (
	"bytes"
	"io"
	"os"
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
	originalDataService := dataService
	dataService = mockDataService
	defer func() { dataService = originalDataService }()

	mockDataService.EXPECT().
		Add(gomock.Any()).
		Return(nil)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "pm"}
	cmd.AddCommand(addCmd)
	cmd.SetArgs([]string{"add", "--type", "login", "--data", `{"username":"user","password":"pass"}`})
	err := cmd.Execute()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Added entry with ID:")
}

func TestListCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataService := NewMockDataServiceInterface(ctrl)
	originalDataService := dataService
	dataService = mockDataService
	defer func() { dataService = originalDataService }()

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

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "pm"}
	cmd.AddCommand(listCmd)
	cmd.SetArgs([]string{"list"})
	err := cmd.Execute()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), testEntries[0].ID)
	assert.Contains(t, buf.String(), testEntries[1].ID)
}

func TestViewCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataService := NewMockDataServiceInterface(ctrl)
	originalDataService := dataService
	dataService = mockDataService
	defer func() { dataService = originalDataService }()

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

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "pm"}
	cmd.AddCommand(viewCmd)
	cmd.SetArgs([]string{"view", testID})
	err := cmd.Execute()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "ID: "+testID)
	assert.Contains(t, buf.String(), "Type: login")
	assert.Contains(t, buf.String(), "Data: "+string(testEntry.Data))
}

func TestDeleteCmd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataService := NewMockDataServiceInterface(ctrl)
	originalDataService := dataService
	dataService = mockDataService
	defer func() { dataService = originalDataService }()

	testID := uuid.NewString()

	mockDataService.EXPECT().
		Delete(testID).
		Return(nil)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "pm"}
	cmd.AddCommand(deleteCmd)
	cmd.SetArgs([]string{"delete", testID})
	err := cmd.Execute()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Entry deleted")
}
