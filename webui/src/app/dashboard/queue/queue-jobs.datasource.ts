import { Apollo } from "apollo-angular";
import {
  BehaviorSubject,
  catchError,
  EMPTY,
  Observable,
  Subscription,
} from "rxjs";
import { CollectionViewer, DataSource } from "@angular/cdk/collections";
import { map } from "rxjs/operators";
import * as generated from "../../graphql/generated";
import { ErrorsService } from "../../errors/errors.service";

const emptyResult = {
  items: [],
  hasNextPage: false,
  totalCount: 0,
  aggregations: {
    queue: [],
    status: [],
  },
};

export class QueueJobsDatasource implements DataSource<generated.QueueJob> {
  private currentRequest = new BehaviorSubject(0);
  private currentSubscription?: Subscription;

  private loadingSubject = new BehaviorSubject(false);
  public loading$ = this.loadingSubject.asObservable();

  public result: generated.QueueJobsQueryResult = emptyResult;
  private resultSubject = new BehaviorSubject<generated.QueueJobsQueryResult>(
    this.result,
  );
  public result$ = this.resultSubject.asObservable();

  public items$ = this.resultSubject.pipe(map((result) => result.items));

  private variables?: generated.QueueJobsQueryVariables;

  constructor(
    private apollo: Apollo,
    private errorsService: ErrorsService,
    queryVariables: Observable<generated.QueueJobsQueryVariables>,
  ) {
    queryVariables.subscribe((variables: generated.QueueJobsQueryVariables) => {
      this.variables = variables;
      this.loadResult(variables);
    });
    this.resultSubject.subscribe((result) => {
      this.result = result;
    });
  }

  connect({}: CollectionViewer): Observable<generated.QueueJob[]> {
    return this.items$;
  }

  disconnect(): void {
    this.resultSubject.complete();
  }

  refresh() {
    if (this.variables) {
      this.loadResult(this.variables);
    }
  }

  private loadResult(variables: generated.QueueJobsQueryVariables): void {
    if (this.currentSubscription) {
      this.currentSubscription.unsubscribe();
      this.currentSubscription = undefined;
    }
    this.loadingSubject.next(true);
    const currentRequest = this.currentRequest.getValue() + 1;
    this.currentRequest.next(currentRequest);
    const result = this.apollo
      .query<generated.QueueJobsQuery, generated.QueueJobsQueryVariables>({
        query: generated.QueueJobsDocument,
        variables,
        fetchPolicy: "no-cache",
      })
      .pipe(map((r) => r.data.queue.jobs))
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(
            `Error loading item results: ${err.message}`,
          );
          return EMPTY;
        }),
      );
    this.currentSubscription = result.subscribe((r) => {
      if (currentRequest === this.currentRequest.getValue()) {
        this.loadingSubject.next(false);
        this.resultSubject.next(r);
      }
    });
  }
}
