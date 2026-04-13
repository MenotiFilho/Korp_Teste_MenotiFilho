import { Injectable } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { ApiError } from '../models/error.model';

@Injectable({ providedIn: 'root' })
export class ApiErrorMapper {
  map(err: any): ApiError {
    if (err instanceof HttpErrorResponse && err.error) {
      const e = err.error as Partial<ApiError>;
      return {
        code: e.code || String(err.status),
        message: e.message || err.message || 'Erro desconhecido',
        details: e.details || null,
        request_id: e.request_id || null,
      } as ApiError;
    }

    return {
      code: 'UNKNOWN',
      message: err?.message || String(err) || 'Erro desconhecido',
      details: null,
      request_id: null,
    } as ApiError;
  }
}
