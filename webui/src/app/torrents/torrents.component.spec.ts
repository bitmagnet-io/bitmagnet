import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TorrentsComponent } from './torrents.component';

describe('TorrentsComponent', () => {
  let component: TorrentsComponent;
  let fixture: ComponentFixture<TorrentsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TorrentsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TorrentsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
