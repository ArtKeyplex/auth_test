package login

import (
	"encoding/base64"
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

func TestLogin(t *testing.T) {
	mockUserService := &MockUserService{
		ValidateCredentialsFunc: func(username, password string) (bool, error) {
			if username == "Никита" && password == "password123" {
				return true, nil
			}
			return false, nil
		},
		GenerateTokenFunc: func(username string) (string, error) {
			return "generatedToken", nil
		},
	}

	handler := &LoginHandler{
		UserService: mockUserService,
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Валидные данные",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("Никита:password123")),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Ошибка парсинга",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("Никитаpassword123")),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Нет заголовка",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Неверный заголовок",
			authHeader:     "Basic invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Неверный пароль",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("Никита:password1234")),
			expectedStatus: http.StatusBadRequest,
		},

		{
			name:           "Неверный парсинг 2",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("Никита:")),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Add("Authorization", tt.authHeader)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.Code)
			}
		})
	}
}
