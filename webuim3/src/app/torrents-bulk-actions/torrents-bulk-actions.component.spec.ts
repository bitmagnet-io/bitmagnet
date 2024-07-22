import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TorrentsBulkActionsComponent } from './torrents-bulk-actions.component';

describe('TorrentsBulkActionsComponent', () => {
  let component: TorrentsBulkActionsComponent;
  let fixture: ComponentFixture<TorrentsBulkActionsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TorrentsBulkActionsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TorrentsBulkActionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
