package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"api-gestor-tareas/internal/domain"

	"github.com/gin-gonic/gin"
)

type tareaServiceMock struct {
	CrearTareaFunc           func(ctx context.Context, descripcion string, fecha time.Time, idUsuario int) (int, error)
	ListarTareasFunc         func(ctx context.Context, idUsuario int, completadas bool) ([]domain.Tarea, error)
	ModificarDescripcionFunc func(ctx context.Context, idTarea int, NuevaDescripcion string) error
	MarcarComoCompletadaFunc func(ctx context.Context, idTarea int, fechaCompletada *time.Time) error
	BuscarFunc               func(ctx context.Context, parametro string, idUsuario int) ([]domain.Tarea, error)
}

func (tsm *tareaServiceMock) CrearTarea(ctx context.Context, descripcion string, fecha time.Time, idUsuario int) (int, error) {
	if tsm.CrearTareaFunc == nil {
		return 0, fmt.Errorf("CrearTareaFunc no implementado")
	}
	return tsm.CrearTareaFunc(ctx, descripcion, fecha, idUsuario)
}

func (tsm *tareaServiceMock) ListarTareas(ctx context.Context, idUsuario int, completadas bool) ([]domain.Tarea, error) {
	if tsm.ListarTareasFunc == nil {
		return nil, fmt.Errorf("ListarTareasFunc no implementado")
	}
	return tsm.ListarTareasFunc(ctx, idUsuario, completadas)
}

func (tsm *tareaServiceMock) ModificarDescripcion(ctx context.Context, idTarea int, nuevaDescripcion string) error {
	if tsm.ModificarDescripcionFunc == nil {
		return fmt.Errorf("ModificarDescripcionFunc no implementado")
	}
	return tsm.ModificarDescripcionFunc(ctx, idTarea, nuevaDescripcion)
}

func (tsm *tareaServiceMock) MarcarComoCompletada(ctx context.Context, idTarea int, fechaCompletada *time.Time) error {
	if tsm.MarcarComoCompletadaFunc == nil {
		return fmt.Errorf("MarcarComoCompletadaFunc no implementado")
	}
	return tsm.MarcarComoCompletadaFunc(ctx, idTarea, fechaCompletada)
}

func (tsm *tareaServiceMock) Buscar(ctx context.Context, parametro string, idUsuario int) ([]domain.Tarea, error) {
	if tsm.BuscarFunc == nil {
		return nil, fmt.Errorf("BuscarFunc no implementado")
	}
	return tsm.BuscarFunc(ctx, parametro, idUsuario)
}

func TestNueva(t *testing.T) {

	tests := []struct {
		name           string
		mockSetup      func() *tareaServiceMock
		descripcion    string
		fecha          string
		id             int
		codigoEsperado int
	}{
		{
			name: "OK",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					CrearTareaFunc: func(ctx context.Context, descripcion string, fecha time.Time, idUsuario int) (int, error) {
						return 1, nil
					},
				}
			},
			descripcion:    "Descripcion de prueba",
			fecha:          "07/04/2026",
			id:             7,
			codigoEsperado: 200,
		},
		{
			name: "Error descripcion requerida",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					CrearTareaFunc: func(ctx context.Context, descripcion string, fecha time.Time, idUsuario int) (int, error) {
						return 0, domain.ErrDescripcionRequerida
					},
				}
			},
			descripcion:    "",
			fecha:          "07/04/2026",
			id:             7,
			codigoEsperado: 400,
		},
		{
			name: "Error ID invalido",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					CrearTareaFunc: func(ctx context.Context, descripcion string, fecha time.Time, idUsuario int) (int, error) {
						return 0, domain.ErrIdNoValido
					},
				}
			},
			descripcion:    "Descripcion",
			fecha:          "07/04/2026",
			id:             -1,
			codigoEsperado: 400,
		},
		{
			name: "Error fecha invalida",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					CrearTareaFunc: func(ctx context.Context, descripcion string, fecha time.Time, idUsuario int) (int, error) {
						return 0, domain.ErrFechaNoValida
					},
				}
			},
			descripcion:    "Descripcion",
			fecha:          "fecha-mal-formada",
			id:             7,
			codigoEsperado: 400,
		},
		{
			name: "Usuario no existe",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					CrearTareaFunc: func(ctx context.Context, descripcion string, fecha time.Time, idUsuario int) (int, error) {
						return 0, domain.ErrUsuarioAsignadoNoexiste
					},
				}
			},
			descripcion:    "Descripcion",
			fecha:          "07/04/2026",
			id:             999,
			codigoEsperado: 404,
		},
		{
			name: "Error interno",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					CrearTareaFunc: func(ctx context.Context, descripcion string, fecha time.Time, idUsuario int) (int, error) {
						return 0, fmt.Errorf("error inesperado")
					},
				}
			},
			descripcion:    "Descripcion",
			fecha:          "07/04/2026",
			id:             7,
			codigoEsperado: 500,
		},
	}

	for _, test := range tests {

		gin.SetMode(gin.TestMode)

		t.Run(test.name, func(t *testing.T) {

			handler := NewTareaHandler(test.mockSetup())

			// armo server
			router := gin.Default()
			router.POST("/test", func(c *gin.Context) {
				c.Set("user_id", test.id) // le mockeo el user_id
				handler.Nueva(c)
			})

			// armo consulta
			body := fmt.Sprintf(`{"descripcion":"%s","fecha":"%s"}`, test.descripcion, test.fecha)

			//realizo la solicitud
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// capturo la respuesta
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.codigoEsperado {
				t.Errorf("Error no esperado.\nSe esperaba: %d \nSe obtuvo: %d \nCuerpo: %s",
					test.codigoEsperado, w.Code, w.Body)
			}
		})
	}
}

