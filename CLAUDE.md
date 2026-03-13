# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Linkme is a linktree-style personal links page built with Go compiled to WebAssembly, deployed on Cloudflare Workers via the `syumai/workers` framework.

## Architecture

Single-file Go app (`main.go`) that:
- Embeds `templates/index.html` and `static/style.css` at compile time via `//go:embed`
- Uses `chi` router with one GET route (`/`) and a catch-all redirect to `/`
- Renders the HTML template with inline CSS (injected as `template.CSS`)
- Runs on Cloudflare Workers via `syumai/workers` (Go → WASM)

The build pipeline (`syumai/workers/cmd/workers-assets-gen`) generates the worker shim in `build/` (runtime.mjs, wasm_exec.js, worker.mjs) alongside the compiled `app.wasm`.

## Commands

```bash
npm run dev      # local dev server (wrangler dev)
npm run build    # build Go to WASM + generate worker assets
npm run deploy   # deploy to Cloudflare Workers (wrangler deploy)
```

## Customization

All content (name, bio, avatar, favicon, links) is configured in the `about` var in `main.go`. The `Link` struct has `Title`, `URL`, and `Icon` (emoji) fields.

## Key Dependencies

- `github.com/go-chi/chi/v5` — HTTP router
- `github.com/syumai/workers` — Go WASM ↔ Cloudflare Workers bridge
- `wrangler` (npm) — Cloudflare Workers CLI for dev/deploy

## Build Target

Go compiles with `GOOS=js GOARCH=wasm`. The wrangler config (`wrangler.jsonc`) points to `build/worker.mjs` as the entrypoint.
