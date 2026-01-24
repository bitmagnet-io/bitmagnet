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
import { TorrentReprocessComponent } from "./torrent-reprocess.component";
import { TargetsService } from "./targets.service";
import { TorrentsSendComponent } from "./torrents-send.component";

@Component({
  selector: "app-torrents-bulk-actions",
  standalone: true,
  imports: [AppModule, TorrentReprocessComponent, TorrentsSendComponent],
  template: `
    <ng-container *transloco="let t">
      <mat-tab-group
        animationDuration="0"
        class="tab-group-bulk-actions"
        [selectedIndex]="selectedTabIndex"
        (focusChange)="selectTab($event.index == 5 ? 0 : $event.index)"
        [mat-stretch-tabs]="false"
      >
        <mat-tab
          [aria-labelledby]="'hidden'"
          class="bulk-tab-placeholder"
        ></mat-tab>

        <mat-tab>
          <ng-template mat-tab-label>
            <mat-icon>content_copy</mat-icon>
            @if (breakpoints.sizeAtLeast("Medium")) {
              <span class="label">{{ t("torrents.copy") }}</span>
            }
          </ng-template>

          <ng-template matTabContent>
            <mat-card>
              <mat-card-actions class="button-row">
                <button
                  mat-stroked-button
                  [disabled]="!selectedItems.length"
                  [matTooltip]="t('general.copy_to_clipboard')"
                  [cdkCopyToClipboard]="getSelectedMagnetLinks()"
                >
                  <mat-icon svgIcon="magnet" />{{ t("torrents.magnet_links") }}
                </button>
                <button
                  mat-stroked-button
                  [disabled]="!selectedItems.length"
                  [matTooltip]="t('general.copy_to_clipboard')"
                  [cdkCopyToClipboard]="getSelectedInfoHashesLines()"
                >
                  <mat-icon>tag</mat-icon>{{ t("torrents.info_hashes") }}
                </button>
              </mat-card-actions>
            </mat-card>
          </ng-template>
        </mat-tab>

        @if (targetsService.targets$ | async; as targets) {
          @if (targets.length > 0) {
            <mat-tab>
              <ng-template mat-tab-label>
                <mat-icon>send</mat-icon>
                @if (breakpoints.sizeAtLeast("Medium")) {
                  <span class="label">{{ t("torrents.send") }}</span>
                }
              </ng-template>

              <ng-template matTabContent>
                <app-torrents-send
                  [index]="index"
                  [infoHashes]="selectedInfoHashes"
                />
              </ng-template>
            </mat-tab>
          }
        }

        <mat-tab>
          <ng-template mat-tab-label>
            <mat-icon>sell</mat-icon>
            @if (breakpoints.sizeAtLeast("Medium")) {
              <span class="label">{{ t("torrents.edit_tags") }}</span>
            }
          </ng-template>

          <ng-template matTabContent>
            <mat-card>
              <mat-form-field class="form-edit-tags" subscriptSizing="dynamic">
                <mat-chip-grid #chipGrid aria-label="Enter tags">
                  @for (tagName of editedTags; let j = $index; track tagName) {
                    <mat-chip-row
                      [editable]="true"
                      (edited)="renameTag(tagName, $event.value)"
                      (removed)="deleteTag(tagName)"
                      [aria-description]="'press enter to edit'"
                    >
                      {{ tagName }}
                      <mat-icon matChipRemove>cancel</mat-icon>
                    </mat-chip-row>
                  }
                </mat-chip-grid>
                <input
                  placeholder="{{ t('torrents.tags.placeholder') }}"
                  [formControl]="newTagCtrl"
                  [matAutocomplete]="auto"
                  [matChipInputFor]="chipGrid"
                  [matChipInputSeparatorKeyCodes]="separatorKeysCodes"
                  (matChipInputTokenEnd)="$event.value && addTag($event.value)"
                  [value]="newTagCtrl.value"
                />
                <mat-autocomplete
                  #auto="matAutocomplete"
                  (optionSelected)="addTag($event.option.viewValue)"
                >
                  @for (tagName of suggestedTags; track tagName) {
                    <mat-option [value]="tagName">{{ tagName }}</mat-option>
                  }
                </mat-autocomplete>
              </mat-form-field>
              <mat-card-actions class="button-row">
                <button
                  mat-stroked-button
                  color="primary"
                  [disabled]="!selectedItems.length"
                  (click)="setTags()"
                  matTooltip="{{ t('torrents.tags.set_tip') }}"
                >
                  {{ t("torrents.tags.set") }}
                </button>
                <button
                  mat-stroked-button
                  color="primary"
                  [disabled]="
                    !selectedItems.length ||
                    (!editedTags.length && !newTagCtrl.value)
                  "
                  (click)="putTags()"
                  matTooltip="{{ t('torrents.tags.put_tip') }}"
                >
                  {{ t("torrents.tags.put") }}
                </button>
                <button
                  mat-stroked-button
                  color="primary"
                  [disabled]="
                    !selectedItems.length ||
                    (!editedTags.length && !newTagCtrl.value)
                  "
                  (click)="deleteTags()"
                  matTooltip="{{ t('torrents.tags.delete_tip') }}"
                >
                  {{ t("torrents.tags.delete") }}
                </button>
              </mat-card-actions>
            </mat-card>
          </ng-template>
        </mat-tab>

        <mat-tab>
          <ng-template mat-tab-label>
            <mat-icon>category</mat-icon>
            @if (breakpoints.sizeAtLeast("Medium")) {
              <span class="label">{{ t("torrents.classification") }}</span>
            }
          </ng-template>

          <ng-template matTabContent>
            <app-torrent-reprocess
              [infoHashes]="selectedInfoHashes"
              (updated)="updated.emit(null)"
            />
          </ng-template>
        </mat-tab>

        <mat-tab>
          <ng-template mat-tab-label>
            <mat-icon>delete_forever</mat-icon>
            @if (breakpoints.sizeAtLeast("Medium")) {
              <span class="label">{{ t("torrents.delete") }}</span>
            }
          </ng-template>

          <ng-template matTabContent>
            <mat-card>
              <mat-card-content>
                <p>
                  <strong>{{ t("torrents.delete_are_you_sure") }}</strong>
                  <br />{{ t("torrents.delete_action_cannot_be_undone") }}.
                </p>
              </mat-card-content>
              <mat-card-actions class="button-row">
                <button
                  mat-stroked-button
                  color="warning"
                  [disabled]="!selectedItems.length"
                  (click)="deleteTorrents()"
                >
                  <mat-icon>delete_forever</mat-icon>{{ t("torrents.delete") }}
                </button>
              </mat-card-actions>
            </mat-card>
          </ng-template>
        </mat-tab>
        @if (selectedTabIndex > 0) {
          <mat-tab>
            <ng-template mat-tab-label>
              <mat-icon style="margin-right: 0">close</mat-icon>
            </ng-template>
          </mat-tab>
        }
      </mat-tab-group>
    </ng-container>
  `,
  styles: [
    `
      mat-tab-group {
        padding-left: 10px;
      }

      .mat-mdc-card {
        margin-bottom: 10px;
      }

      button {
        margin-right: 10px;
      }

      p {
        margin-top: 0;
      }

      ::ng-deep .mdc-tab {
        &[aria-labelledby="hidden"] {
          display: none;
        }
      }
    `,
  ],
})
export class TorrentsBulkActionsComponent implements OnInit {
  private graphQLService = inject(GraphQLService);
  private errorsService = inject(ErrorsService);

  targetsService = inject(TargetsService);

  breakpoints = inject(BreakpointsService);

  @Input() index: string | null = null;
  @Input() selectedItems$: Observable<generated.TorrentContent[]> =
    new Observable();
  @Output() updated = new EventEmitter<null>();

  readonly separatorKeysCodes = [ENTER, COMMA] as const;
  selectedTabIndex = 0;
  newTagCtrl = new FormControl<string>("");
  editedTags = Array<string>();
  suggestedTags = Array<string>();
  selectedItems = new Array<generated.TorrentContent>();
  selectedInfoHashes = new Array<string>();

  ngOnInit() {
    this.selectedItems$.subscribe((items) => {
      this.selectedItems = items;
      this.selectedInfoHashes = items.map((i) => i.infoHash);
    });
    this.newTagCtrl.reset();
  }

  selectTab(index: number): void {
    this.selectedTabIndex = index;
  }

  getSelectedMagnetLinks(): string {
    return this.selectedItems.map((i) => i.torrent.magnetUri).join("\n");
  }

  getSelectedInfoHashesLines(): string {
    return this.selectedInfoHashes.join("\n");
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
