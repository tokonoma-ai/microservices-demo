# User

A microservices-demo service that covers user account storage, including cards and addresses.

**Stack:** Go 1.24, Go kit, MongoDB

## Build

```bash
# Build from repo root
./bin/build --kind   # or --eks

# Or build locally
make build
```

## Run

```bash
# Start MongoDB first
docker-compose up -d user-db

# Run the service
./bin/user -port=8080 -database=mongodb -mongo-host=localhost:27017
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

Test user account passwords can be found in the comments in `users-db-test/scripts/customer-insert.js`.

```bash
# List customers
curl http://localhost:8080/customers

# List cards
curl http://localhost:8080/cards

# List addresses
curl http://localhost:8080/addresses

# Login
curl http://localhost:8080/login

# Register
curl http://localhost:8080/register
```

## API Spec

See the API Spec [here](http://microservices-demo.github.io/api/index?url=https://raw.githubusercontent.com/microservices-demo/user/master/apispec/user.json).
