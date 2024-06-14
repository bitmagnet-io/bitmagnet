import { BehaviorSubject, map, Observable } from 'rxjs';
import * as generated from '../graphql/generated';
import { PageEvent } from '../paginator/paginator.types';

type FacetInput<TValue extends string = string> = {
  aggregate: boolean;
  filter?: TValue[];
};

export type ContentTypeSelection = generated.ContentType | 'null' | null;

export type TorrentSearchControls = {
  limit: number;
  page: number;
  queryString?: string;
  contentType: ContentTypeSelection;
  orderBy?: {
    field: generated.TorrentContentOrderBy;
    descending: boolean;
  };
  facets: {
    genre: FacetInput;
    language: FacetInput<generated.Language>;
    torrentFileType: FacetInput<generated.FileType>;
    torrentSource: FacetInput;
    torrentTag: FacetInput;
    videoResolution: FacetInput<generated.VideoResolution>;
    videoSource: FacetInput<generated.VideoSource>;
  };
};

const controlsToQueryVariables = (
  ctrl: TorrentSearchControls,
): generated.TorrentContentSearchQueryVariables => ({
  query: {
    queryString: ctrl.queryString,
    limit: ctrl.limit,
    page: ctrl.page,
    totalCount: true,
    hasNextPage: true,
    facets: {
      contentType: {
        aggregate: true,
        filter: ctrl.contentType
          ? [ctrl.contentType === 'null' ? null : ctrl.contentType]
          : undefined,
      },
      ...ctrl.facets,
    },
  },
});

export class TorrentsSearchController {
  private controlsSubject: BehaviorSubject<TorrentSearchControls>;
  controls$: Observable<TorrentSearchControls>;

  params$: Observable<generated.TorrentContentSearchQueryVariables>;

  constructor(initialControls: TorrentSearchControls) {
    this.controlsSubject = new BehaviorSubject(initialControls);
    this.controls$ = this.controlsSubject.asObservable();
    this.params$ = this.controls$.pipe(map(controlsToQueryVariables));
  }

  update(fn: (c: TorrentSearchControls) => TorrentSearchControls) {
    const ctrl = this.controlsSubject.getValue();
    const next = fn(ctrl);
    if (JSON.stringify(ctrl) !== JSON.stringify(next)) {
      this.controlsSubject.next(next);
    }
  }

  handlePageEvent(event: PageEvent) {
    this.update((ctrl) => ({
      ...ctrl,
      limit: event.pageSize,
      page: event.page,
    }));
  }

  selectContentType(ct: ContentTypeSelection) {
    this.update((ctrl) => ({
      ...ctrl,
      contentType: ct,
    }));
  }
}
