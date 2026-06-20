# Arquitectura del proyecto

## Stack tecnologico

| Tecnologia | Version | Para que se usa |
|------------|---------|-----------------|
| **Go** | 1.22+ | Lenguaje principal. Compilado, tipado estaticamente, excelente para backends |
| **PostgreSQL** | 16 | Base de datos relacional donde se persisten los datos |
| **GORM** | v1.25 | ORM de Go. Traduce structs de Go a tablas SQL y viceversa |
| **chi** | v5.1 | Router HTTP ligero y rapido. Define las rutas de la API REST |
| **godotenv** | v1.5 | Carga variables de entorno desde el archivo `.env` |
| **Docker Compose** | - | Levanta PostgreSQL localmente sin instalarlo directamente |

---

## Estructura de carpetas

```
proyecto-go/
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── domain/
│   ├── repository/
│   ├── service/
│   ├── handler/
│   └── middleware/
│
├── pkg/
│   ├── database/
│   └── response/
│
├── migrations/
├── .env.example
├── docker-compose.yml
└── Makefile
```

---

## Que hace cada carpeta

### `cmd/api/main.go`
El **punto de entrada** del programa. Es el unico archivo con `func main()`.

Su trabajo es conectar todas las piezas:
1. Carga la configuracion
2. Conecta a la base de datos
3. Crea el repositorio, el servicio y el handler
4. Registra las rutas
5. Arranca el servidor HTTP

```go
// Ejemplo simplificado de lo que hace main.go
cfg  := config.Load()
db   := database.NewPostgresDB(cfg.DSN())
repo := repository.NewGormTaskRepository(db)
svc  := service.NewTaskService(repo)
h    := handler.NewTaskHandler(svc)
// ... registrar rutas y arrancar servidor
```

> Regla: `main.go` no tiene logica de negocio. Solo ensambla piezas.

---

### `internal/config/`
Lee las **variables de entorno** (o el archivo `.env`) y las expone como una struct tipada.

```go
type Config struct {
    DBHost     string  // DB_HOST
    DBPort     string  // DB_PORT
    DBUser     string  // DB_USER
    DBPassword string  // DB_PASSWORD
    DBName     string  // DB_NAME
    ServerPort string  // SERVER_PORT
}
```

Ventaja: ningun otro archivo del proyecto llama a `os.Getenv()` directamente. Todo pasa por aqui. Si cambia el nombre de una variable de entorno, solo se toca este archivo.

---

### `internal/domain/`
Contiene las **estructuras de datos** que representan las entidades del negocio.

```go
type Task struct {
    ID          uint      // identificador unico (generado por la DB)
    Title       string    // titulo de la tarea
    Description string    // descripcion opcional
    Completed   bool      // estado: pendiente o completada
    CreatedAt   time.Time // timestamp automatico de GORM
    UpdatedAt   time.Time // timestamp automatico de GORM
}
```

Tambien define los **DTOs** (Data Transfer Objects): los structs que representan lo que llega en el body JSON de cada peticion HTTP.

```go
type CreateTaskRequest struct { Title, Description string }
type UpdateTaskRequest struct { Title, Description string; Completed bool }
```

> Regla: `domain/` no importa ningun otro paquete interno. Es la base de todo.

---

### `internal/repository/`
Responsable de **guardar y recuperar datos** de la base de datos.

Define la **interfaz** que describe que operaciones existen:

```go
type TaskRepository interface {
    FindAll() ([]domain.Task, error)
    FindByID(id uint) (domain.Task, error)
    Create(task *domain.Task) error
    Update(task *domain.Task) error
    Delete(id uint) error
}
```

Y la **implementacion con GORM** que traduce esas operaciones a SQL:

```go
// GORM convierte esto en: SELECT * FROM tasks WHERE id = ?
func (r *gormTaskRepository) FindByID(id uint) (domain.Task, error) {
    var task domain.Task
    r.db.First(&task, id)
    return task, nil
}
```

> Por que una interfaz? Porque el service no sabe ni le importa si los datos vienen de PostgreSQL, MongoDB o un archivo. Solo conoce la interfaz.

---

### `internal/service/`
Contiene la **logica de negocio** de la aplicacion.

Aqui van las reglas que no son ni HTTP ni base de datos. Por ejemplo:
- "No se puede crear una tarea sin titulo"
- "Al actualizar, primero verificar que la tarea exista"

```go
func (s *taskService) Create(req domain.CreateTaskRequest) (domain.Task, error) {
    if req.Title == "" {
        return domain.Task{}, ErrInvalidTitle  // regla de negocio
    }
    task := domain.Task{Title: req.Title, Description: req.Description}
    s.repo.Create(&task)  // delega el guardado al repositorio
    return task, nil
}
```

Tambien define su propia interfaz `TaskService` que el handler usa.

> Regla: el service nunca importa `net/http`. No sabe nada de HTTP.

---

### `internal/handler/`
Los **controladores HTTP**. Traducen peticiones HTTP en llamadas al service.

