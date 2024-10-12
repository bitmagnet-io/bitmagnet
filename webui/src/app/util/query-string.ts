import { Params } from "@angular/router";

export const stringListParam = (
  params: Params,
  key: string,
): string[] | undefined => {
  const str = stringParam(params, key);
  const list = str
    ?.split(",")
    .map((str) => str.trim())
    .filter(Boolean);
  return list?.length ? Array.from(new Set(list)).sort() : undefined;
};

export const stringParam = (
  params: Params,
  key: string,
): string | undefined => {
  return typeof params[key] === "string"
    ? decodeURIComponent(params[key]) || undefined
    : undefined;
};

export const intParam = (params: Params, key: string): number | undefined => {
  if (params && params[key] && /^\d+$/.test(params[key] as string)) {
    return parseInt(params[key] as string);
  }
  return undefined;
};
