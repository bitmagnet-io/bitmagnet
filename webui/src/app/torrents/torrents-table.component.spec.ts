import { ComponentFixture, TestBed } from "@angular/core/testing";

import { inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { Observable } from "rxjs";
import { appConfig } from "../app.config";
import { ErrorsService } from "../errors/errors.service";
import { AppModule } from "../app.module";
import { TorrentsSearchDatasource } from "./torrents-search.datasource";
import { TorrentsTableComponent } from "./torrents-table.component";

describe("TorrentsTableComponent", () => {
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
      component.dataSource = new TorrentsSearchDatasource(
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
