import {
  ChangeDetectionStrategy,
  Component,
  inject,
  OnDestroy,
  OnInit,
} from "@angular/core";
import { ActivatedRoute, Params, Router } from "@angular/router";
import {
  BehaviorSubject,
  combineLatestWith,
  Observable,
  Subscription,
} from "rxjs";
import { map } from "rxjs/operators";
import { TranslocoService } from "@jsverse/transloco";
import { FormControl } from "@angular/forms";
import { SelectionModel } from "@angular/cdk/collections";
import { Apollo } from "apollo-angular";
import { ErrorsService } from "../errors/errors.service";
import { GraphQLModule } from "../graphql/graphql.module";
import { PaginatorComponent } from "../paginator/paginator.component";
import { BreakpointsService } from "../layout/breakpoints.service";
import * as generated from "../graphql/generated";
import { intParam, stringListParam, stringParam } from "../util/query-string";
import { AppModule } from "../app.module";
import { DocumentTitleComponent } from "../layout/document-title.component";
import { IntEstimatePipe } from "../pipes/int-estimate.pipe";
import { TimeFrameSelectorComponent } from "../dates/time-frame-selector.component";
import { TimeFrame } from "../dates/parse-timeframe";
import { TorrentsBulkActionsComponent } from "./torrents-bulk-actions.component";
import { contentTypeList, contentTypeMap } from "./content-types";
import {
  allColumns,
  compactColumns,
  TorrentsTableComponent,
} from "./torrents-table.component";
import {
  emptyResult,
  TorrentsSearchDatasource,
} from "./torrents-search.datasource";
import {
  ContentTypeSelection,
  defaultOrderBy,
  FacetInfo,
  facets,
  inactiveFacet,
  isDefaultOrdering,
  orderByOptions,
  OrderBySelection,
  TorrentSearchControls,
  TorrentSelection,
  TorrentsSearchController,
  TorrentTab,
  torrentTabNames,
  TorrentTabSelection,
} from "./torrents-search.controller";

