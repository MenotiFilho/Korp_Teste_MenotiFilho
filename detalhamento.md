# Análise Técnica - Sistema de Emissão de Notas Fiscais

**Projeto:** Korp_Teste_MenotiFilho  
**Data:** 14 de abril de 2026  

## 1. Visão Geral do Sistema

O Sistema de Emissão de Notas Fiscais é uma aplicação full-stack composta por:
- **Frontend:** Angular 19 com Angular Material
- **Backend:** 2 microsserviços em Go 1.23
- **Banco de dados:** 2 instâncias PostgreSQL 16

## 2. Funcionalidades Implementadas

### 2.1 Cadastro de Produtos
- **Campos:** Código (único), Descrição, Saldo (quantidade em estoque)
- **Persistência:** PostgreSQL dedicado ao microsserviço de estoque
- **Operações:** CRUD completo (criação, listagem, atualização, exclusão)
- **Validações:** Código único, saldo não negativo

### 2.2 Cadastro de Notas Fiscais
- **Campos:** Numeração sequencial automática, Status (ABERTA/FECHADA), Múltiplos produtos com quantidades
- **Persistência:** PostgreSQL dedicado ao microsserviço de faturamento
- **Operações:** CRUD completo (criação, listagem, atualização, exclusão)
- **Regras de negócio:**
  - Status inicial sempre ABERTA
  - Numeração sequencial via sequence PostgreSQL
  - Itens com quantidade positiva obrigatória

### 2.3 Impressão de Notas Fiscais
- **Fluxo:**
  1. Frontend exibe botão de impressão intuitivo
  2. Clique exibe indicador de processamento (MatProgressSpinner)
  3. Validação: nota deve estar ABERTA
  4. Chamada HTTP para ms-faturamento
  5. ms-faturamento chama ms-estoque para baixa de estoque
  6. Em sucesso: status atualizado para FECHADA
  7. Em falha: nota permanece ABERTA com feedback ao usuário
- **Proteções:**
  - Idempotency-Key obrigatório (`invoice-print-{id}`)
  - Circuit breaker entre microsserviços
  - Retry com backoff para erros transientes
  - Reconciliation Job para tratar inconsistencias geradas por falhas na rede

## 3. Detalhamento Técnico

### 3.1 Ciclos de Vida do Angular

#### Hooks Utilizados:

**ngOnInit:**
- **Componentes:** DashboardComponent, ProdutosListComponent, NotasListComponent, NotaDetailComponent, ProdutoFormComponent, NotaFormComponent
- **Propósito:** Inicialização de dados via chamadas HTTP, configuração de subscriptions reativas
- **Exemplo (NotaFormComponent):**
  - Arquivo: `apps/frontend/src/app/features/notas/nota-form/nota-form.component.ts:47-54`

  ```typescript
  ngOnInit(): void {
    this.carregarProdutos();
    this.drawer.state$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((state) => {
      if (state.open && state.component === 'nota-form') {
        this.carregarProdutos();
      }
    });
  }
  ```

**Router.events no constructor do AppComponent:**
- **Propósito:** Filtrar eventos de navegação para atualizar título da página e fechar drawer
- **Uso de operadores RxJS:** `filter` e `map`
- **Arquivo:** `apps/frontend/src/app/app.component.ts:42-50`

  ```typescript
  this.router.events
    .pipe(
      filter((event) => event instanceof NavigationEnd),
      map((event) => (event as NavigationEnd).urlAfterRedirects)
    )
    .subscribe((url) => {
      this.pageTitle = this.routeTitles[url] || 'Início';
      this.drawer.close();
    });
  ```

### 3.2 Uso de RxJS

- **ngOnInit**: utilizado em todos os componentes para carregar dados da API
  e inicializar subscriptions ao DrawerService.

**Observable em Services:**
- **ProdutoService:** Todos os métodos retornam `Observable<T>`
  - Arquivo: `apps/frontend/src/app/core/services/produto.service.ts:13-15`

  ```typescript
  listAll(): Observable<Produto[]> {
    return this.http.get<Produto[]>(`${this.base}/api/v1/produtos`);
  }
  ```
- **NotaService:** Padrão idêntico para todas as operações CRUD
  - Arquivo: `apps/frontend/src/app/core/services/nota.service.ts:13-19`

