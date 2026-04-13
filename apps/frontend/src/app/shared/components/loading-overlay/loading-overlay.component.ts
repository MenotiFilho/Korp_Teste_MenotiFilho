import { Component, Input } from '@angular/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@Component({
  selector: 'app-loading-overlay',
  standalone: true,
  imports: [MatProgressSpinnerModule],
  templateUrl: './loading-overlay.component.html',
  styleUrls: ['./loading-overlay.component.scss'],
})
export class LoadingOverlayComponent {
  @Input() message = 'Processando...';
  private _visible = false;

  @Input()
  set visible(value: boolean) {
    if (!value && this._visible) {
      setTimeout(() => {
        this._visible = false;
      }, 1000);
    } else {
      this._visible = value;
    }
  }

  get visible(): boolean {
    return this._visible;
  }
}