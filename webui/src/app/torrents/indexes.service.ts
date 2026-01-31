import { inject, Injectable } from "@angular/core";
import { Apollo } from "apollo-angular";
import { map } from "rxjs/internal/operators/map";
import { BehaviorSubject, take } from "rxjs";
import * as generated from "../graphql/generated";
import { filterComplete } from "../graphql/util/filter-complete";

@Injectable({ providedIn: "root" })
export class IndexesService {
  private apollo = inject(Apollo);

  private indexes = new BehaviorSubject<generated.IndexesQuery["index"]>({
    default: "",
    infos: [],
  });

  indexes$ = this.indexes.asObservable();

  constructor() {
    this.watchQuery();
  }

  private watchQuery() {
    this.apollo
      .watchQuery<generated.IndexesQuery>({
        query: generated.IndexesDocument,
        fetchPolicy: "cache-first",
      })
      .valueChanges.pipe(
        filterComplete(),
        take(1),
        map((result) => this.indexes.next(result.data.index)),
      )
      .subscribe();
  }
}