**BehaviorSubject para Gerenciamento de Estado:**
- **DrawerService:** Gerencia estado reativo do drawer lateral
  - Arquivo: `apps/frontend/src/app/shared/services/drawer.service.ts:14-28`

  ```typescript
  private state = new BehaviorSubject<DrawerState>({
    open: false,
    width: 424,
    component: '',
  });

  state$ = this.state.asObservable();

  open(component: string, width = 424): void {
    this.state.next({ open: true, width, component });
  }
  ```

**Padrão .subscribe() com auto-cleanup:**
- Componentes usam `subscribe()` com `next/error` handlers
- Auto-cleanup via `takeUntilDestroyed(this.destroyRef)` no `.pipe()`
- Elimina necessidade de `OnDestroy`, `Subscription` e `unsubscribe()` manual

**Operadores .pipe():**
- **AppComponent:** `filter` e `map` para processamento de eventos de rota
  - Arquivo: `apps/frontend/src/app/app.component.ts:42-46`

  ```typescript
  .pipe(
    filter((event) => event instanceof NavigationEnd),
    map((event) => (event as NavigationEnd).urlAfterRedirects)
  )
  ```

### 3.3 Outras Bibliotecas Utilizadas

#### Frontend:

| Biblioteca | Versão | Finalidade |
|------------|--------|------------|
| `rxjs` | ~7.8.0 | Programação reativa |
| `zone.js` | ~0.15.0 | Detecção de mudanças do Angular |
| `@angular/forms` | 19.2.0 | Formulários reativos |
| `@angular/router` | 19.2.0 | Roteamento com lazy loading |

#### Backend (Go):

| Biblioteca | Finalidade |
|------------|------------|
| `jackc/pgx/v5` | Driver PostgreSQL |
| `rs/cors` | Middleware CORS |
| `stretchr/testify` | Assertions e mocks para testes |

### 3.4 Biblioteca de Componentes Visuais 

Utilizei a angular material para agilizar o desenvolvimento e manter uma padronizacao dos componentes, com ela também usei a angular/animation e a angular/cdk.

| Biblioteca | Versão | Finalidade |
|------------|--------|------------|
| `@angular/material` | 19.2.19 | Componentes UI Material Design |
| `@angular/animations` | 19.2.20 | Animações Material Design |
| `@angular/cdk` | 19.2.19 | Toolkit base para Angular Material |

### 3.5 Gerenciamento de Dependências no Golang

**Go Modules:**
- Cada microsserviço possui `go.mod` e `go.sum` próprios (arquitetura multi-module)
- Versão Go: 1.23.0
- Módulos:
  - `github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque`
  - `github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento`

**Dependências principais:**
- Arquivo: `apps/ms-estoque/go.mod:5-10`
```go
require (
    github.com/jackc/pgx/v5 v5.7.2  // Driver PostgreSQL
    github.com/rs/cors v1.11.0       // Middleware CORS
    github.com/stretchr/testify v1.8.1 // Testes unitários
)
```

### 3.6 Frameworks no Golang

**Nenhum framework HTTP externo utilizado.**

- Biblioteca padrão: `net/http` com `http.ServeMux` (Go 1.22+)

**Middlewares customizados:**
- `Recover`: Captura panics e retorna 500 INTERNAL_ERROR
- `RequestID`: Gera X-Request-ID para correlação de logs
- `MaxBodyBytes`: Limita tamanho do payload
- `Logger`: Logs estruturados em JSON
- `CORS`: Configuração via `rs/cors`

### 3.7 Tratamento de Erros e Exceções no Backend

#### Formato Padrão de Erro JSON com request_id para rastreabilidade:
```json
{
  "code": "INSUFFICIENT_STOCK",
  "message": "saldo insuficiente",
  "details": null,
  "request_id": "a1b2c3d4-..."
}
```

#### Códigos de Erro:

