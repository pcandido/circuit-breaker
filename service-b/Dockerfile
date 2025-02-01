# Etapa 1: Build
FROM golang:1.23.2 AS build
WORKDIR /app
COPY . .
RUN go build -o main -tags netgo -ldflags '-w -extldflags "-static"' .

# Etapa 2: Imagem executável
FROM alpine
WORKDIR /app
COPY static static
COPY --from=build /app/main .
EXPOSE 8080
CMD ["./main"]