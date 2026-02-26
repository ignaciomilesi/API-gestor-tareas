package db

import (
	"api-gestor-tareas/internal/models"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tareasManager struct {
	db *pgxpool.Pool
}

func NewTareasManager(dataBase *pgxpool.Pool) *tareasManager {
	return &tareasManager{
		db: dataBase,
	}
}

// retorna el id del usuario generado
func (tm *tareasManager) RegistrarTarea(ctx context.Context, newTarea models.Tarea) (int, error) {

	query := `INSERT INTO tareas(titulo, fecha_creacion, completada, fecha_completada, id_usuario)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id`
	var id int

	err := tm.db.QueryRow(ctx, query,
		newTarea.Titulo,
		newTarea.Fecha_creacion,
		newTarea.Completada,
		newTarea.Fecha_completada,
		newTarea.Id_usuario).Scan(&id)

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {

			if pgErr.Code == "23503" && pgErr.ConstraintName == "usuario_asignado" {
				return 0, ErrUsuarioAsignadoNoexiste
			}
			return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
		}

	}

	return id, nil
}

func (tm *tareasManager) Listar(ctx context.Context, IdUsuario int, completadas bool) ([]models.Tarea, error) {

	query := `SELECT id, titulo, fecha_creacion, completada, fecha_completada FROM tareas
				WHERE id_usuario = $1
				AND completada = ` + strconv.FormatBool(!completadas)

	var tareas []models.Tarea

	filas, err := tm.db.Query(ctx, query, IdUsuario)

	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}
	defer filas.Close()

	for filas.Next() {
		var t models.Tarea
		err := filas.Scan(&t.Id, &t.Titulo, &t.Fecha_creacion, &t.Completada, &t.Fecha_completada)
		if err != nil {
			return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
		}
		tareas = append(tareas, t)
	}

	return tareas, nil

}

func (tm *tareasManager) ModificarTitulo(ctx context.Context, IdTarea int, nuevoTitulo string) error {

	queryUpdate := `UPDATE tareas
				SET titulo = $1
				WHERE id = $2
				RETURNING id`
	var id int

	err := tm.db.QueryRow(ctx, queryUpdate, nuevoTitulo, IdTarea).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTareaNoExiste
	}

	if err != nil {
		return fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return nil
}

func (tm *tareasManager) MarcarComoCompletada(ctx context.Context, IdTarea int, fechaCompletada *time.Time) error {

	queryUpdate := `UPDATE tareas
				SET completada = true,
				fecha_completada = $1
				WHERE id = $2
				RETURNING id`
	var id int

	err := tm.db.QueryRow(ctx, queryUpdate, fechaCompletada, IdTarea).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTareaNoExiste
	}

	if err != nil {
		return fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return nil
}

func (tm *tareasManager) BuscarEnTitulo(ctx context.Context, palabraABuscar string, IdUsuario int) ([]models.Tarea, error) {

	query := `SELECT id, titulo, fecha_creacion, completada, fecha_completada FROM tareas
	WHERE titulo ILIKE '%' || $1 || '%'
	AND id_usuario = $2`

	var tareas []models.Tarea

	filas, err := tm.db.Query(ctx, query, palabraABuscar, IdUsuario)

	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}
	defer filas.Close()

	for filas.Next() {
		var t models.Tarea
		err := filas.Scan(&t.Id, &t.Titulo, &t.Fecha_creacion, &t.Completada, &t.Fecha_completada)
		if err != nil {
			return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
		}
		tareas = append(tareas, t)
	}

	return tareas, nil

}
