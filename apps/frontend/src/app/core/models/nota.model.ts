export interface NotaItem {
  id: number;
  produto_codigo: string;
  quantidade: number;
}

export interface Nota {
  id: number;
  numero: number;
  status: 'ABERTA' | 'FECHADA';
  // timestamp when the nota was created (optional in some fixtures)
  criado_em?: string;
  itens: NotaItem[];
}
