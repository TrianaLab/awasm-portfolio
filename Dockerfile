# Stage 1: Build the WebAssembly binary
FROM golang:1.26.0 AS builder

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown
ARG GOOS=js
ARG GOARCH=wasm

ENV GOOS=${GOOS}
ENV GOARCH=${GOARCH}

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build VERSION=${VERSION}

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine

LABEL org.opencontainers.image.title=awasm-portfolio
LABEL org.opencontainers.image.description="WebAssembly-powered Kubernetes-style portfolio console"
LABEL org.opencontainers.image.source=https://github.com/TrianaLab/awasm-portfolio
LABEL org.opencontainers.image.licenses=MIT

COPY --from=builder /app/web /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:80/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
