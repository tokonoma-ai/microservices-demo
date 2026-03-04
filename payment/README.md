# Payment

A microservices-demo service that provides payment authorization.

**Stack:** Go 1.24, Go kit

## Build

```bash
# Build from repo root
./bin/build --kind   # or --eks

# Or build locally
go build ./cmd/paymentsvc/
```

## Run

```bash
./paymentsvc
```

## Check

```bash
curl http://localhost:8082/health
```

## Use

Authorize a payment:

```bash
curl -H "Content-Type: application/json" -X POST -d'{"Amount":40}' http://localhost:8082/paymentAuth
{"authorised":true}
```

Amounts over $1,000 are declined. Amounts <= $0 return "Invalid payment amount".

## API Spec

See the API Spec [here](http://microservices-demo.github.io/api/index?url=https://raw.githubusercontent.com/microservices-demo/payment/master/api-spec/payment.json).
