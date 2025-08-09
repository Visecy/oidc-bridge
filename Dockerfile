FROM golang:1.22-alpine AS builder

WORKDIR /app

# 设置Go模块代理以加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o oidc-bridge cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/oidc-bridge .

EXPOSE 8080

CMD ["./oidc-bridge"]