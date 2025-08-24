package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"regexp"
)

func GetProjectPath() string {
	projectPath := os.Getenv("PROJECT_PATH")

	if projectPath != "" {
		return projectPath
	}

	currentDir, _ := os.Getwd()

	// Создаем путь к родительской директории
	// parentDir := filepath.Join(currentDir, "../../../")

	return currentDir
}

func DecodeJson(body []byte, result any) error {
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	return nil
}

func HashPassword(password string) (string, error) {
	// Генерация случайной соли (важно для безопасности)
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Конкатенация пароля и соли
	passwordWithSalt := password + string(salt)

	// Хеширование с использованием SHA256
	hash := sha256.Sum256([]byte(passwordWithSalt))

	// Кодирование хеша и соли в base64 для хранения в базе данных
	hashedPassword := base64.StdEncoding.EncodeToString(append(hash[:], salt...))

	return hashedPassword, nil
}

func CheckPassword(hashedPassword, password string) (bool, error) {
	// Декодирование хешированного пароля и соли из base64
	decoded, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to decode hashed password: %w", err)
	}

	// Извлечение хеша и соли
	hash := decoded[:32] // Первые 32 байта - хеш
	salt := decoded[32:] // Остальные байты - соль

	// Конкатенация пароля и соли
	passwordWithSalt := password + string(salt)

	// Хеширование введенного пароля с той же солью
	calculatedHash := sha256.Sum256([]byte(passwordWithSalt))

	// Сравнение хешей
	return fmt.Sprintf("%x", calculatedHash) == fmt.Sprintf("%x", hash), nil
}

func ResponseJson(response interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
	}
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func GenerateToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err) // Обработка ошибки
		}
		token[i] = charset[randomIndex.Int64()]
	}
	return string(token)
}
