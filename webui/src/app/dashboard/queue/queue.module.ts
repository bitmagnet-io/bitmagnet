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
  MatCardHeader,
  MatCardTitle,
} from "@angular/material/card";
import { MatGridTile } from "@angular/material/grid-list";
import { MatMenu, MatMenuItem } from "@angular/material/menu";
import { GraphQLModule } from "../../graphql/graphql.module";
import { ChartComponent } from "../../charting/chart.component";

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
    ChartComponent,
  ],
})
export class QueueModule {}
