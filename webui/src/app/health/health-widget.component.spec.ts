import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HealthWidgetComponent } from './health-widget.component';

describe('HealthComponent', () => {
  let component: HealthWidgetComponent;
  let fixture: ComponentFixture<HealthWidgetComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HealthWidgetComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(HealthWidgetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
