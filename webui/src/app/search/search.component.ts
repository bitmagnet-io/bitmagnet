import {
  AfterContentInit,
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
} from "@angular/core";
import { catchError, EMPTY, tap } from "rxjs";
import { FormControl } from "@angular/forms";
import {
  animate,
  state,
  style,
  transition,
  trigger,
} from "@angular/animations";
import { COMMA, ENTER } from "@angular/cdk/keycodes";
import { SelectionModel } from "@angular/cdk/collections";
import * as generated from "../graphql/generated";
import { GraphQLService } from "../graphql/graphql.service";
import { AppErrorsService } from "../app-errors.service";
import normalizeTagInput from "../util/normalizeTagInput";
import { SearchEngine } from "./search.engine";

@Component({
  selector: "app-search",
  templateUrl: "./search.component.html",
  styleUrls: ["./search.component.scss"],
  animations: [
    trigger("detailExpand", [
      state("collapsed", style({ height: "0px", minHeight: "0" })),
      state("expanded", style({ height: "*" })),
      transition(
        "expanded <=> collapsed",
        animate("225ms cubic-bezier(0.4, 0.0, 0.2, 1)"),
      ),
    ]),
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SearchComponent implements AfterContentInit, AfterViewInit {
  search: SearchEngine = new SearchEngine(
    this.graphQLService,
    this.errorsService,
  );
  displayedColumns = [
    "select",
    "summary",
    "size",
    "publishedAt",
    "peers",
    "magnet",
  ];
  queryString = new FormControl("");

  items = Array<generated.TorrentContent>();

  contentType = new FormControl<generated.ContentType | "null" | undefined>(
    undefined,
  );

  readonly separatorKeysCodes = [ENTER, COMMA] as const;

  expandedTorrentContentId: string | undefined;

  selectedItems = new SelectionModel<generated.TorrentContent>(true, []);
  selectedTabIndex = 0;

  newTagCtrl = new FormControl<string>("");
  editedTags = Array<string>();
  suggestedTags = Array<string>();

  constructor(
    private graphQLService: GraphQLService,
    private errorsService: AppErrorsService,
  ) {
    this.search.items$.subscribe((items) => {
      this.items = items;
      this.selectedItems.setSelection(
        ...items.filter(({ id }) =>
          this.selectedItems.selected.some(({ id: itemId }) => itemId === id),
        ),
      );
    });
  }

  ngAfterContentInit() {
    this.loadResult();
  }

  ngAfterViewInit() {
    this.contentType.valueChanges.subscribe((value) => {
      this.search.selectContentType(value);
    });
    this.newTagCtrl.valueChanges.subscribe((value) => {
      value &&
        this.newTagCtrl.setValue(normalizeTagInput(value), {
          emitEvent: false,
        });
      this.updateSuggestedTags();
    });
    this.updateSuggestedTags();
  }

  loadResult(cached = true) {
    this.search.loadResult(cached);
  }

  /**
   * Workaround for untyped table cell definitions
   */
  item(item: generated.TorrentContent): generated.TorrentContent {
    return item;
  }

  originalOrder(): number {
    return 0;
  }

  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    return this.items.every((i) => this.selectedItems.isSelected(i));
  }

  /** Selects all rows if they are not all selected; otherwise clear selection. */
  toggleAllRows() {
    if (this.isAllSelected()) {
      this.selectedItems.clear();
      return;
    }
    this.selectedItems.select(...this.items);
  }

  /** The label for the checkbox on the passed row */
  checkboxLabel(row?: generated.TorrentContent): string {
    if (!row) {
      return `${this.isAllSelected() ? "deselect" : "select"} all`;
    }
    return `${this.selectedItems.isSelected(row) ? "deselect" : "select"} ${
      row.torrent.name
    }`;
  }

  selectTab(index: number): void {
    this.selectedTabIndex = index;
  }

  addTag(tagName: string) {
    if (!this.editedTags.includes(tagName)) {
      this.editedTags.push(tagName);
    }
    this.newTagCtrl.reset();
    this.updateSuggestedTags();
  }

  deleteTag(tagName: string) {
    this.editedTags = this.editedTags.filter((t) => t !== tagName);
    this.updateSuggestedTags();
  }

  renameTag(fromTagName: string, toTagName: string) {
    this.editedTags = this.editedTags.map((t) =>
      t === fromTagName ? toTagName : t,
    );
    this.updateSuggestedTags();
  }

  putTags() {
    const infoHashes = this.selectedItems.selected.map((i) => i.infoHash);
    if (!infoHashes.length) {
      return;
    }
    this.newTagCtrl.value && this.addTag(this.newTagCtrl.value);
    return this.graphQLService
      .torrentPutTags({
        infoHashes,
        tagNames: this.editedTags,
      })
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(`Error putting tags: ${err.message}`);
          return EMPTY;
        }),
      )
      .pipe(
        tap(() => {
          this.search.loadResult(false);
        }),
      )
      .subscribe();
  }

  setTags() {
    const infoHashes = this.selectedItems.selected.map((i) => i.infoHash);
    if (!infoHashes.length) {
      return;
    }
    this.newTagCtrl.value && this.addTag(this.newTagCtrl.value);
    return this.graphQLService
      .torrentSetTags({
        infoHashes,
        tagNames: this.editedTags,
      })
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(`Error setting tags: ${err.message}`);
          return EMPTY;
        }),
      )
      .pipe(
        tap(() => {
          this.search.loadResult(false);
        }),
      )
      .subscribe();
  }

  deleteTags() {
    const infoHashes = this.selectedItems.selected.map((i) => i.infoHash);
    if (!infoHashes.length) {
      return;
    }
    this.newTagCtrl.value && this.addTag(this.newTagCtrl.value);
    return this.graphQLService
      .torrentDeleteTags({
        infoHashes,
        tagNames: this.editedTags,
      })
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(`Error deleting tags: ${err.message}`);
          return EMPTY;
        }),
      )
      .pipe(
        tap(() => {
          this.search.loadResult(false);
        }),
      )
      .subscribe();
  }

  private updateSuggestedTags() {
    return this.graphQLService
      .torrentSuggestTags({
        query: {
          prefix: this.newTagCtrl.value,
          exclusions: this.editedTags,
        },
      })
      .pipe(
        tap((result) => {
          this.suggestedTags.splice(
            0,
            this.suggestedTags.length,
            ...result.suggestions.map((t) => t.name),
          );
        }),
      )
      .subscribe();
  }

  selectedInfoHashes(): string[] {
    return this.selectedItems.selected.map((i) => i.infoHash);
  }

  deleteTorrents(infoHashes: string[]) {
    this.graphQLService
      .torrentDelete({ infoHashes })
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(
            `Error deleting torrents: ${err.message}`,
          );
          return EMPTY;
        }),
      )
      .pipe(
        tap(() => {
          this.search.loadResult(false);
        }),
      )
      .subscribe();
  }

  toggleTorrentContentId(id: string) {
    if (this.expandedTorrentContentId === id) {
      this.expandedTorrentContentId = undefined;
    } else {
      this.expandedTorrentContentId = id;
    }
  }
}
