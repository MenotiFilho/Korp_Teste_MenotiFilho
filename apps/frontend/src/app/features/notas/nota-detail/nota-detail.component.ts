import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { Subscription } from 'rxjs';
import { DrawerService } from '../../../shared/services/drawer.service';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { NotaService } from '../../../core/services/nota.service';
import { ProdutoService } from '../../../core/services/produto.service';
import { LoadingOverlayComponent } from '../../../shared/components/loading-overlay/loading-overlay.component';
import { StatusBadgeComponent } from '../../../shared/components/status-badge/status-badge.component';
import { Nota } from '../../../core/models/nota.model';
import { NotaItemsTableComponent } from '../../../shared/components/nota-items-table/nota-items-table.component';

@Component({
  selector: 'app-nota-detail',
  standalone: true,
  imports: [MatIconModule, MatButtonModule, StatusBadgeComponent, LoadingOverlayComponent, NotaItemsTableComponent],
  templateUrl: './nota-detail.component.html',
  styleUrls: ['./nota-detail.component.scss'],
})
export class NotaDetailComponent implements OnInit, OnDestroy {
  nota: Nota | null = null;
  loading = false;
  erroCarregamento = '';
  private produtoMapCache = new Map<string, string>();
  private sub!: Subscription;

  constructor(
    private drawer: DrawerService,
    private snackbar: SnackbarService,
    private notaService: NotaService,
    private produtoService: ProdutoService
  ) {}

  ngOnInit(): void {
    this.sub = this.drawer.state$.subscribe((state) => {
      if (state.open && state.component.startsWith('nota-detail-')) {
        const notaId = parseInt(state.component.replace('nota-detail-', ''), 10);
        this.carregarNota(notaId);
      }
    });
  }

  ngOnDestroy(): void {
    this.sub?.unsubscribe();
  }

  private carregarNota(id: number): void {
    this.erroCarregamento = '';
    this.nota = null;

    this.produtoService.listAll().subscribe({
      next: (produtos) => {
        this.produtoMapCache.clear();
        produtos.forEach((p) => this.produtoMapCache.set(p.codigo, p.descricao));
      },
      error: () => {},
    });

    this.notaService.listAll().subscribe({
      next: (notas) => {
        this.nota = notas.find((n) => n.id === id) ?? null;
        if (!this.nota) {
          this.erroCarregamento = 'Nota não encontrada';
        }
      },
      error: () => {
        this.erroCarregamento = 'Não foi possível carregar a nota';
      },
    });
  }

  get produtoMap(): Map<string, string> {
    return this.produtoMapCache;
  }

  get showRemover(): boolean {
    return false;
  }

  formatarNumero(num: number): string {
    return String(num).padStart(5, '0');
  }

  formatarData(iso: string | undefined): string {
    if (!iso) return '';
    const d = new Date(iso);
    return d.toLocaleDateString('pt-BR', {
      day: '2-digit', month: '2-digit', year: 'numeric',
      hour: '2-digit', minute: '2-digit',
    });
  }

  fechar(): void {
    this.drawer.close();
  }

  imprimir(): void {
    if (!this.nota || this.nota.status !== 'ABERTA') return;
    this.loading = true;
    this.notaService.print(this.nota.id).subscribe({
      next: () => {
        this.loading = false;
        this.snackbar.success('Nota impressa com sucesso!');
        this.carregarNota(this.nota!.id);
      },
      error: () => {
        this.loading = false;
        this.snackbar.error('Erro ao imprimir nota');
      },
    });
  }
}
