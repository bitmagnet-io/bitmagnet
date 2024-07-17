import * as generated from "../graphql/generated";
import {BehaviorSubject, debounceTime, Observable} from "rxjs";

type OrderBySelection = {
  field: generated.QueueJobsOrderByField;
  descending: boolean;
};

export type QueueJobsControls = {
  limit: number;
  page: number;
  queues?: string[];
  statuses?: generated.QueueJobStatus[];
  orderBy: OrderBySelection;
};

const initialControls: QueueJobsControls = {
  limit: 20,
  page: 1,
  orderBy: {
    field: "ran_at",
    descending: true,
  }
}

export class QueueJobsController {
  private controlsSubject: BehaviorSubject<QueueJobsControls>;
  controls$: Observable<QueueJobsControls>;

  private variablesSubject: BehaviorSubject<generated.QueueJobsQueryVariables>;
  variables$: Observable<generated.QueueJobsQueryVariables>;

  constructor(ctrl: QueueJobsControls = initialControls) {
    this.controlsSubject = new BehaviorSubject(ctrl);
    this.controls$ = this.controlsSubject.asObservable();
    this.variablesSubject = new BehaviorSubject(
      controlsToQueryVariables(ctrl),
    );
    this.variables$ = this.variablesSubject.asObservable();
    this.controls$.pipe(debounceTime(100)).subscribe((ctrl) => {
      const currentParams = this.variablesSubject.getValue();
      const nextParams = controlsToQueryVariables(ctrl);
      if (JSON.stringify(currentParams) !== JSON.stringify(nextParams)) {
        this.variablesSubject.next(nextParams);
      }
    });
  }
}

const controlsToQueryVariables = (ctrl: QueueJobsControls) : generated.QueueJobsQueryVariables => ({
  input: {
    limit: ctrl.limit,
    page: ctrl.page,
    orderBy: [
      ctrl.orderBy,
      ...ctrl.orderBy.field === "ran_at" ? [{ field: "created_at" as const, descending: ctrl.orderBy.descending }] : [],
    ],
    queues: ctrl.queues,
    statuses: ctrl.statuses,
  }
})
