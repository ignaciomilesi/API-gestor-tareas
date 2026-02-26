package db

import (
	"api-gestor-tareas/internal/models"
	"errors"
	"testing"
	"time"
)

func TestRegistrarTarea(t *testing.T) {

	tests := []struct {
		name          string
		tareaTest     models.Tarea
		errorEsperado error
	}{
		{
			name: "Registro de tarea ok",
			tareaTest: models.Tarea{
				Titulo:         "Tarea de prueba2",
				Fecha_creacion: time.Now(),
				Completada:     true,
				Id_usuario:     1,
			},
			errorEsperado: nil,
		},
		{
			name: "Asignar tarea a un usuario inexistente",
			tareaTest: models.Tarea{
				Titulo:         "Tarea de prueba",
				Fecha_creacion: time.Now(),
				Completada:     false,
				Id_usuario:     9999,
			},
			errorEsperado: ErrUsuarioAsignadoNoexiste,
		},
	}

	db := conectarBaseDeDatosTest()
	tm := NewTareasManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			_, err := tm.RegistrarTarea(ctx, test.tareaTest)
			if !errors.Is(err, test.errorEsperado) {

				t.Errorf(" -- Error no esperado: se esperaba %t pero se obtuvo %t\n\n",
					test.errorEsperado, err)
			}
		})
	}

}

func TestListar(t *testing.T) {

	tests := []struct {
		name          string
		IdUsuario     int
		completada    bool
		errorEsperado error
	}{
		{
			name:          "Traer listado tareas pendientes",
			IdUsuario:     1,
			errorEsperado: nil,
		},
		{
			name:          "Traer listado tareas completadas",
			IdUsuario:     1,
			completada:    true,
			errorEsperado: nil,
		},
		{
			name:          "Solicitar listado de usuario inexistente",
			IdUsuario:     999,
			errorEsperado: nil,
		},
	}

	db := conectarBaseDeDatosTest()
	tm := NewTareasManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			_, err := tm.Listar(ctx, test.IdUsuario, test.completada)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf(" -- Error no esperado: se esperaba %t pero se obtuvo %t\n\n",
					test.errorEsperado, err)
			}
		})
	}

}

func TestModificarTitulo(t *testing.T) {

	tests := []struct {
		name          string
		IdTarea       int
		errorEsperado error
	}{
		{
			name:          "Modificar tarea Ok",
			IdTarea:       5,
			errorEsperado: nil,
		},
		{
			name:          "Modificar tarea inexistente",
			IdTarea:       999,
			errorEsperado: ErrTareaNoExiste,
		},
	}

	db := conectarBaseDeDatosTest()
	tm := NewTareasManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			err := tm.ModificarTitulo(ctx, test.IdTarea, "Texto modificado")

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado:\n  - se esperaba: %t\n  - se obtuvo %t\n\n",
					test.errorEsperado, err)
			}
		})
	}

}

func TestMarcarComoCompletada(t *testing.T) {

	tests := []struct {
		name          string
		IdTarea       int
		errorEsperado error
	}{
		{
			name:          "Modificar tarea Ok",
			IdTarea:       3,
			errorEsperado: nil,
		},
		{
			name:          "Modificar tarea inexistente",
			IdTarea:       999,
			errorEsperado: ErrTareaNoExiste,
		},
	}

	db := conectarBaseDeDatosTest()
	tm := NewTareasManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			hoy := time.Now()
			err := tm.MarcarComoCompletada(ctx, test.IdTarea, &hoy)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado:\n  - se esperaba: %t\n  - se obtuvo %t\n\n",
					test.errorEsperado, err)
			}
		})
	}

}

func TestBuscarEnTitulo(t *testing.T) {

	tests := []struct {
		name               string
		palabraBuscada     string
		idUsuario          int
		errorEsperado      error
		seEsperaResultados bool
	}{
		{
			name:               "Buscar tarea Ok",
			palabraBuscada:     "prueba",
			idUsuario:          1,
			errorEsperado:      nil,
			seEsperaResultados: true,
		},
		{
			name:               "Buscar tarea en usuario no valido",
			palabraBuscada:     "prueba",
			idUsuario:          999,
			errorEsperado:      nil,
			seEsperaResultados: false,
		},
		{
			name:               "No se encuentra palabra",
			palabraBuscada:     "zzz",
			idUsuario:          1,
			errorEsperado:      nil,
			seEsperaResultados: false,
		},
	}

	db := conectarBaseDeDatosTest()
	tm := NewTareasManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			lista, err := tm.BuscarEnTitulo(ctx, test.palabraBuscada, test.idUsuario)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado:\n  - se esperaba: %t\n  - se obtuvo: %t\n\n",
					test.errorEsperado, err)
			}

			if len(lista) > 0 && !test.seEsperaResultados {
				t.Errorf("Error:\n  - se esperaban resultados?: %t\n  - se obtuvo %b resultados\n\n",
					test.seEsperaResultados, len(lista))
			}
		})
	}

}
