import { Component, OnInit } from '@angular/core';
import { RouterLink } from '@angular/router';
import { DecimalPipe } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MockDataService } from '../../core/services/mock-data.service';
import { Nota } from '../../core/models/nota.model';
import { Produto } from '../../core/models/produto.model';
import { StatusBadgeComponent } from '../../shared/components/status-badge/status-badge.component';
import { PageHeaderComponent } from '../../shared/components/page-header/page-header.component';

@Component({
  selector: 'app-dashboard',
  imports: [
    RouterLink,
    DecimalPipe,
    MatIconModule,
    MatTableModule,
    StatusBadgeComponent,
    PageHeaderComponent,
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss',
})
export class DashboardComponent implements OnInit {
  totalProdutos = 0;
  notasAbertas = 0;
  notasFechadas = 0;
  ultimasNotas: Nota[] = [];
  produtosBaixoEstoque: Produto[] = [];
  displayedColumns = ['numero', 'itens', 'criado_em', 'status'];
  estoqueColumns = ['codigo', 'descricao', 'saldo'];

  constructor(private mockData: MockDataService) {}

  ngOnInit(): void {
    const produtos = this.mockData.getProdutos();
    this.totalProdutos = produtos.length;
    this.notasAbertas = this.mockData.countByStatus('ABERTA');
    this.notasFechadas = this.mockData.countByStatus('FECHADA');
    this.ultimasNotas = this.mockData.getUltimasNotas(6);
    this.produtosBaixoEstoque = this.mockData.getProdutosBaixoEstoque(6);
  }

  formatarData(iso: string | undefined): string {
    if (!iso) return '';
    const d = new Date(iso);
    return d.toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    });
  }
}
