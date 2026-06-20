# 🔦 Faro

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3.4-42B883?style=flat-square&logo=vuedotjs&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat-square&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat-square&logo=docker&logoColor=white)
![JWT](https://img.shields.io/badge/Auth-JWT-black?style=flat-square&logo=jsonwebtokens)

Gestión de tareas full-stack. Backend Go con arquitectura en capas, JWT y GORM. Frontend Vue 3 con componentes reutilizables. Levanta con un comando.

---

## 🚀 Inicio rápido

```bash
docker-compose up --build
```

| Servicio | URL |
|---|---|
| 🌐 Frontend | http://localhost:3000 |
| ⚙️ API | http://localhost:8080/api/v1 |
| 🗄️ PostgreSQL | localhost:**5433** |

> Rebuild con cambios: `docker-compose down && docker-compose build --no-cache && docker-compose up -d`

---

## 👥 Usuarios seed

Se crean automáticamente al primer arranque. Contraseña de todos: **`123456789`**

| | Email |
|---|---|
| 🟢 | admin@faro.app |
| 🔵 | juan.perez@faro.app |
| 🟣 | maria.garcia@faro.app |
| 🔴 | carlos.lopez@faro.app |
| 🟡 | ana.martinez@faro.app |

Categorías seed: `Trabajo` `Personal` `Urgente` `Aprendizaje` `Ideas`

---

## 📡 API

Base URL: `http://localhost:8080/api/v1`

### 🔓 Auth — público

| Método | Ruta | |
|---|---|---|
| `POST` | `/auth/register` | Crear cuenta · devuelve JWT |
| `POST` | `/auth/login` | Iniciar sesión · devuelve JWT |

### ✅ Tasks — `JWT`

| Método | Ruta | |
|---|---|---|
| `GET` | `/tasks` | Listar |
| `POST` | `/tasks` | Crear |
| `GET` | `/tasks/{id}` | Obtener |
| `PUT` | `/tasks/{id}` | Actualizar / completar |
| `DELETE` | `/tasks/{id}` | Eliminar |

### 📝 Notes — `JWT` · privadas por usuario

| Método | Ruta | |
|---|---|---|
| `GET` | `/notes` | Mis notas |
| `POST` | `/notes` | Crear |
| `GET` | `/notes/{id}` | Obtener |
| `PUT` | `/notes/{id}` | Actualizar |
| `DELETE` | `/notes/{id}` | Eliminar |

### 🏷️ Categories — `JWT`

| Método | Ruta | |
|---|---|---|
| `GET` | `/categories` | Listar |
| `POST` | `/categories` | Crear · `name` + `color` hex |
| `PUT` | `/categories/{id}` | Actualizar |
| `DELETE` | `/categories/{id}` | Eliminar |

---

## 🏗️ Arquitectura

```
Request → Router → Auth Middleware → Handler → Service → Repository → PostgreSQL
```

| Capa | Responsabilidad |
|---|---|
| **Handler** | Parsea request, escribe respuesta JSON |
| **Service** | Validaciones de negocio |
| **Repository** | Queries GORM — interfaz desacoplada |
| **Middleware** | Valida JWT, inyecta `userID` en contexto |

<details>
<summary>📁 Estructura de carpetas</summary>

```
proyecto-go/
├── backend/
│   ├── cmd/api/main.go          # Wiring principal
│   ├── internal/
│   │   ├── domain/              # Task · Note · User · Category
│   │   ├── handler/             # HTTP controllers
│   │   ├── middleware/          # JWT · CORS · Logger
│   │   ├── repository/          # GORM implementations
│   │   ├── seed/                # Datos iniciales
│   │   └── service/             # Lógica de negocio
│   └── Dockerfile               # Multi-stage → scratch
│
├── frontend/
│   ├── src/
│   │   ├── components/common/   # BaseButton · BaseModal · ...
│   │   ├── stores/              # Pinia auth
│   │   └── views/               # Tasks · Notes · Categories · Auth
│   └── Dockerfile               # node:20 → nginx:alpine
│
├── docker-compose.yml
└── check-ports.ps1              # Valida puertos antes de levantar
```

</details>

---

## 🧩 Componentes UI

`BaseButton` · `BaseInput` · `BaseTextarea` · `BaseAlert` · `BaseModal` · `BaseBadge` · `EmptyState` · `PageHeader`
