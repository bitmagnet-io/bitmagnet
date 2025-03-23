export type TimeFrame = {
  startDate: Date;
  endDate: Date;
  expression: string;
  isValid: boolean;
  error?: string;
};

/**
 * Keywords for relative times
 */
type UnitSingular =
  | "second"
  | "minute"
  | "hour"
  | "day"
  | "week"
  | "month"
  | "year";
type UnitPlural =
  | "seconds"
  | "minutes"
  | "hours"
  | "days"
  | "weeks"
  | "months"
  | "years";
type Unit = UnitSingular | UnitPlural;

// Unit abbreviations
type UnitAbbreviation = "s" | "m" | "h" | "d" | "w" | "mo" | "y";

// Map of abbreviations to full unit names
const unitMap: Record<UnitAbbreviation, Unit> = {
  s: "seconds",
  m: "minutes",
  h: "hours",
  d: "days",
  w: "weeks",
  mo: "months",
  y: "years",
};

// Special time frames
type SpecialTimeFrame =
  | "today"
  | "yesterday"
  | "this week"
  | "last week"
  | "this month"
  | "last month"
  | "this year"
  | "last year";

/**
 * Creates a time frame with error
 */
function errorTimeFrame(expression: string, error: string): TimeFrame {
  return {
    startDate: new Date(),
    endDate: new Date(),
    expression,
    isValid: false,
    error,
  };
}

/**
 * Parse a relative time expression like "3h", "2d", "1w", etc.
 * @param expression The expression to parse
 */
function parseRelativeTime(expression: string): TimeFrame {
  // Extract numeric value and unit from expression
  const match = expression.match(/^(\d+)([a-z]+)$/i);
  if (!match) {
    return errorTimeFrame(
      expression,
      "Invalid relative time format. Expected format: '3h', '2d', etc.",
    );
  }

  const value = parseInt(match[1], 10);
  const unitAbbr = match[2].toLowerCase() as UnitAbbreviation;

  if (!(unitAbbr in unitMap)) {
    return errorTimeFrame(
      expression,
      `Unknown time unit: '${unitAbbr}'. Valid units: s, m, h, d, w, mo, y`,
    );
  }

  const unit = unitMap[unitAbbr];
  const endDate = new Date();
  const startDate = new Date();

  // Calculate the start date based on the unit and value
  switch (unit) {
    case "seconds":
      startDate.setSeconds(startDate.getSeconds() - value);
      break;
    case "minutes":
      startDate.setMinutes(startDate.getMinutes() - value);
      break;
    case "hours":
      startDate.setHours(startDate.getHours() - value);
      break;
    case "days":
      startDate.setDate(startDate.getDate() - value);
      break;
    case "weeks":
      startDate.setDate(startDate.getDate() - value * 7);
      break;
    case "months":
      startDate.setMonth(startDate.getMonth() - value);
      break;
    case "years":
      startDate.setFullYear(startDate.getFullYear() - value);
      break;
  }

  return {
    startDate,
    endDate,
    expression,
    isValid: true,
  };
}

/**
 * Parse a special time expression like "today", "this week", etc.
 * @param expression The expression to parse
 */
