# Orders

A microservices-demo service that provides ordering capabilities.

**Stack:** Java 17, Spring Boot 2.7.18, MongoDB

## Build

```bash
# Build from repo root
./bin/build --kind   # or --eks

# Or build locally with Maven
mvn -DskipTests package
```

## Run

```bash
mvn spring-boot:run
```

## Use

```bash
curl http://localhost:8082
```

## API Spec

See the API Spec [here](http://microservices-demo.github.io/api/index?url=https://raw.githubusercontent.com/microservices-demo/orders/master/api-spec/orders.json).
