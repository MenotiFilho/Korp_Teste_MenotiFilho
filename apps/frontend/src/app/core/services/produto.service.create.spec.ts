import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { ProdutoService } from './produto.service';
import { environment } from '../../../environments/environment';

describe('ProdutoService create', () => {
  let service: ProdutoService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({ imports: [HttpClientTestingModule], providers: [ProdutoService] });
    service = TestBed.inject(ProdutoService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => httpMock.verify());

  it('should post create produto', () => {
    const payload = { codigo: 'P-999', descricao: 'Teste', saldo: 5 };
    service.create(payload as any).subscribe();
    const req = httpMock.expectOne(`${environment.apiUrl}/api/v1/produtos`);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(payload);
    req.flush({ id: 100, ...payload });
  });
});
