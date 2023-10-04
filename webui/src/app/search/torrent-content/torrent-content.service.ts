import { Injectable } from "@angular/core";
import { map, Observable } from "rxjs";
import * as generated from "../../graphql/generated";

@Injectable()
export class TorrentContentService {
  constructor(private readonly gql: generated.TorrentContentSearchGQL) {}

  search(
    input: generated.TorrentContentSearchQueryVariables,
  ): Observable<generated.TorrentContentResult> {
    return this.gql
      .fetch(input, {
        fetchPolicy: "no-cache",
      })
      .pipe(map((r) => r.data.search.torrentContent));
  }
}
