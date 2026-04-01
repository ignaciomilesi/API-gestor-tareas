package services

import (
	"api-gestor-tareas/internal/domain"
	"errors"
	"fmt"

	"context"
	"testing"
)

type userManagerDbMock struct {
	GenerarNuevoUsuarioFunc  func(cxt context.Context, nuevoUsario domain.Usuario) (int, error)
	ModificarContraseñaFunc  func(cxt context.Context, id int, nuevoPassword string) error
	BuscarUsuarioPorMailFunc func(cxt context.Context, mail string) (*domain.Usuario, error)
}

func (umm *userManagerDbMock) GenerarNuevoUsuario(cxt context.Context, nuevoUsario domain.Usuario) (int, error) {
	if umm.GenerarNuevoUsuarioFunc == nil {
		return 0, fmt.Errorf("GenerarNuevoUsuarioFunc no implementado")
	}
	return umm.GenerarNuevoUsuarioFunc(cxt, nuevoUsario)
}

func (umm *userManagerDbMock) ModificarContraseña(cxt context.Context, id int, nuevoPassword string) error {
	if umm.ModificarContraseñaFunc == nil {
		return fmt.Errorf("ModificarContraseñaFunc no implementado")
	}
	return umm.ModificarContraseñaFunc(cxt, id, nuevoPassword)
}

func (umm *userManagerDbMock) BuscarUsuarioPorMail(cxt context.Context, mail string) (*domain.Usuario, error) {
	if umm.BuscarUsuarioPorMailFunc == nil {
		return nil, fmt.Errorf("BuscarUsuarioPorMailFunc no implementado")
	}
	return umm.BuscarUsuarioPorMailFunc(cxt, mail)
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
		{
			name: "Password corto",
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					GenerarNuevoUsuarioFunc: func(cxt context.Context, nuevoUsario domain.Usuario) (int, error) {
						return 0, nil
					},
				}
			},
			email:         "email_de_prueba@prueba.com",
			password:      "12345",
			errorEsperado: domain.ErrPasswordCorto,
		},
		{
			name: "Usuario duplicado",
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					GenerarNuevoUsuarioFunc: func(cxt context.Context, nuevoUsario domain.Usuario) (int, error) {
						return 0, domain.ErrEmailDuplicado
					},
				}
			},
			email:         "email_de_prueba@prueba.com",
			password:      "1234567",
			errorEsperado: domain.ErrEmailDuplicado,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			service := NewUserService(test.mockSetup())
			ctx := t.Context()

			_, err := service.CrearUsuario(ctx, test.email, test.password)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}

}

func TestModificarContraseña(t *testing.T) {

	tests := []struct {
		name          string
		mockSetup     func() *userManagerDbMock
		id            int
		password      string
		errorEsperado error
	}{
		{
			name:          "Id no valido",
			id:            -1,
			password:      "123456",
			errorEsperado: domain.ErrIdNoValido,
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					ModificarContraseñaFunc: func(cxt context.Context, id int, nuevoPassword string) error {
						return nil
					},
				}
			},
		},
		{
			name:          "Password en blanco",
			id:            1,
			password:      "",
			errorEsperado: domain.ErrPasswordRequerido,
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					ModificarContraseñaFunc: func(cxt context.Context, id int, nuevoPassword string) error {
						return nil
					},
				}
			},
		},
		{
			name:          "Password corto",
			id:            1,
			password:      "1234",
			errorEsperado: domain.ErrPasswordCorto,
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					ModificarContraseñaFunc: func(cxt context.Context, id int, nuevoPassword string) error {
						return nil
					},
				}
			},
		},
		{
			name:          "Id no encontrado en la base de datos",
			id:            1,
			password:      "123456",
			errorEsperado: domain.ErrIdNoEncontrado,
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					ModificarContraseñaFunc: func(cxt context.Context, id int, nuevoPassword string) error {
						return domain.ErrIdNoEncontrado
					},
				}
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			service := NewUserService(test.mockSetup())
			ctx := t.Context()

			err := service.ModificarContraseña(ctx, test.id, test.password)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}

}

func TestObtenerId(t *testing.T) {

	tests := []struct {
		name          string
		mockSetup     func() *userManagerDbMock
		email         string
		password      string
		errorEsperado error
	}{
		{
			name:          "Email en blanco",
			email:         "",
			password:      "123456",
			errorEsperado: domain.ErrEmailRequerido,
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					BuscarUsuarioPorMailFunc: func(cxt context.Context, mail string) (*domain.Usuario, error) {
						return nil, nil
					},
				}
			},
		},
		{
			name:          "Password en blanco",
			email:         "email_de_prueba@prueba.com",
			password:      "",
			errorEsperado: domain.ErrPasswordRequerido,
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					BuscarUsuarioPorMailFunc: func(cxt context.Context, mail string) (*domain.Usuario, error) {
						return nil, nil
					},
				}
			},
		},
		{
			name:          "Password corto",
			email:         "email_de_prueba@prueba.com",
			password:      "1234",
			errorEsperado: domain.ErrPasswordCorto,
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					BuscarUsuarioPorMailFunc: func(cxt context.Context, mail string) (*domain.Usuario, error) {
						return nil, nil
					},
				}
			},
		},
		{
			name:          "Email no encontrado",
			email:         "email_de_prueba@prueba.com",
			password:      "123456",
			errorEsperado: domain.ErrEmailNoEncontrado,
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					BuscarUsuarioPorMailFunc: func(cxt context.Context, mail string) (*domain.Usuario, error) {
						return nil, domain.ErrEmailNoEncontrado
					},
				}
			},
		},
		{
			name:          "Password incorrecto",
			email:         "email_de_prueba@prueba.com",
			password:      "123456",
			errorEsperado: domain.ErrPasswordIncorrecto,
			mockSetup: func() *userManagerDbMock {
				return &userManagerDbMock{
					BuscarUsuarioPorMailFunc: func(cxt context.Context, mail string) (*domain.Usuario, error) {
						return &domain.Usuario{
							Password_hash: "abc",
						}, nil
					},
				}
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			service := NewUserService(test.mockSetup())
			ctx := t.Context()

			_, err := service.ObtenerId(ctx, test.email, test.password)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}

}
