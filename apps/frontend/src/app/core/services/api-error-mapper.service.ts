import { Injectable } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { ApiError } from '../models/error.model';

const FRIENDLY_MESSAGES: Record<number, string> = {
  0:   'Serviço indisponível. Tente novamente mais tarde.',
  408: 'Tempo de requisição esgotado. Tente novamente.',
  429: 'Muitas requisições. Aguarde um momento e tente novamente.',
  502: 'Serviço temporariamente indisponível. Tente novamente.',
  503: 'Serviço em manutenção. Tente novamente mais tarde.',
  504: 'Tempo de resposta esgotado. Tente novamente mais tarde.',
};

@Injectable({ providedIn: 'root' })
export class ApiErrorMapper {
  map(err: any): ApiError {
    if (err instanceof HttpErrorResponse) {
      if (err.status === 0) {
        return {
          code: 'SERVICE_UNAVAILABLE',
          message: FRIENDLY_MESSAGES[0],
          details: null,
          request_id: null,
        } as ApiError;
      }

      if (err.error && typeof err.error === 'object' && !('type' in err.error)) {
        const e = err.error as Partial<ApiError>;
        return {
          code: e.code || String(err.status),
          message: e.message || err.message || 'Erro desconhecido',
          details: e.details || null,
          request_id: e.request_id || null,
        } as ApiError;
      }

      const friendly = FRIENDLY_MESSAGES[err.status];
      return {
        code: String(err.status),
        message: friendly || err.statusText || 'Erro inesperado',
        details: null,
        request_id: null,
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
