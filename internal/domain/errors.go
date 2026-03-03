package domain

import "errors"

var (
	// Errores de repositories/db
	ErrEmailDuplicado          = errors.New("El email ya se encuentra registrado")
	ErrEmailNoEncontrado       = errors.New("El email pasado no fue encontrado")
	ErrIdNoEncontrado          = errors.New("El id pasado no fue encontrado")
	ErrPasswordIncorrecto      = errors.New("La contraseña no es correcta")
	ErrUsuarioAsignadoNoexiste = errors.New("El usuario asignado a la tarea no existe")
	ErrTareaNoExiste           = errors.New("La tarea buscada no existe")

	// Errores de service
	ErrEmailRequerido    = errors.New("El campo email no puede estar en blanco")
	ErrPasswordRequerido = errors.New("El campo password no puede estar en blanco")
	ErrPasswordCorto     = errors.New("El password es demasiado corto")
	ErrIdNovalido        = errors.New("El id pasado debe ser mayor a 0")
	ErrUsuarioExiste     = errors.New("El usuario ya existe")
	ErrUsuarioNoExiste   = errors.New("El usuario no existe")
)
