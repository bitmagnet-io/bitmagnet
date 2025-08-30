// node_modules/date-fns/locale/ar/_lib/formatDistance.js
var formatDistanceLocale = {
  lessThanXSeconds: {
    one: "\u0623\u0642\u0644 \u0645\u0646 \u062B\u0627\u0646\u064A\u0629",
    two: "\u0623\u0642\u0644 \u0645\u0646 \u062B\u0627\u0646\u064A\u062A\u064A\u0646",
    threeToTen: "\u0623\u0642\u0644 \u0645\u0646 {{count}} \u062B\u0648\u0627\u0646\u064A",
    other: "\u0623\u0642\u0644 \u0645\u0646 {{count}} \u062B\u0627\u0646\u064A\u0629"
  },
  xSeconds: {
    one: "\u062B\u0627\u0646\u064A\u0629 \u0648\u0627\u062D\u062F\u0629",
    two: "\u062B\u0627\u0646\u064A\u062A\u0627\u0646",
    threeToTen: "{{count}} \u062B\u0648\u0627\u0646\u064A",
    other: "{{count}} \u062B\u0627\u0646\u064A\u0629"
  },
  halfAMinute: "\u0646\u0635\u0641 \u062F\u0642\u064A\u0642\u0629",
  lessThanXMinutes: {
    one: "\u0623\u0642\u0644 \u0645\u0646 \u062F\u0642\u064A\u0642\u0629",
    two: "\u0623\u0642\u0644 \u0645\u0646 \u062F\u0642\u064A\u0642\u062A\u064A\u0646",
    threeToTen: "\u0623\u0642\u0644 \u0645\u0646 {{count}} \u062F\u0642\u0627\u0626\u0642",
    other: "\u0623\u0642\u0644 \u0645\u0646 {{count}} \u062F\u0642\u064A\u0642\u0629"
  },
  xMinutes: {
    one: "\u062F\u0642\u064A\u0642\u0629 \u0648\u0627\u062D\u062F\u0629",
    two: "\u062F\u0642\u064A\u0642\u062A\u0627\u0646",
    threeToTen: "{{count}} \u062F\u0642\u0627\u0626\u0642",
    other: "{{count}} \u062F\u0642\u064A\u0642\u0629"
  },
  aboutXHours: {
    one: "\u0633\u0627\u0639\u0629 \u0648\u0627\u062D\u062F\u0629 \u062A\u0642\u0631\u064A\u0628\u0627\u064B",
    two: "\u0633\u0627\u0639\u062A\u064A\u0646 \u062A\u0642\u0631\u064A\u0628\u0627",
    threeToTen: "{{count}} \u0633\u0627\u0639\u0627\u062A \u062A\u0642\u0631\u064A\u0628\u0627\u064B",
    other: "{{count}} \u0633\u0627\u0639\u0629 \u062A\u0642\u0631\u064A\u0628\u0627\u064B"
  },
  xHours: {
    one: "\u0633\u0627\u0639\u0629 \u0648\u0627\u062D\u062F\u0629",
    two: "\u0633\u0627\u0639\u062A\u0627\u0646",
    threeToTen: "{{count}} \u0633\u0627\u0639\u0627\u062A",
    other: "{{count}} \u0633\u0627\u0639\u0629"
  },
  xDays: {
    one: "\u064A\u0648\u0645 \u0648\u0627\u062D\u062F",
    two: "\u064A\u0648\u0645\u0627\u0646",
    threeToTen: "{{count}} \u0623\u064A\u0627\u0645",
    other: "{{count}} \u064A\u0648\u0645"
  },
  aboutXWeeks: {
    one: "\u0623\u0633\u0628\u0648\u0639 \u0648\u0627\u062D\u062F \u062A\u0642\u0631\u064A\u0628\u0627",
    two: "\u0623\u0633\u0628\u0648\u0639\u064A\u0646 \u062A\u0642\u0631\u064A\u0628\u0627",
    threeToTen: "{{count}} \u0623\u0633\u0627\u0628\u064A\u0639 \u062A\u0642\u0631\u064A\u0628\u0627",
    other: "{{count}} \u0623\u0633\u0628\u0648\u0639\u0627 \u062A\u0642\u0631\u064A\u0628\u0627"
  },
  xWeeks: {
    one: "\u0623\u0633\u0628\u0648\u0639 \u0648\u0627\u062D\u062F",
    two: "\u0623\u0633\u0628\u0648\u0639\u0627\u0646",
    threeToTen: "{{count}} \u0623\u0633\u0627\u0628\u064A\u0639",
    other: "{{count}} \u0623\u0633\u0628\u0648\u0639\u0627"
  },
  aboutXMonths: {
    one: "\u0634\u0647\u0631 \u0648\u0627\u062D\u062F \u062A\u0642\u0631\u064A\u0628\u0627\u064B",
    two: "\u0634\u0647\u0631\u064A\u0646 \u062A\u0642\u0631\u064A\u0628\u0627",
    threeToTen: "{{count}} \u0623\u0634\u0647\u0631 \u062A\u0642\u0631\u064A\u0628\u0627",
    other: "{{count}} \u0634\u0647\u0631\u0627 \u062A\u0642\u0631\u064A\u0628\u0627\u064B"
  },
  xMonths: {
    one: "\u0634\u0647\u0631 \u0648\u0627\u062D\u062F",
    two: "\u0634\u0647\u0631\u0627\u0646",
    threeToTen: "{{count}} \u0623\u0634\u0647\u0631",
    other: "{{count}} \u0634\u0647\u0631\u0627"
  },
  aboutXYears: {
    one: "\u0633\u0646\u0629 \u0648\u0627\u062D\u062F\u0629 \u062A\u0642\u0631\u064A\u0628\u0627\u064B",
    two: "\u0633\u0646\u062A\u064A\u0646 \u062A\u0642\u0631\u064A\u0628\u0627",
    threeToTen: "{{count}} \u0633\u0646\u0648\u0627\u062A \u062A\u0642\u0631\u064A\u0628\u0627\u064B",
    other: "{{count}} \u0633\u0646\u0629 \u062A\u0642\u0631\u064A\u0628\u0627\u064B"
  },
  xYears: {
    one: "\u0633\u0646\u0629 \u0648\u0627\u062D\u062F",
    two: "\u0633\u0646\u062A\u0627\u0646",
    threeToTen: "{{count}} \u0633\u0646\u0648\u0627\u062A",
    other: "{{count}} \u0633\u0646\u0629"
  },
  overXYears: {
    one: "\u0623\u0643\u062B\u0631 \u0645\u0646 \u0633\u0646\u0629",
    two: "\u0623\u0643\u062B\u0631 \u0645\u0646 \u0633\u0646\u062A\u064A\u0646",
    threeToTen: "\u0623\u0643\u062B\u0631 \u0645\u0646 {{count}} \u0633\u0646\u0648\u0627\u062A",
    other: "\u0623\u0643\u062B\u0631 \u0645\u0646 {{count}} \u0633\u0646\u0629"
  },
  almostXYears: {
    one: "\u0645\u0627 \u064A\u0642\u0627\u0631\u0628 \u0633\u0646\u0629 \u0648\u0627\u062D\u062F\u0629",
    two: "\u0645\u0627 \u064A\u0642\u0627\u0631\u0628 \u0633\u0646\u062A\u064A\u0646",
    threeToTen: "\u0645\u0627 \u064A\u0642\u0627\u0631\u0628 {{count}} \u0633\u0646\u0648\u0627\u062A",
    other: "\u0645\u0627 \u064A\u0642\u0627\u0631\u0628 {{count}} \u0633\u0646\u0629"
  }
};
var formatDistance = (token, count, options) => {
  const usageGroup = formatDistanceLocale[token];
  let result;
  if (typeof usageGroup === "string") {
    result = usageGroup;
  } else if (count === 1) {
    result = usageGroup.one;
  } else if (count === 2) {
    result = usageGroup.two;
  } else if (count <= 10) {
    result = usageGroup.threeToTen.replace("{{count}}", String(count));
  } else {
    result = usageGroup.other.replace("{{count}}", String(count));
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return "\u062E\u0644\u0627\u0644 " + result;
    } else {
      return "\u0645\u0646\u0630 " + result;
    }
  }
  return result;
};

// node_modules/date-fns/locale/_lib/buildFormatLongFn.js
function buildFormatLongFn(args) {
  return (options = {}) => {
    const width = options.width ? String(options.width) : args.defaultWidth;
    const format = args.formats[width] || args.formats[args.defaultWidth];
    return format;
  };
}

