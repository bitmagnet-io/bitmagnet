import { inject, Injectable } from "@angular/core";
import { Apollo } from "apollo-angular";
import { BehaviorSubject, combineLatest, map } from "rxjs";
import { ConfigParam } from "../graphql/generated";
import * as generated from "../graphql/generated";
import {
  PluginsInfoWithContent,
  PluginsService,
} from "../plugins/plugins.service";
import { filterComplete } from "../graphql/util/filter-complete";

const pollInterval = 10000;

type ConfigParamsState = {
  pending: boolean;
  params: Record<string, ConfigParam>;
};

const initialConfigParamsState: ConfigParamsState = {
  pending: false,
  params: {},
};

export type ConfigState = ConfigParamsState & {
  plugins: PluginsInfoWithContent;
};

const initialConfigState: ConfigState = {
  ...initialConfigParamsState,
  plugins: {},
};

@Injectable({
  providedIn: "root",
})
export class ConfigService {
  private apollo = inject(Apollo);
  private plugins = inject(PluginsService);

  private configParamsSubject = new BehaviorSubject<ConfigParamsState>(
    initialConfigParamsState,
  );

  private configSubject = new BehaviorSubject<ConfigState>(initialConfigState);

  config$ = this.configSubject.asObservable();

  constructor() {
    combineLatest({
      params: this.configParamsSubject.asObservable(),
      plugins: this.plugins.pluginsInfoWithContent$,
    }).subscribe((result) => {
      this.configSubject.next({
        ...result.params,
        plugins: result.plugins,
      });
    });
    this.watchQuery();
  }

  private watchQuery() {
    this.apollo
      .watchQuery<
        generated.ConfigParamsQuery,
        generated.ConfigParamsQueryVariables
      >({
        query: generated.ConfigParamsDocument,
        fetchPolicy: "no-cache",
        pollInterval,
      })
      .valueChanges.pipe(filterComplete())
      .subscribe((result) => {
        this.configParamsSubject.next({
          pending: result.data.config.pending,
          params: result.data.config.params.reduce(
            (acc, next) => ({
              ...acc,
              [next.ref]: next,
            }),
            {},
          ),
        });
      });
  }

  public save(ref: string, value: unknown) {
    this.apollo
      .mutate<
        generated.ConfigSaveMutation,
        generated.ConfigSaveMutationVariables
      >({
        mutation: generated.ConfigSaveDocument,
        variables: {
          ref,
          value,
        },
      })
      .pipe(
        map((result) => {
          const current = this.configParamsSubject.getValue();
          const { pending, value, source } = result.data!.config.save;
          this.configParamsSubject.next({
            pending: current.pending || pending,
            params: {
              ...current.params,
              [ref]: {
                ...current.params[ref],
                pending,
                value,
                source,
              },
            },
          });
        }),
      )
      .subscribe();
  }

  public delete(ref: string) {
    this.apollo
      .mutate<
        generated.ConfigDeleteMutation,
        generated.ConfigDeleteMutationVariables
      >({
        mutation: generated.ConfigDeleteDocument,
        variables: {
          ref,
        },
      })
      .pipe(
        map((result) => {
          const current = this.configParamsSubject.getValue();
          const { pending, value, source } = result.data!.config.delete;
          this.configParamsSubject.next({
            pending: current.pending || pending,
            params: {
              ...current.params,
              [ref]: {
                ...current.params[ref],
                pending,
                value,
                source,
              },
            },
          });
        }),
      )
      .subscribe();
  }
}
