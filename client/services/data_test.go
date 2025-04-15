package services

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/client/models"
)

func TestDataService_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		setupMock     func(*MockStorageInterface)
		entry         *models.DataEntry
		expectedError string
	}{
		{
			name: "successful add",
			setupMock: func(storage *MockStorageInterface) {
				storage.EXPECT().Add(gomock.Any()).Return(nil)
			},
			entry: &models.DataEntry{ID: "test1"},
		},
		{
			name: "failed to add",
			setupMock: func(storage *MockStorageInterface) {
				storage.EXPECT().Add(gomock.Any()).Return(errors.New("storage error"))
			},
			entry:         &models.DataEntry{ID: "test1"},
			expectedError: "failed to save locally: storage error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewMockStorageInterface(ctrl)
			tt.setupMock(storageMock)

			service := &DataService{
				storage: storageMock,
			}

			err := service.Add(tt.entry)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDataService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().Unix()
	testData := []*models.DataEntry{
		{
			ID:        "test1",
			Type:      models.Login,
			Data:      []byte("data1"),
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name           string
		setupMock      func(*MockStorageInterface)
		expectedResult []*models.DataEntry
		expectedError  string
	}{
		{
			name: "successful list",
			setupMock: func(storage *MockStorageInterface) {
				storage.EXPECT().GetAll().Return(testData, nil)
			},
			expectedResult: testData,
		},
		{
			name: "failed to list",
			setupMock: func(storage *MockStorageInterface) {
				storage.EXPECT().GetAll().Return(nil, errors.New("storage error"))
			},
			expectedError: "failed to get local data: storage error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewMockStorageInterface(ctrl)
			tt.setupMock(storageMock)

			service := &DataService{
				storage: storageMock,
			}

			result, err := service.List()

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestDataService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testEntry := &models.DataEntry{
		ID:   "test1",
		Type: models.Login,
		Data: []byte("test data"),
	}

	tests := []struct {
		name           string
		setupMock      func(*MockStorageInterface)
		id             string
		expectedResult *models.DataEntry
		expectedError  error
	}{
		{
			name: "successful get",
			setupMock: func(storage *MockStorageInterface) {
				storage.EXPECT().Get("test1").Return(testEntry, nil)
			},
			id:             "test1",
			expectedResult: testEntry,
		},
		{
			name: "not found",
			setupMock: func(storage *MockStorageInterface) {
				storage.EXPECT().Get("missing").Return(nil, errors.New("not found"))
			},
			id:            "missing",
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewMockStorageInterface(ctrl)
			tt.setupMock(storageMock)

			service := &DataService{
				storage: storageMock,
			}

			result, err := service.Get(tt.id)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestDataService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		setupMock     func(*MockStorageInterface)
		id            string
		expectedError string
	}{
		{
			name: "successful delete",
			setupMock: func(storage *MockStorageInterface) {
				storage.EXPECT().Delete("test1").Return(nil)
			},
			id: "test1",
		},
		{
			name: "failed to delete",
			setupMock: func(storage *MockStorageInterface) {
				storage.EXPECT().Delete("missing").Return(errors.New("not found"))
			},
			id:            "missing",
			expectedError: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewMockStorageInterface(ctrl)
			tt.setupMock(storageMock)

			service := &DataService{
				storage: storageMock,
			}

			err := service.Delete(tt.id)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDataService_SyncWithServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().Unix()
	testData := []*models.DataEntry{
		{
			ID:        "test1",
			Type:      models.Login,
			Data:      []byte("data1"),
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name          string
		setupMocks    func(*MockStorageInterface, *MockHTTPClientInterface)
		expectedError string
	}{
		{
			name: "successful sync",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAll().Return(testData, nil)
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
		{
			name: "failed to get entries",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAll().Return(nil, errors.New("storage error"))
			},
			expectedError: "storage error",
		},
		{
			name: "failed to get auth token",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAll().Return(testData, nil)
				storage.EXPECT().GetAuthToken().Return("", errors.New("token error"))
			},
			expectedError: "authentication required: token error",
		},
		{
			name: "request failed",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAll().Return(testData, nil)
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)
				client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("http error"))
			},
			expectedError: "request failed: http error",
		},
		{
			name: "failed to read response",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAll().Return(testData, nil)
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(&errorReader{}),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedError: "failed to read response: simulated read error",
		},
		{
			name: "server error with message",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAll().Return(testData, nil)
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewBufferString(`{"message": "invalid data"}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedError: "server error: invalid data",
		},
		{
			name: "unexpected status code",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAll().Return(testData, nil)
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedError: "unexpected status code: 500",
		},
		{
			name: "failed to decode response",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAll().Return(testData, nil)
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedError: "failed to decode response: invalid character",
		},
		{
			name: "sync failed",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAll().Return(testData, nil)
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"success": false}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedError: "sync failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewMockStorageInterface(ctrl)
			httpClientMock := NewMockHTTPClientInterface(ctrl)

			tt.setupMocks(storageMock, httpClientMock)

			service := &DataService{
				storage: storageMock,
				client:  httpClientMock,
				baseURL: "http://test",
			}

			err := service.SyncWithServer()

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}
