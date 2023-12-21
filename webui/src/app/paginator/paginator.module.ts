import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatOptionModule } from "@angular/material/core";
import { MatSelectModule } from "@angular/material/select";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatTooltipModule } from "@angular/material/tooltip";
import { MatDialogModule } from "@angular/material/dialog";
import { MatMenuModule } from "@angular/material/menu";
import { MatAutocompleteModule } from "@angular/material/autocomplete";
import { MatInputModule } from "@angular/material/input";
import { PaginatorComponent } from "./paginator.component";

@NgModule({
  declarations: [PaginatorComponent],
  imports: [
    CommonModule,
    MatFormFieldModule,
    MatOptionModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatTooltipModule,
  ],
  exports: [PaginatorComponent],
})
export class PaginatorModule {}
