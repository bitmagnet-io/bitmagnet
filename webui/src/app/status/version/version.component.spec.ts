import { ComponentFixture, TestBed } from "@angular/core/testing";

import { StatusModule } from "../status.module";
import { VersionComponent } from "./version.component";

describe("VersionComponent", () => {
  let component: VersionComponent;
  let fixture: ComponentFixture<VersionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StatusModule],
      declarations: [VersionComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(VersionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
