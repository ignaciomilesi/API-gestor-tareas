# API-gestor-tareas
API de gestión de tareas (Task Manager). Funciones: sign up y Log in.  Crear, listar, actualizar y borrar tareas.

Estructura de carpetas
```text
- api/
  - main.go        → arranque de la app

- internal/
  - handlers/      → HTTP / Gin (requests, responses, status codes)
  - services/      → lógica de negocio y casos de uso
  - repositories/
    - db/ → acceso a datos (SQL, DB)
  - models/        → structs del dominio (entidades)
  - middleware/    → middlewares (auth, logs, permisos)

- config/          → configuración y variables de entorno
docker/          → imágenes y config Docker
```
