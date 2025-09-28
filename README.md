# ReminderApp Backend

Backend en **Go** con MongoDB para la aplicación ReminderApp.

## 🚀 Tecnologías
- Go
- Gin (framework HTTP)
- MongoDB
- JWT para autenticación

## 📂 Estructura
- `cmd/server` → entrypoint (main.go)
- `internal/config` → configuración y conexión DB
- `internal/domain` → entidades de negocio (User, Medication, Reminder)
- `internal/repository` → capa de acceso a Mongo
- `internal/service` → lógica de negocio
- `internal/handler` → controladores HTTP
- `api/routes.go` → definición de rutas
- `pkg` → librerías reutilizables (logger, middleware, etc.)

## 🔧 Uso
```bash
go run cmd/server/main.go
