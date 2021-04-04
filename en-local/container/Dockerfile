# Generar imagen de API Rest para incidencias en dos pasos
FROM golang:latest AS builder

# Disponibilizar dependencias
RUN go get github.com/gorilla/mux

# Compilar API y generar ejecutable
WORKDIR /app
COPY ./*.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o ./api_rest_incidencias .

FROM alpine:latest
#FROM scratch
RUN apk --no-cache add ca-certificates

# Copiar artefacto generado por el builder en el paso anterior
COPY --from=builder /app/api_rest_incidencias ./
RUN chmod +x ./api_rest_incidencias

# Configurar el puerto de escucha para la aplciación.
EXPOSE 8080

# Ejecutar API
CMD ["./api_rest_incidencias"]
