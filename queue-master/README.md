# Queue Master

A microservices-demo service that reads from the shipping queue and processes shipment requests.

**Stack:** Java 17, Spring Boot 2.7.18, RabbitMQ

## Build

```bash
# Build from repo root
./bin/build --kind   # or --eks
```

## Related Services

- [orders](../orders/) - Creates orders and sends to shipping queue
- [shipping](../shipping/) - Shipping service
