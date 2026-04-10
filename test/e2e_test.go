package test

import (
	"api-gestor-tareas/internal/handlers"
	"api-gestor-tareas/internal/middleware"
	"api-gestor-tareas/internal/repositories/dataBase"
	"api-gestor-tareas/internal/services"
	"fmt"
	"log"
	"os"

	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func TestEndToEndLoginYCrearTarea(t *testing.T) {

	gin.SetMode(gin.TestMode)

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	dburl := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable&search_path=test",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"))

	db, err := dataBase.NewGestorDb(ctx, dburl)
	if err != nil {
		log.Fatal("Error conectar base de datos", err)
	}

	userManager := dataBase.NewUserManager(db)
	tareaManager := dataBase.NewTareasManager(db)

	userService := services.NewUserService(userManager)
	tareaService := services.NewTareaService(tareaManager)

	jwtSecret := "test_secret"

	userHandler := handlers.NewUserHandler(userService, jwtSecret)
	tareaHandler := handlers.NewTareaHandler(tareaService)

	// router
	r := gin.New()

	r.POST("/signup", userHandler.Signup)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	auth.POST("/tareas", tareaHandler.Nueva)

	// SIGNUP

	signupBody := `{
		"email": "test@test.com",
		"password": "123456"
	}`

	req, _ := http.NewRequest("POST", "/signup", strings.NewReader(signupBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 && w.Code != 201 {
		t.Fatalf("signup fallo: %d - %s", w.Code, w.Body.String())
	}

	// LOGIN

	loginBody := `{
		"email": "test@test.com",
		"password": "123456"
	}`

	req, _ = http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("login fallo: %d - %s", w.Code, w.Body.String())
	}

	var loginResp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	if err != nil {
		t.Fatalf("error parseando login: %v", err)
	}

	token := loginResp["token"]
	if token == "" {
		t.Fatal("no vino token")
	}

	// CREAR TAREA

	tareaBody := `{
		"descripcion": "mi tarea test",
		"fecha": "10/04/2026"
	}`

	req, _ = http.NewRequest("POST", "/api/tareas", strings.NewReader(tareaBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 && w.Code != 201 {
		t.Fatalf("crear tarea fallo: %d - %s", w.Code, w.Body.String())
	}

	var tareaResp map[string]int
	err = json.Unmarshal(w.Body.Bytes(), &tareaResp)
	if err != nil {
		t.Fatalf("error parseando tarea: %v", err)
	}

	if tareaResp["id"] == 0 {
		t.Fatal("id de tarea invalido")
	}
}
