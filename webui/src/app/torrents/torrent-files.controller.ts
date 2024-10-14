import { BehaviorSubject, debounceTime, Observable } from "rxjs";
import * as generated from "../graphql/generated";
import { PageEvent } from "../paginator/paginator.types";

export class TorrentFilesController {
  private controlsSubject: BehaviorSubject<TorrentFilesControls>;
  controls$: Observable<TorrentFilesControls>;

  private variablesSubject: BehaviorSubject<generated.TorrentFilesQueryVariables>;
  variables$: Observable<generated.TorrentFilesQueryVariables>;

  constructor(infoHash: string) {
    const ctrl: TorrentFilesControls = {
      infoHash,
      limit: 10,
      page: 1,
    };
    this.controlsSubject = new BehaviorSubject(ctrl);
    this.controls$ = this.controlsSubject.asObservable();
    this.controls$.pipe(debounceTime(100)).subscribe((ctrl) => {
      const currentParams = this.variablesSubject.getValue();
      const nextParams = controlsToQueryVariables(ctrl);
      if (JSON.stringify(currentParams) !== JSON.stringify(nextParams)) {
        this.variablesSubject.next(nextParams);
      }
    });
    this.variablesSubject = new BehaviorSubject(controlsToQueryVariables(ctrl));
    this.variables$ = this.variablesSubject.asObservable();
  }

  update(fn: (c: TorrentFilesControls) => TorrentFilesControls) {
    const ctrl = this.controlsSubject.getValue();
    const next = fn(ctrl);
    if (JSON.stringify(ctrl) !== JSON.stringify(next)) {
      this.controlsSubject.next(next);
    }
  }

  handlePageEvent(event: PageEvent) {
    this.update((ctrl) => ({
      ...ctrl,
      limit: event.pageSize,
      page: event.page,
    }));
  }
}

export type TorrentFilesControls = {
  infoHash: string;
  limit: number;
  page: number;
};

const controlsToQueryVariables = (
  ctrl: TorrentFilesControls,
): generated.TorrentFilesQueryVariables => ({
  input: {
    infoHashes: [ctrl.infoHash],
    limit: ctrl.limit,
    page: ctrl.page,
    totalCount: true,
    hasNextPage: false,
  },
});
