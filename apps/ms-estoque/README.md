# ms-estoque

Microservico de estoque em Go (`net/http`) com persistencia PostgreSQL.

## Banco local para desenvolvimento

Suba o banco dedicado do estoque:

```bash
docker compose -f ../../infra/docker-compose.yml up -d db-estoque
```

Subir banco + ms-estoque via Docker Compose:

```bash
docker compose -f ../../infra/docker-compose.yml up -d --build db-estoque ms-estoque
```

## Migrations

Por padrao, o `DB_URL` aponta para:

`postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable`

Comandos:

```bash
make migrate-up
make migrate-version
make migrate-down
```

Criar nova migration:

```bash
make migrate-create NAME=nome_da_migration
```

Forcar versao (uso de recuperacao):

```bash
make migrate-force VERSION=1
```

## Testes

Testes completos (inclui integracao de repository por padrao):

```bash
go test ./...
```

Observacoes:

- os testes de integracao usam por padrao `postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable`
- garanta que o Postgres esteja no ar e as migrations aplicadas antes de rodar

Testes apenas unitarios (pula integracao):

```bash
go test ./... -short
```

Testes de integracao do repository com URL customizada:

```bash
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable" go test ./internal/repository -v
```

## Executar API local

Com banco em Docker e migrations aplicadas:

```bash
DB_URL="postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable" go run ./cmd/ms-estoque
```

Parar ambiente Docker:

```bash
docker compose -f ../../infra/docker-compose.yml down
```

Variaveis opcionais de hardening HTTP:

- `HTTP_READ_HEADER_TIMEOUT_SEC` (default: `5`)
- `HTTP_MAX_HEADER_BYTES` (default: `1048576`)
- `HTTP_MAX_BODY_BYTES` (default: `1048576`)
- `HTTP_READ_TIMEOUT_SEC` (default: `10`)
- `HTTP_WRITE_TIMEOUT_SEC` (default: `10`)
- `HTTP_IDLE_TIMEOUT_SEC` (default: `30`)

Endpoints disponiveis:

- `GET /health`
- `POST /api/v1/produtos`
- `GET /api/v1/produtos`
- `POST /api/v1/estoque/baixa`

Exemplo de criacao de produto:

```bash
curl -i -X POST "http://localhost:8081/api/v1/produtos" \
  -H "Content-Type: application/json" \
  -d '{"codigo":"P-001","descricao":"Produto 1","saldo":10}'
```

Exemplo de listagem de produtos:

```bash
curl -i "http://localhost:8081/api/v1/produtos"
```

Exemplo de baixa de estoque:

```bash
curl -i -X POST "http://localhost:8081/api/v1/estoque/baixa" \
  -H "Content-Type: application/json" \
  -d '{"itens":[{"codigo":"P-001","quantidade":2}]}'
```
