import { Component, Input } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';

export interface NotaItem {
  id: number;
  produto_codigo: string;
  quantidade: number;
}

@Component({
  selector: 'app-nota-items-table',
  standalone: true,
  imports: [CommonModule, MatTableModule],
  templateUrl: './nota-items-table.component.html',
  styleUrl: './nota-items-table.component.scss',
})
export class NotaItemsTableComponent {
  @Input() itens: NotaItem[] = [];
  @Input() produtoMap: Map<string, string> = new Map();

  displayedColumns = ['codigo', 'produto', 'quantidade', 'remover'];

  getNomeProduto(codigo: string): string {
    return this.produtoMap.get(codigo) ?? codigo;
  }

  getTotalItens(): number {
    return this.itens.reduce((sum, i) => sum + i.quantidade, 0);
  }

  hasRemover(): boolean {
    return this.displayedColumns.includes('remover');
  }
}