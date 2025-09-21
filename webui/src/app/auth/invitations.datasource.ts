import { inject } from "@angular/core";
import { Apollo } from "apollo-angular";
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
import * as generated from "../graphql/generated";

const defaultLimit = 10;

const initialInput: generated.ListInvitationsInput = {
  pagination: {
    limit: defaultLimit,
  },
};

const pollInterval = 30000;

export class InvitationsDatasource implements DataSource<generated.Invitation> {
  private apollo = inject(Apollo);
  private input = new BehaviorSubject<generated.ListInvitationsInput>(
    initialInput,
  );
  private query = this.apollo.watchQuery<
    generated.InvitationsQuery,
    generated.InvitationsQueryVariables
  >({
    query: generated.InvitationsDocument,
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
    map((result) => result.data.auth.listInvitations),
  );

  invitations$ = this.result$.pipe(map((result) => result.invitations));

  constructor() {
    this.input$.subscribe((input) => {
      void this.query.setVariables({ input });
      void this.query.refetch();
    });
  }

  connect({}: CollectionViewer): Observable<generated.Invitation[]> {
    this.subscription?.unsubscribe();
    this.subscription = this.result$.subscribe();
    return this.invitations$;
  }

  disconnect(): void {
    this.subscription?.unsubscribe();
    this.subscription = undefined;
  }

  refresh() {
    void this.query.refetch();
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
