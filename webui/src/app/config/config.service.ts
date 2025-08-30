import { inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import { BehaviorSubject, combineLatest, map } from "rxjs";
import { ConfigParam } from "../graphql/generated";
import * as generated from "../graphql/generated";
import { PluginsInfo, PluginsService } from "../plugins/plugins.service";

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
  plugins: PluginsInfo;
};

const initialConfigState: ConfigState = {
  ...initialConfigParamsState,
  plugins: {},
};

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
      plugins: this.plugins.pluginsInfo$,
    }).subscribe((result) => {
      this.configSubject.next({
        ...result.params,
        plugins: result.plugins,
      });
    });
    this.query();
  }

  private query() {
    this.apollo
      .watchQuery<
        generated.ConfigParamsQuery,
        generated.ConfigParamsQueryVariables
      >({
        query: generated.ConfigParamsDocument,
        fetchPolicy: "no-cache",
        pollInterval,
      })
      .valueChanges.subscribe((result) => {
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
          this.configParamsSubject.next({
            pending: current.pending || result.data!.config.save.pending,
            params: {
              ...current.params,
              [ref]: result.data!.config.save,
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
          this.configParamsSubject.next({
            pending: current.pending || result.data!.config.delete.pending,
            params: {
              ...current.params,
              [ref]: result.data!.config.delete,
            },
          });
        }),
      )
      .subscribe();
  }
}
