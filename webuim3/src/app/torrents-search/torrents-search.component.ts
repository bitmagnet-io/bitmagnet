import {
  ChangeDetectionStrategy,
  Component,
  inject,
  OnInit,
  ViewChild,
} from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule } from '@angular/material/sort';
import { AsyncPipe, DecimalPipe, KeyValuePipe } from '@angular/common';
import { MatIcon } from '@angular/material/icon';
import { MatIconButton, MatMiniFabButton } from '@angular/material/button';
import { MatListItem, MatNavList } from '@angular/material/list';
import {
  MatDrawer,
  MatDrawerContainer,
  MatDrawerContent,
} from '@angular/material/sidenav';
import { MatToolbar } from '@angular/material/toolbar';
import {
  ActivatedRoute,
  Params,
  Router,
  RouterLinkActive,
} from '@angular/router';
import { BehaviorSubject, combineLatestWith, Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import {
  MatExpansionPanel,
  MatExpansionPanelHeader,
  MatExpansionPanelTitle,
} from '@angular/material/expansion';
import { MatRadioButton } from '@angular/material/radio';
import { TranslocoDirective, TranslocoService } from '@jsverse/transloco';
import { MatCheckbox } from '@angular/material/checkbox';
import { MatFormField, MatLabel } from '@angular/material/form-field';
import { MatInput } from '@angular/material/input';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatTooltip } from '@angular/material/tooltip';
import { MatOption, MatSelect } from '@angular/material/select';
import { SelectionModel } from '@angular/cdk/collections';
import { MatDivider } from '@angular/material/divider';
import { GraphQLService } from '../graphql/graphql.service';
import { ErrorsService } from '../errors/errors.service';
import {
  allColumns,
  compactColumns,
  TorrentsTableComponent,
} from '../torrents-table/torrents-table.component';
import { GraphQLModule } from '../graphql/graphql.module';
import { PaginatorComponent } from '../paginator/paginator.component';
import { contentTypeList, contentTypeMap } from '../taxonomy/content-types';
import { BreakpointsService } from '../layout/breakpoints.service';
import * as generated from '../graphql/generated';
import { TorrentsBulkActionsComponent } from '../torrents-bulk-actions/torrents-bulk-actions.component';
import { intParam, stringListParam, stringParam } from '../util/query-string';
import {
  emptyResult,
  TorrentsSearchDatasource,
} from './torrents-search.datasource';
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
} from './torrents-search.controller';

@Component({
  selector: 'app-torrents-search',
  templateUrl: './torrents-search.component.html',
  styleUrl: './torrents-search.component.scss',
  standalone: true,
  imports: [
    MatTableModule,
    MatSortModule,
    AsyncPipe,
    MatIcon,
    MatIconButton,
    MatLabel,
    MatListItem,
    MatNavList,
    MatToolbar,
    RouterLinkActive,
    MatDrawerContainer,
    MatDrawer,
    MatDrawerContent,
    TorrentsTableComponent,
    GraphQLModule,
    PaginatorComponent,
    MatExpansionPanel,
    MatExpansionPanelTitle,
    MatExpansionPanelHeader,
    DecimalPipe,
    KeyValuePipe,
    MatRadioButton,
    TranslocoDirective,
    MatCheckbox,
    MatFormField,
    MatInput,
    ReactiveFormsModule,
    MatTooltip,
    MatSelect,
    MatOption,
    MatMiniFabButton,
    MatDivider,
    TorrentsBulkActionsComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TorrentsSearchComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private graphQLService = inject(GraphQLService);
  private errorsService = inject(ErrorsService);
  private transloco = inject(TranslocoService);
  breakpoints = inject(BreakpointsService);

  @ViewChild(PaginatorComponent) paginator: PaginatorComponent;
  dataSource: TorrentsSearchDatasource;

  controller: TorrentsSearchController;

  controls: TorrentSearchControls;

  contentTypes = contentTypeList;
  orderByOptions = orderByOptions;

  facets$: Observable<FacetInfo<any, any>[]>;

  allColumns = allColumns;
  compactColumns = compactColumns;

  queryString = new FormControl('');

  result = emptyResult;

  selection = new SelectionModel<string>(true, []);
  private selectedItemsSubject = new BehaviorSubject<
    generated.TorrentContent[]
  >([]);
  selectedItems$ = this.selectedItemsSubject.asObservable();

  constructor() {
    this.controls = {
      ...initControls,
      language: this.transloco.getActiveLang(),
    };
    this.controller = new TorrentsSearchController(this.controls);
    this.dataSource = new TorrentsSearchDatasource(
      this.graphQLService,
      this.errorsService,
      this.controller.params$,
    );
    this.controller.controls$.subscribe((ctrl) => {
      this.controls = ctrl;
    });
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
              controls.contentType !== 'null' &&
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
    this.dataSource.result$.subscribe((result) => {
      this.result = result;
      const infoHashes = new Set(result.items.map(({ infoHash }) => infoHash));
      this.selection.deselect(
        ...this.selection.selected.filter(
          (infoHash) => !infoHashes.has(infoHash),
        ),
      );
    });
    // a bit of a hack to force an update on language switch:
    this.transloco.events$.subscribe(() =>
      this.controller.selectLanguage(this.transloco.getActiveLang()),
    );
  }

  ngOnInit(): void {
    this.route.queryParams.subscribe((params) => {
      const queryString = stringParam(params, 'query');
      this.queryString.setValue(queryString ?? null);
      this.controller.update((ctrl) => {
        const activeFacets = stringListParam(params, 'facets');
        let orderBy = ctrl.orderBy;
        if (queryString) {
          if (queryString !== ctrl.queryString) {
            orderBy = defaultQueryOrderBy;
          }
        } else if (orderBy.field === 'relevance') {
          orderBy = defaultOrderBy;
        }
        return {
          ...ctrl,
          queryString,
          orderBy,
          contentType: contentTypeParam(params, 'content_type'),
          limit: intParam(params, 'limit') ?? ctrl.limit,
          page: intParam(params, 'page') ?? ctrl.page,
          facets: facets.reduce<TorrentSearchControls['facets']>(
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
    });
    this.controller.controls$.subscribe((ctrl) => {
      let page: number | undefined = ctrl.page;
      let limit: number | undefined = ctrl.limit;
      if (page === 1) {
        page = undefined;
      }
      if (limit === defaultLimit) {
        limit = undefined;
      }
      return this.router.navigate([], {
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
        queryParamsHandling: 'merge',
      });
    });
    this.selection.changed.subscribe((selection) => {
      const infoHashes = new Set(selection.source.selected);
      this.selectedItemsSubject.next(
        this.result.items.filter((i) => infoHashes.has(i.infoHash)),
      );
    });
  }
}

const defaultLimit = 20;

const initControls: TorrentSearchControls = {
  language: 'en',
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
  ctrl: TorrentSearchControls['facets'],
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
    facets: activeFacets.length ? activeFacets.join(',') : undefined,
    ...Object.fromEntries(
      Object.entries(filters).map(([k, values]) => [
        k,
        encodeURIComponent(values.join(',')),
      ]),
    ),
  };
};