@Component({
  selector: "app-torrents-search",
  templateUrl: "./torrents-search.component.html",
  styleUrl: "./torrents-search.component.scss",
  standalone: true,
  imports: [
    AppModule,
    DocumentTitleComponent,
    GraphQLModule,
    PaginatorComponent,
    TorrentsBulkActionsComponent,
    TorrentsTableComponent,
    IntEstimatePipe,
    TimeFrameSelectorComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TorrentsSearchComponent implements OnInit, OnDestroy {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private apollo = inject(Apollo);
  private errorsService = inject(ErrorsService);
  private transloco = inject(TranslocoService);
  breakpoints = inject(BreakpointsService);

  dataSource: TorrentsSearchDatasource;

  controller: TorrentsSearchController;

  controls = initControls;

  contentTypes = contentTypeList;
  orderByOptions = orderByOptions;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  facets$: Observable<FacetInfo<any, any>[]>;

  allColumns = allColumns;
  compactColumns = compactColumns;

  queryString = new FormControl("");
  timeFrameExpression = "";

  result = emptyResult;

  multiSelection = new SelectionModel<string>(true, []);

  private selectedItemsSubject = new BehaviorSubject<
    generated.TorrentContent[]
  >([]);
  selectedItems$ = this.selectedItemsSubject.asObservable();

  private subscriptions = Array<Subscription>();

  constructor() {
    this.controller = new TorrentsSearchController(this.controls);
    this.dataSource = new TorrentsSearchDatasource(
      this.apollo,
      this.errorsService,
      this.controller.params$,
    );
    this.subscriptions.push(
      this.controller.controls$.subscribe((ctrl) => {
        this.controls = ctrl;
      }),
    );
    this.facets$ = this.controller.controls$.pipe(
      combineLatestWith(this.dataSource.result$),
      map(([controls, result]) =>
        facets.map((f) => ({
          ...f,
          ...f.extractInput(controls.facets),
          relevant:
            !f.contentTypes ||
            !!(
              controls.contentType &&
              controls.contentType !== "null" &&
              f.contentTypes.includes(controls.contentType)
            ),
          aggregations: f
            .extractAggregations(result.aggregations)
            .map((agg) => ({
              ...agg,
              label: f.resolveLabel(agg, this.transloco),
            })),
        })),
      ),
    );
    this.subscriptions.push(
      this.dataSource.result$.subscribe((result) => {
        this.result = result;
        const infoHashes = new Set(
          result.items.map(({ infoHash }) => infoHash),
        );
        this.multiSelection.deselect(
          ...this.multiSelection.selected.filter(
            (infoHash) => !infoHashes.has(infoHash),
          ),
        );
      }),
    );
  }

  // Handle time frame selection change
  onTimeFrameChanged(timeFrame: TimeFrame): void {
    if (timeFrame.isValid) {
      this.timeFrameExpression = timeFrame.expression;
      this.controller.setPublishedAt(this.timeFrameExpression);
      this.dataSource.refresh();
    }
  }

  ngOnInit(): void {
    this.subscriptions.push(
      this.route.queryParams.subscribe((params) => {
        this.queryString.setValue(stringParam(params, "query") ?? null);

        // Get time frame
        const timeFrame = stringParam(params, "published_at");
        if (timeFrame) {
          this.timeFrameExpression = timeFrame;
        } else {
          this.timeFrameExpression = "";
        }

        // Update controller with all params
        this.controller.update(() => paramsToControls(params));
      }),
      this.controller.controls$.subscribe((ctrl) => {
        void this.router.navigate([], {
          relativeTo: this.route,
          queryParams: controlsToParams(ctrl),
          queryParamsHandling: "replace",
        });
      }),
      this.multiSelection.changed.subscribe((selection) => {
        const infoHashes = new Set(selection.source.selected);
        this.selectedItemsSubject.next(
          this.result.items.filter((i) => infoHashes.has(i.infoHash)),
        );
      }),
    );
  }

  ngOnDestroy() {
    this.subscriptions.forEach((subscription) => subscription.unsubscribe());
    this.subscriptions = new Array<Subscription>();
  }
}

const defaultLimit = 20;

const initControls: TorrentSearchControls = {
  page: 1,
  limit: defaultLimit,
  contentType: null,
  orderBy: defaultOrderBy,
  facets: {
    genre: inactiveFacet,
    language: inactiveFacet,
    fileType: inactiveFacet,
    torrentSource: inactiveFacet,
    torrentTag: inactiveFacet,
    videoResolution: inactiveFacet,
    videoSource: inactiveFacet,
  },
};

const paramsToControls = (params: Params): TorrentSearchControls => {
  const queryString = stringParam(params, "query");
  const activeFacets = stringListParam(params, "facets");
  const publishedAt = stringParam(params, "published_at");

  let selectedTorrent: TorrentSelection | undefined;
  const selectedTorrentParam = stringParam(params, "torrent");
  if (selectedTorrentParam) {
    let torrentTabSelection: TorrentTabSelection;
    const strTab = stringParam(params, "tab");
    if (torrentTabNames.includes(strTab as TorrentTab)) {
      torrentTabSelection = strTab as TorrentTab;
    }
    selectedTorrent = {
      infoHash: selectedTorrentParam,
      tab: torrentTabSelection,
    };
  }
  return {
    queryString,
    orderBy: orderByParam(params, !!queryString),
    contentType: contentTypeParam(params),
    limit: intParam(params, "limit") ?? defaultLimit,
    page: intParam(params, "page") ?? 1,
    selectedTorrent,
    publishedAt,
    facets: facets.reduce<TorrentSearchControls["facets"]>((acc, facet) => {
      const active = activeFacets?.includes(facet.key) ?? false;
      const filter = stringListParam(params, facet.key);
      return facet.patchInput(acc, {
        active,
        filter,
      });
    }, initControls.facets),
  };
};

const controlsToParams = (ctrl: TorrentSearchControls): Params => {
  let page: number | undefined = ctrl.page;
  let limit: number | undefined = ctrl.limit;
  if (page === 1) {
    page = undefined;
  }
  if (limit === defaultLimit) {
    limit = undefined;
  }
  const orderBy = isDefaultOrdering(ctrl) ? undefined : ctrl.orderBy;
  let desc: string | undefined;
  if (orderBy) {
    desc = orderBy.descending ? "1" : "0";
  }
  return {
    query: ctrl.queryString ? encodeURIComponent(ctrl.queryString) : undefined,
    page,
    limit,
    content_type: ctrl.contentType,
    order: orderBy?.field,
    desc,
    published_at: ctrl.publishedAt,
    ...(ctrl.selectedTorrent
      ? {
          torrent: ctrl.selectedTorrent.infoHash,
          tab: ctrl.selectedTorrent.tab ?? undefined,
        }
      : {}),
    ...flattenFacets(ctrl.facets),
  };
};

const contentTypeParam = (params: Params): ContentTypeSelection => {
  const str = stringParam(params, "content_type");
  return str && str in contentTypeMap ? (str as ContentTypeSelection) : null;
};

const orderByParam = (params: Params, hasQuery: boolean): OrderBySelection => {
  let desc: boolean | null = null;
  const strDesc = stringParam(params, "desc");
  if (strDesc === "1") {
    desc = true;
  } else if (strDesc === "0") {
    desc = false;
  }
  const field = stringParam(params, "order");
  for (const opt of orderByOptions) {
    if (opt.field === field) {
      return {
        field,
        descending: desc ?? opt.descending,
      };
    }
  }
  return {
    field: hasQuery ? "relevance" : "published_at",
    descending: desc ?? true,
  };
};

const flattenFacets = (
  ctrl: TorrentSearchControls["facets"],
): Record<string, unknown> => {
  const [activeFacets, filters] = facets.reduce<
    [string[], Record<string, string[]>]
  >(
    (acc, f) => {
      const input = f.extractInput(ctrl);
      if (input.active) {
        return [
          [...acc[0], f.key],
          input.filter
            ? {
                ...acc[1],
                [f.key]: input.filter,
              }
            : acc[1],
        ];
      } else {
        return acc;
      }
    },
    [[], {}],
  );
  return {
    facets: activeFacets.length ? activeFacets.join(",") : undefined,
    ...Object.fromEntries(
      Object.entries(filters).map(([k, values]) => [
        k,
        encodeURIComponent(values.join(",")),
      ]),
    ),
  };
};
