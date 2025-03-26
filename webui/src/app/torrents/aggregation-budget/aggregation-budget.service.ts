import { inject, Injectable } from "@angular/core";
import { Apollo } from "apollo-angular";
import { BehaviorSubject, Observable, map } from "rxjs";
import {
  AggregationStatsDocument,
  AggregationStatsQuery,
  AggregationStatsQueryVariables,
} from "../../graphql/generated";

export type AggregationBudgetSelection = number | null;

const initialBudget: AggregationBudgetSelection = null;

@Injectable({ providedIn: "root" })
export class AggregationBudgetService {
  private apollo = inject(Apollo);

  budget = initialBudget;
  highBudget = initialBudget;
  private budget$ = new BehaviorSubject<AggregationBudgetSelection>(
    this.budget,
  );

  constructor() {
    this.getCost().subscribe({
      next: (highBudget: AggregationBudgetSelection) => {
        this.highBudget = highBudget;
      },
    });
  }

  public setBudget(budget: AggregationBudgetSelection) {
    this.budget = budget;
    this.budget$.next(budget);
  }

  public getBudget(): Observable<AggregationBudgetSelection> {
    this.setBudget(initialBudget);
    return this.budget$.asObservable();
  }

  public getCost(): Observable<AggregationBudgetSelection> {
    return this.apollo
      .query<AggregationStatsQuery, AggregationStatsQueryVariables>({
        query: AggregationStatsDocument,
        fetchPolicy,
      })
      .pipe(map((r) => r.data.aggregationInfo.highCost));
  }
}

const fetchPolicy = "no-cache";
