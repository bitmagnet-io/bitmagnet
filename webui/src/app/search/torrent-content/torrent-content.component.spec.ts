import { ComponentFixture, TestBed } from "@angular/core/testing";
import { SearchModule } from "../search.module";
import { AppModule } from "../../app.module";
import { TorrentContentComponent } from "./torrent-content.component";

describe("TorrentContentComponent", () => {
  let component: TorrentContentComponent;
  let fixture: ComponentFixture<TorrentContentComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AppModule, SearchModule],
      declarations: [TorrentContentComponent],
    });
    fixture = TestBed.createComponent(TorrentContentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
