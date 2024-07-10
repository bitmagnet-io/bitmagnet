import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueueEnqueueReprocessTorrentsBatchDialog } from './queue-enqueue-reprocess-torrents-batch-dialog.component';

describe('QueueEnqueueReprocessTorrentsBatchDialogComponent', () => {
  let component: QueueEnqueueReprocessTorrentsBatchDialog;
  let fixture: ComponentFixture<QueueEnqueueReprocessTorrentsBatchDialog>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [QueueEnqueueReprocessTorrentsBatchDialog],
    }).compileComponents();

    fixture = TestBed.createComponent(QueueEnqueueReprocessTorrentsBatchDialog);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
