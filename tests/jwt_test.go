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

// TestNewToken_HappyPath checks whether the token was created correctly
// given valid input data.
func TestNewToken_HappyPath(t *testing.T) {
	// Preparing test data
	user := models.User{
		ID:       1,
		Email:    "user@example.com",
		Username: "user",
		PassHash: []byte("password"),
	}
	app := models.App{
		ID:     12345,
		Secret: "test-secret",
	}
	duration := time.Hour

	// Call the token creation method
	tokenString, err := jwtlib.NewToken(user, app, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Parsing a token using a secret key
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any,
		error) {
		return []byte(app.Secret), nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok, "claims must be of type jwt.MapClaims")

	// Check the claims content
	assert.Equal(t, float64(user.ID), claims["uid"].(float64))
	assert.Equal(t, user.Email, claims["email"].(string))
	assert.Equal(t, float64(app.ID), claims["app_id"].(float64))
	assert.Equal(t, user.Username, claims["username"].(string))

	// Checking the token expiration time (exp)
	exp := int64(claims["exp"].(float64))
	expectedExp := time.Now().Add(duration).Unix()
	// Take into account a small delay in execution (delta of 2 seconds)
	assert.InDelta(t, expectedExp, exp, 2,
		"The expiration time should be close to now+duration")
}

// TestNewToken_NegativeDuration checks that if duration is negative,
// exp ends up in the past.
func TestNewToken_NegativeDuration(t *testing.T) {
	user := models.User{
		ID:    12345,
		Email: "user@example.com",
	}
	app := models.App{
		ID:     1,
		Secret: "test-secret",
	}
	duration := -time.Minute // Negative duration

	tokenString, err := jwtlib.NewToken(user, app, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Disable claims validation to avoid the "token is expired" error
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{},
		error) {
		return []byte(app.Secret), nil
	}, jwt.WithoutClaimsValidation())
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)

	exp := int64(claims["exp"].(float64))
	assert.True(t, exp < time.Now().Unix(),
		"The expiration time must be in the past")
}

// TestNewToken_EmptySecret tests the function's behavior when the secret is empty.
func TestNewToken_EmptySecret(t *testing.T) {
	user := models.User{
		ID:    12345,
		Email: "user@example.com",
	}
	app := models.App{
		ID:     1,
		Secret: "", // empty secret
	}
	duration := time.Hour

	tokenString, err := jwtlib.NewToken(user, app, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{},
		error) {
		return []byte(app.Secret), nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)

	// Check that even with an empty secret, the required data is saved in claims
	assert.Equal(t, float64(user.ID), claims["uid"].(float64))
	assert.Equal(t, user.Email, claims["email"].(string))
	assert.Equal(t, float64(app.ID), claims["app_id"].(float64))
}
