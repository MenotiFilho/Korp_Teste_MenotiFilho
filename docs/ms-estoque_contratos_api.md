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
- Se qualquer item falhar, a transacao e revertida.
- Nao e permitido saldo negativo.

## Comandos de validacao usados na Task 12

```bash
docker compose -f infra/docker-compose.yml up -d db-estoque ms-estoque
DB_URL="postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable" make -C apps/ms-estoque migrate-up
go test ./... -C apps/ms-estoque
go test ./... -short -C apps/ms-estoque
curl -i http://localhost:8081/health
curl -i -X POST http://localhost:8081/api/v1/produtos -H "Content-Type: application/json" -d '{"codigo":"P-100","descricao":"Produto 100","saldo":10}'
curl -i http://localhost:8081/api/v1/produtos
curl -i -X POST http://localhost:8081/api/v1/estoque/baixa -H "Content-Type: application/json" -d '{"itens":[{"codigo":"P-100","quantidade":3}]}'
curl -i http://localhost:8081/api/v1/produtos
docker compose -f infra/docker-compose.yml down
```
