import { inject, Injectable } from "@angular/core";
import { Apollo } from "apollo-angular";
import { map } from "rxjs/internal/operators/map";
import { BehaviorSubject, catchError, take } from "rxjs";
import * as generated from "../graphql/generated";
import { filterComplete } from "../graphql/util/filter-complete";
import { TranslocoService } from "@jsverse/transloco";
import { ErrorsService } from "../errors/errors.service";

@Injectable({ providedIn: "root" })
export class TargetsService {
  private apollo = inject(Apollo);
  private transloco = inject(TranslocoService);
  private errors = inject(ErrorsService);

  private targets = new BehaviorSubject<
    generated.TargetsQuery["target"]["targets"]
  >([]);

  targets$ = this.targets.asObservable();

  private watchQuery = this.apollo.watchQuery<generated.TargetsQuery>({
    query: generated.TargetsDocument,
    fetchPolicy: "cache-first",
    pollInterval: 30_000,
  });

  constructor() {
    this.watchQuery.valueChanges
      .pipe(
        filterComplete(),
        map((result) => this.targets.next(result.data.target.targets)),
      )
      .subscribe();
    this.transloco.langChanges$.subscribe(() => {
      this.watchQuery.refetch();
    });
  }

  send(request: SendRequest) {
    this.apollo
      .mutate<
        generated.SendTorrentMutation,
        generated.SendTorrentMutationVariables
      >({
        mutation: generated.SendTorrentDocument,
        variables: {
          input: {
            index: request.index,
            infoHashes: request.infoHashes,
            target: request.target,
            data: request.data,
          },
        },
      })
      .pipe(
        take(1),
        catchError((error) => {
          this.errors.addError(error);
          throw error;
        }),
      )
      .subscribe();
  }
}

export type SendRequest = {
  index?: string;
  infoHashes: string[];
  target: string;
  data: unknown;
};
