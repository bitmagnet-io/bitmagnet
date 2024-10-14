import { Component } from "@angular/core";
import { RouterOutlet } from "@angular/router";

@Component({
  selector: "app-torrents",
  standalone: true,
  imports: [RouterOutlet],
  templateUrl: "./torrents.component.html",
  styleUrl: "./torrents.component.scss",
})
export class TorrentsComponent {}
