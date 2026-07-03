# Ambiente de Desenvolvimento AuxiTalk

Este é o setup oficial para testar o AuxiTalk Core + Plugins de forma visual e completa.

## Pré-requisitos

- Go 1.22+
- `templ` CLI instalado (`go install github.com/a-h/templ/cmd/templ@latest`)
- (opcional) Node.js se quiser testar plugins TypeScript

## Como rodar tudo

1. No diretório `auxitalk-core`:

```bash
./dev.sh build      # compila core + plugins
./dev.sh run        # sobe core + dashboard juntos
```

2. Acesse o dashboard:

```
http://localhost:8080
```

## O que acontece ao rodar `dev.sh run`

- Compila o core (`auxitalkd`)
- Compila o dashboard
- Compila o plugin OpenAI
- Inicia o core com a config `configs/auxitalk.dev.json`
- Inicia o dashboard em `http://localhost:8080`
- Cria SQLite em `auxitalk-dev.db`

## Configuração de desenvolvimento

O arquivo `configs/auxitalk.dev.json` já contém:

- modo `dev`
- persistência SQLite ativada
- plugins OpenAI e WhatsApp (WhatsApp desabilitado por padrão)
- workflow de demonstração

## Comandos úteis

```bash
./dev.sh run          # sobe tudo
./dev.sh core         # só o core
./dev.sh dashboard    # só o dashboard
./dev.sh stop         # mata tudo
./dev.sh clean        # remove builds
```

## Visualização

No dashboard você verá:

- Plugins carregados
- Eventos em tempo real
- Ações pendentes / aprovadas / executadas
- Workflows ativos

Isso permite testar o fluxo completo:
Evento → Workflow → Ação → Aprovação (ou auto-execução no modo dev) → Execução → Resultado no dashboard.
