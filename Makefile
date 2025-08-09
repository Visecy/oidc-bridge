all: build

output:
	mkdir -p output

conf:
	mkdir -p conf

.PHONY: build
build: output
	go build -o output/oidc-bridge cmd/main.go

.PHONY: test
test:
	go test ./tests/ -v

.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml

.PHONY: keygen
keygen: conf
	openssl genrsa -out conf/private.key 2048
	openssl rsa -in conf/private.key -pubout -out conf/public.key
	@echo "private key and public key generated"

.PHONY: clean
clean:
	rm output/*
