import * as generated from "../graphql/generated";

export const objectActionKey = (objectAction: generated.AuthObjectAction) =>
  `${objectAction.namespace}::${objectAction.object}::${objectAction.action}`;
