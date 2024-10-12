import { ComponentFixture, TestBed } from "@angular/core/testing";

import { MatIconRegistry } from "@angular/material/icon";
import { DomSanitizer } from "@angular/platform-browser";
import { appConfig } from "../../app.config";
import { QueueDashboardComponent } from "./queue-dashboard.component";

describe("QueueCardComponent", () => {
  let component: QueueDashboardComponent;
  let fixture: ComponentFixture<QueueDashboardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();

    fixture = TestBed.createComponent(QueueDashboardComponent);
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
