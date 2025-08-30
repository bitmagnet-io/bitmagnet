import {
  ChartComponent,
  createThemeColor,
  format
} from "./chunk-BMBEAU42.js";
import {
  ThemeInfoService
} from "./chunk-DCY4KWPQ.js";
import {
  resolveDateLocale
} from "./chunk-3D6CEWET.js";
import {
  ErrorsService
} from "./chunk-75G4HS47.js";
import {
  DocumentTitleComponent
} from "./chunk-OOXMQI6S.js";
import {
  BreakpointsService
} from "./chunk-NQ6E5D5R.js";
import {
  Apollo,
  AppModule,
  GraphQLModule,
  MatCard,
  MatCardContent,
  MatCardHeader,
  MatCardTitle,
  MatFormField,
  MatGridList,
  MatGridTile,
  MatIcon,
  MatIconButton,
  MatInput,
  MatOption,
  MatProgressBar,
  MatSelect,
  MatToolbar,
  MatTooltip,
  TorrentMetricsDocument,
  TranslocoDirective,
  TranslocoService
} from "./chunk-WWRDQTKJ.js";
import "./chunk-VSVMRYN2.js";
import "./chunk-Y2ZC5Z2X.js";
import {
  AsyncPipe,
  BehaviorSubject,
  EMPTY,
  __spreadProps,
  __spreadValues,
  catchError,
  debounceTime,
  inject,
  map,
  ɵsetClassDebugInfo,
  ɵɵStandaloneFeature,
  ɵɵadvance,
  ɵɵclassMap,
  ɵɵdefineComponent,
  ɵɵdefineInjectable,
  ɵɵelement,
  ɵɵelementContainerEnd,
  ɵɵelementContainerStart,
  ɵɵelementEnd,
  ɵɵelementStart,
  ɵɵgetCurrentView,
  ɵɵlistener,
  ɵɵnextContext,
  ɵɵpipe,
  ɵɵpipeBind1,
  ɵɵproperty,
  ɵɵpureFunction0,
  ɵɵpureFunction2,
  ɵɵrepeater,
  ɵɵrepeaterCreate,
  ɵɵrepeaterTrackByIdentity,
  ɵɵresetView,
  ɵɵrestoreView,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate,
  ɵɵtextInterpolate1
} from "./chunk-DMMUMX3A.js";

// src/app/dashboard/torrents/torrent-metrics.utils.ts
var createResult = (params, rawResult) => {
  const { bucketParams, earliestBucket } = createBucketParams(params, rawResult);
  const sources = Object.entries(rawResult.torrent.metrics.buckets.reduce((acc, next) => {
    if (next.source !== (params.source ?? next.source)) {
      return acc;
    }
    let bucket = normalizeBucket(next.bucket, bucketParams);
    if (earliestBucket && earliestBucket.index > bucket.index) {
      bucket = void 0;
    }
    if (!bucket) {
      return acc;
    }
    const currentEventBuckets = acc[next.source] ?? [];
    return __spreadProps(__spreadValues({}, acc), {
      [next.source]: {
        created: !next.updated ? __spreadProps(__spreadValues({}, currentEventBuckets.created), {
          [bucket.key]: {
            count: next.count + (currentEventBuckets.created?.[bucket.key]?.count ?? 0),
            startTime: bucket.start
          }
        }) : currentEventBuckets.created,
        updated: next.updated ? __spreadProps(__spreadValues({}, currentEventBuckets.updated), {
          [bucket.key]: {
            count: next.count + (currentEventBuckets.updated?.[bucket.key]?.count ?? 0),
            startTime: bucket.start
          }
        }) : currentEventBuckets.updated
      }
    });
  }, {})).map(([source, eventBuckets]) => {
    let events;
    if (Object.keys(eventBuckets).length) {
      const bucketDates = Array();
      const buckets = fromEntries(Array("created", "updated").flatMap((event) => {
        const entries = fromEntries(Object.entries(eventBuckets[event] ?? {}).filter(([, v]) => v?.count).sort(([a], [b]) => parseInt(a) < parseInt(b) ? 1 : -1));
        const keys = Object.keys(entries);
        if (!keys.length) {
          return [];
        }
        const earliestBucket2 = parseInt(keys[0]);
        const latestBucket = parseInt(keys[keys.length - 1]);
        bucketDates.push(earliestBucket2, latestBucket);
        return [
          [
            event,
            {
              earliestBucket: earliestBucket2,
              latestBucket,
              entries
            }
          ]
        ];
      }));
      bucketDates.sort();
      events = {
        bucketDuration: bucketParams.duration,
        earliestBucket: bucketDates[0],
        latestBucket: bucketDates[bucketDates.length - 1],
        eventBuckets: buckets
      };
    }
    return {
      source,
      events,
      isEmpty: !events?.eventBuckets
    };
  });
  let bucketSpan;
  const earliestFoundBucket = sources.flatMap((q) => q.events ? [q.events.earliestBucket] : []).sort()[0];
  const latestFoundBucket = sources.flatMap((q) => q.events ? [q.events.latestBucket] : []).sort().reverse()[0];
  if (earliestFoundBucket && latestFoundBucket) {
    bucketSpan = {
      earliestBucket: earliestFoundBucket,
      latestBucket: latestFoundBucket
    };
  }
  return {
    params: __spreadProps(__spreadValues({}, params), {
      buckets: bucketParams
    }),
    sourceSummaries: sources,
    bucketSpan,
    availableSources: rawResult.torrent.listSources.sources.map((s) => ({
      key: s.key,
      name: s.name
    }))
  };
};
var fromEntries = (entries) => Object.fromEntries(entries);
var createBucketParams = (params, rawResult) => {
  const duration = params.buckets.duration === "AUTO" ? "hour" : params.buckets.duration;
  let multiplier = params.buckets.multiplier === "AUTO" ? 1 : params.buckets.multiplier;
  const timeframe = params.buckets.timeframe;
  const now = /* @__PURE__ */ new Date();
  const nowBucket = normalizeBucket(now, { duration, multiplier });
  const startBucket = normalizeBucket(now.getTime() - 1e3 * timeframeLengths[timeframe], {
    duration,
    multiplier
  });
  const allBuckets = [
    startBucket,
    ...rawResult.torrent.metrics.buckets.flatMap((b) => [
      normalizeBucket(b.bucket, { duration, multiplier })
    ]),
    nowBucket
  ].filter((b) => b.index >= startBucket.index).sort((a, b) => a.index - b.index);
  const minBucket = allBuckets[0];
  const maxBucket = allBuckets[allBuckets.length - 1];
  if (params.buckets.multiplier === "AUTO") {
    const targetSpan = 20;
    const span = maxBucket.index - minBucket.index;
    multiplier = Math.min(60, Math.max(Math.floor(span / (targetSpan * 5)) * 5, 1));
  }
  return {
    bucketParams: {
      duration,
      multiplier,
      timeframe
    },
    earliestBucket: normalizeBucket(now.getTime() - 1e3 * timeframeLengths[timeframe], {
      duration,
      multiplier
    }),
    latestBucket: normalizeBucket(Math.max(now.getTime(), maxBucket.start.getTime()), { duration, multiplier })
  };
};
var normalizeBucket = (rawDate, params) => {
  const date = new Date(rawDate);
  const msMultiplier = 1e3 * durationSeconds[params.duration] * params.multiplier;
  const baseNumber = Math.floor(date.getTime() / msMultiplier);
  return {
    key: `${baseNumber}`,
    index: baseNumber,
    start: new Date(baseNumber * msMultiplier)
  };
};

