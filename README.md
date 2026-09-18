# Provar

Provar is a local, git-native end-to-end testing tool for web applications. A user journey becomes a visual flow you sketch on a canvas and keep in git; an AI agent then picks the selectors and writes the browser automation underneath, so a test is not invalidated every time a button moves. When the UI drifts, recompile the affected test and the agent rewrites the broken step. Tests run on your machine or in CI, with no cloud account.

![The Provar editor canvas showing a login flow that branches into a personal task list and a work task list, with one branch passing and the other failed](apps/provar-web/static/screenshot.png)

Provar is in early access: the desktop editor and the CLI are usable today, and the [roadmap](docs/ROADMAP.md) tracks what comes next.

## Install

Download the desktop app — editor, engine and runtime in one package — from [provar.se](https://provar.se). macOS builds are available now, Windows and Linux are coming.

For a terminal, the CLI installs on macOS and Linux:

```bash
curl -fsSL https://provar.se/install.sh | bash
```

## First run

Scaffold a project that points at a live sample application, add a provider API key, then compile and run it:

```bash
provar setup my-app --sample
# put your provider API key in ~/.provar/settings.yml
cd my-app
provar compile .
provar run .
```

`provar --help` lists every command — `setup`, `list`, `validate`, `compile`, `run`, `test`, `doctor`, `clean`, `accept-baseline` — and each command documents its own flags.

## Concepts

- **project** — a directory holding the application under test and its `.provar/` configuration.
- **file** — one test: a YAML list of actions under `.provar/tests/`, with the `.test.yml` extension. A file is a directed acyclic graph, so signed-out and signed-in are two paths through one file rather than two tests.
- **action** — one user-intent step, such as "register as a new user". Compilation turns that intent into step-by-step instructions, and running executes them against a real browser.

These concepts and the full product specification live in [docs/PRODUCT.md](docs/PRODUCT.md).

## Repository layout

| Path | Contents |
| --- | --- |
| `apps/provar-app` | Graphical editor for authoring and debugging tests visually. |
| `apps/provar-cli` | The `provar` command-line app for compiling and running tests headlessly. |
| `apps/provar-api` | The backend service the desktop app talks to. |
| `apps/provar-web` | The product website and guide, published at [provar.se](https://provar.se). |
| `libs/domain` | Domain models and configuration management. |
| `libs/engine` | Compilation and the execution runtime. |
| `libs/models` | One streaming interface over the supported model providers. |
| `libs/logger` | Logging shared by the apps and libraries. |
| `docs` | Product, system, design and roadmap documents, plus architecture decision records. |

Architecture and the boundaries between these packages: [docs/SYSTEMS.md](docs/SYSTEMS.md).

## Build from source

Go 1.26 or newer is required. The frontends are installed with Bun inside their own directories, and the desktop shell needs its platform toolchain (Wails).

```bash
npm run all          # go fmt, go vet, go test, then build the CLI and API binaries
npm run build:cli    # bin/provar
npm run build:api    # bin/provar-api
npm run run          # build the CLI, then launch it
npm run tidy         # after changing dependencies
```

For a frontend, from its own directory: `bun install && bun run dev`.

## Documentation

- [PRODUCT.md](docs/PRODUCT.md) — what Provar is and who it is for.
- [SYSTEMS.md](docs/SYSTEMS.md) — monorepo architecture and library boundaries.
- [DESIGN.md](docs/DESIGN.md) — the UI design system.
- [ROADMAP.md](docs/ROADMAP.md) and [TODOS.md](docs/TODOS.md) — current state and next work.
- [docs/adrs](docs/adrs) — architecture decision records.
- [provar.se/docs](https://provar.se/docs) — quickstart, authoring, running, CI and troubleshooting guides.

## Contributing

Branches are lowercase with hyphens, history stays linear, and commits follow Conventional Commits without scopes. Before opening a pull request run `go fmt ./...` and `go vet ./...` — or `npm run all`, which runs those two and then the tests and both builds. [AGENTS.md](AGENTS.md) holds the full conventions, with the coding, testing and writing guides under `.agents/skills`.

## License

MIT © Aivara. See [LICENSE](LICENSE).
