import { BehaviorSubject, debounceTime, Observable } from "rxjs";
import { TranslocoService } from "@jsverse/transloco";
import * as generated from "../graphql/generated";
import { PageEvent } from "../paginator/paginator.types";

export type FacetInput<TValue = unknown> = {
  active: boolean;
  filter?: TValue[];
};

export type ContentTypeSelection = generated.ContentType | "null" | null;

export const torrentTabNames = [
  "files",
  "tags",
  "reprocess",
  "delete",
] as const;

export type TorrentTab = (typeof torrentTabNames)[number];

export type TorrentTabSelection = TorrentTab | undefined;

export type TorrentSelection = {
  infoHash: string;
  tab: TorrentTabSelection;
};

const compareTorrentSelection = (
  a?: TorrentSelection,
  b?: TorrentSelection,
): boolean => {
  if (a && b) {
    return a.infoHash === b.infoHash && a.tab === b.tab;
  }
  return a === b;
};

export type TorrentSearchControls = {
  index?: string;
  limit: number;
  page: number;
  queryString?: string;
  contentType: ContentTypeSelection;
  orderBy: OrderBySelection;
  facets: {
    genre: FacetInput<string>;
    language: FacetInput<generated.Language>;
    fileType: FacetInput<generated.FileType>;
    torrentSource: FacetInput<string>;
    torrentTag: FacetInput<string>;
    videoResolution: FacetInput<generated.VideoResolution>;
    videoSource: FacetInput<generated.VideoSource>;
  };
  selectedTorrent?: TorrentSelection;
};

const controlsToQueryVariables = (
  ctrl: TorrentSearchControls,
): generated.TorrentContentSearchQueryVariables => ({
  input: {
    index: ctrl.index,
    queryString: ctrl.queryString,
    limit: ctrl.limit,
    page: ctrl.page,
    totalCount: true,
    hasNextPage: true,
    orderBy: [
      {
        // todo: check field vs key
        field: ctrl.orderBy.key,
        descending: ctrl.orderBy.descending,
      },
    ],
    facets: [
      {
        key: "content_type",
        aggregate: true,
        filter:
          ctrl.contentType && ctrl.contentType !== "null"
            ? [ctrl.contentType as string]
            : undefined,
      },
      ...(ctrl.facets.genre.active
        ? Array<generated.FacetInput>({
            key: "content_genre",
            aggregate: true,
            filter: ctrl.facets.genre.filter,
          })
        : []),
      ...(ctrl.facets.language.active
        ? Array<generated.FacetInput>({
            key: "language",
            aggregate: true,
            filter: ctrl.facets.language.filter,
          })
        : []),
      ...(ctrl.facets.fileType.active
        ? Array<generated.FacetInput>({
            key: "file_type",
            aggregate: true,
            filter: ctrl.facets.fileType.filter,
          })
        : []),
      ...(ctrl.facets.torrentSource.active
        ? Array<generated.FacetInput>({
            key: "torrent_source",
            aggregate: true,
            filter: ctrl.facets.torrentSource.filter,
          })
        : []),
      ...(ctrl.facets.torrentTag.active
        ? Array<generated.FacetInput>({
            key: "tag",
            aggregate: true,
            filter: ctrl.facets.torrentTag.filter,
          })
        : []),
      ...(ctrl.facets.videoResolution.active
        ? Array<generated.FacetInput>({
            key: "video_resolution",
            aggregate: true,
            filter: ctrl.facets.videoResolution.filter,
          })
        : []),
      ...(ctrl.facets.videoSource.active
        ? Array<generated.FacetInput>({
            key: "video_source",
            aggregate: true,
            filter: ctrl.facets.videoSource.filter,
          })
        : []),
    ],
  },
});

export const inactiveFacet = {
  active: false,
};

export class TorrentsSearchController {
  private controlsSubject: BehaviorSubject<TorrentSearchControls>;
  controls$: Observable<TorrentSearchControls>;

  private paramsSubject: BehaviorSubject<generated.TorrentContentSearchQueryVariables>;
  params$: Observable<generated.TorrentContentSearchQueryVariables>;

  private selectionSubject: BehaviorSubject<TorrentSelection | undefined>;
  selection$: Observable<TorrentSelection | undefined>;

