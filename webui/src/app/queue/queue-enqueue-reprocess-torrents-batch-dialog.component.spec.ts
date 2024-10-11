import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueueEnqueueReprocessTorrentsBatchDialogComponent } from './queue-enqueue-reprocess-torrents-batch-dialog.component';
import {appConfig} from "../app.config";
import {AppModule} from "../app.module";
import {DialogConfig, DialogRef} from "@angular/cdk/dialog";
import {OverlayRef} from "@angular/cdk/overlay";
import {MatDialogRef} from "@angular/material/dialog";

describe('QueueEnqueueReprocessTorrentsBatchDialogComponent', () => {
  let component: QueueEnqueueReprocessTorrentsBatchDialogComponent;
  let fixture: ComponentFixture<QueueEnqueueReprocessTorrentsBatchDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      providers: [
        ...appConfig.providers,
        {provide: MatDialogRef, useValue: {}}
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(QueueEnqueueReprocessTorrentsBatchDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
