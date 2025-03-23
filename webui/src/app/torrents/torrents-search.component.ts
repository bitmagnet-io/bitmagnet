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
  minSizeControl = new FormControl<number | null>(null);
  maxSizeControl = new FormControl<number | null>(null);
  minSizeUnitControl = new FormControl<string>("MiB");
  maxSizeUnitControl = new FormControl<string>("MiB");

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

  // Helper function to convert size to bytes based on unit
  private sizeToBytes(size: number | null, unit: string): number | undefined {
    if (size === null) {
      return undefined;
    }

    let bytes: number;
    switch (unit) {
      // Standard SI units (KB, MB, GB, TB) - using 1000-based units
      case "KB":
        bytes = Math.floor(size * 1000);
        break;
      case "MB":
        bytes = Math.floor(size * 1000 * 1000);
        break;
      case "GB":
        bytes = Math.floor(size * 1000 * 1000) * 1000;
        break;
      case "TB":
        bytes = Math.floor(size * 1000 * 1000) * 1000 * 1000;
        break;

      // Binary units (KiB, MiB, GiB, TiB) - using 1024-based units
      case "KiB":
        bytes = Math.floor(size * 1024);
        break;
      case "MiB":
        bytes = Math.floor(size * 1024 * 1024);
        break;
      case "GiB":
        bytes = Math.floor(size * 1024 * 1024) * 1024;
        break;
      case "TiB":
        bytes = Math.floor(size * 1024 * 1024) * 1024 * 1024;
        break;
      default:
        bytes = size;
    }

    return bytes;
  }

  updateSizeFilter(): void {
    const minSizeBytes = this.sizeToBytes(
      this.minSizeControl.value,
      this.minSizeUnitControl.value || "MiB",
    );
    const maxSizeBytes = this.sizeToBytes(
      this.maxSizeControl.value,
      this.maxSizeUnitControl.value || "MiB",
    );

    // Update controller
    this.controller.setSizeRange(minSizeBytes, maxSizeBytes);

    // Force a refresh
    this.dataSource.refresh();
  }

  clearSizeFilter(): void {
    // Reset form controls
    this.minSizeControl.setValue(null);
    this.maxSizeControl.setValue(null);
    this.minSizeUnitControl.setValue("MiB");
    this.maxSizeUnitControl.setValue("MiB");

    // Update controller
    this.controller.setSizeRange(undefined, undefined);

    // Force a refresh
    this.dataSource.refresh();
  }

  ngOnInit(): void {
    this.subscriptions.push(
      this.route.queryParams.subscribe((params) => {
        // Update query string
        this.queryString.setValue(stringParam(params, "query") ?? null);

        // Get size values
        const minSize = intParam(params, "min_size");
        const maxSize = intParam(params, "max_size");

        // Get units or defaults
        const minSizeUnit = stringParam(params, "min_size_unit") || "MiB";
        const maxSizeUnit = stringParam(params, "max_size_unit") || "MiB";

        // Set form values
        if (minSize !== undefined) {
          this.minSizeControl.setValue(minSize);
          if (
            ["KB", "MB", "GB", "TB", "KiB", "MiB", "GiB", "TiB"].includes(
              minSizeUnit,
            )
          ) {
            this.minSizeUnitControl.setValue(minSizeUnit);
          }
        } else {
          this.minSizeControl.setValue(null);
        }

        if (maxSize !== undefined) {
          this.maxSizeControl.setValue(maxSize);
          if (
            ["KB", "MB", "GB", "TB", "KiB", "MiB", "GiB", "TiB"].includes(
              maxSizeUnit,
            )
          ) {
            this.maxSizeUnitControl.setValue(maxSizeUnit);
          }
        } else {
          this.maxSizeControl.setValue(null);
        }

        // Update controller with all params
        this.controller.update(() => paramsToControls(params));
      }),
      this.controller.controls$.subscribe((ctrl) => {
        void this.router.navigate([], {
          relativeTo: this.route,
          queryParams: controlsToParams(
            ctrl,
            this.minSizeUnitControl.value || "MiB",
            this.maxSizeUnitControl.value || "MiB",
          ),
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

  // Handle size range parameters
  const minSize = intParam(params, "min_size");
  const maxSize = intParam(params, "max_size");
  const minSizeUnit = stringParam(params, "min_size_unit") || "MiB";
  const maxSizeUnit = stringParam(params, "max_size_unit") || "MiB";

  let minSizeBytes, maxSizeBytes;

  // Convert min size to bytes
  if (minSize !== undefined) {
    switch (minSizeUnit) {
      // Standard SI units (KB, MB, GB, TB) - using 1000-based units
      case "KB":
        minSizeBytes = minSize * 1000;
        break;
      case "MB":
        minSizeBytes = minSize * 1000 * 1000;
        break;
      case "GB":
        // For GB values, calculate more carefully to avoid integer overflow
        minSizeBytes = minSize * 1000 * 1000 * 1000;
        break;
      case "TB":
        // For TB values, calculate even more carefully
        minSizeBytes = minSize * 1000 * 1000 * 1000 * 1000;
        break;

      // Binary units (KiB, MiB, GiB, TiB) - using 1024-based units
      case "KiB":
        minSizeBytes = minSize * 1024;
        break;
      case "MiB":
        minSizeBytes = minSize * 1024 * 1024;
        break;
      case "GiB":
        // For GiB values, calculate more carefully to avoid integer overflow
        minSizeBytes = minSize * 1024 * 1024 * 1024;
        break;
      case "TiB":
        // For TiB values, calculate even more carefully
        minSizeBytes = minSize * 1024 * 1024 * 1024 * 1024;
        break;
      default:
        minSizeBytes = minSize * 1024 * 1024; // Default to MiB
    }
  }

  // Convert max size to bytes
  if (maxSize !== undefined) {
    switch (maxSizeUnit) {
      // Standard SI units (KB, MB, GB, TB) - using 1000-based units
      case "KB":
        maxSizeBytes = maxSize * 1000;
        break;
      case "MB":
        maxSizeBytes = maxSize * 1000 * 1000;
        break;
      case "GB":
        // For GB values, calculate more carefully to avoid integer overflow
        maxSizeBytes = maxSize * 1000 * 1000 * 1000;
        break;
      case "TB":
        // For TB values, calculate even more carefully
        maxSizeBytes = maxSize * 1000 * 1000 * 1000 * 1000;
        break;

      // Binary units (KiB, MiB, GiB, TiB) - using 1024-based units
      case "KiB":
        maxSizeBytes = maxSize * 1024;
        break;
      case "MiB":
        maxSizeBytes = maxSize * 1024 * 1024;
        break;
      case "GiB":
        // For GiB values, calculate more carefully to avoid integer overflow
        maxSizeBytes = maxSize * 1024 * 1024 * 1024;
        break;
      case "TiB":
        // For TiB values, calculate even more carefully
        maxSizeBytes = maxSize * 1024 * 1024 * 1024 * 1024;
        break;
      default:
        maxSizeBytes = maxSize * 1024 * 1024; // Default to MiB
    }
  }

  const sizeRange =
    minSize || maxSize
      ? {
          min: minSizeBytes,
          max: maxSizeBytes,
        }
      : undefined;

  return {
    queryString,
    orderBy: orderByParam(params, !!queryString),
    contentType: contentTypeParam(params),
    limit: intParam(params, "limit") ?? defaultLimit,
    page: intParam(params, "page") ?? 1,
    selectedTorrent,
    sizeRange,
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

const controlsToParams = (
  ctrl: TorrentSearchControls,
  minSizeUnit = "MiB",
  maxSizeUnit = "MiB",
): Params => {
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

  // Handle size range params
  let minSize: number | undefined;
  let maxSize: number | undefined;

  if (ctrl.sizeRange) {
    // Convert bytes back to the selected unit, handling large numbers carefully
    if (ctrl.sizeRange.min) {
      // Convert min bytes to selected unit
      switch (minSizeUnit) {
        // Standard SI units (KB, MB, GB, TB) - using 1000-based units
        case "KB":
          minSize = Math.round(ctrl.sizeRange.min / 1000);
          break;
        case "MB":
          minSize = Math.round(ctrl.sizeRange.min / (1000 * 1000));
          break;
        case "GB":
          // More careful division for large numbers
          minSize = Math.round(ctrl.sizeRange.min / 1000 / (1000 * 1000));
          break;
        case "TB":
          // Even more careful division
          minSize = Math.round(
            ctrl.sizeRange.min / (1000 * 1000) / (1000 * 1000),
          );
          break;

        // Binary units (KiB, MiB, GiB, TiB) - using 1024-based units
        case "KiB":
          minSize = Math.round(ctrl.sizeRange.min / 1024);
          break;
        case "MiB":
          minSize = Math.round(ctrl.sizeRange.min / (1024 * 1024));
          break;
        case "GiB":
          // More careful division for large numbers
          minSize = Math.round(ctrl.sizeRange.min / 1024 / (1024 * 1024));
          break;
        case "TiB":
          // Even more careful division
          minSize = Math.round(
            ctrl.sizeRange.min / (1024 * 1024) / (1024 * 1024),
          );
          break;
        default:
          minSize = ctrl.sizeRange.min;
      }
    }

    if (ctrl.sizeRange.max) {
      // Convert max bytes to selected unit
      switch (maxSizeUnit) {
        // Standard SI units (KB, MB, GB, TB) - using 1000-based units
        case "KB":
          maxSize = Math.round(ctrl.sizeRange.max / 1000);
          break;
        case "MB":
          maxSize = Math.round(ctrl.sizeRange.max / (1000 * 1000));
          break;
        case "GB":
          // More careful division for large numbers
          maxSize = Math.round(ctrl.sizeRange.max / 1000 / (1000 * 1000));
          break;
        case "TB":
          // Even more careful division
          maxSize = Math.round(
            ctrl.sizeRange.max / (1000 * 1000) / (1000 * 1000),
          );
          break;

        // Binary units (KiB, MiB, GiB, TiB) - using 1024-based units
        case "KiB":
          maxSize = Math.round(ctrl.sizeRange.max / 1024);
          break;
        case "MiB":
          maxSize = Math.round(ctrl.sizeRange.max / (1024 * 1024));
          break;
        case "GiB":
          // More careful division for large numbers
          maxSize = Math.round(ctrl.sizeRange.max / 1024 / (1024 * 1024));
          break;
        case "TiB":
          // Even more careful division
          maxSize = Math.round(
            ctrl.sizeRange.max / (1024 * 1024) / (1024 * 1024),
          );
          break;
        default:
          maxSize = ctrl.sizeRange.max;
      }
    }
  }

  // Only include size unit params if we have size values
  const sizeParams =
    minSize || maxSize
      ? {
          min_size: minSize,
          max_size: maxSize,
          min_size_unit: minSize ? minSizeUnit : undefined,
          max_size_unit: maxSize ? maxSizeUnit : undefined,
        }
      : {};

  return {
    query: ctrl.queryString ? encodeURIComponent(ctrl.queryString) : undefined,
    page,
    limit,
    content_type: ctrl.contentType,
    order: orderBy?.field,
    desc,
    ...sizeParams,
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
