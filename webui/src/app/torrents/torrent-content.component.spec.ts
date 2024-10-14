import { ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../app.config";
import { TorrentContentComponent } from "./torrent-content.component";

describe("TorrentContentComponent", () => {
  let component: TorrentContentComponent;
  let fixture: ComponentFixture<TorrentContentComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();

    fixture = TestBed.createComponent(TorrentContentComponent);
    component = fixture.componentInstance;
    const date = new Date().toISOString();
    const infoHash = "aaaaaaaaaaaaaaaaaaaa";
    component.torrentContent = {
      id: "test",
      infoHash,
      title: "Test",
      torrent: {
        name: "Test",
        infoHash,
        size: 10,
        filesStatus: "no_info",
        hasFilesInfo: false,
        magnetUri: `magnet:?xt=urn:btih:${infoHash}`,
        sources: [],
        tagNames: [],
        createdAt: date,
        updatedAt: date,
      },
      publishedAt: date,
      createdAt: date,
      updatedAt: date,
    };
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
