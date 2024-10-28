import { Component, inject, OnInit } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { Apollo } from "apollo-angular";
import * as generated from "../graphql/generated";
import { GraphQLModule } from "../graphql/graphql.module";
import { AppModule } from "../app.module";
import { TorrentContentComponent } from "./torrent-content.component";
import { contentTypeInfo } from "./content-types";
import { TorrentChipsComponent } from "./torrent-chips.component";

@Component({
  selector: "app-torrent-permalink",
  standalone: true,
  imports: [
    AppModule,
    GraphQLModule,
    TorrentContentComponent,
    TorrentChipsComponent,
  ],
  templateUrl: "./torrent-permalink.component.html",
  styleUrl: "./torrent-permalink.component.scss",
})
export class TorrentPermalinkComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private apollo = inject(Apollo);
  loading = true;
  torrentContent: generated.TorrentContent;

  ngOnInit() {
    this.loading = true;
    this.route.paramMap.subscribe((params) => {
      const infoHash = params.get("infoHash");
      if (typeof infoHash !== "string" || !/^[0-9a-f]{40}$/.test(infoHash)) {
        return this.notFound();
      }
      this.apollo
        .query<
          generated.TorrentContentSearchQuery,
          generated.TorrentContentSearchQueryVariables
        >({
          query: generated.TorrentContentSearchDocument,
          variables: {
            input: {
              infoHashes: [infoHash],
            },
          },
          fetchPolicy: "no-cache",
        })
        .subscribe((result) => {
          const items = result.data.torrentContent.search.items;
          if (items.length === 0) {
            return this.notFound();
          }
          this.torrentContent = items[0];
          this.loading = false;
        });
    });
  }

  private notFound() {
    void this.router.navigate(["/not-found"], {
      skipLocationChange: true,
    });
  }

  protected readonly contentTypeInfo = contentTypeInfo;
}
