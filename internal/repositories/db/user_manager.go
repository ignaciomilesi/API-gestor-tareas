package db

import (
	"api-gestor-tareas/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userManager struct {
	db *pgxpool.Pool
}

func NewUserManager(dataBase *pgxpool.Pool) *userManager {
	return &userManager{
		db: dataBase,
	}
}

func (um *userManager) GenerarNuevoUsuario(ctx context.Context, newUsuario models.Usuario) error {
	query := `INSERT INTO usuarios(email, password_hash)
				VALUES ($1, $2)`

	if _, err := um.db.Exec(ctx, query, newUsuario.Email, newUsuario.Password_hash); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {

			if pgErr.Code == "23505" && pgErr.ConstraintName == "email_unico" {
				return fmt.Errorf("No se crea el usuario, el mail ya existe: %v", err)
			}
			return fmt.Errorf("No se crea el usuario, error inesperado: %v", err)
		}
	}

	return nil
}

func (um *userManager) ComprobarExisteUsuario(ctx context.Context, usuario models.Usuario) bool {
	query := `SELECT password_hash FROM usuarios
				WHERE email = $1`

	var passwordHashEncontrado string

	err := um.db.QueryRow(ctx, query, usuario.Email).Scan(&passwordHashEncontrado)

	if err != nil {
		return false // si devuelve error, es que no existe el usuario en el db
	}

	return passwordHashEncontrado == usuario.Password_hash

}

func (um *userManager) ModifcarContraseña(ctx context.Context, usuario models.Usuario, nuevoPasswordHash string) error {

	queryUpdate := `UPDATE usuarios
				SET password_hash = $1
				WHERE email = $2`

	if _, err := um.db.Exec(ctx, queryUpdate, nuevoPasswordHash, usuario.Email); err != nil {
		return fmt.Errorf("Error al actualizar el password_hash: %v", err)
	}

	return nil
}
