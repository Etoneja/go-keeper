# go-keeper - Secure Password Manager

## Quick Start

### 1. Clone and setup
```bash
git clone https://github.com/Etoneja/go-keeper.git
cd go-keeper
make deps
```

### 2. Configuration
```bash
cat > .env << EOF
GOKEEPER_DB_NAME=gokeeper
GOKEEPER_DB_USER=gokeeper  
GOKEEPER_DB_PASSWORD=gokeeper
GOKEEPER_DB_HOST=localhost
GOKEEPER_DB_PORT=5432

GOKEEPER_JWT_SECRET=changeme

GOKEEPER_SERVER_ADDR=127.0.0.1:50051

GOKEEPER_DB_PATH=./dev/fake.db
GOKEEPER_LOGIN=login
GOKEEPER_PASSWORD=password
EOF
```

### 3. Build binaries
```bash
make build-all
```

### 4. Server setup and run (optional)
```bash
# Start database
docker-compose up -d postgres

# Run migrations
make migrate-up

# Build and run server
./bin/keepersrv
```

### 4. Client usage
```bash
./bin/keeperctl version
./bin/keeperctl --help
```

## Client Commands

The client can work in two modes:
* Local mode (no server required)
* Server mode (with running server)

```bash
Zero-Knowledge secret manager

Usage:
  keeperctl [command]

Available Commands:
  add         Add a new secret
  delete      Delete secret by UUID
  get         Get secret by UUID
  help        Help about any command
  init        Initialize local storage
  list        List all secrets
  register    Register new user
  sync        Sync with remote storage
  version     Show version information

Flags:
  -h, --help   help for keeperctl

Use "keeperctl [command] --help" for more information about a command.
```
