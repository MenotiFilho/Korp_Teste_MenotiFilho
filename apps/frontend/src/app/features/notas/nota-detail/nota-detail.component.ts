import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { Subscription } from 'rxjs';
import { DrawerService } from '../../../shared/services/drawer.service';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { NotaService } from '../../../core/services/nota.service';
import { ProdutoService } from '../../../core/services/produto.service';
import { LoadingOverlayComponent } from '../../../shared/components/loading-overlay/loading-overlay.component';
import { StatusBadgeComponent } from '../../../shared/components/status-badge/status-badge.component';
import { NotaItemAdderComponent, NotaItemAdicionado } from '../../../shared/components/nota-item-adder/nota-item-adder.component';
import { Nota } from '../../../core/models/nota.model';
import { Produto } from '../../../core/models/produto.model';
import { ApiErrorMapper } from '../../../core/services/api-error-mapper.service';

interface EditableItem {
  produto_codigo: string;
  quantidade: number;
}

@Component({
  selector: 'app-nota-detail',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatIconModule,
    MatButtonModule,
    StatusBadgeComponent,
    LoadingOverlayComponent,
    NotaItemAdderComponent,
  ],
  templateUrl: './nota-detail.component.html',
  styleUrls: ['./nota-detail.component.scss'],
})
export class NotaDetailComponent implements OnInit, OnDestroy {
  nota: Nota | null = null;
  loading = false;
  saving = false;
  erroCarregamento = '';
  editando = false;

  produtos: Produto[] = [];
  itensEditaveis: EditableItem[] = [];
  produtoMap = new Map<string, string>();

  private sub!: Subscription;

  constructor(
    private drawer: DrawerService,
    private snackbar: SnackbarService,
    private notaService: NotaService,
    private produtoService: ProdutoService,
    private apiErrorMapper: ApiErrorMapper
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
    this.editando = false;

    this.produtoService.listAll().subscribe({
      next: (produtos) => {
        this.produtos = produtos;
        this.produtoMap.clear();
        produtos.forEach((p) => this.produtoMap.set(p.codigo, p.descricao));
      },
      error: () => {},
    });

    this.notaService.listAll().subscribe({
      next: (notas) => {
        this.nota = notas.find((n) => n.id === id) ?? null;
        if (this.nota) {
          this.editando = this.nota.status === 'ABERTA';
          this.itensEditaveis = this.nota.itens.map((i) => ({
            produto_codigo: i.produto_codigo,
            quantidade: i.quantidade,
          }));
        } else {
          this.erroCarregamento = 'Nota não encontrada';
        }
      },
      error: () => {
        this.erroCarregamento = 'Não foi possível carregar a nota';
      },
    });
  }

  getNomeProduto(codigo: string): string {
    return this.produtoMap.get(codigo) ?? codigo;
  }

  getTotalItens(): number {
    return this.itensEditaveis.reduce((sum, i) => sum + i.quantidade, 0);
  }

  formatarNumero(num: number): string {
    return String(num).padStart(5, '0');
  }

  onItemAdicionado(event: NotaItemAdicionado): void {
    const existente = this.itensEditaveis.find(
      (i) => i.produto_codigo === event.produto_codigo
    );
    if (existente) {
      existente.quantidade += event.quantidade;
    } else {
      this.itensEditaveis.push({
        produto_codigo: event.produto_codigo,
        quantidade: event.quantidade,
      });
    }
  }

  removerItem(index: number): void {
    this.itensEditaveis.splice(index, 1);
  }

  salvar(): void {
    if (!this.nota || this.itensEditaveis.length === 0) return;
    this.saving = true;
    const payload = this.itensEditaveis.map((i) => ({
      produto_codigo: i.produto_codigo,
      quantidade: i.quantidade,
    }));
    this.notaService.update(this.nota.id, payload).subscribe({
      next: () => {
        this.saving = false;
        this.snackbar.success('Nota atualizada com sucesso!');
        this.carregarNota(this.nota!.id);
      },
      error: (err) => {
        this.saving = false;
        const mapped = this.apiErrorMapper.map(err);
        this.snackbar.error('Falha ao atualizar nota: ' + mapped.message);
      },
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
