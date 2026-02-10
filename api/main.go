package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default() // creo Gin con el middleware default

	// devolvemos un json
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// encendemos el server, por defecto escucha el localhost:8080 on Windows
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
