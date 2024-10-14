import { waitForAsync, ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../app.config";
import { TorrentsSearchComponent } from "./torrents-search.component";

describe("TableComponent", () => {
  let component: TorrentsSearchComponent;
  let fixture: ComponentFixture<TorrentsSearchComponent>;

  beforeEach(waitForAsync(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TorrentsSearchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should compile", () => {
    expect(component).toBeTruthy();
  });
});
