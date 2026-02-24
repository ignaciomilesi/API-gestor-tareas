package db

import (
	"api-gestor-tareas/internal/models"
	"testing"
	"time"
)

func TestRegistrarTarea(t *testing.T) {

	tests := []struct {
		name          string
		tareaTest     models.Tarea
		errorEsperado bool
	}{
		{
			name: "Registro de tarea ok",
			tareaTest: models.Tarea{
				Titulo:           "Tarea de prueba",
				Fecha_creacion:   time.Now(),
				Completada:       false,
				Fecha_completada: nil,
				Id_usuario:       1,
			},
			errorEsperado: false,
		},
	}

	db := conectarBaseDeDatosTest()
	tm := NewTareasManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			err := tm.RegistrarTarea(ctx, test.tareaTest)
			if (err != nil) && !test.errorEsperado {

				t.Errorf("Error no esperado. Detalle: %v\n\n", err)
			}
		})
	}

}
