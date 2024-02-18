import {
  AfterContentInit,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
} from "@angular/core";
import { BehaviorSubject, catchError, EMPTY, tap } from "rxjs";
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
import * as generated from "../../graphql/generated";
import { GraphQLService } from "../../graphql/graphql.service";
import { AppErrorsService } from "../../app-errors.service";
import { TorrentContentSearchEngine } from "./torrent-content-search.engine";

@Component({
  selector: "app-torrent-content",
  templateUrl: "./torrent-content.component.html",
  styleUrls: ["./torrent-content.component.scss"],
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
})
export class TorrentContentComponent
  implements AfterContentInit, AfterViewInit
{
  search: TorrentContentSearchEngine = new TorrentContentSearchEngine(
    this.graphQLService,
    this.errorsService,
  );
  displayedColumns = ["select", "summary", "size", "peers", "magnet"];
  queryString = new FormControl("");

  items = Array<generated.TorrentContent>();

  contentType = new FormControl<generated.ContentType | "null" | undefined>(
    undefined,
  );

  readonly separatorKeysCodes = [ENTER, COMMA] as const;

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

  getAttribute(
    item: generated.TorrentContent,
    key: string,
    source?: string,
  ): string | undefined {
    return item.content?.attributes?.find(
      (a) => a.key === key && (source === undefined || a.source === source),
    )?.value;
  }

  getCollections(
    item: generated.TorrentContent,
    type: string,
  ): string[] | undefined {
    const collections = item.content?.collections
      ?.filter((a) => a.type === type)
      .map((a) => a.name);
    return collections?.length ? collections.sort() : undefined;
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

  expandedItem = new (class {
    private itemSubject = new BehaviorSubject<
      generated.TorrentContent | undefined
    >(undefined);
    newTagCtrl = new FormControl<string>("");
    private editedTags = Array<string>();
    public readonly suggestedTags = Array<string>();
    public selectedTabIndex = 0;

    constructor(private ds: TorrentContentComponent) {
      ds.search.items$.subscribe((items) => {
        const item = this.itemSubject.getValue();
        if (!item) {
          return;
        }
        const nextItem = items.find((i) => i.id === item.id);
        this.editedTags = nextItem?.torrent.tagNames ?? [];
        this.itemSubject.next(nextItem);
      });
      this.newTagCtrl.valueChanges.subscribe((value) => {
        if (value) {
          value = normalizeTagInput(value);
          this.newTagCtrl.setValue(value, { emitEvent: false });
        }
        return ds.graphQLService
          .torrentSuggestTags({
            query: {
              prefix: value,
              exclusions: this.itemSubject.getValue()?.torrent.tagNames,
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
      });
    }

    get id(): string | undefined {
      return this.itemSubject.getValue()?.id;
    }

    toggle(id?: string): void {
      if (id === this.id) {
        id = undefined;
      }
      const nextItem = this.ds.items.find((i) => i.id === id);
      const current = this.itemSubject.getValue();
      if (current?.id !== id) {
        this.itemSubject.next(nextItem);
        this.editedTags = nextItem?.torrent.tagNames ?? [];
        this.newTagCtrl.reset();
        this.selectedTabIndex = 0;
      }
    }

    selectTab(index: number): void {
      this.selectedTabIndex = index;
    }

    addTag(tagName: string) {
      this.editTags((tags) => [...tags, tagName]);
      this.saveTags();
    }

    renameTag(oldTagName: string, newTagName: string) {
      this.editTags((tags) =>
        tags.map((t) => (t === oldTagName ? newTagName : t)),
      );
      this.saveTags();
    }

    deleteTag(tagName: string) {
      this.editTags((tags) => tags.filter((t) => t !== tagName));
      this.saveTags();
    }

    private editTags(fn: (tagNames: string[]) => string[]) {
      const item = this.itemSubject.getValue();
      if (!item) {
        return;
      }
      this.editedTags = fn(this.editedTags);
      this.newTagCtrl.reset();
    }

    saveTags(): void {
      const expanded = this.itemSubject.getValue();
      if (!expanded) {
        return;
      }
      this.ds.graphQLService
        .torrentSetTags({
          infoHashes: [expanded.infoHash],
          tagNames: this.editedTags,
        })
        .pipe(
          catchError((err: Error) => {
            this.ds.errorsService.addError(`Error saving tags: ${err.message}`);
            return EMPTY;
          }),
        )
        .pipe(
          tap(() => {
            this.editedTags = [];
            this.ds.search.loadResult(false);
          }),
        )
        .subscribe();
    }

    delete() {
      const expanded = this.itemSubject.getValue();
      if (!expanded) {
        return;
      }
      this.ds.deleteTorrents([expanded.infoHash]);
    }
  })(this);
}

const normalizeTagInput = (value: string): string =>
  value
    .toLowerCase()
    .replaceAll(/[^a-z0-9\-]/g, "-")
    .replace(/^-+/, "")
    .replaceAll(/-+/g, "-");
