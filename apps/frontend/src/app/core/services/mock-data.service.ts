import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { Produto } from '../models/produto.model';
import { Nota } from '../models/nota.model';

@Injectable({
  providedIn: 'root',
})
export class MockDataService {
  private produtos: Produto[] = [
    { id: 1, codigo: 'PAR-001', descricao: 'Parafuso Sextavado M8', saldo: 150 },
    { id: 2, codigo: 'POR-002', descricao: 'Porca Sextavada M10', saldo: 3 },
    { id: 3, codigo: 'ARR-003', descricao: 'Arruela Lisa 8mm', saldo: 80 },
    { id: 4, codigo: 'PAR-004', descricao: 'Parafuso Allen M6x30', saldo: 200 },
    { id: 5, codigo: 'POR-005', descricao: 'Porca Mariposa M8', saldo: 1 },
    { id: 6, codigo: 'ARR-006', descricao: 'Arruela de Pressão 10mm', saldo: 45 },
    { id: 7, codigo: 'PAR-007', descricao: 'Parafuso Phillips M4x16', saldo: 0 },
    { id: 8, codigo: 'POR-008', descricao: 'Porca Nylon M6', saldo: 120 },
    { id: 9, codigo: 'PAR-009', descricao: 'Parafuso Soberbo M10x40', saldo: 75 },
    { id: 10, codigo: 'POR-010', descricao: 'Porca Cega M12', saldo: 60 },
    { id: 11, codigo: 'ARR-011', descricao: 'Arruela de Encosto 12mm', saldo: 200 },
    { id: 12, codigo: 'PAR-012', descricao: 'Parafuso Autobrocante 4.2x19', saldo: 500 },
    { id: 13, codigo: 'POR-013', descricao: 'Porca Borboleta M6', saldo: 35 },
    { id: 14, codigo: 'ARR-014', descricao: 'Arruela Dentada 8mm', saldo: 90 },
    { id: 15, codigo: 'PAR-015', descricao: 'Parafuso Sextavado M12x50', saldo: 40 },
    { id: 16, codigo: 'POR-016', descricao: 'Porca Allen M8', saldo: 0 },
    { id: 17, codigo: 'ARR-017', descricao: 'Arruela Lisa 6mm', saldo: 150 },
    { id: 18, codigo: 'PAR-018', descricao: 'Parafuso Rosca Madeira 3.5x25', saldo: 300 },
    { id: 19, codigo: 'POR-019', descricao: 'Porca Sextavada M6', saldo: 2 },
    { id: 20, codigo: 'ARR-020', descricao: 'Arruela de Segurança 10mm', saldo: 65 },
    { id: 21, codigo: 'PAR-021', descricao: 'Parafuso Allen M5x20', saldo: 180 },
    { id: 22, codigo: 'POR-022', descricao: 'Porca Prisioneira M8', saldo: 50 },
    { id: 23, codigo: 'ARR-023', descricao: 'Arruela Ondulada 8mm', saldo: 110 },
    { id: 24, codigo: 'PAR-024', descricao: 'Parafuso Philips M3x10', saldo: 400 },
    { id: 25, codigo: 'POR-025', descricao: 'Porca Sextavada M14', saldo: 4 },
    { id: 26, codigo: 'ARR-026', descricao: 'Arruela Tolex 12mm', saldo: 70 },
    { id: 27, codigo: 'PAR-027', descricao: 'Parafuso Canhão M10x60', saldo: 25 },
    { id: 28, codigo: 'POR-028', descricao: 'Porca Auto Travante M10', saldo: 95 },
    { id: 29, codigo: 'BUC-029', descricao: 'Bucha Nylon M6', saldo: 4 },
    { id: 30, codigo: 'ARR-030', descricao: 'Arruela Pressão 6mm', saldo: 1 },
  ];

