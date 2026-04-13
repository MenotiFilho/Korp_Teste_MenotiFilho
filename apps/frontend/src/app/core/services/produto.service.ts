import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Produto } from '../models/produto.model';
import { environment } from '../../../environments/environment';

@Injectable({ providedIn: 'root' })
export class ProdutoService {
  private base = environment.estoqueUrl;

  constructor(private http: HttpClient) {}

  listAll(): Observable<Produto[]> {
    return this.http.get<Produto[]>(`${this.base}/api/v1/produtos`);
  }

  listLowStock(limit = 6): Observable<Produto[]> {
    return this.http.get<Produto[]>(`${this.base}/api/v1/produtos/baixo-estoque?limit=${limit}`);
  }

  create(prod: Omit<Produto, 'id'>) {
    return this.http.post<Produto>(`${this.base}/api/v1/produtos`, prod, {
      headers: new HttpHeaders({ 'Content-Type': 'application/json' }),
    });
  }

  update(id: number, prod: Partial<Produto>) {
    return this.http.put<Produto>(`${this.base}/api/v1/produtos/${id}`, prod, {
      headers: new HttpHeaders({ 'Content-Type': 'application/json' }),
    });
  }

  delete(id: number) {
    return this.http.delete<void>(`${this.base}/api/v1/produtos/${id}`);
  }
}
