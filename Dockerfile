# Etapa 1: compilar la aplicación
FROM golang:alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o servidor .

# Etapa 2: preparar la imagen de ejecución
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/servidor .
COPY templates ./templates
COPY static ./static

EXPOSE 8080

CMD ["./servidor"]