# Etapa 1: Build
FROM golang:1.23.2 AS build
WORKDIR /app
COPY . .
RUN go build -o main -tags netgo -ldflags '-w -extldflags "-static"' .

# Etapa 2: Imagem executável
FROM alpine
WORKDIR /app
COPY --from=build /app/main .
CMD ["./main"]