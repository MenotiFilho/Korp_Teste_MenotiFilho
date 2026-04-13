import { Component, OnInit, OnDestroy } from '@angular/core';
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
import { Subscription } from 'rxjs';
import { DrawerService } from '../../../shared/services/drawer.service';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { MockDataService } from '../../../core/services/mock-data.service';
import { ProdutoService } from '../../../core/services/produto.service';
import { Produto } from '../../../core/models/produto.model';

@Component({
  selector: 'app-produto-form',
  imports: [
    ReactiveFormsModule,
    MatIconModule,
    MatButtonModule,
    MatInputModule,
    MatFormFieldModule,
  ],
  templateUrl: './produto-form.component.html',
  styleUrl: './produto-form.component.scss',
})
export class ProdutoFormComponent implements OnInit, OnDestroy {
  form: FormGroup;
  editMode = false;
  produtoEdicao: Produto | null = null;
  private sub!: Subscription;

  constructor(
    private fb: FormBuilder,
    private drawer: DrawerService,
    private snackbar: SnackbarService,
    private mockData: MockDataService,
    private produtoService: ProdutoService
  ) {
    this.form = this.fb.group({
      codigo: ['', [Validators.required, Validators.maxLength(20)]],
      descricao: ['', [Validators.required, Validators.maxLength(100)]],
      saldo: [0, [Validators.required, Validators.min(0)]],
    });
  }

  ngOnInit(): void {
    this.sub = this.drawer.state$.subscribe((state) => {
      if (state.open) {
        if (state.component.startsWith('produto-edit-')) {
          const id = parseInt(state.component.replace('produto-edit-', ''), 10);
          const produto = this.mockData.getProdutos().find((p) => p.id === id);
          if (produto) {
            this.editMode = true;
            this.produtoEdicao = produto;
            this.form.patchValue(produto);
          }
        } else if (state.component === 'produto-form') {
          this.editMode = false;
          this.produtoEdicao = null;
          this.form.reset({ codigo: '', descricao: '', saldo: 0 });
        }
      }
    });
  }

  ngOnDestroy(): void {
    this.sub?.unsubscribe();
  }

  fechar(): void {
    this.form.reset({ codigo: '', descricao: '', saldo: 0 });
    this.editMode = false;
    this.produtoEdicao = null;
    this.drawer.close();
  }

  salvar(): void {
    if (this.form.valid) {
      if (this.editMode && this.produtoEdicao) {
        this.snackbar.success('Produto atualizado com sucesso!');
      } else {
        const payload = this.form.value;
        this.produtoService.create(payload).subscribe({
          next: () => this.snackbar.success('Produto cadastrado com sucesso!'),
          error: () => {
            // fallback to mock when backend is unavailable
            this.mockData.addProduto(payload);
            this.snackbar.success('Produto cadastrado (modo offline)');
          }
        });
      }
      this.fechar();
    }
  }
}
