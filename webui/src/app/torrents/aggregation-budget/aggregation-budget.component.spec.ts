import { ComponentFixture, TestBed } from "@angular/core/testing";

import { AggregationBudgetComponent } from "./aggregation-budget.component";

describe("AggregationBudgetDialogComponent", () => {
  let component: AggregationBudgetComponent;
  let fixture: ComponentFixture<AggregationBudgetComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AggregationBudgetComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(AggregationBudgetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
