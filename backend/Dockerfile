# Etapa 1: compilar
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Instalar dependencias del sistema para CGO (no las necesitamos, pero por si acaso)
RUN apk --no-cache add ca-certificates

# Copiar fuentes y descargar dependencias
COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o api ./cmd/api

# Etapa 2: imagen mínima de ejecución
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/api /api
EXPOSE 8080
CMD ["/api"]
