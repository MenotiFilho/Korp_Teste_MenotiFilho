# ms-estoque - Contratos de API (handoff para ms-faturamento)

Este documento consolida os contratos HTTP atualmente disponiveis no `ms-estoque`.

Base URL local:

- `http://localhost:8081`

Prefixo de API:

- `/api/v1`

## Padrao de erro

Todos os erros retornam JSON no formato:

```json
{
  "code": "CODIGO_DO_ERRO",
  "message": "mensagem amigavel",
  "details": null,
  "request_id": "id-da-requisicao"
}
```

## 1) Healthcheck

### Request

- `GET /health`

### Response (200)

- body vazio

## 2) Criar produto

### Request

- `POST /api/v1/produtos`
- `Content-Type: application/json`

Body:

```json
{
  "codigo": "P-100",
  "descricao": "Produto 100",
  "saldo": 10
}
```

### Response (201)

```json
{
  "id": 3,
  "codigo": "P-100",
  "descricao": "Produto 100",
  "saldo": 10
}
```

### Erros comuns

- `400 INVALID_JSON`
- `400 VALIDATION_ERROR`
- `409 PRODUCT_CODIGO_ALREADY_EXISTS`
- `413 PAYLOAD_TOO_LARGE`
- `500 INTERNAL_ERROR`

## 3) Listar produtos

### Request

- `GET /api/v1/produtos`

### Response (200)

```json
[
  {
    "id": 1,
    "codigo": "P-001",
    "descricao": "Produto 1",
    "saldo": 1
  },
  {
    "id": 3,
    "codigo": "P-100",
    "descricao": "Produto 100",
    "saldo": 7
  }
]
```

Observacao:

- retorno ordenado por `id` ascendente.

## 4) Baixa de estoque

### Request

- `POST /api/v1/estoque/baixa`
- `Content-Type: application/json`
- `Idempotency-Key: <chave-unica-obrigatoria>`

Body:

```json
{
  "itens": [
    {
      "codigo": "P-100",
      "quantidade": 3
    }
  ]
}
```

### Response (200)

- body vazio

### Erros de negocio

- `400 VALIDATION_ERROR` (header `Idempotency-Key` ausente)
- `400 VALIDATION_ERROR` (item invalido)
- `404 PRODUCT_NOT_FOUND`
- `409 INSUFFICIENT_STOCK`

### Erros tecnicos

- `400 INVALID_JSON`
- `413 PAYLOAD_TOO_LARGE`
- `500 INTERNAL_ERROR`

## Regras de negocio relevantes para integracao

- A baixa e transacional no banco.
- O servico bloqueia a linha do produto durante baixa (`FOR UPDATE`).
- O servico ordena os itens por `codigo` antes de lock para reduzir risco de deadlock em baixa concorrente.
- Se qualquer item falhar, a transacao e revertida.
- Nao e permitido saldo negativo.
- `Idempotency-Key` e obrigatorio no endpoint de baixa.
- Repeticao da mesma baixa com a mesma `Idempotency-Key` e no-op (nao debita novamente).
- Produtos deletados (soft delete) nao aparecem nas listagens.

## 5) Atualizar produto

### Request

- `PUT /api/v1/produtos/{id}`
- `Content-Type: application/json`

Body:

```json
{
  "descricao": "Produto 100 Atualizado",
  "saldo": 25
}
```

Observacao:

- apenas `descricao` e `saldo` podem ser atualizados.
- `codigo` e imutavel.

### Response (200)

```json
{
  "id": 3,
  "codigo": "P-100",
  "descricao": "Produto 100 Atualizado",
  "saldo": 25
}
```

### Erros comuns

- `400 VALIDATION_ERROR` (id invalido)
- `400 INVALID_JSON`
- `400 VALIDATION_ERROR` (descricao vazia ou saldo negativo)
- `404 PRODUCT_NOT_FOUND`
- `413 PAYLOAD_TOO_LARGE`
- `500 INTERNAL_ERROR`

## 6) Deletar produto (soft delete)

### Request

- `DELETE /api/v1/produtos/{id}`

### Response (204)

- body vazio

### Erros comuns

- `400 VALIDATION_ERROR` (id invalido)
- `404 PRODUCT_NOT_FOUND`
- `500 INTERNAL_ERROR`

### Observacao

- soft delete: o registro e marcado com `deleted_at`, nao removido fisicamente.
- produtos deletados nao aparecem em `GET /api/v1/produtos`.

## Comandos de validacao usados na Task 12

```bash
docker compose -f infra/docker-compose.yml up -d db-estoque ms-estoque
DB_URL="postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable" make -C apps/ms-estoque migrate-up
go test -C apps/ms-estoque ./...
go test -C apps/ms-estoque ./... -short
curl -i http://localhost:8081/health
curl -i -X POST http://localhost:8081/api/v1/produtos -H "Content-Type: application/json" -d '{"codigo":"P-100","descricao":"Produto 100","saldo":10}'
curl -i http://localhost:8081/api/v1/produtos
curl -i -X POST http://localhost:8081/api/v1/estoque/baixa -H "Content-Type: application/json" -H "Idempotency-Key: baixa-p100-001" -d '{"itens":[{"codigo":"P-100","quantidade":3}]}'
curl -i -X POST http://localhost:8081/api/v1/estoque/baixa -H "Content-Type: application/json" -H "Idempotency-Key: baixa-p100-001" -d '{"itens":[{"codigo":"P-100","quantidade":3}]}'
curl -i http://localhost:8081/api/v1/produtos
docker compose -f infra/docker-compose.yml down
```
