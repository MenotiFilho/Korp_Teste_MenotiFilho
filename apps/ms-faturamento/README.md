# ms-faturamento

Microservico de faturamento em Go (`net/http`) com persistencia PostgreSQL.

## Banco local para desenvolvimento

Suba o banco dedicado do faturamento:

```bash
docker compose -f ../../infra/docker-compose.yml up -d db-faturamento
```

Subir banco + ms-faturamento via Docker Compose:

```bash
docker compose -f ../../infra/docker-compose.yml up -d --build db-faturamento ms-faturamento
```

## Migrations

Por padrao, o `DB_URL` aponta para:

`postgres://postgres:postgres@localhost:5434/faturamento?sslmode=disable`

Comandos:

```bash
make migrate-up
make migrate-version
make migrate-down
make migrate-force VERSION=1
```

Criar nova migration:

```bash
make migrate-create NAME=nome_da_migration
```

## Testes

Testes completos (inclui integracao de repository por padrao):

```bash
go test ./...
```

Testes apenas unitarios (pula integracao):

```bash
go test ./... -short
```

## Executar API local

Com banco em Docker e migrations aplicadas:

```bash
DB_URL="postgres://postgres:postgres@localhost:5434/faturamento?sslmode=disable" go run ./cmd/ms-faturamento
```

Parar ambiente Docker:

```bash
docker compose -f ../../infra/docker-compose.yml down
```

## Variaveis opcionais de hardening HTTP

- `HTTP_READ_HEADER_TIMEOUT_SEC` (default: `5`)
- `HTTP_MAX_HEADER_BYTES` (default: `1048576`)
- `HTTP_MAX_BODY_BYTES` (default: `1048576`)
- `HTTP_READ_TIMEOUT_SEC` (default: `10`)
- `HTTP_WRITE_TIMEOUT_SEC` (default: `10`)
- `HTTP_IDLE_TIMEOUT_SEC` (default: `30`)

## Variaveis de integracao

- `ESTOQUE_URL` (default local: `http://localhost:8081`) - URL base do ms-estoque

Observacao em Docker Compose:

- usar `http://ms-estoque:8081` para comunicacao entre containers na mesma rede

## Endpoints disponiveis

- `GET /health`
- `POST /api/v1/notas`
- `GET /api/v1/notas`
- `POST /api/v1/notas/{id}/imprimir`
