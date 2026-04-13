import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { NotaService } from './nota.service';
import { environment } from '../../../environments/environment';

describe('NotaService create', () => {
  let service: NotaService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({ imports: [HttpClientTestingModule], providers: [NotaService] });
    service = TestBed.inject(NotaService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => httpMock.verify());

  it('should post create nota', () => {
    const itens = [{ produto_codigo: 'P-001', quantidade: 2 }];
    service.create(itens).subscribe();
    const req = httpMock.expectOne(`${environment.apiUrl}/api/v1/notas`);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual({ itens });
    req.flush({ id: 50, numero: 101, status: 'ABERTA', itens });
  });
});
