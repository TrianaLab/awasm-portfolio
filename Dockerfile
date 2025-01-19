# Stage 1: Build the WebAssembly binary
FROM golang:1.23.4 AS builder

# ---- Build Arguments ----
# These can be overridden at build time using --build-arg
ARG VERSION
ARG GOOS=js
ARG GOARCH=wasm
ARG APP_WASM=web/assets/app.wasm

# ---- Environment Variables ----
# Set environment variables based on build arguments
ENV GOOS=${GOOS}
ENV GOARCH=${GOARCH}

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the WebAssembly binary with version information
RUN make build VERSION=${VERSION}

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine

# Copy the built WebAssembly assets from the builder stage to Nginx's html directory
COPY --from=builder /app/web /usr/share/nginx/html

# Expose port 80 to allow external access
EXPOSE 80

# Start Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
