import { Injectable } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { ErrorModel } from '../models/error.model';

@Injectable({ providedIn: 'root' })
export class ApiErrorMapper {
  map(err: any): ErrorModel {
    if (err instanceof HttpErrorResponse && err.error) {
      const e = err.error as Partial<ErrorModel>;
      return {
        code: e.code || String(err.status),
        message: e.message || err.message || 'Erro desconhecido',
        details: e.details || null,
        request_id: e.request_id || null,
      } as ErrorModel;
    }

    return {
      code: 'UNKNOWN',
      message: err?.message || String(err) || 'Erro desconhecido',
      details: null,
      request_id: null,
    } as ErrorModel;
  }
}
