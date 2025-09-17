import { inject, Injectable } from "@angular/core";
import { Apollo } from "apollo-angular";
import * as generated from "../graphql/generated";
import {
  catchError,
  from,
  map,
  BehaviorSubject,
  Observable,
  filter,
  take,
  retry,
  EMPTY,
  distinctUntilChanged,
} from "rxjs";
import { AuthTokenService } from "./auth-token.service";
import { newEnforcer, ObjectAction as _objectAction } from "./enforcer";
import { ApolloError } from "@apollo/client/errors";

const pollInterval = 10000;

type LoginResult = generated.LoginMutation["self"]["login"];
type RegisterResult = generated.RegisterMutation["self"]["register"];

export type ObjectAction = _objectAction;

@Injectable({ providedIn: "root" })
export class AuthService {
  private apollo = inject(Apollo);
  private tokenService = inject(AuthTokenService);
  private watchQuery = this.apollo.watchQuery<
    generated.IdentityQuery,
    generated.IdentityQueryVariables
  >({
    query: generated.IdentityDocument,
    fetchPolicy: "no-cache",
    pollInterval,
  });
  private selfSubject = new BehaviorSubject<generated.Self | null>(null);
  private loginErrorSubject = new BehaviorSubject<string | null>(null);
  private registerErrorSubject = new BehaviorSubject<string | null>(null);

  self$ = this.selfSubject.asObservable().pipe(filter((self) => !!self));

  enforcer$ = this.self$.pipe(
    map(({ permissions }) => newEnforcer(permissions)),
  );

  loginError$ = this.loginErrorSubject
    .asObservable()
    .pipe(distinctUntilChanged());

  registerError$ = this.registerErrorSubject
    .asObservable()
    .pipe(distinctUntilChanged());

  token$ = this.tokenService.token$;

  constructor() {
    this.watchQuery.valueChanges
      .pipe(
        map((result) => {
          this.selfSubject.next(result.data.self.identity);
          return result;
        }),
        retry({
          delay: 5000,
        }),
      )
      .subscribe();
  }

  public login(username: string, password: string): Observable<LoginResult> {
    this.resetLoginError();
    return this.apollo
      .mutate<generated.LoginMutation, generated.LoginMutationVariables>({
        mutation: generated.LoginDocument,
        variables: {
          username,
          password,
        },
      })
      .pipe(
        map((result) => {
          const login = result.data!.self.login;
          this.resetLoginError();
          this.tokenService.setToken(login.token);
          this.selfSubject.next({
            user: login.user,
            roles: [login.user.role],
            permissions: login.permissions,
          });

          return login;
        }),
        catchError((error) => {
          if (error instanceof ApolloError) {
            this.loginErrorSubject.next(error.message);
          }
          return EMPTY;
        }),
        take(1),
      );
  }

  public register(input: generated.RegisterInput): Observable<RegisterResult> {
    this.resetRegisterError();
    return this.apollo
      .mutate<generated.RegisterMutation, generated.RegisterMutationVariables>({
        mutation: generated.RegisterDocument,
        variables: {
          input,
        },
      })
      .pipe(
        map((result) => {
          return result.data!.self.register;
        }),
        catchError((error) => {
          if (error instanceof ApolloError) {
            this.registerErrorSubject.next(error.message);
          }
          return EMPTY;
        }),
        take(1),
      );
  }

  public resetLoginError() {
    this.loginErrorSubject.next(null);
  }

  public resetRegisterError() {
    this.registerErrorSubject.next(null);
  }

  public logout(): Observable<generated.Self> {
    this.tokenService.clearToken();
    return from(this.watchQuery.refetch()).pipe(
      take(1),
      map((result) => {
        const self = result.data.self.identity;
        this.selfSubject.next(self);
        return self;
      }),
    );
  }

  public enforce(...objectActions: ObjectAction[]): Observable<boolean> {
    return this.enforcer$.pipe(
      map((enforcer) => objectActions.every(enforcer)),
    );
  }

  public hasRole(role: string): Observable<boolean> {
    return this.self$.pipe(map((self) => self.roles.includes(role)));
  }

  public passwordEntropy(
    password: string,
  ): Observable<generated.PasswordEntropyResult> {
    return this.apollo
      .query<
        generated.PasswordEntropyQuery,
        generated.PasswordEntropyQueryVariables
      >({
        query: generated.PasswordEntropyDocument,
        variables: {
          password,
        },
        fetchPolicy: "no-cache",
      })
      .pipe(map((result) => result.data.self.passwordEntropy));
  }
}
