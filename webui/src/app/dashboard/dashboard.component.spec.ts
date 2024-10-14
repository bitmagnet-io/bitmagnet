import { ComponentFixture, TestBed } from "@angular/core/testing";

import { DomSanitizer } from "@angular/platform-browser";
import { MatIconRegistry } from "@angular/material/icon";
import { appConfig } from "../app.config";
import { DashboardComponent } from "./dashboard.component";

describe("DashboardComponent", () => {
  let component: DashboardComponent;
  let fixture: ComponentFixture<DashboardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();

    fixture = TestBed.createComponent(DashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    TestBed.inject(MatIconRegistry).addSvgIcon(
      "queue",
      TestBed.inject(DomSanitizer).bypassSecurityTrustResourceUrl(
        "/fake/icon.svg",
      ),
    );
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
