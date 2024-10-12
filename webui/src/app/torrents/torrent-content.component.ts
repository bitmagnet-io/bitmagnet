import { Component, EventEmitter, inject, Input, Output } from "@angular/core";
import { catchError, EMPTY, tap } from "rxjs";
import { FormControl } from "@angular/forms";
import { COMMA, ENTER } from "@angular/cdk/keycodes";
import { NgOptimizedImage } from "@angular/common";
import { TranslocoService } from "@jsverse/transloco";
import { FilesizePipe } from "../pipes/filesize.pipe";
import * as generated from "../graphql/generated";
import normalizeTagInput from "../util/normalizeTagInput";
import { GraphQLService } from "../graphql/graphql.service";
import { ErrorsService } from "../errors/errors.service";
import { BreakpointsService } from "../layout/breakpoints.service";
import { TimeAgoPipe } from "../pipes/time-ago.pipe";
import { AppModule } from "../app.module";
import { TorrentFilesTableComponent } from "./torrent-files-table.component";

@Component({
  selector: "app-torrent-content",
  templateUrl: "./torrent-content.component.html",
  styleUrl: "./torrent-content.component.scss",
  standalone: true,
  imports: [
    AppModule,
    FilesizePipe,
    NgOptimizedImage,
    TimeAgoPipe,
    TorrentFilesTableComponent,
  ],
})
export class TorrentContentComponent {
  breakpoints = inject(BreakpointsService);

  @Input() torrentContent: generated.TorrentContent;
  @Input() heading = true;
  @Input() size = true;
  @Input() peers = true;
  @Input() published = true;

  @Output() updated = new EventEmitter<null>();

  newTagCtrl = new FormControl<string>("");
  private editedTags = Array<string>();
  public readonly suggestedTags = Array<string>();
  public selectedTabIndex = 0;

  readonly separatorKeysCodes = [ENTER, COMMA] as const;

  transloco = inject(TranslocoService);

  constructor(
    private graphQLService: GraphQLService,
    private errorsService: ErrorsService,
  ) {
    this.newTagCtrl.valueChanges.subscribe((value) => {
      if (value) {
        value = normalizeTagInput(value);
        this.newTagCtrl.setValue(value, { emitEvent: false });
      }
      return graphQLService
        .torrentSuggestTags({
          input: {
            prefix: value,
            exclusions: this.torrentContent.torrent.tagNames,
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
    this.editedTags = fn(this.editedTags);
    this.newTagCtrl.reset();
  }

  saveTags(): void {
    this.graphQLService
      .torrentSetTags({
        infoHashes: [this.torrentContent.infoHash],
        tagNames: this.editedTags,
      })
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(`Error saving tags: ${err.message}`);
          return EMPTY;
        }),
      )
      .pipe(
        tap(() => {
          this.editedTags = [];
          this.updated.emit(null);
        }),
      )
      .subscribe();
  }

  delete() {
    this.graphQLService
      .torrentDelete({ infoHashes: [this.torrentContent.infoHash] })
      .pipe(
        catchError((err: Error) => {
          this.errorsService.addError(`Error deleting torrent: ${err.message}`);
          return EMPTY;
        }),
      )
      .pipe(
        tap(() => {
          this.updated.emit(null);
        }),
      )
      .subscribe();
  }

  getAttribute(key: string, source?: string): string | undefined {
    return this.torrentContent.content?.attributes?.find(
      (a) => a.key === key && (source === undefined || a.source === source),
    )?.value;
  }

  getCollections(type: string): string[] | undefined {
    const collections = this.torrentContent.content?.collections
      ?.filter((a) => a.type === type)
      .map((a) => a.name);
    return collections?.length ? collections.sort() : undefined;
  }

  filesCount(): number | undefined {
    if (this.torrentContent.torrent.filesStatus === "single") {
      return 1;
    }
    return this.torrentContent.torrent.filesCount ?? undefined;
  }
}
