# CLAUDE.md

Claude and other AI coding agents should follow `AGENTS.md` as the source of truth for this repository.

## Quick start

Read in this order:

1. `AGENTS.md`
2. `docs/product/vision.md`
3. `README.md`
4. `docs/workflow-engine.md`
5. relevant package files for the task

## Default verification

Run before completing code changes:

```sh
go test ./...
```

Run `gofmt` on changed Go files.

## Important constraints

- Do not push unless explicitly requested.
- Do not add real side effects without approval.
- Keep secrets out of logs and commits.
- Preserve JSON-RPC stdio protocol for plugins.
- Prefer small, testable slices.

## Project framing

AuxiTalk is not only a chat assistant. It is an event-driven automation runtime for workflows across chats, terminals, APIs, plugins, dashboards, and AI agents.
