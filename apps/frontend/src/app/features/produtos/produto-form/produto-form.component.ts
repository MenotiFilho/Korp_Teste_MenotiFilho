import { Component } from '@angular/core';
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
import { DrawerService } from '../../../shared/services/drawer.service';
import { SnackbarService } from '../../../shared/services/snackbar.service';
import { MockDataService } from '../../../core/services/mock-data.service';
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
export class ProdutoFormComponent {
  form: FormGroup;
  editMode = false;
  produtoEdicao: Produto | null = null;

  constructor(
    private fb: FormBuilder,
    private drawer: DrawerService,
    private snackbar: SnackbarService,
    private mockData: MockDataService
  ) {
    this.form = this.fb.group({
      codigo: ['', [Validators.required, Validators.maxLength(20)]],
      descricao: ['', [Validators.required, Validators.maxLength(100)]],
      saldo: [0, [Validators.required, Validators.min(0)]],
    });
  }

  setProduto(produto: Produto): void {
    this.editMode = true;
    this.produtoEdicao = produto;
    this.form.patchValue(produto);
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
        this.mockData.addProduto(this.form.value);
        this.snackbar.success('Produto cadastrado com sucesso!');
      }
      this.fechar();
    }
  }
}
