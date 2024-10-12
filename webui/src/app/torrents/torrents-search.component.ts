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
  defaultQueryOrderBy,
  FacetInfo,
  facets,
  inactiveFacet,
  orderByOptions,
  TorrentSearchControls,
  TorrentsSearchController,
} from "./torrents-search.controller";

@Component({
  selector: "app-torrents-search",
  templateUrl: "./torrents-search.component.html",
  styleUrl: "./torrents-search.component.scss",
  standalone: true,
  imports: [
    AppModule,
    GraphQLModule,
    PaginatorComponent,
    TorrentsBulkActionsComponent,
    TorrentsTableComponent,
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

  controls: TorrentSearchControls;

  contentTypes = contentTypeList;
  orderByOptions = orderByOptions;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  facets$: Observable<FacetInfo<any, any>[]>;

  allColumns = allColumns;
  compactColumns = compactColumns;

  queryString = new FormControl("");

  result = emptyResult;

  selection = new SelectionModel<string>(true, []);
  private selectedItemsSubject = new BehaviorSubject<
    generated.TorrentContent[]
  >([]);
  selectedItems$ = this.selectedItemsSubject.asObservable();

  private subscriptions = Array<Subscription>();

  constructor() {
    this.controls = {
      ...initControls,
      language: this.transloco.getActiveLang(),
    };
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
        this.selection.deselect(
          ...this.selection.selected.filter(
            (infoHash) => !infoHashes.has(infoHash),
          ),
        );
      }),
    );
  }

  ngOnInit(): void {
    this.subscriptions.push(
      this.route.queryParams.subscribe((params) => {
        const queryString = stringParam(params, "query");
        this.queryString.setValue(queryString ?? null);
        this.controller.update((ctrl) => {
          const activeFacets = stringListParam(params, "facets");
          let orderBy = ctrl.orderBy;
          if (queryString) {
            if (queryString !== ctrl.queryString) {
              orderBy = defaultQueryOrderBy;
            }
          } else if (orderBy.field === "relevance") {
            orderBy = defaultOrderBy;
          }
          return {
            ...ctrl,
            queryString,
            orderBy,
            contentType: contentTypeParam(params, "content_type"),
            limit: intParam(params, "limit") ?? ctrl.limit,
            page: intParam(params, "page") ?? ctrl.page,
            facets: facets.reduce<TorrentSearchControls["facets"]>(
              (acc, facet) => {
                const active = activeFacets?.includes(facet.key) ?? false;
                const filter = stringListParam(params, facet.key);
                return facet.patchInput(acc, {
                  active,
                  filter,
                });
              },
              ctrl.facets,
            ),
          };
        });
      }),
      this.controller.controls$.subscribe((ctrl) => {
        let page: number | undefined = ctrl.page;
        let limit: number | undefined = ctrl.limit;
        if (page === 1) {
          page = undefined;
        }
        if (limit === defaultLimit) {
          limit = undefined;
        }
        void this.router.navigate([], {
          relativeTo: this.route,
          queryParams: {
            query: ctrl.queryString
              ? encodeURIComponent(ctrl.queryString)
              : undefined,
            page,
            limit,
            content_type: ctrl.contentType,
            ...flattenFacets(ctrl.facets),
          },
          queryParamsHandling: "merge",
        });
      }),
      this.selection.changed.subscribe((selection) => {
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
  language: "en",
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

const contentTypeParam = (
  params: Params,
  key: string,
): ContentTypeSelection => {
  const str = stringParam(params, key);
  return str && str in contentTypeMap ? (str as ContentTypeSelection) : null;
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
