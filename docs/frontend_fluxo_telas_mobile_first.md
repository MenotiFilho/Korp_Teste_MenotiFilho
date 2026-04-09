# Frontend - Fluxo de Telas Fechado (Mobile First)

Documento consolidado para prototipacao no Pencil, com foco em UX do fluxo de impressao de notas fiscais consumindo `ms-estoque` e `ms-faturamento`.

## 1) Objetivo UX

- Tornar a acao de impressao extremamente visivel e intuitiva no mobile.
- Evitar erro de operacao (impressao acidental ou em status invalido).
- Dar feedback claro de processamento, sucesso e falha.
- Garantir consistencia visual entre status da nota e estado do botao.

## 2) Decisoes fechadas

- Navegacao principal no mobile por barra inferior (`Produtos`, `Notas`).
- Acao principal de impressao em barra fixa inferior na tela de detalhe da nota.
- Confirmacao de impressao via bottom sheet.
- Conteudo da confirmacao com itens e quantidades.
- Para status diferente de `ABERTA`, botao visivel porem desabilitado com motivo.
- Ao confirmar impressao, mostrar indicador de processamento no botao.
- Ao finalizar com sucesso, manter usuario na mesma tela e atualizar status para `FECHADA`.

## 3) Arquitetura de navegacao (mobile first)

## Rotas

- `/produtos`
- `/notas`
- `/notas/:id`

## Estrutura de navegacao

- Barra inferior fixa:
  - Aba 1: `Produtos`
  - Aba 2: `Notas`
- A tela de detalhe da nota (`/notas/:id`) e aberta a partir da lista de notas.

## 4) Fluxo de telas fechado

## Tela A - Lista de Notas (`/notas`)

Elementos:

- Lista em cards com: numero da nota, status (`ABERTA`/`FECHADA`), quantidade de itens.
- CTA secundario: `Nova nota`.
- Toque no card abre detalhe da nota.

Estados:

- Loading: skeleton/cards placeholder.
- Empty: mensagem "Nenhuma nota fiscal cadastrada" + CTA `Criar primeira nota`.
- Error: mensagem amigavel + CTA `Tentar novamente`.

## Tela B - Detalhe da Nota (`/notas/:id`)

Elementos:

- Cabecalho com numero da nota e chip de status.
- Resumo dos itens da nota (produto + quantidade).
- Barra fixa inferior com CTA principal: `Imprimir nota`.

Regra do CTA por status:

- `ABERTA`: botao ativo.
- `FECHADA`: botao desabilitado com texto de apoio: "Apenas notas abertas podem ser impressas".

## Tela C - Confirmacao de Impressao (Bottom Sheet)

Conteudo:

- Titulo: `Confirmar impressao da nota`.
- Numero e status atual da nota.
- Lista resumida: itens e quantidades.
- Acao primaria: `Confirmar impressao`.
- Acao secundaria: `Cancelar`.

Comportamento:

- Fechar sem efeito ao cancelar.
- Iniciar processamento ao confirmar.

## Tela D - Processamento (estado na Tela B)

Comportamento esperado:

- Botao da barra fixa troca para loading (`Imprimindo...` + spinner).
- Desabilitar novo clique durante processamento.
- Manter usuario no contexto da nota.

## Tela E - Pos-processamento (resultado na Tela B)

Sucesso:

- Atualizar status para `FECHADA`.
- Desabilitar CTA de impressao com motivo visivel.
- Exibir feedback curto: `Nota fechada com sucesso`.

Falha:

- Manter status atual da nota (tipicamente `ABERTA`).
- Exibir mensagem de erro clara (com opcao de nova tentativa).
- Manter CTA ativo para retry apenas se nota continuar `ABERTA`.

## 5) Regras de negocio refletidas na UX

- Nao permitir impressao de nota com status diferente de `ABERTA`.
- Ao imprimir com sucesso, refletir status `FECHADA` imediatamente.
- Atualizacao de saldo de produtos ocorre no backend no fluxo de impressao e deve ser refletida no frontend apos refresh dos dados.

