# OIDC Bridge

[中文版本](README_cn.md)

This is an OIDC bridge service that converts existing OAuth 2.0 services into services compliant with the OpenID Connect protocol.

## Features

- Discovery endpoint (/.well-known/openid-configuration)
- Authorization endpoint (/authorize)
- Token endpoint (/token)
- UserInfo endpoint (/userinfo)
- JWKS endpoint (/.well-known/jwks.json)

## Configuration

The configuration file is `config.yaml`, which includes the following configuration items:

- `op_authorize_url`: OP's authorization endpoint URL
- `op_token_url`: OP's token endpoint URL
- `op_userinfo_url`: OP's userinfo endpoint URL
- `issuer`: Issuer identifier
- `id_token_lifetime`: ID Token lifetime (seconds)
- `nonce_cache_ttl`: Nonce cache TTL (seconds)
- `id_token_signing_alg`: ID Token signing algorithm
- `scope_mapping`: Scope mapping
- `user_attribute_mapping`: User attribute mapping
- `redis_addr`: Redis address (optional, if not provided or connection fails, the service will fall back to local memory cache)
- `private_key_path`: Private key path
- `public_key_path`: Public key path

## Deployment

### Prerequisites

Before deploying the service, you need to clone the repo and generate RSA key pairs for signing ID tokens:

```bash
# Clone the repository
cd /opt
git clone https://github.com/Visecy/oidc-bridge.git
cd oidc-bridge

# Generate a private key
make keygen
```

### Configuration File Guide

Create a `config.yaml` file with the following content:

```yaml
# OP endpoints
op_authorize_url: "https://your-op.com/oauth/authorize"
op_token_url: "https://your-op.com/oauth/token"
op_userinfo_url: "https://your-op.com/oauth/userinfo"

# Issuer identifier
issuer: "https://your-oidc-bridge.com"

# ID Token settings
id_token_lifetime: 3600  # 1 hour
nonce_cache_ttl: 600    # 10 minutes
id_token_signing_alg: "RS256"

# Scope mapping
scope_mapping:
  openid: "profile email"
  profile: "name picture"
  email: "email"

# User attribute mapping
user_attribute_mapping:
  sub: "user_id"
  name: "full_name"
  email: "email_address"
  picture: "avatar_url"

# Redis address (optional)
# redis_addr: "localhost:6379"

# Key paths
private_key_path: "/path/to/private.key"
public_key_path: "/path/to/public.key"
```

### Local Deployment

1. Install Go 1.22 or higher
2. Run `go mod tidy` to install dependencies
3. Run `make build` to compile the project
4. Run `./output/oidc-bridge` to start the service

You can specify a custom configuration file, key paths, and port using command line arguments or environment variables:

```bash
# Using command line arguments
./output/oidc-bridge --config=/opt/oidc-bridge/config.yaml --private-key=/opt/oidc-bridge/private.key --public-key=/opt/oidc-bridge/public.key --port=8080

# Using environment variables
CONFIG_FILE=/opt/oidc-bridge/conf/config.yaml PRIVATE_KEY_PATH=/opt/oidc-bridge/conf/private.key PUBLIC_KEY_PATH=/opt/oidc-bridge/conf/public.key ./output/oidc-bridge
```

### Docker Deployment

1. Build the image: `docker build -t oidc-bridge .`
2. Run the container: `docker run -p 8080:8080 -v /opt/oidc-bridge/conf:/root/conf oidc-bridge --config=/root/conf/config.yaml --private-key=/root/conf/private.key --public-key=/root/conf/public.key`

### Docker Compose Deployment

Create a `docker-compose.yml` file with the following content:

```yaml
version: '3.8'

services:
  oidc-bridge:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./conf:/root/conf
    environment:
      - REDIS_ADDR=redis:6379
      - CONFIG_FILE=/root/conf/config.yaml
      - PRIVATE_KEY_PATH=/root/conf/private.key
      - PUBLIC_KEY_PATH=/root/conf/public.key
      - GIN_MOD=release
    depends_on:
      - redis

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```

Then run the services using:

```bash
docker-compose up -d
```

## Testing

### Basic Testing

You can perform basic testing with the following commands:

```bash
# Get Discovery document
curl http://localhost:8080/.well-known/openid-configuration

# Get JWKS
curl http://localhost:8080/.well-known/jwks.json
```

### Unit Testing

The project includes a comprehensive unit test suite covering all major modules.

Run all tests:

```bash
go test ./tests/...
```

Run tests for specific modules:

```bash
# Run handler module tests
go test ./tests/*_test.go

# Run service module tests
go test ./tests/*_service_test.go
```

Note: Some tests may require Redis service running at localhost:6379 and valid key files.