import { TimeAgoPipe } from './time-ago.pipe';
import {TestBed} from "@angular/core/testing";
import {appConfig} from "../app.config";

describe('TimeAgoPipe', () => {
  let pipe: TimeAgoPipe;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();
    TestBed.runInInjectionContext(() => {
      pipe = new TimeAgoPipe()
    })
  });

  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });
});
