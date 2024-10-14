import { ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../app.config";
import { QueueChartAdapterTotals } from "../dashboard/queue/queue-chart-adapter.totals";
import { Result } from "../dashboard/queue/queue-metrics.types";
import { ChartComponent } from "./chart.component";

describe("QueueComponent", () => {
  let component: ChartComponent<Result>;
  let fixture: ComponentFixture<ChartComponent<Result>>;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();

    fixture = TestBed.createComponent(ChartComponent<Result>);
    component = fixture.componentInstance;
    TestBed.runInInjectionContext(() => {
      component.adapter = new QueueChartAdapterTotals();
    });
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
