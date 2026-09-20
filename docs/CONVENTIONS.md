# Andromeda-Gateway conventions

## Language

- Source comments and exported godoc: **English**
- User-facing docs (README): English or Chinese is fine; keep one style per file

## Layout

```
cmd/andromeda/     # process entry only
internal/config/   # env/flags loading
internal/driver/   # Driver interface + concrete drivers
docs/              # design and conventions
```

- No business plugins in this repo
- Do not import `Plugin-Collections` plugins directly; Lucy owns plugin selection

## Driver rules

- Implement `Driver` in `internal/driver/<name>/`
- Declare only **real** capabilities; never stub unsupported APIs
- Start path: config → `lucy.Bootstrap` → `protocol.Activate(caps)` → `driver.Start`

## Dependencies

- Prefer depending on `Protocol-ConvertTool` for adaptors (`zerobot`, `qqopen`)
- Depend on `Lucy` for application bootstrap
- Do not depend on `botgo` outside the `qqopen` driver (or Protocol qqopen package once moved)

## Go style

- Module path: `github.com/Hafuunano/Andromeda-Gateway`
- Pass `context.Context` on Start/Stop
- Return errors; do not panic for expected failure (bad config, bind port)
- Tabs for Go; `gofmt` / `goimports` required before merge
