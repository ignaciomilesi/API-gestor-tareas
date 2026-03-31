package handlers

import (
	"api-gestor-tareas/internal/domain"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type userServiceMock struct {
	CrearUsuariofunc        func(ctx context.Context, email string, password string) (int, error)
	ModificarContraseñafunc func(ctx context.Context, id int, password string) error
	ObtenerIdfunc           func(ctx context.Context, email string, password string) (int, error)
}

func (usm *userServiceMock) CrearUsuario(ctx context.Context, email string, password string) (int, error) {
	if usm.CrearUsuariofunc == nil {
		return 0, fmt.Errorf("CrearUsuarioFunc no implementado")
	}
	return usm.CrearUsuariofunc(ctx, email, password)
}

func (usm *userServiceMock) ModificarContraseña(ctx context.Context, id int, password string) error {
	if usm.ModificarContraseñafunc == nil {
		return fmt.Errorf("ModificarContraseñaFunc no implementado")
	}
	return usm.ModificarContraseñafunc(ctx, id, password)
}

func (usm *userServiceMock) ObtenerId(ctx context.Context, email string, password string) (int, error) {
	if usm.ObtenerIdfunc == nil {
		return 0, fmt.Errorf("ObtenerIdFunc no implementado")
	}
	return usm.ObtenerIdfunc(ctx, email, password)
}

func TestSingin(t *testing.T) {

	tests := []struct {
		name           string
		mockSetup      func() *userServiceMock
		email          string
		password       string
		codigoEsperado int
	}{
		{
			name: "OK",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					CrearUsuariofunc: func(ctx context.Context, email string, password string) (int, error) {
						return 1, nil
					},
				}
			},
			email:          "email_prueba@gmail.com",
			password:       "123456789",
			codigoEsperado: 200,
		},
		{
			name: "Error email requerido",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					CrearUsuariofunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, domain.ErrEmailRequerido
					},
				}
			},
			email:          "",
			password:       "123456789",
			codigoEsperado: 400,
		},
		{
			name: "Error password requerido",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					CrearUsuariofunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, domain.ErrPasswordRequerido
					},
				}
			},
			email:          "test@mail.com",
			password:       "",
			codigoEsperado: 400,
		},
		{
			name: "Error password corto",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					CrearUsuariofunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, domain.ErrPasswordCorto
					},
				}
			},
			email:          "test@mail.com",
			password:       "123",
			codigoEsperado: 400,
		},
		{
			name: "Error email duplicado",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					CrearUsuariofunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, domain.ErrEmailDuplicado
					},
				}
			},
			email:          "test@mail.com",
			password:       "123456",
			codigoEsperado: 409,
		},
		{
			name: "Error interno",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					CrearUsuariofunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, errors.New("algo explotó")
					},
				}
			},
			email:          "test@mail.com",
			password:       "123456",
			codigoEsperado: 500,
		},
	}

	for _, test := range tests {

		gin.SetMode(gin.TestMode)

		t.Run(test.name, func(t *testing.T) {

			handler := NewUserHandler(test.mockSetup())

			// armo server
			router := gin.Default()
			router.POST("/singin", handler.Singin)

			// armo consulta
			body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, test.email, test.password)

			//realizo la solicitud
			req := httptest.NewRequest(http.MethodPost, "/singin", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// capturo la respuesta
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.codigoEsperado {
				t.Errorf("Error no esperado.\nSe esperaba: \n --- %d \nse obtuvo: \n --- %d",
					test.codigoEsperado, w.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {

	tests := []struct {
		name           string
		mockSetup      func() *userServiceMock
		email          string
		password       string
		codigoEsperado int
	}{
		{
			name: "OK",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ObtenerIdfunc: func(ctx context.Context, email string, password string) (int, error) {
						return 1, nil
					},
				}
			},
			email:          "email_prueba@gmail.com",
			password:       "123456789",
			codigoEsperado: 200,
		},
		{
			name: "Error email requerido",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ObtenerIdfunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, domain.ErrEmailRequerido
					},
				}
			},
			email:          "",
			password:       "123456",
			codigoEsperado: 400,
		},
		{
			name: "Error password requerido",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ObtenerIdfunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, domain.ErrPasswordRequerido
					},
				}
			},
			email:          "test@mail.com",
			password:       "",
			codigoEsperado: 400,
		},
		{
			name: "Error email no encontrado",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ObtenerIdfunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, domain.ErrEmailNoEncontrado
					},
				}
			},
			email:          "noexiste@mail.com",
			password:       "123456",
			codigoEsperado: 401,
		},
		{
			name: "Error password incorrecto",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ObtenerIdfunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, domain.ErrPasswordIncorrecto
					},
				}
			},
			email:          "test@mail.com",
			password:       "wrongpass",
			codigoEsperado: 401,
		},
		{
			name: "Error interno",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ObtenerIdfunc: func(ctx context.Context, email string, password string) (int, error) {
						return 0, errors.New("fallo inesperado")
					},
				}
			},
			email:          "test@mail.com",
			password:       "123456",
			codigoEsperado: 500,
		},
	}

	for _, test := range tests {

		gin.SetMode(gin.TestMode)

		t.Run(test.name, func(t *testing.T) {

			handler := NewUserHandler(test.mockSetup())

			// armo server
			router := gin.Default()
			router.POST("/login", handler.Login)

			// armo consulta
			body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, test.email, test.password)

			//realizo la solicitud
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// capturo la respuesta
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.codigoEsperado {
				t.Errorf("Error no esperado.\nSe esperaba: \n --- %d \nse obtuvo: \n --- %d",
					test.codigoEsperado, w.Code)
			}
		})
	}
}

