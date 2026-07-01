# AuxiTalk Plugin OpenAI — Planning Document

> Plano de arquitetura, escopo, responsabilidades e requisitos do plugin oficial OpenAI-compatible.

---

## 1. Objetivo do plugin

Criar um plugin oficial que permita ao AuxiTalk Core chamar qualquer API compatível com OpenAI Chat Completions.

Objetivo principal:

- Expor a capability `ai.complete` de forma confiável, segura e configurável.
- Permitir que o AuxiTalk sugira respostas usando modelos OpenAI, Groq, Ollama, Together, LM Studio, ou qualquer endpoint compatível.

---

## 2. Escopo do plugin

### In scope (primeira versão)

- Implementar `ai.complete`
- Suporte a APIs OpenAI-compatible (baseURL customizável)
- Configuração por variáveis de ambiente
- Suporte a `model`, `temperature`, `max_tokens`
- Timeout por requisição
- Logging básico via stderr
- Validação de configuração
- Documentação bilíngue
- Testes básicos

### Out of scope (futura versão)

- Suporte a function calling / tools
- Suporte a streaming
- Suporte a embeddings
- Cache de respostas
- Rate limiting interno
- Retry automático
- Multimodal (imagens, áudio)
- Fine-tuning
- Suporte a múltiplos provedores simultâneos

---

## 3. Capability principal

### `ai.complete`

Responsabilidade:

- Receber prompt + parâmetros
- Chamar a API OpenAI-compatible
- Retornar resposta do modelo

Input esperado:

```json
{
  "prompt": "string",
  "model": "string (opcional)",
  "temperature": "number (opcional)",
  "max_tokens": "number (opcional)"
}
```

Output esperado:

```json
{
  "text": "string",
  "model": "string",
  "usage": {
    "prompt_tokens": "number",
    "completion_tokens": "number",
    "total_tokens": "number"
  }
}
```

---

## 4. Responsabilidades do plugin

### Deve fazer

- Carregar configuração
- Validar API key e baseURL
- Fazer requisições HTTP
- Tratar erros da API
- Retornar resposta estruturada
- Logar erros e chamadas em stderr
- Respeitar timeout configurado

### Não deve fazer

- Executar ações sensíveis sem usar `action.request`
- Armazenar segredos no repositório
- Fazer cache de respostas
- Implementar lógica de retry complexa (primeira versão)
- Expor endpoints HTTP próprios

---

## 5. Requisitos técnicos

### Linguagem

- Go (prioridade por consumo de recursos)

### Dependências

- `github.com/sashabaranov/go-openai` (ou cliente oficial quando disponível)
- `net/http` padrão
- `context` para timeout
- `encoding/json`

### Configuração

O plugin deve suportar configuração via variáveis de ambiente:

```txt
OPENAI_API_KEY
OPENAI_BASE_URL (opcional, default https://api.openai.com/v1)
OPENAI_MODEL (opcional, default gpt-4o-mini)
OPENAI_TIMEOUT (opcional, default 60s)
OPENAI_TEMPERATURE (opcional)
OPENAI_MAX_TOKENS (opcional)
```

### Segurança

- Nunca commitar `.env` ou chaves
- Validar presença de `OPENAI_API_KEY`
- Usar `context.WithTimeout` em todas as chamadas
- Não logar API key
- Declarar permissão `ai.complete` no manifest

---

## 6. Estrutura proposta do repositório

```txt
AuxiTalk/plugin-openai/
  README.md
  README.pt-BR.md
  plugin.json
  go.mod
  go.sum
  .gitignore
  .env.example

  cmd/
    plugin/
      main.go

  internal/
    openai/
      client.go
      config.go
      types.go

  docs/
    configuration.md
    capabilities.md
    security.md

  testdata/
    (exemplos de requests/responses)
```

---

## 7. Fluxo de execução

```txt
Core chama capability "ai.complete"
  -> Plugin recebe input
  -> Plugin valida configuração
  -> Plugin monta request OpenAI
  -> Plugin chama API com timeout
  -> Plugin recebe resposta
  -> Plugin retorna resultado estruturado
  -> Core recebe resposta
```

---

## 8. Tratamento de erros

O plugin deve retornar erro quando:

- API key não configurada
- BaseURL inválida
- Timeout atingido
- API retorna erro (rate limit, invalid model, etc.)
- Resposta inválida

Erros devem ser retornados de forma clara para o core poder decidir o que fazer.

---

## 9. Documentação obrigatória

- `README.md` (inglês)
- `README.pt-BR.md` (português)
- `docs/configuration.md`
- `docs/capabilities.md`
- `docs/security.md`

---

## 10. Critérios de aceite da primeira versão

O plugin estará pronto quando:

- `go build` gerar binário funcional
- `plugin.handshake` responder corretamente
- `ai.complete` funcionar com OpenAI oficial
- `ai.complete` funcionar com Ollama (OpenAI compatible)
- Configuração por env funcionar
- Erros forem tratados e retornados
- Documentação bilíngue existir
- Testes básicos passarem

---

## 11. Roadmap futuro (v2+)

- Suporte a function calling
- Suporte a streaming
- Cache inteligente de respostas
- Retry com backoff
- Métricas de uso
- Suporte a múltiplos modelos simultâneos
- Plugin de embeddings separado

---

## 12. Próximos passos após planejamento

1. Criar repositório `AuxiTalk/plugin-openai`
2. Inicializar Go module
3. Implementar configuração
4. Implementar cliente OpenAI
5. Implementar `ai.complete`
6. Criar manifest e runtime
7. Adicionar testes básicos
8. Escrever documentação
9. Publicar primeira versão

---

## 13. Decisões tomadas

- Linguagem: **Go**
- Capability principal: `ai.complete`
- Protocolo: JSON-RPC 2.0 sobre stdio
- Configuração: variáveis de ambiente
- Primeira versão: sem streaming, sem function calling
- Foco: simplicidade, segurança e compatibilidade ampla

---

**Status:** Planejamento concluído. Aguardando início da implementação.
