import { ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../app.config";
import { PaginatorComponent } from "./paginator.component";

describe("PaginatorComponent", () => {
  let component: PaginatorComponent;
  let fixture: ComponentFixture<PaginatorComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule(appConfig);
    fixture = TestBed.createComponent(PaginatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