| Código HTTP | Código Interno | Descrição |
|-------------|----------------|-----------|
| 400 | INVALID_JSON | JSON mal formatado |
| 400 | VALIDATION_ERROR | Dados de entrada inválidos |
| 404 | PRODUCT_NOT_FOUND | Produto não encontrado |
| 404 | INVOICE_NOT_FOUND | Nota fiscal não encontrada |
| 409 | PRODUCT_CODIGO_ALREADY_EXISTS | Código de produto duplicado |
| 409 | INSUFFICIENT_STOCK | Saldo insuficiente |
| 409 | INVOICE_NOT_ABERTA | Nota fiscal não está ABERTA |
| 413 | PAYLOAD_TOO_LARGE | Payload excede limite |
| 500 | INTERNAL_ERROR | Erro interno do servidor |
| 503 | ESTOQUE_UNAVAILABLE | Microsserviço de estoque indisponível |
| 500 | PRINT_STATUS_UPDATE_FAILED | Falha parcial na impressão |

#### Mecanismos de Proteção:

**Recover Middleware:**
- Arquivo: `apps/ms-estoque/internal/middleware/recover.go:10-21`

```go
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "panic", rec)
				httpapi.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
```

**Circuit Breaker (ms-faturamento → ms-estoque):**
- Estados: Closed (normal), Open (bloqueado), HalfOpen (teste)
- Configuração: 3 falhas consecutivas abrem circuito, 10s para half-open
- Retry: 1 tentativa extra apenas para erros de conexão (timeout, connection refused)

**Idempotência:**
- Endpoint `POST /api/v1/estoque/baixa` exige header `Idempotency-Key`
- Chave determinística: `invoice-print-{id}` para fluxo de impressão
- Constraint UNIQUE no banco para replay no-op

**Logs Estruturados:**
- Formato JSON com correlação por `request_id`
- Middleware gera `X-Request-ID` automaticamente

**Reconciliation Job:**
- Problema: Se ms-estoque baixa estoque com sucesso mas a resposta é perdida (network failure), a nota fica ABERTA com estoque já reduzido
- Solução: Background job no ms-faturamento que detecta e corrige automaticamente
- Frequência: Roda a cada 30 segundos como goroutine
- Threshold: Considera notas ABERTA com mais de 2 minutos de idade
- Fluxo:
  1. Query: `SELECT notas WHERE status='ABERTA' AND created_at < NOW() - 2min`
  2. Para cada nota: chama `GET /api/v1/estoque/baixas/{key}` no ms-estoque
  3. Se a idempotency key existe (estoque já baixou): `UPDATE status = 'FECHADA'`
  4. Se não existe: pula (estoque ainda não foi processado)
- Endpoint ms-estoque: `GET /api/v1/estoque/baixas/{idempotency_key}` — retorna 200 se existe, 404 se não
- Implementação:
  - `apps/ms-faturamento/internal/service/reconciliation_service.go`
  - `apps/ms-estoque/internal/httpapi/stock_handler.go:CheckIdempotencyKey`
  - `apps/ms-estoque/internal/repository/product_repository.go:IdempotencyKeyExists`

**Concorrência:**
- Problema: Dois usarios tentando dar baixa no mesmo momento em um mesmo objeto
- Solução: O sistema usa lock pessimista a nível de linha no banco de dados:
- Arquivo: `apps/ms-estoque/internal/repository/product_repository.go:221`

```go
// SELECT com FOR UPDATE - bloqueia a linha até o commit
err := tx.QueryRowContext(ctx, "SELECT saldo FROM produtos WHERE codigo = $1 FOR UPDATE", codigo).Scan(&saldo)
```

### 3.8 Tratamento de Erros no Frontend

**ApiErrorMapper Service:**
- Traduz erros HTTP para mensagens amigáveis
- Mapeamento de status codes comuns (0, 408, 429, 502, 503, 504)
- Exibe mensagens via MatSnackBar

**Mensagens amigáveis:**
- Arquivo: `apps/frontend/src/app/core/services/api-error-mapper.service.ts:5-12`

```typescript
const FRIENDLY_MESSAGES: Record<number, string> = {
  0:   'Serviço indisponível. Tente novamente mais tarde.',
  408: 'Tempo de requisição esgotado. Tente novamente.',
  429: 'Muitas requisições. Aguarde um momento e tente novamente.',
  502: 'Serviço temporariamente indisponível. Tente novamente.',
  503: 'Serviço em manutenção. Tente novamente mais tarde.',
  504: 'Tempo de resposta esgotado. Tente novamente mais tarde.',
};
```