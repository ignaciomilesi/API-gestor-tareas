package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"api-gestor-tareas/internal/models"
)

func conectarBaseDeDatosTest() *pgxpool.Pool {
	if err := godotenv.Overload("../../../.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	dburl := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable&search_path=test",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"))

	db, err := NewGestorDb(dburl)
	if err != nil {
		log.Fatal("Error conectar base de datos", err)
	}
	return db
}
func TestGenerarNuevoUsuario(t *testing.T) {

	tests := []struct {
		name          string
		usuarioTest   models.Usuario
		errorEsperado bool
	}{
		{
			name: "Email duplicado",
			usuarioTest: models.Usuario{
				Email:         "prueba@prueba.com",
				Password_hash: "4567",
			},
			errorEsperado: true,
		},
	}

	db := conectarBaseDeDatosTest()
	um := NewUserManager(db)

	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			err := um.GenerarNuevoUsuario(ctx, test.usuarioTest)
			if (err != nil) && !test.errorEsperado {

				t.Errorf("Error no esperado. Detalle: %v\n\n", err)
			}
		})
	}

}

func TestComprobarExisteUsuario(t *testing.T) {

	tests := []struct {
		name              string
		usuarioTest       models.Usuario
		resultadoEsperado bool
	}{
		{
			name: "mail correcto, password correcto",
			usuarioTest: models.Usuario{
				Email:         "prueba@prueba.com",
				Password_hash: "4567",
			},
			resultadoEsperado: true,
		},
		{
			name: "mail incorrecto, password correcto",
			usuarioTest: models.Usuario{
				Email:         "mailPruebaIncorrecto@prueba.com",
				Password_hash: "4567",
			},
			resultadoEsperado: false,
		},
		{
			name: "mail correcto, password incorrecto",
			usuarioTest: models.Usuario{
				Email:         "prueba@prueba.com",
				Password_hash: "password incorrecto",
			},
			resultadoEsperado: false,
		},
	}

	db := conectarBaseDeDatosTest()
	um := NewUserManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			resultado := um.ComprobarExisteUsuario(ctx, test.usuarioTest)
			if resultado != test.resultadoEsperado {

				t.Errorf(" -- Error no esperado: se esperaba %t pero se obtuvo %t\n\n",
					test.resultadoEsperado, resultado)
			}
		})
	}

}

func TestModifcarContraseña(t *testing.T) {
	tests := []struct {
		name          string
		usuarioTest   models.Usuario
		errorEsperado bool
	}{
		{
			name: "usuario incorrecto",
			usuarioTest: models.Usuario{
				Email:         "emailIncorecto@prueba.com",
				Password_hash: "4567",
			},
			errorEsperado: true,
		},
		{
			name: "actualización de password",
			usuarioTest: models.Usuario{
				Email:         "emailPruebaModif@prueba.com",
				Password_hash: "asdf",
			},
			errorEsperado: false,
		},
	}

	db := conectarBaseDeDatosTest()
	um := NewUserManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			err := um.ModifcarContraseña(ctx, test.usuarioTest, "4567")

			if (err != nil) && !test.errorEsperado {

				t.Errorf("Error no esperado. Detalle: %v\n\n", err)
			}

		})
	}

}
