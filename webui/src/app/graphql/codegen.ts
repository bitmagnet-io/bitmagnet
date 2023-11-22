import type { CodegenConfig } from "@graphql-codegen/cli";
import * as path from "path"

const config: CodegenConfig = {
  schema: path.resolve(__dirname, "../../../../internal/gql/schema/**/*.graphqls"),
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
          Void: "null",
          Year: "number",
        },
      },
    },
  },
};

export default config;
