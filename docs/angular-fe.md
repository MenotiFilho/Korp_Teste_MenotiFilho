# Plano Angular Frontend - Korp Teste MenotiFilho

## Visao Geral

Projeto Angular standalone com Angular Material (MD3) que consome os microservicos `ms-estoque` e `ms-faturamento`.
Estrategia: **mock first, integration later** — construir toda a UI com dados mockados primeiro, depois conectar aos backends reais.

---

## 1. Stack Tecnica

| Item | Escolha |
|------|---------|
| Framework | Angular 19 (standalone components) |
| UI | Angular Material + CDK (Material Design 3) |
| CSS | Tema MD3 custom (dark) + variaveis CSS |
| HTTP | `HttpClient` nativo do Angular |
| Forms | `ReactiveFormsModule` |
| State | Services com `BehaviorSubject` (sem NgRx) |
| Router | `@angular/router` com lazy loading |
| Testes | Jasmine + Karma (padrao Angular) |

---

## 2. Estrutura de Pastas

```
apps/frontend/
├── src/
│   ├── app/
│   │   ├── core/                    # Servicos, interceptors, guards
│   │   │   ├── services/
│   │   │   │   ├── produto.service.ts
│   │   │   │   ├── nota.service.ts
│   │   │   │   └── mock-data.service.ts
│   │   │   ├── interceptors/
│   │   │   │   ├── error.interceptor.ts
│   │   │   │   └── mock.interceptor.ts
│   │   │   └── models/
│   │   │       ├── produto.model.ts
│   │   │       ├── nota.model.ts
│   │   │       └── error.model.ts
│   │   ├── shared/                  # Componentes reutilizaveis
│   │   │   ├── components/
│   │   │   │   ├── sidebar/
│   │   │   │   ├── top-bar/
│   │   │   │   ├── page-header/
│   │   │   │   ├── status-badge/
│   │   │   │   ├── confirm-dialog/
│   │   │   │   └── loading-overlay/
│   │   │   └── shared.module.ts
│   │   ├── features/                # Telas por feature
│   │   │   ├── dashboard/
│   │   │   │   ├── dashboard.component.ts
│   │   │   │   ├── dashboard.component.html
│   │   │   │   └── dashboard.component.scss
│   │   │   ├── produtos/
│   │   │   │   ├── produtos-list/
│   │   │   │   ├── produto-form/
│   │   │   │   └── produtos.module.ts
│   │   │   └── notas/
│   │   │       ├── notas-list/
│   │   │       ├── nota-detail/
│   │   │       ├── nota-form/
│   │   │       └── notas.module.ts
│   │   ├── app.component.ts         # Shell com sidenav
│   │   ├── app.component.html
│   │   ├── app.component.scss
│   │   ├── app.routes.ts            # Rotas principais
│   │   └── app.config.ts            # Config standalone
│   ├── environments/
│   │   ├── environment.ts           # Dev (localhost)
│   │   └── environment.prod.ts      # Prod
│   ├── styles/
│   │   ├── _theme.scss              # Tokens MD3 + variaveis CSS
│   │   ├── _typography.scss         # Tipografia
│   │   └── _spacing.scss            # Espacamento
│   ├── styles.scss                  # Global imports
│   └── index.html
├── angular.json
├── package.json
└── tsconfig.json
```

---

## 3. Rotas

| Rota | Componente | Lazy? |
|------|-----------|-------|
| `/` | redirect → `/dashboard` | — |
| `/dashboard` | `DashboardComponent` | Sim |
| `/produtos` | `ProdutosListComponent` | Sim |
| `/notas` | `NotasListComponent` | Sim |

Formularios e detalhes abrem em **drawer lateral** (mat-sidenav end) ou **side panel**, sem rotas proprias.

---

## 4. Fases de Implementacao

### Fase 1: Scaffold + Theme (mock)

**Objetivo:** Projeto Angular criado com tema MD3 dark e shell de layout.

