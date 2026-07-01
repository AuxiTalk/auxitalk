# Mock AI

Provides a fake `ai.complete` capability that returns a canned response.

## Purpose

Allows testing the capability router and suggestion flow without calling real LLMs.

## Capability

- `ai.complete`

## Usage

Register this plugin and call `ai.complete` through the capability router.
