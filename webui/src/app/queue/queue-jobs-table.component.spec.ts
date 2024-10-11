import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueueJobsTableComponent } from './queue-jobs-table.component';
import {appConfig} from "../app.config";
import {QueueJobsDatasource} from "./queue-jobs.datasource";
import {inject} from "@angular/core";
import {Apollo} from "apollo-angular";
import {ErrorsService} from "../errors/errors.service";
import {Observable} from "rxjs";
import {AppModule} from "../app.module";

describe('QueueJobsTableComponent', () => {
  let component: QueueJobsTableComponent;
  let fixture: ComponentFixture<QueueJobsTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      ...appConfig,
      imports: [AppModule],
    }).compileComponents();

    fixture = TestBed.createComponent(QueueJobsTableComponent);
    component = fixture.componentInstance;
    TestBed.runInInjectionContext(() => {
      component.dataSource = new QueueJobsDatasource(inject(Apollo), inject(ErrorsService), new Observable());
    })
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
