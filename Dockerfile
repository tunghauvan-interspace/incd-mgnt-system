FROM node:18-alpine AS frontend-builder

WORKDIR /app/web/frontend

# Copy frontend package files
COPY web/frontend/package*.json ./
RUN npm ci

# Copy frontend source code and build
COPY web/frontend/ ./
RUN npm run build

FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend assets
COPY --from=frontend-builder /app/web/static ./web/static

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from backend-builder
COPY --from=backend-builder /app/main .

# Copy web assets including built frontend
COPY --from=backend-builder /app/web ./web

# Copy migrations directory
COPY --from=backend-builder /app/migrations ./migrations

# Create directory for templates
RUN mkdir -p web/templates web/static

EXPOSE 8080

CMD ["./main"]