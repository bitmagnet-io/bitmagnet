import { ComponentFixture, TestBed } from "@angular/core/testing";

import { appConfig } from "../app.config";
import { AppModule } from "../app.module";
import { TorrentFilesTableComponent } from "./torrent-files-table.component";

describe("TorrentFilesTableComponent", () => {
  let component: TorrentFilesTableComponent;
  let fixture: ComponentFixture<TorrentFilesTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      ...appConfig,
      imports: [AppModule],
    }).compileComponents();

    fixture = TestBed.createComponent(TorrentFilesTableComponent);
    component = fixture.componentInstance;
    component.torrent = {
      name: "test",
      infoHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      magnetUri: "magnet:?xt=urn:btih:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      size: 100,
      filesStatus: "single",
      hasFilesInfo: true,
      sources: [],
      tagNames: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