function parseSpecialTimeFrame(expression: string): TimeFrame {
  const now = new Date();
  const startDate = new Date(now);
  const endDate = new Date(now);

  // Default - end of today
  endDate.setHours(23, 59, 59, 999);

  switch (expression.toLowerCase()) {
    case "today":
      startDate.setHours(0, 0, 0, 0);
      break;

    case "yesterday":
      startDate.setDate(startDate.getDate() - 1);
      startDate.setHours(0, 0, 0, 0);
      endDate.setDate(endDate.getDate() - 1);
      break;

    case "this week":
      // Start of week (Sunday or Monday depending on locale)
      const dayOfWeek = startDate.getDay(); // 0 is Sunday
      const diff = dayOfWeek === 0 ? 6 : dayOfWeek - 1; // Adjust for Monday as start of week
      startDate.setDate(startDate.getDate() - diff);
      startDate.setHours(0, 0, 0, 0);
      break;

    case "last week":
      // Start of previous week
      const currentDayOfWeek = startDate.getDay();
      const diffToStartOfThisWeek =
        currentDayOfWeek === 0 ? 6 : currentDayOfWeek - 1;
      startDate.setDate(startDate.getDate() - diffToStartOfThisWeek - 7);
      startDate.setHours(0, 0, 0, 0);

      // End of previous week
      endDate.setDate(endDate.getDate() - diffToStartOfThisWeek - 1);
      endDate.setHours(23, 59, 59, 999);
      break;

    case "this month":
      startDate.setDate(1);
      startDate.setHours(0, 0, 0, 0);

      // End of month
      endDate.setMonth(endDate.getMonth() + 1);
      endDate.setDate(0);
      break;

    case "last month":
      // Start of previous month
      startDate.setMonth(startDate.getMonth() - 1);
      startDate.setDate(1);
      startDate.setHours(0, 0, 0, 0);

      // End of previous month
      endDate.setDate(0);
      break;

    case "this year":
      startDate.setMonth(0, 1);
      startDate.setHours(0, 0, 0, 0);

      // End of year
      endDate.setMonth(11, 31);
      break;

    case "last year":
      // Start of previous year
      startDate.setFullYear(startDate.getFullYear() - 1);
      startDate.setMonth(0, 1);
      startDate.setHours(0, 0, 0, 0);

      // End of previous year
      endDate.setFullYear(endDate.getFullYear() - 1);
      endDate.setMonth(11, 31);
      break;

    default:
      return errorTimeFrame(
        expression,
        `Unknown special time frame: '${expression}'`,
      );
  }

  return {
    startDate,
    endDate,
    expression,
    isValid: true,
  };
}

/**
 * Try to parse a date string in various formats
 */
function tryParseDate(dateStr: string): Date | null {
  // Try to parse ISO dates like "2023-01-01"
  const isoDate = new Date(dateStr);
  if (!isNaN(isoDate.getTime())) {
    return isoDate;
  }

  // Try to parse more human-readable formats like "Jan 1, 2023"
  const formats = [
    // Add more date formats as needed
    /^(\w{3})\s+(\d{1,2}),?\s+(\d{4})$/, // Jan 1, 2023
    /^(\d{1,2})\s+(\w{3})\s+(\d{4})$/, // 1 Jan 2023
    /^(\d{1,2})\/(\d{1,2})\/(\d{4})$/, // MM/DD/YYYY or DD/MM/YYYY
    /^(\d{1,2})-(\d{1,2})-(\d{4})$/, // MM-DD-YYYY or DD-MM-YYYY
  ];

  // Try each format
  for (const format of formats) {
    const match = dateStr.match(format);
    if (match) {
      const parsed = new Date(dateStr);
      if (!isNaN(parsed.getTime())) {
        return parsed;
      }
    }
  }

  return null;
}

/**
 * Parse an absolute time range with a start and end date
 */
function parseAbsoluteTimeRange(expression: string): TimeFrame {
  const parts = expression.split(" to ");
  if (parts.length !== 2) {
    return errorTimeFrame(
      expression,
      "Invalid absolute time range. Expected format: 'start to end'",
    );
  }

  const startStr = parts[0].trim();
  const endStr = parts[1].trim();

  const startDate = tryParseDate(startStr);
  const endDate = tryParseDate(endStr);

  if (!startDate) {
    return errorTimeFrame(
      expression,
      `Could not parse start date: '${startStr}'`,
    );
  }

  if (!endDate) {
    return errorTimeFrame(expression, `Could not parse end date: '${endStr}'`);
  }

  // For dates without time components, set end date to end of day
  if (
    endDate.getHours() === 0 &&
    endDate.getMinutes() === 0 &&
    endDate.getSeconds() === 0
  ) {
    endDate.setHours(23, 59, 59, 999);
  }

  return {
    startDate,
    endDate,
    expression,
    isValid: true,
  };
}

