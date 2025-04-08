# Password Manager

[![Go Version](https://img.shields.io/github/go-mod/go-version/GlebRadaev/password-manager)](https://golang.org/)
[![CI](https://github.com/GlebRadaev/password-manager/actions/workflows/password-manager.yml/badge.svg?branch=main)](https://github.com/GlebRadaev/password-manager/actions/workflows/password-manager.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/GlebRadaev/password-manager)](https://goreportcard.com/report/github.com/GlebRadaev/password-manager)
[![Codecov](https://codecov.io/github/GlebRadaev/password-manager/graph/badge.svg?token=IHLJM4LT2Y)](https://codecov.io/github/GlebRadaev/password-manager)

**Password Manager** — A client-server system for secure storage of passwords, text/binary data with cross-device synchronization.

## Features

- User registration, authentication and authorization
- Secure storage for:
  - Login/password pairs
  - Arbitrary text data
  - Binary data
  - Bank card details
- Cross-device synchronization
- CLI client for Windows, Linux and macOS

## System Architecture

```mermaid
graph LR
    A[CLI Client] -->|HTTP/8079| B[Gateway]
    B -->|gRPC/9090| C[AuthService]
    B -->|gRPC/9091| D[DataService]
    B -->|gRPC/9092| E[SyncService]
    E -->|gRPC/9091| D
    style A fill:#e1f5fe,stroke:#039be5,color:#01579b
    style B fill:#e8f5e9,stroke:#43a047,color:#2e7d32
    style C fill:#ffebee,stroke:#e53935,color:#c62828
    style D fill:#f3e5f5,stroke:#8e24aa,color:#6a1b9a
    style E fill:#fff8e1,stroke:#ffb300,color:#ff8f00
    linkStyle default stroke:#616161,stroke-width:2px
```

### Component Interaction:

1. **Client Side**:

   - CLI client communicates exclusively with Gateway via **HTTP** (port `:8079`)
   - Request format: REST/JSON
   - Gateway serves as single entry point (API Gateway)

2. **Server Side**:

   - Gateway routes requests to respective services via **gRPC**:
     - Authentication → AuthService(`:9090`)
     - Data operations → DataService (`:9091`)
     - Synchronization → SyncService (`:9092`)
   - All inter-service communication uses gRPC

3. **Port Mapping**:

| Service         | HTTP Port | gRPC Port | Primary Functions               |
| --------------- | --------- | --------- | ------------------------------- |
| **Gateway**     | `:8079`   | `:9089`   | Request routing, API entrypoint |
| **AuthService** | `:8080`   | `:9090`   | User management                 |
| **DataService** | `:8081`   | `:9091`   | Secure data storage             |
| **SyncService** | `:8082`   | `:9092`   | Cross-device synchronization    |

**Notes**:

- Core business logic implemented via gRPC
- Clients CANNOT access services directly (only through Gateway)

## Configuration

Main configuration file: `local.config.yaml` in root directory

## Installation

### Requirements

- Go 1.21+
- Docker и Docker Compose
- Protobuf Compiler

### 1. Build Project

```bash
# Clone repository
git clone https://github.com/GlebRadaev/password-manager.git
cd password-manager

# Install dependencies
make bin
make gen

# Build client for all platforms
make build-all
```

### 2. Start Server

```bash
# Build and launch via Docker
make run
```

### 3. Using the CLI Client

Compiled binaries are located in bin/ directory. Example usage:

```bash
# Show available commands
./bin/pm-linux-amd64 --help

# Example registration
./bin/pm-linux-amd64 register --username test --password test1234 --email test@test.com
```

### Documentation

- SSwagger UI available after launch at configured port: `http://localhost:{{port}}/swagger/`
- OpenAPI 3.0 specification: docs/swagger/
