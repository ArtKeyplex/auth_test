package verify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockUserService struct {
	ValidateCredentialsFunc func(username, password string) (bool, error)
	GenerateTokenFunc       func(username string) (string, error)
	RefreshTokenFunc        func(username string) (string, error)
}

func (m *MockUserService) ValidateCredentials(username, password string) (bool, error) {
	return m.ValidateCredentialsFunc(username, password)
}

func (m *MockUserService) GenerateToken(username string) (string, error) {
	return m.GenerateTokenFunc(username)
}

func (m *MockUserService) RefreshToken(username string) (string, error) {
	return m.RefreshTokenFunc(username)
}

func TestVerify(t *testing.T) {
	mockUserService := &MockUserService{
		RefreshTokenFunc: func(username string) (string, error) {
			return "newToken", nil
		},
		GenerateTokenFunc: func(username string) (string, error) {
			return "generatedToken", nil
		},
	}

	handler := &VerifyHandler{
		UserService: mockUserService,
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Валидные данные",
			authHeader:     "Bearer " + "newToken",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Пустой/неверный заголовок",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Add("Authorization", tt.authHeader)
			}
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.Code)
			}
		})
	}
}
