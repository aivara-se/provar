# Provar Web

The provar-web application is the product website: the landing page, the downloads portal and the user documentation guide, served at [provar.se](https://provar.se).

## Development

```bash
bun install
bun run dev     # development server
bun run lint    # prettier --check .
bun run check   # svelte-kit sync + svelte-check
bun run build   # static site
```

Documentation pages live under `src/routes/docs/`, and the conventions they follow are written up in `src/routes/docs/STYLE.md`.
