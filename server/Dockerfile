FROM golang:alpine as builder

WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest

LABEL MAINTAINER="sales-pilot"
ENV TZ=Asia/Shanghai
RUN apk update && apk add --no-cache tzdata openntpd \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

COPY --from=0 /app/server ./
COPY --from=0 /app/resource ./resource/
COPY --from=0 /app/config.yaml ./

EXPOSE 8888
ENTRYPOINT ["./server"]
