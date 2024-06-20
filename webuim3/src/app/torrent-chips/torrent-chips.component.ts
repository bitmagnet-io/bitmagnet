import {Component, Input} from '@angular/core';
import * as generated from "../graphql/generated";
import {TranslocoDirective} from "@jsverse/transloco";
import {MatChip, MatChipAvatar, MatChipSet} from "@angular/material/chips";
import {MatIcon} from "@angular/material/icon";

@Component({
  selector: 'app-torrent-chips',
  standalone: true,
  imports: [
    TranslocoDirective,
    MatChip,
    MatChipAvatar,
    MatChipSet,
    MatIcon
  ],
  templateUrl: './torrent-chips.component.html',
  styleUrl: './torrent-chips.component.scss'
})
export class TorrentChipsComponent {
  @Input() torrentContent: generated.TorrentContent;
}