// src/app/dashboard/torrents/torrent-metrics.constants.ts
var defaultBucketParams = {
  duration: "minute",
  multiplier: 1,
  timeframe: "hours_1"
};
var resolutionNames = ["day", "hour", "minute"];
var durationSeconds = {
  minute: 60,
  hour: 60 * 60,
  day: 60 * 60 * 24
};
var emptyParams = {
  buckets: defaultBucketParams,
  autoRefresh: "off"
};
var emptyRawResult = {
  torrent: {
    metrics: {
      buckets: []
    },
    listSources: {
      sources: [
        {
          key: "dht",
          name: "DHT"
        }
      ]
    }
  }
};
var eventNames = ["created", "updated"];
var timeframeNames = [
  "minutes_15",
  "minutes_30",
  "hours_1",
  "hours_6",
  "hours_12",
  "days_1",
  "weeks_1"
];
var timeframeLengths = {
  minutes_15: 60 * 15,
  minutes_30: 60 * 30,
  hours_1: 60 * 60,
  hours_6: 60 * 60 * 6,
  hours_12: 60 * 60 * 12,
  days_1: 60 * 60 * 24,
  weeks_1: 60 * 60 * 24 * 7
};
var autoRefreshIntervalNames = [
  "off",
  "seconds_10",
  "seconds_30",
  "minutes_1",
  "minutes_5"
];
var autoRefreshIntervals = {
  off: null,
  seconds_10: 10,
  seconds_30: 30,
  minutes_1: 60,
  minutes_5: 60 * 5
};
var emptyResult = createResult(emptyParams, emptyRawResult);