// node_modules/date-fns/locale/ar/_lib/formatLong.js
var dateFormats = {
  full: "EEEE\u060C do MMMM y",
  long: "do MMMM y",
  medium: "d MMM y",
  short: "dd/MM/yyyy"
};
var timeFormats = {
  full: "HH:mm:ss",
  long: "HH:mm:ss",
  medium: "HH:mm:ss",
  short: "HH:mm"
};
var dateTimeFormats = {
  full: "{{date}} '\u0639\u0646\u062F \u0627\u0644\u0633\u0627\u0639\u0629' {{time}}",
  long: "{{date}} '\u0639\u0646\u062F \u0627\u0644\u0633\u0627\u0639\u0629' {{time}}",
  medium: "{{date}}, {{time}}",
  short: "{{date}}, {{time}}"
};
var formatLong = {
  date: buildFormatLongFn({
    formats: dateFormats,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/ar/_lib/formatRelative.js
var formatRelativeLocale = {
  lastWeek: "eeee '\u0627\u0644\u0645\u0627\u0636\u064A \u0639\u0646\u062F \u0627\u0644\u0633\u0627\u0639\u0629' p",
  yesterday: "'\u0627\u0644\u0623\u0645\u0633 \u0639\u0646\u062F \u0627\u0644\u0633\u0627\u0639\u0629' p",
  today: "'\u0627\u0644\u064A\u0648\u0645 \u0639\u0646\u062F \u0627\u0644\u0633\u0627\u0639\u0629' p",
  tomorrow: "'\u063A\u062F\u0627 \u0639\u0646\u062F \u0627\u0644\u0633\u0627\u0639\u0629' p",
  nextWeek: "eeee '\u0627\u0644\u0642\u0627\u062F\u0645 \u0639\u0646\u062F \u0627\u0644\u0633\u0627\u0639\u0629' p",
  other: "P"
};
var formatRelative = (token) => formatRelativeLocale[token];

// node_modules/date-fns/locale/_lib/buildLocalizeFn.js
function buildLocalizeFn(args) {
  return (value, options) => {
    const context = options?.context ? String(options.context) : "standalone";
    let valuesArray;
    if (context === "formatting" && args.formattingValues) {
      const defaultWidth = args.defaultFormattingWidth || args.defaultWidth;
      const width = options?.width ? String(options.width) : defaultWidth;
      valuesArray = args.formattingValues[width] || args.formattingValues[defaultWidth];
    } else {
      const defaultWidth = args.defaultWidth;
      const width = options?.width ? String(options.width) : args.defaultWidth;
      valuesArray = args.values[width] || args.values[defaultWidth];
    }
    const index = args.argumentCallback ? args.argumentCallback(value) : value;
    return valuesArray[index];
  };
}

// node_modules/date-fns/locale/ar/_lib/localize.js
var eraValues = {
  narrow: ["\u0642", "\u0628"],
  abbreviated: ["\u0642.\u0645.", "\u0628.\u0645."],
  wide: ["\u0642\u0628\u0644 \u0627\u0644\u0645\u064A\u0644\u0627\u062F", "\u0628\u0639\u062F \u0627\u0644\u0645\u064A\u0644\u0627\u062F"]
};
var quarterValues = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["\u06311", "\u06312", "\u06313", "\u06314"],
  wide: ["\u0627\u0644\u0631\u0628\u0639 \u0627\u0644\u0623\u0648\u0644", "\u0627\u0644\u0631\u0628\u0639 \u0627\u0644\u062B\u0627\u0646\u064A", "\u0627\u0644\u0631\u0628\u0639 \u0627\u0644\u062B\u0627\u0644\u062B", "\u0627\u0644\u0631\u0628\u0639 \u0627\u0644\u0631\u0627\u0628\u0639"]
};
var monthValues = {
  narrow: ["\u064A", "\u0641", "\u0645", "\u0623", "\u0645", "\u064A", "\u064A", "\u0623", "\u0633", "\u0623", "\u0646", "\u062F"],
  abbreviated: ["\u064A\u0646\u0627\u064A\u0631", "\u0641\u0628\u0631\u0627\u064A\u0631", "\u0645\u0627\u0631\u0633", "\u0623\u0628\u0631\u064A\u0644", "\u0645\u0627\u064A\u0648", "\u064A\u0648\u0646\u064A\u0648", "\u064A\u0648\u0644\u064A\u0648", "\u0623\u063A\u0633\u0637\u0633", "\u0633\u0628\u062A\u0645\u0628\u0631", "\u0623\u0643\u062A\u0648\u0628\u0631", "\u0646\u0648\u0641\u0645\u0628\u0631", "\u062F\u064A\u0633\u0645\u0628\u0631"],
  wide: ["\u064A\u0646\u0627\u064A\u0631", "\u0641\u0628\u0631\u0627\u064A\u0631", "\u0645\u0627\u0631\u0633", "\u0623\u0628\u0631\u064A\u0644", "\u0645\u0627\u064A\u0648", "\u064A\u0648\u0646\u064A\u0648", "\u064A\u0648\u0644\u064A\u0648", "\u0623\u063A\u0633\u0637\u0633", "\u0633\u0628\u062A\u0645\u0628\u0631", "\u0623\u0643\u062A\u0648\u0628\u0631", "\u0646\u0648\u0641\u0645\u0628\u0631", "\u062F\u064A\u0633\u0645\u0628\u0631"]
};
var dayValues = {
  narrow: ["\u062D", "\u0646", "\u062B", "\u0631", "\u062E", "\u062C", "\u0633"],
  short: ["\u0623\u062D\u062F", "\u0627\u062B\u0646\u064A\u0646", "\u062B\u0644\u0627\u062B\u0627\u0621", "\u0623\u0631\u0628\u0639\u0627\u0621", "\u062E\u0645\u064A\u0633", "\u062C\u0645\u0639\u0629", "\u0633\u0628\u062A"],
  abbreviated: ["\u0623\u062D\u062F", "\u0627\u062B\u0646\u064A\u0646", "\u062B\u0644\u0627\u062B\u0627\u0621", "\u0623\u0631\u0628\u0639\u0627\u0621", "\u062E\u0645\u064A\u0633", "\u062C\u0645\u0639\u0629", "\u0633\u0628\u062A"],
  wide: ["\u0627\u0644\u0623\u062D\u062F", "\u0627\u0644\u0627\u062B\u0646\u064A\u0646", "\u0627\u0644\u062B\u0644\u0627\u062B\u0627\u0621", "\u0627\u0644\u0623\u0631\u0628\u0639\u0627\u0621", "\u0627\u0644\u062E\u0645\u064A\u0633", "\u0627\u0644\u062C\u0645\u0639\u0629", "\u0627\u0644\u0633\u0628\u062A"]
};
var dayPeriodValues = {
  narrow: {
    am: "\u0635",
    pm: "\u0645",
    morning: "\u0627\u0644\u0635\u0628\u0627\u062D",
    noon: "\u0627\u0644\u0638\u0647\u0631",
    afternoon: "\u0628\u0639\u062F \u0627\u0644\u0638\u0647\u0631",
    evening: "\u0627\u0644\u0645\u0633\u0627\u0621",
    night: "\u0627\u0644\u0644\u064A\u0644",
    midnight: "\u0645\u0646\u062A\u0635\u0641 \u0627\u0644\u0644\u064A\u0644"
  },
  abbreviated: {
    am: "\u0635",
    pm: "\u0645",
    morning: "\u0627\u0644\u0635\u0628\u0627\u062D",
    noon: "\u0627\u0644\u0638\u0647\u0631",
    afternoon: "\u0628\u0639\u062F \u0627\u0644\u0638\u0647\u0631",
    evening: "\u0627\u0644\u0645\u0633\u0627\u0621",
    night: "\u0627\u0644\u0644\u064A\u0644",
    midnight: "\u0645\u0646\u062A\u0635\u0641 \u0627\u0644\u0644\u064A\u0644"
  },
  wide: {
    am: "\u0635",
    pm: "\u0645",
    morning: "\u0627\u0644\u0635\u0628\u0627\u062D",
    noon: "\u0627\u0644\u0638\u0647\u0631",
    afternoon: "\u0628\u0639\u062F \u0627\u0644\u0638\u0647\u0631",
    evening: "\u0627\u0644\u0645\u0633\u0627\u0621",
    night: "\u0627\u0644\u0644\u064A\u0644",
    midnight: "\u0645\u0646\u062A\u0635\u0641 \u0627\u0644\u0644\u064A\u0644"
  }
};
var formattingDayPeriodValues = {
  narrow: {
    am: "\u0635",
    pm: "\u0645",
    morning: "\u0641\u064A \u0627\u0644\u0635\u0628\u0627\u062D",
    noon: "\u0627\u0644\u0638\u0647\u0631",
    afternoon: "\u0628\u0639\u062F \u0627\u0644\u0638\u0647\u0631",
    evening: "\u0641\u064A \u0627\u0644\u0645\u0633\u0627\u0621",
    night: "\u0641\u064A \u0627\u0644\u0644\u064A\u0644",
    midnight: "\u0645\u0646\u062A\u0635\u0641 \u0627\u0644\u0644\u064A\u0644"
  },
  abbreviated: {
    am: "\u0635",
    pm: "\u0645",
    morning: "\u0641\u064A \u0627\u0644\u0635\u0628\u0627\u062D",
    noon: "\u0627\u0644\u0638\u0647\u0631",
    afternoon: "\u0628\u0639\u062F \u0627\u0644\u0638\u0647\u0631",
    evening: "\u0641\u064A \u0627\u0644\u0645\u0633\u0627\u0621",
    night: "\u0641\u064A \u0627\u0644\u0644\u064A\u0644",
    midnight: "\u0645\u0646\u062A\u0635\u0641 \u0627\u0644\u0644\u064A\u0644"
  },
  wide: {
    am: "\u0635",
    pm: "\u0645",
    morning: "\u0641\u064A \u0627\u0644\u0635\u0628\u0627\u062D",
    noon: "\u0627\u0644\u0638\u0647\u0631",
    afternoon: "\u0628\u0639\u062F \u0627\u0644\u0638\u0647\u0631",
    evening: "\u0641\u064A \u0627\u0644\u0645\u0633\u0627\u0621",
    night: "\u0641\u064A \u0627\u0644\u0644\u064A\u0644",
    midnight: "\u0645\u0646\u062A\u0635\u0641 \u0627\u0644\u0644\u064A\u0644"
  }
};
var ordinalNumber = (num) => String(num);
var localize = {
  ordinalNumber,
  era: buildLocalizeFn({
    values: eraValues,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues,
    defaultWidth: "wide",
    formattingValues: formattingDayPeriodValues,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/_lib/buildMatchPatternFn.js
function buildMatchPatternFn(args) {
  return (string, options = {}) => {
    const matchResult = string.match(args.matchPattern);
    if (!matchResult) return null;
    const matchedString = matchResult[0];
    const parseResult = string.match(args.parsePattern);
    if (!parseResult) return null;
    let value = args.valueCallback ? args.valueCallback(parseResult[0]) : parseResult[0];
    value = options.valueCallback ? options.valueCallback(value) : value;
    const rest = string.slice(matchedString.length);
    return {
      value,
      rest
    };
  };
}

// node_modules/date-fns/locale/_lib/buildMatchFn.js
function buildMatchFn(args) {
  return (string, options = {}) => {
    const width = options.width;
    const matchPattern = width && args.matchPatterns[width] || args.matchPatterns[args.defaultMatchWidth];
    const matchResult = string.match(matchPattern);
    if (!matchResult) {
      return null;
    }
    const matchedString = matchResult[0];
    const parsePatterns = width && args.parsePatterns[width] || args.parsePatterns[args.defaultParseWidth];
    const key = Array.isArray(parsePatterns) ? findIndex(parsePatterns, (pattern) => pattern.test(matchedString)) : (
      // [TODO] -- I challenge you to fix the type
      findKey(parsePatterns, (pattern) => pattern.test(matchedString))
    );
    let value;
    value = args.valueCallback ? args.valueCallback(key) : key;
    value = options.valueCallback ? (
      // [TODO] -- I challenge you to fix the type
      options.valueCallback(value)
    ) : value;
    const rest = string.slice(matchedString.length);
    return {
      value,
      rest
    };
  };
}
function findKey(object, predicate) {
  for (const key in object) {
    if (Object.prototype.hasOwnProperty.call(object, key) && predicate(object[key])) {
      return key;
    }
  }
  return void 0;
}
function findIndex(array, predicate) {
  for (let key = 0; key < array.length; key++) {
    if (predicate(array[key])) {
      return key;
    }
  }
  return void 0;
}

// node_modules/date-fns/locale/ar/_lib/match.js
var matchOrdinalNumberPattern = /^(\d+)(th|st|nd|rd)?/i;
var parseOrdinalNumberPattern = /\d+/i;
var matchEraPatterns = {
  narrow: /[قب]/,
  abbreviated: /[قب]\.م\./,
  wide: /(قبل|بعد) الميلاد/
};
var parseEraPatterns = {
  any: [/قبل/, /بعد/]
};
var matchQuarterPatterns = {
  narrow: /^[1234]/i,
  abbreviated: /ر[1234]/,
  wide: /الربع (الأول|الثاني|الثالث|الرابع)/
};
var parseQuarterPatterns = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns = {
  narrow: /^[أيفمسند]/,
  abbreviated: /^(يناير|فبراير|مارس|أبريل|مايو|يونيو|يوليو|أغسطس|سبتمبر|أكتوبر|نوفمبر|ديسمبر)/,
  wide: /^(يناير|فبراير|مارس|أبريل|مايو|يونيو|يوليو|أغسطس|سبتمبر|أكتوبر|نوفمبر|ديسمبر)/
};
var parseMonthPatterns = {
  narrow: [/^ي/i, /^ف/i, /^م/i, /^أ/i, /^م/i, /^ي/i, /^ي/i, /^أ/i, /^س/i, /^أ/i, /^ن/i, /^د/i],
  any: [/^يناير/i, /^فبراير/i, /^مارس/i, /^أبريل/i, /^مايو/i, /^يونيو/i, /^يوليو/i, /^أغسطس/i, /^سبتمبر/i, /^أكتوبر/i, /^نوفمبر/i, /^ديسمبر/i]
};
var matchDayPatterns = {
  narrow: /^[حنثرخجس]/i,
  short: /^(أحد|اثنين|ثلاثاء|أربعاء|خميس|جمعة|سبت)/i,
  abbreviated: /^(أحد|اثنين|ثلاثاء|أربعاء|خميس|جمعة|سبت)/i,
  wide: /^(الأحد|الاثنين|الثلاثاء|الأربعاء|الخميس|الجمعة|السبت)/i
};
var parseDayPatterns = {
  narrow: [/^ح/i, /^ن/i, /^ث/i, /^ر/i, /^خ/i, /^ج/i, /^س/i],
  wide: [/^الأحد/i, /^الاثنين/i, /^الثلاثاء/i, /^الأربعاء/i, /^الخميس/i, /^الجمعة/i, /^السبت/i],
  any: [/^أح/i, /^اث/i, /^ث/i, /^أر/i, /^خ/i, /^ج/i, /^س/i]
};
var matchDayPeriodPatterns = {
  narrow: /^(ص|م|منتصف الليل|الظهر|بعد الظهر|في الصباح|في المساء|في الليل)/,
  any: /^(ص|م|منتصف الليل|الظهر|بعد الظهر|في الصباح|في المساء|في الليل)/
};
var parseDayPeriodPatterns = {
  any: {
    am: /^ص/,
    pm: /^م/,
    midnight: /منتصف الليل/,
    noon: /الظهر/,
    afternoon: /بعد الظهر/,
    morning: /في الصباح/,
    evening: /في المساء/,
    night: /في الليل/
  }
};
var match = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern,
    parsePattern: parseOrdinalNumberPattern,
    valueCallback: (value) => parseInt(value, 10)
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/ar.js
var ar = {
  code: "ar",
  formatDistance,
  formatLong,
  formatRelative,
  localize,
  match,
  options: {
    weekStartsOn: 6,
    firstWeekContainsDate: 1
  }
};

// node_modules/date-fns/locale/de/_lib/formatDistance.js
var formatDistanceLocale2 = {
  lessThanXSeconds: {
    standalone: {
      one: "weniger als 1 Sekunde",
      other: "weniger als {{count}} Sekunden"
    },
    withPreposition: {
      one: "weniger als 1 Sekunde",
      other: "weniger als {{count}} Sekunden"
    }
  },
  xSeconds: {
    standalone: {
      one: "1 Sekunde",
      other: "{{count}} Sekunden"
    },
    withPreposition: {
      one: "1 Sekunde",
      other: "{{count}} Sekunden"
    }
  },
  halfAMinute: {
    standalone: "eine halbe Minute",
    withPreposition: "einer halben Minute"
  },
  lessThanXMinutes: {
    standalone: {
      one: "weniger als 1 Minute",
      other: "weniger als {{count}} Minuten"
    },
    withPreposition: {
      one: "weniger als 1 Minute",
      other: "weniger als {{count}} Minuten"
    }
  },
  xMinutes: {
    standalone: {
      one: "1 Minute",
      other: "{{count}} Minuten"
    },
    withPreposition: {
      one: "1 Minute",
      other: "{{count}} Minuten"
    }
  },
  aboutXHours: {
    standalone: {
      one: "etwa 1 Stunde",
      other: "etwa {{count}} Stunden"
    },
    withPreposition: {
      one: "etwa 1 Stunde",
      other: "etwa {{count}} Stunden"
    }
  },
  xHours: {
    standalone: {
      one: "1 Stunde",
      other: "{{count}} Stunden"
    },
    withPreposition: {
      one: "1 Stunde",
      other: "{{count}} Stunden"
    }
  },
  xDays: {
    standalone: {
      one: "1 Tag",
      other: "{{count}} Tage"
    },
    withPreposition: {
      one: "1 Tag",
      other: "{{count}} Tagen"
    }
  },
  aboutXWeeks: {
    standalone: {
      one: "etwa 1 Woche",
      other: "etwa {{count}} Wochen"
    },
    withPreposition: {
      one: "etwa 1 Woche",
      other: "etwa {{count}} Wochen"
    }
  },
  xWeeks: {
    standalone: {
      one: "1 Woche",
      other: "{{count}} Wochen"
    },
    withPreposition: {
      one: "1 Woche",
      other: "{{count}} Wochen"
    }
  },
  aboutXMonths: {
    standalone: {
      one: "etwa 1 Monat",
      other: "etwa {{count}} Monate"
    },
    withPreposition: {
      one: "etwa 1 Monat",
      other: "etwa {{count}} Monaten"
    }
  },
  xMonths: {
    standalone: {
      one: "1 Monat",
      other: "{{count}} Monate"
    },
    withPreposition: {
      one: "1 Monat",
      other: "{{count}} Monaten"
    }
  },
  aboutXYears: {
    standalone: {
      one: "etwa 1 Jahr",
      other: "etwa {{count}} Jahre"
    },
    withPreposition: {
      one: "etwa 1 Jahr",
      other: "etwa {{count}} Jahren"
    }
  },
  xYears: {
    standalone: {
      one: "1 Jahr",
      other: "{{count}} Jahre"
    },
    withPreposition: {
      one: "1 Jahr",
      other: "{{count}} Jahren"
    }
  },
  overXYears: {
    standalone: {
      one: "mehr als 1 Jahr",
      other: "mehr als {{count}} Jahre"
    },
    withPreposition: {
      one: "mehr als 1 Jahr",
      other: "mehr als {{count}} Jahren"
    }
  },
  almostXYears: {
    standalone: {
      one: "fast 1 Jahr",
      other: "fast {{count}} Jahre"
    },
    withPreposition: {
      one: "fast 1 Jahr",
      other: "fast {{count}} Jahren"
    }
  }
};
var formatDistance2 = (token, count, options) => {
  let result;
  const tokenValue = options?.addSuffix ? formatDistanceLocale2[token].withPreposition : formatDistanceLocale2[token].standalone;
  if (typeof tokenValue === "string") {
    result = tokenValue;
  } else if (count === 1) {
    result = tokenValue.one;
  } else {
    result = tokenValue.other.replace("{{count}}", String(count));
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return "in " + result;
    } else {
      return "vor " + result;
    }
  }
  return result;
};

// node_modules/date-fns/locale/de/_lib/formatLong.js
var dateFormats2 = {
  full: "EEEE, do MMMM y",
  // Montag, 7. Januar 2018
  long: "do MMMM y",
  // 7. Januar 2018
  medium: "do MMM y",
  // 7. Jan. 2018
  short: "dd.MM.y"
  // 07.01.2018
};
var timeFormats2 = {
  full: "HH:mm:ss zzzz",
  long: "HH:mm:ss z",
  medium: "HH:mm:ss",
  short: "HH:mm"
};
var dateTimeFormats2 = {
  full: "{{date}} 'um' {{time}}",
  long: "{{date}} 'um' {{time}}",
  medium: "{{date}} {{time}}",
  short: "{{date}} {{time}}"
};
var formatLong2 = {
  date: buildFormatLongFn({
    formats: dateFormats2,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats2,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats2,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/de/_lib/formatRelative.js
var formatRelativeLocale2 = {
  lastWeek: "'letzten' eeee 'um' p",
  yesterday: "'gestern um' p",
  today: "'heute um' p",
  tomorrow: "'morgen um' p",
  nextWeek: "eeee 'um' p",
  other: "P"
};
var formatRelative2 = (token, _date, _baseDate, _options) => formatRelativeLocale2[token];

// node_modules/date-fns/locale/de/_lib/localize.js
var eraValues2 = {
  narrow: ["v.Chr.", "n.Chr."],
  abbreviated: ["v.Chr.", "n.Chr."],
  wide: ["vor Christus", "nach Christus"]
};
var quarterValues2 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["Q1", "Q2", "Q3", "Q4"],
  wide: ["1. Quartal", "2. Quartal", "3. Quartal", "4. Quartal"]
};
var monthValues2 = {
  narrow: ["J", "F", "M", "A", "M", "J", "J", "A", "S", "O", "N", "D"],
  abbreviated: ["Jan", "Feb", "M\xE4r", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Dez"],
  wide: ["Januar", "Februar", "M\xE4rz", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"]
};
var formattingMonthValues = {
  narrow: monthValues2.narrow,
  abbreviated: ["Jan.", "Feb.", "M\xE4rz", "Apr.", "Mai", "Juni", "Juli", "Aug.", "Sep.", "Okt.", "Nov.", "Dez."],
  wide: monthValues2.wide
};
var dayValues2 = {
  narrow: ["S", "M", "D", "M", "D", "F", "S"],
  short: ["So", "Mo", "Di", "Mi", "Do", "Fr", "Sa"],
  abbreviated: ["So.", "Mo.", "Di.", "Mi.", "Do.", "Fr.", "Sa."],
  wide: ["Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"]
};
var dayPeriodValues2 = {
  narrow: {
    am: "vm.",
    pm: "nm.",
    midnight: "Mitternacht",
    noon: "Mittag",
    morning: "Morgen",
    afternoon: "Nachm.",
    evening: "Abend",
    night: "Nacht"
  },
  abbreviated: {
    am: "vorm.",
    pm: "nachm.",
    midnight: "Mitternacht",
    noon: "Mittag",
    morning: "Morgen",
    afternoon: "Nachmittag",
    evening: "Abend",
    night: "Nacht"
  },
  wide: {
    am: "vormittags",
    pm: "nachmittags",
    midnight: "Mitternacht",
    noon: "Mittag",
    morning: "Morgen",
    afternoon: "Nachmittag",
    evening: "Abend",
    night: "Nacht"
  }
};
var formattingDayPeriodValues2 = {
  narrow: {
    am: "vm.",
    pm: "nm.",
    midnight: "Mitternacht",
    noon: "Mittag",
    morning: "morgens",
    afternoon: "nachm.",
    evening: "abends",
    night: "nachts"
  },
  abbreviated: {
    am: "vorm.",
    pm: "nachm.",
    midnight: "Mitternacht",
    noon: "Mittag",
    morning: "morgens",
    afternoon: "nachmittags",
    evening: "abends",
    night: "nachts"
  },
  wide: {
    am: "vormittags",
    pm: "nachmittags",
    midnight: "Mitternacht",
    noon: "Mittag",
    morning: "morgens",
    afternoon: "nachmittags",
    evening: "abends",
    night: "nachts"
  }
};
var ordinalNumber2 = (dirtyNumber) => {
  const number = Number(dirtyNumber);
  return number + ".";
};
var localize2 = {
  ordinalNumber: ordinalNumber2,
  era: buildLocalizeFn({
    values: eraValues2,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues2,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues2,
    formattingValues: formattingMonthValues,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues2,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues2,
    defaultWidth: "wide",
    formattingValues: formattingDayPeriodValues2,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/de/_lib/match.js
var matchOrdinalNumberPattern2 = /^(\d+)(\.)?/i;
var parseOrdinalNumberPattern2 = /\d+/i;
var matchEraPatterns2 = {
  narrow: /^(v\.? ?Chr\.?|n\.? ?Chr\.?)/i,
  abbreviated: /^(v\.? ?Chr\.?|n\.? ?Chr\.?)/i,
  wide: /^(vor Christus|vor unserer Zeitrechnung|nach Christus|unserer Zeitrechnung)/i
};
var parseEraPatterns2 = {
  any: [/^v/i, /^n/i]
};
var matchQuarterPatterns2 = {
  narrow: /^[1234]/i,
  abbreviated: /^q[1234]/i,
  wide: /^[1234](\.)? Quartal/i
};
var parseQuarterPatterns2 = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns2 = {
  narrow: /^[jfmasond]/i,
  abbreviated: /^(j[aä]n|feb|mär[z]?|apr|mai|jun[i]?|jul[i]?|aug|sep|okt|nov|dez)\.?/i,
  wide: /^(januar|februar|märz|april|mai|juni|juli|august|september|oktober|november|dezember)/i
};
var parseMonthPatterns2 = {
  narrow: [/^j/i, /^f/i, /^m/i, /^a/i, /^m/i, /^j/i, /^j/i, /^a/i, /^s/i, /^o/i, /^n/i, /^d/i],
  any: [/^j[aä]/i, /^f/i, /^mär/i, /^ap/i, /^mai/i, /^jun/i, /^jul/i, /^au/i, /^s/i, /^o/i, /^n/i, /^d/i]
};
var matchDayPatterns2 = {
  narrow: /^[smdmf]/i,
  short: /^(so|mo|di|mi|do|fr|sa)/i,
  abbreviated: /^(son?|mon?|die?|mit?|don?|fre?|sam?)\.?/i,
  wide: /^(sonntag|montag|dienstag|mittwoch|donnerstag|freitag|samstag)/i
};
var parseDayPatterns2 = {
  any: [/^so/i, /^mo/i, /^di/i, /^mi/i, /^do/i, /^f/i, /^sa/i]
};
var matchDayPeriodPatterns2 = {
  narrow: /^(vm\.?|nm\.?|Mitternacht|Mittag|morgens|nachm\.?|abends|nachts)/i,
  abbreviated: /^(vorm\.?|nachm\.?|Mitternacht|Mittag|morgens|nachm\.?|abends|nachts)/i,
  wide: /^(vormittags|nachmittags|Mitternacht|Mittag|morgens|nachmittags|abends|nachts)/i
};
var parseDayPeriodPatterns2 = {
  any: {
    am: /^v/i,
    pm: /^n/i,
    midnight: /^Mitte/i,
    noon: /^Mitta/i,
    morning: /morgens/i,
    afternoon: /nachmittags/i,
    // will never be matched. Afternoon is matched by `pm`
    evening: /abends/i,
    night: /nachts/i
    // will never be matched. Night is matched by `pm`
  }
};
var match2 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern2,
    parsePattern: parseOrdinalNumberPattern2,
    valueCallback: (value) => parseInt(value)
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns2,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns2,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns2,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns2,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns2,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns2,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns2,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns2,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns2,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPeriodPatterns2,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/de.js
var de = {
  code: "de",
  formatDistance: formatDistance2,
  formatLong: formatLong2,
  formatRelative: formatRelative2,
  localize: localize2,
  match: match2,
  options: {
    weekStartsOn: 1,
    firstWeekContainsDate: 4
  }
};

// node_modules/date-fns/locale/en-US/_lib/formatDistance.js
var formatDistanceLocale3 = {
  lessThanXSeconds: {
    one: "less than a second",
    other: "less than {{count}} seconds"
  },
  xSeconds: {
    one: "1 second",
    other: "{{count}} seconds"
  },
  halfAMinute: "half a minute",
  lessThanXMinutes: {
    one: "less than a minute",
    other: "less than {{count}} minutes"
  },
  xMinutes: {
    one: "1 minute",
    other: "{{count}} minutes"
  },
  aboutXHours: {
    one: "about 1 hour",
    other: "about {{count}} hours"
  },
  xHours: {
    one: "1 hour",
    other: "{{count}} hours"
  },
  xDays: {
    one: "1 day",
    other: "{{count}} days"
  },
  aboutXWeeks: {
    one: "about 1 week",
    other: "about {{count}} weeks"
  },
  xWeeks: {
    one: "1 week",
    other: "{{count}} weeks"
  },
  aboutXMonths: {
    one: "about 1 month",
    other: "about {{count}} months"
  },
  xMonths: {
    one: "1 month",
    other: "{{count}} months"
  },
  aboutXYears: {
    one: "about 1 year",
    other: "about {{count}} years"
  },
  xYears: {
    one: "1 year",
    other: "{{count}} years"
  },
  overXYears: {
    one: "over 1 year",
    other: "over {{count}} years"
  },
  almostXYears: {
    one: "almost 1 year",
    other: "almost {{count}} years"
  }
};
var formatDistance3 = (token, count, options) => {
  let result;
  const tokenValue = formatDistanceLocale3[token];
  if (typeof tokenValue === "string") {
    result = tokenValue;
  } else if (count === 1) {
    result = tokenValue.one;
  } else {
    result = tokenValue.other.replace("{{count}}", count.toString());
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return "in " + result;
    } else {
      return result + " ago";
    }
  }
  return result;
};

// node_modules/date-fns/locale/en-US/_lib/formatLong.js
var dateFormats3 = {
  full: "EEEE, MMMM do, y",
  long: "MMMM do, y",
  medium: "MMM d, y",
  short: "MM/dd/yyyy"
};
var timeFormats3 = {
  full: "h:mm:ss a zzzz",
  long: "h:mm:ss a z",
  medium: "h:mm:ss a",
  short: "h:mm a"
};
var dateTimeFormats3 = {
  full: "{{date}} 'at' {{time}}",
  long: "{{date}} 'at' {{time}}",
  medium: "{{date}}, {{time}}",
  short: "{{date}}, {{time}}"
};
var formatLong3 = {
  date: buildFormatLongFn({
    formats: dateFormats3,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats3,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats3,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/en-US/_lib/formatRelative.js
var formatRelativeLocale3 = {
  lastWeek: "'last' eeee 'at' p",
  yesterday: "'yesterday at' p",
  today: "'today at' p",
  tomorrow: "'tomorrow at' p",
  nextWeek: "eeee 'at' p",
  other: "P"
};
var formatRelative3 = (token, _date, _baseDate, _options) => formatRelativeLocale3[token];

// node_modules/date-fns/locale/en-US/_lib/localize.js
var eraValues3 = {
  narrow: ["B", "A"],
  abbreviated: ["BC", "AD"],
  wide: ["Before Christ", "Anno Domini"]
};
var quarterValues3 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["Q1", "Q2", "Q3", "Q4"],
  wide: ["1st quarter", "2nd quarter", "3rd quarter", "4th quarter"]
};
var monthValues3 = {
  narrow: ["J", "F", "M", "A", "M", "J", "J", "A", "S", "O", "N", "D"],
  abbreviated: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
  wide: ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
};
var dayValues3 = {
  narrow: ["S", "M", "T", "W", "T", "F", "S"],
  short: ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"],
  abbreviated: ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
  wide: ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
};
var dayPeriodValues3 = {
  narrow: {
    am: "a",
    pm: "p",
    midnight: "mi",
    noon: "n",
    morning: "morning",
    afternoon: "afternoon",
    evening: "evening",
    night: "night"
  },
  abbreviated: {
    am: "AM",
    pm: "PM",
    midnight: "midnight",
    noon: "noon",
    morning: "morning",
    afternoon: "afternoon",
    evening: "evening",
    night: "night"
  },
  wide: {
    am: "a.m.",
    pm: "p.m.",
    midnight: "midnight",
    noon: "noon",
    morning: "morning",
    afternoon: "afternoon",
    evening: "evening",
    night: "night"
  }
};
var formattingDayPeriodValues3 = {
  narrow: {
    am: "a",
    pm: "p",
    midnight: "mi",
    noon: "n",
    morning: "in the morning",
    afternoon: "in the afternoon",
    evening: "in the evening",
    night: "at night"
  },
  abbreviated: {
    am: "AM",
    pm: "PM",
    midnight: "midnight",
    noon: "noon",
    morning: "in the morning",
    afternoon: "in the afternoon",
    evening: "in the evening",
    night: "at night"
  },
  wide: {
    am: "a.m.",
    pm: "p.m.",
    midnight: "midnight",
    noon: "noon",
    morning: "in the morning",
    afternoon: "in the afternoon",
    evening: "in the evening",
    night: "at night"
  }
};
var ordinalNumber3 = (dirtyNumber, _options) => {
  const number = Number(dirtyNumber);
  const rem100 = number % 100;
  if (rem100 > 20 || rem100 < 10) {
    switch (rem100 % 10) {
      case 1:
        return number + "st";
      case 2:
        return number + "nd";
      case 3:
        return number + "rd";
    }
  }
  return number + "th";
};
var localize3 = {
  ordinalNumber: ordinalNumber3,
  era: buildLocalizeFn({
    values: eraValues3,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues3,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues3,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues3,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues3,
    defaultWidth: "wide",
    formattingValues: formattingDayPeriodValues3,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/en-US/_lib/match.js
var matchOrdinalNumberPattern3 = /^(\d+)(th|st|nd|rd)?/i;
var parseOrdinalNumberPattern3 = /\d+/i;
var matchEraPatterns3 = {
  narrow: /^(b|a)/i,
  abbreviated: /^(b\.?\s?c\.?|b\.?\s?c\.?\s?e\.?|a\.?\s?d\.?|c\.?\s?e\.?)/i,
  wide: /^(before christ|before common era|anno domini|common era)/i
};
var parseEraPatterns3 = {
  any: [/^b/i, /^(a|c)/i]
};
var matchQuarterPatterns3 = {
  narrow: /^[1234]/i,
  abbreviated: /^q[1234]/i,
  wide: /^[1234](th|st|nd|rd)? quarter/i
};
var parseQuarterPatterns3 = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns3 = {
  narrow: /^[jfmasond]/i,
  abbreviated: /^(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)/i,
  wide: /^(january|february|march|april|may|june|july|august|september|october|november|december)/i
};
var parseMonthPatterns3 = {
  narrow: [/^j/i, /^f/i, /^m/i, /^a/i, /^m/i, /^j/i, /^j/i, /^a/i, /^s/i, /^o/i, /^n/i, /^d/i],
  any: [/^ja/i, /^f/i, /^mar/i, /^ap/i, /^may/i, /^jun/i, /^jul/i, /^au/i, /^s/i, /^o/i, /^n/i, /^d/i]
};
var matchDayPatterns3 = {
  narrow: /^[smtwf]/i,
  short: /^(su|mo|tu|we|th|fr|sa)/i,
  abbreviated: /^(sun|mon|tue|wed|thu|fri|sat)/i,
  wide: /^(sunday|monday|tuesday|wednesday|thursday|friday|saturday)/i
};
var parseDayPatterns3 = {
  narrow: [/^s/i, /^m/i, /^t/i, /^w/i, /^t/i, /^f/i, /^s/i],
  any: [/^su/i, /^m/i, /^tu/i, /^w/i, /^th/i, /^f/i, /^sa/i]
};
var matchDayPeriodPatterns3 = {
  narrow: /^(a|p|mi|n|(in the|at) (morning|afternoon|evening|night))/i,
  any: /^([ap]\.?\s?m\.?|midnight|noon|(in the|at) (morning|afternoon|evening|night))/i
};
var parseDayPeriodPatterns3 = {
  any: {
    am: /^a/i,
    pm: /^p/i,
    midnight: /^mi/i,
    noon: /^no/i,
    morning: /morning/i,
    afternoon: /afternoon/i,
    evening: /evening/i,
    night: /night/i
  }
};
var match3 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern3,
    parsePattern: parseOrdinalNumberPattern3,
    valueCallback: (value) => parseInt(value, 10)
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns3,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns3,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns3,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns3,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns3,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns3,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns3,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns3,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns3,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns3,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/en-US.js
var enUS = {
  code: "en-US",
  formatDistance: formatDistance3,
  formatLong: formatLong3,
  formatRelative: formatRelative3,
  localize: localize3,
  match: match3,
  options: {
    weekStartsOn: 0,
    firstWeekContainsDate: 1
  }
};

// node_modules/date-fns/locale/es/_lib/formatDistance.js
var formatDistanceLocale4 = {
  lessThanXSeconds: {
    one: "menos de un segundo",
    other: "menos de {{count}} segundos"
  },
  xSeconds: {
    one: "1 segundo",
    other: "{{count}} segundos"
  },
  halfAMinute: "medio minuto",
  lessThanXMinutes: {
    one: "menos de un minuto",
    other: "menos de {{count}} minutos"
  },
  xMinutes: {
    one: "1 minuto",
    other: "{{count}} minutos"
  },
  aboutXHours: {
    one: "alrededor de 1 hora",
    other: "alrededor de {{count}} horas"
  },
  xHours: {
    one: "1 hora",
    other: "{{count}} horas"
  },
  xDays: {
    one: "1 d\xEDa",
    other: "{{count}} d\xEDas"
  },
  aboutXWeeks: {
    one: "alrededor de 1 semana",
    other: "alrededor de {{count}} semanas"
  },
  xWeeks: {
    one: "1 semana",
    other: "{{count}} semanas"
  },
  aboutXMonths: {
    one: "alrededor de 1 mes",
    other: "alrededor de {{count}} meses"
  },
  xMonths: {
    one: "1 mes",
    other: "{{count}} meses"
  },
  aboutXYears: {
    one: "alrededor de 1 a\xF1o",
    other: "alrededor de {{count}} a\xF1os"
  },
  xYears: {
    one: "1 a\xF1o",
    other: "{{count}} a\xF1os"
  },
  overXYears: {
    one: "m\xE1s de 1 a\xF1o",
    other: "m\xE1s de {{count}} a\xF1os"
  },
  almostXYears: {
    one: "casi 1 a\xF1o",
    other: "casi {{count}} a\xF1os"
  }
};
var formatDistance4 = (token, count, options) => {
  let result;
  const tokenValue = formatDistanceLocale4[token];
  if (typeof tokenValue === "string") {
    result = tokenValue;
  } else if (count === 1) {
    result = tokenValue.one;
  } else {
    result = tokenValue.other.replace("{{count}}", count.toString());
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return "en " + result;
    } else {
      return "hace " + result;
    }
  }
  return result;
};

// node_modules/date-fns/locale/es/_lib/formatLong.js
var dateFormats4 = {
  full: "EEEE, d 'de' MMMM 'de' y",
  long: "d 'de' MMMM 'de' y",
  medium: "d MMM y",
  short: "dd/MM/y"
};
var timeFormats4 = {
  full: "HH:mm:ss zzzz",
  long: "HH:mm:ss z",
  medium: "HH:mm:ss",
  short: "HH:mm"
};
var dateTimeFormats4 = {
  full: "{{date}} 'a las' {{time}}",
  long: "{{date}} 'a las' {{time}}",
  medium: "{{date}}, {{time}}",
  short: "{{date}}, {{time}}"
};
var formatLong4 = {
  date: buildFormatLongFn({
    formats: dateFormats4,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats4,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats4,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/es/_lib/formatRelative.js
var formatRelativeLocale4 = {
  lastWeek: "'el' eeee 'pasado a la' p",
  yesterday: "'ayer a la' p",
  today: "'hoy a la' p",
  tomorrow: "'ma\xF1ana a la' p",
  nextWeek: "eeee 'a la' p",
  other: "P"
};
var formatRelativeLocalePlural = {
  lastWeek: "'el' eeee 'pasado a las' p",
  yesterday: "'ayer a las' p",
  today: "'hoy a las' p",
  tomorrow: "'ma\xF1ana a las' p",
  nextWeek: "eeee 'a las' p",
  other: "P"
};
var formatRelative4 = (token, date, _baseDate, _options) => {
  if (date.getHours() !== 1) {
    return formatRelativeLocalePlural[token];
  } else {
    return formatRelativeLocale4[token];
  }
};

// node_modules/date-fns/locale/es/_lib/localize.js
var eraValues4 = {
  narrow: ["AC", "DC"],
  abbreviated: ["AC", "DC"],
  wide: ["antes de cristo", "despu\xE9s de cristo"]
};
var quarterValues4 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["T1", "T2", "T3", "T4"],
  wide: ["1\xBA trimestre", "2\xBA trimestre", "3\xBA trimestre", "4\xBA trimestre"]
};
var monthValues4 = {
  narrow: ["e", "f", "m", "a", "m", "j", "j", "a", "s", "o", "n", "d"],
  abbreviated: ["ene", "feb", "mar", "abr", "may", "jun", "jul", "ago", "sep", "oct", "nov", "dic"],
  wide: ["enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"]
};
var dayValues4 = {
  narrow: ["d", "l", "m", "m", "j", "v", "s"],
  short: ["do", "lu", "ma", "mi", "ju", "vi", "s\xE1"],
  abbreviated: ["dom", "lun", "mar", "mi\xE9", "jue", "vie", "s\xE1b"],
  wide: ["domingo", "lunes", "martes", "mi\xE9rcoles", "jueves", "viernes", "s\xE1bado"]
};
var dayPeriodValues4 = {
  narrow: {
    am: "a",
    pm: "p",
    midnight: "mn",
    noon: "md",
    morning: "ma\xF1ana",
    afternoon: "tarde",
    evening: "tarde",
    night: "noche"
  },
  abbreviated: {
    am: "AM",
    pm: "PM",
    midnight: "medianoche",
    noon: "mediodia",
    morning: "ma\xF1ana",
    afternoon: "tarde",
    evening: "tarde",
    night: "noche"
  },
  wide: {
    am: "a.m.",
    pm: "p.m.",
    midnight: "medianoche",
    noon: "mediodia",
    morning: "ma\xF1ana",
    afternoon: "tarde",
    evening: "tarde",
    night: "noche"
  }
};
var formattingDayPeriodValues4 = {
  narrow: {
    am: "a",
    pm: "p",
    midnight: "mn",
    noon: "md",
    morning: "de la ma\xF1ana",
    afternoon: "de la tarde",
    evening: "de la tarde",
    night: "de la noche"
  },
  abbreviated: {
    am: "AM",
    pm: "PM",
    midnight: "medianoche",
    noon: "mediodia",
    morning: "de la ma\xF1ana",
    afternoon: "de la tarde",
    evening: "de la tarde",
    night: "de la noche"
  },
  wide: {
    am: "a.m.",
    pm: "p.m.",
    midnight: "medianoche",
    noon: "mediodia",
    morning: "de la ma\xF1ana",
    afternoon: "de la tarde",
    evening: "de la tarde",
    night: "de la noche"
  }
};
var ordinalNumber4 = (dirtyNumber, _options) => {
  const number = Number(dirtyNumber);
  return number + "\xBA";
};
var localize4 = {
  ordinalNumber: ordinalNumber4,
  era: buildLocalizeFn({
    values: eraValues4,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues4,
    defaultWidth: "wide",
    argumentCallback: (quarter) => Number(quarter) - 1
  }),
  month: buildLocalizeFn({
    values: monthValues4,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues4,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues4,
    defaultWidth: "wide",
    formattingValues: formattingDayPeriodValues4,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/es/_lib/match.js
var matchOrdinalNumberPattern4 = /^(\d+)(º)?/i;
var parseOrdinalNumberPattern4 = /\d+/i;
var matchEraPatterns4 = {
  narrow: /^(ac|dc|a|d)/i,
  abbreviated: /^(a\.?\s?c\.?|a\.?\s?e\.?\s?c\.?|d\.?\s?c\.?|e\.?\s?c\.?)/i,
  wide: /^(antes de cristo|antes de la era com[uú]n|despu[eé]s de cristo|era com[uú]n)/i
};
var parseEraPatterns4 = {
  any: [/^ac/i, /^dc/i],
  wide: [/^(antes de cristo|antes de la era com[uú]n)/i, /^(despu[eé]s de cristo|era com[uú]n)/i]
};
var matchQuarterPatterns4 = {
  narrow: /^[1234]/i,
  abbreviated: /^T[1234]/i,
  wide: /^[1234](º)? trimestre/i
};
var parseQuarterPatterns4 = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns4 = {
  narrow: /^[efmajsond]/i,
  abbreviated: /^(ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic)/i,
  wide: /^(enero|febrero|marzo|abril|mayo|junio|julio|agosto|septiembre|octubre|noviembre|diciembre)/i
};
var parseMonthPatterns4 = {
  narrow: [/^e/i, /^f/i, /^m/i, /^a/i, /^m/i, /^j/i, /^j/i, /^a/i, /^s/i, /^o/i, /^n/i, /^d/i],
  any: [/^en/i, /^feb/i, /^mar/i, /^abr/i, /^may/i, /^jun/i, /^jul/i, /^ago/i, /^sep/i, /^oct/i, /^nov/i, /^dic/i]
};
var matchDayPatterns4 = {
  narrow: /^[dlmjvs]/i,
  short: /^(do|lu|ma|mi|ju|vi|s[áa])/i,
  abbreviated: /^(dom|lun|mar|mi[ée]|jue|vie|s[áa]b)/i,
  wide: /^(domingo|lunes|martes|mi[ée]rcoles|jueves|viernes|s[áa]bado)/i
};
var parseDayPatterns4 = {
  narrow: [/^d/i, /^l/i, /^m/i, /^m/i, /^j/i, /^v/i, /^s/i],
  any: [/^do/i, /^lu/i, /^ma/i, /^mi/i, /^ju/i, /^vi/i, /^sa/i]
};
var matchDayPeriodPatterns4 = {
  narrow: /^(a|p|mn|md|(de la|a las) (mañana|tarde|noche))/i,
  any: /^([ap]\.?\s?m\.?|medianoche|mediodia|(de la|a las) (mañana|tarde|noche))/i
};
var parseDayPeriodPatterns4 = {
  any: {
    am: /^a/i,
    pm: /^p/i,
    midnight: /^mn/i,
    noon: /^md/i,
    morning: /mañana/i,
    afternoon: /tarde/i,
    evening: /tarde/i,
    night: /noche/i
  }
};
var match4 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern4,
    parsePattern: parseOrdinalNumberPattern4,
    valueCallback: function(value) {
      return parseInt(value, 10);
    }
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns4,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns4,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns4,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns4,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns4,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns4,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns4,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns4,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns4,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns4,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/es.js
var es = {
  code: "es",
  formatDistance: formatDistance4,
  formatLong: formatLong4,
  formatRelative: formatRelative4,
  localize: localize4,
  match: match4,
  options: {
    weekStartsOn: 1,
    firstWeekContainsDate: 1
  }
};

// node_modules/date-fns/locale/fr/_lib/formatDistance.js
var formatDistanceLocale5 = {
  lessThanXSeconds: {
    one: "moins d\u2019une seconde",
    other: "moins de {{count}} secondes"
  },
  xSeconds: {
    one: "1 seconde",
    other: "{{count}} secondes"
  },
  halfAMinute: "30 secondes",
  lessThanXMinutes: {
    one: "moins d\u2019une minute",
    other: "moins de {{count}} minutes"
  },
  xMinutes: {
    one: "1 minute",
    other: "{{count}} minutes"
  },
  aboutXHours: {
    one: "environ 1 heure",
    other: "environ {{count}} heures"
  },
  xHours: {
    one: "1 heure",
    other: "{{count}} heures"
  },
  xDays: {
    one: "1 jour",
    other: "{{count}} jours"
  },
  aboutXWeeks: {
    one: "environ 1 semaine",
    other: "environ {{count}} semaines"
  },
  xWeeks: {
    one: "1 semaine",
    other: "{{count}} semaines"
  },
  aboutXMonths: {
    one: "environ 1 mois",
    other: "environ {{count}} mois"
  },
  xMonths: {
    one: "1 mois",
    other: "{{count}} mois"
  },
  aboutXYears: {
    one: "environ 1 an",
    other: "environ {{count}} ans"
  },
  xYears: {
    one: "1 an",
    other: "{{count}} ans"
  },
  overXYears: {
    one: "plus d\u2019un an",
    other: "plus de {{count}} ans"
  },
  almostXYears: {
    one: "presqu\u2019un an",
    other: "presque {{count}} ans"
  }
};
var formatDistance5 = (token, count, options) => {
  let result;
  const form = formatDistanceLocale5[token];
  if (typeof form === "string") {
    result = form;
  } else if (count === 1) {
    result = form.one;
  } else {
    result = form.other.replace("{{count}}", String(count));
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return "dans " + result;
    } else {
      return "il y a " + result;
    }
  }
  return result;
};

// node_modules/date-fns/locale/fr/_lib/formatLong.js
var dateFormats5 = {
  full: "EEEE d MMMM y",
  long: "d MMMM y",
  medium: "d MMM y",
  short: "dd/MM/y"
};
var timeFormats5 = {
  full: "HH:mm:ss zzzz",
  long: "HH:mm:ss z",
  medium: "HH:mm:ss",
  short: "HH:mm"
};
var dateTimeFormats5 = {
  full: "{{date}} '\xE0' {{time}}",
  long: "{{date}} '\xE0' {{time}}",
  medium: "{{date}}, {{time}}",
  short: "{{date}}, {{time}}"
};
var formatLong5 = {
  date: buildFormatLongFn({
    formats: dateFormats5,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats5,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats5,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/fr/_lib/formatRelative.js
var formatRelativeLocale5 = {
  lastWeek: "eeee 'dernier \xE0' p",
  yesterday: "'hier \xE0' p",
  today: "'aujourd\u2019hui \xE0' p",
  tomorrow: "'demain \xE0' p'",
  nextWeek: "eeee 'prochain \xE0' p",
  other: "P"
};
var formatRelative5 = (token, _date, _baseDate, _options) => formatRelativeLocale5[token];

// node_modules/date-fns/locale/fr/_lib/localize.js
var eraValues5 = {
  narrow: ["av. J.-C", "ap. J.-C"],
  abbreviated: ["av. J.-C", "ap. J.-C"],
  wide: ["avant J\xE9sus-Christ", "apr\xE8s J\xE9sus-Christ"]
};
var quarterValues5 = {
  narrow: ["T1", "T2", "T3", "T4"],
  abbreviated: ["1er trim.", "2\xE8me trim.", "3\xE8me trim.", "4\xE8me trim."],
  wide: ["1er trimestre", "2\xE8me trimestre", "3\xE8me trimestre", "4\xE8me trimestre"]
};
var monthValues5 = {
  narrow: ["J", "F", "M", "A", "M", "J", "J", "A", "S", "O", "N", "D"],
  abbreviated: ["janv.", "f\xE9vr.", "mars", "avr.", "mai", "juin", "juil.", "ao\xFBt", "sept.", "oct.", "nov.", "d\xE9c."],
  wide: ["janvier", "f\xE9vrier", "mars", "avril", "mai", "juin", "juillet", "ao\xFBt", "septembre", "octobre", "novembre", "d\xE9cembre"]
};
var dayValues5 = {
  narrow: ["D", "L", "M", "M", "J", "V", "S"],
  short: ["di", "lu", "ma", "me", "je", "ve", "sa"],
  abbreviated: ["dim.", "lun.", "mar.", "mer.", "jeu.", "ven.", "sam."],
  wide: ["dimanche", "lundi", "mardi", "mercredi", "jeudi", "vendredi", "samedi"]
};
var dayPeriodValues5 = {
  narrow: {
    am: "AM",
    pm: "PM",
    midnight: "minuit",
    noon: "midi",
    morning: "mat.",
    afternoon: "ap.m.",
    evening: "soir",
    night: "mat."
  },
  abbreviated: {
    am: "AM",
    pm: "PM",
    midnight: "minuit",
    noon: "midi",
    morning: "matin",
    afternoon: "apr\xE8s-midi",
    evening: "soir",
    night: "matin"
  },
  wide: {
    am: "AM",
    pm: "PM",
    midnight: "minuit",
    noon: "midi",
    morning: "du matin",
    afternoon: "de l\u2019apr\xE8s-midi",
    evening: "du soir",
    night: "du matin"
  }
};
var ordinalNumber5 = (dirtyNumber, options) => {
  const number = Number(dirtyNumber);
  const unit = options?.unit;
  if (number === 0) return "0";
  const feminineUnits = ["year", "week", "hour", "minute", "second"];
  let suffix;
  if (number === 1) {
    suffix = unit && feminineUnits.includes(unit) ? "\xE8re" : "er";
  } else {
    suffix = "\xE8me";
  }
  return number + suffix;
};
var LONG_MONTHS_TOKENS = ["MMM", "MMMM"];
var localize5 = {
  preprocessor: (date, parts) => {
    if (date.getDate() === 1) return parts;
    const hasLongMonthToken = parts.some((part) => part.isToken && LONG_MONTHS_TOKENS.includes(part.value));
    if (!hasLongMonthToken) return parts;
    return parts.map((part) => part.isToken && part.value === "do" ? {
      isToken: true,
      value: "d"
    } : part);
  },
  ordinalNumber: ordinalNumber5,
  era: buildLocalizeFn({
    values: eraValues5,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues5,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues5,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues5,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues5,
    defaultWidth: "wide"
  })
};

// node_modules/date-fns/locale/fr/_lib/match.js
var matchOrdinalNumberPattern5 = /^(\d+)(ième|ère|ème|er|e)?/i;
var parseOrdinalNumberPattern5 = /\d+/i;
var matchEraPatterns5 = {
  narrow: /^(av\.J\.C|ap\.J\.C|ap\.J\.-C)/i,
  abbreviated: /^(av\.J\.-C|av\.J-C|apr\.J\.-C|apr\.J-C|ap\.J-C)/i,
  wide: /^(avant Jésus-Christ|après Jésus-Christ)/i
};
var parseEraPatterns5 = {
  any: [/^av/i, /^ap/i]
};
var matchQuarterPatterns5 = {
  narrow: /^T?[1234]/i,
  abbreviated: /^[1234](er|ème|e)? trim\.?/i,
  wide: /^[1234](er|ème|e)? trimestre/i
};
var parseQuarterPatterns5 = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns5 = {
  narrow: /^[jfmasond]/i,
  abbreviated: /^(janv|févr|mars|avr|mai|juin|juill|juil|août|sept|oct|nov|déc)\.?/i,
  wide: /^(janvier|février|mars|avril|mai|juin|juillet|août|septembre|octobre|novembre|décembre)/i
};
var parseMonthPatterns5 = {
  narrow: [/^j/i, /^f/i, /^m/i, /^a/i, /^m/i, /^j/i, /^j/i, /^a/i, /^s/i, /^o/i, /^n/i, /^d/i],
  any: [/^ja/i, /^f/i, /^mar/i, /^av/i, /^ma/i, /^juin/i, /^juil/i, /^ao/i, /^s/i, /^o/i, /^n/i, /^d/i]
};
var matchDayPatterns5 = {
  narrow: /^[lmjvsd]/i,
  short: /^(di|lu|ma|me|je|ve|sa)/i,
  abbreviated: /^(dim|lun|mar|mer|jeu|ven|sam)\.?/i,
  wide: /^(dimanche|lundi|mardi|mercredi|jeudi|vendredi|samedi)/i
};
var parseDayPatterns5 = {
  narrow: [/^d/i, /^l/i, /^m/i, /^m/i, /^j/i, /^v/i, /^s/i],
  any: [/^di/i, /^lu/i, /^ma/i, /^me/i, /^je/i, /^ve/i, /^sa/i]
};
var matchDayPeriodPatterns5 = {
  narrow: /^(a|p|minuit|midi|mat\.?|ap\.?m\.?|soir|nuit)/i,
  any: /^([ap]\.?\s?m\.?|du matin|de l'après[-\s]midi|du soir|de la nuit)/i
};
var parseDayPeriodPatterns5 = {
  any: {
    am: /^a/i,
    pm: /^p/i,
    midnight: /^min/i,
    noon: /^mid/i,
    morning: /mat/i,
    afternoon: /ap/i,
    evening: /soir/i,
    night: /nuit/i
  }
};
var match5 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern5,
    parsePattern: parseOrdinalNumberPattern5,
    valueCallback: (value) => parseInt(value)
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns5,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns5,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns5,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns5,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns5,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns5,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns5,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns5,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns5,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns5,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/fr.js
var fr = {
  code: "fr",
  formatDistance: formatDistance5,
  formatLong: formatLong5,
  formatRelative: formatRelative5,
  localize: localize5,
  match: match5,
  options: {
    weekStartsOn: 1,
    firstWeekContainsDate: 4
  }
};

// node_modules/date-fns/locale/hi/_lib/localize.js
var numberValues = {
  locale: {
    1: "\u0967",
    2: "\u0968",
    3: "\u0969",
    4: "\u096A",
    5: "\u096B",
    6: "\u096C",
    7: "\u096D",
    8: "\u096E",
    9: "\u096F",
    0: "\u0966"
  },
  number: {
    "\u0967": "1",
    "\u0968": "2",
    "\u0969": "3",
    "\u096A": "4",
    "\u096B": "5",
    "\u096C": "6",
    "\u096D": "7",
    "\u096E": "8",
    "\u096F": "9",
    "\u0966": "0"
  }
};
var eraValues6 = {
  narrow: ["\u0908\u0938\u093E-\u092A\u0942\u0930\u094D\u0935", "\u0908\u0938\u094D\u0935\u0940"],
  abbreviated: ["\u0908\u0938\u093E-\u092A\u0942\u0930\u094D\u0935", "\u0908\u0938\u094D\u0935\u0940"],
  wide: ["\u0908\u0938\u093E-\u092A\u0942\u0930\u094D\u0935", "\u0908\u0938\u0935\u0940 \u0938\u0928"]
};
var quarterValues6 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["\u0924\u093F1", "\u0924\u093F2", "\u0924\u093F3", "\u0924\u093F4"],
  wide: ["\u092A\u0939\u0932\u0940 \u0924\u093F\u092E\u093E\u0939\u0940", "\u0926\u0942\u0938\u0930\u0940 \u0924\u093F\u092E\u093E\u0939\u0940", "\u0924\u0940\u0938\u0930\u0940 \u0924\u093F\u092E\u093E\u0939\u0940", "\u091A\u094C\u0925\u0940 \u0924\u093F\u092E\u093E\u0939\u0940"]
};
var monthValues6 = {
  narrow: ["\u091C", "\u092B\u093C", "\u092E\u093E", "\u0905", "\u092E\u0908", "\u091C\u0942", "\u091C\u0941", "\u0905\u0917", "\u0938\u093F", "\u0905\u0915\u094D\u091F\u0942", "\u0928", "\u0926\u093F"],
  abbreviated: ["\u091C\u0928", "\u092B\u093C\u0930", "\u092E\u093E\u0930\u094D\u091A", "\u0905\u092A\u094D\u0930\u0948\u0932", "\u092E\u0908", "\u091C\u0942\u0928", "\u091C\u0941\u0932", "\u0905\u0917", "\u0938\u093F\u0924", "\u0905\u0915\u094D\u091F\u0942", "\u0928\u0935", "\u0926\u093F\u0938"],
  wide: ["\u091C\u0928\u0935\u0930\u0940", "\u092B\u093C\u0930\u0935\u0930\u0940", "\u092E\u093E\u0930\u094D\u091A", "\u0905\u092A\u094D\u0930\u0948\u0932", "\u092E\u0908", "\u091C\u0942\u0928", "\u091C\u0941\u0932\u093E\u0908", "\u0905\u0917\u0938\u094D\u0924", "\u0938\u093F\u0924\u0902\u092C\u0930", "\u0905\u0915\u094D\u091F\u0942\u092C\u0930", "\u0928\u0935\u0902\u092C\u0930", "\u0926\u093F\u0938\u0902\u092C\u0930"]
};
var dayValues6 = {
  narrow: ["\u0930", "\u0938\u094B", "\u092E\u0902", "\u092C\u0941", "\u0917\u0941", "\u0936\u0941", "\u0936"],
  short: ["\u0930", "\u0938\u094B", "\u092E\u0902", "\u092C\u0941", "\u0917\u0941", "\u0936\u0941", "\u0936"],
  abbreviated: ["\u0930\u0935\u093F", "\u0938\u094B\u092E", "\u092E\u0902\u0917\u0932", "\u092C\u0941\u0927", "\u0917\u0941\u0930\u0941", "\u0936\u0941\u0915\u094D\u0930", "\u0936\u0928\u093F"],
  wide: ["\u0930\u0935\u093F\u0935\u093E\u0930", "\u0938\u094B\u092E\u0935\u093E\u0930", "\u092E\u0902\u0917\u0932\u0935\u093E\u0930", "\u092C\u0941\u0927\u0935\u093E\u0930", "\u0917\u0941\u0930\u0941\u0935\u093E\u0930", "\u0936\u0941\u0915\u094D\u0930\u0935\u093E\u0930", "\u0936\u0928\u093F\u0935\u093E\u0930"]
};
var dayPeriodValues6 = {
  narrow: {
    am: "\u092A\u0942\u0930\u094D\u0935\u093E\u0939\u094D\u0928",
    pm: "\u0905\u092A\u0930\u093E\u0939\u094D\u0928",
    midnight: "\u092E\u0927\u094D\u092F\u0930\u093E\u0924\u094D\u0930\u093F",
    noon: "\u0926\u094B\u092A\u0939\u0930",
    morning: "\u0938\u0941\u092C\u0939",
    afternoon: "\u0926\u094B\u092A\u0939\u0930",
    evening: "\u0936\u093E\u092E",
    night: "\u0930\u093E\u0924"
  },
  abbreviated: {
    am: "\u092A\u0942\u0930\u094D\u0935\u093E\u0939\u094D\u0928",
    pm: "\u0905\u092A\u0930\u093E\u0939\u094D\u0928",
    midnight: "\u092E\u0927\u094D\u092F\u0930\u093E\u0924\u094D\u0930\u093F",
    noon: "\u0926\u094B\u092A\u0939\u0930",
    morning: "\u0938\u0941\u092C\u0939",
    afternoon: "\u0926\u094B\u092A\u0939\u0930",
    evening: "\u0936\u093E\u092E",
    night: "\u0930\u093E\u0924"
  },
  wide: {
    am: "\u092A\u0942\u0930\u094D\u0935\u093E\u0939\u094D\u0928",
    pm: "\u0905\u092A\u0930\u093E\u0939\u094D\u0928",
    midnight: "\u092E\u0927\u094D\u092F\u0930\u093E\u0924\u094D\u0930\u093F",
    noon: "\u0926\u094B\u092A\u0939\u0930",
    morning: "\u0938\u0941\u092C\u0939",
    afternoon: "\u0926\u094B\u092A\u0939\u0930",
    evening: "\u0936\u093E\u092E",
    night: "\u0930\u093E\u0924"
  }
};
var formattingDayPeriodValues5 = {
  narrow: {
    am: "\u092A\u0942\u0930\u094D\u0935\u093E\u0939\u094D\u0928",
    pm: "\u0905\u092A\u0930\u093E\u0939\u094D\u0928",
    midnight: "\u092E\u0927\u094D\u092F\u0930\u093E\u0924\u094D\u0930\u093F",
    noon: "\u0926\u094B\u092A\u0939\u0930",
    morning: "\u0938\u0941\u092C\u0939",
    afternoon: "\u0926\u094B\u092A\u0939\u0930",
    evening: "\u0936\u093E\u092E",
    night: "\u0930\u093E\u0924"
  },
  abbreviated: {
    am: "\u092A\u0942\u0930\u094D\u0935\u093E\u0939\u094D\u0928",
    pm: "\u0905\u092A\u0930\u093E\u0939\u094D\u0928",
    midnight: "\u092E\u0927\u094D\u092F\u0930\u093E\u0924\u094D\u0930\u093F",
    noon: "\u0926\u094B\u092A\u0939\u0930",
    morning: "\u0938\u0941\u092C\u0939",
    afternoon: "\u0926\u094B\u092A\u0939\u0930",
    evening: "\u0936\u093E\u092E",
    night: "\u0930\u093E\u0924"
  },
  wide: {
    am: "\u092A\u0942\u0930\u094D\u0935\u093E\u0939\u094D\u0928",
    pm: "\u0905\u092A\u0930\u093E\u0939\u094D\u0928",
    midnight: "\u092E\u0927\u094D\u092F\u0930\u093E\u0924\u094D\u0930\u093F",
    noon: "\u0926\u094B\u092A\u0939\u0930",
    morning: "\u0938\u0941\u092C\u0939",
    afternoon: "\u0926\u094B\u092A\u0939\u0930",
    evening: "\u0936\u093E\u092E",
    night: "\u0930\u093E\u0924"
  }
};
var ordinalNumber6 = (dirtyNumber, _options) => {
  const number = Number(dirtyNumber);
  return numberToLocale(number);
};
function localeToNumber(locale) {
  const enNumber = locale.toString().replace(/[१२३४५६७८९०]/g, function(match14) {
    return numberValues.number[match14];
  });
  return Number(enNumber);
}
function numberToLocale(enNumber) {
  return enNumber.toString().replace(/\d/g, function(match14) {
    return numberValues.locale[match14];
  });
}
var localize6 = {
  ordinalNumber: ordinalNumber6,
  era: buildLocalizeFn({
    values: eraValues6,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues6,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues6,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues6,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues6,
    defaultWidth: "wide",
    formattingValues: formattingDayPeriodValues5,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/hi/_lib/formatDistance.js
var formatDistanceLocale6 = {
  lessThanXSeconds: {
    one: "\u0967 \u0938\u0947\u0915\u0902\u0921 \u0938\u0947 \u0915\u092E",
    // CLDR #1310
    other: "{{count}} \u0938\u0947\u0915\u0902\u0921 \u0938\u0947 \u0915\u092E"
  },
  xSeconds: {
    one: "\u0967 \u0938\u0947\u0915\u0902\u0921",
    other: "{{count}} \u0938\u0947\u0915\u0902\u0921"
  },
  halfAMinute: "\u0906\u0927\u093E \u092E\u093F\u0928\u091F",
  lessThanXMinutes: {
    one: "\u0967 \u092E\u093F\u0928\u091F \u0938\u0947 \u0915\u092E",
    other: "{{count}} \u092E\u093F\u0928\u091F \u0938\u0947 \u0915\u092E"
  },
  xMinutes: {
    one: "\u0967 \u092E\u093F\u0928\u091F",
    // CLDR #1307
    other: "{{count}} \u092E\u093F\u0928\u091F"
  },
  aboutXHours: {
    one: "\u0932\u0917\u092D\u0917 \u0967 \u0918\u0902\u091F\u093E",
    other: "\u0932\u0917\u092D\u0917 {{count}} \u0918\u0902\u091F\u0947"
  },
  xHours: {
    one: "\u0967 \u0918\u0902\u091F\u093E",
    // CLDR #1304
    other: "{{count}} \u0918\u0902\u091F\u0947"
    // CLDR #4467
  },
  xDays: {
    one: "\u0967 \u0926\u093F\u0928",
    // CLDR #1286
    other: "{{count}} \u0926\u093F\u0928"
  },
  aboutXWeeks: {
    one: "\u0932\u0917\u092D\u0917 \u0967 \u0938\u092A\u094D\u0924\u093E\u0939",
    other: "\u0932\u0917\u092D\u0917 {{count}} \u0938\u092A\u094D\u0924\u093E\u0939"
  },
  xWeeks: {
    one: "\u0967 \u0938\u092A\u094D\u0924\u093E\u0939",
    other: "{{count}} \u0938\u092A\u094D\u0924\u093E\u0939"
  },
  aboutXMonths: {
    one: "\u0932\u0917\u092D\u0917 \u0967 \u092E\u0939\u0940\u0928\u093E",
    other: "\u0932\u0917\u092D\u0917 {{count}} \u092E\u0939\u0940\u0928\u0947"
  },
  xMonths: {
    one: "\u0967 \u092E\u0939\u0940\u0928\u093E",
    other: "{{count}} \u092E\u0939\u0940\u0928\u0947"
  },
  aboutXYears: {
    one: "\u0932\u0917\u092D\u0917 \u0967 \u0935\u0930\u094D\u0937",
    other: "\u0932\u0917\u092D\u0917 {{count}} \u0935\u0930\u094D\u0937"
    // CLDR #4823
  },
  xYears: {
    one: "\u0967 \u0935\u0930\u094D\u0937",
    other: "{{count}} \u0935\u0930\u094D\u0937"
  },
  overXYears: {
    one: "\u0967 \u0935\u0930\u094D\u0937 \u0938\u0947 \u0905\u0927\u093F\u0915",
    other: "{{count}} \u0935\u0930\u094D\u0937 \u0938\u0947 \u0905\u0927\u093F\u0915"
  },
  almostXYears: {
    one: "\u0932\u0917\u092D\u0917 \u0967 \u0935\u0930\u094D\u0937",
    other: "\u0932\u0917\u092D\u0917 {{count}} \u0935\u0930\u094D\u0937"
  }
};
var formatDistance6 = (token, count, options) => {
  let result;
  const tokenValue = formatDistanceLocale6[token];
  if (typeof tokenValue === "string") {
    result = tokenValue;
  } else if (count === 1) {
    result = tokenValue.one;
  } else {
    result = tokenValue.other.replace("{{count}}", numberToLocale(count));
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return result + "\u092E\u0947 ";
    } else {
      return result + " \u092A\u0939\u0932\u0947";
    }
  }
  return result;
};

// node_modules/date-fns/locale/hi/_lib/formatLong.js
var dateFormats6 = {
  full: "EEEE, do MMMM, y",
  // CLDR #1787
  long: "do MMMM, y",
  // CLDR #1788
  medium: "d MMM, y",
  // CLDR #1789
  short: "dd/MM/yyyy"
  // CLDR #1790
};
var timeFormats6 = {
  full: "h:mm:ss a zzzz",
  // CLDR #1791
  long: "h:mm:ss a z",
  // CLDR #1792
  medium: "h:mm:ss a",
  // CLDR #1793
  short: "h:mm a"
  // CLDR #1794
};
var dateTimeFormats6 = {
  full: "{{date}} '\u0915\u094B' {{time}}",
  // CLDR #1795
  long: "{{date}} '\u0915\u094B' {{time}}",
  // CLDR #1796
  medium: "{{date}}, {{time}}",
  // CLDR #1797
  short: "{{date}}, {{time}}"
  // CLDR #1798
};
var formatLong6 = {
  date: buildFormatLongFn({
    formats: dateFormats6,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats6,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats6,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/hi/_lib/formatRelative.js
var formatRelativeLocale6 = {
  lastWeek: "'\u092A\u093F\u091B\u0932\u0947' eeee p",
  yesterday: "'\u0915\u0932' p",
  today: "'\u0906\u091C' p",
  tomorrow: "'\u0915\u0932' p",
  nextWeek: "eeee '\u0915\u094B' p",
  other: "P"
};
var formatRelative6 = (token, _date, _baseDate, _options) => formatRelativeLocale6[token];

// node_modules/date-fns/locale/hi/_lib/match.js
var matchOrdinalNumberPattern6 = /^[०१२३४५६७८९]+/i;
var parseOrdinalNumberPattern6 = /^[०१२३४५६७८९]+/i;
var matchEraPatterns6 = {
  narrow: /^(ईसा-पूर्व|ईस्वी)/i,
  abbreviated: /^(ईसा\.?\s?पूर्व\.?|ईसा\.?)/i,
  wide: /^(ईसा-पूर्व|ईसवी पूर्व|ईसवी सन|ईसवी)/i
};
var parseEraPatterns6 = {
  any: [/^b/i, /^(a|c)/i]
};
var matchQuarterPatterns6 = {
  narrow: /^[1234]/i,
  abbreviated: /^ति[1234]/i,
  wide: /^[1234](पहली|दूसरी|तीसरी|चौथी)? तिमाही/i
};
var parseQuarterPatterns6 = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns6 = {
  // eslint-disable-next-line no-misleading-character-class
  narrow: /^[जफ़माअप्मईजूनजुअगसिअक्तनदि]/i,
  abbreviated: /^(जन|फ़र|मार्च|अप्|मई|जून|जुल|अग|सित|अक्तू|नव|दिस)/i,
  wide: /^(जनवरी|फ़रवरी|मार्च|अप्रैल|मई|जून|जुलाई|अगस्त|सितंबर|अक्तूबर|नवंबर|दिसंबर)/i
};
var parseMonthPatterns6 = {
  narrow: [/^ज/i, /^फ़/i, /^मा/i, /^अप्/i, /^मई/i, /^जू/i, /^जु/i, /^अग/i, /^सि/i, /^अक्तू/i, /^न/i, /^दि/i],
  any: [/^जन/i, /^फ़/i, /^मा/i, /^अप्/i, /^मई/i, /^जू/i, /^जु/i, /^अग/i, /^सि/i, /^अक्तू/i, /^नव/i, /^दिस/i]
};
var matchDayPatterns6 = {
  // eslint-disable-next-line no-misleading-character-class
  narrow: /^[रविसोममंगलबुधगुरुशुक्रशनि]/i,
  short: /^(रवि|सोम|मंगल|बुध|गुरु|शुक्र|शनि)/i,
  abbreviated: /^(रवि|सोम|मंगल|बुध|गुरु|शुक्र|शनि)/i,
  wide: /^(रविवार|सोमवार|मंगलवार|बुधवार|गुरुवार|शुक्रवार|शनिवार)/i
};
var parseDayPatterns6 = {
  narrow: [/^रवि/i, /^सोम/i, /^मंगल/i, /^बुध/i, /^गुरु/i, /^शुक्र/i, /^शनि/i],
  any: [/^रवि/i, /^सोम/i, /^मंगल/i, /^बुध/i, /^गुरु/i, /^शुक्र/i, /^शनि/i]
};
var matchDayPeriodPatterns6 = {
  narrow: /^(पू|अ|म|द.\?|सु|दो|शा|रा)/i,
  any: /^(पूर्वाह्न|अपराह्न|म|द.\?|सु|दो|शा|रा)/i
};
var parseDayPeriodPatterns6 = {
  any: {
    am: /^पूर्वाह्न/i,
    pm: /^अपराह्न/i,
    midnight: /^मध्य/i,
    noon: /^दो/i,
    morning: /सु/i,
    afternoon: /दो/i,
    evening: /शा/i,
    night: /रा/i
  }
};
var match6 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern6,
    parsePattern: parseOrdinalNumberPattern6,
    valueCallback: localeToNumber
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns6,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns6,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns6,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns6,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns6,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns6,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns6,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns6,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns6,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns6,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/hi.js
var hi = {
  code: "hi",
  formatDistance: formatDistance6,
  formatLong: formatLong6,
  formatRelative: formatRelative6,
  localize: localize6,
  match: match6,
  options: {
    weekStartsOn: 0,
    firstWeekContainsDate: 4
  }
};

// node_modules/date-fns/locale/ja/_lib/formatDistance.js
var formatDistanceLocale7 = {
  lessThanXSeconds: {
    one: "1\u79D2\u672A\u6E80",
    other: "{{count}}\u79D2\u672A\u6E80",
    oneWithSuffix: "\u7D041\u79D2",
    otherWithSuffix: "\u7D04{{count}}\u79D2"
  },
  xSeconds: {
    one: "1\u79D2",
    other: "{{count}}\u79D2"
  },
  halfAMinute: "30\u79D2",
  lessThanXMinutes: {
    one: "1\u5206\u672A\u6E80",
    other: "{{count}}\u5206\u672A\u6E80",
    oneWithSuffix: "\u7D041\u5206",
    otherWithSuffix: "\u7D04{{count}}\u5206"
  },
  xMinutes: {
    one: "1\u5206",
    other: "{{count}}\u5206"
  },
  aboutXHours: {
    one: "\u7D041\u6642\u9593",
    other: "\u7D04{{count}}\u6642\u9593"
  },
  xHours: {
    one: "1\u6642\u9593",
    other: "{{count}}\u6642\u9593"
  },
  xDays: {
    one: "1\u65E5",
    other: "{{count}}\u65E5"
  },
  aboutXWeeks: {
    one: "\u7D041\u9031\u9593",
    other: "\u7D04{{count}}\u9031\u9593"
  },
  xWeeks: {
    one: "1\u9031\u9593",
    other: "{{count}}\u9031\u9593"
  },
  aboutXMonths: {
    one: "\u7D041\u304B\u6708",
    other: "\u7D04{{count}}\u304B\u6708"
  },
  xMonths: {
    one: "1\u304B\u6708",
    other: "{{count}}\u304B\u6708"
  },
  aboutXYears: {
    one: "\u7D041\u5E74",
    other: "\u7D04{{count}}\u5E74"
  },
  xYears: {
    one: "1\u5E74",
    other: "{{count}}\u5E74"
  },
  overXYears: {
    one: "1\u5E74\u4EE5\u4E0A",
    other: "{{count}}\u5E74\u4EE5\u4E0A"
  },
  almostXYears: {
    one: "1\u5E74\u8FD1\u304F",
    other: "{{count}}\u5E74\u8FD1\u304F"
  }
};
var formatDistance7 = (token, count, options) => {
  options = options || {};
  let result;
  const tokenValue = formatDistanceLocale7[token];
  if (typeof tokenValue === "string") {
    result = tokenValue;
  } else if (count === 1) {
    if (options.addSuffix && tokenValue.oneWithSuffix) {
      result = tokenValue.oneWithSuffix;
    } else {
      result = tokenValue.one;
    }
  } else {
    if (options.addSuffix && tokenValue.otherWithSuffix) {
      result = tokenValue.otherWithSuffix.replace("{{count}}", String(count));
    } else {
      result = tokenValue.other.replace("{{count}}", String(count));
    }
  }
  if (options.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return result + "\u5F8C";
    } else {
      return result + "\u524D";
    }
  }
  return result;
};

// node_modules/date-fns/locale/ja/_lib/formatLong.js
var dateFormats7 = {
  full: "y\u5E74M\u6708d\u65E5EEEE",
  long: "y\u5E74M\u6708d\u65E5",
  medium: "y/MM/dd",
  short: "y/MM/dd"
};
var timeFormats7 = {
  full: "H\u6642mm\u5206ss\u79D2 zzzz",
  long: "H:mm:ss z",
  medium: "H:mm:ss",
  short: "H:mm"
};
var dateTimeFormats7 = {
  full: "{{date}} {{time}}",
  long: "{{date}} {{time}}",
  medium: "{{date}} {{time}}",
  short: "{{date}} {{time}}"
};
var formatLong7 = {
  date: buildFormatLongFn({
    formats: dateFormats7,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats7,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats7,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/ja/_lib/formatRelative.js
var formatRelativeLocale7 = {
  lastWeek: "\u5148\u9031\u306Eeeee\u306Ep",
  yesterday: "\u6628\u65E5\u306Ep",
  today: "\u4ECA\u65E5\u306Ep",
  tomorrow: "\u660E\u65E5\u306Ep",
  nextWeek: "\u7FCC\u9031\u306Eeeee\u306Ep",
  other: "P"
};
var formatRelative7 = (token, _date, _baseDate, _options) => {
  return formatRelativeLocale7[token];
};

// node_modules/date-fns/locale/ja/_lib/localize.js
var eraValues7 = {
  narrow: ["BC", "AC"],
  abbreviated: ["\u7D00\u5143\u524D", "\u897F\u66A6"],
  wide: ["\u7D00\u5143\u524D", "\u897F\u66A6"]
};
var quarterValues7 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["Q1", "Q2", "Q3", "Q4"],
  wide: ["\u7B2C1\u56DB\u534A\u671F", "\u7B2C2\u56DB\u534A\u671F", "\u7B2C3\u56DB\u534A\u671F", "\u7B2C4\u56DB\u534A\u671F"]
};
var monthValues7 = {
  narrow: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"],
  abbreviated: ["1\u6708", "2\u6708", "3\u6708", "4\u6708", "5\u6708", "6\u6708", "7\u6708", "8\u6708", "9\u6708", "10\u6708", "11\u6708", "12\u6708"],
  wide: ["1\u6708", "2\u6708", "3\u6708", "4\u6708", "5\u6708", "6\u6708", "7\u6708", "8\u6708", "9\u6708", "10\u6708", "11\u6708", "12\u6708"]
};
var dayValues7 = {
  narrow: ["\u65E5", "\u6708", "\u706B", "\u6C34", "\u6728", "\u91D1", "\u571F"],
  short: ["\u65E5", "\u6708", "\u706B", "\u6C34", "\u6728", "\u91D1", "\u571F"],
  abbreviated: ["\u65E5", "\u6708", "\u706B", "\u6C34", "\u6728", "\u91D1", "\u571F"],
  wide: ["\u65E5\u66DC\u65E5", "\u6708\u66DC\u65E5", "\u706B\u66DC\u65E5", "\u6C34\u66DC\u65E5", "\u6728\u66DC\u65E5", "\u91D1\u66DC\u65E5", "\u571F\u66DC\u65E5"]
};
var dayPeriodValues7 = {
  narrow: {
    am: "\u5348\u524D",
    pm: "\u5348\u5F8C",
    midnight: "\u6DF1\u591C",
    noon: "\u6B63\u5348",
    morning: "\u671D",
    afternoon: "\u5348\u5F8C",
    evening: "\u591C",
    night: "\u6DF1\u591C"
  },
  abbreviated: {
    am: "\u5348\u524D",
    pm: "\u5348\u5F8C",
    midnight: "\u6DF1\u591C",
    noon: "\u6B63\u5348",
    morning: "\u671D",
    afternoon: "\u5348\u5F8C",
    evening: "\u591C",
    night: "\u6DF1\u591C"
  },
  wide: {
    am: "\u5348\u524D",
    pm: "\u5348\u5F8C",
    midnight: "\u6DF1\u591C",
    noon: "\u6B63\u5348",
    morning: "\u671D",
    afternoon: "\u5348\u5F8C",
    evening: "\u591C",
    night: "\u6DF1\u591C"
  }
};
var formattingDayPeriodValues6 = {
  narrow: {
    am: "\u5348\u524D",
    pm: "\u5348\u5F8C",
    midnight: "\u6DF1\u591C",
    noon: "\u6B63\u5348",
    morning: "\u671D",
    afternoon: "\u5348\u5F8C",
    evening: "\u591C",
    night: "\u6DF1\u591C"
  },
  abbreviated: {
    am: "\u5348\u524D",
    pm: "\u5348\u5F8C",
    midnight: "\u6DF1\u591C",
    noon: "\u6B63\u5348",
    morning: "\u671D",
    afternoon: "\u5348\u5F8C",
    evening: "\u591C",
    night: "\u6DF1\u591C"
  },
  wide: {
    am: "\u5348\u524D",
    pm: "\u5348\u5F8C",
    midnight: "\u6DF1\u591C",
    noon: "\u6B63\u5348",
    morning: "\u671D",
    afternoon: "\u5348\u5F8C",
    evening: "\u591C",
    night: "\u6DF1\u591C"
  }
};
var ordinalNumber7 = (dirtyNumber, options) => {
  const number = Number(dirtyNumber);
  const unit = String(options?.unit);
  switch (unit) {
    case "year":
      return `${number}\u5E74`;
    case "quarter":
      return `\u7B2C${number}\u56DB\u534A\u671F`;
    case "month":
      return `${number}\u6708`;
    case "week":
      return `\u7B2C${number}\u9031`;
    case "date":
      return `${number}\u65E5`;
    case "hour":
      return `${number}\u6642`;
    case "minute":
      return `${number}\u5206`;
    case "second":
      return `${number}\u79D2`;
    default:
      return `${number}`;
  }
};
var localize7 = {
  ordinalNumber: ordinalNumber7,
  era: buildLocalizeFn({
    values: eraValues7,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues7,
    defaultWidth: "wide",
    argumentCallback: (quarter) => Number(quarter) - 1
  }),
  month: buildLocalizeFn({
    values: monthValues7,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues7,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues7,
    defaultWidth: "wide",
    formattingValues: formattingDayPeriodValues6,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/ja/_lib/match.js
var matchOrdinalNumberPattern7 = /^第?\d+(年|四半期|月|週|日|時|分|秒)?/i;
var parseOrdinalNumberPattern7 = /\d+/i;
var matchEraPatterns7 = {
  narrow: /^(B\.?C\.?|A\.?D\.?)/i,
  abbreviated: /^(紀元[前後]|西暦)/i,
  wide: /^(紀元[前後]|西暦)/i
};
var parseEraPatterns7 = {
  narrow: [/^B/i, /^A/i],
  any: [/^(紀元前)/i, /^(西暦|紀元後)/i]
};
var matchQuarterPatterns7 = {
  narrow: /^[1234]/i,
  abbreviated: /^Q[1234]/i,
  wide: /^第[1234一二三四１２３４]四半期/i
};
var parseQuarterPatterns7 = {
  any: [/(1|一|１)/i, /(2|二|２)/i, /(3|三|３)/i, /(4|四|４)/i]
};
var matchMonthPatterns7 = {
  narrow: /^([123456789]|1[012])/,
  abbreviated: /^([123456789]|1[012])月/i,
  wide: /^([123456789]|1[012])月/i
};
var parseMonthPatterns7 = {
  any: [/^1\D/, /^2/, /^3/, /^4/, /^5/, /^6/, /^7/, /^8/, /^9/, /^10/, /^11/, /^12/]
};
var matchDayPatterns7 = {
  narrow: /^[日月火水木金土]/,
  short: /^[日月火水木金土]/,
  abbreviated: /^[日月火水木金土]/,
  wide: /^[日月火水木金土]曜日/
};
var parseDayPatterns7 = {
  any: [/^日/, /^月/, /^火/, /^水/, /^木/, /^金/, /^土/]
};
var matchDayPeriodPatterns7 = {
  any: /^(AM|PM|午前|午後|正午|深夜|真夜中|夜|朝)/i
};
var parseDayPeriodPatterns7 = {
  any: {
    am: /^(A|午前)/i,
    pm: /^(P|午後)/i,
    midnight: /^深夜|真夜中/i,
    noon: /^正午/i,
    morning: /^朝/i,
    afternoon: /^午後/i,
    evening: /^夜/i,
    night: /^深夜/i
  }
};
var match7 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern7,
    parsePattern: parseOrdinalNumberPattern7,
    valueCallback: function(value) {
      return parseInt(value, 10);
    }
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns7,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns7,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns7,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns7,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns7,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns7,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns7,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns7,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns7,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns7,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/ja.js
var ja = {
  code: "ja",
  formatDistance: formatDistance7,
  formatLong: formatLong7,
  formatRelative: formatRelative7,
  localize: localize7,
  match: match7,
  options: {
    weekStartsOn: 0,
    firstWeekContainsDate: 1
  }
};

// node_modules/date-fns/locale/nl/_lib/formatDistance.js
var formatDistanceLocale8 = {
  lessThanXSeconds: {
    one: "minder dan een seconde",
    other: "minder dan {{count}} seconden"
  },
  xSeconds: {
    one: "1 seconde",
    other: "{{count}} seconden"
  },
  halfAMinute: "een halve minuut",
  lessThanXMinutes: {
    one: "minder dan een minuut",
    other: "minder dan {{count}} minuten"
  },
  xMinutes: {
    one: "een minuut",
    other: "{{count}} minuten"
  },
  aboutXHours: {
    one: "ongeveer 1 uur",
    other: "ongeveer {{count}} uur"
  },
  xHours: {
    one: "1 uur",
    other: "{{count}} uur"
  },
  xDays: {
    one: "1 dag",
    other: "{{count}} dagen"
  },
  aboutXWeeks: {
    one: "ongeveer 1 week",
    other: "ongeveer {{count}} weken"
  },
  xWeeks: {
    one: "1 week",
    other: "{{count}} weken"
  },
  aboutXMonths: {
    one: "ongeveer 1 maand",
    other: "ongeveer {{count}} maanden"
  },
  xMonths: {
    one: "1 maand",
    other: "{{count}} maanden"
  },
  aboutXYears: {
    one: "ongeveer 1 jaar",
    other: "ongeveer {{count}} jaar"
  },
  xYears: {
    one: "1 jaar",
    other: "{{count}} jaar"
  },
  overXYears: {
    one: "meer dan 1 jaar",
    other: "meer dan {{count}} jaar"
  },
  almostXYears: {
    one: "bijna 1 jaar",
    other: "bijna {{count}} jaar"
  }
};
var formatDistance8 = (token, count, options) => {
  let result;
  const tokenValue = formatDistanceLocale8[token];
  if (typeof tokenValue === "string") {
    result = tokenValue;
  } else if (count === 1) {
    result = tokenValue.one;
  } else {
    result = tokenValue.other.replace("{{count}}", String(count));
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return "over " + result;
    } else {
      return result + " geleden";
    }
  }
  return result;
};

// node_modules/date-fns/locale/nl/_lib/formatLong.js
var dateFormats8 = {
  full: "EEEE d MMMM y",
  long: "d MMMM y",
  medium: "d MMM y",
  short: "dd-MM-y"
};
var timeFormats8 = {
  full: "HH:mm:ss zzzz",
  long: "HH:mm:ss z",
  medium: "HH:mm:ss",
  short: "HH:mm"
};
var dateTimeFormats8 = {
  full: "{{date}} 'om' {{time}}",
  long: "{{date}} 'om' {{time}}",
  medium: "{{date}}, {{time}}",
  short: "{{date}}, {{time}}"
};
var formatLong8 = {
  date: buildFormatLongFn({
    formats: dateFormats8,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats8,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats8,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/nl/_lib/formatRelative.js
var formatRelativeLocale8 = {
  lastWeek: "'afgelopen' eeee 'om' p",
  yesterday: "'gisteren om' p",
  today: "'vandaag om' p",
  tomorrow: "'morgen om' p",
  nextWeek: "eeee 'om' p",
  other: "P"
};
var formatRelative8 = (token, _date, _baseDate, _options) => formatRelativeLocale8[token];

// node_modules/date-fns/locale/nl/_lib/localize.js
var eraValues8 = {
  narrow: ["v.C.", "n.C."],
  abbreviated: ["v.Chr.", "n.Chr."],
  wide: ["voor Christus", "na Christus"]
};
var quarterValues8 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["K1", "K2", "K3", "K4"],
  wide: ["1e kwartaal", "2e kwartaal", "3e kwartaal", "4e kwartaal"]
};
var monthValues8 = {
  narrow: ["J", "F", "M", "A", "M", "J", "J", "A", "S", "O", "N", "D"],
  abbreviated: ["jan.", "feb.", "mrt.", "apr.", "mei", "jun.", "jul.", "aug.", "sep.", "okt.", "nov.", "dec."],
  wide: ["januari", "februari", "maart", "april", "mei", "juni", "juli", "augustus", "september", "oktober", "november", "december"]
};
var dayValues8 = {
  narrow: ["Z", "M", "D", "W", "D", "V", "Z"],
  short: ["zo", "ma", "di", "wo", "do", "vr", "za"],
  abbreviated: ["zon", "maa", "din", "woe", "don", "vri", "zat"],
  wide: ["zondag", "maandag", "dinsdag", "woensdag", "donderdag", "vrijdag", "zaterdag"]
};
var dayPeriodValues8 = {
  narrow: {
    am: "AM",
    pm: "PM",
    midnight: "middernacht",
    noon: "het middaguur",
    morning: "'s ochtends",
    afternoon: "'s middags",
    evening: "'s avonds",
    night: "'s nachts"
  },
  abbreviated: {
    am: "AM",
    pm: "PM",
    midnight: "middernacht",
    noon: "het middaguur",
    morning: "'s ochtends",
    afternoon: "'s middags",
    evening: "'s avonds",
    night: "'s nachts"
  },
  wide: {
    am: "AM",
    pm: "PM",
    midnight: "middernacht",
    noon: "het middaguur",
    morning: "'s ochtends",
    afternoon: "'s middags",
    evening: "'s avonds",
    night: "'s nachts"
  }
};
var ordinalNumber8 = (dirtyNumber, _options) => {
  const number = Number(dirtyNumber);
  return number + "e";
};
var localize8 = {
  ordinalNumber: ordinalNumber8,
  era: buildLocalizeFn({
    values: eraValues8,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues8,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues8,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues8,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues8,
    defaultWidth: "wide"
  })
};

// node_modules/date-fns/locale/nl/_lib/match.js
var matchOrdinalNumberPattern8 = /^(\d+)e?/i;
var parseOrdinalNumberPattern8 = /\d+/i;
var matchEraPatterns8 = {
  narrow: /^([vn]\.? ?C\.?)/,
  abbreviated: /^([vn]\. ?Chr\.?)/,
  wide: /^((voor|na) Christus)/
};
var parseEraPatterns8 = {
  any: [/^v/, /^n/]
};
var matchQuarterPatterns8 = {
  narrow: /^[1234]/i,
  abbreviated: /^K[1234]/i,
  wide: /^[1234]e kwartaal/i
};
var parseQuarterPatterns8 = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns8 = {
  narrow: /^[jfmasond]/i,
  abbreviated: /^(jan.|feb.|mrt.|apr.|mei|jun.|jul.|aug.|sep.|okt.|nov.|dec.)/i,
  wide: /^(januari|februari|maart|april|mei|juni|juli|augustus|september|oktober|november|december)/i
};
var parseMonthPatterns8 = {
  narrow: [/^j/i, /^f/i, /^m/i, /^a/i, /^m/i, /^j/i, /^j/i, /^a/i, /^s/i, /^o/i, /^n/i, /^d/i],
  any: [/^jan/i, /^feb/i, /^m(r|a)/i, /^apr/i, /^mei/i, /^jun/i, /^jul/i, /^aug/i, /^sep/i, /^okt/i, /^nov/i, /^dec/i]
};
var matchDayPatterns8 = {
  narrow: /^[zmdwv]/i,
  short: /^(zo|ma|di|wo|do|vr|za)/i,
  abbreviated: /^(zon|maa|din|woe|don|vri|zat)/i,
  wide: /^(zondag|maandag|dinsdag|woensdag|donderdag|vrijdag|zaterdag)/i
};
var parseDayPatterns8 = {
  narrow: [/^z/i, /^m/i, /^d/i, /^w/i, /^d/i, /^v/i, /^z/i],
  any: [/^zo/i, /^ma/i, /^di/i, /^wo/i, /^do/i, /^vr/i, /^za/i]
};
var matchDayPeriodPatterns8 = {
  any: /^(am|pm|middernacht|het middaguur|'s (ochtends|middags|avonds|nachts))/i
};
var parseDayPeriodPatterns8 = {
  any: {
    am: /^am/i,
    pm: /^pm/i,
    midnight: /^middernacht/i,
    noon: /^het middaguur/i,
    morning: /ochtend/i,
    afternoon: /middag/i,
    evening: /avond/i,
    night: /nacht/i
  }
};
var match8 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern8,
    parsePattern: parseOrdinalNumberPattern8,
    valueCallback: (value) => parseInt(value, 10)
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns8,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns8,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns8,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns8,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns8,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns8,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns8,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns8,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns8,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns8,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/nl.js
var nl = {
  code: "nl",
  formatDistance: formatDistance8,
  formatLong: formatLong8,
  formatRelative: formatRelative8,
  localize: localize8,
  match: match8,
  options: {
    weekStartsOn: 1,
    firstWeekContainsDate: 4
  }
};

// node_modules/date-fns/locale/pt/_lib/formatDistance.js
var formatDistanceLocale9 = {
  lessThanXSeconds: {
    one: "menos de um segundo",
    other: "menos de {{count}} segundos"
  },
  xSeconds: {
    one: "1 segundo",
    other: "{{count}} segundos"
  },
  halfAMinute: "meio minuto",
  lessThanXMinutes: {
    one: "menos de um minuto",
    other: "menos de {{count}} minutos"
  },
  xMinutes: {
    one: "1 minuto",
    other: "{{count}} minutos"
  },
  aboutXHours: {
    one: "aproximadamente 1 hora",
    other: "aproximadamente {{count}} horas"
  },
  xHours: {
    one: "1 hora",
    other: "{{count}} horas"
  },
  xDays: {
    one: "1 dia",
    other: "{{count}} dias"
  },
  aboutXWeeks: {
    one: "aproximadamente 1 semana",
    other: "aproximadamente {{count}} semanas"
  },
  xWeeks: {
    one: "1 semana",
    other: "{{count}} semanas"
  },
  aboutXMonths: {
    one: "aproximadamente 1 m\xEAs",
    other: "aproximadamente {{count}} meses"
  },
  xMonths: {
    one: "1 m\xEAs",
    other: "{{count}} meses"
  },
  aboutXYears: {
    one: "aproximadamente 1 ano",
    other: "aproximadamente {{count}} anos"
  },
  xYears: {
    one: "1 ano",
    other: "{{count}} anos"
  },
  overXYears: {
    one: "mais de 1 ano",
    other: "mais de {{count}} anos"
  },
  almostXYears: {
    one: "quase 1 ano",
    other: "quase {{count}} anos"
  }
};
var formatDistance9 = (token, count, options) => {
  let result;
  const tokenValue = formatDistanceLocale9[token];
  if (typeof tokenValue === "string") {
    result = tokenValue;
  } else if (count === 1) {
    result = tokenValue.one;
  } else {
    result = tokenValue.other.replace("{{count}}", String(count));
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return "daqui a " + result;
    } else {
      return "h\xE1 " + result;
    }
  }
  return result;
};

// node_modules/date-fns/locale/pt/_lib/formatLong.js
var dateFormats9 = {
  full: "EEEE, d 'de' MMMM 'de' y",
  long: "d 'de' MMMM 'de' y",
  medium: "d 'de' MMM 'de' y",
  short: "dd/MM/y"
};
var timeFormats9 = {
  full: "HH:mm:ss zzzz",
  long: "HH:mm:ss z",
  medium: "HH:mm:ss",
  short: "HH:mm"
};
var dateTimeFormats9 = {
  full: "{{date}} '\xE0s' {{time}}",
  long: "{{date}} '\xE0s' {{time}}",
  medium: "{{date}}, {{time}}",
  short: "{{date}}, {{time}}"
};
var formatLong9 = {
  date: buildFormatLongFn({
    formats: dateFormats9,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats9,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats9,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/pt/_lib/formatRelative.js
var formatRelativeLocale9 = {
  lastWeek: (date) => {
    const weekday = date.getDay();
    const last = weekday === 0 || weekday === 6 ? "\xFAltimo" : "\xFAltima";
    return "'" + last + "' eeee '\xE0s' p";
  },
  yesterday: "'ontem \xE0s' p",
  today: "'hoje \xE0s' p",
  tomorrow: "'amanh\xE3 \xE0s' p",
  nextWeek: "eeee '\xE0s' p",
  other: "P"
};
var formatRelative9 = (token, date, _baseDate, _options) => {
  const format = formatRelativeLocale9[token];
  if (typeof format === "function") {
    return format(date);
  }
  return format;
};

// node_modules/date-fns/locale/pt/_lib/localize.js
var eraValues9 = {
  narrow: ["aC", "dC"],
  abbreviated: ["a.C.", "d.C."],
  wide: ["antes de Cristo", "depois de Cristo"]
};
var quarterValues9 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["T1", "T2", "T3", "T4"],
  wide: ["1\xBA trimestre", "2\xBA trimestre", "3\xBA trimestre", "4\xBA trimestre"]
};
var monthValues9 = {
  narrow: ["j", "f", "m", "a", "m", "j", "j", "a", "s", "o", "n", "d"],
  abbreviated: ["jan", "fev", "mar", "abr", "mai", "jun", "jul", "ago", "set", "out", "nov", "dez"],
  wide: ["janeiro", "fevereiro", "mar\xE7o", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"]
};
var dayValues9 = {
  narrow: ["d", "s", "t", "q", "q", "s", "s"],
  short: ["dom", "seg", "ter", "qua", "qui", "sex", "s\xE1b"],
  abbreviated: ["dom", "seg", "ter", "qua", "qui", "sex", "s\xE1b"],
  wide: ["domingo", "segunda-feira", "ter\xE7a-feira", "quarta-feira", "quinta-feira", "sexta-feira", "s\xE1bado"]
};
var dayPeriodValues9 = {
  narrow: {
    am: "AM",
    pm: "PM",
    midnight: "meia-noite",
    noon: "meio-dia",
    morning: "manh\xE3",
    afternoon: "tarde",
    evening: "noite",
    night: "madrugada"
  },
  abbreviated: {
    am: "AM",
    pm: "PM",
    midnight: "meia-noite",
    noon: "meio-dia",
    morning: "manh\xE3",
    afternoon: "tarde",
    evening: "noite",
    night: "madrugada"
  },
  wide: {
    am: "AM",
    pm: "PM",
    midnight: "meia-noite",
    noon: "meio-dia",
    morning: "manh\xE3",
    afternoon: "tarde",
    evening: "noite",
    night: "madrugada"
  }
};
var formattingDayPeriodValues7 = {
  narrow: {
    am: "AM",
    pm: "PM",
    midnight: "meia-noite",
    noon: "meio-dia",
    morning: "da manh\xE3",
    afternoon: "da tarde",
    evening: "da noite",
    night: "da madrugada"
  },
  abbreviated: {
    am: "AM",
    pm: "PM",
    midnight: "meia-noite",
    noon: "meio-dia",
    morning: "da manh\xE3",
    afternoon: "da tarde",
    evening: "da noite",
    night: "da madrugada"
  },
  wide: {
    am: "AM",
    pm: "PM",
    midnight: "meia-noite",
    noon: "meio-dia",
    morning: "da manh\xE3",
    afternoon: "da tarde",
    evening: "da noite",
    night: "da madrugada"
  }
};
var ordinalNumber9 = (dirtyNumber, _options) => {
  const number = Number(dirtyNumber);
  return number + "\xBA";
};
var localize9 = {
  ordinalNumber: ordinalNumber9,
  era: buildLocalizeFn({
    values: eraValues9,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues9,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues9,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues9,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues9,
    defaultWidth: "wide",
    formattingValues: formattingDayPeriodValues7,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/pt/_lib/match.js
var matchOrdinalNumberPattern9 = /^(\d+)(º|ª)?/i;
var parseOrdinalNumberPattern9 = /\d+/i;
var matchEraPatterns9 = {
  narrow: /^(ac|dc|a|d)/i,
  abbreviated: /^(a\.?\s?c\.?|a\.?\s?e\.?\s?c\.?|d\.?\s?c\.?|e\.?\s?c\.?)/i,
  wide: /^(antes de cristo|antes da era comum|depois de cristo|era comum)/i
};
var parseEraPatterns9 = {
  any: [/^ac/i, /^dc/i],
  wide: [/^(antes de cristo|antes da era comum)/i, /^(depois de cristo|era comum)/i]
};
var matchQuarterPatterns9 = {
  narrow: /^[1234]/i,
  abbreviated: /^T[1234]/i,
  wide: /^[1234](º|ª)? trimestre/i
};
var parseQuarterPatterns9 = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns9 = {
  narrow: /^[jfmasond]/i,
  abbreviated: /^(jan|fev|mar|abr|mai|jun|jul|ago|set|out|nov|dez)/i,
  wide: /^(janeiro|fevereiro|março|abril|maio|junho|julho|agosto|setembro|outubro|novembro|dezembro)/i
};
var parseMonthPatterns9 = {
  narrow: [/^j/i, /^f/i, /^m/i, /^a/i, /^m/i, /^j/i, /^j/i, /^a/i, /^s/i, /^o/i, /^n/i, /^d/i],
  any: [/^ja/i, /^f/i, /^mar/i, /^ab/i, /^mai/i, /^jun/i, /^jul/i, /^ag/i, /^s/i, /^o/i, /^n/i, /^d/i]
};
var matchDayPatterns9 = {
  narrow: /^[dstq]/i,
  short: /^(dom|seg|ter|qua|qui|sex|s[áa]b)/i,
  abbreviated: /^(dom|seg|ter|qua|qui|sex|s[áa]b)/i,
  wide: /^(domingo|segunda-?\s?feira|terça-?\s?feira|quarta-?\s?feira|quinta-?\s?feira|sexta-?\s?feira|s[áa]bado)/i
};
var parseDayPatterns9 = {
  narrow: [/^d/i, /^s/i, /^t/i, /^q/i, /^q/i, /^s/i, /^s/i],
  any: [/^d/i, /^seg/i, /^t/i, /^qua/i, /^qui/i, /^sex/i, /^s[áa]/i]
};
var matchDayPeriodPatterns9 = {
  narrow: /^(a|p|meia-?\s?noite|meio-?\s?dia|(da) (manh[ãa]|tarde|noite|madrugada))/i,
  any: /^([ap]\.?\s?m\.?|meia-?\s?noite|meio-?\s?dia|(da) (manh[ãa]|tarde|noite|madrugada))/i
};
var parseDayPeriodPatterns9 = {
  any: {
    am: /^a/i,
    pm: /^p/i,
    midnight: /^meia/i,
    noon: /^meio/i,
    morning: /manh[ãa]/i,
    afternoon: /tarde/i,
    evening: /noite/i,
    night: /madrugada/i
  }
};
var match9 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern9,
    parsePattern: parseOrdinalNumberPattern9,
    valueCallback: (value) => parseInt(value, 10)
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns9,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns9,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns9,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns9,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns9,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns9,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns9,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns9,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns9,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns9,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/pt.js
var pt = {
  code: "pt",
  formatDistance: formatDistance9,
  formatLong: formatLong9,
  formatRelative: formatRelative9,
  localize: localize9,
  match: match9,
  options: {
    weekStartsOn: 1,
    firstWeekContainsDate: 4
  }
};

// node_modules/date-fns/locale/ru/_lib/formatDistance.js
function declension(scheme, count) {
  if (scheme.one !== void 0 && count === 1) {
    return scheme.one;
  }
  const rem10 = count % 10;
  const rem100 = count % 100;
  if (rem10 === 1 && rem100 !== 11) {
    return scheme.singularNominative.replace("{{count}}", String(count));
  } else if (rem10 >= 2 && rem10 <= 4 && (rem100 < 10 || rem100 > 20)) {
    return scheme.singularGenitive.replace("{{count}}", String(count));
  } else {
    return scheme.pluralGenitive.replace("{{count}}", String(count));
  }
}
function buildLocalizeTokenFn(scheme) {
  return (count, options) => {
    if (options?.addSuffix) {
      if (options.comparison && options.comparison > 0) {
        if (scheme.future) {
          return declension(scheme.future, count);
        } else {
          return "\u0447\u0435\u0440\u0435\u0437 " + declension(scheme.regular, count);
        }
      } else {
        if (scheme.past) {
          return declension(scheme.past, count);
        } else {
          return declension(scheme.regular, count) + " \u043D\u0430\u0437\u0430\u0434";
        }
      }
    } else {
      return declension(scheme.regular, count);
    }
  };
}
var formatDistanceLocale10 = {
  lessThanXSeconds: buildLocalizeTokenFn({
    regular: {
      one: "\u043C\u0435\u043D\u044C\u0448\u0435 \u0441\u0435\u043A\u0443\u043D\u0434\u044B",
      singularNominative: "\u043C\u0435\u043D\u044C\u0448\u0435 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u044B",
      singularGenitive: "\u043C\u0435\u043D\u044C\u0448\u0435 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434",
      pluralGenitive: "\u043C\u0435\u043D\u044C\u0448\u0435 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434"
    },
    future: {
      one: "\u043C\u0435\u043D\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 \u0441\u0435\u043A\u0443\u043D\u0434\u0443",
      singularNominative: "\u043C\u0435\u043D\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0443",
      singularGenitive: "\u043C\u0435\u043D\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u044B",
      pluralGenitive: "\u043C\u0435\u043D\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434"
    }
  }),
  xSeconds: buildLocalizeTokenFn({
    regular: {
      singularNominative: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0430",
      singularGenitive: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u044B",
      pluralGenitive: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434"
    },
    past: {
      singularNominative: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0443 \u043D\u0430\u0437\u0430\u0434",
      singularGenitive: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u044B \u043D\u0430\u0437\u0430\u0434",
      pluralGenitive: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434 \u043D\u0430\u0437\u0430\u0434"
    },
    future: {
      singularNominative: "\u0447\u0435\u0440\u0435\u0437 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0443",
      singularGenitive: "\u0447\u0435\u0440\u0435\u0437 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u044B",
      pluralGenitive: "\u0447\u0435\u0440\u0435\u0437 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434"
    }
  }),
  halfAMinute: (_count, options) => {
    if (options?.addSuffix) {
      if (options.comparison && options.comparison > 0) {
        return "\u0447\u0435\u0440\u0435\u0437 \u043F\u043E\u043B\u043C\u0438\u043D\u0443\u0442\u044B";
      } else {
        return "\u043F\u043E\u043B\u043C\u0438\u043D\u0443\u0442\u044B \u043D\u0430\u0437\u0430\u0434";
      }
    }
    return "\u043F\u043E\u043B\u043C\u0438\u043D\u0443\u0442\u044B";
  },
  lessThanXMinutes: buildLocalizeTokenFn({
    regular: {
      one: "\u043C\u0435\u043D\u044C\u0448\u0435 \u043C\u0438\u043D\u0443\u0442\u044B",
      singularNominative: "\u043C\u0435\u043D\u044C\u0448\u0435 {{count}} \u043C\u0438\u043D\u0443\u0442\u044B",
      singularGenitive: "\u043C\u0435\u043D\u044C\u0448\u0435 {{count}} \u043C\u0438\u043D\u0443\u0442",
      pluralGenitive: "\u043C\u0435\u043D\u044C\u0448\u0435 {{count}} \u043C\u0438\u043D\u0443\u0442"
    },
    future: {
      one: "\u043C\u0435\u043D\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 \u043C\u0438\u043D\u0443\u0442\u0443",
      singularNominative: "\u043C\u0435\u043D\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 {{count}} \u043C\u0438\u043D\u0443\u0442\u0443",
      singularGenitive: "\u043C\u0435\u043D\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 {{count}} \u043C\u0438\u043D\u0443\u0442\u044B",
      pluralGenitive: "\u043C\u0435\u043D\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 {{count}} \u043C\u0438\u043D\u0443\u0442"
    }
  }),
  xMinutes: buildLocalizeTokenFn({
    regular: {
      singularNominative: "{{count}} \u043C\u0438\u043D\u0443\u0442\u0430",
      singularGenitive: "{{count}} \u043C\u0438\u043D\u0443\u0442\u044B",
      pluralGenitive: "{{count}} \u043C\u0438\u043D\u0443\u0442"
    },
    past: {
      singularNominative: "{{count}} \u043C\u0438\u043D\u0443\u0442\u0443 \u043D\u0430\u0437\u0430\u0434",
      singularGenitive: "{{count}} \u043C\u0438\u043D\u0443\u0442\u044B \u043D\u0430\u0437\u0430\u0434",
      pluralGenitive: "{{count}} \u043C\u0438\u043D\u0443\u0442 \u043D\u0430\u0437\u0430\u0434"
    },
    future: {
      singularNominative: "\u0447\u0435\u0440\u0435\u0437 {{count}} \u043C\u0438\u043D\u0443\u0442\u0443",
      singularGenitive: "\u0447\u0435\u0440\u0435\u0437 {{count}} \u043C\u0438\u043D\u0443\u0442\u044B",
      pluralGenitive: "\u0447\u0435\u0440\u0435\u0437 {{count}} \u043C\u0438\u043D\u0443\u0442"
    }
  }),
  aboutXHours: buildLocalizeTokenFn({
    regular: {
      singularNominative: "\u043E\u043A\u043E\u043B\u043E {{count}} \u0447\u0430\u0441\u0430",
      singularGenitive: "\u043E\u043A\u043E\u043B\u043E {{count}} \u0447\u0430\u0441\u043E\u0432",
      pluralGenitive: "\u043E\u043A\u043E\u043B\u043E {{count}} \u0447\u0430\u0441\u043E\u0432"
    },
    future: {
      singularNominative: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u0447\u0430\u0441",
      singularGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u0447\u0430\u0441\u0430",
      pluralGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u0447\u0430\u0441\u043E\u0432"
    }
  }),
  xHours: buildLocalizeTokenFn({
    regular: {
      singularNominative: "{{count}} \u0447\u0430\u0441",
      singularGenitive: "{{count}} \u0447\u0430\u0441\u0430",
      pluralGenitive: "{{count}} \u0447\u0430\u0441\u043E\u0432"
    }
  }),
  xDays: buildLocalizeTokenFn({
    regular: {
      singularNominative: "{{count}} \u0434\u0435\u043D\u044C",
      singularGenitive: "{{count}} \u0434\u043D\u044F",
      pluralGenitive: "{{count}} \u0434\u043D\u0435\u0439"
    }
  }),
  aboutXWeeks: buildLocalizeTokenFn({
    regular: {
      singularNominative: "\u043E\u043A\u043E\u043B\u043E {{count}} \u043D\u0435\u0434\u0435\u043B\u0438",
      singularGenitive: "\u043E\u043A\u043E\u043B\u043E {{count}} \u043D\u0435\u0434\u0435\u043B\u044C",
      pluralGenitive: "\u043E\u043A\u043E\u043B\u043E {{count}} \u043D\u0435\u0434\u0435\u043B\u044C"
    },
    future: {
      singularNominative: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u043D\u0435\u0434\u0435\u043B\u044E",
      singularGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u043D\u0435\u0434\u0435\u043B\u0438",
      pluralGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u043D\u0435\u0434\u0435\u043B\u044C"
    }
  }),
  xWeeks: buildLocalizeTokenFn({
    regular: {
      singularNominative: "{{count}} \u043D\u0435\u0434\u0435\u043B\u044F",
      singularGenitive: "{{count}} \u043D\u0435\u0434\u0435\u043B\u0438",
      pluralGenitive: "{{count}} \u043D\u0435\u0434\u0435\u043B\u044C"
    }
  }),
  aboutXMonths: buildLocalizeTokenFn({
    regular: {
      singularNominative: "\u043E\u043A\u043E\u043B\u043E {{count}} \u043C\u0435\u0441\u044F\u0446\u0430",
      singularGenitive: "\u043E\u043A\u043E\u043B\u043E {{count}} \u043C\u0435\u0441\u044F\u0446\u0435\u0432",
      pluralGenitive: "\u043E\u043A\u043E\u043B\u043E {{count}} \u043C\u0435\u0441\u044F\u0446\u0435\u0432"
    },
    future: {
      singularNominative: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u043C\u0435\u0441\u044F\u0446",
      singularGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u043C\u0435\u0441\u044F\u0446\u0430",
      pluralGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u043C\u0435\u0441\u044F\u0446\u0435\u0432"
    }
  }),
  xMonths: buildLocalizeTokenFn({
    regular: {
      singularNominative: "{{count}} \u043C\u0435\u0441\u044F\u0446",
      singularGenitive: "{{count}} \u043C\u0435\u0441\u044F\u0446\u0430",
      pluralGenitive: "{{count}} \u043C\u0435\u0441\u044F\u0446\u0435\u0432"
    }
  }),
  aboutXYears: buildLocalizeTokenFn({
    regular: {
      singularNominative: "\u043E\u043A\u043E\u043B\u043E {{count}} \u0433\u043E\u0434\u0430",
      singularGenitive: "\u043E\u043A\u043E\u043B\u043E {{count}} \u043B\u0435\u0442",
      pluralGenitive: "\u043E\u043A\u043E\u043B\u043E {{count}} \u043B\u0435\u0442"
    },
    future: {
      singularNominative: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u0433\u043E\u0434",
      singularGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u0433\u043E\u0434\u0430",
      pluralGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u0438\u0442\u0435\u043B\u044C\u043D\u043E \u0447\u0435\u0440\u0435\u0437 {{count}} \u043B\u0435\u0442"
    }
  }),
  xYears: buildLocalizeTokenFn({
    regular: {
      singularNominative: "{{count}} \u0433\u043E\u0434",
      singularGenitive: "{{count}} \u0433\u043E\u0434\u0430",
      pluralGenitive: "{{count}} \u043B\u0435\u0442"
    }
  }),
  overXYears: buildLocalizeTokenFn({
    regular: {
      singularNominative: "\u0431\u043E\u043B\u044C\u0448\u0435 {{count}} \u0433\u043E\u0434\u0430",
      singularGenitive: "\u0431\u043E\u043B\u044C\u0448\u0435 {{count}} \u043B\u0435\u0442",
      pluralGenitive: "\u0431\u043E\u043B\u044C\u0448\u0435 {{count}} \u043B\u0435\u0442"
    },
    future: {
      singularNominative: "\u0431\u043E\u043B\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 {{count}} \u0433\u043E\u0434",
      singularGenitive: "\u0431\u043E\u043B\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 {{count}} \u0433\u043E\u0434\u0430",
      pluralGenitive: "\u0431\u043E\u043B\u044C\u0448\u0435, \u0447\u0435\u043C \u0447\u0435\u0440\u0435\u0437 {{count}} \u043B\u0435\u0442"
    }
  }),
  almostXYears: buildLocalizeTokenFn({
    regular: {
      singularNominative: "\u043F\u043E\u0447\u0442\u0438 {{count}} \u0433\u043E\u0434",
      singularGenitive: "\u043F\u043E\u0447\u0442\u0438 {{count}} \u0433\u043E\u0434\u0430",
      pluralGenitive: "\u043F\u043E\u0447\u0442\u0438 {{count}} \u043B\u0435\u0442"
    },
    future: {
      singularNominative: "\u043F\u043E\u0447\u0442\u0438 \u0447\u0435\u0440\u0435\u0437 {{count}} \u0433\u043E\u0434",
      singularGenitive: "\u043F\u043E\u0447\u0442\u0438 \u0447\u0435\u0440\u0435\u0437 {{count}} \u0433\u043E\u0434\u0430",
      pluralGenitive: "\u043F\u043E\u0447\u0442\u0438 \u0447\u0435\u0440\u0435\u0437 {{count}} \u043B\u0435\u0442"
    }
  })
};
var formatDistance10 = (token, count, options) => {
  return formatDistanceLocale10[token](count, options);
};

// node_modules/date-fns/locale/ru/_lib/formatLong.js
var dateFormats10 = {
  full: "EEEE, d MMMM y '\u0433.'",
  long: "d MMMM y '\u0433.'",
  medium: "d MMM y '\u0433.'",
  short: "dd.MM.y"
};
var timeFormats10 = {
  full: "H:mm:ss zzzz",
  long: "H:mm:ss z",
  medium: "H:mm:ss",
  short: "H:mm"
};
var dateTimeFormats10 = {
  any: "{{date}}, {{time}}"
};
var formatLong10 = {
  date: buildFormatLongFn({
    formats: dateFormats10,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats10,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats10,
    defaultWidth: "any"
  })
};

// node_modules/date-fns/constants.js
var daysInYear = 365.2425;
var maxTime = Math.pow(10, 8) * 24 * 60 * 60 * 1e3;
var minTime = -maxTime;
var millisecondsInWeek = 6048e5;
var millisecondsInDay = 864e5;
var minutesInMonth = 43200;
var minutesInDay = 1440;
var secondsInHour = 3600;
var secondsInDay = secondsInHour * 24;
var secondsInWeek = secondsInDay * 7;
var secondsInYear = secondsInDay * daysInYear;
var secondsInMonth = secondsInYear / 12;
var secondsInQuarter = secondsInMonth * 3;
var constructFromSymbol = Symbol.for("constructDateFrom");

// node_modules/date-fns/constructFrom.js
function constructFrom(date, value) {
  if (typeof date === "function") return date(value);
  if (date && typeof date === "object" && constructFromSymbol in date) return date[constructFromSymbol](value);
  if (date instanceof Date) return new date.constructor(value);
  return new Date(value);
}

// node_modules/date-fns/_lib/normalizeDates.js
function normalizeDates(context, ...dates) {
  const normalize = constructFrom.bind(null, context || dates.find((date) => typeof date === "object"));
  return dates.map(normalize);
}

// node_modules/date-fns/_lib/defaultOptions.js
var defaultOptions = {};
function getDefaultOptions() {
  return defaultOptions;
}

// node_modules/date-fns/toDate.js
function toDate(argument, context) {
  return constructFrom(context || argument, argument);
}

// node_modules/date-fns/startOfWeek.js
function startOfWeek(date, options) {
  const defaultOptions2 = getDefaultOptions();
  const weekStartsOn = options?.weekStartsOn ?? options?.locale?.options?.weekStartsOn ?? defaultOptions2.weekStartsOn ?? defaultOptions2.locale?.options?.weekStartsOn ?? 0;
  const _date = toDate(date, options?.in);
  const day = _date.getDay();
  const diff = (day < weekStartsOn ? 7 : 0) + day - weekStartsOn;
  _date.setDate(_date.getDate() - diff);
  _date.setHours(0, 0, 0, 0);
  return _date;
}

// node_modules/date-fns/isSameWeek.js
function isSameWeek(laterDate, earlierDate, options) {
  const [laterDate_, earlierDate_] = normalizeDates(options?.in, laterDate, earlierDate);
  return +startOfWeek(laterDate_, options) === +startOfWeek(earlierDate_, options);
}

// node_modules/date-fns/locale/ru/_lib/formatRelative.js
var accusativeWeekdays = ["\u0432\u043E\u0441\u043A\u0440\u0435\u0441\u0435\u043D\u044C\u0435", "\u043F\u043E\u043D\u0435\u0434\u0435\u043B\u044C\u043D\u0438\u043A", "\u0432\u0442\u043E\u0440\u043D\u0438\u043A", "\u0441\u0440\u0435\u0434\u0443", "\u0447\u0435\u0442\u0432\u0435\u0440\u0433", "\u043F\u044F\u0442\u043D\u0438\u0446\u0443", "\u0441\u0443\u0431\u0431\u043E\u0442\u0443"];
function lastWeek(day) {
  const weekday = accusativeWeekdays[day];
  switch (day) {
    case 0:
      return "'\u0432 \u043F\u0440\u043E\u0448\u043B\u043E\u0435 " + weekday + " \u0432' p";
    case 1:
    case 2:
    case 4:
      return "'\u0432 \u043F\u0440\u043E\u0448\u043B\u044B\u0439 " + weekday + " \u0432' p";
    case 3:
    case 5:
    case 6:
      return "'\u0432 \u043F\u0440\u043E\u0448\u043B\u0443\u044E " + weekday + " \u0432' p";
  }
}
function thisWeek(day) {
  const weekday = accusativeWeekdays[day];
  if (day === 2) {
    return "'\u0432\u043E " + weekday + " \u0432' p";
  } else {
    return "'\u0432 " + weekday + " \u0432' p";
  }
}
function nextWeek(day) {
  const weekday = accusativeWeekdays[day];
  switch (day) {
    case 0:
      return "'\u0432 \u0441\u043B\u0435\u0434\u0443\u044E\u0449\u0435\u0435 " + weekday + " \u0432' p";
    case 1:
    case 2:
    case 4:
      return "'\u0432 \u0441\u043B\u0435\u0434\u0443\u044E\u0449\u0438\u0439 " + weekday + " \u0432' p";
    case 3:
    case 5:
    case 6:
      return "'\u0432 \u0441\u043B\u0435\u0434\u0443\u044E\u0449\u0443\u044E " + weekday + " \u0432' p";
  }
}
var formatRelativeLocale10 = {
  lastWeek: (date, baseDate, options) => {
    const day = date.getDay();
    if (isSameWeek(date, baseDate, options)) {
      return thisWeek(day);
    } else {
      return lastWeek(day);
    }
  },
  yesterday: "'\u0432\u0447\u0435\u0440\u0430 \u0432' p",
  today: "'\u0441\u0435\u0433\u043E\u0434\u043D\u044F \u0432' p",
  tomorrow: "'\u0437\u0430\u0432\u0442\u0440\u0430 \u0432' p",
  nextWeek: (date, baseDate, options) => {
    const day = date.getDay();
    if (isSameWeek(date, baseDate, options)) {
      return thisWeek(day);
    } else {
      return nextWeek(day);
    }
  },
  other: "P"
};
var formatRelative10 = (token, date, baseDate, options) => {
  const format = formatRelativeLocale10[token];
  if (typeof format === "function") {
    return format(date, baseDate, options);
  }
  return format;
};

// node_modules/date-fns/locale/ru/_lib/localize.js
var eraValues10 = {
  narrow: ["\u0434\u043E \u043D.\u044D.", "\u043D.\u044D."],
  abbreviated: ["\u0434\u043E \u043D. \u044D.", "\u043D. \u044D."],
  wide: ["\u0434\u043E \u043D\u0430\u0448\u0435\u0439 \u044D\u0440\u044B", "\u043D\u0430\u0448\u0435\u0439 \u044D\u0440\u044B"]
};
var quarterValues10 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["1-\u0439 \u043A\u0432.", "2-\u0439 \u043A\u0432.", "3-\u0439 \u043A\u0432.", "4-\u0439 \u043A\u0432."],
  wide: ["1-\u0439 \u043A\u0432\u0430\u0440\u0442\u0430\u043B", "2-\u0439 \u043A\u0432\u0430\u0440\u0442\u0430\u043B", "3-\u0439 \u043A\u0432\u0430\u0440\u0442\u0430\u043B", "4-\u0439 \u043A\u0432\u0430\u0440\u0442\u0430\u043B"]
};
var monthValues10 = {
  narrow: ["\u042F", "\u0424", "\u041C", "\u0410", "\u041C", "\u0418", "\u0418", "\u0410", "\u0421", "\u041E", "\u041D", "\u0414"],
  abbreviated: ["\u044F\u043D\u0432.", "\u0444\u0435\u0432.", "\u043C\u0430\u0440\u0442", "\u0430\u043F\u0440.", "\u043C\u0430\u0439", "\u0438\u044E\u043D\u044C", "\u0438\u044E\u043B\u044C", "\u0430\u0432\u0433.", "\u0441\u0435\u043D\u0442.", "\u043E\u043A\u0442.", "\u043D\u043E\u044F\u0431.", "\u0434\u0435\u043A."],
  wide: ["\u044F\u043D\u0432\u0430\u0440\u044C", "\u0444\u0435\u0432\u0440\u0430\u043B\u044C", "\u043C\u0430\u0440\u0442", "\u0430\u043F\u0440\u0435\u043B\u044C", "\u043C\u0430\u0439", "\u0438\u044E\u043D\u044C", "\u0438\u044E\u043B\u044C", "\u0430\u0432\u0433\u0443\u0441\u0442", "\u0441\u0435\u043D\u0442\u044F\u0431\u0440\u044C", "\u043E\u043A\u0442\u044F\u0431\u0440\u044C", "\u043D\u043E\u044F\u0431\u0440\u044C", "\u0434\u0435\u043A\u0430\u0431\u0440\u044C"]
};
var formattingMonthValues2 = {
  narrow: ["\u042F", "\u0424", "\u041C", "\u0410", "\u041C", "\u0418", "\u0418", "\u0410", "\u0421", "\u041E", "\u041D", "\u0414"],
  abbreviated: ["\u044F\u043D\u0432.", "\u0444\u0435\u0432.", "\u043C\u0430\u0440.", "\u0430\u043F\u0440.", "\u043C\u0430\u044F", "\u0438\u044E\u043D.", "\u0438\u044E\u043B.", "\u0430\u0432\u0433.", "\u0441\u0435\u043D\u0442.", "\u043E\u043A\u0442.", "\u043D\u043E\u044F\u0431.", "\u0434\u0435\u043A."],
  wide: ["\u044F\u043D\u0432\u0430\u0440\u044F", "\u0444\u0435\u0432\u0440\u0430\u043B\u044F", "\u043C\u0430\u0440\u0442\u0430", "\u0430\u043F\u0440\u0435\u043B\u044F", "\u043C\u0430\u044F", "\u0438\u044E\u043D\u044F", "\u0438\u044E\u043B\u044F", "\u0430\u0432\u0433\u0443\u0441\u0442\u0430", "\u0441\u0435\u043D\u0442\u044F\u0431\u0440\u044F", "\u043E\u043A\u0442\u044F\u0431\u0440\u044F", "\u043D\u043E\u044F\u0431\u0440\u044F", "\u0434\u0435\u043A\u0430\u0431\u0440\u044F"]
};
var dayValues10 = {
  narrow: ["\u0412", "\u041F", "\u0412", "\u0421", "\u0427", "\u041F", "\u0421"],
  short: ["\u0432\u0441", "\u043F\u043D", "\u0432\u0442", "\u0441\u0440", "\u0447\u0442", "\u043F\u0442", "\u0441\u0431"],
  abbreviated: ["\u0432\u0441\u043A", "\u043F\u043D\u0434", "\u0432\u0442\u0440", "\u0441\u0440\u0434", "\u0447\u0442\u0432", "\u043F\u0442\u043D", "\u0441\u0443\u0431"],
  wide: ["\u0432\u043E\u0441\u043A\u0440\u0435\u0441\u0435\u043D\u044C\u0435", "\u043F\u043E\u043D\u0435\u0434\u0435\u043B\u044C\u043D\u0438\u043A", "\u0432\u0442\u043E\u0440\u043D\u0438\u043A", "\u0441\u0440\u0435\u0434\u0430", "\u0447\u0435\u0442\u0432\u0435\u0440\u0433", "\u043F\u044F\u0442\u043D\u0438\u0446\u0430", "\u0441\u0443\u0431\u0431\u043E\u0442\u0430"]
};
var dayPeriodValues10 = {
  narrow: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u043E\u043B\u043D.",
    noon: "\u043F\u043E\u043B\u0434.",
    morning: "\u0443\u0442\u0440\u043E",
    afternoon: "\u0434\u0435\u043D\u044C",
    evening: "\u0432\u0435\u0447.",
    night: "\u043D\u043E\u0447\u044C"
  },
  abbreviated: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u043E\u043B\u043D.",
    noon: "\u043F\u043E\u043B\u0434.",
    morning: "\u0443\u0442\u0440\u043E",
    afternoon: "\u0434\u0435\u043D\u044C",
    evening: "\u0432\u0435\u0447.",
    night: "\u043D\u043E\u0447\u044C"
  },
  wide: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u043E\u043B\u043D\u043E\u0447\u044C",
    noon: "\u043F\u043E\u043B\u0434\u0435\u043D\u044C",
    morning: "\u0443\u0442\u0440\u043E",
    afternoon: "\u0434\u0435\u043D\u044C",
    evening: "\u0432\u0435\u0447\u0435\u0440",
    night: "\u043D\u043E\u0447\u044C"
  }
};
var formattingDayPeriodValues8 = {
  narrow: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u043E\u043B\u043D.",
    noon: "\u043F\u043E\u043B\u0434.",
    morning: "\u0443\u0442\u0440\u0430",
    afternoon: "\u0434\u043D\u044F",
    evening: "\u0432\u0435\u0447.",
    night: "\u043D\u043E\u0447\u0438"
  },
  abbreviated: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u043E\u043B\u043D.",
    noon: "\u043F\u043E\u043B\u0434.",
    morning: "\u0443\u0442\u0440\u0430",
    afternoon: "\u0434\u043D\u044F",
    evening: "\u0432\u0435\u0447.",
    night: "\u043D\u043E\u0447\u0438"
  },
  wide: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u043E\u043B\u043D\u043E\u0447\u044C",
    noon: "\u043F\u043E\u043B\u0434\u0435\u043D\u044C",
    morning: "\u0443\u0442\u0440\u0430",
    afternoon: "\u0434\u043D\u044F",
    evening: "\u0432\u0435\u0447\u0435\u0440\u0430",
    night: "\u043D\u043E\u0447\u0438"
  }
};
var ordinalNumber10 = (dirtyNumber, options) => {
  const number = Number(dirtyNumber);
  const unit = options?.unit;
  let suffix;
  if (unit === "date") {
    suffix = "-\u0435";
  } else if (unit === "week" || unit === "minute" || unit === "second") {
    suffix = "-\u044F";
  } else {
    suffix = "-\u0439";
  }
  return number + suffix;
};
var localize10 = {
  ordinalNumber: ordinalNumber10,
  era: buildLocalizeFn({
    values: eraValues10,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues10,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues10,
    defaultWidth: "wide",
    formattingValues: formattingMonthValues2,
    defaultFormattingWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues10,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues10,
    defaultWidth: "any",
    formattingValues: formattingDayPeriodValues8,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/ru/_lib/match.js
var matchOrdinalNumberPattern10 = /^(\d+)(-?(е|я|й|ое|ье|ая|ья|ый|ой|ий|ый))?/i;
var parseOrdinalNumberPattern10 = /\d+/i;
var matchEraPatterns10 = {
  narrow: /^((до )?н\.?\s?э\.?)/i,
  abbreviated: /^((до )?н\.?\s?э\.?)/i,
  wide: /^(до нашей эры|нашей эры|наша эра)/i
};
var parseEraPatterns10 = {
  any: [/^д/i, /^н/i]
};
var matchQuarterPatterns10 = {
  narrow: /^[1234]/i,
  abbreviated: /^[1234](-?[ыои]?й?)? кв.?/i,
  wide: /^[1234](-?[ыои]?й?)? квартал/i
};
var parseQuarterPatterns10 = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns10 = {
  narrow: /^[яфмаисонд]/i,
  abbreviated: /^(янв|фев|март?|апр|ма[йя]|июн[ья]?|июл[ья]?|авг|сент?|окт|нояб?|дек)\.?/i,
  wide: /^(январ[ья]|феврал[ья]|марта?|апрел[ья]|ма[йя]|июн[ья]|июл[ья]|августа?|сентябр[ья]|октябр[ья]|октябр[ья]|ноябр[ья]|декабр[ья])/i
};
var parseMonthPatterns10 = {
  narrow: [/^я/i, /^ф/i, /^м/i, /^а/i, /^м/i, /^и/i, /^и/i, /^а/i, /^с/i, /^о/i, /^н/i, /^я/i],
  any: [/^я/i, /^ф/i, /^мар/i, /^ап/i, /^ма[йя]/i, /^июн/i, /^июл/i, /^ав/i, /^с/i, /^о/i, /^н/i, /^д/i]
};
var matchDayPatterns10 = {
  narrow: /^[впсч]/i,
  short: /^(вс|во|пн|по|вт|ср|чт|че|пт|пя|сб|су)\.?/i,
  abbreviated: /^(вск|вос|пнд|пон|втр|вто|срд|сре|чтв|чет|птн|пят|суб).?/i,
  wide: /^(воскресень[ея]|понедельника?|вторника?|сред[аы]|четверга?|пятниц[аы]|суббот[аы])/i
};
var parseDayPatterns10 = {
  narrow: [/^в/i, /^п/i, /^в/i, /^с/i, /^ч/i, /^п/i, /^с/i],
  any: [/^в[ос]/i, /^п[он]/i, /^в/i, /^ср/i, /^ч/i, /^п[ят]/i, /^с[уб]/i]
};
var matchDayPeriodPatterns10 = {
  narrow: /^([дп]п|полн\.?|полд\.?|утр[оа]|день|дня|веч\.?|ноч[ьи])/i,
  abbreviated: /^([дп]п|полн\.?|полд\.?|утр[оа]|день|дня|веч\.?|ноч[ьи])/i,
  wide: /^([дп]п|полночь|полдень|утр[оа]|день|дня|вечера?|ноч[ьи])/i
};
var parseDayPeriodPatterns10 = {
  any: {
    am: /^дп/i,
    pm: /^пп/i,
    midnight: /^полн/i,
    noon: /^полд/i,
    morning: /^у/i,
    afternoon: /^д[ен]/i,
    evening: /^в/i,
    night: /^н/i
  }
};
var match10 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern10,
    parsePattern: parseOrdinalNumberPattern10,
    valueCallback: (value) => parseInt(value, 10)
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns10,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns10,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns10,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns10,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns10,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns10,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns10,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns10,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns10,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPeriodPatterns10,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/ru.js
var ru = {
  code: "ru",
  formatDistance: formatDistance10,
  formatLong: formatLong10,
  formatRelative: formatRelative10,
  localize: localize10,
  match: match10,
  options: {
    weekStartsOn: 1,
    firstWeekContainsDate: 1
  }
};

// node_modules/date-fns/locale/tr/_lib/formatDistance.js
var formatDistanceLocale11 = {
  lessThanXSeconds: {
    one: "bir saniyeden az",
    other: "{{count}} saniyeden az"
  },
  xSeconds: {
    one: "1 saniye",
    other: "{{count}} saniye"
  },
  halfAMinute: "yar\u0131m dakika",
  lessThanXMinutes: {
    one: "bir dakikadan az",
    other: "{{count}} dakikadan az"
  },
  xMinutes: {
    one: "1 dakika",
    other: "{{count}} dakika"
  },
  aboutXHours: {
    one: "yakla\u015F\u0131k 1 saat",
    other: "yakla\u015F\u0131k {{count}} saat"
  },
  xHours: {
    one: "1 saat",
    other: "{{count}} saat"
  },
  xDays: {
    one: "1 g\xFCn",
    other: "{{count}} g\xFCn"
  },
  aboutXWeeks: {
    one: "yakla\u015F\u0131k 1 hafta",
    other: "yakla\u015F\u0131k {{count}} hafta"
  },
  xWeeks: {
    one: "1 hafta",
    other: "{{count}} hafta"
  },
  aboutXMonths: {
    one: "yakla\u015F\u0131k 1 ay",
    other: "yakla\u015F\u0131k {{count}} ay"
  },
  xMonths: {
    one: "1 ay",
    other: "{{count}} ay"
  },
  aboutXYears: {
    one: "yakla\u015F\u0131k 1 y\u0131l",
    other: "yakla\u015F\u0131k {{count}} y\u0131l"
  },
  xYears: {
    one: "1 y\u0131l",
    other: "{{count}} y\u0131l"
  },
  overXYears: {
    one: "1 y\u0131ldan fazla",
    other: "{{count}} y\u0131ldan fazla"
  },
  almostXYears: {
    one: "neredeyse 1 y\u0131l",
    other: "neredeyse {{count}} y\u0131l"
  }
};
var formatDistance11 = (token, count, options) => {
  let result;
  const tokenValue = formatDistanceLocale11[token];
  if (typeof tokenValue === "string") {
    result = tokenValue;
  } else if (count === 1) {
    result = tokenValue.one;
  } else {
    result = tokenValue.other.replace("{{count}}", count.toString());
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return result + " sonra";
    } else {
      return result + " \xF6nce";
    }
  }
  return result;
};

// node_modules/date-fns/locale/tr/_lib/formatLong.js
var dateFormats11 = {
  full: "d MMMM y EEEE",
  long: "d MMMM y",
  medium: "d MMM y",
  short: "dd.MM.yyyy"
};
var timeFormats11 = {
  full: "HH:mm:ss zzzz",
  long: "HH:mm:ss z",
  medium: "HH:mm:ss",
  short: "HH:mm"
};
var dateTimeFormats11 = {
  full: "{{date}} 'saat' {{time}}",
  long: "{{date}} 'saat' {{time}}",
  medium: "{{date}}, {{time}}",
  short: "{{date}}, {{time}}"
};
var formatLong11 = {
  date: buildFormatLongFn({
    formats: dateFormats11,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats11,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats11,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/tr/_lib/formatRelative.js
var formatRelativeLocale11 = {
  lastWeek: "'ge\xE7en hafta' eeee 'saat' p",
  yesterday: "'d\xFCn saat' p",
  today: "'bug\xFCn saat' p",
  tomorrow: "'yar\u0131n saat' p",
  nextWeek: "eeee 'saat' p",
  other: "P"
};
var formatRelative11 = (token, _date, _baseDate, _options) => formatRelativeLocale11[token];

// node_modules/date-fns/locale/tr/_lib/localize.js
var eraValues11 = {
  narrow: ["M\xD6", "MS"],
  abbreviated: ["M\xD6", "MS"],
  wide: ["Milattan \xD6nce", "Milattan Sonra"]
};
var quarterValues11 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["1\xC7", "2\xC7", "3\xC7", "4\xC7"],
  wide: ["\u0130lk \xE7eyrek", "\u0130kinci \xC7eyrek", "\xDC\xE7\xFCnc\xFC \xE7eyrek", "Son \xE7eyrek"]
};
var monthValues11 = {
  narrow: ["O", "\u015E", "M", "N", "M", "H", "T", "A", "E", "E", "K", "A"],
  abbreviated: ["Oca", "\u015Eub", "Mar", "Nis", "May", "Haz", "Tem", "A\u011Fu", "Eyl", "Eki", "Kas", "Ara"],
  wide: ["Ocak", "\u015Eubat", "Mart", "Nisan", "May\u0131s", "Haziran", "Temmuz", "A\u011Fustos", "Eyl\xFCl", "Ekim", "Kas\u0131m", "Aral\u0131k"]
};
var dayValues11 = {
  narrow: ["P", "P", "S", "\xC7", "P", "C", "C"],
  short: ["Pz", "Pt", "Sa", "\xC7a", "Pe", "Cu", "Ct"],
  abbreviated: ["Paz", "Pzt", "Sal", "\xC7ar", "Per", "Cum", "Cts"],
  wide: ["Pazar", "Pazartesi", "Sal\u0131", "\xC7ar\u015Famba", "Per\u015Fembe", "Cuma", "Cumartesi"]
};
var dayPeriodValues11 = {
  narrow: {
    am: "\xF6\xF6",
    pm: "\xF6s",
    midnight: "gy",
    noon: "\xF6",
    morning: "sa",
    afternoon: "\xF6s",
    evening: "ak",
    night: "ge"
  },
  abbreviated: {
    am: "\xD6\xD6",
    pm: "\xD6S",
    midnight: "gece yar\u0131s\u0131",
    noon: "\xF6\u011Fle",
    morning: "sabah",
    afternoon: "\xF6\u011Fleden sonra",
    evening: "ak\u015Fam",
    night: "gece"
  },
  wide: {
    am: "\xD6.\xD6.",
    pm: "\xD6.S.",
    midnight: "gece yar\u0131s\u0131",
    noon: "\xF6\u011Fle",
    morning: "sabah",
    afternoon: "\xF6\u011Fleden sonra",
    evening: "ak\u015Fam",
    night: "gece"
  }
};
var formattingDayPeriodValues9 = {
  narrow: {
    am: "\xF6\xF6",
    pm: "\xF6s",
    midnight: "gy",
    noon: "\xF6",
    morning: "sa",
    afternoon: "\xF6s",
    evening: "ak",
    night: "ge"
  },
  abbreviated: {
    am: "\xD6\xD6",
    pm: "\xD6S",
    midnight: "gece yar\u0131s\u0131",
    noon: "\xF6\u011Flen",
    morning: "sabahleyin",
    afternoon: "\xF6\u011Fleden sonra",
    evening: "ak\u015Famleyin",
    night: "geceleyin"
  },
  wide: {
    am: "\xF6.\xF6.",
    pm: "\xF6.s.",
    midnight: "gece yar\u0131s\u0131",
    noon: "\xF6\u011Flen",
    morning: "sabahleyin",
    afternoon: "\xF6\u011Fleden sonra",
    evening: "ak\u015Famleyin",
    night: "geceleyin"
  }
};
var ordinalNumber11 = (dirtyNumber, _options) => {
  const number = Number(dirtyNumber);
  return number + ".";
};
var localize11 = {
  ordinalNumber: ordinalNumber11,
  era: buildLocalizeFn({
    values: eraValues11,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues11,
    defaultWidth: "wide",
    argumentCallback: (quarter) => Number(quarter) - 1
  }),
  month: buildLocalizeFn({
    values: monthValues11,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues11,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues11,
    defaultWidth: "wide",
    formattingValues: formattingDayPeriodValues9,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/tr/_lib/match.js
var matchOrdinalNumberPattern11 = /^(\d+)(\.)?/i;
var parseOrdinalNumberPattern11 = /\d+/i;
var matchEraPatterns11 = {
  narrow: /^(mö|ms)/i,
  abbreviated: /^(mö|ms)/i,
  wide: /^(milattan önce|milattan sonra)/i
};
var parseEraPatterns11 = {
  any: [/(^mö|^milattan önce)/i, /(^ms|^milattan sonra)/i]
};
var matchQuarterPatterns11 = {
  narrow: /^[1234]/i,
  abbreviated: /^[1234]ç/i,
  wide: /^((i|İ)lk|(i|İ)kinci|üçüncü|son) çeyrek/i
};
var parseQuarterPatterns11 = {
  any: [/1/i, /2/i, /3/i, /4/i],
  abbreviated: [/1ç/i, /2ç/i, /3ç/i, /4ç/i],
  wide: [/^(i|İ)lk çeyrek/i, /(i|İ)kinci çeyrek/i, /üçüncü çeyrek/i, /son çeyrek/i]
};
var matchMonthPatterns11 = {
  narrow: /^[oşmnhtaek]/i,
  abbreviated: /^(oca|şub|mar|nis|may|haz|tem|ağu|eyl|eki|kas|ara)/i,
  wide: /^(ocak|şubat|mart|nisan|mayıs|haziran|temmuz|ağustos|eylül|ekim|kasım|aralık)/i
};
var parseMonthPatterns11 = {
  narrow: [/^o/i, /^ş/i, /^m/i, /^n/i, /^m/i, /^h/i, /^t/i, /^a/i, /^e/i, /^e/i, /^k/i, /^a/i],
  any: [/^o/i, /^ş/i, /^mar/i, /^n/i, /^may/i, /^h/i, /^t/i, /^ağ/i, /^ey/i, /^ek/i, /^k/i, /^ar/i]
};
var matchDayPatterns11 = {
  narrow: /^[psçc]/i,
  short: /^(pz|pt|sa|ça|pe|cu|ct)/i,
  abbreviated: /^(paz|pzt|sal|çar|per|cum|cts)/i,
  wide: /^(pazar(?!tesi)|pazartesi|salı|çarşamba|perşembe|cuma(?!rtesi)|cumartesi)/i
};
var parseDayPatterns11 = {
  narrow: [/^p/i, /^p/i, /^s/i, /^ç/i, /^p/i, /^c/i, /^c/i],
  any: [/^pz/i, /^pt/i, /^sa/i, /^ça/i, /^pe/i, /^cu/i, /^ct/i],
  wide: [/^pazar(?!tesi)/i, /^pazartesi/i, /^salı/i, /^çarşamba/i, /^perşembe/i, /^cuma(?!rtesi)/i, /^cumartesi/i]
};
var matchDayPeriodPatterns11 = {
  narrow: /^(öö|ös|gy|ö|sa|ös|ak|ge)/i,
  any: /^(ö\.?\s?[ös]\.?|öğleden sonra|gece yarısı|öğle|(sabah|öğ|akşam|gece)(leyin))/i
};
var parseDayPeriodPatterns11 = {
  any: {
    am: /^ö\.?ö\.?/i,
    pm: /^ö\.?s\.?/i,
    midnight: /^(gy|gece yarısı)/i,
    noon: /^öğ/i,
    morning: /^sa/i,
    afternoon: /^öğleden sonra/i,
    evening: /^ak/i,
    night: /^ge/i
  }
};
var match11 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern11,
    parsePattern: parseOrdinalNumberPattern11,
    valueCallback: function(value) {
      return parseInt(value, 10);
    }
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns11,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns11,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns11,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns11,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns11,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns11,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns11,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns11,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns11,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns11,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/tr.js
var tr = {
  code: "tr",
  formatDistance: formatDistance11,
  formatLong: formatLong11,
  formatRelative: formatRelative11,
  localize: localize11,
  match: match11,
  options: {
    weekStartsOn: 1,
    firstWeekContainsDate: 1
  }
};

// node_modules/date-fns/locale/uk/_lib/formatDistance.js
function declension2(scheme, count) {
  if (scheme.one !== void 0 && count === 1) {
    return scheme.one;
  }
  const rem10 = count % 10;
  const rem100 = count % 100;
  if (rem10 === 1 && rem100 !== 11) {
    return scheme.singularNominative.replace("{{count}}", String(count));
  } else if (rem10 >= 2 && rem10 <= 4 && (rem100 < 10 || rem100 > 20)) {
    return scheme.singularGenitive.replace("{{count}}", String(count));
  } else {
    return scheme.pluralGenitive.replace("{{count}}", String(count));
  }
}
function buildLocalizeTokenFn2(scheme) {
  return (count, options) => {
    if (options && options.addSuffix) {
      if (options.comparison && options.comparison > 0) {
        if (scheme.future) {
          return declension2(scheme.future, count);
        } else {
          return "\u0437\u0430 " + declension2(scheme.regular, count);
        }
      } else {
        if (scheme.past) {
          return declension2(scheme.past, count);
        } else {
          return declension2(scheme.regular, count) + " \u0442\u043E\u043C\u0443";
        }
      }
    } else {
      return declension2(scheme.regular, count);
    }
  };
}
var halfAtMinute = (_, options) => {
  if (options && options.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return "\u0437\u0430 \u043F\u0456\u0432\u0445\u0432\u0438\u043B\u0438\u043D\u0438";
    } else {
      return "\u043F\u0456\u0432\u0445\u0432\u0438\u043B\u0438\u043D\u0438 \u0442\u043E\u043C\u0443";
    }
  }
  return "\u043F\u0456\u0432\u0445\u0432\u0438\u043B\u0438\u043D\u0438";
};
var formatDistanceLocale12 = {
  lessThanXSeconds: buildLocalizeTokenFn2({
    regular: {
      one: "\u043C\u0435\u043D\u0448\u0435 \u0441\u0435\u043A\u0443\u043D\u0434\u0438",
      singularNominative: "\u043C\u0435\u043D\u0448\u0435 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0438",
      singularGenitive: "\u043C\u0435\u043D\u0448\u0435 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434",
      pluralGenitive: "\u043C\u0435\u043D\u0448\u0435 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434"
    },
    future: {
      one: "\u043C\u0435\u043D\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 \u0441\u0435\u043A\u0443\u043D\u0434\u0443",
      singularNominative: "\u043C\u0435\u043D\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0443",
      singularGenitive: "\u043C\u0435\u043D\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0438",
      pluralGenitive: "\u043C\u0435\u043D\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434"
    }
  }),
  xSeconds: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0430",
      singularGenitive: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0438",
      pluralGenitive: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434"
    },
    past: {
      singularNominative: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0443 \u0442\u043E\u043C\u0443",
      singularGenitive: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0438 \u0442\u043E\u043C\u0443",
      pluralGenitive: "{{count}} \u0441\u0435\u043A\u0443\u043D\u0434 \u0442\u043E\u043C\u0443"
    },
    future: {
      singularNominative: "\u0437\u0430 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0443",
      singularGenitive: "\u0437\u0430 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434\u0438",
      pluralGenitive: "\u0437\u0430 {{count}} \u0441\u0435\u043A\u0443\u043D\u0434"
    }
  }),
  halfAMinute: halfAtMinute,
  lessThanXMinutes: buildLocalizeTokenFn2({
    regular: {
      one: "\u043C\u0435\u043D\u0448\u0435 \u0445\u0432\u0438\u043B\u0438\u043D\u0438",
      singularNominative: "\u043C\u0435\u043D\u0448\u0435 {{count}} \u0445\u0432\u0438\u043B\u0438\u043D\u0438",
      singularGenitive: "\u043C\u0435\u043D\u0448\u0435 {{count}} \u0445\u0432\u0438\u043B\u0438\u043D",
      pluralGenitive: "\u043C\u0435\u043D\u0448\u0435 {{count}} \u0445\u0432\u0438\u043B\u0438\u043D"
    },
    future: {
      one: "\u043C\u0435\u043D\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 \u0445\u0432\u0438\u043B\u0438\u043D\u0443",
      singularNominative: "\u043C\u0435\u043D\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 {{count}} \u0445\u0432\u0438\u043B\u0438\u043D\u0443",
      singularGenitive: "\u043C\u0435\u043D\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 {{count}} \u0445\u0432\u0438\u043B\u0438\u043D\u0438",
      pluralGenitive: "\u043C\u0435\u043D\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 {{count}} \u0445\u0432\u0438\u043B\u0438\u043D"
    }
  }),
  xMinutes: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "{{count}} \u0445\u0432\u0438\u043B\u0438\u043D\u0430",
      singularGenitive: "{{count}} \u0445\u0432\u0438\u043B\u0438\u043D\u0438",
      pluralGenitive: "{{count}} \u0445\u0432\u0438\u043B\u0438\u043D"
    },
    past: {
      singularNominative: "{{count}} \u0445\u0432\u0438\u043B\u0438\u043D\u0443 \u0442\u043E\u043C\u0443",
      singularGenitive: "{{count}} \u0445\u0432\u0438\u043B\u0438\u043D\u0438 \u0442\u043E\u043C\u0443",
      pluralGenitive: "{{count}} \u0445\u0432\u0438\u043B\u0438\u043D \u0442\u043E\u043C\u0443"
    },
    future: {
      singularNominative: "\u0437\u0430 {{count}} \u0445\u0432\u0438\u043B\u0438\u043D\u0443",
      singularGenitive: "\u0437\u0430 {{count}} \u0445\u0432\u0438\u043B\u0438\u043D\u0438",
      pluralGenitive: "\u0437\u0430 {{count}} \u0445\u0432\u0438\u043B\u0438\u043D"
    }
  }),
  aboutXHours: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u0433\u043E\u0434\u0438\u043D\u0438",
      singularGenitive: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u0433\u043E\u0434\u0438\u043D",
      pluralGenitive: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u0433\u043E\u0434\u0438\u043D"
    },
    future: {
      singularNominative: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u0433\u043E\u0434\u0438\u043D\u0443",
      singularGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u0433\u043E\u0434\u0438\u043D\u0438",
      pluralGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u0433\u043E\u0434\u0438\u043D"
    }
  }),
  xHours: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "{{count}} \u0433\u043E\u0434\u0438\u043D\u0443",
      singularGenitive: "{{count}} \u0433\u043E\u0434\u0438\u043D\u0438",
      pluralGenitive: "{{count}} \u0433\u043E\u0434\u0438\u043D"
    }
  }),
  xDays: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "{{count}} \u0434\u0435\u043D\u044C",
      singularGenitive: "{{count}} \u0434\u043Di",
      pluralGenitive: "{{count}} \u0434\u043D\u0456\u0432"
    }
  }),
  aboutXWeeks: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u0442\u0438\u0436\u043D\u044F",
      singularGenitive: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u0442\u0438\u0436\u043D\u0456\u0432",
      pluralGenitive: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u0442\u0438\u0436\u043D\u0456\u0432"
    },
    future: {
      singularNominative: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u0442\u0438\u0436\u0434\u0435\u043D\u044C",
      singularGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u0442\u0438\u0436\u043D\u0456",
      pluralGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u0442\u0438\u0436\u043D\u0456\u0432"
    }
  }),
  xWeeks: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "{{count}} \u0442\u0438\u0436\u0434\u0435\u043D\u044C",
      singularGenitive: "{{count}} \u0442\u0438\u0436\u043D\u0456",
      pluralGenitive: "{{count}} \u0442\u0438\u0436\u043D\u0456\u0432"
    }
  }),
  aboutXMonths: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u043C\u0456\u0441\u044F\u0446\u044F",
      singularGenitive: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u043C\u0456\u0441\u044F\u0446\u0456\u0432",
      pluralGenitive: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u043C\u0456\u0441\u044F\u0446\u0456\u0432"
    },
    future: {
      singularNominative: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u043C\u0456\u0441\u044F\u0446\u044C",
      singularGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u043C\u0456\u0441\u044F\u0446\u0456",
      pluralGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u043C\u0456\u0441\u044F\u0446\u0456\u0432"
    }
  }),
  xMonths: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "{{count}} \u043C\u0456\u0441\u044F\u0446\u044C",
      singularGenitive: "{{count}} \u043C\u0456\u0441\u044F\u0446\u0456",
      pluralGenitive: "{{count}} \u043C\u0456\u0441\u044F\u0446\u0456\u0432"
    }
  }),
  aboutXYears: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u0440\u043E\u043A\u0443",
      singularGenitive: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u0440\u043E\u043A\u0456\u0432",
      pluralGenitive: "\u0431\u043B\u0438\u0437\u044C\u043A\u043E {{count}} \u0440\u043E\u043A\u0456\u0432"
    },
    future: {
      singularNominative: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u0440\u0456\u043A",
      singularGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u0440\u043E\u043A\u0438",
      pluralGenitive: "\u043F\u0440\u0438\u0431\u043B\u0438\u0437\u043D\u043E \u0437\u0430 {{count}} \u0440\u043E\u043A\u0456\u0432"
    }
  }),
  xYears: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "{{count}} \u0440\u0456\u043A",
      singularGenitive: "{{count}} \u0440\u043E\u043A\u0438",
      pluralGenitive: "{{count}} \u0440\u043E\u043A\u0456\u0432"
    }
  }),
  overXYears: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "\u0431\u0456\u043B\u044C\u0448\u0435 {{count}} \u0440\u043E\u043A\u0443",
      singularGenitive: "\u0431\u0456\u043B\u044C\u0448\u0435 {{count}} \u0440\u043E\u043A\u0456\u0432",
      pluralGenitive: "\u0431\u0456\u043B\u044C\u0448\u0435 {{count}} \u0440\u043E\u043A\u0456\u0432"
    },
    future: {
      singularNominative: "\u0431\u0456\u043B\u044C\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 {{count}} \u0440\u0456\u043A",
      singularGenitive: "\u0431\u0456\u043B\u044C\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 {{count}} \u0440\u043E\u043A\u0438",
      pluralGenitive: "\u0431\u0456\u043B\u044C\u0448\u0435, \u043D\u0456\u0436 \u0437\u0430 {{count}} \u0440\u043E\u043A\u0456\u0432"
    }
  }),
  almostXYears: buildLocalizeTokenFn2({
    regular: {
      singularNominative: "\u043C\u0430\u0439\u0436\u0435 {{count}} \u0440\u0456\u043A",
      singularGenitive: "\u043C\u0430\u0439\u0436\u0435 {{count}} \u0440\u043E\u043A\u0438",
      pluralGenitive: "\u043C\u0430\u0439\u0436\u0435 {{count}} \u0440\u043E\u043A\u0456\u0432"
    },
    future: {
      singularNominative: "\u043C\u0430\u0439\u0436\u0435 \u0437\u0430 {{count}} \u0440\u0456\u043A",
      singularGenitive: "\u043C\u0430\u0439\u0436\u0435 \u0437\u0430 {{count}} \u0440\u043E\u043A\u0438",
      pluralGenitive: "\u043C\u0430\u0439\u0436\u0435 \u0437\u0430 {{count}} \u0440\u043E\u043A\u0456\u0432"
    }
  })
};
var formatDistance12 = (token, count, options) => {
  options = options || {};
  return formatDistanceLocale12[token](count, options);
};

