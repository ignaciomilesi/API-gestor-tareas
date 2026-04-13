# API Gestor de Tareas

API REST desarrollada en Go para la gestión de usuarios y tareas, con autenticación basada en JWT y arquitectura en capas.

Proyecto backend desarrollado como práctica avanzada de API REST en Go.

## Features

* Registro y autenticación de usuarios (JWT)
* CRUD completo de tareas
* Búsqueda y filtrado
* Arquitectura limpia por capas


## Tecnologías

* **Go (Golang)**
* **Gin** — framework HTTP
* **PostgreSQL** — base de datos
* **pgx** — driver de PostgreSQL
* **JWT** — autenticación
* **bcrypt** — hash de contraseñas


## Arquitectura

El proyecto sigue una arquitectura en capas:

```
handlers → services → repository → database
```

### Capas

* **handlers** → Manejan HTTP (request/response)
* **services** → Lógica de negocio
* **repository** → Acceso a datos
* **middleware** → Autenticación y concerns transversales



## Flujo de autenticación

1. Usuario se registra → `/signup`
2. Hace login → `/login`
3. Recibe un token JWT
4. Usa el token en requests protegidas:

```http
Authorization: Bearer <token>
```

El middleware valida el token y extrae el `user_id`.



## Estructura del proyecto

```text
api/
  └── main.go              # Entry point

internal/
  ├── handlers/           # HTTP (Gin)
  ├── services/           # Lógica de negocio
  ├── repositories/
  │     └── db/           # Acceso a datos
  ├── models/             # Entidades del dominio
  └── middleware/         # Auth, logs, etc.

config/                   # Configuración
```


## Base de datos

El script para la creación de las tablas se encuentra en:

```text
tablas_db.txt
```

### Resumen de tablas

#### usuarios

Almacena la información de los usuarios del sistema.

* `id` → identificador único (PK)
* `email` → correo electrónico (único)
* `password_hash` → contraseña hasheada

---

#### tareas

Almacena las tareas asociadas a cada usuario.

* `id` → identificador de la tarea
* `titulo` → descripción de la tarea
* `fecha_creacion` → fecha de creación
* `completada` → estado de la tarea (true/false)
* `fecha_completada` → fecha en que se completó (nullable)
* `id_usuario` → referencia al usuario (FK → usuarios.id)



## Endpoints

### Públicos

#### POST `/signup` : Crear usuario

```json
{
  "email": "test@test.com",
  "password": "123456"
}
```

---

#### POST `/login` : Login de usuario

```json
{
  "email": "test@test.com",
  "password": "123456"
}
```

**Response**

```json
{
  "token": "jwt_token"
}
```

---

### Privados (`/api`)

#### POST `/api/tareas` - Crear tarea

```json
{
  "descripcion": "Mi tarea",
  "fecha": "10/04/2026"
}
```

---

#### GET `/api/tareas` - Listar tareas

Query params:

```
?completadas=true
```

---

#### PUT `/api/tareas` - Actualizar descripción

```json
{
  "id_tarea": 1,
  "nueva_descripcion": "Nueva descripción"
}
```

---

#### POST `/api/tareas/finalizar` - Marcar tarea como completada

```json
{
  "id_tarea": 1,
  "fecha": "10/04/2026"
}
```

---

#### GET `/api/tareas/buscar` - Buscar tareas

Query params:

```
?parametro_busqueda=texto
```

---

#### PUT `/api/user/password` - Actualizar contraseña

```json
{
  "password": "NuevoPassword"
}
```


## Testing

### Tipos de tests

* **Unitarios**

  * hanndlers
  * services
  * repository

* **Integración / E2E**

  * Signup , Login (obtención de JWT),  Crear tarea autenticado



