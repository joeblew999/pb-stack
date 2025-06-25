# Hosting Strategy for xtask + NATS

**Global Distributed Development Platform Deployment**

## 🎯 Vision

Deploy xtask + NATS infrastructure globally to create a **unified development platform** that runs consistently across any OS and provides low-latency access worldwide through strategic hosting.

## 🚀 Fly.io as Primary Platform

### Why Fly.io is Perfect
- ✅ **BGP Anycast** - Global edge deployment with automatic routing
- ✅ **Multi-process apps** - Run NATS + xtask in single deployment
- ✅ **Cross-platform** - Same containers run everywhere
- ✅ **Global regions** - 30+ regions worldwide for low latency
- ✅ **Private networking** - Secure NATS cluster communication

### Architecture Benefits
```
Developer (anywhere)
    ↓ connects to nearest
Fly.io Edge Region
    ↓ routes via BGP Anycast
NATS Cluster Node
    ↓ coordinates with
Global xtask Network
```

## 🔧 Deployment Architecture

### Multi-Process Fly.io App with Security Stack
```toml
# fly.toml
[build]
  dockerfile = "Dockerfile"

[processes]
  nats = "nats-server --config /etc/nats/nats.conf"
  xtask = "xtask --nats-mode --cluster-url nats://localhost:4222"
  caddy = "caddy run --config /etc/caddy/Caddyfile"
  authelia = "authelia --config /etc/authelia/configuration.yml"
  web = "xtask-web --port 8080"

[services]
  # NATS (internal only)
  [[services.ports]]
    port = 4222
    handlers = ["tcp"]
    internal_port = 4222

  # HTTPS via Caddy (public)
  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]
    internal_port = 443

  # HTTP redirect (public)
  [[services.ports]]
    port = 80
    handlers = ["http"]
    internal_port = 80

  # NATS leaf nodes (authenticated)
  [[services.ports]]
    port = 7422
    handlers = ["tcp"]
    internal_port = 7422

[[services.tcp_checks]]
  port = 4222
  interval = "15s"
  timeout = "10s"

[[services.http_checks]]
  port = 443
  path = "/health"
  interval = "15s"
  tls_skip_verify = false

[[services.http_checks]]
  port = 9091
  path = "/api/health"
  interval = "15s"
```

### Global NATS Cluster Configuration
```yaml
# nats.conf
server_name: "xtask-${FLY_REGION}"
port: 4222
http_port: 8222

jetstream: {
  store_dir: "/data/jetstream"
  max_memory_store: 1GB
  max_file_store: 10GB
}

cluster: {
  name: "xtask-global"
  port: 6222
  routes: [
    "nats://xtask-ams.internal:6222"
    "nats://xtask-dfw.internal:6222"
    "nats://xtask-hkg.internal:6222"
    "nats://xtask-syd.internal:6222"
  ]
}

# Leaf node configuration for external xtask instances
leafnodes: {
  port: 7422
  authorization: {
    users: [
      {user: "xtask", password: "$XTASK_PASSWORD"}
    ]
  }
}
```

### Caddy Configuration
```caddyfile
# /etc/caddy/Caddyfile
{
  # Global options
  auto_https on
  email admin@xtask.dev

  # Enable metrics
  servers {
    metrics
  }
}

# Main xtask platform
xtask.fly.dev {
  # Forward auth to Authelia
  forward_auth authelia:9091 {
    uri /api/verify?rd=https://auth.xtask.fly.dev
    copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
  }

  # xtask web interface
  reverse_proxy /web/* xtask-web:8080

  # NATS monitoring (authenticated)
  reverse_proxy /nats/* localhost:8222

  # Health checks (public)
  reverse_proxy /health localhost:8080

  # API endpoints (authenticated)
  reverse_proxy /api/* xtask-web:8080

  # Static assets
  file_server /static/* {
    root /var/www
  }

  # Security headers
  header {
    Strict-Transport-Security "max-age=31536000; includeSubDomains"
    X-Content-Type-Options "nosniff"
    X-Frame-Options "DENY"
    X-XSS-Protection "1; mode=block"
    Referrer-Policy "strict-origin-when-cross-origin"
  }
}

# Authentication portal
auth.xtask.fly.dev {
  reverse_proxy authelia:9091

  # Security headers for auth portal
  header {
    Strict-Transport-Security "max-age=31536000; includeSubDomains"
    X-Content-Type-Options "nosniff"
    X-Frame-Options "SAMEORIGIN"
    X-XSS-Protection "1; mode=block"
  }
}

# Regional endpoints
xtask-{env.FLY_REGION}.xtask.fly.dev {
  # Same configuration as main, but region-specific
  forward_auth authelia:9091 {
    uri /api/verify?rd=https://auth.xtask.fly.dev
    copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
  }

  reverse_proxy /web/* xtask-web:8080
  reverse_proxy /nats/* localhost:8222
  reverse_proxy /health localhost:8080
  reverse_proxy /api/* xtask-web:8080
}

# NATS leaf node endpoint (TLS + auth)
nats.xtask.fly.dev:7422 {
  reverse_proxy localhost:7422
}
```

