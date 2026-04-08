package handlers

import (
	"api-gestor-tareas/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type tareasServiceInterface interface {
	// Parámetros:
	// 		- descripción, fechaCreación, idUsuario
	// Salida:
	// 		- id de la tarea generada
	// Errores que puede devuelve:
	// 		-  ErrDescripcionRequerida
	// 		-  ErrIdNoValido
	// 		-  ErrFechaNoValida
	// 		-  ErrUsuarioAsignadoNoexiste
	CrearTarea(context.Context, string, time.Time, int) (int, error)

	// Parámetros:
	// 		- idUsuario
	// 		- incluir completadas
	// Salida:
	// 		- lista de tareas
	// Errores que puede devuelve:
	// 		-  ErrIdNoValido
	ListarTareas(context.Context, int, bool) ([]domain.Tarea, error)

	// Parámetros:
	// 		- idTarea
	// 		- nueva descripción
	// Salida:
	// 		-
	// Errores que puede devuelve:
	//		-  ErrIdNoValido
	//		-  ErrTareaNoExiste
	ModificarDescripcion(context.Context, int, string) error

	// Parámetros:
	// 		- idTarea
	// 		- fecha
	// Salida:
	// 		-
	// Errores que puede devuelve:
	//		-  ErrIdNoValido
	//		-  ErrFechaNoValida
	//		-  ErrTareaNoExiste
	MarcarComoCompletada(context.Context, int, *time.Time) error

	// Parámetros:
	// 		- Parámetro de búsqueda
	// 		- fecha
	// Salida:
	// 		- Listado de tareas encontradas
	// Errores que puede devuelve:
	//		-  ErrIdNoValido
	//		-  ErrParametroDeBusquedaNoValido
	Buscar(context.Context, string, int) ([]domain.Tarea, error)
}

type tareasHandler struct {
	tareaService tareasServiceInterface
}

func NewTareaHandler(nuevoTareaService tareasServiceInterface) *tareasHandler {
	return &tareasHandler{
		tareaService: nuevoTareaService,
	}
}

func (th *tareasHandler) Nueva(c *gin.Context) {

	var req struct {
		Descrip   string `json:"descripcion" binding:"required"`
		Fecha     string `json:"fecha" binding:"required"`
		IdUsuario int    `json:"id_usuario" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	fechaParseada, err := time.Parse("02/01/2006", req.Fecha)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "fecha inválida, formato esperado DD/MM/YYYY",
			"detalle": err,
		})
		return
	}

	id, err := th.tareaService.CrearTarea(c.Request.Context(), req.Descrip, fechaParseada, req.IdUsuario)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDescripcionRequerida),
			errors.Is(err, domain.ErrIdNoValido),
			errors.Is(err, domain.ErrFechaNoValida):
			c.JSON(400, gin.H{
				"error":   "Campos requeridos no válidos",
				"detalle": err,
			})

		case errors.Is(err, domain.ErrUsuarioAsignadoNoexiste):
			c.JSON(404, gin.H{
				"error":   "Usuario no valido",
				"detalle": err,
			})

		default:
			c.JSON(500, gin.H{
				"error":   "error interno",
				"detalle": err,
			})
			fmt.Println(err.Error())
		}
		return
	}

	c.JSON(200, gin.H{"id": id})
}

func (th *tareasHandler) Listar(c *gin.Context) {

	var req struct {
		IdUsuario   int  `json:"id_usuario" binding:"required"`
		Completadas bool `json:"completadas" `
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	lista, err := th.tareaService.ListarTareas(c.Request.Context(), req.IdUsuario, req.Completadas)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrIdNoValido):
			c.JSON(400, gin.H{
				"error":   "Campos requeridos no válidos",
				"detalle": err,
			})

		default:
			c.JSON(500, gin.H{
				"error":   "error interno",
				"detalle": err,
			})
			fmt.Println(err.Error())
		}
		return
	}

	c.JSON(200, gin.H{"tareas:": lista})
}

func (th *tareasHandler) Modificar(c *gin.Context) {

	var req struct {
		IdTarea     int    `json:"id_tarea" binding:"required"`
		Descripcion string `json:"nueva_descripcion" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	err := th.tareaService.ModificarDescripcion(c.Request.Context(), req.IdTarea, req.Descripcion)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrIdNoValido):
			c.JSON(400, gin.H{
				"error":   "Campos requeridos no válidos",
				"detalle": err,
			})
		case errors.Is(err, domain.ErrTareaNoExiste):
			c.JSON(404, gin.H{
				"error":   "Id solicitado no existe",
				"detalle": err,
			})

		default:
			c.JSON(500, gin.H{
				"error":   "error interno",
				"detalle": err,
			})
			fmt.Println(err.Error())
		}
		return
	}

	c.JSON(200, gin.H{"Resp": "Cambiado con éxito"})
}

func (th *tareasHandler) Finalizar(c *gin.Context) {

	var req struct {
		IdTarea int    `json:"id_tarea" binding:"required"`
		Fecha   string `json:"fecha" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	fechaParseada, err := time.Parse("02/01/2006", req.Fecha)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "fecha inválida, formato esperado DD/MM/YYYY",
			"detalle": err,
		})
		return
	}

	err = th.tareaService.MarcarComoCompletada(c.Request.Context(), req.IdTarea, &fechaParseada)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrIdNoValido),
			errors.Is(err, domain.ErrFechaNoValida):
			c.JSON(400, gin.H{
				"error":   "Campos requeridos no válidos",
				"detalle": err,
			})
		case errors.Is(err, domain.ErrTareaNoExiste):
			c.JSON(404, gin.H{
				"error":   "Id solicitado no existe",
				"detalle": err,
			})

		default:
			c.JSON(500, gin.H{
				"error":   "error interno",
				"detalle": err,
			})
			fmt.Println(err.Error())
		}
		return
	}

	c.JSON(200, gin.H{"Resp": "Tarea finalizada"})
}

func (th *tareasHandler) Buscar(c *gin.Context) {

	var req struct {
		IdUsuario         int    `json:"id_usuario" binding:"required"`
		ParametroBusqueda string `json:"parametro_busqueda" binding:"required" `
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos (json)",
			"detalle": err,
		})
		return
	}

	lista, err := th.tareaService.Buscar(c.Request.Context(), req.ParametroBusqueda, req.IdUsuario)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrIdNoValido),
			errors.Is(err, domain.ErrParametroDeBusquedaNoValido):
			c.JSON(400, gin.H{
				"error":   "Campos requeridos no válidos",
				"detalle": err,
			})

		default:
			c.JSON(500, gin.H{
				"error":   "error interno",
				"detalle": err,
			})
			fmt.Println(err.Error())
		}
		return
	}

	c.JSON(200, gin.H{"tareas:": lista})
}
