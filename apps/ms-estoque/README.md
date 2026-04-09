# ms-estoque

Microservico de estoque em Go (`net/http`) com persistencia PostgreSQL.

## Banco local para desenvolvimento

Suba o banco dedicado do estoque:

```bash
docker compose -f ../../infra/docker-compose.yml up -d db-estoque
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
