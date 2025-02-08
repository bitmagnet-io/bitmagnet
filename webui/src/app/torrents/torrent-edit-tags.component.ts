import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnInit,
  Output,
} from "@angular/core";
import { FormControl } from "@angular/forms";
import { TranslocoService } from "@jsverse/transloco";
import { catchError, EMPTY, tap } from "rxjs";
import { COMMA, ENTER } from "@angular/cdk/keycodes";
import * as generated from "../graphql/generated";
import { AppModule } from "../app.module";
import { GraphQLService } from "../graphql/graphql.service";
import { ErrorsService } from "../errors/errors.service";
import normalizeTagInput from "../util/normalizeTagInput";

@Component({
  selector: "app-torrent-edit-tags",
  standalone: true,
  imports: [AppModule],
  templateUrl: "./torrent-edit-tags.component.html",
  styleUrl: "./torrent-edit-tags.component.scss",
})
export class TorrentEditTagsComponent implements OnInit {
  @Input() torrentContent: generated.TorrentContent;

  newTagCtrl = new FormControl<string>("");
  protected editedTags = Array<string>();
  public readonly suggestedTags = Array<string>();

  transloco = inject(TranslocoService);
  grapql = inject(GraphQLService);
  errors = inject(ErrorsService);

  readonly separatorKeysCodes = [ENTER, COMMA] as const;

  @Output() updated = new EventEmitter<null>();

  ngOnInit() {
    this.newTagCtrl.valueChanges.subscribe((value) => {
      if (value) {
        value = normalizeTagInput(value);
        this.newTagCtrl.setValue(value, { emitEvent: false });
      }
      return this.grapql
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
    this.editedTags = this.torrentContent.torrent.tagNames;
    this.newTagCtrl.reset();
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
    this.grapql
      .torrentSetTags({
        infoHashes: [this.torrentContent.infoHash],
        tagNames: this.editedTags,
      })
      .pipe(
        catchError((err: Error) => {
          this.errors.addError(`Error saving tags: ${err.message}`);
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
}
