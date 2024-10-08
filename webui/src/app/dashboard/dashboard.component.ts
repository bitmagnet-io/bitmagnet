import { Component, inject } from '@angular/core';
import { AsyncPipe, DecimalPipe } from '@angular/common';
import { MatCheckbox } from '@angular/material/checkbox';
import { MatDivider } from '@angular/material/divider';
import {
  MatDrawer,
  MatDrawerContainer,
  MatDrawerContent,
} from '@angular/material/sidenav';
import {
  MatExpansionPanel,
  MatExpansionPanelHeader,
  MatExpansionPanelTitle,
} from '@angular/material/expansion';
import { MatFormField, MatLabel } from '@angular/material/form-field';
import { MatIcon } from '@angular/material/icon';
import {
  MatAnchor,
  MatIconButton,
  MatMiniFabButton,
} from '@angular/material/button';
import { MatInput } from '@angular/material/input';
import { MatOption } from '@angular/material/core';
import { MatSelect } from '@angular/material/select';
import { ReactiveFormsModule } from '@angular/forms';
import { TranslocoDirective } from '@jsverse/transloco';
import { MatTooltip } from '@angular/material/tooltip';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { MatMenu, MatMenuItem } from '@angular/material/menu';
import { BreakpointsService } from '../layout/breakpoints.service';
import { TorrentsTableComponent } from '../torrents-table/torrents-table.component';
import { TorrentsBulkActionsComponent } from '../torrents-bulk-actions/torrents-bulk-actions.component';
import { PaginatorComponent } from '../paginator/paginator.component';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    AsyncPipe,
    DecimalPipe,
    MatCheckbox,
    MatDivider,
    MatDrawer,
    MatDrawerContainer,
    MatDrawerContent,
    MatExpansionPanel,
    MatExpansionPanelHeader,
    MatExpansionPanelTitle,
    MatFormField,
    MatIcon,
    MatIconButton,
    MatInput,
    MatLabel,
    MatMiniFabButton,
    MatOption,
    MatSelect,
    PaginatorComponent,
    ReactiveFormsModule,
    TorrentsBulkActionsComponent,
    TorrentsTableComponent,
    TranslocoDirective,
    MatTooltip,
    RouterOutlet,
    MatMenu,
    MatMenuItem,
    RouterLink,
    MatAnchor,
    RouterLinkActive,
    TranslocoDirective,
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss',
})
export class DashboardComponent {
  breakpoints = inject(BreakpointsService);
}