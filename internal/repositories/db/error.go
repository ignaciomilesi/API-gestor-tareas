package db

import "errors"

var (
	ErrUsuarioExiste      = errors.New("El usuario ya existe")
	ErrUsuarioNoExiste    = errors.New("El usuario no existe")
	ErrPasswordIncorrecto = errors.New("La contraseña no es correcta")
)
