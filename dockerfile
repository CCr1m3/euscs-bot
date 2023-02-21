FROM node:18-alpine AS frontend
WORKDIR /build
# Copy web source files
COPY web .
# Get dependencies
RUN npm install
# Build static web app files
RUN npm run build

FROM golang:1.20-alpine AS backend
WORKDIR /build
# Copy source files
COPY server server
COPY internal internal
COPY go.mod .
COPY go.sum .
# Get go packages
RUN go mod download
COPY web/web.go web/web.go
COPY --from=frontend /build/dist web/dist
# Build shinpuru backend
RUN go build -o ./bin/omega-strikers-bot ./server/main.go

FROM alpine:3
WORKDIR /app
COPY --from=backend /build/bin .

EXPOSE 9000
ENTRYPOINT ["/app/omega-strikers-bot"]