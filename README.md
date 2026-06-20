# proyecto-go

API REST de gestión de tareas en Go con arquitectura en capas, GORM + PostgreSQL, chi router y pruebas unitarias completas. Diseñada para aprender cómo se estructura un backend real.

## Requisitos

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) (para levantar PostgreSQL)

## Inicio rápido

```bash
# 1. Instalar dependencias
go mod tidy

# 2. Copiar variables de entorno
cp .env.example .env

# 3. Levantar PostgreSQL
docker-compose up -d

# 4. Correr el servidor (las migraciones se aplican automáticamente)
go run ./cmd/api
```

El servidor arranca en `http://localhost:8080`.

## Comandos útiles

```bash
make run              # Correr el servidor
make test             # Pruebas unitarias (sin DB)
make test-cover       # Pruebas con cobertura
make test-integration # Pruebas de integración (requiere DB)
make docker-up        # Levantar PostgreSQL
make docker-down      # Detener PostgreSQL
make build            # Compilar binario en bin/api
make tidy             # go mod tidy
```

## Endpoints

| Método | Ruta                   | Descripción             | Status    |
|--------|------------------------|-------------------------|-----------|
| GET    | /api/v1/tasks/         | Listar todas las tareas | 200       |
| POST   | /api/v1/tasks/         | Crear tarea             | 201       |
| GET    | /api/v1/tasks/{id}     | Obtener tarea por ID    | 200 / 404 |
| PUT    | /api/v1/tasks/{id}     | Actualizar tarea        | 200 / 404 |
| DELETE | /api/v1/tasks/{id}     | Eliminar tarea          | 204 / 404 |

## Ejemplos con curl

**Crear una tarea:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks/ \
  -H "Content-Type: application/json" \
  -d '{"title": "Aprender Go", "description": "Estudiar el tour de Go"}'
```

**Listar tareas:**
```bash
curl http://localhost:8080/api/v1/tasks/
```

**Obtener por ID:**
```bash
curl http://localhost:8080/api/v1/tasks/1
```

**Actualizar:**
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Aprender Go", "description": "Completado!", "completed": true}'
```

**Eliminar:**
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1
```

## Estructura del proyecto

```
proyecto-go/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada y wiring de dependencias
├── internal/
│   ├── config/
│   │   └── config.go            # Configuración desde variables de entorno
│   ├── domain/
│   │   └── task.go              # Entidad Task y DTOs de request
│   ├── repository/
│   │   ├── task_repository.go   # Interfaz + implementación GORM
│   │   └── task_repository_test.go  # Tests de integración (requieren DB)
│   ├── service/
│   │   ├── task_service.go      # Lógica de negocio
│   │   └── task_service_test.go # Tests unitarios con mock del repositorio
│   ├── handler/
│   │   ├── task_handler.go      # Controladores HTTP (chi)
│   │   └── task_handler_test.go # Tests unitarios con mock del servicio
│   └── middleware/
│       └── logger.go            # Middleware de logging HTTP
├── pkg/
│   ├── database/
│   │   └── postgres.go          # Conexión GORM + pool de conexiones
│   └── response/
│       └── response.go          # Helpers JSON para respuestas HTTP
├── migrations/
│   └── 001_create_tasks.sql     # Esquema SQL documentado
├── .env.example                 # Plantilla de variables de entorno
├── docker-compose.yml           # PostgreSQL local
├── Makefile                     # Comandos del proyecto
├── go.mod
├── ARCHITECTURE.md              # Explicación de la arquitectura
└── README.md
```

Consulta [ARCHITECTURE.md](ARCHITECTURE.md) para entender el patrón de capas en detalle.