## 6) Estados e microcopy recomendados

- Botao ativo: `Imprimir nota`
- Loading: `Imprimindo...`
- Sucesso: `Nota fechada com sucesso`
- Bloqueio por status: `Apenas notas abertas podem ser impressas`
- Erro generico: `Nao foi possivel concluir a impressao. Tente novamente.`

## 7) Heuristicas UX aplicadas

- Visibilidade do status: chip sempre presente no topo da nota.
- Prevencao de erro: confirmacao antes da acao irreversivel de fechamento.
- Feedback imediato: loading, sucesso e erro no mesmo contexto.
- Controle do usuario: cancelar antes de confirmar; retry em falha.
- Consistencia: mesmos termos de status em toda a interface (`ABERTA`, `FECHADA`).

## 8) Responsividade (mobile first)

- Breakpoint base mobile como referencia principal de layout.
- Areas de toque amplas para CTA principal e acoes da bottom sheet.
- Barra fixa inferior preservada em viewport pequeno.
- Em desktop/tablet, manter mesma logica e ampliar area de conteudo, sem mudar semantica do fluxo.

## 8.1) Adaptacao para navegador de PC (sem perder mobile first)

- Desktop deve manter o mesmo fluxo funcional (lista -> detalhe -> confirmar -> processar -> resultado), mudando apenas distribuicao visual.
- Em largura maior, usar layout em duas colunas quando fizer sentido:
  - coluna principal com detalhes da nota;
  - coluna lateral com resumo rapido e status.
- A barra inferior do mobile pode evoluir para cabecalho superior com navegacao horizontal no desktop.
- CTA de impressao deve continuar muito visivel no desktop:
  - preferencialmente fixo no rodape do card de detalhe ou na lateral de acoes;
  - evitar esconder em menu de contexto.
- Bottom sheet no mobile pode virar modal central no desktop, mantendo o mesmo conteudo e decisoes.
- Aumentar largura util de listas e tabelas sem aumentar ruido visual:
  - mais espacamento horizontal;
  - manter hierarquia de informacao e escaneabilidade.
- Preservar feedback de estados no desktop com o mesmo texto e semantica (loading, sucesso, erro, bloqueio por status).

## 8.2) Breakpoints recomendados para prototipo

- Mobile: ate 767px.
- Tablet: 768px a 1023px.
- Desktop: 1024px ou mais.
- Regra de ouro: nao criar fluxo alternativo por breakpoint; apenas adaptar composicao da interface.

## 9) Integracao entre microsservicos no fluxo de impressao

- Frontend chama `ms-faturamento` para imprimir nota.
- `ms-faturamento` valida status e chama `ms-estoque` para baixa.
- Em sucesso na baixa, `ms-faturamento` fecha a nota.
- Frontend atualiza tela da nota e lista de notas apos resposta.
- Frontend tambem deve atualizar lista de produtos/saldos quando o usuario acessar `Produtos` apos impressao.

## 10) Lacuna em aberto para fechar prototipo

Decisao pendente no resumo de confirmacao:

- Opcao 1: mostrar apenas itens e quantidades.
- Opcao 2: mostrar itens, quantidades e saldo atual/posterior.
- Opcao 3: mostrar saldo apenas em casos de estoque baixo.

Recomendacao UX: **Opcao 2** quando houver espaco visual suficiente; em telas muito pequenas, usar **Opcao 3**.

## 11) Checklist de prototipacao no Pencil

Use esta checklist na ordem para construir e validar o prototipo.

### 11.1 Preparacao de base visual

- [ ] Definir grid mobile base (largura de referencia e espacamentos principais).
- [ ] Definir tipografia de titulos, corpo e microcopy de feedback.
- [ ] Definir paleta para status (`ABERTA`, `FECHADA`), acao primaria, erro e sucesso.
- [ ] Definir componentes reutilizaveis: card de nota, chip de status, botao primario, bottom sheet, toast/snackbar.
- [ ] Garantir area de toque confortavel para botoes e itens clicaveis.

