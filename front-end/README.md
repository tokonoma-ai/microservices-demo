# Front-end

Front-end application written in [Node.js](https://nodejs.org/en/) that puts together all of the microservices under [microservices-demo](https://github.com/tokonoma-ai/microservices-demo).

**Stack:** Node.js 20, Express

## Build

```bash
# Build from repo root
./bin/build --kind   # or --eks

# Or install locally
npm install
```

## Run

```bash
npm start
```

## Test

```bash
make test
```

## Use

```bash
curl http://localhost:8079
```
