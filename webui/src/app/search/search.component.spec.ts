import { ComponentFixture, TestBed } from "@angular/core/testing";
import { AppModule } from "../app.module";
import { SearchModule } from "./search.module";
import { SearchComponent } from "./search.component";

describe("SearchComponent", () => {
  let component: SearchComponent;
  let fixture: ComponentFixture<SearchComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AppModule, SearchModule],
      declarations: [SearchComponent],
    });
    fixture = TestBed.createComponent(SearchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
