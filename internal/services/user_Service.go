package services

import (
	"api-gestor-tareas/config"
	"api-gestor-tareas/internal/domain"
	"fmt"

	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userManagerDbInterface interface {
	// Parámetros:
	// 		- Usuario a generar
	// Salida:
	// 		- Id del usuario generado
	// Errores que puede devuelve:
	// 		- ErrEmailDuplicado
	GenerarNuevoUsuario(context.Context, domain.Usuario) (int, error)

	// Parámetros:
	// 		- id del usuario a modificar
	//		- nuevo password_hash
	// Errores que puede devuelve:
	// 		- ErrIdNoEncontrado
	ModificarContraseña(context.Context, int, string) error

	// Parámetros:
	// 		- email del usuario a verificar
	// Salida:
	// 		- Usuario encontrado
	// Errores que puede devuelve:
	// 		- ErrEmailNoEncontrado
	BuscarUsuarioPorMail(context.Context, string) (*domain.Usuario, error)
}

type userService struct {
	userManagerDb userManagerDbInterface
}

func NewUserService(userManager userManagerDbInterface) *userService {
	return &userService{
		userManagerDb: userManager,
	}
}

func (us *userService) CrearUsuario(ctx context.Context, email, password string) (int, error) {

	emailTrimSpace := strings.TrimSpace(email)
	passwordTrimSpace := strings.TrimSpace(password)

	if emailTrimSpace == "" {
		return 0, domain.ErrEmailRequerido
	}

	if passwordTrimSpace == "" {
		return 0, domain.ErrPasswordRequerido
	}

	if len(passwordTrimSpace) < config.LargoMinimoPassword {
		return 0, domain.ErrPasswordCorto
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(passwordTrimSpace),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	usuario := domain.Usuario{
		Email:         emailTrimSpace,
		Password_hash: string(hashedPassword),
	}

	return us.userManagerDb.GenerarNuevoUsuario(ctx, usuario)
}

func (us *userService) ModificarContraseña(ctx context.Context, id int, password string) error {

	if id < 1 {
		return domain.ErrIdNoValido
	}

	passwordTrimSpace := strings.TrimSpace(password)

	if passwordTrimSpace == "" {
		return domain.ErrPasswordRequerido
	}

	if len(passwordTrimSpace) < config.LargoMinimoPassword {
		return domain.ErrPasswordCorto
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(passwordTrimSpace),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return us.userManagerDb.ModificarContraseña(ctx, id, string(hashedPassword))

}

func (us *userService) ObtenerId(ctx context.Context, email, password string) (int, error) {

	emailTrimSpace := strings.TrimSpace(email)
	passwordTrimSpace := strings.TrimSpace(password)

	if emailTrimSpace == "" {
		return 0, domain.ErrEmailRequerido
	}

	if passwordTrimSpace == "" {
		return 0, domain.ErrPasswordRequerido
	}

	usuarioEncontrado, err := us.userManagerDb.BuscarUsuarioPorMail(ctx, emailTrimSpace)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usuarioEncontrado.Password_hash), []byte(passwordTrimSpace))
	if err != nil {
		return 0, domain.ErrPasswordIncorrecto
	}

	return *usuarioEncontrado.Id, nil

}
