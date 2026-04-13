# Detalhamento Tecnico

## Ciclos de Vida do Angular

- **ngOnInit**: utilizado em todos os componentes para carregar dados da API
  e inicializar subscriptions ao DrawerService.
- **ngOnDestroy**: utilizado para cancelar subscriptions (unsubscribe) e
  evitar memory leaks em componentes que escutam o estado do drawer.

---

## RxJS

RxJS e utilizado de forma centralizada no frontend:

- **Observable**: todos os metodos de ProdutoService e NotaService retornam
  Observable\<T\> via HttpClient.
- **BehaviorSubject**: DrawerService gerencia o estado do painel lateral
  (aberto/fechado, componente ativo) atraves de um BehaviorSubject exposto
  como state$.
- **.subscribe()**: consumo de observables nos componentes com padrao de
  cleanup em ngOnDestroy.
- **.pipe() + operadores**: AppComponent utiliza filter e map para reagir
  a eventos de navegacao finalizada do Router.

---

## Outras Bibliotecas

| Biblioteca | Finalidade |
|------------|------------|
| @angular/material | Componentes visuais seguindo Material Design |
| @angular/cdk | Foundation de acessibilidade e comportamento para o Angular Material |
| @angular/forms | Formularios reativos e template-driven |
| @angular/router | Roteamento com lazy loading |
| @angular/animations | Animacoes requeridas pelo Angular Material |
| rxjs | Programacao reativa, Observables, BehaviorSubject, operadores de stream |
| zone.js | Change detection do Angular |

---

## Componentes Visuais

A biblioteca de componentes visuais utilizada foi o **Angular Material**,
que forneceu: tabelas de dados (MatTable), dialogos modais (MatDialog),
seletores (MatSelect), icones (MatIcon), tooltips (MatTooltip),
paginacao (MatPaginator), spinners de loading (MatProgressSpinner),
painel lateral (MatSidenav), campos de formulario (MatFormField/MatInput),
notificacoes toast (MatSnackBar) e botoes (MatButton).

---

## Gerenciamento de Dependencias no Golang

Gerenciamento via **Go Modules** (go.mod), o gerenciador nativo do Go.

Dependencias utilizadas:

| Dependencia | Finalidade |
|-------------|------------|
| pgx/v5 | Driver PostgreSQL, utilizado via interface database/sql |
| rs/cors | Middleware CORS para permitir requisicoes do frontend |
| testify | Assercoes em testes unitarios e de integracao |

---

## Frameworks no Golang

Nenhum framework HTTP externo. Ambos os microsservicos utilizam apenas a
biblioteca padrao net/http do Go com o ServeMux aprimorado (Go 1.22+)
que suporta roteamento com metodo HTTP e parametros de path nativamente.

Os middlewares (Recover, RequestID, MaxBodyBytes, Logger com slog, CORS)
sao implementados como funcoes que envolvem http.Handler, aplicados em
cadeia manualmente.

---

## Tratamento de Erros e Excecoes no Backend

Ambos os microsservicos retornam erros em formato JSON padronizado:

```json
{
  "code": "INSUFFICIENT_STOCK",
  "message": "saldo insuficiente",
  "details": null,
  "request_id": "a1b2c3d4-..."
}
```

| Campo | Finalidade |
|-------|------------|
| code | Identificador programatico do erro |
| message | Descricao legivel do erro |
| details | Contexto adicional (opcional) |
| request_id | Identificador unico da requisicao para rastreamento |

Codigos de erro utilizados: INVALID_JSON, VALIDATION_ERROR,
PAYLOAD_TOO_LARGE, PRODUCT_NOT_FOUND, INVOICE_NOT_FOUND,
INSUFFICIENT_STOCK, INVOICE_NOT_ABERTA, INTERNAL_ERROR.

O middleware Recover captura panics (excecoes nao tratadas) e retorna
INTERNAL_ERROR com HTTP 500, impedindo que o processo seja derrubado.

Erros de validacao de dominio sao definidos como sentinel errors e
verificados com errors.Is() nos handlers.
