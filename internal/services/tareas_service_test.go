package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"api-gestor-tareas/internal/domain"
)

type tareaManagerDbMock struct {
	BuscarEnTituloFunc       func(context.Context, string, int) ([]domain.Tarea, error)
	ListarFunc               func(context.Context, int, bool) ([]domain.Tarea, error)
	MarcarComoCompletadaFunc func(context.Context, int, *time.Time) error
	ModificarTituloFunc      func(context.Context, int, string) error
	RegistrarTareaFunc       func(context.Context, domain.Tarea) (int, error)
}

func (tmm *tareaManagerDbMock) BuscarEnTitulo(cxt context.Context, parametroABuscar string, idUsuario int) ([]domain.Tarea, error) {
	if tmm.BuscarEnTituloFunc == nil {
		return nil, fmt.Errorf("BuscarEnTituloFunc no implementado")
	}
	return tmm.BuscarEnTituloFunc(cxt, parametroABuscar, idUsuario)
}

func (tmm *tareaManagerDbMock) Listar(ctx context.Context, idUsuario int, soloPendientes bool) ([]domain.Tarea, error) {
	if tmm.ListarFunc == nil {
		return nil, fmt.Errorf("ListarFunc no implementado")
	}
	return tmm.ListarFunc(ctx, idUsuario, soloPendientes)
}

func (tmm *tareaManagerDbMock) MarcarComoCompletada(ctx context.Context, idTarea int, fecha *time.Time) error {
	if tmm.MarcarComoCompletadaFunc == nil {
		return fmt.Errorf("MarcarComoCompletadaFunc no implementado")
	}
	return tmm.MarcarComoCompletadaFunc(ctx, idTarea, fecha)
}

func (tmm *tareaManagerDbMock) ModificarTitulo(ctx context.Context, idTarea int, nuevoTitulo string) error {
	if tmm.ModificarTituloFunc == nil {
		return fmt.Errorf("ModificarTituloFunc no implementado")
	}
	return tmm.ModificarTituloFunc(ctx, idTarea, nuevoTitulo)
}

func (tmm *tareaManagerDbMock) RegistrarTarea(ctx context.Context, tarea domain.Tarea) (int, error) {
	if tmm.RegistrarTareaFunc == nil {
		return 0, fmt.Errorf("RegistrarTareaFunc no implementado")
	}
	return tmm.RegistrarTareaFunc(ctx, tarea)
}

func TestCrearTarea(t *testing.T) {

	tests := []struct {
		name          string
		mockSetup     func() *tareaManagerDbMock
		descripcion   string
		fechaCreacion time.Time
		idUsuario     int
		errorEsperado error
	}{
		{name: "Todo Ok",
			mockSetup: func() *tareaManagerDbMock {
				return &tareaManagerDbMock{
					RegistrarTareaFunc: func(ctx context.Context, t domain.Tarea) (int, error) {
						return 0, nil
					},
				}
			},
			descripcion:   "Descripción de prueba ",
			fechaCreacion: time.Date(2001, time.March, 27, 0, 0, 0, 0, time.UTC),
			idUsuario:     1,
			errorEsperado: nil,
		},
		{name: "Descripción en blanco",
			mockSetup: func() *tareaManagerDbMock {
				return &tareaManagerDbMock{
					RegistrarTareaFunc: func(ctx context.Context, t domain.Tarea) (int, error) {
						return 0, nil
					},
				}
			},
			descripcion:   "",
			fechaCreacion: time.Date(2001, time.March, 27, 0, 0, 0, 0, time.UTC),
			idUsuario:     1,
			errorEsperado: domain.ErrDescripcionRequerida,
		},
		{name: "Id no valido",
			mockSetup: func() *tareaManagerDbMock {
				return &tareaManagerDbMock{
					RegistrarTareaFunc: func(ctx context.Context, t domain.Tarea) (int, error) {
						return 0, nil
					},
				}
			},
			descripcion:   "Descripción de prueba ",
			fechaCreacion: time.Date(2001, time.March, 27, 0, 0, 0, 0, time.UTC),
			idUsuario:     -1,
			errorEsperado: domain.ErrIdNoValido,
		},
		{name: "Fecha no valida",
			mockSetup: func() *tareaManagerDbMock {
				return &tareaManagerDbMock{
					RegistrarTareaFunc: func(ctx context.Context, t domain.Tarea) (int, error) {
						return 0, nil
					},
				}
			},
			descripcion:   "Descripción de prueba ",
			fechaCreacion: time.Date(2101, time.March, 27, 0, 0, 0, 0, time.UTC),
			idUsuario:     1,
			errorEsperado: domain.ErrFechaNoValida,
		},
		{name: "Usuario no valido",
			mockSetup: func() *tareaManagerDbMock {
				return &tareaManagerDbMock{
					RegistrarTareaFunc: func(ctx context.Context, t domain.Tarea) (int, error) {
						return 0, domain.ErrUsuarioAsignadoNoexiste
					},
				}
			},
			descripcion:   "Descripción de prueba ",
			fechaCreacion: time.Date(2001, time.March, 27, 0, 0, 0, 0, time.UTC),
			idUsuario:     1,
			errorEsperado: domain.ErrUsuarioAsignadoNoexiste,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			service := NewTareaService(test.mockSetup())
			ctx := t.Context()

			_, err := service.CrearTarea(ctx, test.descripcion, test.fechaCreacion, test.idUsuario)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestListarTareas(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func() *tareaManagerDbMock
		idUsuario     int
		errorEsperado error
	}{
		{name: "Todo Ok",
			mockSetup: func() *tareaManagerDbMock {
				return &tareaManagerDbMock{
					ListarFunc: func(context.Context, int, bool) ([]domain.Tarea, error) {
						return nil, nil
					},
				}
			},
			idUsuario:     1,
			errorEsperado: nil,
		},
		{name: "Id no valido",
			mockSetup: func() *tareaManagerDbMock {
				return &tareaManagerDbMock{
					ListarFunc: func(context.Context, int, bool) ([]domain.Tarea, error) {
						return nil, nil
					},
				}
			},
			idUsuario:     -1,
			errorEsperado: domain.ErrIdNoValido,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			service := NewTareaService(test.mockSetup())
			ctx := t.Context()

			_, err := service.ListarTareas(ctx, test.idUsuario, true)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}

}
