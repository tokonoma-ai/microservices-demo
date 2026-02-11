# Kustomize overlays (kind)

- **default** — Deploys Sock Shop with official `weaveworksdemos/*` images from Docker Hub (default tags).
- **dev** — Same as default but all app images use tag `dev`. Use after building images locally and loading into kind:
  - `docker build -t weaveworksdemos/<service>:dev .` (local only; does **not** push to Docker Hub)
  - `kind load docker-image weaveworksdemos/<service>:dev`

## Deploy from repo root

```bash
./bin/deploy           # overlay: default (official images)
./bin/deploy dev       # overlay: dev (your locally built images, tag "dev")
```
