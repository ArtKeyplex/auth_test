package service

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type UserService interface {
	ValidateCredentials(username, password string) (bool, error)
	GenerateToken(username string) (string, error)
	RefreshToken(token string) (string, error)
}

type User struct {
	Username string
	Password string
}

type UserStore interface {
	Get(username string) (*User, error)
}

type userService struct {
	store     UserStore
	jwtSecret []byte
}

func NewUserService(store UserStore, jwtSecret []byte) UserService {
	return &userService{
		store:     store,
		jwtSecret: jwtSecret, // секретный ключ для шифрования
	}
}

// ValidateCredentials проверяет есть ли юзер с таким паролем
func (us *userService) ValidateCredentials(username, password string) (bool, error) {
	user, err := us.store.Get(username)
	if err != nil {
		return false, ErrUserNotFound
	}
	if user.Password != password {
		return false, nil
	}
	return true, nil
}

// GenerateToken генерирует новый токен
func (us *userService) GenerateToken(username string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(60 * time.Minute).Unix(), // 3600 sec = 1 hour
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(us.jwtSecret)
	if err != nil {
		return "", ErrTokenSigning
	}
	return signedToken, nil
}

// RefreshToken проверяет токен и выдаёт новый с обновлённым exp
func (us *userService) RefreshToken(tokenStr string) (string, error) {
	// парсим и валидируем токен
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnsupportedSigningMethod
		}
		return us.jwtSecret, nil
	})
	if err != nil {
		return "", ErrTokenValidation
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}

	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return "", ErrTokenExpired
	}

	// Создаём новый токен с обновленным exp
	signedToken, err := us.GenerateToken(claims.Subject)
	if err != nil {
		return "", ErrTokenGeneration
	}
	return signedToken, nil
}
