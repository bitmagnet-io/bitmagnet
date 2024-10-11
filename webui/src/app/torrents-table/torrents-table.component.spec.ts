import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TorrentsTableComponent } from './torrents-table.component';
import {appConfig} from "../app.config";
import {TorrentsSearchDatasource} from "../torrents-search/torrents-search.datasource";
import {inject} from "@angular/core";
import {Apollo} from "apollo-angular";
import {ErrorsService} from "../errors/errors.service";
import {Observable} from "rxjs";
import {AppModule} from "../app.module";

describe('TorrentsTableComponent', () => {
  let component: TorrentsTableComponent;
  let fixture: ComponentFixture<TorrentsTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      ...appConfig,
      imports: [AppModule],
    }).compileComponents();

    fixture = TestBed.createComponent(TorrentsTableComponent);
    component = fixture.componentInstance;
    TestBed.runInInjectionContext(() => {
      component.dataSource = new TorrentsSearchDatasource(inject(Apollo), inject(ErrorsService), new Observable());
    })
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
