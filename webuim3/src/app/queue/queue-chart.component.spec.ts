import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueueChartComponent } from './queue-chart.component';

describe('QueueComponent', () => {
  let component: QueueChartComponent;
  let fixture: ComponentFixture<QueueChartComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [QueueChartComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(QueueChartComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