// src/app/dashboard/torrents/torrent-metrics.controller.ts
var TorrentMetricsController = class {
  constructor(apollo, initParams = emptyParams, errorsService) {
    this.apollo = apollo;
    this.errorsService = errorsService;
    this.rawResultSubject = new BehaviorSubject(emptyRawResult);
    this.resultSubject = new BehaviorSubject(emptyResult);
    this.result$ = this.resultSubject.asObservable();
    this.loadingSubject = new BehaviorSubject(false);
    this.paramsSubject = new BehaviorSubject(initParams);
    this.params$ = this.paramsSubject.asObservable();
    this.variablesSubject = new BehaviorSubject(createVariables(initParams));
    this.paramsSubject.pipe(debounceTime(50)).subscribe((params) => {
      const variables = this.variablesSubject.getValue();
      const nextVariables = createVariables(params);
      if (JSON.stringify(variables) !== JSON.stringify(nextVariables)) {
        this.variablesSubject.next(nextVariables);
      } else {
        this.resultSubject.next(createResult(params, this.rawResultSubject.getValue()));
      }
    });
    this.variablesSubject.pipe(debounceTime(50)).subscribe((variables) => this.request(variables));
    this.rawResultSubject.subscribe((rawResult) => {
      const params = this.paramsSubject.getValue();
      this.resultSubject.next(createResult(params, rawResult));
      this.setInterval(params.autoRefresh);
    });
  }
  setInterval(interval) {
    clearTimeout(this.refreshTimeout);
    const delay = autoRefreshIntervals[interval ?? this.params.autoRefresh];
    if (delay) {
      this.refreshTimeout = setTimeout(() => {
        this.refresh();
      }, delay * 1e3);
    }
  }
  get params() {
    return this.paramsSubject.getValue();
  }
  get bucketDuration() {
    const d = this.params.buckets.duration;
    if (d === "AUTO") {
      return "hour";
    }
    return d;
  }
  get bucketMultiplier() {
    return this.resultSubject.getValue().params.buckets.multiplier ?? this.params.buckets.multiplier;
  }
  get loading() {
    return this.loadingSubject.getValue();
  }
  setTimeframe(timeframe) {
    this.updateParams((p) => __spreadProps(__spreadValues({}, p), {
      buckets: __spreadProps(__spreadValues({}, p.buckets), {
        timeframe
      })
    }));
  }
  setSource(source) {
    this.updateParams((p) => __spreadProps(__spreadValues({}, p), {
      source: source ?? void 0
    }));
  }
  setBucketDuration(duration, multiplier) {
    this.updateParams((p) => __spreadProps(__spreadValues({}, p), {
      buckets: __spreadProps(__spreadValues({}, p.buckets), {
        duration,
        multiplier: multiplier ?? "AUTO"
      })
    }));
  }
  setBucketMultiplier(multiplier) {
    this.updateParams((p) => __spreadProps(__spreadValues({}, p), {
      buckets: __spreadProps(__spreadValues({}, p.buckets), {
        multiplier
      })
    }));
  }
  setEvent(event) {
    this.updateParams((p) => __spreadProps(__spreadValues({}, p), {
      event: event ?? void 0
    }));
  }
  setAutoRefreshInterval(autoRefreshInterval) {
    this.updateParams((p) => __spreadProps(__spreadValues({}, p), {
      autoRefresh: autoRefreshInterval
    }));
  }
  updateParams(fn) {
    this.paramsSubject.next(fn(this.params));
  }
  refresh() {
    this.variablesSubject.next(this.variablesSubject.getValue());
  }
  request(variables) {
    clearTimeout(this.refreshTimeout);
    this.loadingSubject.next(true);
    return this.apollo.query({
      query: TorrentMetricsDocument,
      variables,
      fetchPolicy: "no-cache"
    }).pipe(map((r) => {
      if (r) {
        this.loadingSubject.next(false);
        this.rawResultSubject.next(r.data);
      }
    })).pipe(catchError((err) => {
      this.errorsService.addError(`Failed to load torrent metrics: ${err.message}`);
      this.loadingSubject.next(false);
      this.setInterval();
      return EMPTY;
    })).subscribe();
  }
};
var createVariables = (params) => ({
  input: {
    bucketDuration: params.buckets.duration === "AUTO" ? "hour" : params.buckets.duration,
    sources: params.source ? [params.source] : void 0,
    startTime: new Date((/* @__PURE__ */ new Date()).getTime() - 1e3 * timeframeLengths[params.buckets.timeframe]).toISOString()
  }
});

// src/app/dashboard/torrents/torrent-chart-adapter.timeline.ts
var eventColors = {
  created: "primary",
  updated: "secondary"
};
var TorrentChartAdapterTimeline = class _TorrentChartAdapterTimeline {
  constructor() {
    this.themeInfo = inject(ThemeInfoService);
    this.transloco = inject(TranslocoService);
  }
  create(result, params) {
    const { colors } = this.themeInfo.info;
    const labels = Array();
    const datasets = [];
    if (result) {
      const nonEmptySources = result.sourceSummaries.filter((q) => !q.isEmpty);
      const nonEmptyBuckets = Array.from(new Set(nonEmptySources.flatMap((q) => q.events ? [q.events.earliestBucket, q.events.latestBucket] : []))).sort();
      const now = /* @__PURE__ */ new Date();
      const minBucket = Math.min(nonEmptyBuckets[0], normalizeBucket(now.getTime() - 1e3 * timeframeLengths[result.params.buckets.timeframe], result.params.buckets).index);
      const maxBucket = Math.max(nonEmptyBuckets[nonEmptyBuckets.length - 1], normalizeBucket(now, result.params.buckets).index);
      if (nonEmptyBuckets.length) {
        for (let i = minBucket; i <= maxBucket; i++) {
          labels.push(this.formatBucketKey(result.params.buckets, i));
        }
        const relevantEvents = eventNames.filter((n) => (result.params.event ?? n) === n);
        for (const source of nonEmptySources) {
          for (const event of relevantEvents) {
            const series = Array();
            for (let i = minBucket; i <= maxBucket; i++) {
              series.push(source.events?.eventBuckets?.[event]?.entries?.[`${i}`]?.count ?? 0);
            }
            datasets.push({
              yAxisID: "yCount",
              label: [source.source, event].join("/"),
              data: series,
              borderColor: colors[createThemeColor(eventColors[event], 50)],
              pointBackgroundColor: colors[createThemeColor(eventColors[event], 20)],
              pointBorderColor: colors[createThemeColor(eventColors[event], 80)],
              pointHoverBackgroundColor: colors[createThemeColor(eventColors[event], 40)],
              pointHoverBorderColor: colors[createThemeColor(eventColors[event], 60)]
            });
          }
        }
      }
    }
    return {
      type: "line",
      options: {
        animation: false,
        responsive: true,
        elements: {
          line: {
            tension: 0.5
          }
        },
        scales: {
          yCount: {
            position: "left",
            ticks: {
              callback: (v) => parseInt(v).toLocaleString(this.transloco.getActiveLang())
            }
          }
        },
        plugins: {
          legend: {
            display: params.legend
          },
          decimation: {
            enabled: true
          }
        }
      },
      data: {
        labels,
        datasets
      }
    };
  }
  formatBucketKey(params, key) {
    let formatStr;
    switch (params.duration) {
      case "day":
        formatStr = "d LLL";
        break;
      case "hour":
        formatStr = "d LLL H:00";
        break;
      case "minute":
        formatStr = "H:mm";
        break;
    }
    return format(1e3 * durationSeconds[params.duration] * params.multiplier * key, formatStr, {
      locale: resolveDateLocale(this.transloco.getActiveLang())
    });
  }
  static {
    this.\u0275fac = function TorrentChartAdapterTimeline_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _TorrentChartAdapterTimeline)();
    };
  }
  static {
    this.\u0275prov = /* @__PURE__ */ \u0275\u0275defineInjectable({ token: _TorrentChartAdapterTimeline, factory: _TorrentChartAdapterTimeline.\u0275fac, providedIn: "root" });
  }
};

