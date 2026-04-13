import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Nota } from '../models/nota.model';
import { environment } from '../../../environments/environment';

@Injectable({ providedIn: 'root' })
export class NotaService {
  private base = environment.apiUrl;

  constructor(private http: HttpClient) {}

  listLatest(limit = 6): Observable<Nota[]> {
    return this.http.get<Nota[]>(`${this.base}/api/v1/notas/ultimas?limit=${limit}`);
  }
}
