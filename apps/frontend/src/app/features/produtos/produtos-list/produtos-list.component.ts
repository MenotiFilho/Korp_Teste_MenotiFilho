import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { Subscription } from 'rxjs';
import { MockDataService } from '../../../core/services/mock-data.service';
import { Produto } from '../../../core/models/produto.model';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { DataTableComponent } from '../../../shared/components/data-table/data-table.component';
import { ConfirmDialogComponent } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { DrawerService } from '../../../shared/services/drawer.service';

@Component({
  selector: 'app-produtos-list',
  imports: [
    MatIconModule,
    MatTableModule,
    MatTooltipModule,
    MatDialogModule,
    PageHeaderComponent,
    DataTableComponent,
  ],
  templateUrl: './produtos-list.component.html',
  styleUrl: './produtos-list.component.scss',
})
export class ProdutosListComponent implements OnInit, OnDestroy {
  displayedColumns = ['codigo', 'descricao', 'saldo', '_editar', '_excluir'];
  searchTerm = '';
  pageSize = 10;
  pageIndex = 0;
  produtos: Produto[] = [];
  produtosFiltrados: Produto[] = [];

  private drawerSub!: Subscription;

  constructor(
    private mockData: MockDataService,
    private dialog: MatDialog,
    private snackbar: SnackbarService,
    private drawer: DrawerService
  ) {}

  ngOnInit(): void {
    this.carregarProdutos();
    this.drawerSub = this.drawer.state$.subscribe((state) => {
      if (!state.open) this.carregarProdutos();
    });
  }

  ngOnDestroy(): void {
    this.drawerSub?.unsubscribe();
  }

  get paginaAtual(): Produto[] {
    const start = this.pageIndex * this.pageSize;
    return this.produtosFiltrados.slice(start, start + this.pageSize);
  }

  carregarProdutos(): void {
    this.produtos = this.mockData.getProdutos();
    this.produtosFiltrados = [...this.produtos];
  }

  onSearch(term: string): void {
    this.searchTerm = term;
    const t = term.trim().toLowerCase();
    this.produtosFiltrados = this.produtos.filter(
      (p) => p.codigo.toLowerCase().includes(t) || p.descricao.toLowerCase().includes(t)
    );
    this.pageIndex = 0;
  }

  onPageChange(event: PageEvent): void {
    this.pageIndex = event.pageIndex;
    this.pageSize = event.pageSize;
  }

  abrirNovo(): void {
    this.drawer.open('produto-form');
  }

  editar(produto: Produto): void {
    this.drawer.open('produto-edit-' + produto.id);
  }

  excluir(produto: Produto): void {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      data: {
        title: 'Excluir Produto',
        message: `Deseja excluir o produto "${produto.descricao}"?`,
        confirmText: 'Excluir',
        cancelText: 'Cancelar',
        danger: true,
      },
    });

    dialogRef.afterClosed().subscribe((confirmed) => {
      if (confirmed) {
        this.produtos = this.produtos.filter((p) => p.id !== produto.id);
        this.onSearch(this.searchTerm);
        this.snackbar.success('Produto excluído com sucesso!');
      }
    });
  }
}