  private notas: Nota[] = [
    {
      id: 1,
      numero: 1,
      status: 'FECHADA',
      criado_em: '2026-04-08T10:30:00',
      itens: [
        { id: 1, produto_codigo: 'PAR-001', quantidade: 10 },
        { id: 2, produto_codigo: 'POR-002', quantidade: 5 },
      ],
    },
    {
      id: 2,
      numero: 2,
      status: 'FECHADA',
      criado_em: '2026-04-08T14:15:00',
      itens: [{ id: 3, produto_codigo: 'ARR-003', quantidade: 20 }],
    },
    {
      id: 3,
      numero: 3,
      status: 'ABERTA',
      criado_em: '2026-04-09T09:00:00',
      itens: [
        { id: 4, produto_codigo: 'PAR-004', quantidade: 15 },
        { id: 5, produto_codigo: 'POR-005', quantidade: 2 },
        { id: 6, produto_codigo: 'ARR-006', quantidade: 8 },
      ],
    },
    {
      id: 4,
      numero: 4,
      status: 'ABERTA',
      criado_em: '2026-04-09T11:45:00',
      itens: [{ id: 7, produto_codigo: 'POR-008', quantidade: 30 }],
    },
    {
      id: 5,
      numero: 5,
      status: 'FECHADA',
      criado_em: '2026-04-09T16:20:00',
      itens: [
        { id: 8, produto_codigo: 'PAR-001', quantidade: 5 },
        { id: 9, produto_codigo: 'ARR-003', quantidade: 10 },
      ],
    },
    {
      id: 6,
      numero: 6,
      status: 'ABERTA',
      criado_em: '2026-04-10T08:30:00',
      itens: [{ id: 10, produto_codigo: 'PAR-004', quantidade: 25 }],
    },
    {
      id: 7,
      numero: 7,
      status: 'FECHADA',
      criado_em: '2026-04-10T10:00:00',
      itens: [
        { id: 11, produto_codigo: 'PAR-001', quantidade: 20 },
        { id: 12, produto_codigo: 'POR-008', quantidade: 10 },
      ],
    },
    {
      id: 8,
      numero: 8,
      status: 'ABERTA',
      criado_em: '2026-04-10T11:30:00',
      itens: [{ id: 13, produto_codigo: 'ARR-003', quantidade: 15 }],
    },
    {
      id: 9,
      numero: 9,
      status: 'FECHADA',
      criado_em: '2026-04-10T14:00:00',
      itens: [
        { id: 14, produto_codigo: 'PAR-009', quantidade: 8 },
        { id: 15, produto_codigo: 'POR-010', quantidade: 4 },
        { id: 16, produto_codigo: 'ARR-011', quantidade: 12 },
      ],
    },
    {
      id: 10,
      numero: 10,
      status: 'ABERTA',
      criado_em: '2026-04-11T09:15:00',
      itens: [{ id: 17, produto_codigo: 'PAR-012', quantidade: 50 }],
    },
    {
      id: 11,
      numero: 11,
      status: 'FECHADA',
      criado_em: '2026-04-11T13:45:00',
      itens: [
        { id: 18, produto_codigo: 'POR-013', quantidade: 6 },
        { id: 19, produto_codigo: 'ARR-014', quantidade: 20 },
      ],
    },
    {
      id: 12,
      numero: 12,
      status: 'ABERTA',
      criado_em: '2026-04-11T16:00:00',
      itens: [
        { id: 20, produto_codigo: 'PAR-015', quantidade: 10 },
      ],
    },
  ];

  private nextProdutoId = 9;
  private nextNotaId = 13;
  private nextNotaNumero = 13;
  private nextItemId = 20;

  produtos$ = new BehaviorSubject<Produto[]>([...this.produtos]);
  notas$ = new BehaviorSubject<Nota[]>([...this.notas]);

  getProdutos(): Produto[] {
    return [...this.produtos];
  }

  getNotas(): Nota[] {
    return [...this.notas];
  }

  getProdutoByCodigo(codigo: string): Produto | undefined {
    return this.produtos.find((p) => p.codigo === codigo);
  }

  getProdutosBaixoEstoque(limite = 6): Produto[] {
    return this.produtos.filter((p) => p.saldo < limite && p.saldo > 0);
  }

  getUltimasNotas(quantidade = 6): Nota[] {
    return [...this.notas]
      .sort((a, b) => b.numero - a.numero)
      .slice(0, quantidade);
  }

  countByStatus(status: 'ABERTA' | 'FECHADA'): number {
    return this.notas.filter((n) => n.status === status).length;
  }

  addProduto(produto: Omit<Produto, 'id'>): Produto {
    const novo: Produto = { ...produto, id: this.nextProdutoId++ };
    this.produtos.push(novo);
    this.produtos$.next([...this.produtos]);
    return novo;
  }

  addNota(itens: { produto_codigo: string; quantidade: number }[]): Nota {
    const notaItens = itens.map((item) => ({
      id: this.nextItemId++,
      produto_codigo: item.produto_codigo,
      quantidade: item.quantidade,
    }));

    const nota: Nota = {
      id: this.nextNotaId++,
      numero: this.nextNotaNumero++,
      status: 'ABERTA',
      criado_em: new Date().toISOString(),
      itens: notaItens,
    };

    this.notas.push(nota);
    this.notas$.next([...this.notas]);
    return nota;
  }

  imprimirNota(id: number): { success: boolean; error?: string } {
    const nota = this.notas.find((n) => n.id === id);
    if (!nota) return { success: false, error: 'Nota não encontrada' };
    if (nota.status !== 'ABERTA')
      return { success: false, error: 'Nota não está ABERTA' };

    for (const item of nota.itens) {
      const produto = this.produtos.find(
        (p) => p.codigo === item.produto_codigo
      );
      if (!produto)
        return {
          success: false,
          error: `Produto ${item.produto_codigo} não encontrado`,
        };
      if (produto.saldo < item.quantidade)
        return {
          success: false,
          error: `Saldo insuficiente para ${item.produto_codigo}`,
        };
    }

    for (const item of nota.itens) {
      const produto = this.produtos.find(
        (p) => p.codigo === item.produto_codigo
      )!;
      produto.saldo -= item.quantidade;
    }

    nota.status = 'FECHADA';
    this.notas$.next([...this.notas]);
    this.produtos$.next([...this.produtos]);
    return { success: true };
  }

  deleteNota(id: number): boolean {
    const index = this.notas.findIndex((n) => n.id === id);
    if (index === -1) return false;
    this.notas.splice(index, 1);
    this.notas$.next([...this.notas]);
    return true;
  }
}
