import { Injectable } from "@angular/core";
import { map, Observable } from "rxjs";
import * as generated from "./generated";

@Injectable()
export class GraphQLService {
  constructor(
    private readonly torrentContentSearchGql: generated.TorrentContentSearchGQL,
    private readonly torrentSetTagsGql: generated.TorrentSetTagsGQL,
    private readonly torrentSuggestTagsGql: generated.TorrentSuggestTagsGQL,
  ) {}

  torrentContentSearch(
    input: generated.TorrentContentSearchQueryVariables,
  ): Observable<generated.TorrentContentSearchResult> {
    return this.torrentContentSearchGql
      .fetch(input, fetchPolicy)
      .pipe(map((r) => r.data.torrentContent.search));
  }

  torrentSetTags(
    input: generated.TorrentSetTagsMutationVariables,
  ): Observable<void> {
    return this.torrentSetTagsGql
      .mutate(input, fetchPolicy)
      .pipe(map(() => void 0));
  }

  torrentSuggestTags(
    input: generated.TorrentSuggestTagsQueryVariables,
  ): Observable<generated.TorrentSuggestTagsResult> {
    return this.torrentSuggestTagsGql
      .fetch(input, fetchPolicy)
      .pipe(map((r) => r.data.torrent.suggestTags));
  }
}

const fetchPolicy = {
  fetchPolicy: "no-cache",
} as const;
