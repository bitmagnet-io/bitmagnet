import {Component} from '@angular/core';
import {MatCard, MatCardContent, MatCardHeader} from "@angular/material/card";
import {TranslocoDirective} from "@jsverse/transloco";
import {MatIcon} from "@angular/material/icon";
import {MatToolbar} from "@angular/material/toolbar";
import {MatAnchor} from "@angular/material/button";
import {TorrentMetricsComponent} from "./torrent-metrics.component";

@Component({
  selector: 'app-torrents-dashboard',
  standalone: true,
  imports: [
    MatCard,
    MatCardHeader,
    TranslocoDirective,
    MatCardContent,
    MatIcon,
    MatToolbar,
    MatAnchor,
    TorrentMetricsComponent,
  ],
  templateUrl: './torrents-dashboard.component.html',
  styleUrl: './torrents-dashboard.component.scss'
})
export class TorrentsDashboardComponent {
}
