package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test_secret"

func generateToken(exp int64, t *testing.T) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{

		"user_id": 123,
		"email":   "mail@prueba.com",
		"exp":     exp,
	})

	tokenString, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("error generando token: %v", err)
	}

	return tokenString
}

func setupRouter() *gin.Engine {

	// armo el ruter
	r := gin.New()

	//indico que use el middelware
	r.Use(AuthMiddleware(testSecret))

	r.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(500, gin.H{"error": "no user_id"})
			return
		}

		c.JSON(200, gin.H{"user_id": userID})
	})

	return r
}

func TestAuthMiddleware(t *testing.T) {

	tests := []struct {
		name           string
		key            string
		value          string
		codigoEsperado int
	}{
		{
			name:           "OK",
			key:            "Authorization",
			value:          "Bearer " + generateToken(time.Now().Add(time.Hour).Unix(), t),
			codigoEsperado: http.StatusOK,
		},
		{
			name:           "Sin token",
			key:            "",
			value:          "",
			codigoEsperado: http.StatusUnauthorized,
		},
		{
			name:           "Token Invalido",
			key:            "Authorization",
			value:          "Bearer tokenNoValido",
			codigoEsperado: http.StatusUnauthorized,
		},
		{
			name:           "Token expirado",
			key:            "Authorization",
			value:          "Bearer " + generateToken(time.Now().Add(-time.Hour).Unix(), t),
			codigoEsperado: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {

		gin.SetMode(gin.TestMode)

		t.Run(test.name, func(t *testing.T) {

			router := setupRouter()

			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set(test.key, test.value)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.codigoEsperado {
				t.Errorf("Error no esperado.\nSe esperaba: \n --- %d \nse obtuvo: \n --- %d",
					test.codigoEsperado, w.Code)
			} else {
				fmt.Println(w.Body.String())
			}
		})
	}
}