  constructor(initialControls: TorrentSearchControls) {
    // initialControls = {
    //   ...initialControls,
    //   index:
    //     (initialControls.index ??
    //       this.browserStorage.get(browserStorageIndexKey)) ||
    //     undefined,
    // };
    this.controlsSubject = new BehaviorSubject(initialControls);
    this.controls$ = this.controlsSubject.asObservable();
    this.paramsSubject = new BehaviorSubject(
      controlsToQueryVariables(initialControls),
    );
    this.params$ = this.paramsSubject.asObservable();
    this.selectionSubject = new BehaviorSubject(
      initialControls.selectedTorrent,
    );
    this.selection$ = this.selectionSubject.asObservable();
    this.controls$.pipe(debounceTime(100)).subscribe((ctrl) => {
      const nextParams = controlsToQueryVariables(ctrl);
      if (
        JSON.stringify(this.paramsSubject.getValue()) !==
        JSON.stringify(nextParams)
      ) {
        this.paramsSubject.next(nextParams);
      }
      if (
        !compareTorrentSelection(
          this.selectionSubject.getValue(),
          ctrl.selectedTorrent,
        )
      ) {
        this.selectionSubject.next(ctrl.selectedTorrent);
      }
    });
  }

  get controls(): TorrentSearchControls {
    return this.controlsSubject.getValue();
  }

  update(fn: (c: TorrentSearchControls) => TorrentSearchControls) {
    const ctrl = this.controlsSubject.getValue();
    const next = fn(ctrl);
    if (JSON.stringify(ctrl) !== JSON.stringify(next)) {
      this.controlsSubject.next(next);
    }
  }

  selectIndex(index?: string) {
    this.update((ctrl) => ({
      ...ctrl,
      index,
    }));
  }

  selectTorrent(infoHash: string, tab?: TorrentTabSelection | null) {
    this.update((ctrl) => {
      if (tab === undefined) {
        tab = ctrl.selectedTorrent?.tab;
      } else if (tab === null) {
        tab = undefined;
      }
      return {
        ...ctrl,
        selectedTorrent: {
          infoHash,
          tab,
        },
      };
    });
  }

