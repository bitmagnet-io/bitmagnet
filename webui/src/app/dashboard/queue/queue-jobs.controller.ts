import { BehaviorSubject, debounceTime, Observable } from "rxjs";
import { TranslocoService } from "@jsverse/transloco";
import * as generated from "../../graphql/generated";
import { PageEvent } from "../../paginator/paginator.types";

export class QueueJobsController {
  private controlsSubject: BehaviorSubject<QueueJobsControls>;
  controls$: Observable<QueueJobsControls>;

  private variablesSubject: BehaviorSubject<generated.QueueJobsQueryVariables>;
  variables$: Observable<generated.QueueJobsQueryVariables>;

  constructor(ctrl: QueueJobsControls = initialControls) {
    this.controlsSubject = new BehaviorSubject(ctrl);
    this.controls$ = this.controlsSubject.asObservable();
    this.variablesSubject = new BehaviorSubject(controlsToQueryVariables(ctrl));
    this.variables$ = this.variablesSubject.asObservable();
    this.controls$.pipe(debounceTime(100)).subscribe((ctrl) => {
      const currentParams = this.variablesSubject.getValue();
      const nextParams = controlsToQueryVariables(ctrl);
      if (JSON.stringify(currentParams) !== JSON.stringify(nextParams)) {
        this.variablesSubject.next(nextParams);
      }
    });
  }

  update(fn: (c: QueueJobsControls) => QueueJobsControls) {
    const ctrl = this.controlsSubject.getValue();
    const next = fn(ctrl);
    if (JSON.stringify(ctrl) !== JSON.stringify(next)) {
      this.controlsSubject.next(next);
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  activateFilter(def: FacetDefinition<any>, filter: string) {
    this.update((ctrl) => {
      const input = def.extractInput(ctrl.facets);
      return {
        ...ctrl,
        page: 1,
        facets: def.patchInput(ctrl.facets, {
          ...input,
          filter: Array.from(
            new Set([...((input.filter as unknown[]) ?? []), filter]),
          ).sort(),
        }),
      };
    });
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  deactivateFilter(def: FacetDefinition<any>, filter: string) {
    this.update((ctrl) => {
      const input = def.extractInput(ctrl.facets);
      const nextFilter: string[] | undefined = input.filter?.filter(
        (value) => value !== filter,
      );
      return {
        ...ctrl,
        page: 1,
        facets: def.patchInput(ctrl.facets, {
          ...input,
          filter: nextFilter?.length ? nextFilter : undefined,
        }),
      };
    });
  }

  selectOrderBy(field: generated.QueueJobsOrderByField) {
    const orderBy = {
      field,
      descending:
        orderByOptions.find((option) => option.field === field)?.descending ??
        false,
    };
    this.update((ctrl) => ({
      ...ctrl,
      orderBy,
      page: 1,
    }));
  }

  toggleOrderByDirection() {
    this.update((ctrl) => ({
      ...ctrl,
      orderBy: {
        ...ctrl.orderBy,
        descending: !ctrl.orderBy.descending,
      },
      page: 1,
    }));
  }

  handlePageEvent(event: PageEvent) {
    this.update((ctrl) => ({
      ...ctrl,
      limit: event.pageSize,
      page: event.page,
    }));
  }
}

const controlsToQueryVariables = (
  ctrl: QueueJobsControls,
): generated.QueueJobsQueryVariables => ({
  input: {
    limit: ctrl.limit,
    page: ctrl.page,
    totalCount: true,
    orderBy: [
      ctrl.orderBy,
      ...(ctrl.orderBy.field !== "created_at"
        ? [
            {
              field: "created_at" as const,
              descending: ctrl.orderBy.descending,
            },
          ]
        : []),
    ],
    queues: ctrl.queues,
    statuses: ctrl.statuses,
    facets: {
      queue: {
        aggregate: true,
        filter: ctrl.facets.queue.filter,
      },
      status: {
        aggregate: true,
        filter: ctrl.facets.status.filter,
      },
    },
  },
});

type OrderBySelection = {
  field: generated.QueueJobsOrderByField;
  descending: boolean;
};

export const orderByOptions: OrderBySelection[] = [
  {
    field: "created_at",
    descending: true,
  },
  {
    field: "ran_at",
    descending: true,
  },
  {
    field: "priority",
    descending: false,
  },
];

export type QueueJobsControls = {
  limit: number;
  page: number;
  queues?: string[];
  statuses?: generated.QueueJobStatus[];
  orderBy: OrderBySelection;
  facets: {
    queue: FacetInput<string>;
    status: FacetInput<generated.QueueJobStatus>;
  };
};

const initialControls: QueueJobsControls = {
  limit: 20,
  page: 1,
  orderBy: {
    field: "ran_at",
    descending: true,
  },
  facets: {
    queue: {},
    status: {},
  },
};

type FacetInput<TValue = unknown> = {
  filter?: TValue[];
};

export type Agg<T> = {
  value: T;
  label: string;
  count: number;
};

export type FacetDefinition<T> = {
  key: string;
  extractInput: (facets: QueueJobsControls["facets"]) => FacetInput<T>;
  patchInput: (
    facets: QueueJobsControls["facets"],
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    input: FacetInput<any>,
  ) => QueueJobsControls["facets"];
  extractAggregations: (
    aggs: generated.QueueJobsQueryResult["aggregations"],
  ) => Array<Agg<T>>;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  resolveLabel: (agg: Agg<any>, t: TranslocoService) => string;
};

export const queueFacet: FacetDefinition<string> = {
  key: "queue",
  extractInput: (f) => f.queue,
  patchInput: (f, i) => ({
    ...f,
    queue: i,
  }),
  extractAggregations: (aggs) => aggs.queue ?? [],
  resolveLabel: (agg) => agg.label,
};

export const statusFacet: FacetDefinition<generated.QueueJobStatus> = {
  key: "status",
  extractInput: (f) => f.status,
  patchInput: (f, i) => ({
    ...f,
    status: i,
  }),
  extractAggregations: (aggs) => aggs.status ?? [],
  resolveLabel: (agg, t) => t.translate("dashboard.queues." + agg.label),
};

export const facets = [queueFacet, statusFacet];

export type FacetInfo<T = unknown> = FacetDefinition<T> & {
  filter?: Array<T>;
  aggregations: Array<Agg<T>>;
};
