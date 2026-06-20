# Faro

> Gestión de tareas full-stack — Go + Vue 3 + Docker

Backend con arquitectura en capas, JWT, GORM y PostgreSQL. Frontend en Vue 3 con componentes reutilizables y Pinia. Todo el stack levanta con un solo comando.

---

## Stack

| Capa | Tecnología |
|---|---|
| Backend | Go 1.22, chi v5, GORM v1.25, JWT HS256, bcrypt |
| Base de datos | PostgreSQL 16 |
| Frontend | Vue 3.4, Vite 5, Pinia, Axios, vue-router |
| Infraestructura | Docker Compose v2, Nginx |

---

## Inicio rápido

```bash
# 1. Verificar que los puertos estén disponibles (requiere PowerShell)
.\check-ports.ps1

# 2. Levantar todo el stack
docker-compose up --build

# 3. Abrir la app
#    Frontend  → http://localhost:3000
#    API REST  → http://localhost:8080/api/v1
#    PostgreSQL → localhost:5433
```

La primera vez tarda ~2 min mientras descarga imágenes y compila. El seed inserta usuarios y categorías automáticamente al arrancar.

> Para forzar rebuild con cambios de código:
> ```bash
> docker-compose down && docker-compose build --no-cache && docker-compose up -d
> ```

---

## Usuarios seed

Al arrancar por primera vez se crean automáticamente:

| Email | Contraseña |
|---|---|
| admin@faro.app | `123456789` |
| juan.perez@faro.app | `123456789` |
| maria.garcia@faro.app | `123456789` |
| carlos.lopez@faro.app | `123456789` |
| ana.martinez@faro.app | `123456789` |

También se crean 5 categorías seed: Trabajo, Personal, Urgente, Aprendizaje, Ideas.

---

## API REST

Base URL: `http://localhost:8080/api/v1`

### Auth (público)

| Método | Ruta | Descripción |
|---|---|---|
| `POST` | `/auth/register` | Crear cuenta — devuelve JWT |
| `POST` | `/auth/login` | Iniciar sesión — devuelve JWT |

### Tasks (JWT requerido)

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/tasks` | Listar tareas |
| `POST` | `/tasks` | Crear tarea |
| `GET` | `/tasks/{id}` | Obtener por ID |
| `PUT` | `/tasks/{id}` | Actualizar (incluye toggle `completed`) |
| `DELETE` | `/tasks/{id}` | Eliminar |

### Notes (JWT — privadas por usuario)

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/notes` | Notas del usuario autenticado |
| `POST` | `/notes` | Crear nota |
| `GET` | `/notes/{id}` | Obtener (solo si pertenece al usuario) |
| `PUT` | `/notes/{id}` | Actualizar |
| `DELETE` | `/notes/{id}` | Eliminar |

### Categories (JWT requerido)

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/categories` | Listar categorías |
| `POST` | `/categories` | Crear — `name` + `color` hex |
| `PUT` | `/categories/{id}` | Actualizar nombre o color |
| `DELETE` | `/categories/{id}` | Eliminar |

---

## Ejemplos con curl

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@faro.app","password":"123456789"}'

# Crear tarea (con token)
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"title":"Estudiar Go","description":"Terminar el tour de Go"}'

# Crear nota privada
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"title":"Ideas","content":"Agregar paginación a la API"}'

# Crear categoría con color
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"name":"Urgente","color":"#EF4444"}'
```

---

## Estructura del proyecto

```
proyecto-go/
├── backend/
│   ├── cmd/api/main.go              # Wiring: repos → services → handlers
│   ├── internal/
│   │   ├── config/                  # Variables de entorno y DSN
│   │   ├── domain/                  # Task · Note · User · Category + DTOs
│   │   ├── handler/                 # Controladores HTTP (chi)
│   │   ├── middleware/              # Auth JWT · CORS · Logger
│   │   ├── repository/              # Interfaces + implementaciones GORM
│   │   ├── seed/                    # Datos iniciales al primer arranque
│   │   └── service/                 # Lógica de negocio + validaciones
│   ├── pkg/
│   │   ├── database/                # Conexión GORM + pool
│   │   ├── response/                # Helpers JSON
│   │   └── token/                   # Generación y validación JWT
│   └── Dockerfile                   # Multi-stage: golang:1.22 → scratch
│
├── frontend/
│   ├── src/
│   │   ├── api/                     # Axios: tasks · notes · categories · auth
│   │   ├── components/
│   │   │   ├── common/              # BaseButton · BaseInput · BaseModal · ...
│   │   │   └── layout/              # AppNavbar
│   │   ├── stores/                  # Pinia: auth store con localStorage
│   │   └── views/                   # Tasks · Notes · Categories · Login · Register
│   ├── nginx.conf                   # SPA routing + proxy /api/ → backend
│   └── Dockerfile                   # node:20 build → nginx:alpine
│
├── docker-compose.yml               # 3 servicios · red interna · puerto 5433
└── check-ports.ps1                  # Valida :5433 · :8080 · :3000
```

---

## Arquitectura del backend

```
Request → chi Router → Middleware (Auth JWT) → Handler → Service → Repository → GORM → PostgreSQL
```

- **Handler**: parsea request, llama al servicio, escribe respuesta JSON
- **Service**: validaciones de negocio, orquesta el repositorio
- **Repository**: interfaz + implementación GORM — facilita testing
- **Middleware Auth**: valida JWT, inyecta `userID` en el contexto

Las notas son privadas por usuario: el repositorio filtra por `user_id` en todas las queries.

---

## Componentes Vue reutilizables

| Componente | Descripción |
|---|---|
| `BaseButton` | Variantes: primary, secondary, ghost, danger · Tamaños: sm, md, lg · Loading spinner |
| `BaseInput` | Label + error + hint + v-model |
| `BaseTextarea` | Igual que BaseInput para áreas de texto |
| `BaseAlert` | Tipos: error, success, warning |
| `BaseModal` | Teleport + Transition + cierre con Escape · Scroll lock |
| `BaseBadge` | Color hex dinámico |
| `EmptyState` | Estado vacío con icono, título y slot de acción |
| `PageHeader` | Título + subtítulo + slot de acciones |

---

## Puertos

| Servicio | Interno (Docker) | Externo (host) |
|---|---|---|
| PostgreSQL | 5432 | **5433** |
| Backend Go | 8080 | 8080 |
| Frontend Nginx | 80 | 3000 |

El puerto 5433 evita conflictos con instalaciones locales de PostgreSQL.
