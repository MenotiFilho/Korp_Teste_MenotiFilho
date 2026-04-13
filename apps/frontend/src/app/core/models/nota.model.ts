export interface NotaItem {
  id: number;
  produto_codigo: string;
  quantidade: number;
}

export interface Nota {
  id: number;
  numero: number;
  status: 'ABERTA' | 'FECHADA';
  itens: NotaItem[];
  criado_em?: string;
}
