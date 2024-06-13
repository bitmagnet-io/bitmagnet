import { CollectionViewer, DataSource } from '@angular/cdk/collections';
import {BehaviorSubject, catchError, EMPTY, Observable, scan, Subscription} from 'rxjs';
import { map } from 'rxjs/operators';
import * as generated from '../graphql/generated';
import { GraphQLService } from '../graphql/graphql.service';
import { ErrorsService } from '../errors/errors.service';

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

const emptyBudgetedCount = {
  count: 0,
  isEstimate: false,
};

export class TorrentsSearchDatasource
  implements DataSource<generated.TorrentContent>
{
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
      }
    })
  )
  public availableContentTypes$: Observable<string[]> = this.resultSubject.pipe(
    scan((acc: generated.ContentType[], next) => Array.from(new Set([
      ...acc,
      ...(next.aggregations.contentType ?? []).flatMap((agg) => agg.value ? [agg.value] : []),
    ])), [])
  )

  public contentTypeCounts$: Observable<Record<string, BudgetedCount>> = this.resultSubject.pipe(
    map((result) => Object.fromEntries((result.aggregations.contentType ?? []).map((ct) => [ct.value, {
      count: ct.count,
      isEstimate: ct.isEstimate,
    }])))
  )

  constructor(
    private graphQLService: GraphQLService,
    private errorsService: ErrorsService,
    searchQueryVariables: Observable<generated.TorrentContentSearchQueryVariables>,
  ) {
    searchQueryVariables.subscribe(
      (variables: generated.TorrentContentSearchQueryVariables) =>
        this.loadResult(variables),
    );
    this.resultSubject.subscribe((result) => {
      this.result = result;
    })
  }

  connect({}: CollectionViewer): Observable<generated.TorrentContent[]> {
    return this.items$;
  }

  disconnect(): void {
    this.resultSubject.complete();
  }

  private loadResult(
    variables: generated.TorrentContentSearchQueryVariables,
  ): void {
    if (this.currentSubscription) {
      this.currentSubscription.unsubscribe();
      this.currentSubscription = undefined;
    }
    const currentRequest = this.currentRequest.getValue() + 1;
    this.currentRequest.next(currentRequest);
    const result = this.graphQLService.torrentContentSearch(variables).pipe(
      catchError((err: Error) => {
        this.errorsService.addError(
          `Error loading item results: ${err.message}`,
        );
        return EMPTY;
      }),
    );
    this.currentSubscription = result.subscribe((r) => {
      if (currentRequest === this.currentRequest.getValue()) {
        this.resultSubject.next(r)
      }
    });
  }
}
