# AuxiTalk Plugin Dashboard — Planning Document

> Planejamento do plugin de dashboard visual oficial do AuxiTalk.

---

## 1. Objetivo

Criar um plugin oficial de dashboard web que permita visualizar e controlar o AuxiTalk Core de forma simples e visual.

O dashboard deve ser:

- Um **plugin `ui`** plugável.
- Um plugin **especial** que vem junto do core para facilitar o uso inicial.
- Acessível via navegador em qualquer dispositivo da rede local.
- Leve, rápido e simples de manter.

---

## 2. Visão de arquitetura

### Modelo atual

- Core se comunica com plugins apenas via **JSON-RPC sobre stdio**.
- Dashboard será um plugin que também usa JSON-RPC.

### Futuro (opcional)

- Criar um plugin `plugin-http-api` que expõe o core via REST.
- Dashboard poderá evoluir para consumir REST futuramente.

Por enquanto, o dashboard vai funcionar usando **JSON-RPC**.

---

## 3. Escopo da primeira versão

### Funcionalidades iniciais (MVP)

- Ver status do core (modo atual, versão, tempo rodando)
- Listar plugins carregados (id, kind, status, capabilities)
- Ver sessões ativas (quantidade, participantes)
- Ver eventos recentes do event bus
- Ver sugestões geradas recentemente
- Aprovar ou rejeitar ações de risco alto/medio
- Ver logs básicos do core

### Funcionalidades fora do escopo inicial

- Editar configurações em tempo real
- Gerenciar memória
- Iniciar/parar/reiniciar plugins
- Visualização avançada de eventos
- Autenticação e usuários
- Temas e customização
- Mobile-first avançado

---

## 4. Stack técnica

- **Linguagem**: Go
- **Templates**: Templ
- **Interatividade**: HTMX
- **Estilo**: CSS minimalista ou Tailwind via CDN (decidir)
- **Servidor HTTP**: `net/http` padrão

### Justificativa

- Go: alinhado com o core e consumo de recursos.
- Templ: templates compilados, type-safe, excelente DX.
- HTMX: permite interatividade sem escrever muito JavaScript.
- Simplicidade: manter o dashboard leve.

---

## 5. Modelo de execução

- O dashboard vai subir um servidor HTTP local (ex: `localhost:8080`).
- Ele será um **plugin especial** que inicia junto com o core.
- O usuário poderá desabilitar ou trocar por outro dashboard futuramente.

### Comunicação

- Dashboard usa JSON-RPC para falar com o core.
- Core envia eventos para o dashboard via JSON-RPC (`event.emit`).

---

## 6. Responsabilidades do dashboard

### Deve fazer

- Subir servidor HTTP
- Renderizar páginas com Templ + HTMX
- Consumir JSON-RPC do core
- Exibir informações em tempo real (via polling ou eventos)
- Permitir ações básicas (aprovar/rejeitar ações)

### Não deve fazer (primeira versão)

- Expor API REST pública
- Gerenciar usuários e autenticação
- Editar configurações do core
- Controlar ciclo de vida de plugins

---

## 7. Estrutura proposta do repositório

```txt
AuxiTalk/plugin-dashboard/
  README.md
  README.pt-BR.md
  plugin.json
  go.mod

  cmd/
    dashboard/
      main.go

  internal/
    server/
      routes.go
      handlers.go
    templates/
      (arquivos .templ)
    rpc/
      client.go
    state/
      store.go

  assets/
    css/
    js/ (se necessário)

  docs/
    usage.md
    development.md
```

---

## 8. Fluxo de dados

```txt
Core
  -> envia eventos via JSON-RPC
  -> responde chamadas de capability

Dashboard
  -> consome eventos
  -> faz polling ou escuta eventos
  -> renderiza páginas com Templ + HTMX
  -> envia ações de volta via JSON-RPC
```

---

## 9. Segurança

- O dashboard deve rodar apenas em localhost por padrão.
- Pode ser configurado para escutar em uma interface específica.
- Não deve expor credenciais ou permitir ações destrutivas sem confirmação.
- Logs via stderr.

---

## 10. Critérios de aceite da primeira versão

- Dashboard sobe servidor HTTP
- Exibe status do core
- Lista plugins carregados
- Lista sessões ativas
- Mostra eventos recentes
- Permite aprovar/rejeitar ações
- Documentação básica

---

## 11. Roadmap futuro

- v2: controle de plugins (start/stop/restart)
- v2: edição de configurações
- v3: autenticação básica
- v3: API REST via plugin separado
- v4: temas, mobile, notificações

---

## 12. Próximos passos

1. Aprovar este planejamento
2. Criar repositório `AuxiTalk/plugin-dashboard`
3. Inicializar Go + Templ
4. Criar estrutura básica do servidor HTTP
5. Implementar páginas iniciais
6. Integrar com JSON-RPC do core
7. Adicionar HTMX para interatividade
8. Testar fluxo completo

---

**Status:** Planejamento em andamento. Aguardando validação.
