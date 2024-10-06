# Stage 1
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cc-service1 .

# Stage 2
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/cc-service1 /app/.env ./

EXPOSE 8080

CMD ["./cc-service1"]