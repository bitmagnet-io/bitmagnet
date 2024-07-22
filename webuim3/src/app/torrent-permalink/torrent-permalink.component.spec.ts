import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TorrentPermalinkComponent } from './torrent-permalink.component';

describe('TorrentPermalinkComponent', () => {
  let component: TorrentPermalinkComponent;
  let fixture: ComponentFixture<TorrentPermalinkComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TorrentPermalinkComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TorrentPermalinkComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
