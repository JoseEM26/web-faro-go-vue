# Go vs Spring Boot

Comparacion directa entre dos de los stacks mas usados para backends. Spring Boot domina el mundo enterprise Java; Go domina infraestructura, microservicios y sistemas de alto rendimiento. Aqui se explica por que son distintos y en que brilla cada uno.

---

## Tabla general de comparacion

| Caracteristica | Go | Spring Boot |
|---|---|---|
| Lenguaje | Go (compilado) | Java / Kotlin (JVM) |
| Tiempo de arranque | ~5-50 ms | ~3-15 segundos |
| Uso de memoria en reposo | ~10-30 MB | ~150-400 MB |
| Binario final | Unico ejecutable estatico | JAR + JVM instalada |
| Concurrencia | Goroutines (muy livianas) | Threads del SO (pesados) |
| Curva de aprendizaje | Baja (spec del lenguaje: 50 paginas) | Alta (anotaciones, contexto, beans, AOP...) |
| Ecosistema | Mas pequeno pero creciendo | Enorme y maduro |
| Tipado | Estatico, inferido | Estatico, verboso |
| Generics | Si (desde Go 1.18) | Si (Java 5+) |
| ORM principal | GORM | Hibernate / JPA |
| Router popular | chi, gin, echo | Spring MVC integrado |
| Testing | Testing stdlib incluida | JUnit + Mockito + Spring Test |
| Deploy | Binario o contenedor minimo | JAR o contenedor con JVM |

---

## Comparaciones de codigo lado a lado

### Hola Mundo — servidor HTTP

**Go** (stdlib pura, sin framework):
```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hola mundo")
    })
    http.ListenAndServe(":8080", nil)
}
```

**Spring Boot:**
```java
@SpringBootApplication
public class App {
    public static void main(String[] args) {
        SpringApplication.run(App.class, args);
    }
}

@RestController
public class HelloController {
    @GetMapping("/")
    public String hello() {
        return "Hola mundo";
    }
}
```

Go no necesita framework para levantar un servidor. La stdlib incluye todo.

---

### Endpoint REST con modelo

**Go** (con chi):
```go
type Task struct {
    ID    uint   `json:"id"`
    Title string `json:"title"`
}

func getTask(w http.ResponseWriter, r *http.Request) {
    task := Task{ID: 1, Title: "Aprender Go"}
    json.NewEncoder(w).Encode(task)
}
```

**Spring Boot:**
```java
@Entity
public class Task {
    @Id @GeneratedValue
    private Long id;
    private String title;
    // getters, setters, constructores...
}

@RestController
@RequestMapping("/tasks")
public class TaskController {
    @GetMapping("/{id}")
    public Task getTask(@PathVariable Long id) {
        return new Task(1L, "Aprender Java");
    }
}
```

Go serializa structs a JSON directamente con tags. Spring Boot necesita anotaciones en la entidad, en el controlador y en el mapeo de rutas.

---

### Manejo de errores

**Go** — explicito, sin excepciones:
```go
task, err := repo.FindByID(id)
if err != nil {
    // el compilador te obliga a manejar el error
    return domain.Task{}, err
}
return task, nil
```

**Spring Boot** — excepciones implicitas:
```java
// Si el repo lanza una excepcion, puede silenciarse accidentalmente
Task task = repo.findById(id)
    .orElseThrow(() -> new ResourceNotFoundException("Task not found"));
return task;
```

En Go los errores son valores que se retornan. El compilador advierte si los ignoras. En Java las excepciones pueden propagarse silenciosamente y aparecer en produccion.

---

### Goroutines vs Threads

**Go** — una goroutine usa ~2 KB de memoria y el runtime las multiplexa:
```go
// Lanzar 10,000 goroutines es completamente normal
for i := 0; i < 10_000; i++ {
    go func(id int) {
        procesarPeticion(id)
    }(i)
}
```

**Spring Boot** — un thread del SO usa ~1 MB de stack:
```java
// 10,000 threads = ~10 GB de memoria, inviable
// Spring WebFlux intenta resolver esto con programacion reactiva (mas compleja)
Flux.range(1, 10_000)
    .flatMap(id -> procesarPeticion(id))
    .subscribe();
```

Go maneja concurrencia masiva de forma natural y sencilla. Spring Boot necesita WebFlux (reactivo) para escalar igual, pero la curva de aprendizaje es mucho mayor.

---

### Channels — comunicacion entre goroutines

Go tiene **channels** como primitiva nativa del lenguaje para pasar datos entre goroutines de forma segura:

