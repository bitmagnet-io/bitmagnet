import { parseTimeFrame, formatTimeFrameDescription } from "./parse-timeframe";

describe("parseTimeFrame", () => {
  it("should handle empty", () => {
    const emptyResult = parseTimeFrame("");
    expect(emptyResult.isValid).toBeTrue();
    expect(emptyResult.expression).toBe("");
  });

  it("should parse relative time expressions", () => {
    const testCases = [
      { input: "1h", expectValid: true },
      { input: "24h", expectValid: true },
      { input: "7d", expectValid: true },
      { input: "2w", expectValid: true },
      { input: "6m", expectValid: true },
      { input: "1y", expectValid: true },
      { input: "0d", expectValid: true },
      { input: "-1d", expectValid: false }, // Negative values not supported
      { input: "1x", expectValid: false }, // Invalid unit
    ];

    for (const tc of testCases) {
      const result = parseTimeFrame(tc.input);
      expect(result.isValid).toBe(tc.expectValid);
      if (tc.expectValid) {
        expect(result.expression).toBe(tc.input);
        expect(result.startDate).toBeDefined();
        expect(result.endDate).toBeDefined();
        expect(result.startDate <= result.endDate).toBeTrue();
      }
    }
  });

  it("should parse special time expressions", () => {
    const specialCases = [
      "today",
      "yesterday",
      "this week",
      "last week",
      "this month",
      "last month",
      "this year",
      "last year",
    ];

    for (const expression of specialCases) {
      const result = parseTimeFrame(expression);
      expect(result.isValid).toBeTrue();
      expect(result.expression).toBe(expression);
      expect(result.startDate).toBeDefined();
      expect(result.endDate).toBeDefined();
      expect(result.startDate <= result.endDate).toBeTrue();
    }
  });

  it("should parse date ranges", () => {
    const result = parseTimeFrame("2023-01-01 to 2023-01-31");
    expect(result.isValid).toBeTrue();
    expect(result.expression).toBe("2023-01-01 to 2023-01-31");

    // Check properties without exact values since some browsers/environments
    // may parse differently
    expect(result.startDate instanceof Date).toBeTrue();
    expect(result.endDate instanceof Date).toBeTrue();

    // Just verify dates are correctly ordered
    expect(result.startDate <= result.endDate).toBeTrue();
  });

  it("should handle human-readable dates and formats", () => {
    // Testing various date formats
    const testCases = [
      { input: "Jan 1, 2023 to Jan 31, 2023", expectValid: true },
      { input: "2023-01-01", expectValid: true }, // Single date
      { input: "Jan 1, 2023", expectValid: true }, // Single date in human format
      { input: "1/1/2023", expectValid: true }, // MM/DD/YYYY
      { input: "2023/01/01", expectValid: true }, // YYYY/MM/DD
      { input: "invalid date", expectValid: false },
    ];

    for (const tc of testCases) {
      const result = parseTimeFrame(tc.input);
      expect(result.isValid).toBe(tc.expectValid);
      if (tc.expectValid) {
        expect(result.expression).toBe(tc.input);
      }
    }
  });

  it("should handle invalid inputs", () => {
    const result = parseTimeFrame("completely invalid");
    expect(result.isValid).toBeFalse();
    expect(result.error).toBeDefined();
  });
});

describe("formatTimeFrameDescription", () => {
  it("should format relative times", () => {
    const timeFrame = parseTimeFrame("7d");
    const result = formatTimeFrameDescription(timeFrame);
    expect(result).toBe("Last 7 days");
  });

  it("should pass through special expressions", () => {
    const timeFrame = parseTimeFrame("today");
    const result = formatTimeFrameDescription(timeFrame);
    expect(result).toBe("today");
  });

  it("should format date ranges nicely", () => {
    const timeFrame = parseTimeFrame("2023-01-01 to 2023-01-31");
    const result = formatTimeFrameDescription(timeFrame);

    // The exact format may depend on locale, but should contain the months and dates
    expect(result).toContain("Jan");
    expect(result).toContain("2023");
    expect(result).toContain("to");
  });

  it("should handle invalid time frames", () => {
    const timeFrame = parseTimeFrame("invalid");
    const result = formatTimeFrameDescription(timeFrame);
    expect(result).toContain("Invalid");
  });
});
