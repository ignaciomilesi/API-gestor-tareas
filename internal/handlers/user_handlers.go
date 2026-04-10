package handlers

import (
	"api-gestor-tareas/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type userServiceInterface interface {
	// Parámetros:
	// 		- mail
	//		- password
	// Salida:
	// 		- Id del usuario generado
	// Errores que puede devuelve:
	// 		-  ErrEmailRequerido
	// 		-  ErrPasswordRequerido
	// 		-  ErrPasswordCorto
	// 		-  ErrEmailDuplicado
	CrearUsuario(context.Context, string, string) (int, error)

	// Parámetros:
	// 		- Id del usuario
	//		- nuevo password
	// Errores que puede devuelve:
	// 		- ErrIdNoEncontrado
	// 		- ErrIdNoValido
	// 		- ErrPasswordRequerido
	// 		- ErrPasswordCorto
	ModificarContraseña(context.Context, int, string) error

	// Parámetros:
	// 		- mail
	//		- password
	// Salida:
	// 		- Id del usuario
	// Errores que puede devuelve:
	// 		- ErrEmailRequerido
	//		- ErrPasswordRequerido
	//		- ErrEmailNoEncontrado
	//		- ErrPasswordIncorrecto
	ObtenerId(context.Context, string, string) (int, error)
}

type userHandler struct {
	userService userServiceInterface
	jwtSecret   string
}

func NewUserHandler(nuevoUserService userServiceInterface, secret string) *userHandler {
	return &userHandler{
		userService: nuevoUserService,
		jwtSecret:   secret,
	}
}

func (h *userHandler) Signin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"` //  no vacío, formato de email válido
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Campos requeridos no válidos"})
		return
	}

	_, err := h.userService.CrearUsuario(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailRequerido),
			errors.Is(err, domain.ErrPasswordRequerido),
			errors.Is(err, domain.ErrPasswordCorto):
			c.JSON(400, gin.H{"error": "Campos requeridos no válidos"})

		case errors.Is(err, domain.ErrEmailDuplicado):
			c.JSON(409, gin.H{"error": "Mail existente"})

		default:
			c.JSON(500, gin.H{"error": "error interno"})
			fmt.Println(err.Error())
		}
		return
	}

	c.JSON(200, gin.H{"resp": "usuario generado"})
}

func (h *userHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"` //  no vacío, formato de email válido
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Campos requeridos no válidos"})
		return
	}

	id, err := h.userService.ObtenerId(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailRequerido),
			errors.Is(err, domain.ErrPasswordRequerido):
			c.JSON(400, gin.H{"error": "Campos requeridos no válidos"})

		case errors.Is(err, domain.ErrEmailNoEncontrado),
			errors.Is(err, domain.ErrPasswordIncorrecto):
			c.JSON(401, gin.H{"error": "Usuario no encontrado"})

		default:
			c.JSON(500, gin.H{"error": "error interno"})
			fmt.Println(err.Error())
		}
		return
	}

	// generar JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   req.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(500, gin.H{"error": "error generando token"})
		return
	}

	// Devuelvo token
	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

func (h *userHandler) ActualizarContraseña(c *gin.Context) {
	var req struct {
		//Id       int    `json:"id" binding:"required"` //  no vacío
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Campos requeridos no válidos"})
		return
	}

	// obtengo el id del contexto
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "usuario no autenticado"})
		return
	}

	// lo transformo en int
	userID, ok := userIDValue.(int)
	if !ok {
		c.JSON(500, gin.H{"error": "error interno"})
		return
	}

	err := h.userService.ModificarContraseña(c.Request.Context(), userID, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrIdNoValido),
			errors.Is(err, domain.ErrPasswordRequerido),
			errors.Is(err, domain.ErrPasswordCorto):
			fmt.Println(err)
			c.JSON(400, gin.H{"error": "Campos requeridos no válidos"})

		case errors.Is(err, domain.ErrIdNoEncontrado):
			c.JSON(401, gin.H{"error": "Usuario no encontrado"})

		default:
			c.JSON(500, gin.H{"error": "error interno"})
			fmt.Println(err.Error())
		}
		return
	}

	c.JSON(200, gin.H{"Resp": "Cambiado con éxito"})
}