  selectContentType(ct: ContentTypeSelection) {
    this.update((ctrl) => ({
      ...ctrl,
      contentType: ct,
      page: 1,
      facets: {
        ...ctrl.facets,
        genre: matchesContentType(ct, genreFacet.contentTypes)
          ? ctrl.facets.genre
          : inactiveFacet,
        videoResolution: matchesContentType(
          ct,
          videoResolutionFacet.contentTypes,
        )
          ? ctrl.facets.videoResolution
          : inactiveFacet,
        videoSource: matchesContentType(ct, videoSourceFacet.contentTypes)
          ? ctrl.facets.videoSource
          : inactiveFacet,
      },
    }));
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  activateFacet(def: FacetDefinition<any, any>) {
    this.update((ctrl) => ({
      ...ctrl,
      facets: def.patchInput(ctrl.facets, {
        ...def.extractInput(ctrl.facets),
        active: true,
      }),
    }));
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  deactivateFacet(def: FacetDefinition<any, any>) {
    this.update((ctrl) => {
      const input = def.extractInput(ctrl.facets);
      return {
        ...ctrl,
        page: input.filter ? 1 : ctrl.page,
        facets: def.patchInput(ctrl.facets, {
          ...input,
          active: false,
          filter: undefined,
        }),
      };
    });
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  activateFilter(def: FacetDefinition<any, any>, filter: string) {
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
  deactivateFilter(def: FacetDefinition<any, any>, filter: string) {
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

  setQueryString(str?: string | null) {
    str = str || undefined;
    this.update((ctrl) => {
      let orderBy = ctrl.orderBy;
      if (str) {
        if (str !== ctrl.queryString) {
          orderBy = defaultQueryOrderBy;
        }
      } else if (orderBy.key === "relevance") {
        orderBy = defaultOrderBy;
      }
      return {
        ...ctrl,
        queryString: str,
        orderBy,
        page: str === ctrl.queryString ? ctrl.page : 1,
      };
    });
  }

  selectOrderBy(key: generated.TorrentContentOrderByField) {
    const orderBy = {
      key,
      descending:
        orderByOptions.find((option) => option.key === key)?.descending ??
        false,
    };
    this.update((ctrl) => ({
      ...ctrl,
      orderBy:
        orderBy.key !== "relevance" || ctrl.queryString
          ? orderBy
          : defaultOrderBy,
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

// todo: Get labels from backend

export type FacetDefinition<T, _allowNull extends boolean = boolean> = {
  key: string;
  icon: string;
  contentTypes?: generated.ContentType[];
  allowNull: _allowNull;
  extractInput: (facets: TorrentSearchControls["facets"]) => FacetInput<T>;
  patchInput: (
    facets: TorrentSearchControls["facets"],
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    input: FacetInput<any>,
  ) => TorrentSearchControls["facets"];
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  resolveLabel: (agg: Agg<any, any>, t: TranslocoService) => string;
};

export const torrentSourceFacet: FacetDefinition<string, false> = {
  key: "torrent_source",
  icon: "mediation",
  allowNull: false,
  extractInput: (f) => f.torrentSource,
  patchInput: (f, i) => ({
    ...f,
    torrentSource: i,
  }),
  resolveLabel: (agg) => agg.label,
};

export const torrentTagFacet: FacetDefinition<string, false> = {
  key: "torrent_tag",
  icon: "sell",
  allowNull: false,
  extractInput: (f) => f.torrentTag,
  patchInput: (f, i) => ({
    ...f,
    torrentTag: i,
  }),
  resolveLabel: (agg) => agg.value as string,
};

export const fileTypeFacet: FacetDefinition<generated.FileType, false> = {
  key: "file_type",
  icon: "file_present",
  allowNull: false,
  extractInput: (f) => f.fileType,
  patchInput: (f, i) => ({
    ...f,
    fileType: i,
  }),
  resolveLabel: (agg, t) => t.translate(`file_types.${agg.value}`),
};

export const languageFacet: FacetDefinition<generated.Language, false> = {
  key: "language",
  icon: "translate",
  allowNull: false,
  extractInput: (f) => f.language,
  patchInput: (f, i) => ({
    ...f,
    language: i,
  }),
  resolveLabel: (agg, t) => t.translate(`languages.${agg.value}`),
};

export const genreFacet: FacetDefinition<string, false> = {
  key: "genre",
  icon: "theater_comedy",
  allowNull: false,
  contentTypes: ["movie", "tv_show"],
  extractInput: (f) => f.genre,
  patchInput: (f, i) => ({
    ...f,
    genre: i,
  }),
  resolveLabel: (agg) => agg.label,
};

export const videoResolutionFacet: FacetDefinition<
  generated.VideoResolution,
  true
> = {
  key: "video_resolution",
  icon: "aspect_ratio",
  allowNull: true,
  contentTypes: ["movie", "tv_show", "xxx"],
  extractInput: (f) => f.videoResolution,
  patchInput: (f, i) => ({
    ...f,
    videoResolution: i,
  }),
  resolveLabel: (agg) => (agg.value as string | undefined)?.slice(1) ?? "?",
};

export const videoSourceFacet: FacetDefinition<generated.VideoSource, true> = {
  key: "video_source",
  icon: "album",
  allowNull: true,
  contentTypes: ["movie", "tv_show", "xxx"],
  extractInput: (f) => f.videoSource,
  patchInput: (f, i) => ({
    ...f,
    videoSource: i,
  }),
  resolveLabel: (agg) => (agg.value as string | undefined) ?? "?",
};

export const facets = [
  torrentSourceFacet,
  torrentTagFacet,
  fileTypeFacet,
  languageFacet,
  genreFacet,
  videoResolutionFacet,
  videoSourceFacet,
];

export type FacetValue<T = unknown, _allowNull extends boolean = true> =
  | T
  | (_allowNull extends true ? null : T);

export type Agg<T, _allowNull extends boolean> = {
  value: FacetValue<T, _allowNull>;
  label: string;
  count: number;
  isEstimate: boolean;
};

export type FacetInfo<
  T = unknown,
  _allowNull extends boolean = true,
> = FacetDefinition<_allowNull> & {
  relevant: boolean;
  active: boolean;
  filter?: Array<FacetValue<T, _allowNull>>;
  aggregations: Array<Agg<T, _allowNull>>;
};

export const orderByOptions: OrderBySelection[] = [
  {
    key: "relevance",
    descending: true,
  },
  {
    key: "published_at",
    descending: true,
  },
  {
    key: "updated_at",
    descending: true,
  },
  {
    key: "size",
    descending: true,
  },
  {
    key: "files_count",
    descending: true,
  },
  {
    key: "seeders",
    descending: true,
  },
  {
    key: "leechers",
    descending: true,
  },
  {
    key: "name",
    descending: false,
  },
];

export const defaultOrderBy = {
  key: "published_at" as const,
  descending: true,
};

export const defaultQueryOrderBy = {
  key: "relevance" as const,
  descending: true,
};

export type OrderBySelection = {
  key: generated.TorrentContentOrderByField;
  descending: boolean;
};

const matchesContentType = (
  selection: ContentTypeSelection,
  cts?: generated.ContentType[],
) => !cts || (selection && cts.includes(selection as generated.ContentType));

export const isDefaultOrdering = (ctrl: TorrentSearchControls): boolean => {
  if (!ctrl.orderBy.descending) {
    return false;
  }
  return ctrl.orderBy.key === (ctrl.queryString ? "relevance" : "published_at");
};
