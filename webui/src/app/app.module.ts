import { NgModule } from "@angular/core";
import { BrowserModule, DomSanitizer } from "@angular/platform-browser";

import { HttpClientModule } from "@angular/common/http";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { MatIconModule, MatIconRegistry } from "@angular/material/icon";
import { MatButtonModule } from "@angular/material/button";
import { MatToolbarModule } from "@angular/material/toolbar";
import { MatGridListModule } from "@angular/material/grid-list";
import { MatMenuModule } from "@angular/material/menu";
import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { GraphQLModule } from "./graphql/graphql.module";
import { SearchModule } from "./search/search.module";
import { AppErrorsService } from "./app-errors.service";

@NgModule({
  declarations: [AppComponent],
  imports: [
    AppRoutingModule,
    BrowserAnimationsModule,
    BrowserModule,
    GraphQLModule,
    HttpClientModule,
    MatButtonModule,
    MatGridListModule,
    MatIconModule,
    MatMenuModule,
    MatToolbarModule,
    SearchModule,
  ],
  providers: [AppErrorsService],
  bootstrap: [AppComponent],
})
export class AppModule {
  constructor(iconRegistry: MatIconRegistry, domSanitizer: DomSanitizer) {
    iconRegistry.setDefaultFontSetClass(
      "material-icons-outlined",
      "material-symbols-outlined",
    );
    iconRegistry.addSvgIcon(
      "magnet",
      domSanitizer.bypassSecurityTrustResourceUrl("assets/magnet.svg"),
    );
  }
}
