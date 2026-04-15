# --- Build Frontend ---
FROM node:lts-alpine AS ui-builder

WORKDIR /app/ui
COPY ui/package*.json ./
RUN npm install

COPY ui/ ./
RUN npm run build

# --- Build Backend ---
FROM golang:1.26-alpine AS go-builder

WORKDIR /app

# Install dependencies for building (if any)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Copy built UI assets to the location expected by go:embed in assets/assets.go
COPY --from=ui-builder /app/ui/dist ./assets/ui/dist

RUN go build -o /app/server ./cmd/server

# --- Final Image ---
FROM alpine:latest

# Install CA certificates for Prometheus API requests (HTTPS)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /
COPY --from=go-builder /app/server /server

# Default port
EXPOSE 8080

ENTRYPOINT ["/server"]
CMD ["-l", ":8080"]
