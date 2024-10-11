import { waitForAsync, ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../app.config";
import { DashboardHomeComponent } from "./dashboard-home.component";

describe("DashboardComponent", () => {
  let component: DashboardHomeComponent;
  let fixture: ComponentFixture<DashboardHomeComponent>;

  beforeEach(waitForAsync(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DashboardHomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should compile", () => {
    expect(component).toBeTruthy();
  });
});
