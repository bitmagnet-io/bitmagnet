import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnInit,
  Output,
} from "@angular/core";
import {
  animate,
  state,
  style,
  transition,
  trigger,
} from "@angular/animations";
import { SelectionModel } from "@angular/cdk/collections";
import { BehaviorSubject } from "rxjs";
import { ActivatedRoute, Router } from "@angular/router";
import { FilesizePipe } from "../pipes/filesize.pipe";
import { TimeAgoPipe } from "../pipes/time-ago.pipe";
import * as generated from "../graphql/generated";
import { BreakpointsService } from "../layout/breakpoints.service";
import { stringParam } from "../util/query-string";
import { AppModule } from "../app.module";
import { TorrentsSearchDatasource } from "./torrents-search.datasource";
import { contentTypeInfo } from "./content-types";
import { TorrentChipsComponent } from "./torrent-chips.component";
import { TorrentContentComponent } from "./torrent-content.component";

@Component({
  selector: "app-torrents-table",
  standalone: true,
  imports: [
    AppModule,
    FilesizePipe,
    TimeAgoPipe,
    TorrentChipsComponent,
    TorrentContentComponent,
  ],
  templateUrl: "./torrents-table.component.html",
  styleUrl: "./torrents-table.component.scss",
  animations: [
    trigger("detailExpand", [
      state("collapsed,void", style({ height: "0px", minHeight: "0" })),
      state("expanded", style({ height: "*" })),
      transition(
        "expanded <=> collapsed",
        animate("225ms cubic-bezier(0.4, 0.0, 0.2, 1)"),
      ),
    ]),
  ],
})
export class TorrentsTableComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  breakpoints = inject(BreakpointsService);

  contentTypeInfo = contentTypeInfo;

  @Input() dataSource: TorrentsSearchDatasource;
  @Input() selection: SelectionModel<string>;
  @Input() displayedColumns: readonly Column[] = allColumns;

  @Output() updated = new EventEmitter<string>();

  items = Array<generated.TorrentContent>();

  expandedId = new BehaviorSubject<string | null>(null);

  ngOnInit() {
    this.dataSource.items$.subscribe((items) => {
      this.items = items;
      // if (items.length) {
      //   const expandedId = this.expandedId.getValue();
      //   if (expandedId && !items.some(({ id }) => id === expandedId)) {
      //     this.expandedId.next(null);
      //   }
      // }
    });
    this.route.queryParams.subscribe((params) => {
      const expandedId = this.expandedId.getValue() ?? undefined;
      const nextExpandedId = stringParam(params, "expanded");
      if (expandedId !== nextExpandedId) {
        this.expandedId.next(nextExpandedId ?? null);
      }
    });
    this.expandedId.subscribe((expandedId) => {
      void this.router.navigate([], {
        relativeTo: this.route,
        queryParams: {
          expanded: expandedId ? encodeURIComponent(expandedId) : undefined,
        },
        queryParamsHandling: "merge",
      });
    });
  }

  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    return this.items.every((i) => this.selection.isSelected(i.infoHash));
  }

  /** Selects all rows if they are not all selected; otherwise clear selection. */
  toggleAllRows() {
    if (this.isAllSelected()) {
      this.selection.clear();
      return;
    }
    this.selection.select(...this.items.map((i) => i.infoHash));
  }

  toggleTorrentContentId(id: string) {
    if (this.expandedId.getValue() === id) {
      this.expandedId.next(null);
    } else {
      this.expandedId.next(id);
    }
  }

  /**
   * Workaround for untyped table cell definitions
   */
  item(item: generated.TorrentContent): generated.TorrentContent {
    return item;
  }
}

export const allColumns = [
  "select",
  "summary",
  "size",
  "publishedAt",
  "peers",
  "magnet",
] as const;

export const compactColumns = ["select", "summary", "size", "magnet"] as const;

export type Column = (typeof allColumns)[number];
