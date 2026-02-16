#!/usr/bin/env python3
"""
Inject one checkout failure for the demo user.
Adds items to cart and submits order; orders service throws a runtime
exception so checkout always fails with HTTP 500.
Run every 15 min via CronJob, or once at demo startup.
"""
import os
import sys

import requests

USER_ID = os.environ.get("DEMO_USER_ID", "57a98d98e4b00679b4a830af")
USER_SVC = os.environ.get("USER_SVC", "http://user.sock-shop.svc.cluster.local")
CARTS_SVC = os.environ.get("CARTS_SVC", "http://carts.sock-shop.svc.cluster.local")
ORDERS_SVC = os.environ.get("ORDERS_SVC", "http://orders.sock-shop.svc.cluster.local")
CATALOGUE_SVC = os.environ.get("CATALOGUE_SVC", "http://catalogue.sock-shop.svc.cluster.local")

ITEM_PRICE = 7.99

CATALOGUE_IDS = [
    "6d62d909-f957-430e-8689-b5129c0bb75e",
    "a0a4f044-b040-410d-8ead-4de0446aec7e",
    "808a2de1-1aaa-4c25-a9b9-6612e8f29a38",
]


def _resolve_url(href):
    """Use USER_SVC for fetches from loadtest ns (user returns http://user/...)."""
    if not href or not href.startswith("http"):
        return f"{USER_SVC}{href}" if href else None
    # user service returns http://user/... ; from loadtest we need FQDN
    if href.startswith("http://user") or href.startswith("http://user/"):
        base = USER_SVC.rstrip("/")
        path = href.replace("http://user", "").replace("http://user/", "/")
        return f"{base}{path}" if path.startswith("/") else f"{base}/{path}"
    return href


def fetch_order_payload():
    """Fetch customer, address, card from user service (same flow as front-end)."""
    r = requests.get(f"{USER_SVC}/customers/{USER_ID}", timeout=5)
    if r.status_code != 200:
        raise RuntimeError(f"customer fetch failed: {r.status_code} {r.text[:200]}")
    data = r.json()
    if data.get("status_code") == 500:
        raise RuntimeError(f"customer error: {data.get('error', 'unknown')}")
    links = data.get("_links", {})
    customer_href = links.get("customer", {}).get("href")
    address_link = links.get("addresses", {}).get("href")
    card_link = links.get("cards", {}).get("href")
    if not all((customer_href, address_link, card_link)):
        raise RuntimeError("customer missing links")

    # Fetch address and card (resolve user -> user.sock-shop.svc for loadtest ns)
    addr_r = requests.get(_resolve_url(address_link), timeout=5)
    card_r = requests.get(_resolve_url(card_link), timeout=5)
    if addr_r.status_code != 200:
        raise RuntimeError(f"address fetch failed: {addr_r.status_code}")
    if card_r.status_code != 200:
        raise RuntimeError(f"card fetch failed: {card_r.status_code}")

    addr_data = addr_r.json()
    card_data = card_r.json()
    addrs = addr_data.get("_embedded", {}).get("address", [])
    cards = card_data.get("_embedded", {}).get("card", [])
    if not addrs or not cards:
        raise RuntimeError("address/card empty")
    addr_self = addrs[0].get("_links", {}).get("self", {}).get("href")
    card_self = cards[0].get("_links", {}).get("self", {}).get("href")
    if not addr_self or not card_self:
        raise RuntimeError("address/card missing self href")

    # Pass hrefs as-is to orders (orders runs in sock-shop, resolves "user")
    return {
        "customer": customer_href,
        "address": addr_self,
        "card": card_self,
        "items": f"{CARTS_SVC}/carts/{USER_ID}/items",
    }


def fetch_item():
    for iid in CATALOGUE_IDS:
        try:
            r = requests.get(f"{CATALOGUE_SVC}/catalogue/{iid}", timeout=3)
            if r.status_code == 200:
                d = r.json()
                return d.get("id", iid), float(d.get("price", ITEM_PRICE))
        except Exception:
            pass
    return CATALOGUE_IDS[0], ITEM_PRICE


def run():
    item_id, price = fetch_item()

    # Clear cart
    try:
        requests.delete(f"{CARTS_SVC}/carts/{USER_ID}", timeout=5)
    except Exception:
        pass

    # Add a few items to the cart (any amount triggers the exception)
    for _ in range(3):
        try:
            r = requests.post(
                f"{CARTS_SVC}/carts/{USER_ID}/items",
                json={"itemId": item_id, "unitPrice": price},
                timeout=5,
            )
            if r.status_code not in (200, 201):
                print(f"add to cart returned {r.status_code}", file=sys.stderr)
        except Exception as e:
            print(f"add to cart failed: {e}", file=sys.stderr)
            sys.exit(1)

    # Fetch order payload from user service (same as front-end)
    try:
        order = fetch_order_payload()
    except Exception as e:
        print(f"fetch order payload failed: {e}", file=sys.stderr)
        sys.exit(1)

    # Submit order (orders service throws exception -> 500)
    try:
        r = requests.post(f"{ORDERS_SVC}/orders", json=order, timeout=10)
        if r.status_code == 500:
            print("checkout failed as expected (orders exception)")
            sys.exit(0)
        if r.status_code == 201:
            print("unexpected: order succeeded", file=sys.stderr)
            sys.exit(1)
        print(f"orders returned {r.status_code}: {r.text[:200]}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"order failed: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    run()
