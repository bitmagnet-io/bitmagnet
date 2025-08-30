import {
  Apollo,
  GraphQLModule,
  HealthCheckDocument,
  MatButton,
  MatCard,
  MatCardContent,
  MatCardHeader,
  MatCardTitle,
  MatDialog,
  MatDialogActions,
  MatDialogClose,
  MatDialogContent,
  MatDialogModule,
  MatDialogTitle,
  MatDivider,
  MatGridTile,
  MatIcon,
  MatIconButton,
  MatMenu,
  MatMenuItem,
  MatTooltip,
  TranslocoDirective
} from "./chunk-WWRDQTKJ.js";
import {
  BehaviorSubject,
  __spreadProps,
  __spreadValues,
  inject,
  map,
  ɵsetClassDebugInfo,
  ɵɵadvance,
  ɵɵclassMap,
  ɵɵconditional,
  ɵɵdefineComponent,
  ɵɵdefineInjector,
  ɵɵdefineNgModule,
  ɵɵelement,
  ɵɵelementContainerEnd,
  ɵɵelementContainerStart,
  ɵɵelementEnd,
  ɵɵelementStart,
  ɵɵgetCurrentView,
  ɵɵlistener,
  ɵɵnextContext,
  ɵɵproperty,
  ɵɵpureFunction1,
  ɵɵreference,
  ɵɵrepeater,
  ɵɵrepeaterCreate,
  ɵɵresetView,
  ɵɵrestoreView,
  ɵɵtemplate,
  ɵɵtemplateRefExtractor,
  ɵɵtext,
  ɵɵtextInterpolate,
  ɵɵtextInterpolate1,
  ɵɵtextInterpolate2
} from "./chunk-DMMUMX3A.js";

// src/app/health/health.service.ts
var statusIcons = {
  error: "error",
  degraded: "warning",
  down: "warning",
  unknown: "pending",
  inactive: "circle",
  up: "check_circle",
  started: "play_circle"
};
var initialResult = {
  status: "unknown",
  checks: [],
  icon: statusIcons.unknown,
  error: null
};
var pollInterval = 1e4;
var HealthService = class {
  constructor() {
    this.apollo = inject(Apollo);
    this.resultSubject = new BehaviorSubject(initialResult);
    this.result$ = this.resultSubject.asObservable();
    this.result = initialResult;
    this.watchQuery();
    this.result$.subscribe((result) => {
      this.result = result;
    });
  }
  watchQuery() {
    this.apollo.watchQuery({
      query: HealthCheckDocument,
      fetchPolicy: "no-cache",
      pollInterval
    }).valueChanges.pipe(map((r) => ({
      status: r.data.health.status === "down" ? "degraded" : r.data.health.status,
      checks: r.data.health.checks.map((c) => __spreadProps(__spreadValues({}, c), {
        icon: statusIcons[c.status]
      })),
      icon: statusIcons[r.data.health.status],
      error: null
    }))).subscribe({
      next: (result) => this.resultSubject.next(result),
      error: (error) => {
        this.resultSubject.next({
          status: "error",
          checks: [],
          error,
          icon: statusIcons.error
        });
        setTimeout(this.watchQuery.bind(this), pollInterval);
      }
    });
  }
};

