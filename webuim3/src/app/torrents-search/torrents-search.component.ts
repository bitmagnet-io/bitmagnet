import {
  ChangeDetectionStrategy,
  Component,
  inject,
  OnInit,
  ViewChild
} from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, MatSort } from '@angular/material/sort';
import {AsyncPipe, DecimalPipe, KeyValuePipe} from '@angular/common';
import { MatIcon } from '@angular/material/icon';
import { MatIconButton } from '@angular/material/button';
import { MatListItem, MatNavList } from '@angular/material/list';
import {
  MatDrawer,
  MatDrawerContainer,
  MatDrawerContent,
} from '@angular/material/sidenav';
import { MatToolbar } from '@angular/material/toolbar';
import {ActivatedRoute, Params, Router, RouterLinkActive} from '@angular/router';
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import {Observable} from 'rxjs';
import { map, shareReplay } from 'rxjs/operators';
import {GraphQLService} from "../graphql/graphql.service";
import {ErrorsService} from "../errors/errors.service";
import {TorrentsSearchDatasource} from "./torrents-search.datasource";
import {TorrentsTableComponent} from "../torrents-table/torrents-table.component";
import {GraphQLModule} from "../graphql/graphql.module";
import {PaginatorComponent} from "../paginator/paginator.component";
import {TorrentsSearchController} from "./torrents-search.controller";
import * as generated from "../graphql/generated";
import {contentTypes} from "../taxonomy/content-types";
import {MatExpansionPanel, MatExpansionPanelHeader, MatExpansionPanelTitle} from "@angular/material/expansion";
import {MatRadioButton} from "@angular/material/radio";

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
    MatRadioButton
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TorrentsSearchComponent implements OnInit {
  private breakpointObserver = inject(BreakpointObserver);

  isHandset$: Observable<boolean> = this.breakpointObserver
    .observe(Breakpoints.Handset)
    .pipe(
      map((result) => result.matches),
      shareReplay(),
    );

  @ViewChild(PaginatorComponent) paginator: PaginatorComponent;
  @ViewChild(MatSort) sort!: MatSort;
  dataSource: TorrentsSearchDatasource

  controller: TorrentsSearchController;

  contentTypes = contentTypes

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    graphQLService: GraphQLService,
    errorsService: ErrorsService,
  ) {
    this.controller = new TorrentsSearchController(initSearchParams);
    this.dataSource = new TorrentsSearchDatasource(graphQLService, errorsService, this.controller.params$)
  }

  ngOnInit(): void {
    this.route.queryParams.subscribe(qParams =>
      this.controller.updateParams((cParams) => ({
          ...cParams,
          query: {
            ...cParams.query,
            limit: intParam(qParams, "limit") ?? cParams.query?.limit ?? defaultLimit,
            page: intParam(qParams, "page") ?? 1,
          }
        })
    ))
    this.controller.params$.subscribe((params) => {
      let page = params.query?.page
      let limit = params.query?.limit
      if (page === 1) {
        page = undefined
      }
      if (limit === defaultLimit) {
        limit = undefined
      }
      return this.router.navigate([], {
        relativeTo: this.route,
        queryParams: {
         page, limit,
        }
      })
      }
    )
  }

  originalOrder(): number {
    return 0;
  }
}

const defaultLimit = 20;

const initSearchParams: generated.TorrentContentSearchQueryVariables = {
  query: {
    limit: defaultLimit,
    hasNextPage: true,
    totalCount: true,
  },
  facets: {
    contentType: {
      aggregate: true,
    },
  },
}

const intParam = (params: Params, key: string): number | undefined => {
  if (params && params[key] && /^\d+$/.test(params[key])) {
    return parseInt(params[key]);
  }
  return undefined
}
