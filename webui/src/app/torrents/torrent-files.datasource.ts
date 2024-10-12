import { CollectionViewer, DataSource } from "@angular/cdk/collections";
import { Apollo } from "apollo-angular";
import {
  BehaviorSubject,
  catchError,
  EMPTY,
  Observable,
  Subscription,
} from "rxjs";
import { map } from "rxjs/operators";
import { ErrorsService } from "../errors/errors.service";
import * as generated from "../graphql/generated";

const emptyResult = {
  items: [],
  hasNextPage: false,
  totalCount: 0,
  aggregations: {
    queue: [],
    status: [],
  },
};

export interface ITorrentFilesDatasource
  extends DataSource<generated.TorrentFile> {
  loading$: Observable<boolean>;
  result$: Observable<generated.TorrentFilesQueryResult>;
  result: generated.TorrentFilesQueryResult;
  items$: Observable<generated.TorrentFile[]>;
}

export class TorrentFilesDatasource implements ITorrentFilesDatasource {
  private currentRequest = new BehaviorSubject(0);
  private currentSubscription?: Subscription;

  private loadingSubject = new BehaviorSubject(false);
  public loading$ = this.loadingSubject.asObservable();

  public result: generated.TorrentFilesQueryResult = emptyResult;
  private resultSubject =
    new BehaviorSubject<generated.TorrentFilesQueryResult>(this.result);
  public result$ = this.resultSubject.asObservable();

  public items$ = this.resultSubject.pipe(map((result) => result.items));

  constructor(
    private apollo: Apollo,
    private errorsService: ErrorsService,
    queryVariables: Observable<generated.TorrentFilesQueryVariables>,
  ) {
    queryVariables.subscribe(
      (variables: generated.TorrentFilesQueryVariables) => {
        this.loadResult(variables);
      },
    );
    this.resultSubject.subscribe((result) => {
      this.result = result;
    });
  }

  connect({}: CollectionViewer): Observable<generated.TorrentFile[]> {
    return this.items$;
  }

  disconnect(): void {
    this.resultSubject.complete();
  }

  private loadResult(variables: generated.TorrentFilesQueryVariables): void {
    if (this.currentSubscription) {
      this.currentSubscription.unsubscribe();
      this.currentSubscription = undefined;
    }
    this.loadingSubject.next(true);
    const currentRequest = this.currentRequest.getValue() + 1;
    this.currentRequest.next(currentRequest);
    const result = this.apollo
      .query<generated.TorrentFilesQuery, generated.TorrentFilesQueryVariables>(
        {
          query: generated.TorrentFilesDocument,
          variables,
          fetchPolicy: "no-cache",
        },
      )
      .pipe(map((r) => r.data.torrent.files))
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

export class TorrentFilesSingleDatasource implements ITorrentFilesDatasource {
  private file: generated.TorrentFile;

  loading$ = new BehaviorSubject(false).asObservable();
  result: generated.TorrentFilesQueryResult;
  result$: Observable<generated.TorrentFilesQueryResult>;
  items$: Observable<generated.TorrentFile[]>;

  constructor(private torrent: generated.Torrent) {
    this.file = {
      infoHash: torrent.infoHash,
      index: 0,
      path: torrent.name,
      size: torrent.size,
      fileType: torrent.fileType,
      extension: torrent.extension,
      createdAt: torrent.createdAt,
      updatedAt: torrent.updatedAt,
    };
    this.result = {
      hasNextPage: false,
      items: [this.file],
      totalCount: 1,
    };
    this.result$ = new BehaviorSubject(this.result).asObservable();
    this.items$ = new BehaviorSubject([this.file]).asObservable();
  }

  connect({}: CollectionViewer): Observable<generated.TorrentFile[]> {
    return this.items$;
  }

  disconnect(): void {}
}