- [x] Criar projeto Angular 19 com `ng new` (standalone, SCSS, routing)
- [x] Instalar Angular Material + CDK
- [x] Configurar tema MD3 dark com tokens do `design-system.pen` (`src/styles/_theme.scss`)
- [x] Criar variaveis CSS custom (`--color-success`, `--color-warning`, spacing, radius)
- [x] Criar `AppComponent` com `mat-sidenav-container` (sidebar + content)
- [x] Criar `SidebarComponent` (logo "Korp Notas" + nav items: Dashboard, Produtos, Notas)
- [x] Criar `TopBarComponent` (breadcrumb "Início / [pagina]" + avatar)
- [x] Configurar rotas com lazy loading (`/dashboard`, `/produtos`, `/notas`)
- [x] Build verificado com sucesso (`ng build`)

### Fase 2: Componentes Shared (mock)

**Objetivo:** Todos os componentes reutilizaveis criados e testados.

- [x] `PageHeaderComponent` — titulo + subtítulo + slot de ações (ng-content)
- [x] `StatusBadgeComponent` — badge ABERTA/FECHADA (cores exatas do .pen: warning/success)
- [x] `ConfirmDialogComponent` — mat-dialog genérico (title, message, confirmText, danger mode)
- [x] `LoadingOverlayComponent` — overlay fixo com mat-progress-spinner
- [x] `SnackbarService` — wrapper MatSnackBar (success, error, info) com estilos globais
- [x] Build verificado com sucesso (`ng build`)

### Fase 3: Dashboard (mock)

**Objetivo:** Tela de dashboard com KPIs e listas usando dados mockados.

- [x] Criar `MockDataService` com dados hardcoded (8 produtos, 6 notas)
- [x] Criar models TypeScript (`Produto`, `Nota`, `NotaItem`, `ApiError`)
- [x] 3x KPI cards: Total Produtos, Notas Abertas (warning), Notas Fechadas (success)
- [x] Tabela últimas 5 notas (mat-table, colunas: Número, Itens, Criado em, Status)
- [x] Lista produtos com estoque baixo (saldo < 5, com link "Ver todos")
- [x] Botões de ação rápida: "Nova Nota" + "Novo Produto" no page header
- [x] Build verificado com sucesso (`ng build`)

### Fase 4: Lista de Produtos (mock)

**Objetivo:** CRUD completo de produtos com drawer, usando mock.

- [x] `ProdutosListComponent` com mat-table + mat-paginator
- [x] Campo busca custom (input + mat-icon search, width 360, height 44)
- [x] Colunas: ID (60px), Código (120px), Descrição (fill), Saldo (100px), Ações (100px)
- [x] `ProdutoFormComponent` em dialog (width 400px, header/body/footer do .pen)
  - Campos: Código (obrigatório, max 20), Descrição (obrigatório, max 100), Saldo (>= 0)
  - Validação inline com mat-error
  - Cria e edita (mesmo form, editMode flag)
- [x] ConfirmDialog para exclusão (danger mode)
- [x] Snackbar de feedback (sucesso)
- [x] Botão editar + excluir por linha (icon buttons com tooltip)
- [x] Saldo destacado: warning se < 5, muted se 0
- [x] Build verificado com sucesso (`ng build`)

### Fase 5: Lista de Notas Fiscais (mock)

**Objetivo:** Listagem, criacao e impressao de notas com mock.

- [ ] `NotasListComponent` com mat-table + mat-paginator + mat-sort
- [ ] Campo busca por numero + dropdown filtro status (Todos/ABERTA/FECHADA)
- [ ] Colunas: ID, Numero, Itens (contagem), Criado em, Status, Acoes
- [ ] Botao imprimir (só visivel se ABERTA) com ConfirmDialog
- [ ] Loading overlay durante "impressao" (setTimeout simulando latencia)
- [ ] `NotaDetailComponent` em side panel (mat-sidenav end)
  - Mostra numero, status, data criacao
  - Tabela de itens (codigo, produto, quantidade)
  - Botao "Imprimir Nota" se ABERTA
- [ ] `NotaFormComponent` em drawer para criar nota
  - mat-select de produtos (mock: "CODIGO - Descricao (Saldo: XX)")
  - Input quantidade (inteiro > 0)
  - Botao "Adicionar Item" → tabela temporaria
  - Botao remover item da lista
  - Validacao: pelo menos 1 item, quantidade <= saldo
- [ ] Snackbar de feedback
- [ ] Verificar: fluxo completo criar/imprimir/detalhar com mock

