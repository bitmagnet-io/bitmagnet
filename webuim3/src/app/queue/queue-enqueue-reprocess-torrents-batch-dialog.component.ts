import {Component, Inject, inject} from '@angular/core';
import {Apollo} from "apollo-angular";
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle
} from "@angular/material/dialog";
import * as generated from "../graphql/generated";
import {QueuePurgeJobsDialog} from "./queue-purge-jobs-dialog.component";
import { availableQueueNames, statusNames } from './queue.constants';
import {MatButton} from "@angular/material/button";
import {MatCard} from "@angular/material/card";
import {MatCheckbox} from "@angular/material/checkbox";
import {MatProgressSpinner} from "@angular/material/progress-spinner";
import {TranslocoDirective} from "@jsverse/transloco";

@Component({
  selector: 'app-queue-enqueue-reprocess-torrents-batch-dialog',
  standalone: true,
  imports: [
    MatButton,
    MatCard,
    MatCheckbox,
    MatDialogActions,
    MatDialogContent,
    MatDialogTitle,
    MatProgressSpinner,
    TranslocoDirective
  ],
  templateUrl: './queue-enqueue-reprocess-torrents-batch-dialog.component.html',
  styleUrl: './queue-enqueue-reprocess-torrents-batch-dialog.component.scss'
})
export class QueueEnqueueReprocessTorrentsBatchDialog {

  apollo = inject(Apollo)
  readonly dialogRef = inject(MatDialogRef<QueuePurgeJobsDialog>);

  protected readonly availableQueueNames = availableQueueNames
  protected readonly statusNames = statusNames

  protected stage: "PENDING" | "REQUESTING" | "DONE" = "PENDING"

  @Inject(MAT_DIALOG_DATA) public data: {onEnqueued: () => void}

  apisDisabled = true;
  localSearchDisabled = true;
  classifierRematch = false;
  contentTypes?: Array<generated.ContentType | null>
  orphans = false;

  handleEnqueue() {
    if (this.stage !== "PENDING") {
      return;
    }
    this.stage = "REQUESTING"
    this.apollo.mutate<generated.QueueEnqueueReprocessTorrentsBatchMutation, generated.QueueEnqueueReprocessTorrentsBatchMutationVariables>({
      mutation:  generated.QueueEnqueueReprocessTorrentsBatchDocument,
      variables: {
        input: {
          apisDisabled: this.apisDisabled,
          localSearchDisabled: this.localSearchDisabled,
          classifierRematch: this.classifierRematch,
          contentTypes: this.contentTypes,
          orphans: this.orphans ? true : undefined,
        }
      }
    }).subscribe(() => {
      this.stage = "DONE"
      this.data.onEnqueued();
    })
  }
}
