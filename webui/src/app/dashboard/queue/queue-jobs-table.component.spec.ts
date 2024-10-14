import { ComponentFixture, TestBed } from "@angular/core/testing";

import { inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { Observable } from "rxjs";
import { ErrorsService } from "../../errors/errors.service";
import { appConfig } from "../../app.config";
import { AppModule } from "../../app.module";
import { QueueJobsDatasource } from "./queue-jobs.datasource";
import { QueueJobsTableComponent } from "./queue-jobs-table.component";

describe("QueueJobsTableComponent", () => {
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
      component.dataSource = new QueueJobsDatasource(
        inject(Apollo),
        inject(ErrorsService),
        new Observable(),
      );
    });
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
