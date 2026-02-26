package models

import (
	"time"
)

type Tarea struct {
	Id               *int
	Titulo           string
	Fecha_creacion   time.Time
	Completada       bool
	Fecha_completada *time.Time
	Id_usuario       int
}

type Usuario struct {
	Id            *int
	Email         string
	Password_hash string
}
