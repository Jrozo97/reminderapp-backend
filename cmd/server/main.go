package main

import (
	"log"
	"os"

	"github.com/Jrozo97/reminderapp-backend/api"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar .env, usando variables del sistema")
	}

	// Crear router
	r := gin.Default()

	// Registrar rutas
	api.RegisterRoutes(r)

	// Puerto desde .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor corriendo en puerto:", port)
	r.Run(":" + port)
}