### Fase 6: Services + Models (estrutura para integracao)

**Objetivo:** Criar services HTTP e models tipados, prontos para integrar.

- [ ] Models TypeScript:
  - `Produto` (id, codigo, descricao, saldo)
  - `Nota` (id, numero, status, itens[])
  - `NotaItem` (id, produto_codigo, quantidade)
  - `ApiError` (code, message, details, request_id)
- [ ] `ProdutoService`:
  - `list()`: GET /api/v1/produtos
  - `create(produto)`: POST /api/v1/produtos
  - ~~`update(id, produto)`: PUT (se existir)~~
  - ~~`delete(id)`: DELETE (se existir)~~
- [ ] `NotaService`:
  - `list()`: GET /api/v1/notas
  - `create(nota)`: POST /api/v1/notas
  - `imprimir(id)`: POST /api/v1/notas/{id}/imprimir
- [ ] `ErrorInterceptor` (HttpInterceptor) para padronizar erros
- [ ] Configurar `environment.ts` com base URLs:
  - `estoqueUrl: 'http://localhost:8081'`
  - `faturamentoUrl: 'http://localhost:8082'`

### Fase 7: Integracao com Microservicos

**Objetivo:** Substituir mock por chamadas HTTP reais, passo a passo.

- [ ] Ativar `HttpClientModule` em `app.config.ts`
- [ ] **Passo 1 - Listar produtos:** Substituir mock por `ProdutoService.list()`
- [ ] **Passo 2 - Criar produto:** Substituir mock por `ProdutoService.create()`
- [ ] **Passo 3 - Listar notas:** Substituir mock por `NotaService.list()`
- [ ] **Passo 4 - Criar nota:** Substituir mock por `NotaService.create()`
- [ ] **Passo 5 - Imprimir nota:** Substituir mock por `NotaService.imprimir()`
  - Testar cenarios de erro: INSUFFICIENT_STOCK, ESTOQUE_UNAVAILABLE
  - Verificar snackbar com mensagem padronizada
- [ ] **Passo 6 - Dashboard:** Usar services reais para KPIs
- [ ] Testar fluxo completo ponta a ponta com docker-compose rodando

---

## 5. Mapeamento Design → Angular Material

| Design (.pen / MD) | Angular Material |
|---------------------|------------------|
| Sidebar | `mat-sidenav-container` + `mat-sidenav` |
| Drawer lateral | `mat-sidenav` (position="end") |
| Side panel | `mat-sidenav` (position="end") |
| Tabela produtos | `mat-table` + `matSort` + `mat-paginator` |
| Tabela notas | `mat-table` + `matSort` + `mat-paginator` |
| Form field | `mat-form-field` + `matInput` |
| Select produto | `mat-select` + `mat-option` |
| Status badge | `mat-chip` (ngClass por status) |
| Botao primario | `mat-flat-button` |
| Botao icone | `mat-icon-button` + `mat-icon` |
| KPI card | `mat-card` |
| Dialog confirmacao | `mat-dialog` |
| Snackbar | `mat-snack-bar` |
| Loading | `mat-progress-spinner` + overlay |
| Tooltip | `mat-tooltip` |
| Divisor | `mat-divider` |
| Lista simples | `mat-list` + `mat-list-item` |

---

## 6. Configuracao de Cores (tokens CSS)

Implementar no `styles/_theme.scss` usando `@use '@angular/material' as mat`:

```scss
$dark-theme: mat.define-theme((
  color: (
    theme-type: dark,
    primary: mat.$indigo-palette,
    tertiary: mat.$violet-palette,
  ),
  typography: (
    brand-family: 'Roboto',
    plain-family: 'Roboto',
  ),
));

// Tokens custom
:root {
  --color-success: #22c55e;
  --color-success-container: #14532d;
  --color-warning: #f59e0b;
  --color-warning-container: #713f12;
  // Spacing
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;
}
```

---

## 7. Comando de Inicializacao

```bash
cd Korp_Teste_MenotiFilho/apps/frontend
ng new . --standalone --routing --style=scss --skip-git
ng add @angular/material --theme=custom --animations=enabled --typography=true
```

Depois ajustar `angular.json` para o projeto standalone dentro de `apps/frontend`.

