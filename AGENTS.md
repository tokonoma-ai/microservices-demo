# AGENTS.md

## Cursor Cloud specific instructions

This is a polyglot microservices monorepo (Sock Shop). Each top-level directory is an independent service. See `DEVELOP.md` for full build/deploy docs and `README.md` for the demo workflow.

### Running the full stack locally

Use Docker Compose with pre-built images (no Kubernetes required):

```bash
docker compose -f deploy/docker-compose/docker-compose.yml up -d
```

The front-end is accessible at `http://localhost:80` via the edge-router (Traefik). All 15 containers must be running for end-to-end functionality.

### Service test commands

| Service | Language | Test command |
|---------|----------|-------------|
| front-end | Node.js 20 | `cd front-end && npm test` |
| payment | Go 1.24 | `cd payment && GOFLAGS=-mod=mod go test ./...` |
| catalogue | Go 1.24 | `cd catalogue && GOFLAGS=-mod=mod go test ./...` (has a pre-existing test bug in `service_test.go:132`) |
| carts | Java 17 | `mvn -f carts/pom.xml test` |
| orders | Java 17 | `mvn -f orders/pom.xml test` |
| shipping | Java 17 | `mvn -f shipping/pom.xml test` |

### Lint commands

| Service type | Command |
|-------------|---------|
| Go services | `cd <service> && go vet ./...` |
| Front-end | `cd front-end && npx jshint server.js` (style warnings are pre-existing) |
| Java services | `mvn -f <service>/pom.xml compile` |

### Gotchas

- Go services require `GOFLAGS=-mod=mod` because `go.sum` files are incomplete in the repo. Without this flag, `go test` fails with "missing go.sum entry" errors.
- The `catalogue` service has a pre-existing test compilation error in `service_test.go:132` (wrong arg count in `Errorf`). This is not a setup issue.
- The `payment` service has a pre-existing `go vet` warning about unbuffered `os.Signal` channel in `cmd/paymentsvc/main.go:85`.
- Docker Compose warns about `MYSQL_ROOT_PASSWORD` not being set; this is harmless as `MYSQL_ALLOW_EMPTY_PASSWORD=true` is configured.
- Docker Compose warns about the `version` attribute being obsolete; this is cosmetic.
- Building from source uses `bin/build --kind` or `bin/build --eks` which requires either a `kind` cluster or AWS ECR. For local dev without Kubernetes, build individual services directly with `docker build`.
- Java services can be built locally with Maven: `mvn -f <service>/pom.xml package -DskipTests`.
- The `payment` service has no `docker/payment/Dockerfile` in the repo; it is built via the `bin/build` script which expects it. Build front-end or use `docker build` on services that have Dockerfiles.
