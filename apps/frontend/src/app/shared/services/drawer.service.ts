import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export interface DrawerState {
  open: boolean;
  width: number;
  component: string;
}

@Injectable({
  providedIn: 'root',
})
export class DrawerService {
  private state = new BehaviorSubject<DrawerState>({
    open: false,
    width: 424,
    component: '',
  });

  state$ = this.state.asObservable();

  open(component: string, width = 424): void {
    this.state.next({ open: true, width, component });
  }

  close(): void {
    this.state.next({ ...this.state.value, open: false });
  }

  get isOpen(): boolean {
    return this.state.value.open;
  }
}
