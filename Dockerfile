### FRONTEND
FROM node:22-bookworm-slim AS frontend
WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

### BACKEND
FROM golang:1.25-alpine AS backend
WORKDIR /backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN go build -o /app/api-service ./cmd/api

### RUNTIME
FROM alpine:latest
WORKDIR /app

COPY --from=frontend /frontend/build/ ./build/
COPY --from=backend /app/api-service .
RUN touch .env

EXPOSE 8080
CMD ["./api-service", "sql"]
