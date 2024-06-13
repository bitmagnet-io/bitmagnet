import {BehaviorSubject, Observable} from "rxjs";
import * as generated from "../graphql/generated";
import {PageEvent} from "../paginator/paginator.types";

export class TorrentsSearchController {

  private paramsSubject: BehaviorSubject<generated.TorrentContentSearchQueryVariables>
  params$: Observable<generated.TorrentContentSearchQueryVariables>
  params: generated.TorrentContentSearchQueryVariables

  private contentTypeSubject: BehaviorSubject<ContentTypeSelection>
  private contentType$ : Observable<ContentTypeSelection>

  constructor(
    initialParams: generated.TorrentContentSearchQueryVariables,
    initialContentType: ContentTypeSelection = null
  ) {
    this.params = initialParams;
    this.paramsSubject = new BehaviorSubject(initialParams);
    this.params$ = this.paramsSubject.asObservable();
    this.paramsSubject.subscribe(params => {
      this.params = params;
    })
    this.contentTypeSubject = new BehaviorSubject<ContentTypeSelection>(initialContentType)
    this.contentType$ = this.contentTypeSubject.asObservable();
    this.contentType$.subscribe((ct) => this.updateParams((params) => ({
      ...params,
      facets: {
        ...params.facets,
        contentType: {
          aggregate: true,
          ...(ct ? { filter: [ct === "null" ? null : ct]} : undefined)
        }
      }
    })))
  }

  updateParams(fn: (p: generated.TorrentContentSearchQueryVariables) => generated.TorrentContentSearchQueryVariables) {
    const params = this.paramsSubject.getValue();
    const nextParams = fn(params);
    if (JSON.stringify(params) !== JSON.stringify(nextParams)) {
      this.paramsSubject.next(nextParams);
    }
  }

  handlePageEvent(event: PageEvent) {
    this.updateParams((params) => ({
      ...params,
      query: {
        ...params.query,
        limit: event.pageSize,
        page: event.page,
      }
    }))
  }

  selectContentType(ct: ContentTypeSelection) {
    this.contentTypeSubject.next(ct);
  }
}

type ContentTypeSelection = generated.ContentType | "null" | null
