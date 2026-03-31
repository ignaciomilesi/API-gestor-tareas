package handlers

import (
	"api-gestor-tareas/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
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
}

func NewUserHandler(nuevoUserService userServiceInterface) *userHandler {
	return &userHandler{
		userService: nuevoUserService,
	}
}

func (h *userHandler) Singin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"` //  no vacío, formato de email válido
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Campos requeridos no válidos"})
		return
	}

	id, err := h.userService.CrearUsuario(c.Request.Context(), req.Email, req.Password)
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

	c.JSON(200, gin.H{"id": id})
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

	c.JSON(200, gin.H{"id": id})
}

func (h *userHandler) ActualizarContraseña(c *gin.Context) {
	var req struct {
		Id       int    `json:"id" binding:"required"` //  no vacío
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Campos requeridos no válidos"})
		return
	}

	err := h.userService.ModificarContraseña(c.Request.Context(), req.Id, req.Password)
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
