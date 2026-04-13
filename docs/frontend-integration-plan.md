# Integração Frontend ↔ Microsserviços (Plano Consolidado)

Este documento consolida o plano para substituir os mocks do frontend Angular por chamadas reais aos microsserviços `ms-estoque` e `ms-faturamento`.

Status: pronto para implementação. O middleware CORS recomendado para os serviços Go é `github.com/rs/cors`.

Resumo
- Objetivo: trocar MockDataService por serviços HTTP reais no frontend.
- Middleware CORS: usar `rs/cors` em `ms-estoque` e `ms-faturamento`.
- Endpoints base (locais):
  - ms-estoque: http://localhost:8081
  - ms-faturamento: http://localhost:8082

Passos concretos

1) Preparação
 - Garantir que os microsserviços rodem nas portas definidas e que os bancos estejam migrados.
 - Confirmar que o frontend roda em `http://localhost:4200` (ajustar conforme necessário).

2) Habilitar CORS nos microsserviços (usar rs/cors)
 - Adicionar dependência: `go get github.com/rs/cors`.
 - Inserir middleware antes do `http.Server` registrar o mux.

Exemplo mínimo (main.go) usando rs/cors:

```go
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/rs/cors"
)

func main() {
    mux := NewMux() // seu http.ServeMux já existente

    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:4200"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization", "Idempotency-Key"},
        ExposedHeaders:   []string{"Content-Type"},
        AllowCredentials: true,
    })

    handler := c.Handler(mux)

    srv := &http.Server{
        Addr:         ":8081", // ajustar por serviço
        Handler:      handler,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Println("starting server on :8081")
    log.Fatal(srv.ListenAndServe())
}
```

Notas:
 - Ajustar `Addr` para o serviço (8081 para estoque, 8082 para faturamento).
 - `AllowedOrigins` pode usar `*` em desenvolvimento, mas preferimos `http://localhost:4200` para segurança.

3) Implementar serviços HTTP no frontend (Angular)
 - Criar `api.config.ts` contendo as base URLs (ou usar `environment.ts`/`environment.prod.ts`). Recomendação: usar `environment.ts` (padrão Angular) se for buildar para ambientes.
 - Criar `ProdutoService` (chama `GET /api/v1/produtos`, `POST /api/v1/produtos`).
 - Criar `NotaService` (chama `GET /api/v1/notas`, `POST /api/v1/notas`, `POST /api/v1/notas/{id}/imprimir`).
 - Criar `ApiErrorMapper` para transformar o objeto de erro padronizado do backend ({code,message,details,request_id}) em um ErrorModel do frontend.

Exemplo de endpoints (conforme contratos):
 - GET http://localhost:8081/api/v1/produtos
 - POST http://localhost:8081/api/v1/produtos
 - GET http://localhost:8082/api/v1/notas
 - POST http://localhost:8082/api/v1/notas
 - POST http://localhost:8082/api/v1/notas/{id}/imprimir

4) Substituições incrementais no frontend
 - Fase A (read-only): substituir listas que consomem MockDataService por chamadas `list()` dos serviços.
 - Fase B (writes): substituir criação de produto/nota para usar os endpoints POST.
 - Fase C (impressão): substituir o fluxo de imprimir para chamar `NotaService.imprimir(id)` e tratar loading/erros.
 - Manter MockDataService presente como fallback até que todas referências sejam removidas; fazer remoção final quando tudo estiver integrado e testado.

5) UX e tratamento de erros
 - Mapear erro vindo do backend e exibir através do `snackbar.service.ts`.
 - Ao imprimir: exibir indicador de processamento; desabilitar botão; em erro, mostrar mensagem clara.

6) Testes recomendados
 - Unit tests com HttpTestingController para `ProdutoService` e `NotaService` (caso sucesso e caso erro padronizado).
 - Testes unitários dos componentes que exibem loading/erro.
 - Teste manual obrigatório do cenário de falha:
   1. Subir serviços e bancos.
   2. Criar nota ABERTA via frontend.
   3. Parar `ms-estoque` (docker stop container).
   4. Tentar imprimir: esperar erro 503 ou `ESTOQUE_UNAVAILABLE` e nota permanecer ABERTA.
   5. Reiniciar `ms-estoque` e tentar imprimir novamente: deve ocorrer sucesso e nota virar FECHADA.

Comandos úteis
 - Subir infra (exemplo):
   - `docker compose -f infra/docker-compose.yml up -d db-estoque db-faturamento ms-estoque ms-faturamento`
 - Migrar bancos:
   - `DB_URL="postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable" make -C apps/ms-estoque migrate-up`
   - `DB_URL="postgres://postgres:postgres@localhost:5434/faturamento?sslmode=disable" make -C apps/ms-faturamento migrate-up`
 - Rodar frontend:
   - `cd apps/frontend && npm install && npm run start` (ou `ng serve` conforme package.json)

Checklist de implementação (para assinalar ao concluir)
 - [ ] rs/cors adicionado e configurado em ms-estoque
 - [ ] rs/cors adicionado e configurado em ms-faturamento
 - [ ] ProdutoService implementado e integrado aos componentes
 - [ ] NotaService implementado e integrado aos componentes
 - [ ] ApiErrorMapper implementado e usado nos componentes
 - [ ] Testes unitários dos serviços adicionados
 - [ ] Cenário de falha documentado e testado (estoque interrompido)
 - [ ] Atualizar `ENTREGA_E_DETALHAMENTO_TECNICO.md` com decisões de integração e ciclos de vida do Angular utilizados

Próximos passos sugeridos
1. Eu posso implementar os arquivos de serviço Angular (ProdutoService, NotaService, api config) e aplicar as alterações em um componente como prova de conceito.
2. Eu posso adicionar o middleware `rs/cors` nos dois serviços Go (patches mínimos em `cmd/*/main.go`), se você quiser que eu faça as mudanças de backend também.

Escolha rápida
- Quer que eu aplique agora:
  - A) Apenas frontend (criar services e atualizar 1-2 componentes)
  - B) Backend + frontend (adicionar rs/cors nos dois microsserviços e implementar frontend)

Informe A ou B e eu começo a aplicar os patches.
