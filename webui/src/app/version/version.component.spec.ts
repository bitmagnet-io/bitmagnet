import { ComponentFixture, TestBed } from "@angular/core/testing";

// import { AppModule } from 'src/app/app.module';
// import { StatusModule } from '../status.module';
import { GraphQLModule } from "../graphql/graphql.module";
import { appConfig } from "../app.config";
import { VersionComponent } from "./version.component";

describe("VersionComponent", () => {
  let component: VersionComponent;
  let fixture: ComponentFixture<VersionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      ...appConfig,
      imports: [GraphQLModule],
    }).compileComponents();

    fixture = TestBed.createComponent(VersionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
