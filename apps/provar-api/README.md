# API

`provar-api` is a WebSocket service that bridges remote clients to the Provar SDK — one process, one port, one connection per client. It exists for clients running on another machine; the desktop editor calls the SDK in-process instead of going through it, as [ADR 006](../../docs/adrs/006-desktop-app-is-a-peer-of-provar-api.md) explains.

## Running

From a checkout of this repository, at its root:

```bash
npm run build:api                       # writes bin/provar-api
./bin/provar-api                        # 127.0.0.1:7741
./bin/provar-api -port 8080             # another port
```

`-addr` (default `127.0.0.1`) and `-port` (default `7741`) are the only flags. The service carries no authentication of its own, so the default is localhost-only — pointing `-addr` elsewhere is a deliberate call, not a setting.

## Protocol

- Endpoint: `ws://127.0.0.1:<port>/v1/ws`, one connection per client.
- Plain-text JSON frames, one message per frame.
- Keepalive is WebSocket ping/pong; there is no application-level ping.
- Wire types fall out of the CLI subcommands and the engine's event types — no transport-only messages. A connection carries requests, the event stream for a job, and control messages (stop, pause, resume) for a job already in flight.

The protocol and the reasoning behind it: [ADR 004](../../docs/adrs/004-provar-api-as-a-websocket-service.md).

## Development

Handlers are keyed by message type and live under `handlers/`, with the v1 surface in `handlers/v1/`. Importing the package is what registers them — `cmd/provar-api` blank-imports `handlers` for exactly that reason. Connection envelope and per-client state are in `envelope.go` and `state.go`; dispatch is `dispatch.go`.
