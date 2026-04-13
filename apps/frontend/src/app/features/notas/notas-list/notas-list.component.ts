import { Component } from '@angular/core';

@Component({
  selector: 'app-notas-list',
  template: `
    <h1 class="page-title">Notas Fiscais</h1>
    <p class="page-subtitle">Gerencie as notas fiscais emitidas</p>
  `,
  styles: [`
    :host {
      display: block;
    }
    .page-title {
      font-family: Inter, sans-serif;
      font-size: 28px;
      font-weight: 500;
      color: var(--color-on-surface);
      margin: 0 0 var(--spacing-xs) 0;
    }
    .page-subtitle {
      font-family: Inter, sans-serif;
      font-size: 14px;
      color: var(--color-on-surface-variant);
      margin: 0 0 var(--spacing-lg) 0;
    }
  `],
})
export class NotasListComponent {}
