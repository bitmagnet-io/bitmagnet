import { Component } from "@angular/core";
import { AppModule } from "../../app.module";
import { TorrentMetricsComponent } from "./torrent-metrics.component";

@Component({
  selector: "app-torrents",
  standalone: true,
  imports: [AppModule, TorrentMetricsComponent],
  templateUrl: "./torrents-dashboard.component.html",
  styleUrl: "./torrents-dashboard.component.scss",
})
export class TorrentsDashboardComponent {}
