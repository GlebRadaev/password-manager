package services

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	reflect "reflect"
	"testing"

	"github.com/GlebRadaev/password-manager/client/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewAuthService(t *testing.T) {
	tests := []struct {
		name            string
		createService   func(t *testing.T) *AuthService
		expectedBaseURL string
		expectedClient  reflect.Type
		expectPath      bool
		expectedPath    string
		expectedError   string
	}{
		{
			name: "default initialization",
			createService: func(t *testing.T) *AuthService {
				return NewAuthService()
			},
			expectedBaseURL: "http://localhost:8079",
			expectedClient:  reflect.TypeOf(&http.Client{}),
			expectPath:      true,
		},
		{
			name: "custom tokenPath with error",
			createService: func(t *testing.T) *AuthService {
				return &AuthService{
					baseURL: "http://localhost:8079",
					client:  &http.Client{},
					tokenPath: func() (string, error) {
						return "", errors.New("path error")
					},
				}
			},
			expectedBaseURL: "http://localhost:8079",
			expectedClient:  reflect.TypeOf(&http.Client{}),
			expectPath:      true,
			expectedError:   "path error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := tt.createService(t)

			assert.Equal(t, tt.expectedBaseURL, service.baseURL, "baseURL mismatch")
			assert.Equal(t, tt.expectedClient, reflect.TypeOf(service.client), "client type mismatch")
			assert.NotNil(t, service.client, "client should not be nil")

			if tt.expectPath {
				path, err := service.tokenPath()

				if tt.expectedError != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedError)
				} else {
					require.NoError(t, err)
					if tt.expectedPath != "" {
						assert.Equal(t, tt.expectedPath, path, "tokenPath mismatch")
					} else {
						assert.Contains(t, path, tokenFileName, "tokenPath should contain token file name")
						assert.True(t, filepath.IsAbs(path), "tokenPath should be absolute")
					}
				}
			}
		})
	}
}
func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		setupMock     func(*MockHTTPClientInterface)
		username      string
		password      string
		email         string
		expected      *models.RegisterResponse
		expectedError string
	}{
		{
			name: "successful registration",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"user_id": "user123",
						"message": "User registered successfully"
					}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			username: "testuser",
			password: "testpass123",
			email:    "test@example.com",
			expected: &models.RegisterResponse{
				UserID:  "user123",
				Message: "User registered successfully",
			},
		},
		{
			name: "failed request",
			setupMock: func(client *MockHTTPClientInterface) {
				client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("network error"))
			},
			username:      "testuser",
			password:      "testpass123",
			email:         "test@example.com",
			expectedError: "request failed: network error",
		},
		{
			name: "server error",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusBadRequest,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"message": "username already exists"
					}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			username:      "testuser",
			password:      "testpass123",
			email:         "test@example.com",
			expectedError: "server error: username already exists",
		},
		{
			name: "invalid response format",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			username:      "testuser",
			password:      "testpass123",
			email:         "test@example.com",
			expectedError: "failed to decode response: invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := NewMockHTTPClientInterface(ctrl)
			tt.setupMock(httpClientMock)

			service := &AuthService{
				client:  httpClientMock,
				baseURL: "http://test",
			}

			result, err := service.Register(tt.username, tt.password, tt.email)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tempDir := t.TempDir()

	tests := []struct {
		name          string
		setupMock     func(*MockHTTPClientInterface)
		username      string
		password      string
		tokenPath     func() (string, error)
		expected      *models.AuthResponse
		expectedError string
	}{
		{
			name: "successful login",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"access_token": "test-token",
						"refresh_token": "refresh-token",
						"expires_in": "3600"
					}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			username: "testuser",
			password: "testpass123",
			tokenPath: func() (string, error) {
				return filepath.Join(tempDir, tokenFileName), nil
			},
			expected: &models.AuthResponse{
				AccessToken:  "test-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    "3600",
			},
		},
		{
			name: "empty token from server",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"access_token": "",
						"refresh_token": "",
						"expires_in": "0"
					}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			username: "testuser",
			password: "testpass123",
			tokenPath: func() (string, error) {
				return filepath.Join(tempDir, tokenFileName), nil
			},
			expectedError: "server returned empty token",
		},
		{
			name: "failed to save token",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"access_token": "test-token",
						"refresh_token": "refresh-token",
						"expires_in": "3600"
					}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			username: "testuser",
			password: "testpass123",
			tokenPath: func() (string, error) {
				return "", errors.New("save error")
			},
			expectedError: "failed to save token: failed to get token path: save error",
		},
		{
			name: "request failed",
			setupMock: func(client *MockHTTPClientInterface) {
				client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("network error"))
			},
			username: "testuser",
			password: "testpass123",
			tokenPath: func() (string, error) {
				return filepath.Join(tempDir, tokenFileName), nil
			},
			expectedError: "request failed: network error",
		},
		{
			name: "invalid response format",
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			username: "testuser",
			password: "testpass123",
			tokenPath: func() (string, error) {
				return filepath.Join(tempDir, tokenFileName), nil
			},
			expectedError: "failed to decode response: invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := NewMockHTTPClientInterface(ctrl)
			tt.setupMock(httpClientMock)

			service := &AuthService{
				client:    httpClientMock,
				baseURL:   "http://test",
				tokenPath: tt.tokenPath,
			}

			result, err := service.Login(tt.username, tt.password)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)

			if tt.expected != nil && tt.expected.AccessToken != "" {
				tokenPath := filepath.Join(tempDir, tokenFileName)
				data, err := os.ReadFile(tokenPath)
				require.NoError(t, err)
				assert.Equal(t, tt.expected.AccessToken, string(data))
			}
		})
	}
}
func TestAuthService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tempDir := t.TempDir()

	tests := []struct {
		name          string
		setup         func(*AuthService)
		expectedError string
	}{
		{
			name: "successful logout",
			setup: func(s *AuthService) {
				tokenPath := filepath.Join(tempDir, tokenFileName)
				err := os.WriteFile(tokenPath, []byte("test-token"), tokenFileMode)
				require.NoError(t, err)
			},
		},
		{
			name: "failed to clear token",
			setup: func(s *AuthService) {
				s.tokenPath = func() (string, error) {
					return "", errors.New("clear error")
				}
			},
			expectedError: "failed to get token path: clear error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &AuthService{
				client:  NewMockHTTPClientInterface(ctrl),
				baseURL: "http://test",
				tokenPath: func() (string, error) {
					return filepath.Join(tempDir, tokenFileName), nil
				},
			}
			tt.setup(service)

			err := service.Logout()

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)

			tokenPath := filepath.Join(tempDir, tokenFileName)
			_, err = os.Stat(tokenPath)
			assert.True(t, os.IsNotExist(err))
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tempDir := t.TempDir()

	tests := []struct {
		name          string
		setup         func(*MockHTTPClientInterface, *AuthService)
		expectedValid bool
		expectedID    string
		expectedError string
	}{
		{
			name: "valid token",
			setup: func(client *MockHTTPClientInterface, s *AuthService) {
				tokenPath := filepath.Join(tempDir, tokenFileName)
				err := os.WriteFile(tokenPath, []byte("valid-token"), tokenFileMode)
				require.NoError(t, err)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"valid": true,
						"user_id": "user123"
					}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedValid: true,
			expectedID:    "user123",
		},
		{
			name: "invalid token",
			setup: func(client *MockHTTPClientInterface, s *AuthService) {
				tokenPath := filepath.Join(tempDir, tokenFileName)
				err := os.WriteFile(tokenPath, []byte("invalid-token"), tokenFileMode)
				require.NoError(t, err)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"valid": false,
						"user_id": ""
					}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedValid: false,
			expectedID:    "",
		},
		{
			name: "empty token",
			setup: func(client *MockHTTPClientInterface, s *AuthService) {
				tokenPath := filepath.Join(tempDir, tokenFileName)
				err := os.WriteFile(tokenPath, []byte(""), tokenFileMode)
				require.NoError(t, err)
				client.EXPECT().Do(gomock.Any()).Times(0)
			},
			expectedValid: false,
			expectedID:    "",
		},
		{
			name: "invalid validation response",
			setup: func(client *MockHTTPClientInterface, s *AuthService) {
				tokenPath := filepath.Join(tempDir, tokenFileName)
				err := os.WriteFile(tokenPath, []byte("valid-token"), tokenFileMode)
				require.NoError(t, err)

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedError: "failed to decode response: invalid character",
		},
		{
			name: "failed to load token",
			setup: func(client *MockHTTPClientInterface, s *AuthService) {
				s.tokenPath = func() (string, error) {
					return "", errors.New("path error")
				}
				client.EXPECT().Do(gomock.Any()).Times(0)
			},
			expectedError: "failed to load token: failed to get token path: path error",
		},
		{
			name: "validation request failed",
			setup: func(client *MockHTTPClientInterface, s *AuthService) {
				tokenPath := filepath.Join(tempDir, tokenFileName)
				err := os.WriteFile(tokenPath, []byte("valid-token"), tokenFileMode)
				require.NoError(t, err)

				client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("network error"))
			},
			expectedError: "validation request failed: request failed: network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := NewMockHTTPClientInterface(ctrl)
			service := &AuthService{
				client:  httpClientMock,
				baseURL: "http://test",
				tokenPath: func() (string, error) {
					return filepath.Join(tempDir, tokenFileName), nil
				},
			}
			tt.setup(httpClientMock, service)

			valid, userID, err := service.ValidateToken()

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedValid, valid)
			assert.Equal(t, tt.expectedID, userID)
		})
	}
}

