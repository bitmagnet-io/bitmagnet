import { Component, Inject, inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import * as generated from "../../graphql/generated";
import { AppModule } from "../../app.module";
import { availableQueueNames, statusNames } from "./queue.constants";
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

  protected readonly availableQueueNames = availableQueueNames;
  protected readonly statusNames = statusNames;

  protected stage: "PENDING" | "REQUESTING" | "DONE" = "PENDING";

  @Inject(MAT_DIALOG_DATA) public data: { onEnqueued?: () => void };

  purge = true;
  apisDisabled = true;
  localSearchDisabled = true;
  classifierRematch = false;
  contentTypes?: Array<generated.ContentType | null>;
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
            contentTypes: this.contentTypes,
            orphans: this.orphans ? true : undefined,
          },
        },
      })
      .subscribe(() => {
        this.stage = "DONE";
        this.data.onEnqueued?.();
      });
  }
}
