import { APOLLO_OPTIONS, ApolloModule } from "apollo-angular";
import { HttpLink } from "apollo-angular/http";
import { NgModule } from "@angular/core";
import { ApolloClientOptions, InMemoryCache } from "@apollo/client/core";
import { graphqlEndpoint } from "../../environments/environment";
import { GraphQLService } from "./graphql.service";

export function createApollo(httpLink: HttpLink): ApolloClientOptions<unknown> {
  return {
    link: httpLink.create({ uri: graphqlEndpoint }),
    cache: new InMemoryCache({
      typePolicies: {
        Query: {
          fields: {
            search: {
              merge(
                existing: Record<string, unknown>,
                incoming: Record<string, unknown>,
              ): Record<string, unknown> {
                return { ...existing, ...incoming };
              },
            },
          },
        },
      },
    }),
  };
}

@NgModule({
  exports: [ApolloModule],
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory: createApollo,
      deps: [HttpLink],
    },
    GraphQLService,
  ],
})
export class GraphQLModule {}
