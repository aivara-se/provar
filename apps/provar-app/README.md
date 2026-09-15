# Editor

The Provar editor is the desktop application for authoring and debugging tests visually: the canvas is the graph of actions in a test file, and compiling or running paints the outcome back onto the nodes. It calls `libs/domain` and `libs/engine` directly rather than talking to `provar-api`, for the reasons in [ADR 006](../../docs/adrs/006-desktop-app-is-a-peer-of-provar-api.md).

## Setup

Beyond Go, the app needs the Wails CLI and the frontend's dependencies. Wails drives the frontend through the commands in `wails.json`, so install and build happen from here:

```bash
wails dev     # development build, hot reload on frontend changes
wails build   # production build
```

The frontend's own scripts (`dev`, `build`, `check`) also run from `frontend/`, which is the quickest way to type-check a UI change without starting the app.

## Layout

- `frontend/` — the editor UI: canvas, panels, modals and stores.
- `frontend/wailsjs/` — generated bindings; do not edit by hand, they are regenerated from the Go side.
- `internal/bindings/` — the Go methods exposed to the frontend.
- `internal/testfile/` — the graph view the canvas renders, and the inverse that saves it back as action lists.
- `internal/menu/` — the application menu.
