import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TorrentFilesTableComponent } from './torrent-files-table.component';

describe('TorrentFilesTableComponent', () => {
  let component: TorrentFilesTableComponent;
  let fixture: ComponentFixture<TorrentFilesTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TorrentFilesTableComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TorrentFilesTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