### 11.2 Navegacao principal

- [ ] Prototipar barra inferior fixa com `Produtos` e `Notas`.
- [ ] Destacar aba ativa de forma clara e consistente.
- [ ] Validar navegacao entre abas sem perder contexto visual.

### 11.3 Tela A - Lista de Notas

- [ ] Criar layout de lista com cards de nota (numero, status, quantidade de itens).
- [ ] Incluir CTA `Nova nota` visivel no contexto da lista.
- [ ] Prototipar estado de loading (skeleton/cards placeholder).
- [ ] Prototipar estado vazio com mensagem e CTA `Criar primeira nota`.
- [ ] Prototipar estado de erro com mensagem amigavel e CTA `Tentar novamente`.
- [ ] Validar que toque no card abre `Tela B - Detalhe da Nota`.

### 11.4 Tela B - Detalhe da Nota

- [ ] Exibir cabecalho com numero da nota e chip de status.
- [ ] Exibir resumo dos itens (produto + quantidade) com leitura facil no mobile.
- [ ] Inserir barra fixa inferior com CTA `Imprimir nota`.
- [ ] Prototipar variacao com nota `ABERTA` (botao ativo).
- [ ] Prototipar variacao com nota `FECHADA` (botao desabilitado + motivo visivel).

### 11.5 Tela C - Bottom sheet de confirmacao

- [ ] Prototipar abertura da bottom sheet ao tocar em `Imprimir nota`.
- [ ] Mostrar numero, status e resumo de itens/quantidades.
- [ ] Inserir acao primaria `Confirmar impressao` e secundaria `Cancelar`.
- [ ] Validar fechamento sem efeito ao cancelar.
- [ ] Validar transicao para processamento ao confirmar.

### 11.6 Tela D - Processamento

- [ ] Prototipar estado de loading no botao (`Imprimindo...` + spinner).
- [ ] Bloquear novo clique durante processamento.
- [ ] Manter usuario na mesma tela de detalhe durante a espera.

### 11.7 Tela E - Resultado

- [ ] Sucesso: atualizar status para `FECHADA` na mesma tela.
- [ ] Sucesso: desabilitar CTA com mensagem de bloqueio por status.
- [ ] Sucesso: mostrar feedback curto `Nota fechada com sucesso`.
- [ ] Falha: manter status atual (ex.: `ABERTA`) e exibir erro claro.
- [ ] Falha: permitir nova tentativa quando a nota permanecer `ABERTA`.

### 11.8 Regra de negocio e consistencia

- [ ] Garantir visualmente que impressao so ocorre para status `ABERTA`.
- [ ] Garantir consistencia de terminologia (`ABERTA`, `FECHADA`) em todas as telas.
- [ ] Garantir alinhamento do fluxo com regra de atualizacao de saldo apos impressao bem-sucedida.

### 11.9 Responsividade e validacao final

- [ ] Validar legibilidade e hierarquia visual em viewport mobile pequena.
- [ ] Validar comportamento da barra fixa inferior sem sobrepor conteudo critico.
- [ ] Validar adaptacao para tablet/desktop mantendo semantica do fluxo.
- [ ] Revisar estados obrigatorios: loading, empty, error, sucesso, bloqueio por status.
- [ ] Revisar clareza do fluxo completo (lista -> detalhe -> confirmar -> processar -> resultado).

### 11.10 Checklist especifico para experiencia no PC

- [ ] Prototipar versao desktop da lista de notas com melhor aproveitamento de largura (sem poluicao).
- [ ] Prototipar detalhe da nota em duas colunas quando houver espaco, mantendo leitura clara.
- [ ] Garantir CTA `Imprimir nota` sempre visivel sem depender de scroll longo.
- [ ] Trocar bottom sheet por modal central no desktop, mantendo os mesmos dados e acoes.
- [ ] Validar alinhamento de chips de status, botoes e mensagens para leitura rapida em monitor.
- [ ] Validar estados de erro e sucesso no desktop sem deslocar o usuario do contexto da nota.
