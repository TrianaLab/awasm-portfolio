# Stage 1: Build the WebAssembly binary
FROM golang:1.23.4 AS builder

# Set environment variables for WebAssembly build
ENV GOOS=js
ENV GOARCH=wasm

# Set the working directory
WORKDIR /app

# Copy the Go modules manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the WebAssembly binary
RUN go build -o web/assets/app.wasm main.go

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine

# Copy the built files to the Nginx html directory
COPY --from=builder /app/web /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
