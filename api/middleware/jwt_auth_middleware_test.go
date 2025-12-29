// Package middleware
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/horaoen/go-backend-clean-architecture/domain"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret-key"

func createTestToken(claims *domain.JwtCustomClaims, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func setupRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(JwtAuthMiddleware(secret))
	r.GET("/protected", func(c *gin.Context) {
		userID := c.GetString("x-user-id")
		c.JSON(http.StatusOK, gin.H{"userId": userID})
	})
	return r
}

func TestJwtAuthMiddleware_ValidToken(t *testing.T) {
	router := setupRouter(testSecret)

	claims := &domain.JwtCustomClaims{
		Name: "Test User",
		ID:   "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := createTestToken(claims, testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "123")
}

func TestJwtAuthMiddleware_ExpiredToken(t *testing.T) {
	router := setupRouter(testSecret)

	claims := &domain.JwtCustomClaims{
		Name: "Test User",
		ID:   "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	token := createTestToken(claims, testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}

func TestJwtAuthMiddleware_MissingBearerPrefix(t *testing.T) {
	router := setupRouter(testSecret)

	claims := &domain.JwtCustomClaims{
		Name: "Test User",
		ID:   "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := createTestToken(claims, testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing or invalid authorization header")
}

func TestJwtAuthMiddleware_InvalidSignature(t *testing.T) {
	router := setupRouter(testSecret)

	claims := &domain.JwtCustomClaims{
		Name: "Test User",
		ID:   "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := createTestToken(claims, "wrong-secret")

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}

func TestJwtAuthMiddleware_MissingHeader(t *testing.T) {
	router := setupRouter(testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing or invalid authorization header")
}
