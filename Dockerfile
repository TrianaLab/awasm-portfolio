# ── Stage 1: build the Go WASM artifact ─────────────────────────────────
FROM golang:1.26.0 AS wasm-builder

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

WORKDIR /app

# Cache go.mod first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Produces web/assets/app.wasm + web/scripts/wasm_exec.js.
RUN make wasm VERSION=${VERSION}

# ── Stage 2: build the Svelte + Vite frontend ──────────────────────────
FROM node:22-alpine AS ui-builder

WORKDIR /app

# Re-use the WASM artifacts produced in stage 1 so Vite assembles
# everything under web/ in one shot.
COPY --from=wasm-builder /app /app

WORKDIR /app/frontend
RUN npm ci --ignore-scripts && npm run build

# ── Stage 3: serve the static bundle ───────────────────────────────────
FROM nginx:stable-alpine

LABEL org.opencontainers.image.title=awasm-portfolio
LABEL org.opencontainers.image.description="WebAssembly-powered Kubernetes-style portfolio console"
LABEL org.opencontainers.image.source=https://github.com/TrianaLab/awasm-portfolio
LABEL org.opencontainers.image.licenses=MIT

COPY --from=ui-builder /app/web /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:80/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
