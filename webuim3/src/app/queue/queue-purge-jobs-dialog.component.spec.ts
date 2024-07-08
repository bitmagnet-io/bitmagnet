import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueuePurgeJobsDialog } from './queue-purge-jobs-dialog.component';

describe('QueuePurgeJobsComponent', () => {
  let component: QueuePurgeJobsDialog;
  let fixture: ComponentFixture<QueuePurgeJobsDialog>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [QueuePurgeJobsDialog]
    })
    .compileComponents();

    fixture = TestBed.createComponent(QueuePurgeJobsDialog);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
