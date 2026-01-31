import { inject, Injectable } from "@angular/core";
import { Apollo } from "apollo-angular";
import { map, BehaviorSubject, take, combineLatest } from "rxjs";
import { TranslocoService } from "@jsverse/transloco";
import * as generated from "../graphql/generated";

export type PluginsInfo = Record<string, generated.PluginInfoFragment>;
export type PluginsLocalizedContent = Record<
  string,
  generated.PluginLocalizedContentFragment
>;

export type PluginInfoWithContent = generated.PluginInfoFragment &
  generated.PluginLocalizedContentFragment;

export type PluginsInfoWithContent = Record<string, PluginInfoWithContent>;

@Injectable({
  providedIn: "root",
})
export class PluginsService {
  private apollo = inject(Apollo);
  private transloco = inject(TranslocoService);

  private pluginsInfo = new BehaviorSubject<PluginsInfo>({});
  private pluginsLocalizedContent =
    new BehaviorSubject<PluginsLocalizedContent>({});

  pluginsInfo$ = this.pluginsInfo.asObservable();
  pluginsLocalizedContent$ = this.pluginsLocalizedContent.asObservable();

  pluginsInfoWithContent$ = combineLatest([
    this.pluginsInfo$,
    this.pluginsLocalizedContent$,
  ]).pipe(
    map(([info, content]): PluginsInfoWithContent => {
      const result: PluginsInfoWithContent = {};
      for (const ref of Object.keys(info)) {
        result[ref] = {
          ...info[ref],
          ...content[ref],
        } as PluginInfoWithContent;
      }
      return result;
    }),
  );

  constructor() {
    this.queryInfos();
    this.queryLocalizedContent();
    this.transloco.langChanges$.subscribe(() => {
      this.queryLocalizedContent();
    });
  }

  private queryInfos() {
    this.apollo
      .query<generated.PluginInfosQuery, generated.PluginInfosQueryVariables>({
        query: generated.PluginInfosDocument,
      })
      .pipe(
        take(1),
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

  private queryLocalizedContent() {
    this.apollo
      .query<
        generated.PluginLocalizedContentQuery,
        generated.PluginLocalizedContentQueryVariables
      >({
        query: generated.PluginLocalizedContentDocument,
      })
      .pipe(
        map(
          (r): PluginsLocalizedContent =>
            r.data!.plugin.list.reduce(
              (acc, next) => ({
                ...acc,
                [next.ref]: next.localizedContent,
              }),
              {},
            ),
        ),
      )
      .subscribe({
        next: (pluginsLocalizedContent) => {
          this.pluginsLocalizedContent.next(pluginsLocalizedContent);
        },
      });
  }
}
