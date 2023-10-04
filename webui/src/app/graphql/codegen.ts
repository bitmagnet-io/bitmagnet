import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  schema: "http://localhost:3333/graphql",
  documents: "./src/app/graphql/{fragments,mutations,queries}/*.graphql",
  generates: {
    "./src/app/graphql/generated/index.ts": {
      plugins: [
        "typescript",
        "typescript-operations",
        "typescript-apollo-angular",
        {
          add: {
            content: "// THIS FILE IS GENERATED, DO NOT EDIT!\n",
          },
        },
      ],
      config: {
        addExplicitOverride: true,
        enumsAsTypes: true,
        skipTypename: false,
        strictScalars: true,
        scalars: {
          Date: "string",
          DateTime: "string",
          Hash20: "string",
          Year: "number",
        },
      },
    },
    "./src/app/graphql/generated/schema.graphql": {
      plugins: [
        "schema-ast",
        {
          add: {
            content: "# THIS FILE IS GENERATED, DO NOT EDIT!\n",
          },
        },
      ],
    },
  },
};

export default config;
