import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { VersionComponent } from "./version/version.component";

@NgModule({
  declarations: [VersionComponent],
  imports: [CommonModule],
  exports: [VersionComponent],
})
export class StatusModule {}
