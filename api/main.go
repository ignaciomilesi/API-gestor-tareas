package main

import (
	"api-gestor-tareas/internal/handlers"
	"api-gestor-tareas/internal/middleware"
	"api-gestor-tareas/internal/repositories/dataBase"
	"api-gestor-tareas/internal/services"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Cargamos el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// generamos el contexto
	ctx := context.Background()

	// generamos los managerDb
	dburl := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"))

	db, err := dataBase.NewGestorDb(ctx, dburl)
	if err != nil {
		log.Fatal("Error conectar base de datos", err)
	}

	userManagerDb := dataBase.NewUserManager(db)
	tareaManagerDb := dataBase.NewTareasManager(db)

	// generamos los services
	userService := services.NewUserService(userManagerDb)
	tareaService := services.NewTareaService(tareaManagerDb)

	//generamos los handlers
	userHandler := handlers.NewUserHandler(userService, os.Getenv("JWT_SECRET"))
	tareasHandler := handlers.NewTareaHandler(tareaService)

	// creo el ruter
	r := gin.Default()

	// públicas (sin middleware)
	r.POST("/login", userHandler.Login)
	r.POST("/signup", userHandler.Signup)

	// protegidas
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET"))) // aplico el middleware

	auth.POST("/tareas", tareasHandler.Nueva)
	auth.GET("/tareas", tareasHandler.Listar)
	auth.PUT("/tareas", tareasHandler.Modificar)
	auth.POST("/tareas/finalizar", tareasHandler.Finalizar)
	auth.GET("/tareas/buscar", tareasHandler.Buscar)

	// encendemos el server, por defecto escucha el localhost:8080 on Windows
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
