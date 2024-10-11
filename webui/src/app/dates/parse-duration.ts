/**
 * @description A module for parsing ISO8601 durations
 * Implementation copied from https://github.com/tolu/ISO8601-duration
 */

/**
 * The pattern used for parsing ISO8601 duration (PnYnMnWnDTnHnMnS).
 */
// PnYnMnWnDTnHnMnS
const numbers = "\\d+";
const fractionalNumbers = "".concat(numbers, "(?:[\\.,]").concat(numbers, ")?");
const datePattern = "("
  .concat(numbers, "Y)?(")
  .concat(numbers, "M)?(")
  .concat(numbers, "W)?(")
  .concat(numbers, "D)?");
const timePattern = "T("
  .concat(fractionalNumbers, "H)?(")
  .concat(fractionalNumbers, "M)?(")
  .concat(fractionalNumbers, "S)?");
const iso8601 = "P(?:".concat(datePattern, "(?:").concat(timePattern, ")?)");
const objMap = [
  "years",
  "months",
  "weeks",
  "days",
  "hours",
  "minutes",
  "seconds",
];

export type Duration = {
  years: number;
  months: number;
  weeks: number;
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
};

const defaultDuration: Duration = {
  years: 0,
  months: 0,
  weeks: 0,
  days: 0,
  hours: 0,
  minutes: 0,
  seconds: 0,
};

/**
 * The ISO8601 regex for matching / testing durations
 */
export const durationPattern = new RegExp(iso8601);

/** Parse PnYnMnDTnHnMnS format to object */
export const parseDuration = function (durationString: string): Duration {
  const matches = durationString.replace(/,/g, ".").match(durationPattern);
  if (!matches) {
    throw new RangeError("invalid duration: ".concat(durationString));
  }
  // Slice away first entry in match-array (the input string)
  const slicedMatches = matches.slice(1);
  if (
    slicedMatches.filter(function (v) {
      return v != null;
    }).length === 0
  ) {
    throw new RangeError("invalid duration: ".concat(durationString));
  }
  // Check only one fraction is used
  if (
    slicedMatches.filter(function (v) {
      return /\./.test(v || "");
    }).length > 1
  ) {
    throw new RangeError("only the smallest unit can be fractional");
  }
  return slicedMatches.reduce(function (prev, next, idx) {
    Object.assign(prev, { [objMap[idx]]: parseFloat(next || "0") || 0 });
    return prev;
  }, {}) as Duration;
};

/** Convert ISO8601 duration object to an end Date. */
const end = function (durationInput: Duration, startDate?: Date): Date {
  if (!startDate) {
    startDate = new Date();
  }
  const duration = Object.assign({}, defaultDuration, durationInput);
  // Create two equal timestamps, add duration to 'then' and return time difference
  const timestamp = startDate.getTime();
  const then = new Date(timestamp);
  then.setFullYear(then.getFullYear() + duration.years);
  then.setMonth(then.getMonth() + duration.months);
  then.setDate(then.getDate() + duration.days);
  // set time as milliseconds to get fractions working for minutes/hours
  const hoursInMs = duration.hours * 3600 * 1000;
  const minutesInMs = duration.minutes * 60 * 1000;
  then.setMilliseconds(
    then.getMilliseconds() + duration.seconds * 1000 + hoursInMs + minutesInMs,
  );
  // Special case weeks
  then.setDate(then.getDate() + duration.weeks * 7);
  return then;
};

/** Convert ISO8601 duration object to seconds */
export const durationToSeconds = function (
  durationInput: Duration,
  startDate?: Date,
): number {
  if (!startDate) {
    startDate = new Date();
  }
  const duration = Object.assign({}, defaultDuration, durationInput);
  const timestamp = startDate.getTime();
  const now = new Date(timestamp);
  const then = end(duration, now);
  // Account for timezone offset between start and end date
  const tzStart = startDate.getTimezoneOffset();
  const tzEnd = then.getTimezoneOffset();
  const tzOffsetSeconds = (tzStart - tzEnd) * 60;
  const seconds = (then.getTime() - now.getTime()) / 1000;
  return seconds + tzOffsetSeconds;
};