Su trabajo en cada request:
1. Leer parametros de la URL o el body JSON
2. Llamar al metodo correspondiente del service
3. Traducir el resultado (o error) a un codigo HTTP y respuesta JSON

```go
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req domain.CreateTaskRequest
    json.NewDecoder(r.Body).Decode(&req)  // 1. leer el body JSON

    task, err := h.svc.Create(req)        // 2. llamar al service
    if err != nil {
        response.Error(w, 400, err.Error()) // 3a. error → 400 Bad Request
        return
    }
    response.JSON(w, 201, task)            // 3b. exito → 201 Created + JSON
}
```

Las rutas se registran en chi:

```
GET    /api/v1/tasks/       → GetAll
POST   /api/v1/tasks/       → Create
GET    /api/v1/tasks/{id}   → GetByID
PUT    /api/v1/tasks/{id}   → Update
DELETE /api/v1/tasks/{id}   → Delete
```

> Regla: el handler nunca accede a la base de datos directamente.

---

### `internal/middleware/`
Funciones que se ejecutan **antes o despues de cada peticion HTTP**, sin que los handlers lo sepan.

El middleware `Logger` registra en consola cada peticion con su metodo, ruta, status code y tiempo de respuesta:

```
POST   /api/v1/tasks/       201 1.2ms
GET    /api/v1/tasks/1      200 0.8ms
DELETE /api/v1/tasks/99     404 0.5ms
```

Para capturar el status code, el middleware envuelve el `http.ResponseWriter` original:

```go
type wrappedWriter struct {
    http.ResponseWriter
    statusCode int  // aqui guarda el status que el handler escribio
}
```

---

### `pkg/database/`
Crea la **conexion a PostgreSQL** usando GORM y configura el **pool de conexiones**.

```go
sqlDB.SetMaxOpenConns(25)          // maximo 25 conexiones abiertas
sqlDB.SetMaxIdleConns(5)           // maximo 5 conexiones inactivas esperando
sqlDB.SetConnMaxLifetime(5 * Minute) // reusar conexion por hasta 5 minutos
```

El pool evita abrir una conexion nueva por cada peticion (caro) y reutiliza las existentes.

---

### `pkg/response/`
Funciones de ayuda para escribir respuestas JSON de forma **consistente** en todos los handlers.

```go
response.JSON(w, 200, task)              // → { "id": 1, "title": "..." }
response.Error(w, 404, "no encontrado") // → { "error": "no encontrado" }
```

Sin esto, cada handler tendria que setear el Content-Type, llamar WriteHeader y json.Encode por separado.

---

### `migrations/`
Archivos SQL que documentan el **esquema de la base de datos**. En este proyecto GORM aplica los cambios automaticamente con `AutoMigrate`, pero los archivos `.sql` son utiles para:
- Entender la estructura de la DB sin leer el codigo Go
- Reproducir el esquema manualmente si es necesario
- Control de versiones del esquema

---

### `.env.example`
Plantilla de las variables de entorno necesarias para correr el proyecto. Se copia a `.env` y se ajustan los valores locales. El `.env` real **nunca se commitea** (esta en `.gitignore`) para no exponer credenciales.

### `docker-compose.yml`
Levanta PostgreSQL en un contenedor Docker sin necesidad de instalarlo. Un solo comando:

```bash
docker-compose up -d
```

### `Makefile`
Atajos para los comandos mas comunes. En lugar de recordar `go test -tags=integration ./...`, simplemente:

```bash
make test-integration
```

---

## Como fluyen los datos en una peticion

```
Cliente (curl / Postman)
    │
    │  POST /api/v1/tasks/  body: {"title": "Aprender Go"}
    ▼
chi router
    │  matchea la ruta y llama al handler
    ▼
handler.Create
    │  decodifica el JSON → CreateTaskRequest
    │  llama svc.Create(req)
    ▼
service.Create
    │  valida que el titulo no este vacio
    │  construye domain.Task
    │  llama repo.Create(&task)
    ▼
gormRepository.Create
    │  ejecuta: INSERT INTO tasks (title, ...) VALUES (...)
    ▼
PostgreSQL
    │  genera el ID, created_at, updated_at
    │  retorna la fila insertada
    ▼
(mismo camino de vuelta)
    ▼
handler.Create
    │  recibe la tarea con ID generado
    │  llama response.JSON(w, 201, task)
    ▼
Cliente recibe: 201 Created  {"id": 1, "title": "Aprender Go", ...}
```

---

## Por que esta separacion de capas

| Sin capas | Con capas |
|-----------|-----------|
| Todo mezclado en un archivo | Cada responsabilidad en su lugar |
| Para testear hay que levantar la DB | Los tests del service corren en milisegundos sin DB |
| Cambiar de PostgreSQL a MongoDB rompe todo | Solo se reescribe el repositorio |
| Un bug de validacion puede estar en cualquier parte | Las validaciones siempre estan en el service |
| Agregar un endpoint requiere entender todo el codigo | Solo se toca el handler y el service |
