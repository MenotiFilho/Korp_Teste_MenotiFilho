import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'inicio',
    pathMatch: 'full',
  },
  {
    path: 'inicio',
    loadComponent: () =>
      import('./features/dashboard/dashboard.component').then(
        (m) => m.DashboardComponent
      ),
    title: 'Início',
  },
  {
    path: 'produtos',
    loadComponent: () =>
      import(
        './features/produtos/produtos-list/produtos-list.component'
      ).then((m) => m.ProdutosListComponent),
    title: 'Produtos',
  },
  {
    path: 'notas',
    loadComponent: () =>
      import('./features/notas/notas-list/notas-list.component').then(
        (m) => m.NotasListComponent
      ),
    title: 'Notas Fiscais',
  },
];
