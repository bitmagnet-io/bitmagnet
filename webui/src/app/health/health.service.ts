import { inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { BehaviorSubject, map } from "rxjs";
import * as generated from "../graphql/generated";

export type HealthStatus =
  | generated.HealthStatus
  | "started"
  | "error"
  | "degraded";

const icons: Record<HealthStatus, string> = {
  error: "error",
  degraded: "warning",
  down: "warning",
  unknown: "pending",
  inactive: "circle",
  up: "check_circle",
  started: "play_circle",
};

type Check = generated.HealthCheck & {
  icon: string;
};

type Worker = generated.Worker & {
  icon: string;
};

type Result = {
  status: HealthStatus;
  checks: Check[];
  workers: Worker[];
  error: Error | null;
  icon: string;
};

const initialResult: Result = {
  status: "unknown",
  checks: [],
  icon: icons.unknown,
  workers: [],
  error: null,
};

const pollInterval = 10000;

export class HealthService {
  private apollo = inject(Apollo);

  private resultSubject = new BehaviorSubject<Result>(initialResult);

  result$ = this.resultSubject.asObservable();

  result = initialResult;

  constructor() {
    this.watchQuery();
    this.result$.subscribe((result) => {
      this.result = result;
    });
  }

  private watchQuery() {
    this.apollo
      .watchQuery<
        generated.HealthCheckQuery,
        generated.HealthCheckQueryVariables
      >({
        query: generated.HealthCheckDocument,
        fetchPolicy: "no-cache",
        pollInterval,
      })
      .valueChanges.pipe(
        map(
          (r): Result => ({
            status:
              r.data.health.status === "down"
                ? "degraded"
                : r.data.health.status,
            checks: r.data.health.checks.map((c) => ({
              ...c,
              icon: icons[c.status],
            })),
            workers: r.data.workers.listAll.workers.map((w) => ({
              ...w,
              icon: icons[w.started ? "started" : "inactive"],
            })),
            icon: icons[r.data.health.status],
            error: null,
          }),
        ),
      )
      .subscribe({
        next: (result) => this.resultSubject.next(result),
        error: (error: Error) => {
          this.resultSubject.next({
            status: "error",
            checks: [],
            workers: [],
            error,
            icon: icons.error,
          });
          setTimeout(this.watchQuery.bind(this), pollInterval);
        },
      });
  }
}
