# OIDC Bridge

这是一个OIDC桥接服务，用于将现有的OAuth2.0服务转换为符合OpenID Connect协议的服务。

## 功能

- Discovery端点 (/.well-known/openid-configuration)
- Authorization端点 (/authorize)
- Token端点 (/token)
- UserInfo端点 (/userinfo)
- JWKS端点 (/.well-known/jwks.json)

## 配置

配置文件为`config.yaml`，包含以下配置项：

- `op_authorize_url`: OP的授权端点URL
- `op_token_url`: OP的Token端点URL
- `op_userinfo_url`: OP的UserInfo端点URL
- `issuer`: Issuer标识
- `id_token_lifetime`: ID Token生命周期（秒）
- `nonce_cache_ttl`: nonce缓存TTL（秒）
- `id_token_signing_alg`: ID Token签名算法
- `scope_mapping`: Scope映射
- `user_attribute_mapping`: 用户属性映射
- `redis_addr`: Redis地址（可选，如果未提供或连接失败，服务将降级到本地内存缓存）
- `private_key_path`: 私钥路径
- `public_key_path`: 公钥路径

## 部署

### 准备工作

在部署服务之前，您需要克隆代码仓库并生成用于签名ID Token的RSA密钥对：

```bash
# 克隆代码仓库
cd /opt
git clone https://github.com/Visecy/oidc-bridge.git
cd oidc-bridge

# 生成私钥
make keygen
```

### 配置文件编写指南

创建一个 `config.yaml` 文件，内容如下：

```yaml
# OP 端点
op_authorize_url: "https://your-op.com/oauth/authorize"
op_token_url: "https://your-op.com/oauth/token"
op_userinfo_url: "https://your-op.com/oauth/userinfo"

# Issuer 标识
issuer: "https://your-oidc-bridge.com"

# ID Token 设置
id_token_lifetime: 3600  # 1 小时
nonce_cache_ttl: 600    # 10 分钟
id_token_signing_alg: "RS256"

# Scope 映射
scope_mapping:
  openid: "profile email"
  profile: "name picture"
  email: "email"

# 用户属性映射
user_attribute_mapping:
  sub: "user_id"
  name: "full_name"
  email: "email_address"
  picture: "avatar_url"

# Redis 地址（可选）
# redis_addr: "localhost:6379"

# 密钥路径
private_key_path: "/path/to/private.key"
public_key_path: "/path/to/public.key"
```

### 本地运行

1. 安装Go 1.22或更高版本
2. 运行`go mod tidy`安装依赖
3. 运行`make build`编译项目
4. 运行`./output/oidc-bridge`启动服务

您可以通过命令行参数或环境变量指定自定义配置文件、密钥路径和端口：

```bash
# 使用命令行参数
./output/oidc-bridge --config=/opt/oidc-bridge/config.yaml --private-key=/opt/oidc-bridge/private.key --public-key=/opt/oidc-bridge/public.key --port=8080

# 使用环境变量
CONFIG_FILE=/opt/oidc-bridge/config.yaml PRIVATE_KEY_PATH=/opt/oidc-bridge/private.key PUBLIC_KEY_PATH=/opt/oidc-bridge/public.key ./output/oidc-bridge
```

### Docker部署

1. 构建镜像: `docker build -t oidc-bridge .`
2. 运行容器: `docker run -p 8080:8080 -v /opt/oidc-bridge/conf:/root/conf oidc-bridge --config=/root/conf/config.yaml --private-key=/root/conf/private.key --public-key=/root/conf/public.key`

### Docker Compose部署

创建一个`docker-compose.yml`文件，内容如下：

```yaml
version: '3.8'

services:
  oidc-bridge:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - /opt/oidc-bridge/config.yaml:/root/config.yaml
      - /opt/oidc-bridge/private.key:/root/private.key
      - /opt/oidc-bridge/public.key:/root/public.key
    environment:
      - REDIS_ADDR=redis:6379
    depends_on:
      - redis

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```

然后使用以下命令运行服务：

```bash
docker-compose up -d
```

## 测试

### 基本测试

可以使用以下命令进行基本测试：

```bash
# 获取Discovery文档
curl http://localhost:8080/.well-known/openid-configuration

# 获取JWKS
curl http://localhost:8080/.well-known/jwks.json
```

### 单元测试

项目包含全面的单元测试套件，覆盖了所有主要模块。

运行所有测试：

```bash
go test ./tests/...
```

运行特定模块的测试：

```bash
# 运行handler模块的测试
go test ./tests/*_test.go

# 运行service模块的测试
go test ./tests/*_service_test.go
```

注意：某些测试可能需要Redis服务运行在localhost:6379，并且需要有效的密钥文件。