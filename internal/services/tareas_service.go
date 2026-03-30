package services

import (
	"api-gestor-tareas/config"
	"api-gestor-tareas/internal/domain"
	"strings"

	"context"
	"time"
)

type tareaManagerDbInterface interface {
	// Parámetros:
	// 		- Palabra a buscar
	// 		- Id del usuario
	// Salida:
	// 		- Array con las tareas encontradas
	// Errores que puede devuelve:
	// 		-
	BuscarEnTitulo(context.Context, string, int) ([]domain.Tarea, error)
	// Parámetros:
	// 		- Id del usuario
	//		- Si se desea incluir las tareas completadas
	// Salida:
	// 		- Array con las tareas encontradas
	// Errores que puede devuelve:
	// 		-
	Listar(context.Context, int, bool) ([]domain.Tarea, error)
	// Parámetros:
	// 		- Id de la tarea
	//		- Fecha de finalización
	// Errores que puede devuelve:
	// 		- ErrTareaNoExiste
	MarcarComoCompletada(context.Context, int, *time.Time) error
	// Parámetros:
	// 		- Id de la tarea
	//		- Nuevo titulo
	// Errores que puede devuelve:
	// 		- ErrTareaNoExiste
	ModificarTitulo(context.Context, int, string) error
	// Parámetros:
	// 		- Tarea a generar
	// Salida:
	// 		- Id de la tarea registrada
	// Errores que puede devuelve:
	// 		- ErrUsuarioAsignadoNoexiste
	RegistrarTarea(context.Context, domain.Tarea) (int, error)
}

type tareaService struct {
	tareaManagerDb tareaManagerDbInterface
}

func NewTareaService(tareaManager tareaManagerDbInterface) tareaService {
	return tareaService{
		tareaManagerDb: tareaManager,
	}
}

// devuelve el id de la tarea generada
func (ts *tareaService) CrearTarea(ctx context.Context, descripcion string, fechaCreacion time.Time, idUsuario int) (int, error) {

	if strings.TrimSpace(descripcion) == "" {
		return 0, domain.ErrDescripcionRequerida
	}

	if idUsuario < 1 {
		return 0, domain.ErrIdNoValido
	}

	if fechaCreacion.After(time.Now()) {
		return 0, domain.ErrFechaNoValida
	}

	nuevaTarea := domain.Tarea{
		Descripcion:    strings.TrimSpace(descripcion),
		Fecha_creacion: fechaCreacion,
		Completada:     false,
		Id_usuario:     idUsuario,
	}

	return ts.tareaManagerDb.RegistrarTarea(ctx, nuevaTarea)

}

func (ts *tareaService) ListarTareas(ctx context.Context, idUsuario int, incluirCompletadas bool) ([]domain.Tarea, error) {

	if idUsuario < 1 {
		return nil, domain.ErrIdNoValido
	}

	return ts.tareaManagerDb.Listar(ctx, idUsuario, incluirCompletadas)
}

func (ts *tareaService) ModificarDescripcion(ctx context.Context, idTarea int, nuevaDescripcion string) error {

	if strings.TrimSpace(nuevaDescripcion) == "" {
		return domain.ErrDescripcionRequerida
	}

	if idTarea < 1 {
		return domain.ErrIdNoValido
	}

	return ts.tareaManagerDb.ModificarTitulo(ctx, idTarea, strings.TrimSpace(nuevaDescripcion))
}

func (ts *tareaService) MarcarComoCompletada(ctx context.Context, idTarea int, fechaCompletada *time.Time) error {

	if idTarea < 1 {
		return domain.ErrIdNoValido
	}

	if fechaCompletada.After(time.Now()) {
		return domain.ErrFechaNoValida
	}

	return ts.tareaManagerDb.MarcarComoCompletada(ctx, idTarea, fechaCompletada)

}

func (ts *tareaService) Buscar(ctx context.Context, parametroABuscar string, idUsuario int) ([]domain.Tarea, error) {

	if idUsuario < 1 {
		return nil, domain.ErrIdNoValido
	}

	if len(strings.TrimSpace(parametroABuscar)) < config.LargoMinimoParametroBusqueda {
		return nil, domain.ErrParametroDeBusquedaNoValido
	}

	return ts.tareaManagerDb.BuscarEnTitulo(ctx, parametroABuscar, idUsuario)
}
