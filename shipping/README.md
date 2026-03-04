# Shipping

A microservices-demo service that provides shipping capabilities.

**Stack:** Java 17, Spring Boot 2.7.18, RabbitMQ

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

## Check

```bash
curl http://localhost:8080/health
```

## Use

```bash
curl http://localhost:8080
```
