import { Component, inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { TranslocoService } from "@jsverse/transloco";
import { SelectionModel } from "@angular/cdk/collections";
import { combineLatestWith, Observable } from "rxjs";
import { map } from "rxjs/operators";
import { PaginatorComponent } from "../../paginator/paginator.component";
import { BreakpointsService } from "../../layout/breakpoints.service";
import { ErrorsService } from "../../errors/errors.service";
import { AppModule } from "../../app.module";
import { QueueJobsTableComponent } from "./queue-jobs-table.component";
import { QueueJobsDatasource } from "./queue-jobs.datasource";
import {
  FacetInfo,
  facets,
  orderByOptions,
  QueueJobsController,
  QueueJobsControls,
} from "./queue-jobs.controller";

@Component({
  selector: "app-queue-jobs",
  standalone: true,
  imports: [AppModule, PaginatorComponent, QueueJobsTableComponent],
  templateUrl: "./queue-jobs.component.html",
  styleUrl: "./queue-jobs.component.scss",
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
