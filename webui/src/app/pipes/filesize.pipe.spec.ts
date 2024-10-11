import { TestBed } from "@angular/core/testing";
import { appConfig } from "../app.config";
import { FilesizePipe } from "./filesize.pipe";

describe("FilesizePipe", () => {
  let pipe: FilesizePipe;

  beforeEach(async () => {
    await TestBed.configureTestingModule(appConfig).compileComponents();
    TestBed.runInInjectionContext(() => {
      pipe = new FilesizePipe();
    });
  });

  it("create an instance", () => {
    expect(pipe).toBeTruthy();
  });
});
