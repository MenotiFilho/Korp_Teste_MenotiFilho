import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { Subscription } from 'rxjs';
import { MockDataService } from '../../../core/services/mock-data.service';
import { NotaService } from '../../../core/services/nota.service';
import { Nota } from '../../../core/models/nota.model';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { StatusBadgeComponent } from '../../../shared/components/status-badge/status-badge.component';
import { DataTableComponent, FilterOption } from '../../../shared/components/data-table/data-table.component';
import { ConfirmDialogComponent } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { DrawerService } from '../../../shared/services/drawer.service';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { LoadingOverlayComponent } from '../../../shared/components/loading-overlay/loading-overlay.component';

@Component({
  selector: 'app-notas-list',
  imports: [
    MatIconModule,
    MatTableModule,
    MatTooltipModule,
    MatDialogModule,
    PageHeaderComponent,
    StatusBadgeComponent,
    DataTableComponent,
    LoadingOverlayComponent,
  ],
  templateUrl: './notas-list.component.html',
  styleUrl: './notas-list.component.scss',
})
export class NotasListComponent implements OnInit, OnDestroy {
  displayedColumns = ['numero', 'itens', 'criado_em', 'status', '_visualizar', '_imprimir', '_excluir'];
  pageSize = 10;
  pageIndex = 0;
  notas: Nota[] = [];
  notasFiltradas: Nota[] = [];
  loading = false;
  searchTerm = '';
  statusFilter = '';

  filterOptions: FilterOption[] = [
    { value: '', label: 'Todos' },
    { value: 'ABERTA', label: 'ABERTA' },
    { value: 'FECHADA', label: 'FECHADA' },
  ];

  private drawerSub!: Subscription;

  constructor(
    private mockData: MockDataService,
    private notaService: NotaService,
    private drawer: DrawerService,
    private dialog: MatDialog,
    private snackbar: SnackbarService
  ) {}

  ngOnInit(): void {
    this.carregarNotas();
    this.drawerSub = this.drawer.state$.subscribe((state) => {
      if (!state.open) this.carregarNotas();
    });
  }

  ngOnDestroy(): void {
    this.drawerSub?.unsubscribe();
  }

  get paginaAtual(): Nota[] {
    const start = this.pageIndex * this.pageSize;
    return this.notasFiltradas.slice(start, start + this.pageSize);
  }

  carregarNotas(): void {
    // Phase A: use backend service to fetch latest notes, fallback to mock on error
    this.notaService.listLatest().subscribe({
      next: (r) => { this.notas = r; this.aplicarFiltros(); },
      error: () => { this.notas = this.mockData.getNotas(); this.aplicarFiltros(); }
    });
  }

  aplicarFiltros(): void {
    let r = [...this.notas];
    if (this.searchTerm.trim()) {
      r = r.filter((n) => String(n.numero).includes(this.searchTerm.trim()));
    }
    if (this.statusFilter) {
      r = r.filter((n) => n.status === this.statusFilter);
    }
    this.notasFiltradas = r;
  }

  onSearch(term: string): void {
    this.searchTerm = term;
    this.pageIndex = 0;
    this.aplicarFiltros();
  }

  onFilter(value: string): void {
    this.statusFilter = value;
    this.pageIndex = 0;
    this.aplicarFiltros();
  }

  onPageChange(event: PageEvent): void {
    this.pageIndex = event.pageIndex;
    this.pageSize = event.pageSize;
  }

  formatarData(iso: string | undefined): string {
    if (!iso) return '';
    return new Date(iso).toLocaleDateString('pt-BR', {
      day: '2-digit', month: '2-digit', year: 'numeric',
    });
  }

  formatarNumero(num: number): string {
    return String(num).padStart(5, '0');
  }

  abrirNovaNota(): void {
    this.drawer.open('nota-form');
  }

  visualizar(nota: Nota): void {
    this.drawer.open('nota-detail-' + nota.id);
  }

  imprimir(nota: Nota): void {
    if (nota.status !== 'ABERTA') return;
    this.loading = true;
    setTimeout(() => {
      const result = this.mockData.imprimirNota(nota.id);
      this.loading = false;
      if (result.success) {
        this.snackbar.success('Nota impressa com sucesso!');
        this.carregarNotas();
      } else {
        this.snackbar.error(result.error || 'Erro ao imprimir nota');
      }
    }, 1500);
  }

  excluir(nota: Nota): void {
    if (nota.status !== 'ABERTA') return;
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      data: {
        title: 'Excluir Nota',
        message: `Deseja excluir a nota #${this.formatarNumero(nota.numero)}?`,
        confirmText: 'Excluir',
        cancelText: 'Cancelar',
        danger: true,
      },
    });

    dialogRef.afterClosed().subscribe((confirmed) => {
      if (confirmed) {
        const deleted = this.mockData.deleteNota(nota.id);
        if (deleted) {
          this.snackbar.success('Nota excluída com sucesso!');
          this.carregarNotas();
        } else {
          this.snackbar.error('Erro ao excluir nota');
        }
      }
    });
  }
}
