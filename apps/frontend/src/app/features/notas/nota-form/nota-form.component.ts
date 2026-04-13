import { Component, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
} from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { DrawerService } from '../../../shared/services/drawer.service';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { MockDataService } from '../../../core/services/mock-data.service';
import { NotaService } from '../../../core/services/nota.service';
import { Produto } from '../../../core/models/produto.model';
import { NotaItemsTableComponent, NotaItem } from '../../../shared/components/nota-items-table/nota-items-table.component';

interface ItemNota {
  id?: number;
  produto_codigo: string;
  produto_descricao: string;
  quantidade: number;
}

@Component({
  selector: 'app-nota-form',
  imports: [
    ReactiveFormsModule,
    MatIconModule,
    MatButtonModule,
    MatInputModule,
    MatFormFieldModule,
    MatSelectModule,
    NotaItemsTableComponent,
  ],
  templateUrl: './nota-form.component.html',
  styleUrl: './nota-form.component.scss',
})
export class NotaFormComponent implements OnInit {
  form: FormGroup;
  produtos: Produto[] = [];
  itens: ItemNota[] = [];
  produtoMap = new Map<string, string>();
  displayedColumns = ['codigo', 'descricao', 'quantidade', 'remover'];

  constructor(
    private fb: FormBuilder,
    private drawer: DrawerService,
    private snackbar: SnackbarService,
    private mockData: MockDataService,
    private notaService: NotaService
  ) {
    this.form = this.fb.group({
      produto: [null, Validators.required],
      quantidade: [1, [Validators.required, Validators.min(1)]],
    });
  }

  ngOnInit(): void {
    this.produtos = this.mockData.getProdutos().filter((p) => p.saldo > 0);
    this.produtos.forEach(p => this.produtoMap.set(p.codigo, p.descricao));
  }

  get produtoSelecionado(): Produto | null {
    return this.form.get('produto')?.value ?? null;
  }

  compareProdutos(p1: Produto, p2: Produto): boolean {
    return p1 && p2 ? p1.id === p2.id : p1 === p2;
  }

  get itensParaTabela(): NotaItem[] {
    return this.itens.map((item, index) => ({
      id: index,
      produto_codigo: item.produto_codigo,
      quantidade: item.quantidade,
    }));
  }

  adicionarItem(): void {
    if (this.form.valid && this.produtoSelecionado) {
      const produto = this.produtoSelecionado;
      const quantidade = this.form.get('quantidade')?.value;

      if (quantidade > produto.saldo) {
        this.snackbar.error(
          `Saldo insuficiente. Disponível: ${produto.saldo}`
        );
        return;
      }

      this.itens.push({
        produto_codigo: produto.codigo,
        produto_descricao: produto.descricao,
        quantidade,
      });

      this.form.patchValue({ produto: null, quantidade: 1 });
    }
  }

  removerItem(index: number): void {
    this.itens.splice(index, 1);
  }

  fechar(): void {
    this.itens = [];
    this.form.reset({ produto: null, quantidade: 1 });
    this.drawer.close();
  }

  salvar(): void {
    if (this.itens.length === 0) {
      this.snackbar.error('Adicione pelo menos um item à nota');
      return;
    }

    const payload = this.itens.map((i) => ({ produto_codigo: i.produto_codigo, quantidade: i.quantidade }));
    this.notaService.create(payload).subscribe({
      next: () => {
        this.snackbar.success('Nota fiscal criada com sucesso!');
        this.fechar();
      },
      error: (err) => {
        const mapped = this.apiErrorMapper.map(err);
        this.snackbar.error('Falha ao criar nota: ' + mapped.message);
      }
    });
  }
}
