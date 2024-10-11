import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueueJobsComponent } from './queue-jobs.component';
import {appConfig} from "../app.config";

describe('QueueJobsComponent', () => {
  let component: QueueJobsComponent;
  let fixture: ComponentFixture<QueueJobsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();

    fixture = TestBed.createComponent(QueueJobsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
