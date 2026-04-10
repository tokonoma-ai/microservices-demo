# Security Notes

This document records known security issues in the Sock Shop demo application. Some are fixed; others are documented here because the fix requires architectural changes beyond the scope of this PR.

## Fixed in this PR

| Issue | Location | Fix |
|---|---|---|
| Hardcoded RabbitMQ `guest`/`guest` credentials | `shipping/`, `queue-master/` | Credentials now read from `RABBITMQ_USER` / `RABBITMQ_PASS` env vars (default `guest` for dev) |
| Hardcoded session secret `'sooper secret'` | `front-end/config.js` | Now read from `SESSION_SECRET` env var; throws in production if unset |
| Missing `sameSite` flag on `logged_in` cookie | `front-end/api/user/index.js` | `sameSite: 'Strict'` added to both login and register cookie responses |
| `custId` query-param impersonation bypass | `front-end/helpers/index.js` | Development-only override removed entirely |
| `WriteHeader` called before `Content-Type` set | `catalogue/transport.go` | Header ordering corrected |
| Duplicate `err.Error()` argument in test | `catalogue/service_test.go` | Extra argument removed |

## Known Remaining Issues (not fixed here)

### SHA-1 password hashing (High)
- **Location**: `user/api/service.go` (`calculatePassHash`), `user/users/users.go` (`NewSalt`)
- **Risk**: SHA-1 is cryptographically broken and unsuitable for password storage. An attacker with the MongoDB dump can crack passwords using rainbow tables or GPU acceleration.
- **Fix needed**: Replace with bcrypt, scrypt, or Argon2id. Requires a migration for existing hashed passwords.

### Plaintext credit card numbers in MongoDB (High)
- **Location**: `user/users/cards.go`, `user/db/mongodb/mongodb.go`, `deploy/kubernetes/manifests/27a-user-db-seed-cm.yaml`
- **Risk**: Card `longNum`, `expires`, and `ccv` fields are stored unencrypted. Only the last 4 digits are masked on read; the full number is persisted.
- **Fix needed**: Encrypt card data at rest (application-level encryption or a secrets manager). This is a demo; real PCI-DSS compliance requires tokenization.

### No CSRF protection (Medium)
- **Location**: `front-end/server.js`, all POST/DELETE routes
- **Risk**: State-changing requests (add to cart, place order, register) carry no CSRF token. An attacker can forge cross-origin requests from a page the user visits while their `md.sid` session cookie is valid.
- **Fix needed**: Add `csurf` middleware (or equivalent) to the Express app and include the token in all mutating forms/AJAX calls.

### Unauthenticated admin-style endpoints (Medium)
- **Location**: `user/api/transport.go` (`GET /customers`, `GET /cards`, `GET /addresses`; `DELETE /customers/:id`)
- **Risk**: Any client that can reach the user service can enumerate all users, cards, and addresses, or delete any user, with no authentication.
- **Fix needed**: Add authentication middleware on the user service. For the demo topology these routes are not exposed externally, but there is no defense-in-depth if the network perimeter is bypassed.

### `logged_in` cookie readable by JavaScript (Low)
- **Location**: `front-end/api/user/index.js`
- **Risk**: The `logged_in` cookie does not have the `httpOnly` flag, so client-side JavaScript can read it. This is an intentional design choice (client.js uses `$.cookie('logged_in')` to gate UI actions) but it means an XSS vulnerability could exfiltrate the cookie.
- **Fix needed**: Refactor client-side auth-state check to use a non-`httpOnly` indicator cookie (or a JS-readable endpoint) separate from the session identifier, then mark the session cookie `httpOnly`.

### Hardcoded MySQL credentials in Kubernetes manifests (Low)
- **Location**: `deploy/kubernetes/manifests/07-catalogue-db-dep.yaml` (`MYSQL_ROOT_PASSWORD: fake_password`), `07a-catalogue-db-seed-cm.yaml` (`default_password`)
- **Risk**: Credentials committed in plaintext to the repository.
- **Fix needed**: Use Kubernetes Secrets (or an external secrets manager) and inject via `secretKeyRef`.

### Grafana admin password in Docker Compose (Low)
- **Location**: `deploy/docker-compose/docker-compose.monitoring.yml` (`GF_SECURITY_ADMIN_PASSWORD=foobar`)
- **Risk**: Grafana admin account secured with a known default credential committed to the repository.
- **Fix needed**: Read from an environment variable or Docker secret.