func TestListar(t *testing.T) {

	tests := []struct {
		name           string
		mockSetup      func() *tareaServiceMock
		id             int
		completadas    bool
		codigoEsperado int
	}{
		{
			name: "OK",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					ListarTareasFunc: func(ctx context.Context, idUsuario int, completadas bool) ([]domain.Tarea, error) {
						return []domain.Tarea{}, nil
					},
				}
			},
			id:             7,
			completadas:    false,
			codigoEsperado: 200,
		},
		{
			name: "Error ID invalido",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					ListarTareasFunc: func(ctx context.Context, idUsuario int, completadas bool) ([]domain.Tarea, error) {
						return nil, domain.ErrIdNoValido
					},
				}
			},
			id:             -1,
			completadas:    false,
			codigoEsperado: 400,
		},
		{
			name: "Error interno",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					ListarTareasFunc: func(ctx context.Context, idUsuario int, completadas bool) ([]domain.Tarea, error) {
						return nil, fmt.Errorf("error inesperado")
					},
				}
			},
			id:             7,
			completadas:    false,
			codigoEsperado: 500,
		},
	}

	for _, test := range tests {

		gin.SetMode(gin.TestMode)

		t.Run(test.name, func(t *testing.T) {

			handler := NewTareaHandler(test.mockSetup())

			// armo server
			router := gin.Default()

			router.POST("/test", func(c *gin.Context) {
				c.Set("user_id", test.id) // le mockeo el user_id
				handler.Listar(c)
			})

			// armo el path con query
			path := fmt.Sprintf(`/test?completadas=%t}`, test.completadas)

			//realizo la solicitud
			req := httptest.NewRequest(http.MethodPost, path, nil)
			req.Header.Set("Content-Type", "application/json")

			// capturo la respuesta
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.codigoEsperado {
				t.Errorf("Error no esperado.\nSe esperaba: %d \nSe obtuvo: %d \nCuerpo: %s",
					test.codigoEsperado, w.Code, w.Body)
			}
		})
	}
}

func TestModificar(t *testing.T) {

	tests := []struct {
		name           string
		mockSetup      func() *tareaServiceMock
		descripcion    string
		id             int
		codigoEsperado int
	}{
		{
			name: "OK",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					ModificarDescripcionFunc: func(ctx context.Context, idTarea int, descripcion string) error {
						return nil
					},
				}
			},
			descripcion:    "Nueva Descripcion de prueba",
			id:             7,
			codigoEsperado: 200,
		},
		{
			name: "Error - Id no válido",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					ModificarDescripcionFunc: func(ctx context.Context, idTarea int, descripcion string) error {
						return domain.ErrIdNoValido
					},
				}
			},
			descripcion:    "Descripcion",
			id:             -1,
			codigoEsperado: 400,
		},
		{
			name: "Error - Tarea no existe",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					ModificarDescripcionFunc: func(ctx context.Context, idTarea int, descripcion string) error {
						return domain.ErrTareaNoExiste
					},
				}
			},
			descripcion:    "Descripcion",
			id:             999,
			codigoEsperado: 404,
		},
		{
			name: "Error - Interno",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					ModificarDescripcionFunc: func(ctx context.Context, idTarea int, descripcion string) error {
						return errors.New("unexpected error")
					},
				}
			},
			descripcion:    "Descripcion",
			id:             1,
			codigoEsperado: 500,
		},
	}

	for _, test := range tests {

		gin.SetMode(gin.TestMode)

		t.Run(test.name, func(t *testing.T) {

			handler := NewTareaHandler(test.mockSetup())

			// armo server
			router := gin.Default()
			router.POST("/test", handler.Modificar)

			// armo consulta
			body := fmt.Sprintf(`{"nueva_descripcion":"%s", "id_tarea":%d}`, test.descripcion, test.id)

			//realizo la solicitud
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// capturo la respuesta
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.codigoEsperado {
				t.Errorf("Error no esperado.\nSe esperaba: %d \nSe obtuvo: %d \nCuerpo: %s",
					test.codigoEsperado, w.Code, w.Body)
			}
		})
	}
}

