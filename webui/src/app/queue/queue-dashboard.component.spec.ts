import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueueDashboardComponent } from './queue-dashboard.component';

describe('QueueCardComponent', () => {
  let component: QueueDashboardComponent;
  let fixture: ComponentFixture<QueueDashboardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [QueueDashboardComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(QueueDashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