```go
resultados := make(chan string, 3)

go func() { resultados <- buscarEnDB() }()
go func() { resultados <- buscarEnCache() }()
go func() { resultados <- buscarEnAPI() }()

// Espera los 3 resultados concurrentemente
for i := 0; i < 3; i++ {
    fmt.Println(<-resultados)
}
```

En Spring Boot esto requiere `CompletableFuture`, `ExecutorService`, o WebFlux con `Flux/Mono`, todo mas verboso y con mas boilerplate.

---

### Variables de entorno / configuracion

**Go** (con godotenv):
```go
// Leer una variable de entorno
port := os.Getenv("PORT")
```

**Spring Boot** (`application.yml` + anotacion):
```yaml
# application.yml
server:
  port: ${PORT:8080}
database:
  url: ${DATABASE_URL}
```
```java
@Value("${server.port}")
private int port;
```

Go usa directamente las variables del sistema operativo. Spring Boot tiene su propio sistema de configuracion en capas (`application.yml`, `application.properties`, perfiles, etc.) que es poderoso pero complejo.

---

## Lo que Go hace mejor que Spring Boot

### 1. Tiempo de arranque ultrarapido

Go compila a binario nativo. No necesita JVM, no hay calentamiento del JIT compiler.

| | Arranque | Util para |
|---|---|---|
| Go | 5-50 ms | Serverless, scaling rapido, Kubernetes |
| Spring Boot | 3-15 seg | Aplicaciones que corren dias sin reiniciar |

En entornos cloud donde los pods se crean y destruyen frecuentemente (AWS Lambda, Cloud Run, Kubernetes HPA), el arranque lento de Spring Boot es un problema real. Go no tiene este problema.

---

### 2. Consumo de memoria 10x menor

Una API REST simple en reposo:

| | Memoria | Descripcion |
|---|---|---|
| Go | ~15 MB | Solo el binario + heap minimo |
| Spring Boot | ~200-350 MB | JVM + Spring context + beans + JIT |

En un cluster de 20 microservicios, la diferencia entre 300 MB x 20 = 6 GB y 15 MB x 20 = 300 MB es la diferencia entre necesitar 3 nodos de servidor o 1.

---

### 3. Un unico binario estatico

Go compila todo en un **unico ejecutable**:

```bash
# Compilar para Linux desde Windows/Mac
GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

# El binario se copia y ejecuta directamente. Sin JVM, sin dependencias.
scp api servidor:/usr/local/bin/
ssh servidor "/usr/local/bin/api"
```

```dockerfile
# Dockerfile de Go: imagen final minima
FROM scratch                     # imagen vacia, 0 MB de base
COPY api /api
CMD ["/api"]
# Imagen final: ~15 MB
```

```dockerfile
# Dockerfile de Spring Boot
FROM eclipse-temurin:21-jre      # ~200 MB solo la JVM
COPY app.jar /app.jar
CMD ["java", "-jar", "/app.jar"]
# Imagen final: ~250-400 MB
```

---

### 4. Concurrencia sin frameworks externos

En Go la concurrencia es parte del lenguaje con `go`, `chan`, `select` y `sync`:

```go
// Worker pool nativo sin ninguna libreria
func procesarTareas(tareas []Tarea) {
    jobs := make(chan Tarea, len(tareas))
    var wg sync.WaitGroup

    // Lanzar 5 workers
    for w := 0; w < 5; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for tarea := range jobs {
                procesar(tarea)
            }
        }()
    }

    // Enviar trabajo
    for _, t := range tareas {
        jobs <- t
    }
    close(jobs)
    wg.Wait()
}
```

En Spring Boot necesitas `@Async`, `ThreadPoolTaskExecutor`, `CompletableFuture` o WebFlux, cada uno con su propia configuracion y forma de manejar errores.

---

### 5. Compilacion detecta mas errores antes de ejecutar

```go
// Esto no compila — variable declarada pero no usada
func main() {
    x := 5  // error: x declared and not used
}

// Esto no compila — import no usado
import "fmt"  // error: "fmt" imported and not used

// Esto no compila — interfaz no satisfecha completamente
type MiRepo struct{}
var _ repository.TaskRepository = MiRepo{}  // falla si falta algun metodo
```

El compilador de Go es muy estricto: variables sin usar, imports sin usar, tipos incompatibles. Muchos bugs que en Java aparecen en runtime, en Go no compilan.

---

### 6. Testing sin dependencias externas

Go incluye un framework de testing en la stdlib. No se necesita JUnit, Mockito ni ninguna libreria extra:

