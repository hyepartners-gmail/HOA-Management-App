# ---------- Build Frontend ----------
FROM node:20 AS frontend
WORKDIR /app
COPY frontend ./frontend
WORKDIR /app/frontend
RUN npm ci
RUN npm run build

# ---------- Build Backend ----------
FROM golang:1.22 AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY backend ./backend
COPY --from=frontend /app/frontend/dist ./frontend_dist

RUN go build -o server ./backend

# ---------- Final Image ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=backend /app/frontend_dist ./frontend_dist

ENV PORT=8080
CMD ["/app/server"]