// src/app/health/health-summary.component.ts
var _forTrack0 = ($index, $item) => $item.key;
function HealthSummaryComponent_ng_container_0_Conditional_1_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "p");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r1 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275textInterpolate2("", t_r1("health.check_failed_with_error"), ": ", ctx_r1.health.result.error, "");
  }
}
function HealthSummaryComponent_ng_container_0_Conditional_2_Conditional_8_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r1 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r1("general.error"));
  }
}
function HealthSummaryComponent_ng_container_0_Conditional_2_For_11_Conditional_8_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "td");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const check_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(check_r3.error);
  }
}
function HealthSummaryComponent_ng_container_0_Conditional_2_For_11_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "tr")(1, "td", 2)(2, "mat-icon");
    \u0275\u0275text(3);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(4, "th", 3);
    \u0275\u0275text(5);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(6, "td");
    \u0275\u0275text(7);
    \u0275\u0275elementEnd();
    \u0275\u0275template(8, HealthSummaryComponent_ng_container_0_Conditional_2_For_11_Conditional_8_Template, 2, 1, "td");
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const check_r3 = ctx.$implicit;
    const t_r1 = \u0275\u0275nextContext(2).$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate(check_r3.icon);
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(t_r1("health.components." + check_r3.key));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(t_r1("health.statuses." + check_r3.status));
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r1.health.result.status === "down" ? 8 : -1);
  }
}
function HealthSummaryComponent_ng_container_0_Conditional_2_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "table", 1)(1, "thead")(2, "tr");
    \u0275\u0275element(3, "th");
    \u0275\u0275elementStart(4, "th");
    \u0275\u0275text(5);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(6, "th");
    \u0275\u0275text(7);
    \u0275\u0275elementEnd();
    \u0275\u0275template(8, HealthSummaryComponent_ng_container_0_Conditional_2_Conditional_8_Template, 2, 1, "th");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(9, "tbody");
    \u0275\u0275repeaterCreate(10, HealthSummaryComponent_ng_container_0_Conditional_2_For_11_Template, 9, 4, "tr", null, _forTrack0);
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const t_r1 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance(5);
    \u0275\u0275textInterpolate(t_r1("health.component"));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(t_r1("general.status"));
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r1.health.result.status === "down" ? 8 : -1);
    \u0275\u0275advance(2);
    \u0275\u0275repeater(ctx_r1.health.result.checks);
  }
}
function HealthSummaryComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275template(1, HealthSummaryComponent_ng_container_0_Conditional_1_Template, 2, 2, "p")(2, HealthSummaryComponent_ng_container_0_Conditional_2_Template, 12, 3, "table", 1);
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r1.health.result.error ? 1 : 2);
  }
}
var HealthSummaryComponent = class _HealthSummaryComponent {
  constructor() {
    this.health = inject(HealthService);
  }
  static {
    this.\u0275fac = function HealthSummaryComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _HealthSummaryComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _HealthSummaryComponent, selectors: [["app-health-summary"]], decls: 1, vars: 0, consts: [[4, "transloco"], [1, "table-health"], [1, "td-icon"], ["scope", "row"]], template: function HealthSummaryComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, HealthSummaryComponent_ng_container_0_Template, 3, 1, "ng-container", 0);
      }
    }, dependencies: [TranslocoDirective, MatIcon], styles: ["\n\n.table-health[_ngcontent-%COMP%] {\n  width: 100%;\n}\n.table-health[_ngcontent-%COMP%]   th[_ngcontent-%COMP%], \n.table-health[_ngcontent-%COMP%]   td[_ngcontent-%COMP%] {\n  padding-right: 20px;\n}\n.table-health[_ngcontent-%COMP%]   th[_ngcontent-%COMP%] {\n  text-align: left;\n}\n.table-health[_ngcontent-%COMP%]   thead[_ngcontent-%COMP%]   th[_ngcontent-%COMP%] {\n  padding-bottom: 10px;\n}\n/*# sourceMappingURL=health-summary.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(HealthSummaryComponent, { className: "HealthSummaryComponent", filePath: "src/app/health/health-summary.component.ts", lineNumber: 10 });
})();

// src/app/health/health-card.component.ts
var _c0 = (a0) => ({ status: a0 });
function HealthCardComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-card", 1)(2, "mat-card-header")(3, "mat-card-title")(4, "h3")(5, "mat-icon");
    \u0275\u0275text(6);
    \u0275\u0275elementEnd();
    \u0275\u0275text(7);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(8, "mat-card-content", 2)(9, "mat-card")(10, "mat-card-content");
    \u0275\u0275element(11, "app-health-summary");
    \u0275\u0275elementEnd()()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r1 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance(6);
    \u0275\u0275textInterpolate(ctx_r1.health.result.icon);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r1("health.bitmagnet_is_status", \u0275\u0275pureFunction1(2, _c0, t_r1("health.statuses." + (ctx_r1.health.result.error ? "down" : ctx_r1.health.result.status)))), " ");
  }
}
var HealthCardComponent = class _HealthCardComponent {
  constructor() {
    this.health = inject(HealthService);
  }
  static {
    this.\u0275fac = function HealthCardComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _HealthCardComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _HealthCardComponent, selectors: [["app-health-card"]], decls: 1, vars: 0, consts: [[4, "transloco"], [1, "dashboard-card", "dashboard-card-health"], [1, "dashboard-card-content"]], template: function HealthCardComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, HealthCardComponent_ng_container_0_Template, 12, 4, "ng-container", 0);
      }
    }, dependencies: [TranslocoDirective, MatIcon, MatCard, MatCardContent, MatCardHeader, MatCardTitle, HealthSummaryComponent], styles: ["\n\n.dashboard-card-health[_ngcontent-%COMP%] {\n  position: absolute;\n  top: 15px;\n  left: 15px;\n  right: 15px;\n  bottom: 15px;\n}\n.dashboard-card-health[_ngcontent-%COMP%]   h3[_ngcontent-%COMP%] {\n  margin-top: 0;\n}\n.dashboard-card-health[_ngcontent-%COMP%]   h3[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  position: relative;\n  top: 3px;\n  margin-left: 4px;\n  margin-right: 6px;\n}\n.dashboard-card-health[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]   mat-card-header[_ngcontent-%COMP%] {\n  margin-bottom: 16px;\n}\n/*# sourceMappingURL=health-card.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(HealthCardComponent, { className: "HealthCardComponent", filePath: "src/app/health/health-card.component.ts", lineNumber: 10 });
})();

