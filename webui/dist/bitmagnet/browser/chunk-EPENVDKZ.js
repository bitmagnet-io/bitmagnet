import {
  TorrentsBulkActionsComponent,
  TorrentsTableComponent,
  allColumns,
  compactColumns
} from "./chunk-QLUBXXM3.js";
import {
  TorrentsSearchController,
  defaultOrderBy,
  facets,
  inactiveFacet,
  isDefaultOrdering,
  orderByOptions,
  torrentTabNames
} from "./chunk-LUZJBAO3.js";
import {
  IntEstimatePipe,
  PaginatorComponent
} from "./chunk-43HRGFU3.js";
import "./chunk-ORIQXXAG.js";
import {
  contentTypeList,
  contentTypeMap
} from "./chunk-UGVUNZOV.js";
import "./chunk-3D6CEWET.js";
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
  DefaultValueAccessor,
  FormControl,
  FormControlDirective,
  GraphQLModule,
  MatCheckbox,
  MatDivider,
  MatDrawer,
  MatDrawerContainer,
  MatDrawerContent,
  MatExpansionPanel,
  MatExpansionPanelHeader,
  MatExpansionPanelTitle,
  MatFormField,
  MatIcon,
  MatIconButton,
  MatInput,
  MatLabel,
  MatMiniFabButton,
  MatOption,
  MatSelect,
  MatTooltip,
  NgControlStatus,
  SelectionModel,
  TorrentContentSearchDocument,
  TranslocoDirective,
  TranslocoService
} from "./chunk-WWRDQTKJ.js";
import "./chunk-VSVMRYN2.js";
import {
  ActivatedRoute,
  Router
} from "./chunk-Y2ZC5Z2X.js";
import {
  AsyncPipe,
  BehaviorSubject,
  EMPTY,
  __spreadProps,
  __spreadValues,
  catchError,
  combineLatestWith,
  inject,
  map,
  scan,
  ɵsetClassDebugInfo,
  ɵɵStandaloneFeature,
  ɵɵadvance,
  ɵɵattribute,
  ɵɵclassMap,
  ɵɵconditional,
  ɵɵdefineComponent,
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
  ɵɵpipeBind2,
  ɵɵproperty,
  ɵɵpureFunction1,
  ɵɵpureFunction4,
  ɵɵreference,
  ɵɵrepeater,
  ɵɵrepeaterCreate,
  ɵɵresetView,
  ɵɵrestoreView,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate,
  ɵɵtextInterpolate1
} from "./chunk-DMMUMX3A.js";

// src/app/util/query-string.ts
var stringListParam = (params, key) => {
  const str = stringParam(params, key);
  const list = str?.split(",").map((str2) => str2.trim()).filter(Boolean);
  return list?.length ? Array.from(new Set(list)).sort() : void 0;
};
var stringParam = (params, key) => {
  return typeof params[key] === "string" ? decodeURIComponent(params[key]) || void 0 : void 0;
};
var intParam = (params, key) => {
  if (params && params[key] && /^\d+$/.test(params[key])) {
    return parseInt(params[key]);
  }
  return void 0;
};

