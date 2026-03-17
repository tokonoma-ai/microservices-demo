## Sock Shop checkout investigation (user `57a98d98e4b00679b4a830b2`)

- Runtime availability:
  - Docker daemon is reachable via `DOCKER_HOST=tcp://127.0.0.1:2375`.
  - Key services started at:
    - `/docker-compose-orders-db-1 2026-03-17T23:22:46.281104527Z`
    - `/docker-compose-orders-1 2026-03-17T23:22:46.579166546Z`
    - `/docker-compose-payment-1 2026-03-17T23:22:46.248100622Z`
    - `/docker-compose-front-end-1 2026-03-17T23:22:46.698977728Z`
  - Source: `container_started_at.txt`, `container_status.txt`

- Orders DB checks:
  - Initial query (after stack start): `{ "user_count" : 0, "last24h_count" : 0, "all_last24h_count" : 0 }`
    - Source: `orders_counts.json`
  - After traffic/repro: `{ "user_count" : 18, "all_count" : 18 }`
    - Source: `orders_counts_after_repro.json`
  - Recent persisted orders for this user are low total amounts (roughly 12.98 to 31.99).
    - Source: `orders_user_last24h_compact.txt`, `orders_for_user_after_repro.json`

- Login identity mapping:
  - Runtime user-db record for `_id=57a98d98e4b00679b4a830b2` is `username: "user"` (not `james_cooper`).
  - Source: `user_57a98_runtime_record.json`

- Reproduction:
  - `GET /login` with `user:password` returned `200` and set cookie.
    - Source: `login_user_headers.txt`, `login_user_body.txt`
  - `POST /cart` succeeded (`201`).
    - Source: `add_cart_headers.txt`
  - `POST /orders` returned `406` with:
    - `PaymentDeclinedException`
    - `Payment declined: amount exceeds 100.00`
    - Source: `order_headers.txt`, `order_body.txt`

- Service logs:
  - `orders` log shows payment response not authorised for this flow:
    - `Received payment response: PaymentResponse{authorised=false, message=Payment declined: amount exceeds 100.00}`
    - Source: `orders_after_repro.log`
  - `payment` log shows failed auth decisions (`method=Authorise result=false`) at matching timestamps.
    - Source: `payment_after_repro.log`
  - `front-end` log shows the same user and downstream 406 responses for `/orders`.
    - Source: `front-end_after_repro.log`

- Concurrent load evidence:
  - `user-sim` (locust) ran during this window and reported:
    - `POST /orders 17 reqs, 5 fails (22.73%)`
    - Error: `406 Client Error: Not Acceptable for url: http://edge-router/orders`
  - Source: `user-sim.log`
