import { Component } from "@angular/core";
import { AppModule } from "../../app.module";
import { DocumentTitleComponent } from "../../layout/document-title.component";
import { TorrentMetricsComponent } from "./torrent-metrics.component";

@Component({
  selector: "app-torrents",
  standalone: true,
  imports: [AppModule, TorrentMetricsComponent, DocumentTitleComponent],
  templateUrl: "./torrents-dashboard.component.html",
  styleUrl: "./torrents-dashboard.component.scss",
})
export class TorrentsDashboardComponent {}
