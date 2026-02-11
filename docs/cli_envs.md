# Inngest CLI & Docker Environment Variables

This file lists the runtime flags and their equivalent `INNGEST_` environment variables. Env names use the `INNGEST_` prefix with the CLI flag name upper‑cased and dashes replaced by underscores (see `cmd/internal/config/config.go`).

## Queue / Executor tuning

| Env var | CLI flag | Default | Purpose / when to change |
| --- | --- | --- | --- |
| `INNGEST_QUEUE_WORKERS` | `--queue-workers` | 100 | Number of executor workers (goroutines) allowed to lease and run queue items. Raise to increase concurrent runs; lower to reduce CPU/memory pressure. |
| `INNGEST_DISABLE_FIFO_FUNCTIONS` | `--disable-fifo-functions` | empty | Comma-separated function IDs to process in parallel within a function partition. Improves throughput for hot functions but relaxes FIFO ordering guarantees. |
| `INNGEST_DISABLE_FIFO_ACCOUNTS` | `--disable-fifo-accounts` | empty | Comma-separated account IDs to process in parallel. Broadens parallelism across all functions in those accounts, with the same FIFO tradeoff. |
| `INNGEST_QUEUE_PEEK_MIN` | `--queue-peek-min` | 300 | Minimum number of items to peek per partition scan. Higher values reduce scheduling round trips but increase per-scan work. |
| `INNGEST_QUEUE_PEEK_MAX` | `--queue-peek-max` | 750 | Maximum number of items to peek per partition scan. Raise for large bursts, but monitor Redis CPU/latency. |
| `INNGEST_QUEUE_MIN_WORKERS_FREE` | `--queue-min-workers-free` | 5 | Minimum free workers required before scanning for new partitions. Lower for throughput, higher for headroom. |
| `INNGEST_QUEUE_PARTITION_LEASE_MS` | `--queue-partition-lease-ms` | 4000 | Partition lease duration in milliseconds. Lower increases re-lease frequency/churn; higher can reduce fairness and increase wait time behind busy partitions. |
| `INNGEST_TICK` | `--tick` | 150 (ms) | Interval at which the executor polls the queue. Shorter ticks reduce latency but increase Redis load. |
| `INNGEST_RETRY_INTERVAL` | `--retry-interval` | 0 | Linear retry backoff (seconds) when linear mode is enabled; leave 0 to use table backoff. |
| `INNGEST_POLL_INTERVAL` | `--poll-interval` | 5 (s) | Polling cadence for app/config discovery. Mostly relevant in dev. |

## Persistence / data stores

| Env var | CLI flag | Default | Purpose |
| --- | --- | --- | --- |
| `INNGEST_REDIS_URI` | `--redis-uri` | in‑memory redis (dev) | Point executor and queue to external Redis. |
| `INNGEST_POSTGRES_URI` | `--postgres-uri` | SQLite (dev) | Use external Postgres for config/history persistence. |
| `INNGEST_POSTGRES_MAX_IDLE_CONNS` | `--postgres-max-idle-conns` | 10 | Max idle DB connections. |
| `INNGEST_POSTGRES_MAX_OPEN_CONNS` | `--postgres-max-open-conns` | 100 | Max open DB connections. |
| `INNGEST_POSTGRES_CONN_MAX_IDLE_TIME` | `--postgres-conn-max-idle-time` | 5 (min) | Idle lifetime for DB conns. |
| `INNGEST_POSTGRES_CONN_MAX_LIFETIME` | `--postgres-conn-max-lifetime` | 30 (min) | Max reuse lifetime for DB conns. |
| `INNGEST_PERSIST` | `--persist` | false | Persist data between restarts (dev only). |
| `INNGEST_SQLITE_DIR` | `--sqlite-dir` | cwd | Directory for SQLite files when persisting in dev. |

## Networking / ports

| Env var | CLI flag | Default | Purpose |
| --- | --- | --- | --- |
| `INNGEST_HOST` | `--host` | 127.0.0.1 | Bind address for HTTP services. |
| `INNGEST_PORT` | `--port` | 8288 (dev default) | Main HTTP/UI/GraphQL port. |
| `INNGEST_CONNECT_GATEWAY_HOST` | `--connect-gateway-host` | 0.0.0.0 | Bind address for Connect gateway. |
| `INNGEST_CONNECT_GATEWAY_PORT` | `--connect-gateway-port` | 8289 | Port for Connect gateway HTTP endpoint. |

## App discovery / SDK

| Env var | CLI flag | Default | Purpose |
| --- | --- | --- | --- |
| `INNGEST_SDK_URL` | `--sdk-url` | none | One or more SDK base URLs to autodiscover functions. Repeat flag or set comma‑separated values. |
| `INNGEST_NO_DISCOVERY` | `--no-discovery` | false | Disable function autodiscovery. |
| `INNGEST_NO_POLL` | `--no-poll` | false | Disable periodic polling for updates (useful when pushing via webhooks only). |

## Auth / security

| Env var | CLI flag | Default | Purpose |
| --- | --- | --- | --- |
| `INNGEST_SIGNING_KEY` | `--signing-key` | none | HMAC signing key for responses where required. |
| `INNGEST_EVENT_KEY` | `--event-key` | none | API key(s) required for ingesting events; repeatable. |

## Docker examples

Run with higher executor concurrency:

```bash
docker run \
  -e INNGEST_QUEUE_WORKERS=1500 \
  -e INNGEST_QUEUE_PEEK_MAX=1500 \
  -e INNGEST_DISABLE_FIFO_FUNCTIONS=8ef2c628-8c4e-41d2-8b27-b89fd8ec5cef \
  -e INNGEST_REDIS_URI=redis://redis:6379 \
  -p 8288:8288 \
  inngest-image
```

Or using flags (equivalent inside Docker entrypoint):

```bash
inngest start \
  --queue-workers 1500 \
  --queue-peek-max 1500 \
  --disable-fifo-functions 8ef2c628-8c4e-41d2-8b27-b89fd8ec5cef \
  --queue-partition-lease-ms 3000 \
  --redis-uri redis://redis:6379
```
