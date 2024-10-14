import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnInit,
  Output,
} from "@angular/core";
import { FormControl } from "@angular/forms";
import { catchError, EMPTY, Observable, tap } from "rxjs";
import { COMMA, ENTER } from "@angular/cdk/keycodes";
import * as generated from "../graphql/generated";
import { BreakpointsService } from "../layout/breakpoints.service";
import { ErrorsService } from "../errors/errors.service";
import { GraphQLService } from "../graphql/graphql.service";
import { AppModule } from "../app.module";

@Component({
  selector: "app-torrents-bulk-actions",
  standalone: true,
  imports: [AppModule],
  templateUrl: "./torrents-bulk-actions.component.html",
  styleUrl: "./torrents-bulk-actions.component.scss",
})
export class TorrentsBulkActionsComponent implements OnInit {
  private graphQLService = inject(GraphQLService);
  private errorsService = inject(ErrorsService);
  breakpoints = inject(BreakpointsService);

  @Input() selectedItems$: Observable<generated.TorrentContent[]> =
    new Observable();
  @Output() updated = new EventEmitter<null>();

  readonly separatorKeysCodes = [ENTER, COMMA] as const;
  selectedTabIndex = 0;
  newTagCtrl = new FormControl<string>("");
  editedTags = Array<string>();
  suggestedTags = Array<string>();
  selectedItems = new Array<generated.TorrentContent>();

  ngOnInit() {
    this.selectedItems$.subscribe((items) => {
      this.selectedItems = items;
    });
  }

  selectTab(index: number): void {
    this.selectedTabIndex = index;
  }

  getSelectedMagnetLinks(): string {
    return this.selectedItems.map((i) => i.torrent.magnetUri).join("\n");
  }

  getSelectedInfoHashes(): string {
    return this.selectedItems.map((i) => i.infoHash).join("\n");
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
    const infoHashes = this.selectedItems.map(({ infoHash }) => infoHash);
    if (!infoHashes.length) {
      return;
    }
    if (this.newTagCtrl.value) {
      this.addTag(this.newTagCtrl.value);
    }
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
          this.updated.emit();
        }),
      )
      .subscribe();
  }

  setTags() {
    const infoHashes = this.selectedItems.map(({ infoHash }) => infoHash);
    if (!infoHashes.length) {
      return;
    }
    if (this.newTagCtrl.value) {
      this.addTag(this.newTagCtrl.value);
    }
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
          this.updated.emit();
        }),
      )
      .subscribe();
  }

  deleteTags() {
    const infoHashes = this.selectedItems.map(({ infoHash }) => infoHash);
    if (!infoHashes.length) {
      return;
    }
    if (this.newTagCtrl.value) {
      this.addTag(this.newTagCtrl.value);
    }
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
          this.updated.emit();
        }),
      )
      .subscribe();
  }

  private updateSuggestedTags() {
    return this.graphQLService
      .torrentSuggestTags({
        input: {
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

  deleteTorrents() {
    const infoHashes = this.selectedItems.map(({ infoHash }) => infoHash);
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
          this.updated.emit();
        }),
      )
      .subscribe();
  }
}