### Authelia Configuration
```yaml
# /etc/authelia/configuration.yml
server:
  host: 0.0.0.0
  port: 9091
  path: ""
  enable_pprof: false
  enable_expvars: false

log:
  level: info
  format: text

theme: auto

jwt_secret: ${AUTHELIA_JWT_SECRET}
default_redirection_url: https://xtask.fly.dev

totp:
  issuer: xtask.fly.dev
  algorithm: sha1
  digits: 6
  period: 30
  skew: 1

authentication_backend:
  password_reset:
    disable: false
  refresh_interval: 5m

  file:
    path: /config/users_database.yml
    password:
      algorithm: argon2id
      iterations: 1
      salt_length: 16
      parallelism: 8
      memory: 64

access_control:
  default_policy: deny

  rules:
    # Public health checks
    - domain: "*.xtask.fly.dev"
      policy: bypass
      resources:
        - "^/health.*"

    # Admin access
    - domain: "*.xtask.fly.dev"
      policy: two_factor
      subject: "group:admins"

    # Developer access
    - domain: "*.xtask.fly.dev"
      policy: one_factor
      subject: "group:developers"
      resources:
        - "^/web.*"
        - "^/api.*"

    # NATS monitoring (admin only)
    - domain: "*.xtask.fly.dev"
      policy: two_factor
      subject: "group:admins"
      resources:
        - "^/nats.*"

session:
  name: authelia_session
  domain: xtask.fly.dev
  same_site: lax
  secret: ${AUTHELIA_SESSION_SECRET}
  expiration: 1h
  inactivity: 5m
  remember_me_duration: 1M

regulation:
  max_retries: 3
  find_time: 2m
  ban_time: 5m

storage:
  encryption_key: ${AUTHELIA_STORAGE_ENCRYPTION_KEY}

  postgres:
    host: ${POSTGRES_HOST}
    port: 5432
    database: authelia
    schema: public
    username: ${POSTGRES_USER}
    password: ${POSTGRES_PASSWORD}
    timeout: 5s

notifier:
  disable_startup_check: false

  smtp:
    host: ${SMTP_HOST}
    port: 587
    timeout: 5s
    username: ${SMTP_USERNAME}
    password: ${SMTP_PASSWORD}
    sender: "xtask <noreply@xtask.fly.dev>"
    startup_check_address: "test@xtask.fly.dev"
```

### User Database
```yaml
# /config/users_database.yml
users:
  admin:
    displayname: "Administrator"
    password: "$argon2id$v=19$m=65536,t=3,p=4$..."  # Generated hash
    email: admin@xtask.fly.dev
    groups:
      - admins
      - developers

  developer:
    displayname: "Developer"
    password: "$argon2id$v=19$m=65536,t=3,p=4$..."  # Generated hash
    email: dev@xtask.fly.dev
    groups:
      - developers

  ci:
    displayname: "CI/CD Service"
    password: "$argon2id$v=19$m=65536,t=3,p=4$..."  # Generated hash
    email: ci@xtask.fly.dev
    groups:
      - automation
```

## 🌐 Global Deployment Strategy

### Regional Distribution
```yaml
# Primary regions for global coverage
regions:
  - ams    # Europe (Amsterdam)
  - dfw    # North America (Dallas)
  - hkg    # Asia (Hong Kong)
  - syd    # Oceania (Sydney)
  - gru    # South America (São Paulo)
  - jnb    # Africa (Johannesburg)

# Each region runs:
# - NATS cluster node
# - xtask coordinator
# - Web interface
# - Monitoring stack
```

