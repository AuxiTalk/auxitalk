# Plugin Registry

The plugin registry lives in `internal/plugins`.

## Responsibilities

- Load `plugin.json` files.
- Parse plugin manifests.
- Validate manifests using `pkg/types.PluginManifest`.
- Store plugin metadata by id.
- Reject duplicate plugin ids.
- Return deterministic plugin lists sorted by id.

## Manifest loading

`LoadManifest(path)` reads a plugin manifest and returns:

- absolute manifest path;
- plugin root directory;
- validated manifest data.

## Registry behavior

The registry is in-memory and concurrency-safe.

Duplicate ids return `ErrPluginAlreadyRegistered`.

## Current scope

This phase only validates and registers plugin metadata.

It does not start plugin processes yet. Process lifecycle belongs to the future plugin supervisor phase.
