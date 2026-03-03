package services

import (
	"api-gestor-tareas/internal/domain"
	"errors"
	"fmt"

	"context"
	"testing"
)

type userManagerDbMock struct {
	GenerarNuevoUsuarioFunc func(cxt context.Context, nuevoUsario domain.Usuario) (int, error)
	ModifcarContraseñaFunc  func(cxt context.Context, id int, nuevoPassword string) error
	ObternerIdFunc          func(cxt context.Context, usuarioAVerificar domain.Usuario) (int, error)
}

func (umm *userManagerDbMock) GenerarNuevoUsuario(cxt context.Context, nuevoUsario domain.Usuario) (int, error) {
	if umm.GenerarNuevoUsuarioFunc == nil {
		return 0, fmt.Errorf("GenerarNuevoUsuarioFunc no implementado")
	}
	return umm.GenerarNuevoUsuarioFunc(cxt, nuevoUsario)
}

func (umm *userManagerDbMock) ModifcarContraseña(cxt context.Context, id int, nuevoPassword string) error {
	if umm.ModifcarContraseñaFunc == nil {
		return fmt.Errorf("ModifcarContraseñaFunc no implementado")
	}
	return umm.ModifcarContraseñaFunc(cxt, id, nuevoPassword)
}

func (umm *userManagerDbMock) ObternerId(cxt context.Context, usuarioAVerificar domain.Usuario) (int, error) {
	if umm.ObternerIdFunc == nil {
		return 0, fmt.Errorf("ObternerIdFunc no implementado")
	}
	return umm.ObternerIdFunc(cxt, usuarioAVerificar)
}

func TestCrearUsuario(t *testing.T) {

	tests := []struct {
		name          string
		mockSetup     func() *userManagerDbMock
		email         string
		password      string
		errorEsperado error
	}{
		{
			name: "Mail en blanco",
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					GenerarNuevoUsuarioFunc: func(cxt context.Context, nuevoUsario domain.Usuario) (int, error) {
						return 0, nil
					},
				}
			},
			email:         "",
			password:      "password_de_prueba",
			errorEsperado: domain.ErrEmailRequerido,
		},
		{
			name: "Password en blanco",
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					GenerarNuevoUsuarioFunc: func(cxt context.Context, nuevoUsario domain.Usuario) (int, error) {
						return 0, nil
					},
				}
			},
			email:         "email_de_prueba@prueba.com",
			password:      "",
			errorEsperado: domain.ErrPasswordRequerido,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			service := NewUserService(test.mockSetup())
			ctx := t.Context()

			err := service.CrearUsuario(ctx, test.email, test.password)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf(" -- Error no esperado: se esperaba %t pero se obtuvo %t\n\n",
					test.errorEsperado, err)
			}
		})
	}

}
