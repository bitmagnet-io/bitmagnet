import { TestBed } from '@angular/core/testing';
import { AppComponent } from './app.component';
import {appConfig} from "./app.config";
import {TranslocoModule} from "@jsverse/transloco";

describe('AppComponent', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      ...appConfig,
      imports: [AppComponent, TranslocoModule],
    }).compileComponents();
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });

  it(`should have the 'bitmagnet' title`, () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app.title).toEqual('bitmagnet');
  });
});
