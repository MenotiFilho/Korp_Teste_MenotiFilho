import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { NotaService } from './nota.service';
import { environment } from '../../../environments/environment';

describe('NotaService print', () => {
  let service: NotaService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({ imports: [HttpClientTestingModule], providers: [NotaService] });
    service = TestBed.inject(NotaService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => httpMock.verify());

  it('should post print with idempotency header', () => {
    service.print(123).subscribe();
    const req = httpMock.expectOne(`${environment.faturamentoUrl}/api/v1/notas/123/imprimir`);
    expect(req.request.method).toBe('POST');
    expect(req.request.headers.get('Idempotency-Key')).toBe('invoice-print-123');
    req.flush(null);
  });
});
