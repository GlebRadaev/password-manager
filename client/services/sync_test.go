package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GlebRadaev/password-manager/client/models"
)

func TestSyncService_Sync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().Unix()
	testData := []byte("test-data")

	tests := []struct {
		name           string
		setupMocks     func(*MockStorageInterface, *MockHTTPClientInterface)
		expectedResult *models.SyncResponse
		expectedError  string
	}{
		{
			name: "successful sync with no conflicts",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)
				storage.EXPECT().GetPendingSyncEntries().Return([]*models.DataEntry{
					{
						ID:        "entry1",
						Type:      models.Login,
						Data:      testData,
						UpdatedAt: now,
					},
				}, nil)

				validateResp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "valid": true,
                        "UserID": "user123"
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					if strings.HasSuffix(req.URL.Path, "/validate-token") {
						return validateResp, nil
					}
					return nil, nil
				})

				syncResp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "success": true,
                        "conflicts": []
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(syncResp, nil)

				storage.EXPECT().UpdateSyncStatus(gomock.Any()).Return(nil)
			},
			expectedResult: &models.SyncResponse{
				Success:   true,
				Conflicts: []models.Conflict{},
			},
		},
		{
			name: "failed to get auth token",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("", errors.New("token error"))
			},
			expectedError: "authentication required: token error",
		},
		{
			name: "token validation failed",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("invalid-token", nil)
				client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("validation error"))
			},
			expectedError: "validation error",
		},
		{
			name: "failed to get pending entries",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				validateResp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "valid": true,
                        "UserID": "user123"
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(validateResp, nil)

				storage.EXPECT().GetPendingSyncEntries().Return(nil, errors.New("storage error"))
			},
			expectedError: "storage error",
		},
		{
			name: "sync with conflicts",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)
				storage.EXPECT().GetPendingSyncEntries().Return([]*models.DataEntry{
					{
						ID:        "entry1",
						Type:      models.Login,
						Data:      testData,
						UpdatedAt: now,
					},
				}, nil).Times(1)

				validateResp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "valid": true,
                        "UserID": "user123"
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					if strings.HasSuffix(req.URL.Path, "/validate-token") {
						return validateResp, nil
					}
					return nil, nil
				})

				syncResp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "success": true,
                        "conflicts": [{
                            "conflict_id": "conflict1",
                            "data_id": "entry1"
                        }]
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(syncResp, nil)
			},
			expectedResult: &models.SyncResponse{
				Success: true,
				Conflicts: []models.Conflict{
					{
						ConflictID: "conflict1",
						DataID:     "entry1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewMockStorageInterface(ctrl)
			httpClientMock := NewMockHTTPClientInterface(ctrl)

			tt.setupMocks(storageMock, httpClientMock)

			service := &SyncService{
				baseURL: "http://test",
				storage: storageMock,
				client:  httpClientMock,
			}

			result, err := service.Sync()

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

func TestSyncService_Resolve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMocks     func(*MockStorageInterface, *MockHTTPClientInterface)
		conflictID     string
		strategy       string
		expectedResult *models.ResolutionResponse
		expectedError  string
	}{
		{
			name: "successful resolution",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := httptest.NewRecorder()
				json.NewEncoder(resp).Encode(models.ResolutionResponse{
					Success: true,
					Message: "conflict resolved",
				})
				client.EXPECT().Do(gomock.Any()).Return(resp.Result(), nil)
			},
			conflictID: "conflict1",
			strategy:   "client",
			expectedResult: &models.ResolutionResponse{
				Success: true,
				Message: "conflict resolved",
			},
		},
		{
			name: "failed to get auth token",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("", errors.New("token error"))
			},
			conflictID:    "conflict1",
			strategy:      "client",
			expectedError: "authentication required: token error",
		},
		{
			name: "server error during resolution",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)
				client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("server error"))
			},
			conflictID:    "conflict1",
			strategy:      "client",
			expectedError: "request failed: server error",
		},
		{
			name: "server error with message",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusBadRequest,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "message": "invalid conflict id"
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			conflictID:    "invalid-conflict",
			strategy:      "client",
			expectedError: "server error: invalid conflict id",
		},
		{
			name: "unexpected status code without message",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			conflictID:    "conflict1",
			strategy:      "client",
			expectedError: "unexpected status code: 500",
		},
		{
			name: "failed to read response",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(&errorReader{}),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			conflictID:    "conflict1",
			strategy:      "client",
			expectedError: "failed to read response: simulated read error",
		},
		{
			name: "failed to decode response",
			setupMocks: func(storage *MockStorageInterface, client *MockHTTPClientInterface) {
				storage.EXPECT().GetAuthToken().Return("valid-token", nil)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			conflictID:    "conflict1",
			strategy:      "client",
			expectedError: "failed to decode response: invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewMockStorageInterface(ctrl)
			httpClientMock := NewMockHTTPClientInterface(ctrl)

			tt.setupMocks(storageMock, httpClientMock)

			service := &SyncService{
				baseURL: "http://test",
				storage: storageMock,
				client:  httpClientMock,
			}

			result, err := service.Resolve(tt.conflictID, tt.strategy)

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

func TestValidateTokenAndGetUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		setupMock     func(*MockHTTPClientInterface)
		token         string
		expectedID    string
		expectedError string
	}{
		{
			name: "valid token",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "valid": true,
                        "UserID": "user123"
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			token:      "valid-token",
			expectedID: "user123",
		},
		{
			name: "invalid token",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "valid": false
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			token:         "invalid-token",
			expectedError: "invalid token",
		},
		{
			name: "empty user ID",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "valid": true,
                        "UserID": ""
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			token:         "empty-user-token",
			expectedError: "server returned empty user_id",
		},
		{
			name: "validation request failed",
			setupMock: func(client *MockHTTPClientInterface) {
				client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("connection error"))
			},
			token:         "error-token",
			expectedError: "validation request failed: connection error",
		},
		{
			name: "invalid status code",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body: io.NopCloser(bytes.NewBufferString(`{
                        "message": "invalid token"
                    }`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			token:         "unauthorized-token",
			expectedError: "invalid token status: 401, response: {",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := NewMockHTTPClientInterface(ctrl)
			tt.setupMock(httpClientMock)

			service := &SyncService{
				baseURL: "http://test",
				client:  httpClientMock,
			}

			userID, err := service.validateTokenAndGetUserID(tt.token)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedID, userID)
		})
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}
