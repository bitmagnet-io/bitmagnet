import { Injectable } from "@angular/core";
import { map, Observable } from "rxjs";
import { Apollo } from "apollo-angular";
import * as generated from "./generated";

@Injectable()
export class GraphQLService {
  constructor(private readonly apollo: Apollo) {}

  torrentContentSearch(
    input: generated.TorrentContentSearchQueryVariables,
  ): Observable<generated.TorrentContentSearchResult> {
    return this.apollo
      .query<
        generated.TorrentContentSearchQuery,
        generated.TorrentContentSearchQueryVariables
      >({
        query: generated.TorrentContentSearchDocument,
        variables: input,
        fetchPolicy,
      })
      .pipe(map((r) => r.data.torrentContent.search));
  }

  torrentDelete(
    input: generated.TorrentDeleteMutationVariables,
  ): Observable<void> {
    return this.apollo
      .mutate<
        generated.TorrentDeleteMutation,
        generated.TorrentDeleteMutationVariables
      >({
        mutation: generated.TorrentDeleteDocument,
        variables: input,
        fetchPolicy,
      })
      .pipe(map(() => void 0));
  }

  torrentPutTags(
    input: generated.TorrentPutTagsMutationVariables,
  ): Observable<void> {
    return this.apollo
      .mutate<
        generated.TorrentPutTagsMutation,
        generated.TorrentPutTagsMutationVariables
      >({
        mutation: generated.TorrentPutTagsDocument,
        variables: input,
        fetchPolicy,
      })
      .pipe(map(() => void 0));
  }

  torrentSetTags(
    input: generated.TorrentSetTagsMutationVariables,
  ): Observable<void> {
    return this.apollo
      .mutate<
        generated.TorrentSetTagsMutation,
        generated.TorrentSetTagsMutationVariables
      >({
        mutation: generated.TorrentSetTagsDocument,
        variables: input,
        fetchPolicy,
      })
      .pipe(map(() => void 0));
  }

  torrentDeleteTags(
    input: generated.TorrentDeleteTagsMutationVariables,
  ): Observable<void> {
    return this.apollo
      .mutate<
        generated.TorrentDeleteTagsMutation,
        generated.TorrentDeleteTagsMutationVariables
      >({
        mutation: generated.TorrentDeleteTagsDocument,
        variables: input,
        fetchPolicy,
      })
      .pipe(map(() => void 0));
  }

  torrentSuggestTags(
    input: generated.TorrentSuggestTagsQueryVariables,
  ): Observable<generated.TorrentSuggestTagsResult> {
    return this.apollo
      .query<
        generated.TorrentSuggestTagsQuery,
        generated.TorrentSuggestTagsQueryVariables
      >({
        query: generated.TorrentSuggestTagsDocument,
        variables: input,
        fetchPolicy,
      })
      .pipe(map((r) => r.data.torrent.suggestTags));
  }

  systemQuery(): Observable<generated.SystemQuery> {
    return this.apollo
      .query<generated.SystemQueryQuery, generated.SystemQueryQueryVariables>({
        query: generated.SystemQueryDocument,
        fetchPolicy,
      })
      .pipe(map((r) => r.data.system));
  }
}

const fetchPolicy = "no-cache" as const;
