import { Component, Input, Output, EventEmitter } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { Produto } from '../../../core/models/produto.model';

export interface NotaItemAdicionado {
  produto_codigo: string;
  quantidade: number;
}

@Component({
  selector: 'app-nota-item-adder',
  standalone: true,
  imports: [
    FormsModule,
    MatIconModule,
    MatSelectModule,
    MatFormFieldModule,
    MatInputModule,
  ],
  templateUrl: './nota-item-adder.component.html',
  styleUrls: ['./nota-item-adder.component.scss'],
})
export class NotaItemAdderComponent {
  @Input() produtos: Produto[] = [];
  @Output() itemAdded = new EventEmitter<NotaItemAdicionado>();

  produtoSelecionado: Produto | null = null;
  quantidade = 1;

  adicionar(): void {
    if (!this.produtoSelecionado || this.quantidade < 1) return;
    this.itemAdded.emit({
      produto_codigo: this.produtoSelecionado.codigo,
      quantidade: this.quantidade,
    });
    this.produtoSelecionado = null;
    this.quantidade = 1;
  }
}
