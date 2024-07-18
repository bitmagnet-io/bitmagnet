import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueueJobsTableComponent } from './queue-jobs-table.component';

describe('QueueJobsTableComponent', () => {
  let component: QueueJobsTableComponent;
  let fixture: ComponentFixture<QueueJobsTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [QueueJobsTableComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(QueueJobsTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
