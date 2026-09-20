# Andromeda-Gateway

Main process for the Lucy bot stack: protocol drivers, lifecycle, and webhook/WS ingress.

Lucy (the application layer) registers plugins and middlewares; this gateway selects one driver and activates handlers by capability.

## Role

| Owns | Does not own |
|------|----------------|
| `main`, config, driver start/stop | Plugin business logic |
| OneBot WS / QQ Open Webhook ingress | Protocol segment details (see Protocol-ConvertTool) |
| `protocol.Activate(caps)` before serving | Persona / which plugins to enable (see Lucy) |

## Drivers (planned)

- `onebot` — ZeroBot WebSocket (migrated from Lucy-QOnebot)
- `qqopen` — QQ Open Platform via [botgo](https://github.com/tencent-connect/botgo) **Webhook** (C2C + group)
- `sandbox` — local test IM (`protocol/sandbox` web UI + WS), open `http://127.0.0.1:9100/sandbox/`

Only one driver is active per process in v1.

## Local layout

Sibling repos under the same workspace:

- `Andromeda-Gateway` (this repo)
- `Lucy`
- `Protocol-ConvertTool`
- `Plugin-Collections`
- `Core-SkillAction`

## Develop (run yourself)

```bash
cp .env.example .env
# DRIVER=onebot|qqopen and credentials

go mod tidy
# if botgo APIs fail to resolve:
#   go get github.com/tencent-connect/botgo@fe31c0dfe469

go run ./cmd/andromeda
```

Local `replace` points at sibling repos under `../`.

## Docs

- Conventions: [docs/CONVENTIONS.md](docs/CONVENTIONS.md)
- Design drafts under `docs/superpowers/` are local-only (gitignored)

## License

AGPL-3.0
