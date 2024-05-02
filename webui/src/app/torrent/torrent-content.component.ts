import { Component, EventEmitter, Input, Output } from "@angular/core";
import { catchError, EMPTY, tap } from "rxjs";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { DecimalPipe, NgOptimizedImage } from "@angular/common";
import {
  MatAutocomplete,
  MatAutocompleteTrigger,
  MatOption,
} from "@angular/material/autocomplete";
import { MatButton } from "@angular/material/button";
import {
  MatCard,
  MatCardActions,
  MatCardContent,
} from "@angular/material/card";
import {
  MatChipGrid,
  MatChipInput,
  MatChipRemove,
  MatChipRow,
} from "@angular/material/chips";
import { MatDivider } from "@angular/material/divider";
import { MatFormField } from "@angular/material/form-field";
import { MatIcon } from "@angular/material/icon";
import {
  MatTab,
  MatTabContent,
  MatTabGroup,
  MatTabLabel,
} from "@angular/material/tabs";
import { MatTooltip } from "@angular/material/tooltip";
import { NgxFilesizeModule } from "ngx-filesize";
import { COMMA, ENTER } from "@angular/cdk/keycodes";
import { CdkCopyToClipboard } from "@angular/cdk/clipboard";
import { AppErrorsService } from "../app-errors.service";
import { GraphQLService } from "../graphql/graphql.service";
import normalizeTagInput from "../util/normalizeTagInput";
import * as generated from "../graphql/generated";

@Component({
  selector: "app-torrent-content",
  standalone: true,
  imports: [
    DecimalPipe,
    MatAutocomplete,
    MatAutocompleteTrigger,
    MatButton,
    MatCard,
    MatCardActions,
    MatCardContent,
    MatChipGrid,
    MatChipInput,
    MatChipRemove,
    MatChipRow,
    MatDivider,
    MatFormField,
    MatIcon,
    MatOption,
    MatTab,
    MatTabContent,
    MatTabGroup,
    MatTabLabel,
    MatTooltip,
    NgOptimizedImage,
    NgxFilesizeModule,
    ReactiveFormsModule,
    CdkCopyToClipboard,
  ],
  templateUrl: "./torrent-content.component.html",
  styleUrl: "./torrent-content.component.scss",
})
export class TorrentContentComponent {
  @Input() torrentContent: generated.TorrentContent;

  @Output() updated = new EventEmitter<null>();

  newTagCtrl = new FormControl<string>("");
  private editedTags = Array<string>();
  public readonly suggestedTags = Array<string>();
  public selectedTabIndex = 0;

  readonly separatorKeysCodes = [ENTER, COMMA] as const;

  constructor(
    private graphQLService: GraphQLService,
    private errorsService: AppErrorsService,
  ) {
    this.newTagCtrl.valueChanges.subscribe((value) => {
      if (value) {
        value = normalizeTagInput(value);
        this.newTagCtrl.setValue(value, { emitEvent: false });
      }
      return graphQLService
        .torrentSuggestTags({
          query: {
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
}
