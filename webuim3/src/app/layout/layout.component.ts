import { Component, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import {
  RouterLink,
  RouterLinkActive,
  RouterOutlet,
} from '@angular/router';
import {
  MatTab,
  MatTabGroup,
  MatTabLink,
  MatTabNav,
  MatTabNavPanel,
} from '@angular/material/tabs';
import { MatMenu, MatMenuItem, MatMenuTrigger } from '@angular/material/menu';
import { MatTooltip } from '@angular/material/tooltip';
import {
  TranslocoDirective,
  TranslocoPipe,
} from '@jsverse/transloco';
import { Observable } from 'rxjs';
import { ThemeManager } from '../themes/theme-manager.service';
import { VersionComponent } from '../version/version.component';
import {TranslateManager} from "../i18n/translate-manager.service";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";
import {map, shareReplay} from "rxjs/operators";
import {BreakpointsService} from "./breakpoints.service";

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
    MatTooltip,
    VersionComponent,
    TranslocoDirective,
    TranslocoPipe,
  ],
})
export class LayoutComponent {
  themeManager = inject(ThemeManager);
  translateManager = inject(TranslateManager);
  breakpoints = inject(BreakpointsService)
}
