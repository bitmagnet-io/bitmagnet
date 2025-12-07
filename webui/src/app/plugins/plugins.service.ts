import { inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { map, BehaviorSubject } from "rxjs";
import { TranslocoService } from "@jsverse/transloco";
import * as generated from "../graphql/generated";

export type PluginsInfo = Record<string, generated.PluginInfoFragment>;

export class PluginsService {
  private apollo = inject(Apollo);
  private transloco = inject(TranslocoService);

  private pluginsInfo = new BehaviorSubject<PluginsInfo>({});

  pluginsInfo$ = this.pluginsInfo.asObservable();

  constructor() {
    this.query();
    this.transloco.langChanges$.subscribe(() => {
      this.query();
    });
  }

  private query() {
    this.apollo
      .query<generated.PluginsQuery, generated.PluginsQueryVariables>({
        query: generated.PluginsDocument,
        fetchPolicy: "no-cache",
      })
      .pipe(
        map(
          (r): PluginsInfo =>
            r.data!.plugin.list.reduce(
              (acc, next) => ({
                ...acc,
                [next.ref]: next,
              }),
              {},
            ),
        ),
      )
      .subscribe({
        next: (pluginsInfo) => {
          this.pluginsInfo.next(pluginsInfo);
        },
      });
  }
}
