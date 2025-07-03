import { Injectable } from "@angular/core";
import { map, Observable } from "rxjs";
import { Apollo } from "apollo-angular";
import * as generated from "./generated";

@Injectable()
export class GraphQLService {
  constructor(private readonly apollo: Apollo) {}

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

  clientSendToTarget(
    input: generated.ClientSendToMutationVariables,
  ): Observable<void> {
    return this.apollo
      .mutate<
        generated.ClientSendToMutation,
        generated.ClientSendToMutationVariables
      >({
        mutation: generated.ClientSendToDocument,
        variables: input,
        fetchPolicy,
      })
      .pipe(map(() => void 0));
  }

  clentSendToConfig(): Observable<generated.ClientSendToConfigQuery> {
    return this.apollo
      .query<generated.SendToConfigQuery, generated.SendToConfigQueryVariables>(
        {
          query: generated.SendToConfigDocument,
        },
      )
      .pipe(map((r) => r.data.sendToConfig));
  }
}

const fetchPolicy = "no-cache";
