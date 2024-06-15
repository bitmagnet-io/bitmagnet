import {
  ChangeDetectionStrategy,
  Component,
  inject, Input,
  OnInit,
  ViewChild,
} from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, MatSort } from '@angular/material/sort';
import { AsyncPipe, DecimalPipe, KeyValuePipe } from '@angular/common';
import { MatIcon } from '@angular/material/icon';
import { MatIconButton } from '@angular/material/button';
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
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { Observable } from 'rxjs';
import { map, shareReplay } from 'rxjs/operators';
import {
  MatExpansionPanel,
  MatExpansionPanelHeader,
  MatExpansionPanelTitle,
} from '@angular/material/expansion';
import { MatRadioButton } from '@angular/material/radio';
import { TranslocoDirective } from '@jsverse/transloco';
import { GraphQLService } from '../graphql/graphql.service';
import { ErrorsService } from '../errors/errors.service';
import { TorrentsTableComponent } from '../torrents-table/torrents-table.component';
import { GraphQLModule } from '../graphql/graphql.module';
import { PaginatorComponent } from '../paginator/paginator.component';
import { contentTypes } from '../taxonomy/content-types';
import {
  ContentTypeSelection,
  TorrentSearchControls,
  TorrentsSearchController,
} from './torrents-search.controller';
import { TorrentsSearchDatasource } from './torrents-search.datasource';
import {BreakpointsService} from "../layout/breakpoints.service";

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
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TorrentsSearchComponent implements OnInit {
  breakpoints = inject(BreakpointsService)

  @ViewChild(PaginatorComponent) paginator: PaginatorComponent;
  @ViewChild(MatSort) sort!: MatSort;
  dataSource: TorrentsSearchDatasource;

  controller: TorrentsSearchController;

  contentTypes = Object.entries(contentTypes).map(([key, info]) => ({
    key: key as keyof typeof contentTypes,
    ...info,
  }));

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    graphQLService: GraphQLService,
    errorsService: ErrorsService,
  ) {
    this.controller = new TorrentsSearchController(initControls);
    this.dataSource = new TorrentsSearchDatasource(
      graphQLService,
      errorsService,
      this.controller.params$,
    );
  }

  ngOnInit(): void {
    this.route.queryParams.subscribe((qParams) =>
      this.controller.update((ctrl) => ({
        ...ctrl,
        contentType: contentTypeParam(qParams, 'contentType'),
        queryString: stringParam(qParams, 'queryString') ?? ctrl.queryString,
        limit: intParam(qParams, 'limit') ?? ctrl.limit,
        page: intParam(qParams, 'page') ?? ctrl.page,
      })),
    );
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
          page,
          limit,
          contentType: ctrl.contentType,
        },
      });
    });
  }
}

const defaultLimit = 20;

const emptyAgg = {
  aggregate: false,
};

const initControls: TorrentSearchControls = {
  page: 1,
  limit: defaultLimit,
  contentType: null,
  facets: {
    genre: emptyAgg,
    language: emptyAgg,
    torrentFileType: emptyAgg,
    torrentSource: emptyAgg,
    torrentTag: emptyAgg,
    videoResolution: emptyAgg,
    videoSource: emptyAgg,
  },
};

const contentTypeParam = (
  params: Params,
  key: string,
): ContentTypeSelection => {
  const str = stringParam(params, 'contentType');
  return str && str in contentTypes ? (str as ContentTypeSelection) : null;
};

const stringParam = (params: Params, key: string): string | undefined => {
  return typeof params[key] === 'string' ? params[key] || undefined : undefined;
};

const intParam = (params: Params, key: string): number | undefined => {
  if (params && params[key] && /^\d+$/.test(params[key])) {
    return parseInt(params[key]);
  }
  return undefined;
};
