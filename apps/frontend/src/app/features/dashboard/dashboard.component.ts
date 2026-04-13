import { Component, DestroyRef, OnInit, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ProdutoService } from '../../core/services/produto.service';
import { NotaService } from '../../core/services/nota.service';
import { Nota } from '../../core/models/nota.model';
import { Produto } from '../../core/models/produto.model';
import { StatusBadgeComponent } from '../../shared/components/status-badge/status-badge.component';
import { PageHeaderComponent } from '../../shared/components/page-header/page-header.component';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    MatIconModule,
    MatTableModule,
    StatusBadgeComponent,
    PageHeaderComponent,
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit {
  private destroyRef = inject(DestroyRef);

  totalProdutos = 0;
  notasAbertas = 0;
  notasFechadas = 0;
  ultimasNotas: Nota[] = [];
  produtosBaixoEstoque: Produto[] = [];
  displayedColumns = ['numero', 'itens', 'status'];
  estoqueColumns = ['codigo', 'descricao', 'saldo'];

  notasError = '';
  estoqueError = '';

  constructor(
    private produtoService: ProdutoService,
    private notaService: NotaService
  ) {}

  ngOnInit(): void {
    this.carregarProdutos();
    this.carregarNotas();
  }

  private carregarProdutos(): void {
    this.produtoService.listAll().pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe({
      next: (produtos) => {
        this.totalProdutos = produtos.length;
        this.estoqueError = '';
      },
      error: () => {
        this.totalProdutos = 0;
        this.estoqueError = 'Não foi possível carregar os produtos';
      },
    });

    this.produtoService.listLowStock(6).pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe({
      next: (produtos) => {
        this.produtosBaixoEstoque = produtos;
        this.estoqueError = '';
      },
      error: () => {
        this.produtosBaixoEstoque = [];
        this.estoqueError = 'Não foi possível carregar os produtos com baixo estoque';
      },
    });
  }

  private carregarNotas(): void {
    this.notaService.listLatest(6).pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe({
      next: (notas) => {
        this.ultimasNotas = notas;
        this.notasAbertas = notas.filter((n) => n.status === 'ABERTA').length;
        this.notasFechadas = notas.filter((n) => n.status === 'FECHADA').length;
        this.notasError = '';
      },
      error: () => {
        this.ultimasNotas = [];
        this.notasAbertas = 0;
        this.notasFechadas = 0;
        this.notasError = 'Não foi possível carregar as notas';
      },
    });
  }

  formatarNumero(num: number): string {
    return num.toString();
  }

  totalItens(nota: Nota): number {
    return nota.itens.reduce((sum, i) => sum + i.quantidade, 0);
  }
}
