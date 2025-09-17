import * as generated from "../graphql/generated";
import picomatch from "picomatch";

export type ObjectAction = {
  namespace?: string;
  object: string;
  action: string;
};

export type Enforcer = (objAct: ObjectAction) => generated.Permission | null;

export const newEnforcer = (permissions: generated.Permission[]): Enforcer => {
  const matchers = permissions
    .sort((a, b) => {
      if (a.core === b.core) {
        return 0;
      }
      return a.core ? -1 : 1;
    })
    .map((permission) => {
      const { namespace, object, action } = permission.objectAction;
      return {
        permission,
        matcher: {
          namespace: picomatch(namespace),
          object: picomatch(object),
          action: picomatch(action),
        },
      };
    });
  return (objAct) =>
    matchers.find(
      ({ matcher }) =>
        matcher.namespace(objAct.namespace ?? "graphql") &&
        matcher.object(objAct.object) &&
        matcher.action(objAct.action),
    )?.permission || null;
};
