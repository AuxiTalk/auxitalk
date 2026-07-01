# AuxiTalk Core

AuxiTalk Core é um runtime modular para assistência em conversas.

Ele observa eventos de conversa, mantém estado normalizado de sessões, constrói contexto compacto, roteia capabilities de plugins e coordena ações aprovadas pelo usuário por meio de um sistema de plugins agnóstico de linguagem.

> English documentation: [README.md](README.md)

## O que é o AuxiTalk

AuxiTalk é a camada de orquestração para um ecossistema futuro de assistência em conversas.

Ele foi projetado para suportar fluxos como:

- observar uma conversa em canais como WhatsApp Web, Telegram, Discord, navegador, email ou outra fonte;
- entender contexto e tom atual;
- sugerir respostas;
- recomendar se deve responder, esperar, ignorar ou pedir mais contexto;
- exibir sugestões por UI, overlay, CLI ou outro plugin de saída;
- executar ações sensíveis apenas por gates controlados pelo usuário;
- aprender com feedback e futuros plugins de memória.

O AuxiTalk Core não é preso a um app, provedor de IA, banco de dados ou UI específica.

## Objetivos de design

- Runtime leve em Go.
- Arquitetura modular.
- Plugins em múltiplas linguagens.
- JSON-RPC 2.0 sobre stdio como protocolo inicial.
- Contratos claros para eventos, sessões, mensagens, sugestões, ações, manifests e capabilities.
- Segurança progressiva para ações sensíveis.
- Desenvolvimento guiado por documentação.
- Fácil de entender por humanos e agentes de IA.

## Arquitetura

```txt
Input Plugin
  -> Event Bus
  -> Session Manager
  -> Context Builder
  -> Action Gate / Policy
  -> Capability Router
  -> AI / Memory / Tool Plugins
  -> Suggestion Event
  -> UI / Output Plugin
  -> User Feedback
  -> Memory Update
```

## Módulos atuais

| Módulo | Caminho | Função |
| --- | --- | --- |
| Runtime | `internal/runtime` | Entrada do ciclo de vida do runtime |
| Config | `internal/config` | Configuração, defaults e modos |
| Events | `internal/events` | Event bus interno pub/sub |
| Plugin registry | `internal/plugins` | Manifest, registry e supervisor inicial |
| JSON-RPC codec | `internal/rpc` | Codec JSON-RPC stdio linha-a-linha |
| Capabilities | `internal/capabilities` | Registro e roteamento de capabilities |
| Sessions | `internal/sessions` | Estado normalizado de conversas |
| Context | `internal/context` | Construção de contexto compacto |
| Actions | `internal/actions` | Action Gate por risco/modo |
| Types | `pkg/types` | Contratos públicos de domínio |
| Protocol | `pkg/protocol` | Tipos públicos JSON-RPC |

## Modos de runtime

```txt
dev      modo permissivo para testes locais rápidos
local    modo local mais seguro com proteção para ações sensíveis
strict   modo restrito com validação mais forte
```

## Sistema de plugins

Plugins rodam como processos externos e se comunicam com o core usando JSON-RPC 2.0 sobre stdio.

Isso permite plugins em Go, TypeScript/Node.js, Python, Rust ou qualquer linguagem que leia stdin e escreva stdout.

Regras importantes:

- `stdout` é reservado para mensagens JSON-RPC.
- `stderr` é reservado para logs humanos.
- cada mensagem JSON-RPC deve ocupar uma linha.
- manifests declaram permissões e capabilities.
- ações sensíveis devem usar o fluxo de action request.

### Exemplo de manifest

```json
{
  "id": "mock-ai",
  "name": "Mock AI",
  "version": "0.1.0",
  "runtime": "node",
  "entry": "index.js",
  "kind": "ai",
  "permissions": [],
  "capabilities": [
    {
      "name": "ai.complete"
    }
  ]
}
```

## Estrutura do repositório

```txt
cmd/                 binários: auxitalkd e auxitalkctl
configs/             configurações de exemplo
docs/                arquitetura, roadmap, decisões e docs de plugins
examples/            exemplos de fluxo
internal/            pacotes privados do core
pkg/                 tipos e protocolo públicos
plugins/             plugins mínimos de exemplo
FINAL_STATUS.md      relatório da fundação inicial
```

## Começando

### Requisitos

- Go instalado.

### Rodar daemon

```sh
go run ./cmd/auxitalkd
```

### Rodar com config

```sh
go run ./cmd/auxitalkd --config configs/auxitalk.example.json
```

### Rodar CLI placeholder

```sh
go run ./cmd/auxitalkctl
```

### Rodar testes

```sh
go test ./...
```

## Status atual

A fundação inicial está completa.

Implementado:

- sistema de configuração;
- tipos centrais;
- event bus interno;
- registry de manifests de plugins;
- protocolo JSON-RPC e codec stdio;
- esqueleto de supervisor de plugins;
- capability router;
- session manager;
- context builder;
- action gate;
- manifests de plugins mock;
- documentação do primeiro loop completo.

## Próximos passos

Trabalhos recomendados:

1. conectar o supervisor a chamadas JSON-RPC reais;
2. tornar o loop mock executável via `auxitalkd`;
3. criar repositórios oficiais de plugins;
4. criar o primeiro plugin de IA;
5. criar o primeiro plugin de memória;
6. criar o primeiro plugin de UI/overlay;
7. criar o primeiro plugin real de entrada de conversa.

## Repositórios relacionados

- `AuxiTalk/auxitalk` — core runtime.
- `AuxiTalk/plugin-template` — templates oficiais e documentação de criação de plugins.

## Documentação

Docs importantes:

- `docs/architecture/core.md`
- `docs/architecture/core-types.md`
- `docs/architecture/configuration.md`
- `docs/architecture/event-bus.md`
- `docs/architecture/capability-router.md`
- `docs/architecture/session-context.md`
- `docs/architecture/action-gate.md`
- `docs/plugins/authoring-guide.md`
- `docs/plugins/protocol-draft.md`
- `docs/plugins/system.md`
- `docs/roadmap/initial-implementation-plan.md`
- `docs/decisions/0001-go-core-jsonrpc-plugins.md`

## Licença

Licença a definir.
