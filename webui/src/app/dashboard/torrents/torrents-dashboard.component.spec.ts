import { ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../../app.config";
import { TorrentsDashboardComponent } from "./torrents-dashboard.component";

describe("QueueCardComponent", () => {
  let component: TorrentsDashboardComponent;
  let fixture: ComponentFixture<TorrentsDashboardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();

    fixture = TestBed.createComponent(TorrentsDashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
