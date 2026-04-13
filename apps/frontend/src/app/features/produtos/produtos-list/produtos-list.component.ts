import { Component, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatPaginator, MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatTooltipModule } from '@angular/material/tooltip';
import { Subscription } from 'rxjs';
import { MockDataService } from '../../../core/services/mock-data.service';
import { Produto } from '../../../core/models/produto.model';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { ConfirmDialogComponent } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { DrawerService } from '../../../shared/services/drawer.service';

@Component({
  selector: 'app-produtos-list',
  imports: [
    FormsModule,
    MatIconModule,
    MatTableModule,
    MatPaginatorModule,
    MatDialogModule,
    MatTooltipModule,
    PageHeaderComponent,
  ],
  templateUrl: './produtos-list.component.html',
  styleUrl: './produtos-list.component.scss',
})
export class ProdutosListComponent implements OnInit, OnDestroy {
  displayedColumns = ['id', 'codigo', 'descricao', 'saldo', 'editar', 'excluir'];
  searchTerm = '';
  pageSize = 10;
  pageIndex = 0;
  produtos: Produto[] = [];
  produtosFiltrados: Produto[] = [];

  @ViewChild(MatPaginator) paginator!: MatPaginator;

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
      if (!state.open) {
        this.carregarProdutos();
      }
    });
  }

  ngOnDestroy(): void {
    this.drawerSub?.unsubscribe();
  }

  get totalRegistros(): number {
    return this.produtosFiltrados.length;
  }

  get paginaAtual(): Produto[] {
    const start = this.pageIndex * this.pageSize;
    return this.produtosFiltrados.slice(start, start + this.pageSize);
  }

  carregarProdutos(): void {
    this.produtos = this.mockData.getProdutos();
    this.produtosFiltrados = [...this.produtos];
    this.pageIndex = 0;
  }

  filtrar(): void {
    const termo = this.searchTerm.trim().toLowerCase();
    this.produtosFiltrados = this.produtos.filter(
      (p) =>
        p.codigo.toLowerCase().includes(termo) ||
        p.descricao.toLowerCase().includes(termo)
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
    this.drawer.open('produto-form');
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
        this.filtrar();
        this.snackbar.success('Produto excluído com sucesso!');
      }
    });
  }
}
