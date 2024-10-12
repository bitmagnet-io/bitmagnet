import { ComponentFixture, TestBed } from "@angular/core/testing";

import { MatDialogRef } from "@angular/material/dialog";
import { appConfig } from "../../app.config";
import { QueuePurgeJobsDialogComponent } from "./queue-purge-jobs-dialog.component";

describe("QueuePurgeJobsDialogComponent", () => {
  let component: QueuePurgeJobsDialogComponent;
  let fixture: ComponentFixture<QueuePurgeJobsDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      providers: [
        ...appConfig.providers,
        { provide: MatDialogRef, useValue: {} },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(QueuePurgeJobsDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
