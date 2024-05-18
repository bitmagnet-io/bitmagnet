import { ComponentFixture, TestBed } from "@angular/core/testing";

import { TorrentContentComponent } from "./torrent-content.component";

describe("TorrentContentComponent", () => {
  let component: TorrentContentComponent;
  let fixture: ComponentFixture<TorrentContentComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TorrentContentComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TorrentContentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
