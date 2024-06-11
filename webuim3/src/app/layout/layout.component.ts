import { Component, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import {Route, RouterLink, RouterLinkActive, RouterOutlet} from "@angular/router";
import {routes} from "../app.routes";
import {MatTab, MatTabGroup, MatTabLink, MatTabNav, MatTabNavPanel} from "@angular/material/tabs";
import {ThemeManager} from "../themes/theme-manager.service";
import {MatMenu, MatMenuItem, MatMenuTrigger} from "@angular/material/menu";
import {MatTooltip} from "@angular/material/tooltip";
import {HttpClientModule} from "@angular/common/http";

@Component({
  selector: 'app-layout',
  templateUrl: './layout.component.html',
  styleUrl: './layout.component.scss',
  standalone: true,
  imports: [
    MatToolbarModule,
    MatButtonModule,
    MatSidenavModule,
    MatListModule,
    MatIconModule,
    AsyncPipe,
    RouterLink,
    RouterLinkActive,
    MatTabGroup,
    MatTab,
    MatTabNavPanel,
    RouterOutlet,
    MatTabNav,
    MatTabLink,
    MatMenuTrigger,
    MatMenuItem,
    MatMenu,
    MatTooltip
  ],
  providers: [ HttpClientModule]
})
export class LayoutComponent {
  rootRoutes = routes.filter(r=>r.path);
  themeManager = inject(ThemeManager);

  routeTitle(route: Route): string {
    return route.title as string;
  }
}