// src/app/dashboard/torrents/torrent-metrics.component.ts
var _forTrack0 = ($index, $item) => $item.key;
var _c0 = () => ["dht"];
function TorrentMetricsComponent_ng_container_0_For_14_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 6);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const name_r3 = ctx.$implicit;
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("value", name_r3);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r4("dashboard.interval." + name_r3));
  }
}
function TorrentMetricsComponent_ng_container_0_For_41_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 6);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const name_r5 = ctx.$implicit;
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("value", name_r5);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r4("dashboard.interval." + name_r5 + "s"));
  }
}
function TorrentMetricsComponent_ng_container_0_For_73_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 6);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const source_r6 = ctx.$implicit;
    \u0275\u0275property("value", source_r6.key);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(source_r6.name);
  }
}
function TorrentMetricsComponent_ng_container_0_For_80_Template(rf, ctx) {
  if (rf & 1) {
    const _r7 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "button", 16);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_For_80_Template_button_click_0_listener() {
      const source_r8 = \u0275\u0275restoreView(_r7).$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.params.source === source_r8 || ctx_r1.torrentMetricsController.setSource(source_r8));
    });
    \u0275\u0275elementStart(1, "mat-icon");
    \u0275\u0275text(2);
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const source_r8 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275classMap(ctx_r1.torrentMetricsController.params.source === source_r8 ? "selected" : "deselected");
    \u0275\u0275property("matTooltip", source_r8);
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(ctx_r1.torrentMetricsController.params.source === source_r8 ? "radio_button_checked" : "radio_button_unchecked");
  }
}
function TorrentMetricsComponent_ng_container_0_For_93_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 6);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const event_r9 = ctx.$implicit;
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("value", event_r9);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r4("dashboard.event." + event_r9));
  }
}
function TorrentMetricsComponent_ng_container_0_For_114_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 6);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const name_r10 = ctx.$implicit;
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("value", name_r10);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r4("dashboard.interval." + name_r10));
  }
}
function TorrentMetricsComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-card")(2, "mat-card-content")(3, "mat-grid-list", 1)(4, "mat-grid-tile", 2)(5, "mat-card", 3)(6, "mat-card-header")(7, "mat-card-title")(8, "h4");
    \u0275\u0275text(9);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(10, "mat-card-content")(11, "mat-form-field", 4)(12, "mat-select", 5);
    \u0275\u0275listener("valueChange", function TorrentMetricsComponent_ng_container_0_Template_mat_select_valueChange_12_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setTimeframe($event));
    });
    \u0275\u0275repeaterCreate(13, TorrentMetricsComponent_ng_container_0_For_14_Template, 2, 2, "mat-option", 6, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(15, "div", 7)(16, "button", 8);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_16_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setTimeframe(ctx_r1.timeframeNames[0]));
    });
    \u0275\u0275elementStart(17, "mat-icon");
    \u0275\u0275text(18, "first_page");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(19, "button", 8);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_19_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setTimeframe(ctx_r1.timeframeNames[ctx_r1.timeframeNames.indexOf(ctx_r1.torrentMetricsController.params.buckets.timeframe) - 1]));
    });
    \u0275\u0275elementStart(20, "mat-icon");
    \u0275\u0275text(21, "navigate_before");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(22, "button", 8);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_22_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setTimeframe(ctx_r1.timeframeNames[ctx_r1.timeframeNames.indexOf(ctx_r1.torrentMetricsController.params.buckets.timeframe) + 1]));
    });
    \u0275\u0275elementStart(23, "mat-icon");
    \u0275\u0275text(24, "navigate_next");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(25, "button", 8);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_25_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setTimeframe(ctx_r1.timeframeNames[ctx_r1.timeframeNames.length - 1]));
    });
    \u0275\u0275elementStart(26, "mat-icon");
    \u0275\u0275text(27, "last_page");
    \u0275\u0275elementEnd()()()()()();
    \u0275\u0275elementStart(28, "mat-grid-tile", 2)(29, "mat-card", 9)(30, "mat-card-header")(31, "mat-card-title")(32, "h4");
    \u0275\u0275text(33);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(34, "mat-card-content")(35, "mat-form-field", 10)(36, "input", 11);
    \u0275\u0275pipe(37, "async");
    \u0275\u0275listener("change", function TorrentMetricsComponent_ng_container_0_Template_input_change_36_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.handleMultiplierEvent($event));
    });
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(38, "mat-form-field", 12)(39, "mat-select", 5);
    \u0275\u0275listener("valueChange", function TorrentMetricsComponent_ng_container_0_Template_mat_select_valueChange_39_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setBucketDuration($event));
    });
    \u0275\u0275repeaterCreate(40, TorrentMetricsComponent_ng_container_0_For_41_Template, 2, 2, "mat-option", 6, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(42, "div", 7)(43, "button", 8);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_43_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setBucketMultiplier(ctx_r1.torrentMetricsController.bucketMultiplier - 1));
    });
    \u0275\u0275elementStart(44, "mat-icon");
    \u0275\u0275text(45, "remove");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(46, "button", 13);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_46_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setBucketMultiplier(ctx_r1.torrentMetricsController.bucketMultiplier + 1));
    });
    \u0275\u0275elementStart(47, "mat-icon");
    \u0275\u0275text(48, "add");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(49, "button", 8);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_49_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setBucketDuration(ctx_r1.resolutionNames[0]));
    });
    \u0275\u0275elementStart(50, "mat-icon");
    \u0275\u0275text(51, "first_page");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(52, "button", 8);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_52_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setBucketDuration(ctx_r1.resolutionNames[ctx_r1.resolutionNames.indexOf(ctx_r1.torrentMetricsController.bucketDuration) - 1]));
    });
    \u0275\u0275elementStart(53, "mat-icon");
    \u0275\u0275text(54, "navigate_before");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(55, "button", 8);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_55_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setBucketDuration(ctx_r1.resolutionNames[ctx_r1.resolutionNames.indexOf(ctx_r1.torrentMetricsController.bucketDuration) + 1]));
    });
    \u0275\u0275elementStart(56, "mat-icon");
    \u0275\u0275text(57, "navigate_next");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(58, "button", 8);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_58_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setBucketDuration(ctx_r1.resolutionNames[ctx_r1.resolutionNames.length - 1]));
    });
    \u0275\u0275elementStart(59, "mat-icon");
    \u0275\u0275text(60, "last_page");
    \u0275\u0275elementEnd()()()()()();
    \u0275\u0275elementStart(61, "mat-grid-tile", 2)(62, "mat-card")(63, "mat-card-header")(64, "mat-card-title")(65, "h4");
    \u0275\u0275text(66);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(67, "mat-card-content")(68, "mat-form-field", 4)(69, "mat-select", 5);
    \u0275\u0275listener("valueChange", function TorrentMetricsComponent_ng_container_0_Template_mat_select_valueChange_69_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setSource($event === "_all" ? null : $event));
    });
    \u0275\u0275elementStart(70, "mat-option", 14);
    \u0275\u0275text(71, "All");
    \u0275\u0275elementEnd();
    \u0275\u0275repeaterCreate(72, TorrentMetricsComponent_ng_container_0_For_73_Template, 2, 2, "mat-option", 6, _forTrack0);
    \u0275\u0275pipe(74, "async");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(75, "div", 15)(76, "button", 16);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_76_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setSource(null));
    });
    \u0275\u0275elementStart(77, "mat-icon", 17);
    \u0275\u0275text(78, "workspaces");
    \u0275\u0275elementEnd()();
    \u0275\u0275repeaterCreate(79, TorrentMetricsComponent_ng_container_0_For_80_Template, 3, 4, "button", 18, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd()()()();
    \u0275\u0275elementStart(81, "mat-grid-tile", 2)(82, "mat-card")(83, "mat-card-header")(84, "mat-card-title")(85, "h4");
    \u0275\u0275text(86);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(87, "mat-card-content")(88, "mat-form-field", 4)(89, "mat-select", 5);
    \u0275\u0275listener("valueChange", function TorrentMetricsComponent_ng_container_0_Template_mat_select_valueChange_89_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setEvent($event === "_all" ? null : $event));
    });
    \u0275\u0275elementStart(90, "mat-option", 14);
    \u0275\u0275text(91, "All");
    \u0275\u0275elementEnd();
    \u0275\u0275repeaterCreate(92, TorrentMetricsComponent_ng_container_0_For_93_Template, 2, 2, "mat-option", 6, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(94, "div", 15)(95, "button", 16);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_95_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setEvent(null));
    });
    \u0275\u0275elementStart(96, "mat-icon", 17);
    \u0275\u0275text(97, "radio_button_checked");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(98, "button", 16);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_98_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.params.event === "created" || ctx_r1.torrentMetricsController.setEvent("created"));
    });
    \u0275\u0275elementStart(99, "mat-icon");
    \u0275\u0275text(100, "add_circle");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(101, "button", 16);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_101_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.params.event === "updated" || ctx_r1.torrentMetricsController.setEvent("updated"));
    });
    \u0275\u0275elementStart(102, "mat-icon");
    \u0275\u0275text(103, "check_circle");
    \u0275\u0275elementEnd()()()()()();
    \u0275\u0275elementStart(104, "mat-grid-tile", 2)(105, "mat-card", 19)(106, "mat-card-header")(107, "mat-card-title")(108, "h4");
    \u0275\u0275text(109);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(110, "mat-card-content")(111, "mat-form-field", 4)(112, "mat-select", 5);
    \u0275\u0275listener("valueChange", function TorrentMetricsComponent_ng_container_0_Template_mat_select_valueChange_112_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.setAutoRefreshInterval($event));
    });
    \u0275\u0275repeaterCreate(113, TorrentMetricsComponent_ng_container_0_For_114_Template, 2, 2, "mat-option", 6, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(115, "div", 15)(116, "button", 16);
    \u0275\u0275listener("click", function TorrentMetricsComponent_ng_container_0_Template_button_click_116_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.torrentMetricsController.refresh());
    });
    \u0275\u0275elementStart(117, "mat-icon");
    \u0275\u0275text(118, "sync");
    \u0275\u0275elementEnd()()()()()()();
    \u0275\u0275elementStart(119, "div", 20);
    \u0275\u0275element(120, "mat-progress-bar", 21);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(121, "mat-grid-list", 1)(122, "mat-grid-tile", 2);
    \u0275\u0275element(123, "app-chart", 22);
    \u0275\u0275elementEnd();
    \u0275\u0275element(124, "mat-grid-tile", 2);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    let tmp_15_0;
    let tmp_27_0;
    let tmp_28_0;
    let tmp_35_0;
    const t_r4 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance(3);
    \u0275\u0275property("cols", ctx_r1.breakpoints.sizeAtLeast("Large") ? 5 : ctx_r1.breakpoints.sizeAtLeast("Medium") ? 3 : ctx_r1.breakpoints.sizeAtLeast("Small") ? 2 : 1);
    \u0275\u0275advance();
    \u0275\u0275property("colspan", 1)("rowspan", 2);
    \u0275\u0275advance(5);
    \u0275\u0275textInterpolate(t_r4("dashboard.metrics.timeframe"));
    \u0275\u0275advance(3);
    \u0275\u0275property("value", ctx_r1.torrentMetricsController.params.buckets.timeframe);
    \u0275\u0275advance();
    \u0275\u0275repeater(ctx_r1.timeframeNames);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", ctx_r1.timeframeNames.indexOf(ctx_r1.torrentMetricsController.params.buckets.timeframe) <= 0);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", ctx_r1.timeframeNames.indexOf(ctx_r1.torrentMetricsController.params.buckets.timeframe) <= 0);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", ctx_r1.timeframeNames.indexOf(ctx_r1.torrentMetricsController.params.buckets.timeframe) >= ctx_r1.timeframeNames.length - 1);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", ctx_r1.timeframeNames.indexOf(ctx_r1.torrentMetricsController.params.buckets.timeframe) >= ctx_r1.timeframeNames.length - 1);
    \u0275\u0275advance(3);
    \u0275\u0275property("colspan", 1)("rowspan", 2);
    \u0275\u0275advance(5);
    \u0275\u0275textInterpolate1(" ", t_r4("dashboard.metrics.resolution"), " ");
    \u0275\u0275advance(3);
    \u0275\u0275property("placeholder", (tmp_15_0 = (tmp_15_0 = \u0275\u0275pipeBind1(37, 57, ctx_r1.torrentMetricsController.result$)) == null ? null : tmp_15_0.params == null ? null : tmp_15_0.params.buckets == null ? null : tmp_15_0.params.buckets.multiplier == null ? null : tmp_15_0.params.buckets.multiplier.toString()) !== null && tmp_15_0 !== void 0 ? tmp_15_0 : "")("value", ctx_r1.torrentMetricsController.params.buckets.multiplier);
    \u0275\u0275advance(3);
    \u0275\u0275property("value", ctx_r1.torrentMetricsController.bucketDuration);
    \u0275\u0275advance();
    \u0275\u0275repeater(ctx_r1.resolutionNames);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", ctx_r1.torrentMetricsController.bucketMultiplier === 1);
    \u0275\u0275advance(6);
    \u0275\u0275property("disabled", ctx_r1.resolutionNames.indexOf(ctx_r1.torrentMetricsController.bucketDuration) <= 0);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", ctx_r1.resolutionNames.indexOf(ctx_r1.torrentMetricsController.bucketDuration) <= 0);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", ctx_r1.resolutionNames.indexOf(ctx_r1.torrentMetricsController.bucketDuration) >= ctx_r1.resolutionNames.length - 1);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", ctx_r1.resolutionNames.indexOf(ctx_r1.torrentMetricsController.bucketDuration) >= ctx_r1.resolutionNames.length - 1);
    \u0275\u0275advance(3);
    \u0275\u0275property("colspan", 1)("rowspan", 2);
    \u0275\u0275advance(5);
    \u0275\u0275textInterpolate(t_r4("torrents.source"));
    \u0275\u0275advance(3);
    \u0275\u0275property("value", (tmp_27_0 = ctx_r1.torrentMetricsController.params.source) !== null && tmp_27_0 !== void 0 ? tmp_27_0 : "_all");
    \u0275\u0275advance(3);
    \u0275\u0275repeater((tmp_28_0 = \u0275\u0275pipeBind1(74, 59, ctx_r1.torrentMetricsController.result$)) == null ? null : tmp_28_0.availableSources);
    \u0275\u0275advance(4);
    \u0275\u0275classMap(ctx_r1.torrentMetricsController.params.source ? "deselected" : "selected");
    \u0275\u0275property("matTooltip", "all");
    \u0275\u0275advance(3);
    \u0275\u0275repeater(\u0275\u0275pureFunction0(61, _c0));
    \u0275\u0275advance(2);
    \u0275\u0275property("colspan", 1)("rowspan", 2);
    \u0275\u0275advance(5);
    \u0275\u0275textInterpolate(t_r4("dashboard.metrics.event"));
    \u0275\u0275advance(3);
    \u0275\u0275property("value", (tmp_35_0 = ctx_r1.torrentMetricsController.params.event) !== null && tmp_35_0 !== void 0 ? tmp_35_0 : "_all");
    \u0275\u0275advance(3);
    \u0275\u0275repeater(ctx_r1.eventNames);
    \u0275\u0275advance(3);
    \u0275\u0275classMap(!ctx_r1.torrentMetricsController.params.event ? "selected" : "deselected");
    \u0275\u0275property("matTooltip", "all");
    \u0275\u0275advance(3);
    \u0275\u0275classMap(ctx_r1.torrentMetricsController.params.event === "created" ? "selected" : "deselected");
    \u0275\u0275property("matTooltip", "created");
    \u0275\u0275advance(3);
    \u0275\u0275classMap(ctx_r1.torrentMetricsController.params.event === "updated" ? "selected" : "deselected");
    \u0275\u0275property("matTooltip", "updated");
    \u0275\u0275advance(3);
    \u0275\u0275property("colspan", 1)("rowspan", 2);
    \u0275\u0275advance(5);
    \u0275\u0275textInterpolate(t_r4("general.refresh"));
    \u0275\u0275advance(3);
    \u0275\u0275property("value", ctx_r1.torrentMetricsController.params.autoRefresh);
    \u0275\u0275advance();
    \u0275\u0275repeater(ctx_r1.autoRefreshIntervalNames);
    \u0275\u0275advance(3);
    \u0275\u0275property("matTooltip", "Refresh");
    \u0275\u0275advance(4);
    \u0275\u0275property("mode", ctx_r1.torrentMetricsController.loading ? "indeterminate" : "determinate")("value", 0);
    \u0275\u0275advance();
    \u0275\u0275property("cols", ctx_r1.breakpoints.sizeAtLeast("Large") ? 2 : 1);
    \u0275\u0275advance();
    \u0275\u0275property("colspan", 1)("rowspan", 6);
    \u0275\u0275advance();
    \u0275\u0275property("title", t_r4("dashboard.metrics.throughput"))("adapter", ctx_r1.timeline)("$data", ctx_r1.torrentMetricsController.result$)("height", 400)("width", 550);
    \u0275\u0275advance();
    \u0275\u0275property("colspan", 1)("rowspan", 5);
  }
}
var TorrentMetricsComponent = class _TorrentMetricsComponent {
  constructor() {
    this.breakpoints = inject(BreakpointsService);
    this.apollo = inject(Apollo);
    this.torrentMetricsController = new TorrentMetricsController(this.apollo, {
      buckets: defaultBucketParams,
      autoRefresh: "seconds_30"
    }, inject(ErrorsService));
    this.timeline = inject(TorrentChartAdapterTimeline);
    this.resolutionNames = resolutionNames;
    this.timeframeNames = timeframeNames;
    this.autoRefreshIntervalNames = autoRefreshIntervalNames;
    this.eventNames = eventNames;
  }
  ngOnDestroy() {
    this.torrentMetricsController.setAutoRefreshInterval("off");
  }
  handleMultiplierEvent(event) {
    const value = event.currentTarget.value;
    this.torrentMetricsController.setBucketMultiplier(/^\d+$/.test(value) ? parseInt(value) : "AUTO");
  }
  static {
    this.\u0275fac = function TorrentMetricsComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _TorrentMetricsComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _TorrentMetricsComponent, selectors: [["app-torrent-metrics"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], ["rowHeight", "100px", 3, "cols"], [3, "colspan", "rowspan"], [1, "form-timeframe"], ["subscriptSizing", "dynamic"], [3, "valueChange", "value"], [3, "value"], [1, "paginator", "actions"], ["mat-icon-button", "", 3, "click", "disabled"], [1, "form-resolution"], ["subscriptSizing", "dynamic", 1, "form-input-multiplier"], ["type", "number", "matInput", "", "min", "1", "step", "1", 3, "change", "placeholder", "value"], ["subscriptSizing", "dynamic", 1, "form-select-duration"], ["mat-icon-button", "", 3, "click"], ["value", "_all"], [1, "actions"], ["mat-icon-button", "", 3, "click", "matTooltip"], ["fontSet", "material-icons"], ["mat-icon-button", "", 3, "class", "matTooltip"], [1, "form-refresh"], [1, "progress-bar-container"], [3, "mode", "value"], [3, "title", "adapter", "$data", "height", "width"]], template: function TorrentMetricsComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, TorrentMetricsComponent_ng_container_0_Template, 125, 62, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatOption, MatIconButton, MatCard, MatCardContent, MatCardHeader, MatCardTitle, MatFormField, MatGridList, MatGridTile, MatIcon, MatInput, MatProgressBar, MatSelect, MatTooltip, TranslocoDirective, AsyncPipe, ChartComponent, GraphQLModule], styles: ["\n\n.actions[_ngcontent-%COMP%] {\n  width: 210px;\n  padding-top: 12px;\n  --mdc-icon-button-state-layer-size: 32px;\n}\n.actions[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  font-size: 22px;\n}\n.actions[_ngcontent-%COMP%]   button[_ngcontent-%COMP%] {\n  margin-right: 0;\n}\n.progress-bar-container[_ngcontent-%COMP%] {\n  width: 100%;\n  height: 10px;\n}\nmat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%] {\n  width: 100%;\n}\nmat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]   mat-card-content[_ngcontent-%COMP%] {\n  min-width: 190px;\n}\nmat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]   h4[_ngcontent-%COMP%] {\n  margin-bottom: 16px;\n  font-size: 18px;\n}\nmat-form-field[_ngcontent-%COMP%] {\n  width: 186px;\n}\n.form-resolution[_ngcontent-%COMP%]   .actions[_ngcontent-%COMP%] {\n  margin-left: -2px;\n}\n.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%] {\n  width: 60px;\n  margin-right: 10px;\n}\n.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[_ngcontent-%COMP%]::-webkit-outer-spin-button, \n.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[_ngcontent-%COMP%]::-webkit-inner-spin-button {\n  -webkit-appearance: none;\n  margin: 0;\n}\n.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[type=number][_ngcontent-%COMP%] {\n  -moz-appearance: textfield;\n}\n.form-resolution[_ngcontent-%COMP%]   .form-select-duration[_ngcontent-%COMP%] {\n  width: 116px;\n}\n/*# sourceMappingURL=torrent-metrics.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(TorrentMetricsComponent, { className: "TorrentMetricsComponent", filePath: "src/app/dashboard/torrents/torrent-metrics.component.ts", lineNumber: 25 });
})();

