# Mock Input

Emits simulated conversation messages for testing the AuxiTalk runtime.

## Purpose

Used during development to inject fake inbound messages without requiring real apps.

## Manifest

See `plugin.json`.

## Usage

This plugin is loaded by the registry and started by the supervisor during tests.

It should emit events via the JSON-RPC `event.emit` method.
