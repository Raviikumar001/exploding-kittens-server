# Exploding Kittens Go Server

Fiber-based API for the Exploding Kittens game.

## Environment

You can use either a Redis URL or discrete connection fields.

- REDIS_URL (preferred in single var form)
- or
  - REDIS_HOST (default: redis)
  - REDIS_PORT (default: 6379)
  - REDIS_DB (default: 0)
  - REDIS_PASSWORD (optional)
- JWT_SECRET (required)
- PORT (default: 8080)

Example `.env` (app.env):

```
PORT=8080
JWT_SECRET=change-me
# Option A
REDIS_URL=redis://localhost:6379/0
# Option B
# REDIS_HOST=localhost
# REDIS_PORT=6379
# REDIS_DB=0
# REDIS_PASSWORD=
```

## Run locally (Go)

```bash
# from this folder
go run server.go
```

## Docker

```bash
# Build image
docker build -t exploding-kittens-api .

# Run container
# (configure Redis via REDIS_URL or REDIS_HOST/PORT/PASSWORD/DB)
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e JWT_SECRET=super-secret \
  -e REDIS_HOST=host.docker.internal \
  -e REDIS_PORT=6379 \
  exploding-kittens-api
```

## Docker Compose (with Redis)

```bash
# from the parent 'server' folder that contains docker-compose.yml
docker compose up --build
```

Services:
- redis: port 6379, with persisted volume
- api: port 8080, configured to use the redis service

Update the `JWT_SECRET` in the compose file before deploying.

## Health check

GET `/` returns a simple JSON message.
