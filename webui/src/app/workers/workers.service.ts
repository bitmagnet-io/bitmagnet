import { inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { BehaviorSubject, map } from "rxjs";
import * as generated from "../graphql/generated";
import { WorkerState } from "../graphql/generated";

const workerStateIcons: Record<WorkerState, string> = {
  error: "error",
  idle: "circle",
  running: "play_circle",
  shutdown: "pending",
  startup: "pending",
};

type Worker = generated.Worker & {
  icon: string;
};

type Result = {
  workers: Worker[];
  workerError: boolean;
  error: Error | null;
};

const initialResult: Result = {
  workers: [],
  workerError: false,
  error: null,
};

const pollInterval = 10000;

export class WorkersService {
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
      .watchQuery<generated.WorkersQuery, generated.WorkersQueryVariables>({
        query: generated.WorkersDocument,
        fetchPolicy: "no-cache",
        pollInterval,
      })
      .valueChanges.pipe(
        map(
          (r): Result => ({
            workers: r.data.worker.listAll.workers.map((w) => ({
              ...w,
              icon: workerStateIcons[w.state],
            })),
            workerError: r.data.worker.listAll.workers.some((w) => w.error),
            error: null,
          }),
        ),
      )
      .subscribe({
        next: (result) => {
          this.resultSubject.next(result);
          setTimeout(this.watchQuery.bind(this), pollInterval);
        },
        error: (error: Error) => {
          this.resultSubject.next({
            workers: [],
            workerError: false,
            error,
          });
          setTimeout(this.watchQuery.bind(this), pollInterval);
        },
      });
  }

  public startWorkers(...refs: string[]) {
    this.apollo
      .mutate<
        generated.WorkersStartMutation,
        generated.WorkersStartMutationVariables
      >({
        mutation: generated.WorkersStartDocument,
        variables: {
          refs,
        },
      })
      .pipe(
        map((result) =>
          this.resultSubject.next({
            workers:
              result.data?.worker.start.workers?.map((w) => ({
                ...w,
                icon: workerStateIcons[w.state],
              })) ?? [],
            workerError:
              result.data?.worker.start.workers?.some((w) => w.error) ?? false,
            error: null,
          }),
        ),
      )
      .subscribe();
  }

  public shutdownWorkers(...refs: string[]) {
    this.apollo
      .mutate<
        generated.WorkersShutdownMutation,
        generated.WorkersShutdownMutationVariables
      >({
        mutation: generated.WorkersShutdownDocument,
        variables: {
          refs,
        },
      })
      .pipe(
        map((result) =>
          this.resultSubject.next({
            workers:
              result.data?.worker.shutdown.workers?.map((w) => ({
                ...w,
                icon: workerStateIcons[w.state],
              })) ?? [],
            workerError:
              result.data?.worker.shutdown.workers?.some((w) => w.error) ??
              false,
            error: null,
          }),
        ),
      )
      .subscribe();
  }

  public restartWorkers(...refs: string[]) {
    this.apollo
      .mutate<
        generated.WorkersRestartMutation,
        generated.WorkersRestartMutationVariables
      >({
        mutation: generated.WorkersRestartDocument,
        variables: {
          refs,
        },
      })
      .pipe(
        map((result) => {
          this.resultSubject.next({
            workers:
              result.data?.worker.restart.workers?.map((w) => ({
                ...w,
                icon: workerStateIcons[w.state],
              })) ?? [],
            workerError:
              result.data?.worker.restart.workers?.some((w) => w.error) ??
              false,
            error: null,
          });
          setTimeout(this.watchQuery.bind(this), 2000);
        }),
      )
      .subscribe();
  }
}
