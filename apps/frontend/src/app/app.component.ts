import { Component } from '@angular/core';
import { RouterOutlet, Router, NavigationEnd } from '@angular/router';
import { MatSidenavModule } from '@angular/material/sidenav';
import { SidebarComponent } from './shared/components/sidebar/sidebar.component';
import { TopBarComponent } from './shared/components/top-bar/top-bar.component';
import { DrawerService } from './shared/services/drawer.service';
import { ProdutoFormComponent } from './features/produtos/produto-form/produto-form.component';
import { NotaFormComponent } from './features/notas/nota-form/nota-form.component';
import { NotaDetailComponent } from './features/notas/nota-detail/nota-detail.component';
import { filter, map } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  imports: [
    RouterOutlet,
    MatSidenavModule,
    SidebarComponent,
    TopBarComponent,
    ProdutoFormComponent,
    NotaFormComponent,
    NotaDetailComponent,
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  pageTitle = 'Início';
  drawerOpen = false;
  drawerComponent = '';

  private routeTitles: Record<string, string> = {
    '/inicio': 'Início',
    '/produtos': 'Produtos',
    '/notas': 'Notas Fiscais',
  };

  constructor(
    private router: Router,
    public drawer: DrawerService
  ) {
    this.router.events
      .pipe(
        filter((event) => event instanceof NavigationEnd),
        map((event) => (event as NavigationEnd).urlAfterRedirects)
      )
      .subscribe((url) => {
        this.pageTitle = this.routeTitles[url] || 'Início';
        this.drawer.close();
      });

    this.drawer.state$.subscribe((state) => {
      this.drawerOpen = state.open;
      this.drawerComponent = state.component;
    });
  }

  onDrawerClosed(): void {
    this.drawer.close();
  }
}
