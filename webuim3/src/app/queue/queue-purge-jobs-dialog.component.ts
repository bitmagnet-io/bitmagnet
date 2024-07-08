import {Component, EventEmitter, Inject, inject, Output} from '@angular/core';
import {GraphQLModule} from "../graphql/graphql.module";
import {Apollo} from "apollo-angular";
import {TranslocoDirective} from "@jsverse/transloco";
import {MatCard, MatCardActions, MatCardContent, MatCardHeader, MatCardTitle} from "@angular/material/card";
import {MatRadioButton, MatRadioGroup} from "@angular/material/radio";
import {MatLabel} from "@angular/material/form-field";
import {availableQueueNames, statusNames} from "./queue.constants";
import * as generated from "../graphql/generated"
import {MatCheckbox, MatCheckboxChange} from "@angular/material/checkbox";
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle
} from "@angular/material/dialog";
import {MatButton} from "@angular/material/button";
import {MatProgressSpinner} from "@angular/material/progress-spinner";

@Component({
  selector: 'app-queue-purge-jobs',
  standalone: true,
  imports: [GraphQLModule, TranslocoDirective, MatCardHeader, MatCard, MatCardTitle, MatCardContent, MatRadioGroup, MatLabel, MatRadioButton, MatCheckbox, MatCardActions, MatButton, MatProgressSpinner, MatDialogTitle, MatDialogContent, MatDialogActions],
  templateUrl: './queue-purge-jobs-dialog.component.html',
  styleUrl: './queue-purge-jobs-dialog.component.scss'
})
export class QueuePurgeJobsDialog {
  apollo = inject(Apollo)
  readonly dialogRef = inject(MatDialogRef<QueuePurgeJobsDialog>);

  queues?: string[]
  statuses?: generated.QueueJobStatus[]

  protected readonly availableQueueNames = availableQueueNames
  protected readonly statusNames = statusNames

  protected stage: "PENDING" | "REQUESTING" | "DONE" = "PENDING"

  @Inject(MAT_DIALOG_DATA) public data: {onPurged: () => void}

  handleQueueEvent(event: MatCheckboxChange) {
    if (event.source.value === "_all") {
      this.queues = undefined;
      return
    }
    if (event.checked) {
      let queues = this.queues ?? [];
      if (!queues.includes(event.source.value as generated.QueueJobStatus)) {
        queues = [...queues, event.source.value]
      }
      if (queues.length === this.availableQueueNames.length) {
        event.source.checked = false;
        this.queues = undefined
      } else {
        this.queues = queues
      }
    } else {
      const queues = this.queues?.filter((q) => q !== event.source.value)
      if (!queues?.length) {
        this.queues = undefined
      } else {
        this.queues = queues
      }
    }
  }

  handleStatusEvent(event: MatCheckboxChange) {
    if (event.source.value === "_all") {
      this.statuses = undefined;
      return
    }
    if (event.checked) {
      let statuses = this.statuses ?? [];
      if (!statuses.includes(event.source.value as generated.QueueJobStatus)) {
        statuses = [...statuses, event.source.value as generated.QueueJobStatus]
      }
      if (statuses.length === this.statusNames.length) {
        event.source.checked = false;
        this.statuses = undefined
      } else {
        this.statuses = statuses
      }
    } else {
      const statuses = this.statuses?.filter((s) => s !== event.source.value)
      if (!statuses?.length) {
        this.statuses = undefined
      } else {
        this.statuses = statuses
      }
    }
  }

  handlePurgeJobs() {
    if (this.stage !== "PENDING") {
      return;
    }
    this.stage = "REQUESTING"
    this.apollo.mutate<generated.QueuePurgeJobsMutation, generated.QueuePurgeJobsMutationVariables>({
      mutation:  generated.QueuePurgeJobsDocument,
      variables: {
        input: {
          queues: this.queues,
          statuses: this.statuses,
        }
      }
    }).subscribe(() => {
      this.stage = "DONE"
      this.data.onPurged();
    })
  }
}
