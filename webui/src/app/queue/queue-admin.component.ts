import { Component, inject } from '@angular/core';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatAnchor } from '@angular/material/button';
import { MatCard, MatCardContent } from '@angular/material/card';
import { QueueEnqueueReprocessTorrentsBatchDialog } from './queue-enqueue-reprocess-torrents-batch-dialog.component';
import { QueuePurgeJobsDialog } from './queue-purge-jobs-dialog.component';

@Component({
  selector: 'app-queue-admin',
  standalone: true,
  imports: [MatDialogModule, MatAnchor, MatCardContent, MatCard],
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
