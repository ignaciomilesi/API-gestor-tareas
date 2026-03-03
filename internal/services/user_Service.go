package services

import (
	"api-gestor-tareas/internal/domain"
	"fmt"

	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userManagerDbInterface interface {
	// Parámetros:
	// 		- Usuario a generar
	// Errores que devuelve:
	// 		- ErrEmailDuplicado
	GenerarNuevoUsuario(context.Context, domain.Usuario) (int, error)

	// Parámetros:
	// 		- id del usuario a modificar
	//		- nuevo password_hash
	// Errores que devuelve:
	// 		- ErrIdNoEncontrado
	ModifcarContraseña(context.Context, int, string) error

	// Parámetros:
	// 		- Usuario al que quiero obtener el id
	// Errores que devuelve:
	// 		- ErrEmailNoEncontrado
	//		- ErrPasswordIncorrecto
	ObternerId(context.Context, domain.Usuario) (int, error)
}

type userService struct {
	userManagerDb userManagerDbInterface
}

func NewUserService(userManager userManagerDbInterface) userService {
	return userService{
		userManagerDb: userManager,
	}
}

func (us *userService) CrearUsuario(ctx context.Context, email, password string) error {

	emailTrimSpace := strings.TrimSpace(email)
	passwordTrimSpace := strings.TrimSpace(password)

	if emailTrimSpace == "" {
		return domain.ErrEmailRequerido
	}

	if passwordTrimSpace == "" {
		return domain.ErrPasswordRequerido
	}

	if len(passwordTrimSpace) < 6 {
		return domain.ErrPasswordCorto
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(passwordTrimSpace),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	usuario := domain.Usuario{
		Email:         emailTrimSpace,
		Password_hash: string(hashedPassword),
	}

	_, erra := us.userManagerDb.GenerarNuevoUsuario(ctx, usuario)

	if erra != nil {
		if errors.Is(err, domain.ErrEmailDuplicado) {
			return domain.ErrUsuarioExiste
		}
		return fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return nil
}

func (us *userService) ModificarContraseña(ctx context.Context, id int, password string) error {

	if id < 1 {
		return domain.ErrIdNovalido
	}

	passwordTrimSpace := strings.TrimSpace(password)

	if passwordTrimSpace == "" {
		return domain.ErrPasswordRequerido
	}

	if len(passwordTrimSpace) < 6 {
		return domain.ErrPasswordCorto
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(passwordTrimSpace),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	err = us.userManagerDb.ModifcarContraseña(ctx, id, string(hashedPassword))

	if err != nil {
		if errors.Is(err, domain.ErrIdNoEncontrado) {
			return domain.ErrIdNovalido
		}
		return fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return nil
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

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(passwordTrimSpace),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	usuarioAComprobar := domain.Usuario{
		Email:         emailTrimSpace,
		Password_hash: string(hashedPassword),
	}

	idEncontrado, erra := us.userManagerDb.ObternerId(ctx, usuarioAComprobar)

	if erra != nil {
		if errors.Is(err, domain.ErrEmailNoEncontrado) || errors.Is(err, domain.ErrPasswordIncorrecto) {
			return 0, domain.ErrUsuarioNoExiste
		}
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return idEncontrado, nil
}
