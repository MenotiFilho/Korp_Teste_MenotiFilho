import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { NotaService } from './nota.service';
import { environment } from '../../../environments/environment';

describe('NotaService', () => {
  let service: NotaService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({ imports: [HttpClientTestingModule], providers: [NotaService] });
    service = TestBed.inject(NotaService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => httpMock.verify());

  it('should request listLatest with default limit', () => {
    service.listLatest().subscribe();
    const req = httpMock.expectOne(`${environment.faturamentoUrl}/api/v1/notas/ultimas?limit=6`);
    expect(req.request.method).toBe('GET');
    req.flush([]);
  });
});
