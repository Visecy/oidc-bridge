FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o oidc-bridge cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/oidc-bridge .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/private.key .
COPY --from=builder /app/public.key .

EXPOSE 8080

CMD ["./oidc-bridge"]