import { Component, inject } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { AppModule } from '../app.module';
import { QueueEnqueueReprocessTorrentsBatchDialog } from './queue-enqueue-reprocess-torrents-batch-dialog.component';
import { QueuePurgeJobsDialog } from './queue-purge-jobs-dialog.component';

@Component({
  selector: 'app-queue-admin',
  standalone: true,
  imports: [AppModule],
  templateUrl: './queue-admin.component.html',
  styleUrl: './queue-admin.component.scss',
})
export class QueueAdminComponent {
  readonly dialog = inject(MatDialog);

  openDialogPurgeJobs() {
    this.dialog.open(QueuePurgeJobsDialog);
  }

  openDialogEnqueueReprocessTorrentsBatch() {
    this.dialog.open(QueueEnqueueReprocessTorrentsBatchDialog);
  }
}