### Scaling Configuration
```toml
# fly.toml scaling
[scaling]
  min_machines = 2
  max_machines = 10

[regions]
  ams = 2    # Europe primary
  dfw = 2    # North America primary
  hkg = 1    # Asia
  syd = 1    # Oceania
  gru = 1    # South America
  jnb = 1    # Africa
```

## 🐳 Container Strategy

### Multi-Stage Dockerfile with Security Stack
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# Build xtask with all tools embedded
RUN go build -o xtask ./cmd/xtask
RUN go build -o xtask-web ./cmd/xtask-web

# Runtime stage
FROM alpine:latest

# Install system dependencies
RUN apk add --no-cache ca-certificates curl unzip

# Install NATS server
RUN curl -L https://github.com/nats-io/nats-server/releases/download/v2.10.4/nats-server-v2.10.4-linux-amd64.zip -o nats.zip \
    && unzip nats.zip \
    && mv nats-server-v2.10.4-linux-amd64/nats-server /usr/local/bin/ \
    && rm -rf nats*

# Install Caddy
RUN curl -L https://github.com/caddyserver/caddy/releases/download/v2.7.6/caddy_2.7.6_linux_amd64.tar.gz -o caddy.tar.gz \
    && tar -xzf caddy.tar.gz \
    && mv caddy /usr/local/bin/ \
    && rm caddy.tar.gz

# Install Authelia
RUN curl -L https://github.com/authelia/authelia/releases/download/v4.37.5/authelia-v4.37.5-linux-amd64.tar.gz -o authelia.tar.gz \
    && tar -xzf authelia.tar.gz \
    && mv authelia-linux-amd64 /usr/local/bin/authelia \
    && rm authelia.tar.gz

# Copy xtask binaries
COPY --from=builder /app/xtask /usr/local/bin/
COPY --from=builder /app/xtask-web /usr/local/bin/

# Configuration files
COPY nats.conf /etc/nats/nats.conf
COPY Caddyfile /etc/caddy/Caddyfile
COPY authelia-config.yml /etc/authelia/configuration.yml
COPY users_database.yml /config/users_database.yml
COPY entrypoint.sh /entrypoint.sh

# Create necessary directories
RUN mkdir -p /data/jetstream /data/caddy /data/authelia /var/www/static \
    && chmod +x /entrypoint.sh

# Create non-root user for security
RUN addgroup -g 1000 xtask \
    && adduser -D -s /bin/sh -u 1000 -G xtask xtask \
    && chown -R xtask:xtask /data /config

# Volumes for persistent data
VOLUME ["/data", "/config"]

# Expose ports
EXPOSE 80 443 4222 6222 7422 8080 8222 9091

ENTRYPOINT ["/entrypoint.sh"]
```

### Enhanced Entrypoint Script
```bash
#!/bin/sh
# entrypoint.sh

# Set region-specific configuration
export NATS_SERVER_NAME="xtask-${FLY_REGION:-local}"
export FLY_REGION="${FLY_REGION:-local}"

# Generate secrets if not provided
if [ -z "$AUTHELIA_JWT_SECRET" ]; then
  export AUTHELIA_JWT_SECRET=$(openssl rand -base64 32)
fi

if [ -z "$AUTHELIA_SESSION_SECRET" ]; then
  export AUTHELIA_SESSION_SECRET=$(openssl rand -base64 32)
fi

if [ -z "$AUTHELIA_STORAGE_ENCRYPTION_KEY" ]; then
  export AUTHELIA_STORAGE_ENCRYPTION_KEY=$(openssl rand -base64 32)
fi

# Wait for dependencies
wait_for_service() {
  local host=$1
  local port=$2
  local timeout=${3:-30}

  echo "Waiting for $host:$port..."
  for i in $(seq 1 $timeout); do
    if nc -z $host $port 2>/dev/null; then
      echo "$host:$port is ready"
      return 0
    fi
    sleep 1
  done
  echo "Timeout waiting for $host:$port"
  return 1
}

