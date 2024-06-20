import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TorrentsTableComponent } from './torrents-table.component';

describe('TorrentsTableComponent', () => {
  let component: TorrentsTableComponent;
  let fixture: ComponentFixture<TorrentsTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TorrentsTableComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TorrentsTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
