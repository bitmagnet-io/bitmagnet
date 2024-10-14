import { CollectionViewer, DataSource } from "@angular/cdk/collections";
import {
  BehaviorSubject,
  catchError,
  EMPTY,
  Observable,
  scan,
  Subscription,
} from "rxjs";
import { map } from "rxjs/operators";
import { Apollo } from "apollo-angular";
import * as generated from "../graphql/generated";
import { ErrorsService } from "../errors/errors.service";

export const emptyResult: generated.TorrentContentSearchResult = {
  items: [],
  totalCount: 0,
  totalCountIsEstimate: false,
  aggregations: {},
};

type BudgetedCount = {
  count: number;
  isEstimate: boolean;
};

export class TorrentsSearchDatasource
  implements DataSource<generated.TorrentContent>
{
  private input: generated.TorrentContentSearchQueryInput;

  private currentRequest = new BehaviorSubject(0);
  private currentSubscription?: Subscription;

  private loadingSubject = new BehaviorSubject(false);
  public loading$ = this.loadingSubject.asObservable();

  public result: generated.TorrentContentSearchResult = emptyResult;
  private resultSubject =
    new BehaviorSubject<generated.TorrentContentSearchResult>(this.result);
  public result$ = this.resultSubject.asObservable();

  public items$ = this.resultSubject.pipe(map((result) => result.items));

  public overallTotalCount$ = this.resultSubject.pipe(
    map((result) => {
      let overallTotalCount = 0;
      let overallIsEstimate = false;
      for (const ct of result.aggregations.contentType ?? []) {
        overallTotalCount += ct.count;
        overallIsEstimate = overallIsEstimate || ct.isEstimate;
      }
      return {
        count: overallTotalCount,
        isEstimate: overallIsEstimate,
      };
    }),
  );
  public availableContentTypes$: Observable<generated.ContentType[]> =
    this.resultSubject.pipe(
      scan(
        (acc: generated.ContentType[], next) =>
          Array.from(
            new Set([
              ...acc,
              ...(next.aggregations.contentType ?? []).flatMap((agg) =>
                agg.value ? [agg.value] : [],
              ),
            ]),
          ),
        [],
      ),
    );

  public contentTypeCounts$: Observable<Record<string, BudgetedCount>> =
    this.resultSubject.pipe(
      map((result) =>
        Object.fromEntries<BudgetedCount>(
          (result.aggregations.contentType ?? []).map((ct) => [
            ct.value as string,
            {
              count: ct.count,
              isEstimate: ct.isEstimate,
            },
          ]),
        ),
      ),
    );

  constructor(
    private apollo: Apollo,
    private errorsService: ErrorsService,
    searchQueryVariables: Observable<generated.TorrentContentSearchQueryVariables>,
  ) {
    searchQueryVariables.subscribe(
      (variables: generated.TorrentContentSearchQueryVariables) => {
        this.input = variables.input;
        this.loadResult({
          input: {
            ...variables.input,
            cached: true,
          },
        });
      },
    );
    this.resultSubject.subscribe((result) => {
      this.result = result;
    });
  }

  connect({}: CollectionViewer): Observable<generated.TorrentContent[]> {
    return this.items$;
  }

  disconnect(): void {
    this.resultSubject.complete();
  }

  refresh() {
    this.loadResult({
      input: {
        ...this.input,
        cached: false,
      },
    });
  }

  private loadResult(
    variables: generated.TorrentContentSearchQueryVariables,
  ): void {
    if (this.currentSubscription) {
      this.currentSubscription.unsubscribe();
      this.currentSubscription = undefined;
    }
    this.loadingSubject.next(true);
    const currentRequest = this.currentRequest.getValue() + 1;
    this.currentRequest.next(currentRequest);
    const result = this.apollo
      .query<
        generated.TorrentContentSearchQuery,
        generated.TorrentContentSearchQueryVariables
      >({
        query: generated.TorrentContentSearchDocument,
        variables,
        fetchPolicy: "no-cache",
      })
      .pipe(map((r) => r.data.torrentContent.search))
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