func TestActualizarContraseña(t *testing.T) {

	tests := []struct {
		name           string
		mockSetup      func() *userServiceMock
		id             int
		password       string
		codigoEsperado int
	}{
		{
			name: "OK",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ModificarContraseñafunc: func(ctx context.Context, id int, password string) error {
						return nil
					},
				}
			},
			id:             1,
			password:       "12345678",
			codigoEsperado: 200,
		},
		{
			name: "Error id no válido",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ModificarContraseñafunc: func(ctx context.Context, id int, password string) error {
						return domain.ErrIdNoValido
					},
				}
			},
			id:             -1,
			password:       "12345678",
			codigoEsperado: 400,
		},
		{
			name: "Error password requerido",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ModificarContraseñafunc: func(ctx context.Context, id int, password string) error {
						return domain.ErrPasswordRequerido
					},
				}
			},
			id:             1,
			password:       "",
			codigoEsperado: 400,
		},
		{
			name: "Error password corto",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ModificarContraseñafunc: func(ctx context.Context, id int, password string) error {
						return domain.ErrPasswordCorto
					},
				}
			},
			id:             1,
			password:       "123",
			codigoEsperado: 400,
		},
		{
			name: "Error id no encontrado",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ModificarContraseñafunc: func(ctx context.Context, id int, password string) error {
						return domain.ErrIdNoEncontrado
					},
				}
			},
			id:             999,
			password:       "12345678",
			codigoEsperado: 401,
		},
		{
			name: "Error interno",
			mockSetup: func() *userServiceMock {
				return &userServiceMock{
					ModificarContraseñafunc: func(ctx context.Context, id int, password string) error {
						return errors.New("fallo inesperado")
					},
				}
			},
			id:             1,
			password:       "12345678",
			codigoEsperado: 500,
		},
	}

	for _, test := range tests {

		gin.SetMode(gin.TestMode)

		t.Run(test.name, func(t *testing.T) {

			handler := NewUserHandler(test.mockSetup())

			// armo server
			router := gin.Default()
			router.POST("/actualizarContraseña", handler.ActualizarContraseña)

			// armo consulta
			body := fmt.Sprintf(`{"id":%d,"password":"%s"}`, test.id, test.password)

			//realizo la solicitud
			req := httptest.NewRequest(http.MethodPost, "/actualizarContraseña", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// capturo la respuesta
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.codigoEsperado {
				t.Errorf("Error no esperado.\nSe esperaba: \n --- %d \nse obtuvo: \n --- %d",
					test.codigoEsperado, w.Code)
			}
		})
	}
}
