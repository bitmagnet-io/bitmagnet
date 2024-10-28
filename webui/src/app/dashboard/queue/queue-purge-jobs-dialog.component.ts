import { Component, Inject, inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { MatCheckboxChange } from "@angular/material/checkbox";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { map } from "rxjs/operators";
import { catchError, EMPTY } from "rxjs";
import * as generated from "../../graphql/generated";
import { GraphQLModule } from "../../graphql/graphql.module";
import { AppModule } from "../../app.module";
import { availableQueueNames, statusNames } from "./queue.constants";

@Component({
  selector: "app-queue-purge-jobs-dialog",
  standalone: true,
  imports: [AppModule, GraphQLModule],
  templateUrl: "./queue-purge-jobs-dialog.component.html",
  styleUrl: "./queue-purge-jobs-dialog.component.scss",
})
export class QueuePurgeJobsDialogComponent {
  apollo = inject(Apollo);
  readonly dialogRef = inject(MatDialogRef<QueuePurgeJobsDialogComponent>);

  queues?: string[];
  statuses?: generated.QueueJobStatus[];

  protected readonly availableQueueNames = availableQueueNames;
  protected readonly statusNames = statusNames;

  protected stage: "PENDING" | "REQUESTING" | "DONE" = "PENDING";

  @Inject(MAT_DIALOG_DATA) public data: { onPurged?: () => void };

  protected error?: Error;

  handleQueueEvent(event: MatCheckboxChange) {
    if (event.source.value === "_all") {
      this.queues = undefined;
      return;
    }
    if (event.checked) {
      let queues = this.queues ?? [];
      if (!queues.includes(event.source.value as generated.QueueJobStatus)) {
        queues = [...queues, event.source.value];
      }
      if (queues.length === this.availableQueueNames.length) {
        event.source.checked = false;
        this.queues = undefined;
      } else {
        this.queues = queues;
      }
    } else {
      const queues = this.queues?.filter((q) => q !== event.source.value);
      if (!queues?.length) {
        this.queues = undefined;
      } else {
        this.queues = queues;
      }
    }
  }

  handleStatusEvent(event: MatCheckboxChange) {
    if (event.source.value === "_all") {
      this.statuses = undefined;
      return;
    }
    if (event.checked) {
      let statuses = this.statuses ?? [];
      if (!statuses.includes(event.source.value as generated.QueueJobStatus)) {
        statuses = [
          ...statuses,
          event.source.value as generated.QueueJobStatus,
        ];
      }
      if (statuses.length === this.statusNames.length) {
        event.source.checked = false;
        this.statuses = undefined;
      } else {
        this.statuses = statuses;
      }
    } else {
      const statuses = this.statuses?.filter((s) => s !== event.source.value);
      if (!statuses?.length) {
        this.statuses = undefined;
      } else {
        this.statuses = statuses;
      }
    }
  }

  handlePurgeJobs() {
    if (this.stage !== "PENDING") {
      return;
    }
    this.stage = "REQUESTING";
    this.apollo
      .mutate<
        generated.QueuePurgeJobsMutation,
        generated.QueuePurgeJobsMutationVariables
      >({
        mutation: generated.QueuePurgeJobsDocument,
        variables: {
          input: {
            queues: this.queues,
            statuses: this.statuses,
          },
        },
      })
      .pipe(
        catchError((err: Error) => {
          this.stage = "DONE";
          this.error = err;
          return EMPTY;
        }),
        map(() => {
          this.stage = "DONE";
          this.data?.onPurged?.();
        }),
      )
      .subscribe();
  }
}