// src/app/torrents/torrents-search.datasource.ts
var emptyResult = {
  items: [],
  totalCount: 0,
  totalCountIsEstimate: false,
  aggregations: {}
};
var TorrentsSearchDatasource = class {
  constructor(apollo, errorsService, searchQueryVariables) {
    this.apollo = apollo;
    this.errorsService = errorsService;
    this.currentRequest = new BehaviorSubject(0);
    this.loadingSubject = new BehaviorSubject(false);
    this.loading$ = this.loadingSubject.asObservable();
    this.result = emptyResult;
    this.resultSubject = new BehaviorSubject(this.result);
    this.result$ = this.resultSubject.asObservable();
    this.items$ = this.resultSubject.pipe(map((result) => result.items));
    this.overallTotalCount$ = this.resultSubject.pipe(map((result) => {
      let overallTotalCount = 0;
      let overallIsEstimate = false;
      for (const ct of result.aggregations.contentType ?? []) {
        overallTotalCount += ct.count;
        overallIsEstimate = overallIsEstimate || ct.isEstimate;
      }
      return {
        count: overallTotalCount,
        isEstimate: overallIsEstimate
      };
    }));
    this.availableContentTypes$ = this.resultSubject.pipe(scan((acc, next) => Array.from(/* @__PURE__ */ new Set([
      ...acc,
      ...(next.aggregations.contentType ?? []).flatMap((agg) => agg.value ? [agg.value] : [])
    ])), []));
    this.contentTypeCounts$ = this.resultSubject.pipe(map((result) => Object.fromEntries((result.aggregations.contentType ?? []).map((ct) => [
      ct.value,
      {
        count: ct.count,
        isEstimate: ct.isEstimate
      }
    ]))));
    searchQueryVariables.subscribe((variables) => {
      this.input = variables.input;
      this.loadResult({
        input: __spreadProps(__spreadValues({}, variables.input), {
          cached: true
        })
      });
    });
    this.resultSubject.subscribe((result) => {
      this.result = result;
    });
  }
  connect({}) {
    return this.items$;
  }
  disconnect() {
    this.resultSubject.complete();
  }
  refresh() {
    this.loadResult({
      input: __spreadProps(__spreadValues({}, this.input), {
        cached: false
      })
    });
  }
  loadResult(variables) {
    if (this.currentSubscription) {
      this.currentSubscription.unsubscribe();
      this.currentSubscription = void 0;
    }
    this.loadingSubject.next(true);
    const currentRequest = this.currentRequest.getValue() + 1;
    this.currentRequest.next(currentRequest);
    const result = this.apollo.query({
      query: TorrentContentSearchDocument,
      variables,
      fetchPolicy: "no-cache"
    }).pipe(map((r) => r.data.torrentContent.search)).pipe(catchError((err) => {
      this.errorsService.addError(`Error loading item results: ${err.message}`);
      return EMPTY;
    }));
    this.currentSubscription = result.subscribe((r) => {
      if (currentRequest === this.currentRequest.getValue()) {
        this.loadingSubject.next(false);
        this.resultSubject.next(r);
      }
    });
  }
};

