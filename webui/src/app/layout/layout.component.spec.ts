import { waitForAsync, ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../app.config";
import { LayoutComponent } from "./layout.component";

describe("LayoutComponent", () => {
  let component: LayoutComponent;
  let fixture: ComponentFixture<LayoutComponent>;

  beforeEach(waitForAsync(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should compile", () => {
    expect(component).toBeTruthy();
  });
});