// src/app/dashboard/torrents/torrents-dashboard.component.ts
var _c02 = (a0, a1) => [a0, a1];
function TorrentsDashboardComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275element(1, "app-document-title", 1);
    \u0275\u0275elementStart(2, "mat-card", 2)(3, "mat-card-header")(4, "mat-toolbar")(5, "h2");
    \u0275\u0275element(6, "mat-icon", 3);
    \u0275\u0275text(7);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(8, "mat-card-content");
    \u0275\u0275element(9, "app-torrent-metrics");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r1 = ctx.$implicit;
    \u0275\u0275advance();
    \u0275\u0275property("parts", \u0275\u0275pureFunction2(2, _c02, t_r1("routes.torrents"), t_r1("routes.dashboard")));
    \u0275\u0275advance(6);
    \u0275\u0275textInterpolate(t_r1("routes.torrents"));
  }
}
var TorrentsDashboardComponent = class _TorrentsDashboardComponent {
  static {
    this.\u0275fac = function TorrentsDashboardComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _TorrentsDashboardComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _TorrentsDashboardComponent, selectors: [["app-torrents"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], [3, "parts"], [1, "dashboard-card"], ["svgIcon", "magnet"]], template: function TorrentsDashboardComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, TorrentsDashboardComponent_ng_container_0_Template, 10, 5, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatCard, MatCardContent, MatCardHeader, MatIcon, MatToolbar, TranslocoDirective, TorrentMetricsComponent, DocumentTitleComponent], styles: ["\n\nmat-card-header[_ngcontent-%COMP%] {\n  flex-wrap: wrap;\n}\nmat-card-header[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%] {\n  font-size: 18px;\n  margin: 0 60px 0 48px;\n  height: 48px;\n  line-height: 48px;\n}\nmat-card-header[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  position: relative;\n  top: 6px;\n  margin-right: 14px;\n  line-height: 1.25rem;\n}\nmat-card-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%] {\n  flex: 0 0 100%;\n}\nmat-card-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   a[_ngcontent-%COMP%] {\n  margin-top: 2px;\n}\nmat-card-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   a[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  margin-right: 12px;\n}\n/*# sourceMappingURL=torrents-dashboard.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(TorrentsDashboardComponent, { className: "TorrentsDashboardComponent", filePath: "src/app/dashboard/torrents/torrents-dashboard.component.ts", lineNumber: 13 });
})();
export {
  TorrentsDashboardComponent
};
//# sourceMappingURL=chunk-YZYDLAWD.js.map
