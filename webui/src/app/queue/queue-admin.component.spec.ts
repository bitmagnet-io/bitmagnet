import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueueAdminComponent } from './queue-admin.component';
import {appConfig} from "../app.config";

describe('QueueAdminComponent', () => {
  let component: QueueAdminComponent;
  let fixture: ComponentFixture<QueueAdminComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();

    fixture = TestBed.createComponent(QueueAdminComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
