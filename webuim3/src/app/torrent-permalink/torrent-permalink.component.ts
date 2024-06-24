import { Component, inject, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import {
  MatCard,
  MatCardAvatar,
  MatCardContent,
  MatCardHeader,
  MatCardSubtitle,
  MatCardTitle,
} from '@angular/material/card';
import { MatIcon } from '@angular/material/icon';
import { MatTooltip } from '@angular/material/tooltip';
import { TranslocoDirective } from '@jsverse/transloco';
import { MatProgressBar } from '@angular/material/progress-bar';
import { GraphQLService } from '../graphql/graphql.service';
import * as generated from '../graphql/generated';
import { TorrentContentComponent } from '../torrent-content/torrent-content.component';
import { GraphQLModule } from '../graphql/graphql.module';
import { contentTypeInfo } from '../taxonomy/content-types';
import { TorrentChipsComponent } from '../torrent-chips/torrent-chips.component';

@Component({
  selector: 'app-torrent-permalink',
  standalone: true,
  imports: [
    TorrentContentComponent,
    GraphQLModule,
    MatCardContent,
    MatCard,
    MatCardTitle,
    MatCardHeader,
    MatCardSubtitle,
    MatIcon,
    MatTooltip,
    TranslocoDirective,
    MatCardAvatar,
    TorrentChipsComponent,
    MatProgressBar,
  ],
  templateUrl: './torrent-permalink.component.html',
  styleUrl: './torrent-permalink.component.scss',
})
export class TorrentPermalinkComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private graphQLService = inject(GraphQLService);
  loading = true;
  found = false;
  torrentContent: generated.TorrentContent;

  ngOnInit() {
    this.loading = true;
    this.route.paramMap.subscribe((params) => {
      this.graphQLService
        .torrentContentSearch({
          query: {
            infoHashes: [params.get('infoHash') as string],
          },
        })
        .subscribe((result) => {
          this.torrentContent = result.items[0];
          this.found = result.items.length > 0;
          this.loading = false;
        });
    });
  }

  protected readonly contentTypeInfo = contentTypeInfo;
}
