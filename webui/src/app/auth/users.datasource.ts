import { inject } from "@angular/core";
import { Apollo } from "apollo-angular";
import * as generated from "../graphql/generated";
import {
  BehaviorSubject,
  distinctUntilChanged,
  map,
  Observable,
  retry,
  Subscription,
} from "rxjs";
import { DataSource } from "@angular/cdk/collections";
import { CollectionViewer } from "@angular/cdk/collections";

const defaultLimit = 10;

const initialInput: generated.ListUsersInput = {
  pagination: {
    limit: defaultLimit,
  },
};

const pollInterval = 30000;

export class UsersDatasource implements DataSource<generated.User> {
  private apollo = inject(Apollo);
  private input = new BehaviorSubject<generated.ListUsersInput>(initialInput);
  private query = this.apollo.watchQuery<
    generated.UsersQuery,
    generated.UsersQueryVariables
  >({
    query: generated.UsersDocument,
    variables: {
      input: this.input.getValue(),
    },
    pollInterval,
    fetchPolicy: "no-cache",
  });
  private subscription?: Subscription;

  input$ = this.input.asObservable().pipe(distinctUntilChanged());

  result$ = this.query.valueChanges.pipe(
    retry({
      delay: 5000,
    }),
    map((result) => result.data.auth.listUsers),
  );

  users$ = this.result$.pipe(map((result) => result.users));

  constructor() {
    this.input$.subscribe((input) => {
      this.query.setVariables({ input });
      this.query.refetch();
    });
  }

  connect({}: CollectionViewer): Observable<generated.User[]> {
    this.subscription?.unsubscribe();
    this.subscription = this.result$.subscribe();
    return this.users$;
  }

  disconnect(): void {
    this.subscription?.unsubscribe();
    this.subscription = undefined;
  }

  get page(): number {
    return this.input.getValue().pagination?.page ?? 1;
  }

  get limit(): number {
    return this.input.getValue().pagination?.limit ?? defaultLimit;
  }

  handlePagination(event: { pageSize?: number; page?: number }) {
    const input = this.input.getValue();
    this.input.next({
      ...input,
      pagination: {
        ...input.pagination,
        limit: event.pageSize ?? defaultLimit,
        page: event.page ?? 1,
      },
    });
  }
}
