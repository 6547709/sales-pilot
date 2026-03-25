# 构建阶段：Debian Bookworm（非 Alpine）
FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# 运行阶段：UBI9 微镜像（无 Alpine）
FROM registry.access.redhat.com/ubi9/ubi-micro:latest
COPY --from=build /out/server /usr/local/bin/sales-pilot-api
EXPOSE 8080
USER 1001:1001
ENTRYPOINT ["/usr/local/bin/sales-pilot-api"]
