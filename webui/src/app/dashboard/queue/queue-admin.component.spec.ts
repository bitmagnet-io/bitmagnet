import { ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../../app.config";
import { QueueAdminComponent } from "./queue-admin.component";

describe("QueueAdminComponent", () => {
  let component: QueueAdminComponent;
  let fixture: ComponentFixture<QueueAdminComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();

    fixture = TestBed.createComponent(QueueAdminComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
