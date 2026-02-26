package db

import "errors"

var (
	ErrUsuarioExiste           = errors.New("El usuario ya existe")
	ErrUsuarioNoExiste         = errors.New("El usuario buscado no existe")
	ErrPasswordIncorrecto      = errors.New("La contraseña no es correcta")
	ErrUsuarioAsignadoNoexiste = errors.New("El usuario asignado a la tarea no existe")
	ErrTareaNoExiste           = errors.New("La tarea buscada no existe")
)
