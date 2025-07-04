import { waitForAsync, ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../app.config";
import { DashboardWorkersComponent } from "./dashboard-workers.component";

describe("DashboardComponent", () => {
  let component: DashboardWorkersComponent;
  let fixture: ComponentFixture<DashboardWorkersComponent>;

  beforeEach(waitForAsync(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DashboardWorkersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should compile", () => {
    expect(component).toBeTruthy();
  });
});