func TestFinalizar(t *testing.T) {

	tests := []struct {
		name           string
		mockSetup      func() *tareaServiceMock
		fecha          string
		id             int
		codigoEsperado int
	}{
		{
			name: "OK",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					MarcarComoCompletadaFunc: func(ctx context.Context, idTarea int, fecha *time.Time) error {
						return nil
					},
				}
			},
			fecha:          "08/04/2026",
			id:             7,
			codigoEsperado: 200,
		},
		{
			name: "Error - Id no válido",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					MarcarComoCompletadaFunc: func(ctx context.Context, idTarea int, fecha *time.Time) error {
						return domain.ErrIdNoValido
					},
				}
			},
			fecha:          "08/04/2026",
			id:             -1,
			codigoEsperado: 400,
		},
		{
			name: "Error - Fecha no válida",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					MarcarComoCompletadaFunc: func(ctx context.Context, idTarea int, fecha *time.Time) error {
						return domain.ErrFechaNoValida
					},
				}
			},
			fecha:          "fecha-invalida",
			id:             7,
			codigoEsperado: 400,
		},
		{
			name: "Error - Tarea no existe",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					MarcarComoCompletadaFunc: func(ctx context.Context, idTarea int, fecha *time.Time) error {
						return domain.ErrTareaNoExiste
					},
				}
			},
			fecha:          "08/04/2026",
			id:             999,
			codigoEsperado: 404,
		},
		{
			name: "Error - Interno",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					MarcarComoCompletadaFunc: func(ctx context.Context, idTarea int, fecha *time.Time) error {
						return errors.New("unexpected error")
					},
				}
			},
			fecha:          "08/04/2026",
			id:             1,
			codigoEsperado: 500,
		},
	}

	for _, test := range tests {

		gin.SetMode(gin.TestMode)

		t.Run(test.name, func(t *testing.T) {

			handler := NewTareaHandler(test.mockSetup())

			// armo server
			router := gin.Default()
			router.POST("/test", handler.Finalizar)

			// armo consulta
			body := fmt.Sprintf(`{"fecha":"%s", "id_tarea":%d}`, test.fecha, test.id)

			//realizo la solicitud
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// capturo la respuesta
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.codigoEsperado {
				t.Errorf("Error no esperado.\nSe esperaba: %d \nSe obtuvo: %d \nCuerpo: %s",
					test.codigoEsperado, w.Code, w.Body)
			}
		})
	}
}

func TestBuscar(t *testing.T) {

	tests := []struct {
		name                string
		mockSetup           func() *tareaServiceMock
		parametroDeBusqueda string
		id                  int
		codigoEsperado      int
	}{
		{
			name: "OK",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					BuscarFunc: func(ctx context.Context, parametro string, idUsuario int) ([]domain.Tarea, error) {
						return []domain.Tarea{}, nil
					},
				}
			},
			parametroDeBusqueda: "BusquedaDePrueba",
			id:                  7,
			codigoEsperado:      200,
		},
		{
			name: "Error - Id no válido",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					BuscarFunc: func(ctx context.Context, parametro string, idUsuario int) ([]domain.Tarea, error) {
						return nil, domain.ErrIdNoValido
					},
				}
			},
			parametroDeBusqueda: "BusquedaDePrueba",
			id:                  -1,
			codigoEsperado:      400,
		},
		{
			name: "Error - Parámetro de búsqueda no válido",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					BuscarFunc: func(ctx context.Context, parametro string, idUsuario int) ([]domain.Tarea, error) {
						return nil, domain.ErrParametroDeBusquedaNoValido
					},
				}
			},
			parametroDeBusqueda: "",
			id:                  7,
			codigoEsperado:      400,
		},
		{
			name: "Error - Interno",
			mockSetup: func() *tareaServiceMock {
				return &tareaServiceMock{
					BuscarFunc: func(ctx context.Context, parametro string, idUsuario int) ([]domain.Tarea, error) {
						return nil, errors.New("unexpected error")
					},
				}
			},
			parametroDeBusqueda: "Busqueda",
			id:                  1,
			codigoEsperado:      500,
		},
	}

	for _, test := range tests {

		gin.SetMode(gin.TestMode)

		t.Run(test.name, func(t *testing.T) {

			handler := NewTareaHandler(test.mockSetup())

			// armo server
			router := gin.Default()
			router.POST("/test", func(c *gin.Context) {
				c.Set("user_id", test.id) // le mockeo el user_id
				handler.Buscar(c)
			})

			// armo el path con query
			path := fmt.Sprintf(`/test?parametro_busqueda=%s}`, test.parametroDeBusqueda)

			//realizo la solicitud
			req := httptest.NewRequest(http.MethodPost, path, nil)
			req.Header.Set("Content-Type", "application/json")

			// capturo la respuesta
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.codigoEsperado {
				t.Errorf("Error no esperado.\nSe esperaba: %d \nSe obtuvo: %d \nCuerpo: %s",
					test.codigoEsperado, w.Code, w.Body)
			}
		})
	}
}
