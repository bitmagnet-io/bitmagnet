import { BehaviorSubject, debounceTime, Observable } from "rxjs";
import { TranslocoService } from "@jsverse/transloco";
import * as generated from "../graphql/generated";
import { PageEvent } from "../paginator/paginator.types";

type FacetInput<TValue = unknown> = {
  active: boolean;
  filter?: TValue[];
};

export type ContentTypeSelection = generated.ContentType | "null" | null;

export type TorrentSearchControls = {
  language: string;
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
};

const controlsToQueryVariables = (
  ctrl: TorrentSearchControls,
): generated.TorrentContentSearchQueryVariables => ({
  input: {
    queryString: ctrl.queryString,
    limit: ctrl.limit,
    page: ctrl.page,
    totalCount: true,
    hasNextPage: true,
    orderBy: [ctrl.orderBy],
    facets: {
      contentType: {
        aggregate: true,
        filter: ctrl.contentType
          ? [ctrl.contentType === "null" ? null : ctrl.contentType]
          : undefined,
      },
      genre: ctrl.facets.genre.active
        ? {
            aggregate: true,
            filter: ctrl.facets.genre.filter,
          }
        : undefined,
      language: ctrl.facets.language.active
        ? {
            aggregate: ctrl.facets.language.active,
            filter: ctrl.facets.language.filter,
          }
        : undefined,
      torrentFileType: ctrl.facets.fileType.active
        ? {
            aggregate: true,
            filter: ctrl.facets.fileType.filter,
          }
        : undefined,
      torrentSource: ctrl.facets.torrentSource.active
        ? {
            aggregate: true,
            filter: ctrl.facets.torrentSource.filter,
          }
        : undefined,
      torrentTag: ctrl.facets.torrentTag.active
        ? {
            aggregate: true,
            filter: ctrl.facets.torrentTag.filter,
          }
        : undefined,
      videoResolution: ctrl.facets.videoResolution.active
        ? {
            aggregate: true,
            filter: ctrl.facets.videoResolution.filter,
          }
        : undefined,
      videoSource: ctrl.facets.videoSource.active
        ? {
            aggregate: true,
            filter: ctrl.facets.videoSource.filter,
          }
        : undefined,
    },
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

  constructor(initialControls: TorrentSearchControls) {
    this.controlsSubject = new BehaviorSubject(initialControls);
    this.controls$ = this.controlsSubject.asObservable();
    this.paramsSubject = new BehaviorSubject(
      controlsToQueryVariables(initialControls),
    );
    this.params$ = this.paramsSubject.asObservable();
    this.controls$.pipe(debounceTime(100)).subscribe((ctrl) => {
      const currentParams = this.paramsSubject.getValue();
      const nextParams = controlsToQueryVariables(ctrl);
      if (JSON.stringify(currentParams) !== JSON.stringify(nextParams)) {
        this.paramsSubject.next(nextParams);
      }
    });
  }

  update(fn: (c: TorrentSearchControls) => TorrentSearchControls) {
    const ctrl = this.controlsSubject.getValue();
    const next = fn(ctrl);
    if (JSON.stringify(ctrl) !== JSON.stringify(next)) {
      this.controlsSubject.next(next);
    }
  }

  selectLanguage(lang: string) {
    this.update((ctrl) => ({
      ...ctrl,
      language: lang,
    }));
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
      } else if (orderBy.field === "relevance") {
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

  selectOrderBy(field: generated.TorrentContentOrderByField) {
    const orderBy = {
      field,
      descending:
        orderByOptions.find((option) => option.field === field)?.descending ??
        false,
    };
    this.update((ctrl) => ({
      ...ctrl,
      orderBy:
        orderBy.field !== "relevance" || ctrl.queryString
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
  extractAggregations: (
    aggs: generated.TorrentContentSearchResult["aggregations"],
  ) => Array<Agg<T, _allowNull>>;
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
  extractAggregations: (aggs) => aggs.torrentSource ?? [],
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
  extractAggregations: (aggs) => aggs.torrentTag ?? [],
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
  extractAggregations: (aggs) => aggs.torrentFileType ?? [],
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
  extractAggregations: (aggs) => aggs.language ?? [],
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
  extractAggregations: (aggs) => aggs.genre ?? [],
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
  extractAggregations: (aggs) =>
    (aggs.videoResolution ?? []).map((agg) => ({
      ...agg,
      value: agg.value ?? null,
    })),
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
  extractAggregations: (aggs) =>
    (aggs.videoSource ?? []).map((agg) => ({
      ...agg,
      value: (agg.value as generated.VideoSource | undefined) ?? null,
    })),
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
    field: "relevance",
    descending: true,
  },
  {
    field: "published_at",
    descending: true,
  },
  {
    field: "updated_at",
    descending: true,
  },
  {
    field: "size",
    descending: true,
  },
  {
    field: "files_count",
    descending: true,
  },
  {
    field: "seeders",
    descending: true,
  },
  {
    field: "leechers",
    descending: true,
  },
  {
    field: "name",
    descending: false,
  },
];

export const defaultOrderBy = {
  field: "published_at" as const,
  descending: true,
};

export const defaultQueryOrderBy = {
  field: "relevance" as const,
  descending: true,
};

type OrderBySelection = {
  field: generated.TorrentContentOrderByField;
  descending: boolean;
};

const matchesContentType = (
  selection: ContentTypeSelection,
  cts?: generated.ContentType[],
) => !cts || (selection && cts.includes(selection as generated.ContentType));
