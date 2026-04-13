import { Component, OnInit } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { DrawerService } from '../../../shared/services/drawer.service';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { ProdutoService } from '../../../core/services/produto.service';
import { NotaService } from '../../../core/services/nota.service';
import { Produto } from '../../../core/models/produto.model';
import { NotaItemsTableComponent, NotaItem } from '../../../shared/components/nota-items-table/nota-items-table.component';
import { NotaItemAdderComponent, NotaItemAdicionado } from '../../../shared/components/nota-item-adder/nota-item-adder.component';
import { ApiErrorMapper } from '../../../core/services/api-error-mapper.service';

interface ItemNota {
  produto_codigo: string;
  produto_descricao: string;
  quantidade: number;
}

@Component({
  selector: 'app-nota-form',
  standalone: true,
  imports: [
    MatIconModule,
    MatButtonModule,
    NotaItemsTableComponent,
    NotaItemAdderComponent,
  ],
  templateUrl: './nota-form.component.html',
  styleUrls: ['./nota-form.component.scss'],
})
export class NotaFormComponent implements OnInit {
  produtos: Produto[] = [];
  itens: ItemNota[] = [];
  produtoMap = new Map<string, string>();

  constructor(
    private drawer: DrawerService,
    private snackbar: SnackbarService,
    private produtoService: ProdutoService,
    private notaService: NotaService,
    private apiErrorMapper: ApiErrorMapper
  ) {}

  ngOnInit(): void {
    this.produtoService.listAll().subscribe({
      next: (produtos) => {
        this.produtos = produtos.filter((p) => p.saldo > 0);
        this.produtoMap.clear();
        this.produtos.forEach((p) => this.produtoMap.set(p.codigo, p.descricao));
      },
      error: () => {
        this.produtos = [];
      },
    });
  }

  get itensParaTabela(): NotaItem[] {
    return this.itens.map((item, index) => ({
      id: index,
      produto_codigo: item.produto_codigo,
      quantidade: item.quantidade,
    }));
  }

  onItemAdicionado(event: NotaItemAdicionado): void {
    const produto = this.produtos.find((p) => p.codigo === event.produto_codigo);
    if (produto) {
      this.itens.push({
        produto_codigo: event.produto_codigo,
        produto_descricao: produto.descricao,
        quantidade: event.quantidade,
      });
    }
  }

  removerItem(index: number): void {
    this.itens.splice(index, 1);
  }

  fechar(): void {
    this.itens = [];
    this.drawer.close();
  }

  salvar(): void {
    if (this.itens.length === 0) {
      this.snackbar.error('Adicione pelo menos um item à nota');
      return;
    }

    const payload = this.itens.map((i) => ({ produto_codigo: i.produto_codigo, quantidade: i.quantidade }));
    this.notaService.create(payload).subscribe({
      next: () => {
        this.snackbar.success('Nota fiscal criada com sucesso!');
        this.fechar();
      },
      error: (err) => {
        const mapped = this.apiErrorMapper.map(err);
        this.snackbar.error('Falha ao criar nota: ' + mapped.message);
      }
    });
  }
}