// src/app/torrents/torrents-search.component.ts
var _forTrack0 = ($index, $item) => $item.key;
var _forTrack1 = ($index, $item) => $item.field;
var _forTrack2 = ($index, $item) => $item.value;
var _c0 = (a0) => ({ x: a0 });
var _c1 = (a0, a1, a2, a3) => [a0, a1, a2, a3];
function TorrentsSearchComponent_ng_container_0_Conditional_18_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "small");
    \u0275\u0275text(1);
    \u0275\u0275pipe(2, "intEstimate");
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const count_r3 = ctx;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", \u0275\u0275pipeBind2(2, 1, count_r3.count, count_r3.isEstimate), " ");
  }
}
function TorrentsSearchComponent_ng_container_0_For_21_Conditional_0_Conditional_4_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "small");
    \u0275\u0275text(1);
    \u0275\u0275pipe(2, "intEstimate");
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const agg_r6 = ctx;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(\u0275\u0275pipeBind2(2, 1, agg_r6.count, agg_r6.isEstimate));
  }
}
function TorrentsSearchComponent_ng_container_0_For_21_Conditional_0_Conditional_6_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "small");
    \u0275\u0275text(1, "0");
    \u0275\u0275elementEnd();
  }
}
function TorrentsSearchComponent_ng_container_0_For_21_Conditional_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r4 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "li", 6);
    \u0275\u0275listener("click", function TorrentsSearchComponent_ng_container_0_For_21_Conditional_0_Template_li_click_0_listener() {
      \u0275\u0275restoreView(_r4);
      const ct_r5 = \u0275\u0275nextContext().$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.controller.selectContentType(ct_r5.key));
    });
    \u0275\u0275elementStart(1, "mat-icon");
    \u0275\u0275text(2);
    \u0275\u0275elementEnd();
    \u0275\u0275text(3);
    \u0275\u0275template(4, TorrentsSearchComponent_ng_container_0_For_21_Conditional_0_Conditional_4_Template, 3, 4, "small");
    \u0275\u0275pipe(5, "async");
    \u0275\u0275template(6, TorrentsSearchComponent_ng_container_0_For_21_Conditional_0_Conditional_6_Template, 2, 0, "small");
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    let tmp_17_0;
    const ct_r5 = \u0275\u0275nextContext().$implicit;
    const t_r7 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275classMap(ctx_r1.controls.contentType === ct_r5.key ? "active" : "");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(ct_r5.icon);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r7("content_types.plural." + ct_r5.key), " ");
    \u0275\u0275advance();
    \u0275\u0275conditional((tmp_17_0 = (tmp_17_0 = \u0275\u0275pipeBind1(5, 5, ctx_r1.dataSource.contentTypeCounts$)) == null ? null : tmp_17_0[ct_r5.key]) ? 4 : 6, tmp_17_0);
  }
}
function TorrentsSearchComponent_ng_container_0_For_21_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275template(0, TorrentsSearchComponent_ng_container_0_For_21_Conditional_0_Template, 7, 7, "li", 24);
    \u0275\u0275pipe(1, "async");
  }
  if (rf & 2) {
    let tmp_13_0;
    const ct_r5 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275conditional(ct_r5.key === "null" || ((tmp_13_0 = \u0275\u0275pipeBind1(1, 1, ctx_r1.dataSource.availableContentTypes$)) == null ? null : tmp_13_0.includes(ct_r5.key)) ? 0 : -1);
  }
}
function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_6_For_2_Template(rf, ctx) {
  if (rf & 1) {
    const _r10 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-checkbox", 30);
    \u0275\u0275listener("change", function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_6_For_2_Template_mat_checkbox_change_0_listener($event) {
      const agg_r11 = \u0275\u0275restoreView(_r10).$implicit;
      const facet_r9 = \u0275\u0275nextContext(3).$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView($event.checked ? ctx_r1.controller.activateFilter(facet_r9, agg_r11.value) : ctx_r1.controller.deactivateFilter(facet_r9, agg_r11.value));
    });
    \u0275\u0275text(1);
    \u0275\u0275elementStart(2, "small");
    \u0275\u0275text(3);
    \u0275\u0275pipe(4, "intEstimate");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const agg_r11 = ctx.$implicit;
    const facet_r9 = \u0275\u0275nextContext(3).$implicit;
    \u0275\u0275property("checked", facet_r9.filter == null ? null : facet_r9.filter.includes(agg_r11.value));
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", agg_r11.label, " ");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(\u0275\u0275pipeBind2(4, 3, agg_r11.count, agg_r11.isEstimate));
  }
}
function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_6_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "section", 27);
    \u0275\u0275repeaterCreate(1, TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_6_For_2_Template, 5, 6, "mat-checkbox", 29, _forTrack2);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const facet_r9 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275repeater(facet_r9.aggregations);
  }
}
function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_7_Conditional_1_For_1_Template(rf, ctx) {
  if (rf & 1) {
    const _r12 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-checkbox", 32);
    \u0275\u0275listener("change", function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_7_Conditional_1_For_1_Template_mat_checkbox_change_0_listener() {
      const agg_r13 = \u0275\u0275restoreView(_r12).$implicit;
      const facet_r9 = \u0275\u0275nextContext(4).$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.controller.activateFilter(facet_r9, agg_r13.value));
    });
    \u0275\u0275text(1);
    \u0275\u0275elementStart(2, "small");
    \u0275\u0275text(3);
    \u0275\u0275pipe(4, "intEstimate");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const agg_r13 = ctx.$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", agg_r13.label, " ");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(\u0275\u0275pipeBind2(4, 2, agg_r13.count, agg_r13.isEstimate));
  }
}
function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_7_Conditional_1_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275repeaterCreate(0, TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_7_Conditional_1_For_1_Template, 5, 5, "mat-checkbox", 31, _forTrack2);
  }
  if (rf & 2) {
    const facet_r9 = \u0275\u0275nextContext(3).$implicit;
    \u0275\u0275repeater(facet_r9.aggregations);
  }
}
function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_7_Conditional_2_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275text(0);
  }
  if (rf & 2) {
    const t_r7 = \u0275\u0275nextContext(4).$implicit;
    \u0275\u0275textInterpolate1(" ", t_r7("general.none"), " ");
  }
}
function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_7_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "section", 28);
    \u0275\u0275template(1, TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_7_Conditional_1_Template, 2, 0)(2, TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_7_Conditional_2_Template, 1, 1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const facet_r9 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275conditional(facet_r9.aggregations.length ? 1 : 2);
  }
}
function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r8 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-expansion-panel", 26);
    \u0275\u0275listener("opened", function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Template_mat_expansion_panel_opened_0_listener() {
      \u0275\u0275restoreView(_r8);
      const facet_r9 = \u0275\u0275nextContext().$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.controller.activateFacet(facet_r9));
    })("closed", function TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Template_mat_expansion_panel_closed_0_listener() {
      \u0275\u0275restoreView(_r8);
      const facet_r9 = \u0275\u0275nextContext().$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.controller.deactivateFacet(facet_r9));
    });
    \u0275\u0275elementStart(1, "mat-expansion-panel-header")(2, "mat-panel-title")(3, "mat-icon");
    \u0275\u0275text(4);
    \u0275\u0275elementEnd();
    \u0275\u0275text(5);
    \u0275\u0275elementEnd()();
    \u0275\u0275template(6, TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_6_Template, 3, 0, "section", 27)(7, TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Conditional_7_Template, 3, 1, "section", 28);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const facet_r9 = \u0275\u0275nextContext().$implicit;
    const t_r7 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("expanded", facet_r9.active);
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(facet_r9.icon);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r7("facets." + facet_r9.key), " ");
    \u0275\u0275advance();
    \u0275\u0275conditional((facet_r9.filter == null ? null : facet_r9.filter.length) ? 6 : 7);
  }
}
function TorrentsSearchComponent_ng_container_0_For_23_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275template(0, TorrentsSearchComponent_ng_container_0_For_23_Conditional_0_Template, 8, 4, "mat-expansion-panel", 25);
  }
  if (rf & 2) {
    const facet_r9 = ctx.$implicit;
    \u0275\u0275conditional(facet_r9.relevant ? 0 : -1);
  }
}
function TorrentsSearchComponent_ng_container_0_Conditional_34_Template(rf, ctx) {
  if (rf & 1) {
    const _r15 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "button", 18);
    \u0275\u0275listener("click", function TorrentsSearchComponent_ng_container_0_Conditional_34_Template_button_click_0_listener() {
      \u0275\u0275restoreView(_r15);
      const ctx_r1 = \u0275\u0275nextContext(2);
      ctx_r1.queryString.reset();
      return \u0275\u0275resetView(ctx_r1.controller.setQueryString(null));
    });
    \u0275\u0275elementStart(1, "mat-icon");
    \u0275\u0275text(2, "close");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const t_r7 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("matTooltip", t_r7("torrents.clear_search"));
  }
}
function TorrentsSearchComponent_ng_container_0_For_41_Conditional_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 33);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const option_r16 = \u0275\u0275nextContext().$implicit;
    const t_r7 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("value", option_r16.field);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r7("torrents.ordering." + option_r16.field), " ");
  }
}
function TorrentsSearchComponent_ng_container_0_For_41_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275template(0, TorrentsSearchComponent_ng_container_0_For_41_Conditional_0_Template, 2, 2, "mat-option", 33);
  }
  if (rf & 2) {
    const option_r16 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275conditional(option_r16.field != "relevance" || ctx_r1.queryString.value ? 0 : -1);
  }
}
function TorrentsSearchComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275element(1, "app-document-title", 2);
    \u0275\u0275elementStart(2, "mat-drawer-container", 3)(3, "mat-drawer", 4, 0)(5, "mat-expansion-panel", 5)(6, "mat-expansion-panel-header")(7, "mat-panel-title")(8, "mat-icon");
    \u0275\u0275text(9, "interests");
    \u0275\u0275elementEnd();
    \u0275\u0275text(10);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(11, "section")(12, "nav")(13, "ul")(14, "li", 6);
    \u0275\u0275listener("click", function TorrentsSearchComponent_ng_container_0_Template_li_click_14_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.controller.selectContentType(null));
    });
    \u0275\u0275elementStart(15, "mat-icon", 7);
    \u0275\u0275text(16, "emergency");
    \u0275\u0275elementEnd();
    \u0275\u0275text(17);
    \u0275\u0275template(18, TorrentsSearchComponent_ng_container_0_Conditional_18_Template, 3, 4, "small");
    \u0275\u0275pipe(19, "async");
    \u0275\u0275elementEnd();
    \u0275\u0275repeaterCreate(20, TorrentsSearchComponent_ng_container_0_For_21_Template, 2, 3, null, null, _forTrack0);
    \u0275\u0275elementEnd()()()();
    \u0275\u0275repeaterCreate(22, TorrentsSearchComponent_ng_container_0_For_23_Template, 1, 1, null, null, _forTrack0);
    \u0275\u0275pipe(24, "async");
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(25, "mat-drawer-content")(26, "div", 8)(27, "div", 9)(28, "button", 10);
    \u0275\u0275listener("click", function TorrentsSearchComponent_ng_container_0_Template_button_click_28_listener() {
      \u0275\u0275restoreView(_r1);
      const drawer_r14 = \u0275\u0275reference(4);
      return \u0275\u0275resetView(drawer_r14.toggle());
    });
    \u0275\u0275elementStart(29, "mat-icon", 11);
    \u0275\u0275text(30);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(31, "div", 12)(32, "mat-form-field", 13)(33, "input", 14);
    \u0275\u0275listener("keyup.enter", function TorrentsSearchComponent_ng_container_0_Template_input_keyup_enter_33_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.controller.setQueryString(ctx_r1.queryString.value));
    });
    \u0275\u0275elementEnd();
    \u0275\u0275template(34, TorrentsSearchComponent_ng_container_0_Conditional_34_Template, 3, 1, "button", 15);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(35, "div", 16)(36, "mat-form-field", 13)(37, "mat-label");
    \u0275\u0275text(38);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(39, "mat-select", 17);
    \u0275\u0275listener("valueChange", function TorrentsSearchComponent_ng_container_0_Template_mat_select_valueChange_39_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.controller.selectOrderBy($event));
    });
    \u0275\u0275repeaterCreate(40, TorrentsSearchComponent_ng_container_0_For_41_Template, 1, 1, null, null, _forTrack1);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(42, "button", 18);
    \u0275\u0275listener("click", function TorrentsSearchComponent_ng_container_0_Template_button_click_42_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.controller.toggleOrderByDirection());
    });
    \u0275\u0275elementStart(43, "mat-icon");
    \u0275\u0275text(44);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(45, "div", 19)(46, "button", 20);
    \u0275\u0275listener("click", function TorrentsSearchComponent_ng_container_0_Template_button_click_46_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.dataSource.refresh());
    });
    \u0275\u0275elementStart(47, "mat-icon");
    \u0275\u0275text(48, "sync");
    \u0275\u0275elementEnd()()()();
    \u0275\u0275element(49, "mat-divider");
    \u0275\u0275elementStart(50, "app-torrents-bulk-actions", 21);
    \u0275\u0275listener("updated", function TorrentsSearchComponent_ng_container_0_Template_app_torrents_bulk_actions_updated_50_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.dataSource.refresh());
    });
    \u0275\u0275elementEnd();
    \u0275\u0275element(51, "mat-divider");
    \u0275\u0275elementStart(52, "app-torrents-table", 22);
    \u0275\u0275listener("updated", function TorrentsSearchComponent_ng_container_0_Template_app_torrents_table_updated_52_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.dataSource.refresh());
    });
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(53, "app-paginator", 23);
    \u0275\u0275listener("paging", function TorrentsSearchComponent_ng_container_0_Template_app_paginator_paging_53_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.controller.handlePageEvent($event));
    });
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    let tmp_3_0;
    let tmp_11_0;
    const t_r7 = ctx.$implicit;
    const drawer_r14 = \u0275\u0275reference(4);
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275property("parts", \u0275\u0275pureFunction4(37, _c1, ctx_r1.controls.queryString, ((tmp_3_0 = ctx_r1.controls.contentType) !== null && tmp_3_0 !== void 0 ? tmp_3_0 : "null") === "null" ? null : t_r7("content_types.plural." + ctx_r1.controls.contentType), ctx_r1.controls.page > 1 ? t_r7("paginator.page_x", \u0275\u0275pureFunction1(35, _c0, ctx_r1.controls.page)) : null, t_r7("routes.torrents")));
    \u0275\u0275advance(2);
    \u0275\u0275property("mode", ctx_r1.breakpoints.sizeAtLeast("Medium") ? "side" : "over")("opened", ctx_r1.breakpoints.sizeAtLeast("Medium"));
    \u0275\u0275attribute("role", ctx_r1.breakpoints.sizeAtLeast("Medium") ? "navigation" : "dialog");
    \u0275\u0275advance(2);
    \u0275\u0275property("expanded", ctx_r1.breakpoints.sizeAtLeast("Medium"));
    \u0275\u0275advance(5);
    \u0275\u0275textInterpolate1(" ", t_r7("facets.content_type"), " ");
    \u0275\u0275advance(4);
    \u0275\u0275classMap(ctx_r1.controls.contentType === null ? "active" : "");
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate1("", t_r7("content_types.plural.all"), " ");
    \u0275\u0275advance();
    \u0275\u0275conditional((tmp_11_0 = \u0275\u0275pipeBind1(19, 31, ctx_r1.dataSource.overallTotalCount$)) ? 18 : -1, tmp_11_0);
    \u0275\u0275advance(2);
    \u0275\u0275repeater(ctx_r1.contentTypes);
    \u0275\u0275advance(2);
    \u0275\u0275repeater(\u0275\u0275pipeBind1(24, 33, ctx_r1.facets$));
    \u0275\u0275advance(6);
    \u0275\u0275property("matTooltip", t_r7("torrents.toggle_drawer"));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(drawer_r14.opened ? "arrow_circle_left" : "arrow_circle_right");
    \u0275\u0275advance(3);
    \u0275\u0275property("placeholder", t_r7("torrents.search"))("formControl", ctx_r1.queryString);
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r1.queryString.value ? 34 : -1);
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(t_r7("torrents.order_by"));
    \u0275\u0275advance();
    \u0275\u0275property("value", ctx_r1.controls.orderBy.field);
    \u0275\u0275advance();
    \u0275\u0275repeater(ctx_r1.orderByOptions);
    \u0275\u0275advance(2);
    \u0275\u0275property("matTooltip", t_r7("torrents.order_direction_toggle"));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(ctx_r1.controls.orderBy.descending ? "arrow_downward" : "arrow_upward");
    \u0275\u0275advance(2);
    \u0275\u0275property("matTooltip", t_r7("torrents.refresh"));
    \u0275\u0275advance(4);
    \u0275\u0275property("selectedItems$", ctx_r1.selectedItems$);
    \u0275\u0275advance(2);
    \u0275\u0275property("dataSource", ctx_r1.dataSource)("controller", ctx_r1.controller)("displayedColumns", ctx_r1.breakpoints.sizeAtLeast("Medium") ? ctx_r1.allColumns : ctx_r1.compactColumns)("multiSelection", ctx_r1.multiSelection);
    \u0275\u0275advance();
    \u0275\u0275property("page", ctx_r1.controls.page)("pageSize", ctx_r1.controls.limit)("pageLength", ctx_r1.dataSource.result.items.length)("totalLength", ctx_r1.dataSource.result.totalCount)("totalIsEstimate", ctx_r1.dataSource.result.totalCountIsEstimate)("hasNextPage", ctx_r1.dataSource.result.hasNextPage);
  }
}
var TorrentsSearchComponent = class _TorrentsSearchComponent {
  constructor() {
    this.route = inject(ActivatedRoute);
    this.router = inject(Router);
    this.apollo = inject(Apollo);
    this.errorsService = inject(ErrorsService);
    this.transloco = inject(TranslocoService);
    this.breakpoints = inject(BreakpointsService);
    this.controls = initControls;
    this.contentTypes = contentTypeList;
    this.orderByOptions = orderByOptions;
    this.allColumns = allColumns;
    this.compactColumns = compactColumns;
    this.queryString = new FormControl("");
    this.result = emptyResult;
    this.multiSelection = new SelectionModel(true, []);
    this.selectedItemsSubject = new BehaviorSubject([]);
    this.selectedItems$ = this.selectedItemsSubject.asObservable();
    this.subscriptions = Array();
    this.controller = new TorrentsSearchController(this.controls);
    this.dataSource = new TorrentsSearchDatasource(this.apollo, this.errorsService, this.controller.params$);
    this.subscriptions.push(this.controller.controls$.subscribe((ctrl) => {
      this.controls = ctrl;
    }));
    this.facets$ = this.controller.controls$.pipe(combineLatestWith(this.dataSource.result$), map(([controls, result]) => facets.map((f) => __spreadProps(__spreadValues(__spreadValues({}, f), f.extractInput(controls.facets)), {
      relevant: !f.contentTypes || !!(controls.contentType && controls.contentType !== "null" && f.contentTypes.includes(controls.contentType)),
      aggregations: f.extractAggregations(result.aggregations).map((agg) => __spreadProps(__spreadValues({}, agg), {
        label: f.resolveLabel(agg, this.transloco)
      }))
    }))));
    this.subscriptions.push(this.dataSource.result$.subscribe((result) => {
      this.result = result;
      const infoHashes = new Set(result.items.map(({ infoHash }) => infoHash));
      this.multiSelection.deselect(...this.multiSelection.selected.filter((infoHash) => !infoHashes.has(infoHash)));
    }));
  }
  ngOnInit() {
    this.subscriptions.push(this.route.queryParams.subscribe((params) => {
      this.queryString.setValue(stringParam(params, "query") ?? null);
      this.controller.update(() => paramsToControls(params));
    }), this.controller.controls$.subscribe((ctrl) => {
      void this.router.navigate([], {
        relativeTo: this.route,
        queryParams: controlsToParams(ctrl),
        queryParamsHandling: "replace"
      });
    }), this.multiSelection.changed.subscribe((selection) => {
      const infoHashes = new Set(selection.source.selected);
      this.selectedItemsSubject.next(this.result.items.filter((i) => infoHashes.has(i.infoHash)));
    }));
  }
  ngOnDestroy() {
    this.subscriptions.forEach((subscription) => subscription.unsubscribe());
    this.subscriptions = new Array();
  }
  static {
    this.\u0275fac = function TorrentsSearchComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _TorrentsSearchComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _TorrentsSearchComponent, selectors: [["app-torrents-search"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [["drawer", ""], [4, "transloco"], [3, "parts"], [1, "drawer-container"], [1, "drawer", 3, "mode", "opened"], [1, "panel-content-type", 3, "expanded"], [3, "click"], ["fontSet", "material-icons"], [1, "search-form"], [1, "form-field-container", "button-container", "button-container-toggle-drawer"], ["type", "button", "mat-icon-button", "", 1, "button-toggle-drawer", 3, "click", "matTooltip"], ["aria-label", "Side nav toggle icon", "fontSet", "material-icons"], [1, "form-field-container", "form-field-container-search-query"], ["subscriptSizing", "dynamic"], ["matInput", "", "autocapitalize", "none", 3, "keyup.enter", "placeholder", "formControl"], ["mat-icon-button", "", 3, "matTooltip"], [1, "form-field-container", "form-field-container-order-by"], [3, "valueChange", "value"], ["mat-icon-button", "", 3, "click", "matTooltip"], [1, "form-field-container", "button-container", "button-container-refresh"], ["mat-mini-fab", "", "color", "primary", 3, "click", "matTooltip"], [3, "updated", "selectedItems$"], [3, "updated", "dataSource", "controller", "displayedColumns", "multiSelection"], [3, "paging", "page", "pageSize", "pageLength", "totalLength", "totalIsEstimate", "hasNextPage"], [3, "class"], [3, "expanded"], [3, "opened", "closed", "expanded"], [1, "filtered"], [1, "unfiltered"], [3, "checked"], [3, "change", "checked"], ["checked", "true"], ["checked", "true", 3, "change"], [3, "value"]], template: function TorrentsSearchComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, TorrentsSearchComponent_ng_container_0_Template, 54, 42, "ng-container", 1);
      }
    }, dependencies: [
      AppModule,
      MatOption,
      MatIconButton,
      MatMiniFabButton,
      MatCheckbox,
      MatDivider,
      MatExpansionPanel,
      MatExpansionPanelHeader,
      MatExpansionPanelTitle,
      MatFormField,
      MatLabel,
      MatIcon,
      MatInput,
      MatSelect,
      MatDrawer,
      MatDrawerContainer,
      MatDrawerContent,
      MatTooltip,
      DefaultValueAccessor,
      NgControlStatus,
      FormControlDirective,
      TranslocoDirective,
      AsyncPipe,
      DocumentTitleComponent,
      GraphQLModule,
      PaginatorComponent,
      TorrentsBulkActionsComponent,
      TorrentsTableComponent,
      IntEstimatePipe
    ], styles: ["\n\n.mat-expansion-panel[_ngcontent-%COMP%] {\n  margin-top: 14px;\n  margin-right: 14px;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   section[_ngcontent-%COMP%] {\n  margin-left: -10px;\n}\n.mat-expansion-panel.panel-content-type[_ngcontent-%COMP%] {\n  margin-top: 20px;\n}\n.mat-expansion-panel.panel-content-type[_ngcontent-%COMP%]   section[_ngcontent-%COMP%] {\n  margin-left: 0;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   ul[_ngcontent-%COMP%] {\n  list-style: none;\n  padding-left: 0;\n  margin: 0;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   mat-panel-title[_ngcontent-%COMP%], \n.mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%] {\n  position: relative;\n  line-height: 40px;\n  padding-left: 40px;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   mat-panel-title[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%], \n.mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  position: absolute;\n  left: 0;\n  top: 8px;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%] {\n  cursor: pointer;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  top: 6px;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%]   small[_ngcontent-%COMP%] {\n  float: right;\n  font-size: 0.8rem;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   mat-checkbox[_ngcontent-%COMP%] {\n  display: block;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   mat-checkbox[_ngcontent-%COMP%]     label {\n  min-width: 220px;\n}\n.mat-expansion-panel[_ngcontent-%COMP%]   mat-checkbox[_ngcontent-%COMP%]   small[_ngcontent-%COMP%] {\n  margin-left: 10px;\n  position: absolute;\n  right: 0;\n}\n.search-form[_ngcontent-%COMP%] {\n  padding-top: 20px;\n  padding-bottom: 10px;\n  position: relative;\n  clear: both;\n  display: flex;\n  flex-wrap: wrap;\n}\n.search-form[_ngcontent-%COMP%]   .form-field-container[_ngcontent-%COMP%] {\n  display: inline-flex;\n  flex-direction: column;\n  position: relative;\n  margin-left: 20px;\n  padding-bottom: 20px;\n}\n.search-form[_ngcontent-%COMP%]   .form-field-container[_ngcontent-%COMP%]   button[_ngcontent-%COMP%] {\n  top: 8px;\n}\n.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-order-by[_ngcontent-%COMP%] {\n  padding-right: 40px;\n}\n.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-order-by[_ngcontent-%COMP%]   button[_ngcontent-%COMP%] {\n  position: absolute;\n  right: 0;\n}\n.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-search-query[_ngcontent-%COMP%] {\n  width: 300px;\n}\n.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-search-query[_ngcontent-%COMP%]   button[_ngcontent-%COMP%] {\n  position: absolute;\n  right: 0;\n}\n.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-search-query[_ngcontent-%COMP%]     .mat-mdc-form-field-infix {\n  padding-right: 50px;\n}\n.search-form[_ngcontent-%COMP%]   .button-container-toggle-direction[_ngcontent-%COMP%] {\n  margin-left: 4px;\n}\napp-paginator[_ngcontent-%COMP%] {\n  float: right;\n  padding-top: 14px;\n  padding-bottom: 20px;\n}\n/*# sourceMappingURL=torrents-search.component.css.map */"], changeDetection: 0 });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(TorrentsSearchComponent, { className: "TorrentsSearchComponent", filePath: "src/app/torrents/torrents-search.component.ts", lineNumber: 73 });
})();
var defaultLimit = 20;
var initControls = {
  page: 1,
  limit: defaultLimit,
  contentType: null,
  orderBy: defaultOrderBy,
  facets: {
    genre: inactiveFacet,
    language: inactiveFacet,
    fileType: inactiveFacet,
    torrentSource: inactiveFacet,
    torrentTag: inactiveFacet,
    videoResolution: inactiveFacet,
    videoSource: inactiveFacet
  }
};
var paramsToControls = (params) => {
  const queryString = stringParam(params, "query");
  const activeFacets = stringListParam(params, "facets");
  let selectedTorrent;
  const selectedTorrentParam = stringParam(params, "torrent");
  if (selectedTorrentParam) {
    let torrentTabSelection;
    const strTab = stringParam(params, "tab");
    if (torrentTabNames.includes(strTab)) {
      torrentTabSelection = strTab;
    }
    selectedTorrent = {
      infoHash: selectedTorrentParam,
      tab: torrentTabSelection
    };
  }
  return {
    queryString,
    orderBy: orderByParam(params, !!queryString),
    contentType: contentTypeParam(params),
    limit: intParam(params, "limit") ?? defaultLimit,
    page: intParam(params, "page") ?? 1,
    selectedTorrent,
    facets: facets.reduce((acc, facet) => {
      const active = activeFacets?.includes(facet.key) ?? false;
      const filter = stringListParam(params, facet.key);
      return facet.patchInput(acc, {
        active,
        filter
      });
    }, initControls.facets)
  };
};
var controlsToParams = (ctrl) => {
  let page = ctrl.page;
  let limit = ctrl.limit;
  if (page === 1) {
    page = void 0;
  }
  if (limit === defaultLimit) {
    limit = void 0;
  }
  const orderBy = isDefaultOrdering(ctrl) ? void 0 : ctrl.orderBy;
  let desc;
  if (orderBy) {
    desc = orderBy.descending ? "1" : "0";
  }
  return __spreadValues(__spreadValues({
    query: ctrl.queryString ? encodeURIComponent(ctrl.queryString) : void 0,
    page,
    limit,
    content_type: ctrl.contentType,
    order: orderBy?.field,
    desc
  }, ctrl.selectedTorrent ? {
    torrent: ctrl.selectedTorrent.infoHash,
    tab: ctrl.selectedTorrent.tab ?? void 0
  } : {}), flattenFacets(ctrl.facets));
};
var contentTypeParam = (params) => {
  const str = stringParam(params, "content_type");
  return str && str in contentTypeMap ? str : null;
};
var orderByParam = (params, hasQuery) => {
  let desc = null;
  const strDesc = stringParam(params, "desc");
  if (strDesc === "1") {
    desc = true;
  } else if (strDesc === "0") {
    desc = false;
  }
  const field = stringParam(params, "order");
  for (const opt of orderByOptions) {
    if (opt.field === field) {
      return {
        field,
        descending: desc ?? opt.descending
      };
    }
  }
  return {
    field: hasQuery ? "relevance" : "published_at",
    descending: desc ?? true
  };
};
var flattenFacets = (ctrl) => {
  const [activeFacets, filters] = facets.reduce((acc, f) => {
    const input = f.extractInput(ctrl);
    if (input.active) {
      return [
        [...acc[0], f.key],
        input.filter ? __spreadProps(__spreadValues({}, acc[1]), {
          [f.key]: input.filter
        }) : acc[1]
      ];
    } else {
      return acc;
    }
  }, [[], {}]);
  return __spreadValues({
    facets: activeFacets.length ? activeFacets.join(",") : void 0
  }, Object.fromEntries(Object.entries(filters).map(([k, values]) => [
    k,
    encodeURIComponent(values.join(","))
  ])));
};
export {
  TorrentsSearchComponent
};
//# sourceMappingURL=chunk-EPENVDKZ.js.map