/**
 * Parse a time frame expression and return a TimeFrame object
 * @param expression The time frame expression
 */
export function parseTimeFrame(expression: string): TimeFrame {
  if (!expression || expression.trim() === "") {
    // Return current date for both start and end when empty
    const now = new Date();
    return {
      startDate: now,
      endDate: now,
      expression: "",
      isValid: true,
    };
  }

  expression = expression.trim();

  // Check if it's a special time frame
  const specialTimeFrames: SpecialTimeFrame[] = [
    "today",
    "yesterday",
    "this week",
    "last week",
    "this month",
    "last month",
    "this year",
    "last year",
  ];

  if (
    specialTimeFrames.includes(expression.toLowerCase() as SpecialTimeFrame)
  ) {
    return parseSpecialTimeFrame(expression);
  }

  // Check if it's a relative time (e.g., "3h", "2d")
  const relativeTimePattern = /^\d+[a-z]+$/i;
  if (relativeTimePattern.test(expression)) {
    return parseRelativeTime(expression);
  }

  // Check if it contains "to" for absolute time range
  if (expression.includes(" to ")) {
    return parseAbsoluteTimeRange(expression);
  }

  // Try to parse as a single date
  const date = tryParseDate(expression);
  if (date) {
    const endDate = new Date(date);
    endDate.setHours(23, 59, 59, 999);

    return {
      startDate: date,
      endDate,
      expression,
      isValid: true,
    };
  }

  // If nothing matches, return an error
  return errorTimeFrame(
    expression,
    "Could not parse time frame. Try formats like '3h', 'today', or 'Jan 1, 2023 to Jan 2, 2023'",
  );
}

/**
 * Get a human-friendly description of a TimeFrame
 */
export function formatTimeFrameDescription(timeFrame: TimeFrame): string {
  if (!timeFrame.isValid) {
    return `Invalid: ${timeFrame.error}`;
  }

  // For special expressions, just return the expression
  const specialExpressions = [
    "today",
    "yesterday",
    "this week",
    "last week",
    "this month",
    "last month",
    "this year",
    "last year",
  ];

  if (specialExpressions.includes(timeFrame.expression.toLowerCase())) {
    return timeFrame.expression;
  }

  // For relative expressions (like 3h, 2d), make it more readable
  const relativeMatch = timeFrame.expression.match(/^(\d+)([a-z]+)$/i);
  if (relativeMatch) {
    const value = relativeMatch[1];
    const unitAbbr = relativeMatch[2].toLowerCase();
    const unit = unitMap[unitAbbr as UnitAbbreviation] || unitAbbr;

    return `Last ${value} ${unit}`;
  }

  // For absolute ranges, format nicely
  const formatDate = (date: Date) => {
    return date.toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return `${formatDate(timeFrame.startDate)} to ${formatDate(timeFrame.endDate)}`;
}

/**
 * Common time frame presets
 */
export const timeFramePresets = [
  { label: "Last hour", value: "1h" },
  { label: "Last 6 hours", value: "6h" },
  { label: "Last 12 hours", value: "12h" },
  { label: "Last 24 hours", value: "24h" },
  { label: "Last 2 days", value: "2d" },
  { label: "Last 7 days", value: "7d" },
  { label: "Last 30 days", value: "30d" },
  { label: "Last 90 days", value: "90d" },
  { label: "Today", value: "today" },
  { label: "Yesterday", value: "yesterday" },
  { label: "This week", value: "this week" },
  { label: "Last week", value: "last week" },
  { label: "This month", value: "this month" },
  { label: "Last month", value: "last month" },
  { label: "This year", value: "this year" },
  { label: "Last year", value: "last year" },
];
