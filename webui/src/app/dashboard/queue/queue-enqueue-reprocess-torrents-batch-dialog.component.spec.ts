import { ComponentFixture, TestBed } from "@angular/core/testing";

import { MatDialogRef } from "@angular/material/dialog";
import { appConfig } from "../../app.config";
import { QueueEnqueueReprocessTorrentsBatchDialogComponent } from "./queue-enqueue-reprocess-torrents-batch-dialog.component";

describe("QueueEnqueueReprocessTorrentsBatchDialogComponent", () => {
  let component: QueueEnqueueReprocessTorrentsBatchDialogComponent;
  let fixture: ComponentFixture<QueueEnqueueReprocessTorrentsBatchDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      providers: [
        ...appConfig.providers,
        { provide: MatDialogRef, useValue: {} },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(
      QueueEnqueueReprocessTorrentsBatchDialogComponent,
    );
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
