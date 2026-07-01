# Official Plugin Template Plan

> Plano do template oficial de plugins do AuxiTalk.  
> Plan for the official AuxiTalk plugin template.

---

## Português (PT-BR)

## Objetivo

Criar um repositório oficial `AuxiTalk/plugin-template` para servir como base simples, segura e bem documentada para qualquer pessoa criar plugins do AuxiTalk.

O template deve ser fácil para humanos e agentes de IA entenderem, copiarem, modificarem e publicarem.

## Linguagem inicial

O primeiro template oficial será em **TypeScript/Node.js**.

Motivos:

- comunidade ampla;
- boa integração com navegador, desktop e automação;
- fácil distribuição;
- bom suporte para JSON-RPC/stdin/stdout;
- amigável para agentes de IA;
- ideal para plugins como WhatsApp Web, overlay e integrações web.

Templates futuros podem incluir:

- Go;
- Python;
- Rust;
- WASM.

## Nome do repositório

```txt
AuxiTalk/plugin-template
```

## Objetivos do template

O template deve fornecer:

- estrutura mínima de plugin;
- manifesto `plugin.json` válido;
- helper JSON-RPC;
- lifecycle básico (`handshake`, `start`, `stop`, `health`);
- emissão de eventos;
- resposta a chamadas de capability;
- logs via stderr;
- README bilíngue;
- exemplos de configuração;
- testes básicos;
- comandos de build/dev/test;
- guia para renomear o template.

## Não objetivos iniciais

O template não deve começar com:

- framework pesado;
- UI;
- banco de dados;
- bundler complexo;
- dependências desnecessárias;
- automação real de apps;
- envio real de mensagens.

Ele deve ser uma fundação limpa.

## Estrutura proposta

```txt
plugin-template/
  README.md
  LICENSE
  package.json
  tsconfig.json
  plugin.json
  .gitignore
  .env.example

  src/
    index.ts
    runtime.ts
    rpc.ts
    logger.ts
    manifest.ts
    capabilities.ts
    events.ts
    types.ts

  test/
    rpc.test.ts
    manifest.test.ts
    runtime.test.ts

  docs/
    development.md
    protocol.md
    publishing.md
    security.md
```

## Arquivos principais

### `plugin.json`

Manifesto padrão editável pelo autor do plugin.

Exemplo inicial:

```json
{
  "id": "my-plugin",
  "name": "My Plugin",
  "version": "0.1.0",
  "runtime": "node",
  "entry": "dist/index.js",
  "kind": "tool",
  "permissions": [
    "event.emit"
  ],
  "capabilities": [
    {
      "name": "example.ping"
    }
  ]
}
```

### `src/index.ts`

Entrada principal.

Responsabilidades:

- carregar manifesto;
- iniciar runtime;
- conectar stdin/stdout;
- registrar handlers;
- enviar logs para stderr.

### `src/rpc.ts`

Helper JSON-RPC mínimo.

Responsabilidades:

- ler mensagens linha-a-linha de stdin;
- escrever mensagens linha-a-linha em stdout;
- validar `jsonrpc: "2.0"`;
- enviar response/result/error;
- rotear métodos recebidos.

### `src/runtime.ts`

Runtime local do plugin.

Responsabilidades:

- responder `plugin.handshake`;
- responder `plugin.start`;
- responder `plugin.stop`;
- responder `plugin.health`;
- lidar com `capability.call`;
- expor função para emitir eventos.

### `src/logger.ts`

Logger simples via stderr.

Regra:

```txt
stdout = somente JSON-RPC
stderr = logs humanos
```

### `src/capabilities.ts`

Local para registrar capabilities do plugin.

Exemplo:

```ts
registerCapability("example.ping", async (input) => {
  return { pong: true, input };
});
```

### `src/events.ts`

Helpers para emitir eventos para o core.

Exemplo:

```ts
emitEvent({
  type: "example.event",
  source: manifest.id,
  payload: { hello: "world" }
});
```

## Métodos JSON-RPC suportados

Core -> plugin:

```txt
plugin.handshake
plugin.start
plugin.stop
plugin.health
capability.call
```

Plugin -> core:

```txt
event.emit
log.write
action.request
memory.query
memory.write
ai.complete
```

O template deve implementar o lado plugin e oferecer helpers para chamadas plugin -> core.

## Fluxo de lifecycle

```txt
core inicia processo
  -> plugin carrega manifesto
  -> core chama plugin.handshake
  -> core chama plugin.start
  -> plugin fica aguardando mensagens
  -> core chama capability.call quando necessário
  -> core chama plugin.health periodicamente
  -> core chama plugin.stop antes de encerrar
```

## Scripts previstos

```json
{
  "scripts": {
    "dev": "tsx src/index.ts",
    "build": "tsc",
    "start": "node dist/index.js",
    "test": "vitest run",
    "typecheck": "tsc --noEmit",
    "lint": "eslint ."
  }
}
```

## Dependências propostas

Produção:

```txt
nenhuma inicialmente, se possível
```

Desenvolvimento:

```txt
typescript
tsx
vitest
eslint
@types/node
```

## Segurança

O template deve ensinar:

- nunca escrever logs no stdout;
- nunca salvar tokens no repositório;
- usar `.env` apenas localmente;
- declarar permissões mínimas;
- validar entradas de capability;
- tratar timeouts e encerramento limpo;
- não executar ações sensíveis diretamente;
- usar `action.request` para ações sensíveis.

## Documentação obrigatória para plugins derivados

Todo plugin criado a partir do template deve documentar:

- o que faz;
- tipo do plugin;
- permissões usadas;
- capabilities fornecidas;
- eventos emitidos;
- configuração;
- variáveis de ambiente;
- exemplos de uso;
- riscos e limitações.

