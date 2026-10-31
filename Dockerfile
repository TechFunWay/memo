# syntax=docker/dockerfile:1

FROM node:20-alpine AS frontend-builder
WORKDIR /build/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.21-alpine AS backend-builder
WORKDIR /build/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/memo .

# CA 根证书：scratch 镜像没有证书包，后端出站 HTTPS（统计/赞赏上报）需要
RUN apk add --no-cache ca-certificates

FROM scratch
WORKDIR /app
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=backend-builder /out/memo /app/memo
COPY --from=frontend-builder /build/frontend/dist/ /app/www/
VOLUME ["/app/data"]
EXPOSE 8905
ENTRYPOINT ["/app/memo", "-port", "8905", "-data-dir", "/app/data", "-web-dir", "/app/www"]