// node_modules/date-fns/locale/uk/_lib/formatLong.js
var dateFormats12 = {
  full: "EEEE, do MMMM y '\u0440.'",
  long: "do MMMM y '\u0440.'",
  medium: "d MMM y '\u0440.'",
  short: "dd.MM.y"
};
var timeFormats12 = {
  full: "H:mm:ss zzzz",
  long: "H:mm:ss z",
  medium: "H:mm:ss",
  short: "H:mm"
};
var dateTimeFormats12 = {
  full: "{{date}} '\u043E' {{time}}",
  long: "{{date}} '\u043E' {{time}}",
  medium: "{{date}}, {{time}}",
  short: "{{date}}, {{time}}"
};
var formatLong12 = {
  date: buildFormatLongFn({
    formats: dateFormats12,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats12,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats12,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/uk/_lib/formatRelative.js
var accusativeWeekdays2 = ["\u043D\u0435\u0434\u0456\u043B\u044E", "\u043F\u043E\u043D\u0435\u0434\u0456\u043B\u043E\u043A", "\u0432\u0456\u0432\u0442\u043E\u0440\u043E\u043A", "\u0441\u0435\u0440\u0435\u0434\u0443", "\u0447\u0435\u0442\u0432\u0435\u0440", "\u043F\u2019\u044F\u0442\u043D\u0438\u0446\u044E", "\u0441\u0443\u0431\u043E\u0442\u0443"];
function lastWeek2(day) {
  const weekday = accusativeWeekdays2[day];
  switch (day) {
    case 0:
    case 3:
    case 5:
    case 6:
      return "'\u0443 \u043C\u0438\u043D\u0443\u043B\u0443 " + weekday + " \u043E' p";
    case 1:
    case 2:
    case 4:
      return "'\u0443 \u043C\u0438\u043D\u0443\u043B\u0438\u0439 " + weekday + " \u043E' p";
  }
}
function thisWeek2(day) {
  const weekday = accusativeWeekdays2[day];
  return "'\u0443 " + weekday + " \u043E' p";
}
function nextWeek2(day) {
  const weekday = accusativeWeekdays2[day];
  switch (day) {
    case 0:
    case 3:
    case 5:
    case 6:
      return "'\u0443 \u043D\u0430\u0441\u0442\u0443\u043F\u043D\u0443 " + weekday + " \u043E' p";
    case 1:
    case 2:
    case 4:
      return "'\u0443 \u043D\u0430\u0441\u0442\u0443\u043F\u043D\u0438\u0439 " + weekday + " \u043E' p";
  }
}
var lastWeekFormat = (dirtyDate, baseDate, options) => {
  const date = toDate(dirtyDate);
  const day = date.getDay();
  if (isSameWeek(date, baseDate, options)) {
    return thisWeek2(day);
  } else {
    return lastWeek2(day);
  }
};
var nextWeekFormat = (dirtyDate, baseDate, options) => {
  const date = toDate(dirtyDate);
  const day = date.getDay();
  if (isSameWeek(date, baseDate, options)) {
    return thisWeek2(day);
  } else {
    return nextWeek2(day);
  }
};
var formatRelativeLocale12 = {
  lastWeek: lastWeekFormat,
  yesterday: "'\u0432\u0447\u043E\u0440\u0430 \u043E' p",
  today: "'\u0441\u044C\u043E\u0433\u043E\u0434\u043D\u0456 \u043E' p",
  tomorrow: "'\u0437\u0430\u0432\u0442\u0440\u0430 \u043E' p",
  nextWeek: nextWeekFormat,
  other: "P"
};
var formatRelative12 = (token, date, baseDate, options) => {
  const format = formatRelativeLocale12[token];
  if (typeof format === "function") {
    return format(date, baseDate, options);
  }
  return format;
};

// node_modules/date-fns/locale/uk/_lib/localize.js
var eraValues12 = {
  narrow: ["\u0434\u043E \u043D.\u0435.", "\u043D.\u0435."],
  abbreviated: ["\u0434\u043E \u043D. \u0435.", "\u043D. \u0435."],
  wide: ["\u0434\u043E \u043D\u0430\u0448\u043E\u0457 \u0435\u0440\u0438", "\u043D\u0430\u0448\u043E\u0457 \u0435\u0440\u0438"]
};
var quarterValues12 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["1-\u0439 \u043A\u0432.", "2-\u0439 \u043A\u0432.", "3-\u0439 \u043A\u0432.", "4-\u0439 \u043A\u0432."],
  wide: ["1-\u0439 \u043A\u0432\u0430\u0440\u0442\u0430\u043B", "2-\u0439 \u043A\u0432\u0430\u0440\u0442\u0430\u043B", "3-\u0439 \u043A\u0432\u0430\u0440\u0442\u0430\u043B", "4-\u0439 \u043A\u0432\u0430\u0440\u0442\u0430\u043B"]
};
var monthValues12 = {
  // ДСТУ 3582:2013
  narrow: ["\u0421", "\u041B", "\u0411", "\u041A", "\u0422", "\u0427", "\u041B", "\u0421", "\u0412", "\u0416", "\u041B", "\u0413"],
  abbreviated: ["\u0441\u0456\u0447.", "\u043B\u044E\u0442.", "\u0431\u0435\u0440\u0435\u0437.", "\u043A\u0432\u0456\u0442.", "\u0442\u0440\u0430\u0432.", "\u0447\u0435\u0440\u0432.", "\u043B\u0438\u043F.", "\u0441\u0435\u0440\u043F.", "\u0432\u0435\u0440\u0435\u0441.", "\u0436\u043E\u0432\u0442.", "\u043B\u0438\u0441\u0442\u043E\u043F.", "\u0433\u0440\u0443\u0434."],
  wide: ["\u0441\u0456\u0447\u0435\u043D\u044C", "\u043B\u044E\u0442\u0438\u0439", "\u0431\u0435\u0440\u0435\u0437\u0435\u043D\u044C", "\u043A\u0432\u0456\u0442\u0435\u043D\u044C", "\u0442\u0440\u0430\u0432\u0435\u043D\u044C", "\u0447\u0435\u0440\u0432\u0435\u043D\u044C", "\u043B\u0438\u043F\u0435\u043D\u044C", "\u0441\u0435\u0440\u043F\u0435\u043D\u044C", "\u0432\u0435\u0440\u0435\u0441\u0435\u043D\u044C", "\u0436\u043E\u0432\u0442\u0435\u043D\u044C", "\u043B\u0438\u0441\u0442\u043E\u043F\u0430\u0434", "\u0433\u0440\u0443\u0434\u0435\u043D\u044C"]
};
var formattingMonthValues3 = {
  narrow: ["\u0421", "\u041B", "\u0411", "\u041A", "\u0422", "\u0427", "\u041B", "\u0421", "\u0412", "\u0416", "\u041B", "\u0413"],
  abbreviated: ["\u0441\u0456\u0447.", "\u043B\u044E\u0442.", "\u0431\u0435\u0440\u0435\u0437.", "\u043A\u0432\u0456\u0442.", "\u0442\u0440\u0430\u0432.", "\u0447\u0435\u0440\u0432.", "\u043B\u0438\u043F.", "\u0441\u0435\u0440\u043F.", "\u0432\u0435\u0440\u0435\u0441.", "\u0436\u043E\u0432\u0442.", "\u043B\u0438\u0441\u0442\u043E\u043F.", "\u0433\u0440\u0443\u0434."],
  wide: ["\u0441\u0456\u0447\u043D\u044F", "\u043B\u044E\u0442\u043E\u0433\u043E", "\u0431\u0435\u0440\u0435\u0437\u043D\u044F", "\u043A\u0432\u0456\u0442\u043D\u044F", "\u0442\u0440\u0430\u0432\u043D\u044F", "\u0447\u0435\u0440\u0432\u043D\u044F", "\u043B\u0438\u043F\u043D\u044F", "\u0441\u0435\u0440\u043F\u043D\u044F", "\u0432\u0435\u0440\u0435\u0441\u043D\u044F", "\u0436\u043E\u0432\u0442\u043D\u044F", "\u043B\u0438\u0441\u0442\u043E\u043F\u0430\u0434\u0430", "\u0433\u0440\u0443\u0434\u043D\u044F"]
};
var dayValues12 = {
  narrow: ["\u041D", "\u041F", "\u0412", "\u0421", "\u0427", "\u041F", "\u0421"],
  short: ["\u043D\u0434", "\u043F\u043D", "\u0432\u0442", "\u0441\u0440", "\u0447\u0442", "\u043F\u0442", "\u0441\u0431"],
  abbreviated: ["\u043D\u0435\u0434", "\u043F\u043E\u043D", "\u0432\u0456\u0432", "\u0441\u0435\u0440", "\u0447\u0442\u0432", "\u043F\u0442\u043D", "\u0441\u0443\u0431"],
  wide: ["\u043D\u0435\u0434\u0456\u043B\u044F", "\u043F\u043E\u043D\u0435\u0434\u0456\u043B\u043E\u043A", "\u0432\u0456\u0432\u0442\u043E\u0440\u043E\u043A", "\u0441\u0435\u0440\u0435\u0434\u0430", "\u0447\u0435\u0442\u0432\u0435\u0440", "\u043F\u2019\u044F\u0442\u043D\u0438\u0446\u044F", "\u0441\u0443\u0431\u043E\u0442\u0430"]
};
var dayPeriodValues12 = {
  narrow: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u0456\u0432\u043D.",
    noon: "\u043F\u043E\u043B.",
    morning: "\u0440\u0430\u043D\u043E\u043A",
    afternoon: "\u0434\u0435\u043D\u044C",
    evening: "\u0432\u0435\u0447.",
    night: "\u043D\u0456\u0447"
  },
  abbreviated: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u0456\u0432\u043D.",
    noon: "\u043F\u043E\u043B.",
    morning: "\u0440\u0430\u043D\u043E\u043A",
    afternoon: "\u0434\u0435\u043D\u044C",
    evening: "\u0432\u0435\u0447.",
    night: "\u043D\u0456\u0447"
  },
  wide: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u0456\u0432\u043D\u0456\u0447",
    noon: "\u043F\u043E\u043B\u0443\u0434\u0435\u043D\u044C",
    morning: "\u0440\u0430\u043D\u043E\u043A",
    afternoon: "\u0434\u0435\u043D\u044C",
    evening: "\u0432\u0435\u0447\u0456\u0440",
    night: "\u043D\u0456\u0447"
  }
};
var formattingDayPeriodValues10 = {
  narrow: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u0456\u0432\u043D.",
    noon: "\u043F\u043E\u043B.",
    morning: "\u0440\u0430\u043D\u043A\u0443",
    afternoon: "\u0434\u043D\u044F",
    evening: "\u0432\u0435\u0447.",
    night: "\u043D\u043E\u0447\u0456"
  },
  abbreviated: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u0456\u0432\u043D.",
    noon: "\u043F\u043E\u043B.",
    morning: "\u0440\u0430\u043D\u043A\u0443",
    afternoon: "\u0434\u043D\u044F",
    evening: "\u0432\u0435\u0447.",
    night: "\u043D\u043E\u0447\u0456"
  },
  wide: {
    am: "\u0414\u041F",
    pm: "\u041F\u041F",
    midnight: "\u043F\u0456\u0432\u043D\u0456\u0447",
    noon: "\u043F\u043E\u043B\u0443\u0434\u0435\u043D\u044C",
    morning: "\u0440\u0430\u043D\u043A\u0443",
    afternoon: "\u0434\u043D\u044F",
    evening: "\u0432\u0435\u0447.",
    night: "\u043D\u043E\u0447\u0456"
  }
};
var ordinalNumber12 = (dirtyNumber, options) => {
  const unit = String(options?.unit);
  const number = Number(dirtyNumber);
  let suffix;
  if (unit === "date") {
    if (number === 3 || number === 23) {
      suffix = "-\u0454";
    } else {
      suffix = "-\u0435";
    }
  } else if (unit === "minute" || unit === "second" || unit === "hour") {
    suffix = "-\u0430";
  } else {
    suffix = "-\u0439";
  }
  return number + suffix;
};
var localize12 = {
  ordinalNumber: ordinalNumber12,
  era: buildLocalizeFn({
    values: eraValues12,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues12,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues12,
    defaultWidth: "wide",
    formattingValues: formattingMonthValues3,
    defaultFormattingWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues12,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues12,
    defaultWidth: "any",
    formattingValues: formattingDayPeriodValues10,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/uk/_lib/match.js
var matchOrdinalNumberPattern12 = /^(\d+)(-?(е|й|є|а|я))?/i;
var parseOrdinalNumberPattern12 = /\d+/i;
var matchEraPatterns12 = {
  narrow: /^((до )?н\.?\s?е\.?)/i,
  abbreviated: /^((до )?н\.?\s?е\.?)/i,
  wide: /^(до нашої ери|нашої ери|наша ера)/i
};
var parseEraPatterns12 = {
  any: [/^д/i, /^н/i]
};
var matchQuarterPatterns12 = {
  narrow: /^[1234]/i,
  abbreviated: /^[1234](-?[иі]?й?)? кв.?/i,
  wide: /^[1234](-?[иі]?й?)? квартал/i
};
var parseQuarterPatterns12 = {
  any: [/1/i, /2/i, /3/i, /4/i]
};
var matchMonthPatterns12 = {
  narrow: /^[слбктчвжг]/i,
  abbreviated: /^(січ|лют|бер(ез)?|квіт|трав|черв|лип|серп|вер(ес)?|жовт|лис(топ)?|груд)\.?/i,
  wide: /^(січень|січня|лютий|лютого|березень|березня|квітень|квітня|травень|травня|червня|червень|липень|липня|серпень|серпня|вересень|вересня|жовтень|жовтня|листопад[а]?|грудень|грудня)/i
};
var parseMonthPatterns12 = {
  narrow: [/^с/i, /^л/i, /^б/i, /^к/i, /^т/i, /^ч/i, /^л/i, /^с/i, /^в/i, /^ж/i, /^л/i, /^г/i],
  any: [/^сі/i, /^лю/i, /^б/i, /^к/i, /^т/i, /^ч/i, /^лип/i, /^се/i, /^в/i, /^ж/i, /^лис/i, /^г/i]
};
var matchDayPatterns12 = {
  narrow: /^[нпвсч]/i,
  short: /^(нд|пн|вт|ср|чт|пт|сб)\.?/i,
  abbreviated: /^(нед|пон|вів|сер|че?тв|птн?|суб)\.?/i,
  wide: /^(неділ[яі]|понеділ[ок][ка]|вівтор[ок][ка]|серед[аи]|четвер(га)?|п\W*?ятниц[яі]|субот[аи])/i
};
var parseDayPatterns12 = {
  narrow: [/^н/i, /^п/i, /^в/i, /^с/i, /^ч/i, /^п/i, /^с/i],
  any: [/^н/i, /^п[он]/i, /^в/i, /^с[ер]/i, /^ч/i, /^п\W*?[ят]/i, /^с[уб]/i]
};
var matchDayPeriodPatterns12 = {
  narrow: /^([дп]п|півн\.?|пол\.?|ранок|ранку|день|дня|веч\.?|ніч|ночі)/i,
  abbreviated: /^([дп]п|півн\.?|пол\.?|ранок|ранку|день|дня|веч\.?|ніч|ночі)/i,
  wide: /^([дп]п|північ|полудень|ранок|ранку|день|дня|вечір|вечора|ніч|ночі)/i
};
var parseDayPeriodPatterns12 = {
  any: {
    am: /^дп/i,
    pm: /^пп/i,
    midnight: /^півн/i,
    noon: /^пол/i,
    morning: /^р/i,
    afternoon: /^д[ен]/i,
    evening: /^в/i,
    night: /^н/i
  }
};
var match12 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern12,
    parsePattern: parseOrdinalNumberPattern12,
    valueCallback: (value) => parseInt(value, 10)
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns12,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns12,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns12,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns12,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns12,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns12,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns12,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns12,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns12,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPeriodPatterns12,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/uk.js
var uk = {
  code: "uk",
  formatDistance: formatDistance12,
  formatLong: formatLong12,
  formatRelative: formatRelative12,
  localize: localize12,
  match: match12,
  options: {
    weekStartsOn: 1,
    firstWeekContainsDate: 1
  }
};

// node_modules/date-fns/locale/zh-CN/_lib/formatDistance.js
var formatDistanceLocale13 = {
  lessThanXSeconds: {
    one: "\u4E0D\u5230 1 \u79D2",
    other: "\u4E0D\u5230 {{count}} \u79D2"
  },
  xSeconds: {
    one: "1 \u79D2",
    other: "{{count}} \u79D2"
  },
  halfAMinute: "\u534A\u5206\u949F",
  lessThanXMinutes: {
    one: "\u4E0D\u5230 1 \u5206\u949F",
    other: "\u4E0D\u5230 {{count}} \u5206\u949F"
  },
  xMinutes: {
    one: "1 \u5206\u949F",
    other: "{{count}} \u5206\u949F"
  },
  xHours: {
    one: "1 \u5C0F\u65F6",
    other: "{{count}} \u5C0F\u65F6"
  },
  aboutXHours: {
    one: "\u5927\u7EA6 1 \u5C0F\u65F6",
    other: "\u5927\u7EA6 {{count}} \u5C0F\u65F6"
  },
  xDays: {
    one: "1 \u5929",
    other: "{{count}} \u5929"
  },
  aboutXWeeks: {
    one: "\u5927\u7EA6 1 \u4E2A\u661F\u671F",
    other: "\u5927\u7EA6 {{count}} \u4E2A\u661F\u671F"
  },
  xWeeks: {
    one: "1 \u4E2A\u661F\u671F",
    other: "{{count}} \u4E2A\u661F\u671F"
  },
  aboutXMonths: {
    one: "\u5927\u7EA6 1 \u4E2A\u6708",
    other: "\u5927\u7EA6 {{count}} \u4E2A\u6708"
  },
  xMonths: {
    one: "1 \u4E2A\u6708",
    other: "{{count}} \u4E2A\u6708"
  },
  aboutXYears: {
    one: "\u5927\u7EA6 1 \u5E74",
    other: "\u5927\u7EA6 {{count}} \u5E74"
  },
  xYears: {
    one: "1 \u5E74",
    other: "{{count}} \u5E74"
  },
  overXYears: {
    one: "\u8D85\u8FC7 1 \u5E74",
    other: "\u8D85\u8FC7 {{count}} \u5E74"
  },
  almostXYears: {
    one: "\u5C06\u8FD1 1 \u5E74",
    other: "\u5C06\u8FD1 {{count}} \u5E74"
  }
};
var formatDistance13 = (token, count, options) => {
  let result;
  const tokenValue = formatDistanceLocale13[token];
  if (typeof tokenValue === "string") {
    result = tokenValue;
  } else if (count === 1) {
    result = tokenValue.one;
  } else {
    result = tokenValue.other.replace("{{count}}", String(count));
  }
  if (options?.addSuffix) {
    if (options.comparison && options.comparison > 0) {
      return result + "\u5185";
    } else {
      return result + "\u524D";
    }
  }
  return result;
};

// node_modules/date-fns/locale/zh-CN/_lib/formatLong.js
var dateFormats13 = {
  full: "y'\u5E74'M'\u6708'd'\u65E5' EEEE",
  long: "y'\u5E74'M'\u6708'd'\u65E5'",
  medium: "yyyy-MM-dd",
  short: "yy-MM-dd"
};
var timeFormats13 = {
  full: "zzzz a h:mm:ss",
  long: "z a h:mm:ss",
  medium: "a h:mm:ss",
  short: "a h:mm"
};
var dateTimeFormats13 = {
  full: "{{date}} {{time}}",
  long: "{{date}} {{time}}",
  medium: "{{date}} {{time}}",
  short: "{{date}} {{time}}"
};
var formatLong13 = {
  date: buildFormatLongFn({
    formats: dateFormats13,
    defaultWidth: "full"
  }),
  time: buildFormatLongFn({
    formats: timeFormats13,
    defaultWidth: "full"
  }),
  dateTime: buildFormatLongFn({
    formats: dateTimeFormats13,
    defaultWidth: "full"
  })
};

// node_modules/date-fns/locale/zh-CN/_lib/formatRelative.js
function checkWeek(date, baseDate, options) {
  const baseFormat = "eeee p";
  if (isSameWeek(date, baseDate, options)) {
    return baseFormat;
  } else if (date.getTime() > baseDate.getTime()) {
    return "'\u4E0B\u4E2A'" + baseFormat;
  }
  return "'\u4E0A\u4E2A'" + baseFormat;
}
var formatRelativeLocale13 = {
  lastWeek: checkWeek,
  // days before yesterday, maybe in this week or last week
  yesterday: "'\u6628\u5929' p",
  today: "'\u4ECA\u5929' p",
  tomorrow: "'\u660E\u5929' p",
  nextWeek: checkWeek,
  // days after tomorrow, maybe in this week or next week
  other: "PP p"
};
var formatRelative13 = (token, date, baseDate, options) => {
  const format = formatRelativeLocale13[token];
  if (typeof format === "function") {
    return format(date, baseDate, options);
  }
  return format;
};

// node_modules/date-fns/locale/zh-CN/_lib/localize.js
var eraValues13 = {
  narrow: ["\u524D", "\u516C\u5143"],
  abbreviated: ["\u524D", "\u516C\u5143"],
  wide: ["\u516C\u5143\u524D", "\u516C\u5143"]
};
var quarterValues13 = {
  narrow: ["1", "2", "3", "4"],
  abbreviated: ["\u7B2C\u4E00\u5B63", "\u7B2C\u4E8C\u5B63", "\u7B2C\u4E09\u5B63", "\u7B2C\u56DB\u5B63"],
  wide: ["\u7B2C\u4E00\u5B63\u5EA6", "\u7B2C\u4E8C\u5B63\u5EA6", "\u7B2C\u4E09\u5B63\u5EA6", "\u7B2C\u56DB\u5B63\u5EA6"]
};
var monthValues13 = {
  narrow: ["\u4E00", "\u4E8C", "\u4E09", "\u56DB", "\u4E94", "\u516D", "\u4E03", "\u516B", "\u4E5D", "\u5341", "\u5341\u4E00", "\u5341\u4E8C"],
  abbreviated: ["1\u6708", "2\u6708", "3\u6708", "4\u6708", "5\u6708", "6\u6708", "7\u6708", "8\u6708", "9\u6708", "10\u6708", "11\u6708", "12\u6708"],
  wide: ["\u4E00\u6708", "\u4E8C\u6708", "\u4E09\u6708", "\u56DB\u6708", "\u4E94\u6708", "\u516D\u6708", "\u4E03\u6708", "\u516B\u6708", "\u4E5D\u6708", "\u5341\u6708", "\u5341\u4E00\u6708", "\u5341\u4E8C\u6708"]
};
var dayValues13 = {
  narrow: ["\u65E5", "\u4E00", "\u4E8C", "\u4E09", "\u56DB", "\u4E94", "\u516D"],
  short: ["\u65E5", "\u4E00", "\u4E8C", "\u4E09", "\u56DB", "\u4E94", "\u516D"],
  abbreviated: ["\u5468\u65E5", "\u5468\u4E00", "\u5468\u4E8C", "\u5468\u4E09", "\u5468\u56DB", "\u5468\u4E94", "\u5468\u516D"],
  wide: ["\u661F\u671F\u65E5", "\u661F\u671F\u4E00", "\u661F\u671F\u4E8C", "\u661F\u671F\u4E09", "\u661F\u671F\u56DB", "\u661F\u671F\u4E94", "\u661F\u671F\u516D"]
};
var dayPeriodValues13 = {
  narrow: {
    am: "\u4E0A",
    pm: "\u4E0B",
    midnight: "\u51CC\u6668",
    noon: "\u5348",
    morning: "\u65E9",
    afternoon: "\u4E0B\u5348",
    evening: "\u665A",
    night: "\u591C"
  },
  abbreviated: {
    am: "\u4E0A\u5348",
    pm: "\u4E0B\u5348",
    midnight: "\u51CC\u6668",
    noon: "\u4E2D\u5348",
    morning: "\u65E9\u6668",
    afternoon: "\u4E2D\u5348",
    evening: "\u665A\u4E0A",
    night: "\u591C\u95F4"
  },
  wide: {
    am: "\u4E0A\u5348",
    pm: "\u4E0B\u5348",
    midnight: "\u51CC\u6668",
    noon: "\u4E2D\u5348",
    morning: "\u65E9\u6668",
    afternoon: "\u4E2D\u5348",
    evening: "\u665A\u4E0A",
    night: "\u591C\u95F4"
  }
};
var formattingDayPeriodValues11 = {
  narrow: {
    am: "\u4E0A",
    pm: "\u4E0B",
    midnight: "\u51CC\u6668",
    noon: "\u5348",
    morning: "\u65E9",
    afternoon: "\u4E0B\u5348",
    evening: "\u665A",
    night: "\u591C"
  },
  abbreviated: {
    am: "\u4E0A\u5348",
    pm: "\u4E0B\u5348",
    midnight: "\u51CC\u6668",
    noon: "\u4E2D\u5348",
    morning: "\u65E9\u6668",
    afternoon: "\u4E2D\u5348",
    evening: "\u665A\u4E0A",
    night: "\u591C\u95F4"
  },
  wide: {
    am: "\u4E0A\u5348",
    pm: "\u4E0B\u5348",
    midnight: "\u51CC\u6668",
    noon: "\u4E2D\u5348",
    morning: "\u65E9\u6668",
    afternoon: "\u4E2D\u5348",
    evening: "\u665A\u4E0A",
    night: "\u591C\u95F4"
  }
};
var ordinalNumber13 = (dirtyNumber, options) => {
  const number = Number(dirtyNumber);
  switch (options?.unit) {
    case "date":
      return number.toString() + "\u65E5";
    case "hour":
      return number.toString() + "\u65F6";
    case "minute":
      return number.toString() + "\u5206";
    case "second":
      return number.toString() + "\u79D2";
    default:
      return "\u7B2C " + number.toString();
  }
};
var localize13 = {
  ordinalNumber: ordinalNumber13,
  era: buildLocalizeFn({
    values: eraValues13,
    defaultWidth: "wide"
  }),
  quarter: buildLocalizeFn({
    values: quarterValues13,
    defaultWidth: "wide",
    argumentCallback: (quarter) => quarter - 1
  }),
  month: buildLocalizeFn({
    values: monthValues13,
    defaultWidth: "wide"
  }),
  day: buildLocalizeFn({
    values: dayValues13,
    defaultWidth: "wide"
  }),
  dayPeriod: buildLocalizeFn({
    values: dayPeriodValues13,
    defaultWidth: "wide",
    formattingValues: formattingDayPeriodValues11,
    defaultFormattingWidth: "wide"
  })
};

// node_modules/date-fns/locale/zh-CN/_lib/match.js
var matchOrdinalNumberPattern13 = /^(第\s*)?\d+(日|时|分|秒)?/i;
var parseOrdinalNumberPattern13 = /\d+/i;
var matchEraPatterns13 = {
  narrow: /^(前)/i,
  abbreviated: /^(前)/i,
  wide: /^(公元前|公元)/i
};
var parseEraPatterns13 = {
  any: [/^(前)/i, /^(公元)/i]
};
var matchQuarterPatterns13 = {
  narrow: /^[1234]/i,
  abbreviated: /^第[一二三四]刻/i,
  wide: /^第[一二三四]刻钟/i
};
var parseQuarterPatterns13 = {
  any: [/(1|一)/i, /(2|二)/i, /(3|三)/i, /(4|四)/i]
};
var matchMonthPatterns13 = {
  narrow: /^(一|二|三|四|五|六|七|八|九|十[二一])/i,
  abbreviated: /^(一|二|三|四|五|六|七|八|九|十[二一]|\d|1[12])月/i,
  wide: /^(一|二|三|四|五|六|七|八|九|十[二一])月/i
};
var parseMonthPatterns13 = {
  narrow: [/^一/i, /^二/i, /^三/i, /^四/i, /^五/i, /^六/i, /^七/i, /^八/i, /^九/i, /^十(?!(一|二))/i, /^十一/i, /^十二/i],
  any: [/^一|1/i, /^二|2/i, /^三|3/i, /^四|4/i, /^五|5/i, /^六|6/i, /^七|7/i, /^八|8/i, /^九|9/i, /^十(?!(一|二))|10/i, /^十一|11/i, /^十二|12/i]
};
var matchDayPatterns13 = {
  narrow: /^[一二三四五六日]/i,
  short: /^[一二三四五六日]/i,
  abbreviated: /^周[一二三四五六日]/i,
  wide: /^星期[一二三四五六日]/i
};
var parseDayPatterns13 = {
  any: [/日/i, /一/i, /二/i, /三/i, /四/i, /五/i, /六/i]
};
var matchDayPeriodPatterns13 = {
  any: /^(上午?|下午?|午夜|[中正]午|早上?|下午|晚上?|凌晨|)/i
};
var parseDayPeriodPatterns13 = {
  any: {
    am: /^上午?/i,
    pm: /^下午?/i,
    midnight: /^午夜/i,
    noon: /^[中正]午/i,
    morning: /^早上/i,
    afternoon: /^下午/i,
    evening: /^晚上?/i,
    night: /^凌晨/i
  }
};
var match13 = {
  ordinalNumber: buildMatchPatternFn({
    matchPattern: matchOrdinalNumberPattern13,
    parsePattern: parseOrdinalNumberPattern13,
    valueCallback: (value) => parseInt(value, 10)
  }),
  era: buildMatchFn({
    matchPatterns: matchEraPatterns13,
    defaultMatchWidth: "wide",
    parsePatterns: parseEraPatterns13,
    defaultParseWidth: "any"
  }),
  quarter: buildMatchFn({
    matchPatterns: matchQuarterPatterns13,
    defaultMatchWidth: "wide",
    parsePatterns: parseQuarterPatterns13,
    defaultParseWidth: "any",
    valueCallback: (index) => index + 1
  }),
  month: buildMatchFn({
    matchPatterns: matchMonthPatterns13,
    defaultMatchWidth: "wide",
    parsePatterns: parseMonthPatterns13,
    defaultParseWidth: "any"
  }),
  day: buildMatchFn({
    matchPatterns: matchDayPatterns13,
    defaultMatchWidth: "wide",
    parsePatterns: parseDayPatterns13,
    defaultParseWidth: "any"
  }),
  dayPeriod: buildMatchFn({
    matchPatterns: matchDayPeriodPatterns13,
    defaultMatchWidth: "any",
    parsePatterns: parseDayPeriodPatterns13,
    defaultParseWidth: "any"
  })
};

// node_modules/date-fns/locale/zh-CN.js
var zhCN = {
  code: "zh-CN",
  formatDistance: formatDistance13,
  formatLong: formatLong13,
  formatRelative: formatRelative13,
  localize: localize13,
  match: match13,
  options: {
    weekStartsOn: 1,
    firstWeekContainsDate: 4
  }
};

// src/app/dates/dates.locales.ts
var locales = {
  ar,
  de,
  en: enUS,
  es,
  fr,
  hi,
  ja,
  nl,
  pt,
  ru,
  tr,
  uk,
  zh: zhCN
};
var resolveDateLocale = (locale) => locales[locale] ?? enUS;

// node_modules/date-fns/_lib/getTimezoneOffsetInMilliseconds.js
function getTimezoneOffsetInMilliseconds(date) {
  const _date = toDate(date);
  const utcDate = new Date(Date.UTC(_date.getFullYear(), _date.getMonth(), _date.getDate(), _date.getHours(), _date.getMinutes(), _date.getSeconds(), _date.getMilliseconds()));
  utcDate.setUTCFullYear(_date.getFullYear());
  return +date - +utcDate;
}

export {
  millisecondsInWeek,
  millisecondsInDay,
  minutesInMonth,
  minutesInDay,
  constructFrom,
  toDate,
  getDefaultOptions,
  startOfWeek,
  getTimezoneOffsetInMilliseconds,
  normalizeDates,
  enUS,
  resolveDateLocale
};
//# sourceMappingURL=chunk-3D6CEWET.js.map
