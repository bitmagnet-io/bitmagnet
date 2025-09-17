import { Component, Inject, inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { MatSelectChange } from "@angular/material/select";
import { catchError, EMPTY } from "rxjs";
import * as generated from "../../graphql/generated";
import { AppModule } from "../../app.module";
import { contentTypeList } from "../../torrents/content-types";
import { ErrorsService } from "../../errors/errors.service";
import { QueuePurgeJobsDialogComponent } from "./queue-purge-jobs-dialog.component";

@Component({
  selector: "app-queue-enqueue-reprocess-torrents-batch-dialog",
  template: `
    <ng-container *transloco="let t">
      <mat-card>
        <h2 mat-dialog-title>
          {{ t("admin.queues.enqueue_torrent_processing_batch") }}
        </h2>
        <mat-dialog-content>
          @if (stage === "PENDING") {
            <section>
              <mat-checkbox
                [checked]="purge"
                (change)="purge = $event.checked"
                >{{ t("admin.queues.purge_queue_jobs") }}</mat-checkbox
              ><br />
              <mat-checkbox
                [checked]="!localSearchDisabled"
                (change)="
                  localSearchDisabled = !$event.checked;
                  apisDisabled = !$event.checked ? true : apisDisabled
                "
                >{{
                  t("torrents.reprocess.match_content_by_local_search")
                }}</mat-checkbox
              ><br />
              <mat-checkbox
                [checked]="!apisDisabled"
                (change)="apisDisabled = !$event.checked"
                >{{
                  t("torrents.reprocess.match_content_by_external_api_search")
                }}</mat-checkbox
              ><br />
              <mat-checkbox
                [checked]="classifierRematch"
                (change)="classifierRematch = $event.checked"
                >{{ t("torrents.reprocess.force_rematch") }}</mat-checkbox
              ><br />
              <mat-checkbox
                [checked]="orphans"
                (change)="
                  orphans = $event.checked;
                  contentTypes = $event.checked ? ['all'] : contentTypes
                "
                >{{
                  t("admin.queues.process_orphaned_torrents_only")
                }}</mat-checkbox
              >
              <br />
              <mat-form-field class="select-content-types">
                <mat-label>{{ t("facets.content_type") }}</mat-label>
                <mat-select
                  (selectionChange)="onContentTypeSelectionChange($event)"
                  [value]="contentTypes"
                  multiple
                >
                  <mat-option value="all">{{ t("general.all") }}</mat-option>
                  @for (contentType of allContentTypes; track contentType.key) {
                    <mat-option [value]="contentType.key">
                      {{ t("content_types.plural." + contentType.key) }}
                    </mat-option>
                  }
                </mat-select>
              </mat-form-field>
            </section>
          } @else if (stage === "REQUESTING") {
            <mat-spinner></mat-spinner>
          } @else if (stage === "DONE") {
            <p>{{ t("admin.queues.jobs_enqueued") }}</p>
          }
        </mat-dialog-content>
        <mat-dialog-actions>
          @if (stage === "PENDING") {
            <button
              mat-stroked-button
              color="warning"
              (click)="handleEnqueue()"
            >
              {{ t("admin.queues.enqueue_jobs") }}
            </button>
          }
          <button mat-stroked-button (click)="dialogRef.close()">
            {{ t("general.dismiss") }}
          </button>
        </mat-dialog-actions>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      mat-grid-tile mat-card {
        width: 100%;
        height: 100%;
      }

      .select-content-types {
        margin-top: 10px;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule],
})
export class QueueEnqueueReprocessTorrentsBatchDialogComponent {
  apollo = inject(Apollo);
  readonly dialogRef = inject(MatDialogRef<QueuePurgeJobsDialogComponent>);
  private errorsService = inject(ErrorsService);

  allContentTypes = contentTypeList;

  protected stage: "PENDING" | "REQUESTING" | "DONE" = "PENDING";

  @Inject(MAT_DIALOG_DATA) public data: { onEnqueued?: () => void };

  purge = true;
  apisDisabled = true;
  localSearchDisabled = true;
  classifierRematch = false;
  contentTypes: Array<generated.ContentType | "null" | "all"> = ["all"];
  orphans = false;

  handleEnqueue() {
    if (this.stage !== "PENDING") {
      return;
    }
    this.stage = "REQUESTING";
    this.apollo
      .mutate<
        generated.QueueEnqueueReprocessTorrentsBatchMutation,
        generated.QueueEnqueueReprocessTorrentsBatchMutationVariables
      >({
        mutation: generated.QueueEnqueueReprocessTorrentsBatchDocument,
        variables: {
          input: {
            purge: this.purge,
            apisDisabled: this.apisDisabled,
            localSearchDisabled: this.localSearchDisabled,
            classifierRematch: this.classifierRematch,
            contentTypes: this.contentTypes.includes("all")
              ? undefined
              : this.contentTypes.map((ct) =>
                  ct === "null" ? null : (ct as generated.ContentType),
                ),
            orphans: this.orphans ? true : undefined,
          },
        },
      })
      .pipe(
        catchError((error: Error) => {
          this.errorsService.addError(error.message);
          this.dialogRef.close();
          return EMPTY;
        }),
      )
      .subscribe(() => {
        this.stage = "DONE";
        this.data.onEnqueued?.();
      });
  }

  onContentTypeSelectionChange(change: MatSelectChange) {
    if (
      !Array.isArray(change.value) ||
      !change.value.length ||
      (change.value.includes("all") &&
        (!this.contentTypes.includes("all") || change.value.length === 1))
    ) {
      this.contentTypes = ["all"];
    } else {
      this.orphans = false;
      this.contentTypes = this.allContentTypes
        .map((ct) => ct.key)
        .filter((ct) => (change.value as string[]).includes(ct));
    }
  }
}
