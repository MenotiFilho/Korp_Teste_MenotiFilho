# ms-faturamento - Contratos de API (handoff para frontend)

Este documento consolida os contratos HTTP do `ms-faturamento`.

Base URL local:

- `http://localhost:8082`

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

## 2) Criar nota fiscal

### Request

- `POST /api/v1/notas`
- `Content-Type: application/json`

Body:

```json
{
  "itens": [
    {
      "produto_codigo": "P-T13",
      "quantidade": 4
    }
  ]
}
```

### Response (201)

```json
{
  "id": 3,
  "numero": 15,
  "status": "ABERTA",
  "itens": [
    {
      "id": 3,
      "produto_codigo": "P-T13",
      "quantidade": 4
    }
  ]
}
```

### Erros comuns

- `400 INVALID_JSON`
- `400 VALIDATION_ERROR`
- `500 INTERNAL_ERROR`

## 3) Listar notas fiscais

### Request

- `GET /api/v1/notas`

### Response (200)

```json
[
  {
    "id": 4,
    "numero": 16,
    "status": "ABERTA",
    "itens": [
      {
        "id": 4,
        "produto_codigo": "P-T13",
        "quantidade": 2
      }
    ]
  }
]
```

## 4) Últimas notas (novo)

### Request

- `GET /api/v1/notas/ultimas`
- optional query params: `limit` (default 6)

### Response (200)

```json
[
  {
    "id": 4,
    "numero": 16,
    "status": "ABERTA",
    "itens": [
      {
        "id": 4,
        "produto_codigo": "P-T13",
        "quantidade": 2
      }
    ]
  }
]
```

Regras:
- retorna notas ordenadas por `numero` desc (mais recente primeiro).
- `limit` default = 6, max = 100

## 4) Imprimir nota fiscal

### Request

- `POST /api/v1/notas/{id}/imprimir`

Exemplo:

- `POST /api/v1/notas/4/imprimir`

### Response sucesso (200)

- body vazio

### Regras de negocio

- so permite impressao de nota com status `ABERTA`
- usa chave de idempotencia por nota no fluxo de impressao (`invoice-print-{id}`)
- ao imprimir com sucesso:
  - chama `ms-estoque` para baixar saldo
  - atualiza status da nota para `FECHADA`
- se o estoque falhar:
  - retorna erro para o frontend
  - nota permanece `ABERTA`

### Erros de negocio

- `400 VALIDATION_ERROR` (id invalido)
- `404 INVOICE_NOT_FOUND`
- `409 INVOICE_NOT_ABERTA`
- `404 PRODUCT_NOT_FOUND_IN_STOCK`
- `409 INSUFFICIENT_STOCK`

### Erros de integracao

- `503 ESTOQUE_UNAVAILABLE`
- `500 INTERNAL_ERROR`

## 5) Atualizar nota fiscal

### Request

- `PUT /api/v1/notas/{id}`
- `Content-Type: application/json`

Body:

```json
{
  "itens": [
    {
      "produto_codigo": "P-T13",
      "quantidade": 6
    },
    {
      "produto_codigo": "P-001",
      "quantidade": 2
    }
  ]
}
```

### Response (200)

```json
{
  "id": 3,
  "numero": 15,
  "status": "ABERTA",
  "itens": [
    {
      "id": 10,
      "produto_codigo": "P-T13",
      "quantidade": 6
    },
    {
      "id": 11,
      "produto_codigo": "P-001",
      "quantidade": 2
    }
  ]
}
```

### Regras de negocio

- so permite atualizacao de nota com status `ABERTA`
- os itens antigos sao substituidos pelos novos
- nota com status `FECHADA` retorna erro `409`

### Erros comuns

- `400 VALIDATION_ERROR` (id invalido)
- `400 INVALID_JSON`
- `400 VALIDATION_ERROR` (itens vazios ou invalidos)
- `404 INVOICE_NOT_FOUND`
- `409 INVOICE_NOT_ABERTA`
- `500 INTERNAL_ERROR`

## 6) Deletar nota fiscal (soft delete)

### Request

- `DELETE /api/v1/notas/{id}`

### Response (204)

- body vazio

### Regras de negocio

- so permite exclusao de nota com status `ABERTA`
- nota com status `FECHADA` retorna erro `409`
- soft delete: o registro e marcado com `deleted_at`, nao removido fisicamente
- notas deletadas nao aparecem em `GET /api/v1/notas`

### Erros comuns

- `400 VALIDATION_ERROR` (id invalido)
- `404 INVOICE_NOT_FOUND`
- `409 INVOICE_NOT_ABERTA`
- `500 INTERNAL_ERROR`

## Fluxo de falha obrigatorio (demonstrado)

Cenario validado na Task 13:

1. nota `ABERTA` criada no faturamento
2. `ms-estoque` interrompido
3. tentativa de imprimir retorna `503 ESTOQUE_UNAVAILABLE`
4. nota permanece `ABERTA`
5. `ms-estoque` religado
6. nova tentativa de imprimir retorna `200`
7. nota passa para `FECHADA`

## Comandos de validacao usados na Task 13

```bash
docker compose -f infra/docker-compose.yml up -d --build db-estoque ms-estoque db-faturamento ms-faturamento
DB_URL="postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable" make -C apps/ms-estoque migrate-up
DB_URL="postgres://postgres:postgres@localhost:5434/faturamento?sslmode=disable" make -C apps/ms-faturamento migrate-up
go test ./... -C apps/ms-faturamento
go test ./... -short -C apps/ms-faturamento
curl -s -X POST http://localhost:8081/api/v1/produtos -H "Content-Type: application/json" -d '{"codigo":"P-T13","descricao":"Produto task13","saldo":10}'
curl -s -X POST http://localhost:8082/api/v1/notas -H "Content-Type: application/json" -d '{"itens":[{"produto_codigo":"P-T13","quantidade":4}]}'
curl -s -X POST http://localhost:8082/api/v1/notas/3/imprimir -w "\n%{http_code}"
curl -s -X POST http://localhost:8082/api/v1/notas -H "Content-Type: application/json" -d '{"itens":[{"produto_codigo":"P-T13","quantidade":2}]}'
docker stop korp-ms-estoque
curl -s -X POST http://localhost:8082/api/v1/notas/4/imprimir -w "\n%{http_code}"
docker start korp-ms-estoque
curl -s -X POST http://localhost:8082/api/v1/notas/4/imprimir -w "\n%{http_code}"
docker compose -f infra/docker-compose.yml down
```
