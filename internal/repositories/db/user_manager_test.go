package db

import (
	"errors"
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
		errorEsperado error
	}{
		{
			name: "Email duplicado",
			usuarioTest: models.Usuario{
				Email:         "prueba@prueba.com",
				Password_hash: "4567",
			},
			errorEsperado: ErrUsuarioExiste,
		},
	}

	db := conectarBaseDeDatosTest()
	um := NewUserManager(db)

	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			_, err := um.GenerarNuevoUsuario(ctx, test.usuarioTest)

			if err == nil {
				t.Errorf("Error. Se esperaba que falle")
			}

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado. Detalle: %v\n\n", err)
			}
		})
	}

}

func TestObternerId(t *testing.T) {

	tests := []struct {
		name          string
		usuarioTest   models.Usuario
		errorEsperado error
	}{
		{
			name: "mail correcto, password correcto",
			usuarioTest: models.Usuario{
				Email:         "prueba@prueba.com",
				Password_hash: "4567",
			},
			errorEsperado: nil,
		},
		{
			name: "mail incorrecto, password correcto",
			usuarioTest: models.Usuario{
				Email:         "mailPruebaIncorrecto@prueba.com",
				Password_hash: "4567",
			},
			errorEsperado: ErrUsuarioNoExiste,
		},
		{
			name: "mail correcto, password incorrecto",
			usuarioTest: models.Usuario{
				Email:         "prueba@prueba.com",
				Password_hash: "password incorrecto",
			},
			errorEsperado: ErrPasswordIncorrecto,
		},
	}

	db := conectarBaseDeDatosTest()
	um := NewUserManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			_, err := um.ObternerId(ctx, test.usuarioTest)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf(" -- Error no esperado: se esperaba %t pero se obtuvo %t\n\n",
					test.errorEsperado, err)
			}
		})
	}

}

func TestModifcarContraseña(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		errorEsperado error
	}{
		{
			name:          "usuario incorrecto",
			id:            999999,
			errorEsperado: ErrUsuarioNoExiste,
		},
		{
			name:          "actualización de password",
			id:            3,
			errorEsperado: nil,
		},
	}

	db := conectarBaseDeDatosTest()
	um := NewUserManager(db)
	ctx := t.Context()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			err := um.ModifcarContraseña(ctx, test.id, "45678")

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf(" -- Error no esperado: se esperaba %t pero se obtuvo %t\n\n",
					test.errorEsperado, err)
			}

		})
	}

}
