import { Component, inject } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { TranslocoDirective, TranslocoService } from '@jsverse/transloco';
import { SelectionModel } from '@angular/cdk/collections';
import {
  MatAnchor,
  MatIconButton,
  MatMiniFabButton,
} from '@angular/material/button';
import {
  MatDrawer,
  MatDrawerContainer,
  MatDrawerContent,
} from '@angular/material/sidenav';
import { MatIcon } from '@angular/material/icon';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { MatTooltip } from '@angular/material/tooltip';
import { combineLatestWith, Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { AsyncPipe, DecimalPipe } from '@angular/common';
import { MatCheckbox } from '@angular/material/checkbox';
import {
  MatExpansionPanel,
  MatExpansionPanelHeader,
  MatExpansionPanelTitle,
} from '@angular/material/expansion';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBar } from '@angular/material/progress-bar';
import { MatFormField, MatLabel } from '@angular/material/form-field';
import { MatOption } from '@angular/material/core';
import { MatSelect } from '@angular/material/select';
import { PaginatorComponent } from '../paginator/paginator.component';
import { BreakpointsService } from '../layout/breakpoints.service';
import { ErrorsService } from '../errors/errors.service';
import { GraphQLModule } from '../graphql/graphql.module';
import { QueueJobsTableComponent } from './queue-jobs-table.component';
import { QueueJobsDatasource } from './queue-jobs.datasource';
import {
  FacetInfo,
  facets,
  orderByOptions,
  QueueJobsController,
  QueueJobsControls,
} from './queue-jobs.controller';

@Component({
  selector: 'app-queue-jobs',
  standalone: true,
  imports: [
    GraphQLModule,
    TranslocoDirective,
    QueueJobsTableComponent,
    MatAnchor,
    MatDrawer,
    MatDrawerContainer,
    MatDrawerContent,
    MatIcon,
    MatIconButton,
    RouterLink,
    RouterLinkActive,
    RouterOutlet,
    MatTooltip,
    AsyncPipe,
    DecimalPipe,
    MatCheckbox,
    MatExpansionPanel,
    MatExpansionPanelHeader,
    MatExpansionPanelTitle,
    MatCardModule,
    MatProgressBar,
    PaginatorComponent,
    MatFormField,
    MatLabel,
    MatOption,
    MatSelect,
    MatMiniFabButton,
  ],
  templateUrl: './queue-jobs.component.html',
  styleUrl: './queue-jobs.component.scss',
})
export class QueueJobsComponent {
  private apollo = inject(Apollo);
  private errorsService = inject(ErrorsService);
  protected breakpoints = inject(BreakpointsService);
  protected transloco = inject(TranslocoService);
  protected controller = new QueueJobsController();
  protected dataSource = new QueueJobsDatasource(
    this.apollo,
    this.errorsService,
    this.controller.variables$,
  );
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  facets$: Observable<FacetInfo<any>[]>;
  protected controls: QueueJobsControls;
  selection = new SelectionModel<string>();

  constructor() {
    this.facets$ = this.controller.controls$.pipe(
      combineLatestWith(this.dataSource.result$),
      map(([controls, result]) =>
        facets.map((f) => ({
          ...f,
          ...f.extractInput(controls.facets),
          aggregations: f
            .extractAggregations(result.aggregations)
            .map((agg) => ({
              ...agg,
              label: f.resolveLabel(agg, this.transloco),
            })),
        })),
      ),
    );
    this.controller.controls$.subscribe((ctrl) => {
      this.controls = ctrl;
    });
  }

  protected readonly orderByOptions = orderByOptions;
}
