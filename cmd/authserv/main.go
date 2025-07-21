package main

import (
	"auth_test/configs"
	"auth_test/internal/service"
	"auth_test/internal/store"

	"fmt"
)

func main() {
	jwtSecret := []byte("СЕКРЕТНЫЙ_КЛЮЧ") // ваш секретный ключ
	fmt.Println("------------------------------------------------------")
	fmt.Println("----------------1. Скелет + конфиг--------------------")
	fmt.Println("------------------------------------------------------")
	config, err := configs.LoadConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Имя приложения:", config.AppName)
	fmt.Println("------------------------------------------------------")
	fmt.Println("---------------3. InMemoryUserStore-------------------")
	fmt.Println("------------------------------------------------------")
	newStore := &store.InMemoryUserStore{
		Users: map[string]*service.User{
			"Никита": &service.User{Username: "Никита", Password: "password123"},
			"Антон":  &service.User{Username: "Антон", Password: "super123pro"},
			"Адольф": &service.User{Username: "Адольф", Password: "always171wet"},
		},
	}

	// Добавление нового пользователя + проверка на существование (т.к. map)
	addName := "Никита"
	_, err = newStore.Add(&service.User{Username: addName, Password: "password1223"})
	if err != nil {
		fmt.Println("Ошибка добавления:", err, "(", addName, ")")
	}

	// Вывод всех пользователей
	for _, user := range newStore.Users {
		fmt.Println(user)
	}

	fmt.Println("------------------------------------------------------")
	fmt.Println("-----------------2. UserService-----------------------")
	fmt.Println("--------------Validate Credentials--------------------")
	fmt.Println("------------------------------------------------------")
	// Создаем сервис
	userSvc := service.NewUserService(newStore, jwtSecret)

	// Проверим пароль 'password123' среди всех пользователей
	password := "password123"

	validateCount := 1
	for username, _ := range newStore.Users {
		fmt.Print(validateCount, ") ")
		validateCount++
		valid, err := userSvc.ValidateCredentials(username, password)
		if err != nil {
			fmt.Println("Ошибка при валидации:", err)
			continue
		}
		if !valid {
			fmt.Println("Неверные учетные данные")
			continue
		}
		fmt.Println("Учетные данные подтверждены!")
	}
	fmt.Println("------------------------------------------------------")
	fmt.Println("-------------------Generate Token---------------------")
	fmt.Println("------------------------------------------------------")

	token, err := userSvc.GenerateToken(addName)
	if err != nil {
		fmt.Println("Ошибка при генерации:", err)
	}
	fmt.Println("Сгенерированный токен:", token)

	fmt.Println("------------------------------------------------------")
	fmt.Println("-------------------Refresh Token----------------------")
	fmt.Println("------------------------------------------------------")
	newToken, err := userSvc.RefreshToken(token)
	if err != nil {
		fmt.Println("Ошибка обновления:", err)
	}
	fmt.Println("Обновленный токен:", newToken)
	fmt.Println("------------------------------------------------------")
}