```go
// Esto corre con: go test ./...
func TestCreate_TituloVacio(t *testing.T) {
    svc := service.NewTaskService(newMockRepo())

    _, err := svc.Create(domain.CreateTaskRequest{Title: ""})

    if !errors.Is(err, service.ErrInvalidTitle) {
        t.Errorf("esperaba ErrInvalidTitle, obtuvo %v", err)
    }
}
```

```go
// Tests de HTTP tampoco necesitan Spring Test ni MockMvc
func TestCreate_Exitoso(t *testing.T) {
    body, _ := json.Marshal(domain.CreateTaskRequest{Title: "Test"})
    req := httptest.NewRequest("POST", "/tasks/", bytes.NewReader(body))
    rec := httptest.NewRecorder()
    router.ServeHTTP(rec, req)

    if rec.Code != http.StatusCreated {
        t.Errorf("esperaba 201, obtuvo %d", rec.Code)
    }
}
```

---

### 7. Sin magia oculta — el codigo hace lo que dice

Spring Boot usa reflexion, proxies dinamicos y anotaciones para hacer "magia" en tiempo de ejecucion:

```java
// @Transactional crea un proxy dinamico en runtime
// Si llamas este metodo desde otra funcion de la misma clase, la transaccion NO funciona
// Este bug es muy comun y dificil de detectar
@Transactional
public void actualizarTarea(Long id) { ... }
```

En Go lo que se ve es lo que pasa. No hay proxies, no hay contexto de Spring, no hay contenedor de beans. Si una funcion hace algo, lo ves en el codigo.

---

### 8. Cross-compilation nativa

Compilar para cualquier sistema operativo y arquitectura desde tu maquina:

```bash
GOOS=linux   GOARCH=amd64  go build -o api-linux    ./cmd/api
GOOS=windows GOARCH=amd64  go build -o api.exe       ./cmd/api
GOOS=darwin  GOARCH=arm64  go build -o api-mac-m1    ./cmd/api
```

Con Spring Boot necesitas una JVM instalada en cada maquina destino. El JAR corre "en cualquier lugar" donde haya Java instalado, pero eso no es lo mismo que un binario nativo.

---

### 9. Velocidad de compilacion

Go fue disenado especificamente para compilar rapido:

```
Proyecto mediano:
  Go:          go build  →  0.5 - 2 segundos
  Spring Boot: mvn package → 15 - 60 segundos
```

Esto impacta directamente en cuanto tiempo esperas al iterar durante el desarrollo.

---

### 10. `defer` — limpieza garantizada de recursos

```go
func leerArchivo(nombre string) error {
    f, err := os.Open(nombre)
    if err != nil {
        return err
    }
    defer f.Close()  // se ejecuta SIEMPRE al salir de la funcion, sin importar como

    // procesar archivo...
    return nil
}
```

`defer` garantiza que los recursos (conexiones, archivos, locks) se liberan siempre, aunque haya un error. En Java `try-finally` y `try-with-resources` hacen algo similar pero con mas verbosidad.

---

## Lo que Spring Boot hace mejor que Go

No todo es a favor de Go. Spring Boot gana en:

| Aspecto | Por que Spring Boot gana |
|---|---|
| **Ecosistema maduro** | Miles de librerias, integraciones con todo (AWS, Kafka, Redis, OAuth...) |
| **ORM avanzado** | Hibernate/JPA es mas poderoso que GORM para modelos complejos |
| **Seguridad** | Spring Security es completo y probado en produccion por decadas |
| **Empresas grandes** | Equipos grandes con Java prefieren Spring por su estandarizacion |
| **Herramientas** | IntelliJ + Spring tiene un soporte IDE excepcional |
| **Comunidad** | Mas tutoriales, Stack Overflow, libros disponibles |

---

## Cuando elegir cada uno

### Elegir Go cuando:

- Microservicios con muchas instancias (el costo de memoria importa)
- Herramientas CLI o utilitarios del sistema
- Servicios con concurrencia masiva (WebSockets, streaming, proxies)
- Entornos serverless (arranque rapido es critico)
- El equipo valora la simplicidad y el codigo explicito
- APIs REST simples y de alto rendimiento

### Elegir Spring Boot cuando:

- Aplicaciones enterprise con logica de negocio muy compleja
- El equipo ya conoce Java y el ecosistema JVM
- Se necesita integracion con muchos sistemas externos
- Se usa Hibernate para modelos de datos complejos con relaciones avanzadas
- Se requiere Spring Security para autenticacion OAuth/SSO compleja

---

## Resumen en una frase

> **Spring Boot** te da todo listo pero con mayor complejidad y mayor consumo de recursos.
> **Go** te da control total, maximo rendimiento y simplicidad, pero construyes mas desde cero.

Ninguno es mejor en absoluto. Son herramientas distintas para contextos distintos.
