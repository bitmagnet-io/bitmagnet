import { ObservableQuery } from "@apollo/client/core";
import { filter, OperatorFunction } from "rxjs";

export function filterComplete<T>(): OperatorFunction<
  ObservableQuery.Result<T, "complete" | "empty" | "streaming" | "partial">,
  ObservableQuery.Result<T, "complete">
> {
  return filter(
    (
      result: ObservableQuery.Result<
        T,
        "complete" | "empty" | "streaming" | "partial"
      >,
    ) => result.dataState === "complete",
  );
}
