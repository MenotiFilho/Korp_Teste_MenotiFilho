import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { Subscription } from 'rxjs';
import { DrawerService } from '../../../shared/services/drawer.service';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { MockDataService } from '../../../core/services/mock-data.service';
import { LoadingOverlayComponent } from '../../../shared/components/loading-overlay/loading-overlay.component';
import { StatusBadgeComponent } from '../../../shared/components/status-badge/status-badge.component';
import { Nota } from '../../../core/models/nota.model';
import { NotaItemsTableComponent } from '../../../shared/components/nota-items-table/nota-items-table.component';

@Component({
  selector: 'app-nota-detail',
  imports: [MatIconModule, MatButtonModule, StatusBadgeComponent, LoadingOverlayComponent, NotaItemsTableComponent],
  templateUrl: './nota-detail.component.html',
  styleUrl: './nota-detail.component.scss',
})
export class NotaDetailComponent implements OnInit, OnDestroy {
  nota: Nota | null = null;
  loading = false;
  private sub!: Subscription;

  constructor(
    private drawer: DrawerService,
    private snackbar: SnackbarService,
    private mockData: MockDataService
  ) {}

  ngOnInit(): void {
    this.sub = this.drawer.state$.subscribe((state) => {
      if (state.open && state.component.startsWith('nota-detail-')) {
        const notaId = parseInt(state.component.replace('nota-detail-', ''), 10);
        this.nota = this.mockData.getNotas().find((n) => n.id === notaId) ?? null;
      }
    });
  }

  ngOnDestroy(): void {
    this.sub?.unsubscribe();
  }

  get produtoMap(): Map<string, string> {
    const map = new Map<string, string>();
    this.mockData.getProdutos().forEach((p) => map.set(p.codigo, p.descricao));
    return map;
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
    setTimeout(() => {
      const result = this.mockData.imprimirNota(this.nota!.id);
      this.loading = false;
      if (result.success) {
        this.snackbar.success('Nota impressa com sucesso!');
        this.nota = this.mockData.getNotas().find((n) => n.id === this.nota!.id) ?? null;
      } else {
        this.snackbar.error(result.error || 'Erro ao imprimir nota');
      }
    }, 1500);
  }
}
