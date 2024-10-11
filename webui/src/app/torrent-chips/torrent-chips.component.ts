import { Component, Input } from "@angular/core";
import * as generated from "../graphql/generated";
import { AppModule } from "../app.module";

@Component({
  selector: "app-torrent-chips",
  standalone: true,
  imports: [AppModule],
  templateUrl: "./torrent-chips.component.html",
  styleUrl: "./torrent-chips.component.scss",
})
export class TorrentChipsComponent {
  @Input() torrentContent: generated.TorrentContent;
}