# Start based on process type
case "$1" in
  "nats")
    echo "Starting NATS server..."
    exec nats-server --config /etc/nats/nats.conf
    ;;

  "xtask")
    echo "Starting xtask coordinator..."
    wait_for_service localhost 4222 30
    exec xtask --nats-mode --cluster-url nats://localhost:4222
    ;;

  "caddy")
    echo "Starting Caddy reverse proxy..."
    # Wait for backend services
    wait_for_service localhost 8080 30  # xtask-web
    wait_for_service localhost 9091 30  # authelia
    exec caddy run --config /etc/caddy/Caddyfile
    ;;

  "authelia")
    echo "Starting Authelia authentication..."
    # Wait for database if using external postgres
    if [ -n "$POSTGRES_HOST" ]; then
      wait_for_service $POSTGRES_HOST 5432 60
    fi
    exec authelia --config /etc/authelia/configuration.yml
    ;;

  "web")
    echo "Starting xtask web interface..."
    wait_for_service localhost 4222 30
    exec xtask-web --port 8080 --nats-url nats://localhost:4222
    ;;

  "all")
    echo "Starting all services..."
    # Start services in dependency order
    nats-server --config /etc/nats/nats.conf &
    sleep 5

    authelia --config /etc/authelia/configuration.yml &
    sleep 3

    xtask-web --port 8080 --nats-url nats://localhost:4222 &
    sleep 2

    xtask --nats-mode --cluster-url nats://localhost:4222 &
    sleep 2

    caddy run --config /etc/caddy/Caddyfile &

    # Wait for all background processes
    wait
    ;;

  *)
    echo "Usage: $0 {nats|xtask|caddy|authelia|web|all}"
    echo ""
    echo "Services:"
    echo "  nats     - NATS server with JetStream"
    echo "  xtask    - xtask coordinator node"
    echo "  caddy    - Reverse proxy with TLS"
    echo "  authelia - Authentication and authorization"
    echo "  web      - xtask web interface"
    echo "  all      - Start all services"
    exit 1
    ;;
esac
```

## 🔧 Cross-Platform Compatibility

### Local Development Integration
```yaml
# docker-compose.yml for local development
version: '3.8'
services:
  nats:
    image: nats:alpine
    ports:
      - "4222:4222"
      - "8222:8222"
    command: ["nats-server", "--jetstream", "--http_port", "8222"]

  xtask:
    build: .
    environment:
      - XTASK_NATS_URL=nats://nats:4222
    depends_on:
      - nats
    command: ["xtask", "--nats-mode"]

  web:
    build: .
    ports:
      - "8080:8080"
    environment:
      - XTASK_NATS_URL=nats://nats:4222
    depends_on:
      - nats
    command: ["xtask-web", "--port", "8080"]
```

### Native OS Support
```yaml
# Same xtask binary runs on:
platforms:
  - linux/amd64    # Fly.io, most CI/CD
  - linux/arm64    # ARM servers, Apple Silicon containers
  - darwin/amd64   # Intel Macs
  - darwin/arm64   # Apple Silicon Macs
  - windows/amd64  # Windows development machines

# Connection to global cluster:
XTASK_NATS_URL=nats://xtask.fly.dev:7422 xtask --leaf-node
```

## 🌟 Advanced Features

### BGP Anycast Benefits
```yaml
# Automatic routing to nearest region
developer_locations:
  - "London, UK" → routes to ams (Amsterdam)
  - "New York, US" → routes to dfw (Dallas)
  - "Tokyo, Japan" → routes to hkg (Hong Kong)
  - "Sydney, Australia" → routes to syd (Sydney)

# Sub-100ms latency globally
latency_targets:
  - Europe: <50ms
  - North America: <30ms
  - Asia: <80ms
  - Oceania: <40ms
```

### Intelligent Load Balancing
```go
// xtask automatically discovers optimal cluster node
func connectToCluster() *nats.Conn {
    // Try regional endpoints in order of preference
    endpoints := []string{
        "nats://xtask.fly.dev:7422",           // Anycast (automatic)
        "nats://xtask-ams.fly.dev:7422",      // Europe
        "nats://xtask-dfw.fly.dev:7422",      // North America
        "nats://xtask-hkg.fly.dev:7422",      // Asia
    }

    for _, endpoint := range endpoints {
        if conn, err := nats.Connect(endpoint); err == nil {
            return conn
        }
    }

    return nil
}
```

### Global State Synchronization
```yaml
# JetStream replication across regions
jetstream_config:
  replicas: 3
  placement:
    cluster: "xtask-global"
    tags: ["region:primary"]

  # Streams for global coordination
  streams:
    - name: "xtask-commands"
      subjects: ["xtask.>"]
      retention: "workqueue"
      replicas: 3

    - name: "xtask-results"
      subjects: ["results.>"]
      retention: "limits"
      max_age: "24h"
      replicas: 2
