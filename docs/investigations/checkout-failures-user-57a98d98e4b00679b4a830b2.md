# Checkout Failures Investigation: User 57a98d98e4b00679b4a830b2

**Date:** 2026-03-20  
**User:** James Cooper (`james_cooper`), ID: `57a98d98e4b00679b4a830b2`  
**Index:** `sockshop` (Tokonoma)  
**Period:** 2026-03-20 00:00 UTC – 2026-03-20 16:02 UTC

---

## Summary

User `57a98d98e4b00679b4a830b2` experienced **12 checkout failures** over a 16-hour window on 2026-03-20. The primary cause is a corrupted catalogue entry (item `bad00002-0000-0000-0000-000000000002`, named "Cursed") with a **negative unit price of -$5**, which makes the cart total negative and causes the payment service to reject the authorization with `"Invalid payment amount"`. A secondary, intermittent issue of payment service connection failures also contributed on 3 of those attempts.

---

## Root Cause

### Issue 1 (Primary): Negative-priced catalogue item

**Catalogue item `bad00002-0000-0000-0000-000000000002`:**

```json
{
  "id": "bad00002-0000-0000-0000-000000000002",
  "name": "Cursed",
  "price": -5,
  "count": 50,
  "tag": ["magic"],
  "description": "This sock has a negative price due to a data entry error. When added to cart with other items, it can cause the order total to become negative, triggering payment validation failures."
}
```

The item's own description documents it as a data entry error. The front-end fetches this item from the catalogue service and passes `unitPrice=-5` directly to the carts service, which accepts it without validation. The resulting negative cart total causes payment authorization to fail.

**Failure chain:**

1. `front-end` adds item `bad00002` with `unitPrice=-5` to cart for user `57a98d98e4b00679b4a830b2`
2. `carts` accepts the item — no negative-price validation
3. `orders` POSTs to `http://payment/paymentAuth` with a negative total
4. `payment` returns HTTP 500: `{"error":"Invalid payment amount","status_code":500}`
5. `orders` logs: `action=new_order customer_url=.../57a98d98e4b00679b4a830b2 status=failed err=...HttpServerErrorException$InternalServerError: 500 Internal Server Error`

**Sample error log (orders service):**
```
2026-03-20 00:03:04.263 ERROR 1 --- [p-nio-80-exec-4] w.w.s.o.controllers.OrdersController
  : action=new_order customer_url=http://user/customers/57a98d98e4b00679b4a830b2
    status=failed
    err=org.springframework.web.client.HttpServerErrorException$InternalServerError:
    500 Internal Server Error: {"error":"Invalid payment amount","status_code":500,"status_text":"Internal Server Error"}
```

### Issue 2 (Secondary): Payment service connection refused

On 3 of the 12 failures (at ~04:27, ~12:00, ~16:01 UTC), the payment service was unreachable:

```
err=org.springframework.web.client.ResourceAccessException:
  I/O error on POST request for "http://payment/paymentAuth":
  Connection refused; nested exception is java.net.ConnectException: Connection refused
```

This is a separate availability issue that compounds the checkout failure rate.

---

## Impact

| Metric | Value |
|--------|-------|
| Failures for user `57a98d98e4b00679b4a830b2` | 12 |
| Caused by "Invalid payment amount" | 9 |
| Caused by payment service connection refused | 3 |
| Total "Invalid payment amount" errors (all users, same window) | 112 |
| Log events involving item `bad00002` (all services) | 693 |
| Services affected by `bad00002` | `front-end` (330), `carts` (231), `catalogue` (132) |
| Failure window | 2026-03-20 00:03 – 16:02 UTC (~16 hours) |

Other affected users also observed in the same window: `57a98d98e4b00679b4a830af`, `57a98d98e4b00679b4a830b5`.

---

## Recommended Fixes

### Immediate

1. **Remove or correct catalogue item `bad00002-0000-0000-0000-000000000002`** — fix the price to a valid positive value or remove the item entirely from the catalogue database.

### Short-term

2. **Add price validation in the carts service** — reject `add_to_cart` requests with `unitPrice <= 0` and return a 400 Bad Request.
3. **Add order total validation in the orders service** — before calling `paymentAuth`, validate that the computed total is positive. Return a clear user-facing error if not.
4. **Investigate payment service restarts** — three connection refusals over 16 hours indicate instability. Review pod crash logs and resource limits for the `payment` deployment.

### Long-term

5. **Add catalogue ingestion validation** — enforce that `price > 0` when adding or updating catalogue items.
6. **Add integration tests** for negative-price cart scenarios and payment service unavailability.

---

## Evidence (Tokonoma Display Files)

- [Error pattern insights (all services)](https://mcp.tokonoma.ai/display/log-patterns-20260321-194211-473773.txt)
- [User 57a98d98e4b00679b4a830b2 error logs](https://mcp.tokonoma.ai/display/search-results-20260321-194238-810272.txt)
- ["Invalid payment amount" error details with stack traces](https://mcp.tokonoma.ai/display/search-results-20260321-194239-711570.txt)
- [Item bad00002 log activity across services](https://mcp.tokonoma.ai/display/search-results-20260321-194240-390953.txt)
- [Item bad00002 by service breakdown](https://mcp.tokonoma.ai/display/aggregate-results-20260321-194302-016535.txt)
