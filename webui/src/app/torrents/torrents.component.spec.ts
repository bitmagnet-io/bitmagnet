import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TorrentsComponent } from './torrents.component';
import {appConfig} from "../app.config";

describe('TorrentsComponent', () => {
  let component: TorrentsComponent;
  let fixture: ComponentFixture<TorrentsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();

    fixture = TestBed.createComponent(TorrentsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
