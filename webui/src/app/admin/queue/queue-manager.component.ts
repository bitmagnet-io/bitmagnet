import { Component, inject } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { AppModule } from "../../app.module";
import { QueueEnqueueReprocessTorrentsBatchDialogComponent } from "./queue-enqueue-reprocess-torrents-batch-dialog.component";
import { QueuePurgeJobsDialogComponent } from "./queue-purge-jobs-dialog.component";

@Component({
  selector: "app-queue-manage",
  template: `
    <ng-container *transloco="let t">
      <app-document-title
        [parts]="[t('routes.admin'), t('routes.queues'), t('routes.admin')]"
      />
      <mat-card>
        <mat-card-content>
          <ul>
            <li>
              <a mat-button (click)="openDialogPurgeJobs()">{{
                t("admin.queues.purge_queue_jobs")
              }}</a>
            </li>
            <li>
              <a
                mat-button
                (click)="openDialogEnqueueReprocessTorrentsBatch()"
                >{{ t("admin.queues.enqueue_torrent_processing_batch") }}</a
              >
            </li>
          </ul>
        </mat-card-content>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      ul {
        list-style-type: none;
        padding-left: 0;
        li {
          margin-bottom: 6px;
        }
      }
    `,
  ],
  standalone: true,
  imports: [AppModule],
})
export class QueueManageComponent {
  readonly dialog = inject(MatDialog);

  openDialogPurgeJobs() {
    this.dialog.open(QueuePurgeJobsDialogComponent);
  }

  openDialogEnqueueReprocessTorrentsBatch() {
    this.dialog.open(QueueEnqueueReprocessTorrentsBatchDialogComponent);
  }
}
