import { Injectable } from "@angular/core";
import { map, Observable } from "rxjs";
import * as generated from "./generated";

@Injectable()
export class GraphQLService {
  constructor(
    private readonly torrentContentSearchGql: generated.TorrentContentSearchGQL,
    private readonly torrentSetTagsGql: generated.TorrentSetTagsGQL,
  ) {}

  torrentContentSearch(
    input: generated.TorrentContentSearchQueryVariables,
  ): Observable<generated.TorrentContentSearchResult> {
    return this.torrentContentSearchGql
      .fetch(input, {
        fetchPolicy: "no-cache",
      })
      .pipe(map((r) => r.data.torrentContent.search));
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
