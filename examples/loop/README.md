# First Complete Loop Example

This example demonstrates the full AuxiTalk runtime loop using only mock plugins.

## Flow

1. `mock-input` emits a simulated inbound message.
2. Event is published to the internal event bus.
3. `SessionManager` updates the session and appends the message.
4. `ContextBuilder` generates a compact context.
5. `CapabilityRouter` calls `mock-ai` via `ai.complete`.
6. Suggestion is created and emitted.
7. `console-output` prints the suggestion.
8. `file-memory` records the feedback event.
9. `ActionGate` decides whether the suggestion can be sent.

## Purpose

Validate that all core components work together before integrating real plugins.

## How to run

This is a conceptual example. In a future version it will be executable via `auxitalkd --example loop`.