## Experiência do desenvolvedor

O fluxo ideal deve ser:

```sh
git clone https://github.com/AuxiTalk/plugin-template my-plugin
cd my-plugin
npm install
npm run dev
npm test
npm run build
```

Depois editar:

- `plugin.json`;
- `package.json`;
- `README.md`;
- handlers em `src/capabilities.ts`.

## Critérios de aceite do template

O repositório `plugin-template` estará pronto quando:

- `npm install` funcionar;
- `npm test` passar;
- `npm run build` gerar `dist/`;
- `plugin.json` for válido para o core;
- plugin responder `plugin.handshake`;
- plugin responder `plugin.health`;
- plugin executar uma capability exemplo;
- README explicar como criar um novo plugin;
- docs explicarem protocolo, segurança e publicação.

## Roadmap do template

### Fase T1 — Fundação

- Criar repo `AuxiTalk/plugin-template`.
- Adicionar TypeScript básico.
- Adicionar `plugin.json`.
- Adicionar README bilíngue.
- Adicionar `.env.example` e `.gitignore`.

### Fase T2 — JSON-RPC helper

- Implementar parser linha-a-linha.
- Implementar request/response/error.
- Implementar router de métodos.
- Testar mensagens inválidas.

### Fase T3 — Runtime de plugin

- Implementar handshake/start/stop/health.
- Implementar capability.call.
- Implementar exemplo `example.ping`.

### Fase T4 — DX e documentação

- Adicionar testes.
- Adicionar docs.
- Adicionar guia de publicação.
- Adicionar guia de segurança.

### Fase T5 — Integração com core

- Testar com `auxitalk-core`.
- Garantir compatibilidade com `plugin.json`.
- Documentar como registrar no core.

---

## English (EN)

## Goal

Create an official `AuxiTalk/plugin-template` repository as a simple, safe, and well-documented foundation for building AuxiTalk plugins.

The template should be easy for humans and AI agents to understand, copy, modify, and publish.

## Initial language

The first official template will use **TypeScript/Node.js**.

Reasons:

- large community;
- good browser, desktop, and automation integration;
- easy distribution;
- good support for JSON-RPC/stdin/stdout;
- AI-agent friendly;
- ideal for plugins such as WhatsApp Web, overlays, and web integrations.

Future templates may include:

- Go;
- Python;
- Rust;
- WASM.

## Repository name

```txt
AuxiTalk/plugin-template
```

## Template goals

The template should provide:

- minimal plugin structure;
- valid `plugin.json` manifest;
- JSON-RPC helper;
- basic lifecycle (`handshake`, `start`, `stop`, `health`);
- event emission;
- capability call handling;
- stderr logging;
- bilingual README;
- configuration examples;
- basic tests;
- build/dev/test commands;
- guide for renaming the template.

## Non-goals

The initial template should not include:

- heavy framework;
- UI;
- database;
- complex bundler;
- unnecessary dependencies;
- real app automation;
- real message sending.

It should remain a clean foundation.

## Proposed structure

```txt
plugin-template/
  README.md
  LICENSE
  package.json
  tsconfig.json
  plugin.json
  .gitignore
  .env.example

  src/
    index.ts
    runtime.ts
    rpc.ts
    logger.ts
    manifest.ts
    capabilities.ts
    events.ts
    types.ts

  test/
    rpc.test.ts
    manifest.test.ts
    runtime.test.ts

  docs/
    development.md
    protocol.md
    publishing.md
    security.md
```

## Main files

### `plugin.json`

Default manifest edited by plugin authors.

### `src/index.ts`

Main entrypoint.

Responsibilities:

- load manifest;
- start runtime;
- connect stdin/stdout;
- register handlers;
- send logs to stderr.

### `src/rpc.ts`

Minimal JSON-RPC helper.

Responsibilities:

- read line-delimited messages from stdin;
- write line-delimited messages to stdout;
- validate `jsonrpc: "2.0"`;
- send response/result/error;
- route received methods.

### `src/runtime.ts`

Local plugin runtime.

Responsibilities:

- respond to `plugin.handshake`;
- respond to `plugin.start`;
- respond to `plugin.stop`;
- respond to `plugin.health`;
- handle `capability.call`;
- expose event emission helpers.

## Supported JSON-RPC methods

Core -> plugin:

```txt
plugin.handshake
plugin.start
plugin.stop
plugin.health
capability.call
```

Plugin -> core:

```txt
event.emit
log.write
action.request
memory.query
memory.write
ai.complete
```

## Lifecycle flow

```txt
core starts process
  -> plugin loads manifest
  -> core calls plugin.handshake
  -> core calls plugin.start
  -> plugin waits for messages
  -> core calls capability.call when needed
  -> core calls plugin.health periodically
  -> core calls plugin.stop before shutdown
```

## Planned scripts

```json
{
  "scripts": {
    "dev": "tsx src/index.ts",
    "build": "tsc",
    "start": "node dist/index.js",
    "test": "vitest run",
    "typecheck": "tsc --noEmit",
    "lint": "eslint ."
  }
}
```

## Security

The template should teach plugin authors to:

- never write logs to stdout;
- never store tokens in the repository;
- use `.env` only locally;
- declare minimal permissions;
- validate capability inputs;
- handle timeouts and clean shutdown;
- never execute sensitive actions directly;
- use `action.request` for sensitive actions.

## Acceptance criteria

The `plugin-template` repository is ready when:

- `npm install` works;
- `npm test` passes;
- `npm run build` generates `dist/`;
- `plugin.json` is valid for the core;
- plugin responds to `plugin.handshake`;
- plugin responds to `plugin.health`;
- plugin executes an example capability;
- README explains how to create a new plugin;
- docs explain protocol, security, and publishing.
