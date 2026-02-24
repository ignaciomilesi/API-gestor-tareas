package db

import (
	"api-gestor-tareas/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// - registrar tarea
// - modificar tarea
// - marcar como completada
// - mostrar lista no completadas y completadas

type tareasManager struct {
	db *pgxpool.Pool
}

func NewTareasManager(dataBase *pgxpool.Pool) *tareasManager {
	return &tareasManager{
		db: dataBase,
	}
}

func (tm *tareasManager) RegistrarTarea(ctx context.Context, newTarea models.Tarea) error {

	if newTarea.Titulo == "" {
		return fmt.Errorf("Error, el campo titulo no puede estar en blanco")
	}

	query := `INSERT INTO usuarios(titulo, fecha_creacion, completada, fecha_completada, id_usuario)
				VALUES ($1, $2, $3, $4, $5)`

	_, err := tm.db.Exec(ctx, query,
		newTarea.Titulo,
		newTarea.Fecha_creacion,
		newTarea.Completada,
		newTarea.Fecha_completada,
		newTarea.Id_usuario)
	if err != nil {
		return fmt.Errorf("Error al registrar la tarea: %v", err)
	}

	return nil
}
