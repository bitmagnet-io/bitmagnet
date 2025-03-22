import { TestBed } from "@angular/core/testing";
import {
  CUSTOM_ELEMENTS_SCHEMA,
  Component,
  NO_ERRORS_SCHEMA,
} from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { TranslocoService, TRANSLOCO_CONFIG } from "@jsverse/transloco";
import { TimeFrameSelectorComponent } from "./time-frame-selector.component";

const translationsMock: Record<string, string> = {
  "dates.title": "Date Filter",
  "dates.time_frame": "Time Frame",
};

@Component({
  template: "",
})
class TestComponent {}

describe("TimeFrameSelectorComponent (No DOM)", () => {
  let component: TimeFrameSelectorComponent;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TestComponent],
      imports: [ReactiveFormsModule],
      schemas: [NO_ERRORS_SCHEMA, CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        TimeFrameSelectorComponent,
        {
          provide: TranslocoService,
          useValue: {
            translate: (key: string) => translationsMock[key] || key,
            selectTranslate: () => (key: string) =>
              translationsMock[key] || key,
            getActiveLang: () => "en",
            load: () => Promise.resolve({}),
            getLangs: () => ["en"],
            events$: {
              pipe: () => ({
                subscribe: () => {},
              }),
            },
          },
        },
        {
          provide: TRANSLOCO_CONFIG,
          useValue: {
            reRenderOnLangChange: false,
            defaultLang: "en",
            availableLangs: ["en"],
            missingHandler: {
              logMissingKey: false,
            },
          },
        },
      ],
    });

    component = TestBed.inject(TimeFrameSelectorComponent);
    // Initialize component manually
    component.ngOnInit();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should initialize with empty time frame", () => {
    expect(component.timeFrameControl.value).toBe("");
  });

  it("should emit when updateTimeFrameAndEmit is called", () => {
    const spy = spyOn(component.timeFrameChanged, "emit");
    component.updateTimeFrameAndEmit("7d");
    expect(spy).toHaveBeenCalled();
  });
});
