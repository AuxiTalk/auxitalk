# Plugin Authoring Guide

This guide explains how to create AuxiTalk plugins.

## Core idea

AuxiTalk plugins are external processes. They communicate with the Go core using JSON-RPC over stdio.

A plugin can be written in any language that can:

- read JSON lines from stdin;
- write JSON lines to stdout;
- write logs to stderr.

## Plugin responsibilities

A plugin should do one thing well.

Examples:

- observe a conversation;
- generate a response suggestion;
- store memory;
- show an overlay;
- send a message after user approval;
- expose an external tool.

## Required files

A plugin should include:

```txt
plugin.json       manifest
README.md         usage and configuration
src/              implementation
```

## Manifest

```json
{
  "id": "example-plugin",
  "name": "Example Plugin",
  "version": "0.1.0",
  "runtime": "node",
  "entry": "src/index.js",
  "kind": "input",
  "permissions": [
    "event.emit"
  ],
  "capabilities": [
    {
      "name": "conversation.observe"
    }
  ]
}
```

## Manifest fields

Required fields:

- `id`: stable unique plugin id;
- `name`: human-readable name;
- `version`: plugin version;
- `runtime`: executable/runtime used to start the plugin;
- `entry`: plugin entrypoint relative to the plugin directory;
- `kind`: plugin category.

Supported kinds:

```txt
input
output
ai
memory
ui
policy
tool
profile
```

Each capability is an object with a required `name` field and optional schemas.

## Communication rules

- stdout is only for JSON-RPC messages.
- stderr is for logs.
- each JSON-RPC message must be one line.
- every request must respect core timeouts.
- large payloads should be avoided.
- sensitive actions must use `action.request`.

## Minimal event

```json
{
  "jsonrpc": "2.0",
  "id": "evt-1",
  "method": "event.emit",
  "params": {
    "type": "message.received",
    "sessionId": "session-1",
    "payload": {
      "text": "hello"
    }
  }
}
```

## Documentation requirement

Every plugin must document:

- what it does;
- required permissions;
- configuration options;
- emitted events;
- provided capabilities;
- security considerations;
- examples of use.
