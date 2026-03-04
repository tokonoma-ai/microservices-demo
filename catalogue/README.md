# Catalogue

A microservices-demo service that provides catalogue/product information.

**Stack:** Go 1.24, Go kit, MySQL

## Build

```bash
# Build from repo root
./bin/build --kind   # or --eks

# Or build locally
go build ./cmd/cataloguesvc/
```

## Run

```bash
./cataloguesvc
```

## Test

```bash
make test
```

## Check

```bash
curl http://localhost:8080/health
```

## Use

```bash
curl http://localhost:8080/catalogue
```

## API Spec

See the API Spec [here](http://microservices-demo.github.io/api/index?url=https://raw.githubusercontent.com/microservices-demo/catalogue/master/api-spec/catalogue.json).
