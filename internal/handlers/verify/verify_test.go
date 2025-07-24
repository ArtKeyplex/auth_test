package verify

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockUserService struct {
	ValidateCredentialsFunc func(username, password string) (bool, error)
	GenerateTokenFunc       func(username string) (string, error)
	RefreshTokenFunc        func(token string) (string, error)
}

func (m *MockUserService) ValidateCredentials(username, password string) (bool, error) {
	return m.ValidateCredentialsFunc(username, password)
}

func (m *MockUserService) GenerateToken(username string) (string, error) {
	return m.GenerateTokenFunc(username)
}

func (m *MockUserService) RefreshToken(token string) (string, error) {
	return m.RefreshTokenFunc(token)
}

func TestVerify(t *testing.T) {
	mockUserService := &MockUserService{
		RefreshTokenFunc: func(token string) (string, error) {
			if token == "Internal" {
				return "newToken", errors.New("Internal server error")
			}
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
		{
			name:           "Внутренняя ошибка",
			authHeader:     "Bearer " + "Internal",
			expectedStatus: http.StatusInternalServerError,
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
