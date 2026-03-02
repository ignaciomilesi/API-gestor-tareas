package services

import (
	"api-gestor-tareas/internal/domain"
	"fmt"

	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userManagerInterface interface {
	// Errores que devuelve:
	// 		ErrEmailDuplicado
	GenerarNuevoUsuario(context.Context, domain.Usuario) (int, error)

	// Errores que devuelve:
	// 		ErrIdNoEncontrado
	ModifcarContraseña(context.Context, int, string) error

	// Errores que devuelve:
	// 		ErrEmailNoEncontrado
	//		ErrPasswordIncorrecto
	ObternerId(context.Context, domain.Usuario) (int, error)
}

type userService struct {
	userManagerDb userManagerInterface
}

func NewUserService(userManager userManagerInterface) userService {
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
			return errors.New("")
		}
		return fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return nil
}
