package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/Jrozo97/reminderapp-backend/api"
	"github.com/Jrozo97/reminderapp-backend/internal/config"
)

// Helper para inicializar el servidor en modo test
func setupRouter() (*gin.Engine, error) {

	err := godotenv.Load("../.env")
	if err != nil {
		return nil, fmt.Errorf("no se pudo cargar .env: %w", err)
	}

	// forzar modo test de gin
	gin.SetMode(gin.TestMode)

	config.ConnectMongo()

	r := gin.Default()
	api.RegisterRoutes(r)
	return r, nil
}

func TestRegisterAndLogin(t *testing.T) {
	router, err := setupRouter()
	if err != nil {
		t.Fatal(err)
	}

	registerPayload := map[string]string{
		"name":     "Test User",
		"email":    "integration@example.com",
		"password": "password123",
	}
	body, err := json.Marshal(registerPayload)
	require.NoError(t, err)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusBadRequest {
		// Created la primera vez, BadRequest si ya existía
		t.Fatalf("Registro falló. Status: %d, Body: %s", w.Code, w.Body.String())
	}

	loginPayload := map[string]string{
		"email":    "integration@example.com",
		"password": "password123",
	}

	body, errLogin := json.Marshal(loginPayload)
	require.NoError(t, errLogin)

	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login falló. Status: %d, Body: %s", w.Code, w.Body.String())
	}

	// parsear respuesta
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	token := resp["token"]
	if token == "" {
		t.Fatal("No se devolvió token en login")
	}
	t.Log("✅ Token recibido:", token[:20], "...")
}
