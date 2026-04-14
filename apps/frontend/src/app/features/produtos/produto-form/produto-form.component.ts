import { Component, DestroyRef, OnInit, inject } from '@angular/core';
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
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { DrawerService } from '../../../shared/services/drawer.service';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { ProdutoService } from '../../../core/services/produto.service';
import { Produto } from '../../../core/models/produto.model';
import { ApiErrorMapper } from '../../../core/services/api-error-mapper.service';

@Component({
  selector: 'app-produto-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatIconModule,
    MatButtonModule,
    MatInputModule,
    MatFormFieldModule,
  ],
  templateUrl: './produto-form.component.html',
  styleUrls: ['./produto-form.component.scss'],
})
export class ProdutoFormComponent implements OnInit {
  private destroyRef = inject(DestroyRef);

  form: FormGroup;
  editMode = false;
  produtoEdicao: Produto | null = null;

  constructor(
    private fb: FormBuilder,
    private drawer: DrawerService,
    private snackbar: SnackbarService,
    private produtoService: ProdutoService,
    private apiErrorMapper: ApiErrorMapper
  ) {
    this.form = this.fb.group({
      codigo: ['', [Validators.required, Validators.maxLength(20)]],
      descricao: ['', [Validators.required, Validators.maxLength(100)]],
      saldo: [0, [Validators.required, Validators.min(0)]],
    });
  }

  ngOnInit(): void {
    this.drawer.state$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((state) => {
      if (state.open) {
        if (state.component.startsWith('produto-edit-')) {
          const id = parseInt(state.component.replace('produto-edit-', ''), 10);
          this.produtoService.listAll().subscribe({
            next: (produtos) => {
              const produto = produtos.find((p) => p.id === id);
              if (produto) {
                this.editMode = true;
                this.produtoEdicao = produto;
                this.form.patchValue(produto);
              }
            },
            error: () => {
              this.snackbar.error('Não foi possível carregar o produto');
            },
          });
        } else if (state.component === 'produto-form') {
          this.editMode = false;
          this.produtoEdicao = null;
          this.form.reset({ codigo: '', descricao: '', saldo: 0 });
        }
      }
    });
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
        const { descricao, saldo } = this.form.value;
        this.produtoService.update(this.produtoEdicao.id, { descricao, saldo }).subscribe({
          next: () => {
            this.snackbar.success('Produto atualizado com sucesso!');
            this.fechar();
          },
          error: (err) => {
            const mapped = this.apiErrorMapper.map(err);
            this.snackbar.error(`Falha ao atualizar produto: ${mapped.message}`);
          },
        });
        return;
      }

      const payload = this.form.value;
      this.produtoService.create(payload).subscribe({
        next: () => {
          this.snackbar.success('Produto cadastrado com sucesso!');
          this.fechar();
        },
        error: (err) => {
          const mapped = this.apiErrorMapper.map(err);
          this.snackbar.error(`Falha ao cadastrar produto: ${mapped.message}`);
        },
      });
    }
  }
}
