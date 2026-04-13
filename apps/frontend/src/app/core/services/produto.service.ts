import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Produto } from '../models/produto.model';
import { environment } from '../../../environments/environment';

@Injectable({ providedIn: 'root' })
export class ProdutoService {
  private base = environment.apiUrl;

  constructor(private http: HttpClient) {}

  listAll(): Observable<Produto[]> {
    return this.http.get<Produto[]>(`${this.base}/api/v1/produtos`);
  }

  listLowStock(limit = 6): Observable<Produto[]> {
    return this.http.get<Produto[]>(`${this.base}/api/v1/produtos/baixo-estoque?limit=${limit}`);
  }
}
