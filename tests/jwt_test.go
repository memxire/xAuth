package tests

import (
	"testing"
	"time"

	"xauth/internal/domain/models"
	jwtlib "xauth/internal/lib/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewToken_HappyPath проверяет корректность создания токена при валидных входных данных.
func TestNewToken_HappyPath(t *testing.T) {
	// Подготавливаем тестовые данные
	user := models.User{
		ID:    12345,
		Email: "user@example.com",
	}
	app := models.App{
		ID:     1,
		Secret: "test-secret",
	}
	duration := time.Hour

	// Вызываем функцию создания токена
	tokenString, err := jwtlib.NewToken(user, app, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Парсим токен с использованием секретного ключа
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(app.Secret), nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok, "claims должны быть типа jwt.MapClaims")

	// Проверяем содержимое claims
	assert.Equal(t, float64(user.ID), claims["uid"].(float64))
	assert.Equal(t, user.Email, claims["email"].(string))
	assert.Equal(t, float64(app.ID), claims["app_id"].(float64))

	// Проверяем время истечения token'а (exp)
	exp := int64(claims["exp"].(float64))
	expectedExp := time.Now().Add(duration).Unix()
	// Учитываем небольшую задержку выполнения (delta в 2 секунды)
	assert.InDelta(t, expectedExp, exp, 2, "Время истечения должно быть близко к now+duration")
}

// TestNewToken_NegativeDuration проверяет, что при отрицательной длительности exp оказывается в прошлом.
func TestNewToken_NegativeDuration(t *testing.T) {
	user := models.User{
		ID:    12345,
		Email: "user@example.com",
	}
	app := models.App{
		ID:     1,
		Secret: "test-secret",
	}
	duration := -time.Minute // Отрицательная длительность

	tokenString, err := jwtlib.NewToken(user, app, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Отключаем валидацию claims, чтобы не получать ошибку "token is expired"
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Secret), nil
	}, jwt.WithoutClaimsValidation())
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)

	exp := int64(claims["exp"].(float64))
	assert.True(t, exp < time.Now().Unix(), "Время истечения должно быть в прошлом")
}

// TestNewToken_EmptySecret проверяет поведение функции, когда секрет пустой.
func TestNewToken_EmptySecret(t *testing.T) {
	user := models.User{
		ID:    12345,
		Email: "user@example.com",
	}
	app := models.App{
		ID:     1,
		Secret: "", // пустой секрет
	}
	duration := time.Hour

	tokenString, err := jwtlib.NewToken(user, app, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Secret), nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)

	// Проверяем, что даже при пустом секрете в claims сохранены нужные данные
	assert.Equal(t, float64(user.ID), claims["uid"].(float64))
	assert.Equal(t, user.Email, claims["email"].(string))
	assert.Equal(t, float64(app.ID), claims["app_id"].(float64))
}
