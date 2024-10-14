import { Component, inject } from "@angular/core";
import { BreakpointsService } from "../layout/breakpoints.service";
import { TorrentsTableComponent } from "../torrents/torrents-table.component";
import { TorrentsBulkActionsComponent } from "../torrents/torrents-bulk-actions.component";
import { PaginatorComponent } from "../paginator/paginator.component";
import { AppModule } from "../app.module";

@Component({
  selector: "app-dashboard",
  standalone: true,
  imports: [
    AppModule,
    PaginatorComponent,
    TorrentsBulkActionsComponent,
    TorrentsTableComponent,
  ],
  templateUrl: "./dashboard.component.html",
  styleUrl: "./dashboard.component.scss",
})
export class DashboardComponent {
  breakpoints = inject(BreakpointsService);
}