// src/app/health/health-widget.component.ts
var _c02 = (a0) => ({ status: a0 });
function HealthWidgetComponent_ng_container_0_ng_template_4_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "h2", 3);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(2, "mat-dialog-content")(3, "mat-card")(4, "mat-card-content");
    \u0275\u0275element(5, "app-health-summary");
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(6, "mat-dialog-actions")(7, "button", 4);
    \u0275\u0275text(8);
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const t_r4 = \u0275\u0275nextContext().$implicit;
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275textInterpolate2(" ", t_r4("health.summary"), ": ", t_r4("health.bitmagnet_is_status", \u0275\u0275pureFunction1(3, _c02, t_r4("health.statuses." + (ctx_r2.health.result.error ? "down" : ctx_r2.health.result.status)))), " ");
    \u0275\u0275advance(7);
    \u0275\u0275textInterpolate1(" ", t_r4("general.dismiss"), " ");
  }
}
function HealthWidgetComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "button", 2);
    \u0275\u0275listener("click", function HealthWidgetComponent_ng_container_0_Template_button_click_1_listener() {
      \u0275\u0275restoreView(_r1);
      const healthDialog_r2 = \u0275\u0275reference(5);
      const ctx_r2 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r2.dialog.open(healthDialog_r2));
    });
    \u0275\u0275elementStart(2, "mat-icon");
    \u0275\u0275text(3);
    \u0275\u0275elementEnd()();
    \u0275\u0275template(4, HealthWidgetComponent_ng_container_0_ng_template_4_Template, 9, 5, "ng-template", null, 0, \u0275\u0275templateRefExtractor);
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r4 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275classMap("health-icon health-icon-" + ctx_r2.health.result.status);
    \u0275\u0275property("matTooltip", t_r4("health.bitmagnet_is_status", \u0275\u0275pureFunction1(4, _c02, t_r4("health.statuses." + ctx_r2.health.result.status))));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(ctx_r2.health.result.icon);
  }
}
var HealthWidgetComponent = class _HealthWidgetComponent {
  constructor() {
    this.health = inject(HealthService);
    this.dialog = inject(MatDialog);
  }
  static {
    this.\u0275fac = function HealthWidgetComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _HealthWidgetComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _HealthWidgetComponent, selectors: [["app-health-widget"]], decls: 1, vars: 0, consts: [["healthDialog", ""], [4, "transloco"], ["mat-icon-button", "", 3, "click", "matTooltip"], ["matDialogTitle", ""], ["mat-button", "", "matDialogClose", "", "color", "primary"]], template: function HealthWidgetComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, HealthWidgetComponent_ng_container_0_Template, 6, 6, "ng-container", 1);
      }
    }, dependencies: [TranslocoDirective, MatIcon, MatTooltip, MatDialogClose, MatDialogTitle, MatDialogActions, MatDialogContent, MatButton, MatIconButton, MatCard, MatCardContent, HealthSummaryComponent], styles: ["\n\nmat-card-header[_ngcontent-%COMP%] {\n  margin-bottom: 16px;\n}\nmat-card[_ngcontent-%COMP%]    + mat-card[_ngcontent-%COMP%] {\n  margin-top: 16px;\n}\n/*# sourceMappingURL=health-widget.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(HealthWidgetComponent, { className: "HealthWidgetComponent", filePath: "src/app/health/health-widget.component.ts", lineNumber: 11 });
})();

// src/app/health/health.module.ts
var HealthModule = class _HealthModule {
  static {
    this.\u0275fac = function HealthModule_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _HealthModule)();
    };
  }
  static {
    this.\u0275mod = /* @__PURE__ */ \u0275\u0275defineNgModule({ type: _HealthModule });
  }
  static {
    this.\u0275inj = /* @__PURE__ */ \u0275\u0275defineInjector({ providers: [HealthService], imports: [
      GraphQLModule,
      MatIcon,
      MatDialogModule,
      MatButton,
      MatIconButton,
      MatCard,
      MatCardHeader,
      MatGridTile,
      MatMenu,
      MatMenuItem,
      MatDivider
    ] });
  }
};

export {
  HealthService,
  HealthCardComponent,
  HealthWidgetComponent,
  HealthModule
};
//# sourceMappingURL=chunk-WNZLSZUT.js.map
