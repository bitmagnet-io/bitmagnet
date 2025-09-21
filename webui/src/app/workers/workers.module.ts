import { NgModule } from "@angular/core";
import { TranslocoDirective } from "@jsverse/transloco";
import { MatIcon } from "@angular/material/icon";
import { MatTooltip } from "@angular/material/tooltip";
import { AsyncPipe } from "@angular/common";
import { MatDialogClose, MatDialogModule } from "@angular/material/dialog";
import { MatButton, MatIconButton } from "@angular/material/button";
import {
  MatCard,
  MatCardContent,
  MatCardFooter,
  MatCardHeader,
  MatCardTitle,
} from "@angular/material/card";
import { MatGridTile } from "@angular/material/grid-list";
import { MatMenu, MatMenuItem } from "@angular/material/menu";
import { GraphQLModule } from "../graphql/graphql.module";
import { WorkersTableComponent } from "./workers-table.component";
import { WorkersService } from "./workers.service";
import { WorkersCardComponent } from "./workers-card.component";
import { WorkersConfirmActionDialogComponent } from "./workers-confirm-action-dialog.component";

@NgModule({
  imports: [
    GraphQLModule,
    TranslocoDirective,
    MatIcon,
    MatTooltip,
    AsyncPipe,
    MatDialogModule,
    MatButton,
    MatDialogClose,
    MatIconButton,
    MatCard,
    MatCardContent,
    MatCardHeader,
    MatCardTitle,
    MatGridTile,
    MatMenu,
    MatMenuItem,
    MatCardFooter,
  ],
  declarations: [
    WorkersCardComponent,
    WorkersTableComponent,
    WorkersConfirmActionDialogComponent,
  ],
  providers: [WorkersService],
  exports: [WorkersTableComponent, WorkersCardComponent],
})
export class WorkersModule {}
