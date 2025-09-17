import { inject, Injectable } from "@angular/core";
import { Apollo } from "apollo-angular";
import { map, Observable } from "rxjs";
import * as generated from "../graphql/generated";

@Injectable({ providedIn: "root" })
export class RolesService {
  private apollo = inject(Apollo);

  objectActions$ = this.apollo
    .query<
      generated.AuthObjectActionsQuery,
      generated.AuthObjectActionsQueryVariables
    >({
      query: generated.AuthObjectActionsDocument,
      fetchPolicy: "no-cache",
    })
    .pipe(map((result) => result.data.auth.listObjectActions));

  listRoles(): Observable<generated.Role[]> {
    return this.apollo
      .query<generated.RolesQuery, generated.RolesQueryVariables>({
        query: generated.RolesDocument,
        fetchPolicy: "no-cache",
      })
      .pipe(map((result) => result.data.auth.listRoles));
  }

  putRole(
    input: generated.PutRoleMutationVariables,
  ): Observable<generated.Role> {
    return this.apollo
      .mutate<generated.PutRoleMutation, generated.PutRoleMutationVariables>({
        mutation: generated.PutRoleDocument,
        variables: input,
      })
      .pipe(map((result) => result.data!.auth.putRole));
  }
}
