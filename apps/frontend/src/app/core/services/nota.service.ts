import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Nota } from '../models/nota.model';
import { environment } from '../../../environments/environment';

@Injectable({ providedIn: 'root' })
export class NotaService {
  private base = environment.faturamentoUrl;

  constructor(private http: HttpClient) {}

  listLatest(limit = 6): Observable<Nota[]> {
    return this.http.get<Nota[]>(`${this.base}/api/v1/notas/ultimas?limit=${limit}`);
  }

  listAll(): Observable<Nota[]> {
    return this.http.get<Nota[]>(`${this.base}/api/v1/notas`);
  }

  create(itens: { produto_codigo: string; quantidade: number }[]) {
    return this.http.post<Nota>(`${this.base}/api/v1/notas`, { itens }, {
      headers: new HttpHeaders({ 'Content-Type': 'application/json' }),
    });
  }

  update(id: number, itens: { produto_codigo: string; quantidade: number }[]) {
    return this.http.put<Nota>(`${this.base}/api/v1/notas/${id}`, { itens }, {
      headers: new HttpHeaders({ 'Content-Type': 'application/json' }),
    });
  }

  print(id: number) {
    const url = `${this.base}/api/v1/notas/${id}/imprimir`;
    const idempotency = `invoice-print-${id}`;
    const headers = new HttpHeaders({ 'Content-Type': 'application/json', 'Idempotency-Key': idempotency });
    return this.http.post<void>(url, null, { headers });
  }

  delete(id: number) {
    return this.http.delete<void>(`${this.base}/api/v1/notas/${id}`);
  }
}
