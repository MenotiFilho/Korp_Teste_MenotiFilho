import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { ProdutoService } from './produto.service';
import { environment } from '../../../environments/environment';

describe('ProdutoService', () => {
  let service: ProdutoService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({ imports: [HttpClientTestingModule], providers: [ProdutoService] });
    service = TestBed.inject(ProdutoService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => httpMock.verify());

  it('should request listAll', () => {
    service.listAll().subscribe();
    const req = httpMock.expectOne(`${environment.estoqueUrl}/api/v1/produtos`);
    expect(req.request.method).toBe('GET');
    req.flush([]);
  });

  it('should request listLowStock with limit', () => {
    service.listLowStock(5).subscribe();
    const req = httpMock.expectOne(`${environment.estoqueUrl}/api/v1/produtos/baixo-estoque?limit=5`);
    expect(req.request.method).toBe('GET');
    req.flush([]);
  });
});