func TestAuthService_tokenMethods(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("save and load token", func(t *testing.T) {
		service := &AuthService{
			tokenPath: func() (string, error) {
				return filepath.Join(tempDir, tokenFileName), nil
			},
		}

		token := "test-token-123"
		err := service.saveToken(token)
		require.NoError(t, err)

		path := filepath.Join(tempDir, tokenFileName)
		stat, err := os.Stat(path)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(tokenFileMode), stat.Mode().Perm())

		loadedToken, err := service.loadToken()
		require.NoError(t, err)
		assert.Equal(t, token, loadedToken)
	})

	t.Run("clear token", func(t *testing.T) {
		service := &AuthService{
			tokenPath: func() (string, error) {
				return filepath.Join(tempDir, tokenFileName), nil
			},
		}

		token := "test-token-456"
		err := service.saveToken(token)
		require.NoError(t, err)

		err = service.clearToken()
		require.NoError(t, err)

		path := filepath.Join(tempDir, tokenFileName)
		_, err = os.Stat(path)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("failed to get token path", func(t *testing.T) {
		service := &AuthService{
			tokenPath: func() (string, error) {
				return "", errors.New("token path error")
			},
		}

		_, err := service.tokenPath()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token path error")
	})
}

func TestAuthService_doRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		method        string
		url           string
		body          interface{}
		setupMock     func(*MockHTTPClientInterface)
		expected      []byte
		expectedError string
	}{
		{
			name:   "successful request",
			method: http.MethodPost,
			url:    "http://test/success",
			body:   map[string]string{"key": "value"},
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"result":"success"}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expected: []byte(`{"result":"success"}`),
		},
		{
			name:   "failed to encode body",
			method: http.MethodPost,
			url:    "http://test/fail",
			body:   make(chan int),
			setupMock: func(client *MockHTTPClientInterface) {
			},
			expectedError: "failed to encode request body: json: unsupported type: chan int",
		},
		{
			name:   "failed to create request",
			method: "\x00",
			url:    "http://test/fail",
			body:   map[string]string{"key": "value"},
			setupMock: func(client *MockHTTPClientInterface) {
			},
			expectedError: "failed to create request: net/http: invalid method",
		},
		{
			name:   "request failed",
			method: http.MethodPost,
			url:    "http://test/fail",
			body:   map[string]string{"key": "value"},
			setupMock: func(client *MockHTTPClientInterface) {
				client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("network error"))
			},
			expectedError: "request failed: network error",
		},
		{
			name:   "server error with message",
			method: http.MethodPost,
			url:    "http://test/error",
			body:   map[string]string{"key": "value"},
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewBufferString(`{"message":"invalid request"}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedError: "server error: invalid request",
		},
		{
			name:   "server error without message",
			method: http.MethodPost,
			url:    "http://test/error",
			body:   map[string]string{"key": "value"},
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedError: "unexpected status code: 500",
		},
		{
			name:   "failed to read response",
			method: http.MethodPost,
			url:    "http://test/fail",
			body:   map[string]string{"key": "value"},
			setupMock: func(client *MockHTTPClientInterface) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       &errorReaderBody{err: errors.New("read error")},
				}
				client.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
			expectedError: "failed to read response: read error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := NewMockHTTPClientInterface(ctrl)
			tt.setupMock(httpClientMock)

			service := &AuthService{
				client: httpClientMock,
			}

			result, err := service.doRequest(tt.method, tt.url, tt.body)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

type errorReaderBody struct {
	err error
}

func (r *errorReaderBody) Read(_ []byte) (int, error) {
	return 0, r.err
}

func (r *errorReaderBody) Close() error {
	return nil
}
