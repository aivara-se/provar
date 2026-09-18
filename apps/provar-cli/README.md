# CLI

The `provar` command-line app compiles test files into runnable scripts and executes them against a running application, which is the surface for CI pipelines and for working from a terminal. It uses the same `libs/engine` as the editor.

## Install

macOS and Linux:

```bash
curl -fsSL https://provar.se/install.sh | bash
```

From a checkout of this repository, at its root:

```bash
npm run build:cli     # writes bin/provar
npm run install:cli   # installs provar onto your PATH
```

## Usage

`provar --help` lists the commands; `provar <command> --help` documents one command's flags.

- `setup` — create a new project directory; `--sample` scaffolds a working test.
- `list` — enumerate the scenarios in a project.
- `validate` — parse-check scenarios without compiling.
- `compile` — turn file definitions into Lua.
- `run` — execute compiled scenarios against a running app.
- `test` — compile and run in one shot.
- `doctor` — diagnose common setup problems.
- `clean` — remove generated artifacts from the project.
- `accept-baseline` — promote the latest screenshots to visual baselines.

```bash
provar setup my-app --sample
cd my-app
provar compile .
provar run .
```

`compile` and `run` need a provider API key in `~/.provar/settings.yml`; `setup`, `list` and `validate` work without one.

Exit codes: `0` success, `1` runtime failure (a failed model call, a missing key, a failing test), `2` usage error, `130` interrupted.

## Development

Each command is one file in `commands/` plus one entry in `commands/registry.go`, which is also the help screen's order. Argument parsing, output formatting, exit codes and signal handling are shared in `helpers/`.
