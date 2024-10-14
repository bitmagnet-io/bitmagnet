import { ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../app.config";
import { HealthWidgetComponent } from "./health-widget.component";
import { HealthModule } from "./health.module";

describe("HealthComponent", () => {
  let component: HealthWidgetComponent;
  let fixture: ComponentFixture<HealthWidgetComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      ...appConfig,
      imports: [HealthModule],
    }).compileComponents();

    fixture = TestBed.createComponent(HealthWidgetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
