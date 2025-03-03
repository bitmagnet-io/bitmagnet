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
  standalone: true,
  imports: [AppModule],
  templateUrl: "./queue-enqueue-reprocess-torrents-batch-dialog.component.html",
  styleUrl: "./queue-enqueue-reprocess-torrents-batch-dialog.component.scss",
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
