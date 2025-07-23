package service

import (
	"errors"
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
		return false, err
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
		return "", err
	}
	return signedToken, nil
}

// RefreshToken проверяет токен и выдаёт новый с обновлённым exp
func (us *userService) RefreshToken(tokenStr string) (string, error) {
	// парсим и валидируем токен
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неподдерживаемый метод подписи")
		}
		return us.jwtSecret, nil
	})
	if err != nil {
		return "", errors.New("ошибка парсинга/валидации")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return "", errors.New("недопустимый токен")
	}

	// Проверяем, что токен ещё не истёк
	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return "", errors.New("токен просрочен")
	}

	// Создаём новый токен с обновленным exp
	newClaims := jwt.StandardClaims{
		Subject:   claims.Subject,
		ExpiresAt: time.Now().Add(60 * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signedToken, err := newToken.SignedString(us.jwtSecret)
	if err != nil {
		return "", errors.New("ошибка кодирования в jwt токен")
	}
	return signedToken, nil
}