```

## 🚀 Deployment Workflow

### Automated Deployment
```yaml
# .github/workflows/deploy.yml
name: Deploy to Fly.io
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Fly.io
        uses: superfly/flyctl-actions/setup-flyctl@master

      - name: Deploy to all regions
        run: |
          flyctl deploy --region ams,dfw,hkg,syd,gru,jnb
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

### Health Monitoring
```yaml
# Comprehensive health checks
monitoring:
  - NATS cluster health
  - xtask node connectivity
  - JetStream replication status
  - Regional latency metrics
  - Resource utilization
  - Error rates and alerts
```

## 🔒 Security & Enterprise Features

### Enterprise-Grade Security Stack
```
Internet → Caddy (TLS Termination) → Authelia (Auth) → xtask Services
```

#### Multi-Factor Authentication
- **TOTP support** - Time-based one-time passwords
- **WebAuthn ready** - Hardware security keys
- **Session management** - Secure session handling
- **Password policies** - Configurable complexity requirements

#### Role-Based Access Control
```yaml
# Fine-grained permissions
access_control:
  rules:
    - domain: "*.xtask.fly.dev"
      policy: two_factor
      subject: "group:admins"
      resources: ["^/nats.*", "^/admin.*"]

    - domain: "*.xtask.fly.dev"
      policy: one_factor
      subject: "group:developers"
      resources: ["^/web.*", "^/api.*"]

    - domain: "*.xtask.fly.dev"
      policy: bypass
      resources: ["^/health.*", "^/metrics.*"]
```

#### Security Headers & TLS
- **HSTS enforcement** - Strict transport security
- **Content security policies** - XSS protection
- **Automatic TLS certificates** - Let's Encrypt integration
- **TLS 1.3 support** - Modern encryption standards

### Compliance & Governance
- **Audit logging** - All authentication events logged
- **Session recording** - Optional command session recording
- **Data residency** - Regional data storage compliance
- **GDPR compliance** - User data protection controls

## 🌟 Benefits of This Architecture

### For Developers
- ✅ **Global access** - Connect from anywhere with low latency
- ✅ **Consistent experience** - Same platform everywhere
- ✅ **Automatic failover** - If one region fails, others continue
- ✅ **Local development** - Same stack runs locally
- ✅ **Secure by default** - Enterprise authentication built-in
- ✅ **Single sign-on** - One login for all xtask services

### For Organizations
- ✅ **Cost effective** - Pay only for what you use
- ✅ **Globally distributed** - No single point of failure
- ✅ **Scalable** - Automatic scaling based on demand
- ✅ **Compliant** - Data residency via regional deployment
- ✅ **Enterprise security** - Multi-factor auth, RBAC, audit logs
- ✅ **Zero-trust architecture** - Every request authenticated

### For Platform Teams
- ✅ **Simple deployment** - Single Fly.io app with security stack
- ✅ **Unified monitoring** - Global visibility with security metrics
- ✅ **Easy scaling** - Add regions as needed
- ✅ **Automated operations** - Self-healing infrastructure
- ✅ **Security automation** - Automated certificate management
- ✅ **Compliance reporting** - Built-in audit trails

### For Security Teams
- ✅ **Centralized authentication** - Single identity provider
- ✅ **Fine-grained access control** - Resource-level permissions
- ✅ **Comprehensive logging** - All access attempts logged
- ✅ **Modern security standards** - TLS 1.3, HSTS, CSP headers
- ✅ **Multi-factor enforcement** - Configurable MFA policies
- ✅ **Session management** - Automatic session timeouts

## 🔮 Future Enhancements

### Edge Computing Integration
- **CDN integration** for binary distribution
- **Edge functions** for lightweight xtask operations
- **IoT device support** via MQTT bridge
- **Mobile app integration** for remote development

### Advanced Networking
- **Private networks** for enterprise customers
- **VPN integration** for secure access
- **Custom domains** for branded experiences
- **SSL/TLS termination** at edge

**This hosting strategy creates a truly global, resilient, and performant platform for distributed development that scales from individual developers to enterprise teams!** 🚀
