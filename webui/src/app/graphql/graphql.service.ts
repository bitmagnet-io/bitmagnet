import { Injectable } from "@angular/core";
import { map, Observable } from "rxjs";
import * as generated from "./generated";

@Injectable()
export class GraphQLService {
  constructor(
    private readonly searchTorrentContentGql: generated.SearchTorrentContentGQL,
    private readonly torrentSetTagsGql: generated.TorrentSetTagsGQL,
  ) {}

  searchTorrentContent(
    input: generated.SearchTorrentContentQueryVariables,
  ): Observable<generated.TorrentContentResult> {
    return this.searchTorrentContentGql
      .fetch(input, {
        fetchPolicy: "no-cache",
      })
      .pipe(map((r) => r.data.search.torrentContent));
  }

  torrentSetTags(
    input: generated.TorrentSetTagsMutationVariables,
  ): Observable<void> {
    return this.torrentSetTagsGql
      .mutate(input, {
        fetchPolicy: "no-cache",
      })
      .pipe(map(() => void 0));
  }
}
