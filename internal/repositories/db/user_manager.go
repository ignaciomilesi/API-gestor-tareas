package db

import (
	"context"
	"errors"
	"fmt"

	"api-gestor-tareas/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type userManager struct {
	db *pgxpool.Pool
}

func NewUserManager(dataBase *pgxpool.Pool) *userManager {
	return &userManager{
		db: dataBase,
	}
}

// retorna el id del usuario generado
func (um *userManager) GenerarNuevoUsuario(ctx context.Context, newUsuario domain.Usuario) (int, error) {
	query := `INSERT INTO usuarios(email, password_hash)
				VALUES ($1, $2)
				RETURNING id;`

	var id int

	err := um.db.QueryRow(ctx, query, newUsuario.Email, newUsuario.Password_hash).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {

			if pgErr.Code == "23505" && pgErr.ConstraintName == "email_unico" {
				return 0, domain.ErrEmailDuplicado
			}
			return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
		}
	}

	return id, nil
}

// verifica la existencia del usuario, devuelve el id del mismo
func (um *userManager) ObternerId(ctx context.Context, usuario domain.Usuario) (int, error) {
	query := `SELECT id, password_hash FROM usuarios
				WHERE email = $1`

	var passwordHashEncontrado string
	var id int

	err := um.db.QueryRow(ctx, query, usuario.Email).Scan(&id, &passwordHashEncontrado)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrEmailNoEncontrado
		}
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHashEncontrado), []byte(usuario.Password_hash))
	if err != nil {
		return 0, domain.ErrPasswordIncorrecto
	}

	return id, nil

}

func (um *userManager) ModifcarContraseña(ctx context.Context, IdUsuario int, nuevoPasswordHash string) error {

	queryUpdate := `UPDATE usuarios
				SET password_hash = $1
				WHERE id = $2
				RETURNING id`
	var id int

	err := um.db.QueryRow(ctx, queryUpdate, nuevoPasswordHash, IdUsuario).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) { // sql.ErrNoRows) {
		return domain.ErrIdNoEncontrado
	}

	if err != nil {
		return fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return nil
}
