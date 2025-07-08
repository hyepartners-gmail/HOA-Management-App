# ---------- Build Frontend ----------
FROM node:24-slim AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend ./
RUN npm run build

# -# ---------- Build Backend ----------
FROM golang:1.24 AS backend
WORKDIR /app

# Copy go.mod + go.sum from root (not inside /backend)
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy actual source code from ./backend
COPY backend ./backend

# Also copy frontend build artifacts
COPY --from=frontend /app/frontend/dist ./frontend_dist

# Build Go binary
RUN GOOS=linux GOARCH=amd64 go build -o server ./backend

# ---------- Final Image ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=backend /app/frontend_dist ./frontend_dist

ENV PORT=8080
CMD ["/app/server"]