---

## 8. Dependencias NPM

| Pacote | Versao | Uso |
|--------|--------|-----|
| `@angular/core` | ^19 | Framework |
| `@angular/material` | ^19 | Componentes UI |
| `@angular/cdk` | ^19 | Layout, overlay |
| `@angular/forms` | ^19 | Formularios reativos |
| `@angular/router` | ^19 | Roteamento |
| `@angular/animations` | ^19 | Animacoes Material |

Sem dependencias externas alem das padroes do Angular.

---

## 9. Alinhamento com desafio.md

### Funcionalidades obrigatorias (cobertas pelo plano)

| Requisito do desafio | Fase do plano |
|----------------------|---------------|
| Cadastro de Produtos (Codigo, Descricao, Saldo) | Fase 4 |
| Cadastro de Notas Fiscais (numero sequencial, status Aberta/Fechada, multiplos produtos) | Fase 5 |
| Impressao de Notas (botao visivel, loading, status → Fechada, baixa de saldo) | Fase 5 + Fase 7 |
| Arquitetura 2 microsservicos (Estoque + Faturamento) | Ja existente (backend) |
| Tratamento de falhas (feedback ao usuario, recuperacao) | Fase 7 Passo 5 |
| Conexao real com banco de dados | Ja existente (backend) |

### Detalhamento tecnico exigido no video

O video de entrega deve cobrir os seguintes topicos (documentar em `ENTREGA_E_DETALHAMENTO_TECNICO.md`):

| Topico | Onde documentar |
|--------|-----------------|
| Ciclos de vida Angular utilizados | `DESIGN_UI_UX.md` → secao "Ciclos de vida Angular e RxJS" |
| Uso de RxJS (como e onde) | `DESIGN_UI_UX.md` → secao "Padroes RxJS" |
| Bibliotecas utilizadas e finalidade | Secao 8 deste plano (Dependencias NPM) |
| Componentes visuais (Angular Material) | Secao 5 deste plano (Mapeamento) |
| Gerenciamento de dependencias Go | Ja documentado em `DECISOES_MICROSSERVICOS_GO.md` |
| Frameworks Go (nenhum — net/http puro) | Ja documentado em `DECISOES_MICROSSERVICOS_GO.md` |
| Tratamento de erros backend (JSON padronizado) | Contratos API `docs/ms-*-contratos_api.md` |

### Requisitos opcionais

| Opcional | Status | Observacao |
|----------|--------|------------|
| Tratamento de Concorrencia (saldo 1, duas notas simultaneas) | Ja implementado no backend (FOR UPDATE + transacao) |
| Uso de Inteligencia Artificial | Nao planejado — avaliar se ha tempo |
| Implementacao de Idempotencia | Ja implementado no backend (Idempotency-Key + invoice-print-{id}) |

---

## 10. Checklist de Entrega

- [ ] Projeto Angular compilando sem erros (`ng build`)
- [ ] Tema MD3 dark aplicado globalmente
- [ ] Shell com sidebar + topbar funcional
- [ ] Dashboard com KPIs e listas
- [ ] Cadastro de produtos (Codigo, Descricao, Saldo) com CRUD
- [ ] Cadastro de notas fiscais (numero, status, multiplos produtos)
- [ ] Impressao de notas (botao, loading, status → Fechada, baixa de saldo)
- [ ] Validacoes de formulario inline
- [ ] Feedback visual (snackbar, loading, dialog)
- [ ] Tratamento de falhas no frontend (erro de estoque indisponivel, saldo insuficiente)
- [ ] Integracao real com ms-estoque e ms-faturamento
- [ ] Video de apresentacao gravado e disponibilizado
- [ ] Detalhamento tecnico atualizado (ciclos de vida, RxJS, bibliotecas)

---

## 11. Referencias

- Design UI/UX: `design/DESIGN_UI_UX.md`
- Design System (.pen): `design/design-system.pen`
- API ms-estoque: `docs/ms-estoque_contratos_api.md`
- API ms-faturamento: `docs/ms-faturamento_contratos_api.md`
- Plano macro: `PLANO_KORP_TESTE_MENOTIFILHO.md`
- Requisitos: `desafio.md`
