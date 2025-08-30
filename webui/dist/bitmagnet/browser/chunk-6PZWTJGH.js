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
  MatButton,
  MatCard,
  MatCardContent,
  MatCardHeader,
  MatDialog,
  MatDialogActions,
  MatDialogContent,
  MatDialogModule,
  MatDialogRef,
  MatDivider,
  MatGridList,
  MatGridTile,
  MatIcon,
  MatIconButton,
  MatMenu,
  MatMenuItem,
  MatToolbar,
  MatTooltip,
  TranslocoDirective,
  WorkersDocument,
  WorkersRestartDocument,
  WorkersShutdownDocument,
  WorkersStartDocument
} from "./chunk-WWRDQTKJ.js";
import "./chunk-VSVMRYN2.js";
import "./chunk-Y2ZC5Z2X.js";
import {
  BehaviorSubject,
  __spreadProps,
  __spreadValues,
  inject,
  map,
  ɵsetClassDebugInfo,
  ɵɵStandaloneFeature,
  ɵɵadvance,
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
  ɵɵpureFunction0,
  ɵɵpureFunction1,
  ɵɵrepeater,
  ɵɵrepeaterCreate,
  ɵɵrepeaterTrackByIdentity,
  ɵɵresetView,
  ɵɵrestoreView,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate,
  ɵɵtextInterpolate1,
  ɵɵtextInterpolate2
} from "./chunk-DMMUMX3A.js";

// src/app/workers/workers.service.ts
var workerStateIcons = {
  error: "error",
  idle: "circle",
  running: "play_circle",
  shutdown: "pending",
  startup: "pending"
};
var initialResult = {
  workers: [],
  workerError: false,
  error: null
};
var pollInterval = 1e4;
var WorkersService = class {
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
      query: WorkersDocument,
      fetchPolicy: "no-cache",
      pollInterval
    }).valueChanges.pipe(map((r) => ({
      workers: r.data.worker.listAll.workers.map((w) => __spreadProps(__spreadValues({}, w), {
        icon: workerStateIcons[w.state]
      })),
      workerError: r.data.worker.listAll.workers.some((w) => w.error),
      error: null
    }))).subscribe({
      next: (result) => {
        this.resultSubject.next(result);
        setTimeout(this.watchQuery.bind(this), pollInterval);
      },
      error: (error) => {
        this.resultSubject.next({
          workers: [],
          workerError: false,
          error
        });
        setTimeout(this.watchQuery.bind(this), pollInterval);
      }
    });
  }
  startWorkers(...keys) {
    this.apollo.mutate({
      mutation: WorkersStartDocument,
      variables: {
        keys
      }
    }).pipe(map((result) => this.resultSubject.next({
      workers: result.data?.worker.start.workers?.map((w) => __spreadProps(__spreadValues({}, w), {
        icon: workerStateIcons[w.state]
      })) ?? [],
      workerError: result.data?.worker.start.workers?.some((w) => w.error) ?? false,
      error: null
    }))).subscribe();
  }
  shutdownWorkers(...keys) {
    this.apollo.mutate({
      mutation: WorkersShutdownDocument,
      variables: {
        keys
      }
    }).pipe(map((result) => this.resultSubject.next({
      workers: result.data?.worker.shutdown.workers?.map((w) => __spreadProps(__spreadValues({}, w), {
        icon: workerStateIcons[w.state]
      })) ?? [],
      workerError: result.data?.worker.shutdown.workers?.some((w) => w.error) ?? false,
      error: null
    }))).subscribe();
  }
  restartWorkers(...keys) {
    this.apollo.mutate({
      mutation: WorkersRestartDocument,
      variables: {
        keys
      }
    }).pipe(map((result) => {
      this.resultSubject.next({
        workers: result.data?.worker.restart.workers?.map((w) => __spreadProps(__spreadValues({}, w), {
          icon: workerStateIcons[w.state]
        })) ?? [],
        workerError: result.data?.worker.restart.workers?.some((w) => w.error) ?? false,
        error: null
      });
      setTimeout(this.watchQuery.bind(this), 2e3);
    })).subscribe();
  }
};

// src/app/workers/workers-confirm-action-dialog.component.ts
var _c0 = () => ["shutdown", "restart"];
function WorkersConfirmActionDialogComponent_ng_container_0_Conditional_4_For_4_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "li");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const name_r2 = ctx.$implicit;
    const t_r3 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("workers.names." + name_r2));
  }
}
function WorkersConfirmActionDialogComponent_ng_container_0_Conditional_4_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "p");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(2, "ul");
    \u0275\u0275repeaterCreate(3, WorkersConfirmActionDialogComponent_ng_container_0_Conditional_4_For_4_Template, 2, 1, "li", null, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    const ctx_r3 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("workers.warnings.required_by"));
    \u0275\u0275advance(2);
    \u0275\u0275repeater(ctx_r3.requiredBy);
  }
}
function WorkersConfirmActionDialogComponent_ng_container_0_Conditional_5_For_4_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "li");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const name_r5 = ctx.$implicit;
    const t_r3 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("workers.names." + name_r5));
  }
}
function WorkersConfirmActionDialogComponent_ng_container_0_Conditional_5_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "p");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(2, "ul");
    \u0275\u0275repeaterCreate(3, WorkersConfirmActionDialogComponent_ng_container_0_Conditional_5_For_4_Template, 2, 1, "li", null, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    const ctx_r3 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("workers.warnings.depends_on"));
    \u0275\u0275advance(2);
    \u0275\u0275repeater(ctx_r3.dependsOn);
  }
}
function WorkersConfirmActionDialogComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-dialog-content")(2, "h3");
    \u0275\u0275text(3);
    \u0275\u0275elementEnd();
    \u0275\u0275template(4, WorkersConfirmActionDialogComponent_ng_container_0_Conditional_4_Template, 5, 1)(5, WorkersConfirmActionDialogComponent_ng_container_0_Conditional_5_Template, 5, 1);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(6, "mat-dialog-actions")(7, "button", 1);
    \u0275\u0275listener("click", function WorkersConfirmActionDialogComponent_ng_container_0_Template_button_click_7_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r3 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r3.cancel());
    });
    \u0275\u0275text(8);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(9, "button", 2);
    \u0275\u0275listener("click", function WorkersConfirmActionDialogComponent_ng_container_0_Template_button_click_9_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r3 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r3.confirm());
    });
    \u0275\u0275text(10);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r3 = ctx.$implicit;
    const ctx_r3 = \u0275\u0275nextContext();
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate2("", t_r3("workers.actions." + ctx_r3.action), " ", t_r3("workers.names." + ctx_r3.worker), "?");
    \u0275\u0275advance();
    \u0275\u0275conditional(\u0275\u0275pureFunction0(6, _c0).includes(ctx_r3.action) && ctx_r3.requiredBy.length > 0 ? 4 : -1);
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r3.action == "start" && ctx_r3.dependsOn.length > 0 ? 5 : -1);
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate(t_r3("general.cancel"));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(t_r3("general.confirm"));
  }
}
var WorkersConfirmActionDialogComponent = class _WorkersConfirmActionDialogComponent {
  constructor() {
    this.workers = inject(WorkersService);
    this.dialogRef = inject(MatDialogRef);
    this.dependsOn = Array();
    this.requiredBy = Array();
  }
  ngOnInit() {
    this.workers.result$.subscribe((result) => {
      const worker = result.workers.find((w) => w.key === this.worker);
      if (worker) {
        this.dependsOn = worker.dependsOn;
        this.requiredBy = worker.requiredBy;
      }
    });
  }
  confirm() {
    switch (this.action) {
      case "start":
        this.workers.startWorkers(this.worker);
        break;
      case "shutdown":
        this.workers.shutdownWorkers(this.worker);
        break;
      case "restart":
        this.workers.restartWorkers(this.worker);
        break;
    }
    this.dialogRef.close();
  }
  cancel() {
    this.dialogRef.close();
  }
  static {
    this.\u0275fac = function WorkersConfirmActionDialogComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _WorkersConfirmActionDialogComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _WorkersConfirmActionDialogComponent, selectors: [["app-workers-confirm-action-dialog"]], inputs: { action: "action", worker: "worker" }, decls: 1, vars: 0, consts: [[4, "transloco"], ["mat-button", "", "color", "secondary", 3, "click"], ["mat-button", "", "color", "warning", 3, "click"]], template: function WorkersConfirmActionDialogComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, WorkersConfirmActionDialogComponent_ng_container_0_Template, 11, 7, "ng-container", 0);
      }
    }, dependencies: [TranslocoDirective, MatDialogActions, MatDialogContent, MatButton], encapsulation: 2 });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(WorkersConfirmActionDialogComponent, { className: "WorkersConfirmActionDialogComponent", filePath: "src/app/workers/workers-confirm-action-dialog.component.ts", lineNumber: 11 });
})();

// src/app/workers/workers-table.component.ts
var _forTrack0 = ($index, $item) => $item.key;
function WorkersTableComponent_ng_container_0_Conditional_9_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 2);
    \u0275\u0275text(1, "Actions");
    \u0275\u0275elementEnd();
  }
}
function WorkersTableComponent_ng_container_0_Conditional_10_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r1 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r1("general.error"));
  }
}
function WorkersTableComponent_ng_container_0_For_13_Conditional_8_Conditional_1_Template(rf, ctx) {
  if (rf & 1) {
    const _r3 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "button", 6);
    \u0275\u0275listener("click", function WorkersTableComponent_ng_container_0_For_13_Conditional_8_Conditional_1_Template_button_click_0_listener() {
      \u0275\u0275restoreView(_r3);
      const worker_r4 = \u0275\u0275nextContext(2).$implicit;
      const ctx_r4 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r4.confirm("shutdown", worker_r4.key));
    });
    \u0275\u0275elementStart(1, "mat-icon");
    \u0275\u0275text(2, "stop");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const worker_r4 = \u0275\u0275nextContext(2).$implicit;
    const t_r1 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("disabled", worker_r4.key == "http_server" || worker_r4.requiredBy.includes("http_server"))("matTooltip", t_r1("workers.actions.shutdown"));
  }
}
function WorkersTableComponent_ng_container_0_For_13_Conditional_8_Conditional_2_Template(rf, ctx) {
  if (rf & 1) {
    const _r6 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "button", 6);
    \u0275\u0275listener("click", function WorkersTableComponent_ng_container_0_For_13_Conditional_8_Conditional_2_Template_button_click_0_listener() {
      \u0275\u0275restoreView(_r6);
      const worker_r4 = \u0275\u0275nextContext(2).$implicit;
      const ctx_r4 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r4.confirm("start", worker_r4.key));
    });
    \u0275\u0275elementStart(1, "mat-icon");
    \u0275\u0275text(2, "play_arrow");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const worker_r4 = \u0275\u0275nextContext(2).$implicit;
    const t_r1 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("disabled", worker_r4.state != "idle")("matTooltip", t_r1("workers.actions.start"));
  }
}
function WorkersTableComponent_ng_container_0_For_13_Conditional_8_Template(rf, ctx) {
  if (rf & 1) {
    const _r2 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td");
    \u0275\u0275template(1, WorkersTableComponent_ng_container_0_For_13_Conditional_8_Conditional_1_Template, 3, 2, "button", 5)(2, WorkersTableComponent_ng_container_0_For_13_Conditional_8_Conditional_2_Template, 3, 2, "button", 5);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(3, "td")(4, "button", 6);
    \u0275\u0275listener("click", function WorkersTableComponent_ng_container_0_For_13_Conditional_8_Template_button_click_4_listener() {
      \u0275\u0275restoreView(_r2);
      const worker_r4 = \u0275\u0275nextContext().$implicit;
      const ctx_r4 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r4.confirm("restart", worker_r4.key));
    });
    \u0275\u0275elementStart(5, "mat-icon");
    \u0275\u0275text(6, "replay");
    \u0275\u0275elementEnd()()();
  }
  if (rf & 2) {
    const worker_r4 = \u0275\u0275nextContext().$implicit;
    const t_r1 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275conditional(worker_r4.state === "running" ? 1 : 2);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", worker_r4.state != "running")("matTooltip", t_r1("workers.actions.restart"));
  }
}
function WorkersTableComponent_ng_container_0_For_13_Conditional_9_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "td");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const worker_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(worker_r4.error);
  }
}
function WorkersTableComponent_ng_container_0_For_13_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "tr")(1, "td", 3)(2, "mat-icon");
    \u0275\u0275text(3);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(4, "th", 4);
    \u0275\u0275text(5);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(6, "td");
    \u0275\u0275text(7);
    \u0275\u0275elementEnd();
    \u0275\u0275template(8, WorkersTableComponent_ng_container_0_For_13_Conditional_8_Template, 7, 3)(9, WorkersTableComponent_ng_container_0_For_13_Conditional_9_Template, 2, 1, "td");
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const worker_r4 = ctx.$implicit;
    const t_r1 = \u0275\u0275nextContext().$implicit;
    const ctx_r4 = \u0275\u0275nextContext();
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate(worker_r4.icon);
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate1(" ", t_r1("workers.names." + worker_r4.key), " ");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate1(" ", t_r1("workers.states." + worker_r4.state), " ");
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r4.actions ? 8 : -1);
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r4.workers.result.workerError ? 9 : -1);
  }
}
function WorkersTableComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "table", 1)(2, "thead")(3, "tr");
    \u0275\u0275element(4, "th");
    \u0275\u0275elementStart(5, "th");
    \u0275\u0275text(6);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(7, "th");
    \u0275\u0275text(8);
    \u0275\u0275elementEnd();
    \u0275\u0275template(9, WorkersTableComponent_ng_container_0_Conditional_9_Template, 2, 0, "th", 2)(10, WorkersTableComponent_ng_container_0_Conditional_10_Template, 2, 1, "th");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(11, "tbody");
    \u0275\u0275repeaterCreate(12, WorkersTableComponent_ng_container_0_For_13_Template, 10, 5, "tr", null, _forTrack0);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r1 = ctx.$implicit;
    const ctx_r4 = \u0275\u0275nextContext();
    \u0275\u0275advance(6);
    \u0275\u0275textInterpolate(t_r1("workers.worker"));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(t_r1("general.status"));
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r4.actions ? 9 : -1);
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r4.workers.result.workerError ? 10 : -1);
    \u0275\u0275advance(2);
    \u0275\u0275repeater(ctx_r4.workers.result.workers);
  }
}
var WorkersTableComponent = class _WorkersTableComponent {
  constructor() {
    this.workers = inject(WorkersService);
    this.dialog = inject(MatDialog);
    this.actions = false;
  }
  confirm(action, worker) {
    const ref = this.dialog.open(WorkersConfirmActionDialogComponent);
    ref.componentInstance.action = action;
    ref.componentInstance.worker = worker;
  }
  static {
    this.\u0275fac = function WorkersTableComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _WorkersTableComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _WorkersTableComponent, selectors: [["app-workers-table"]], inputs: { actions: "actions" }, decls: 1, vars: 0, consts: [[4, "transloco"], [1, "table-workers"], ["colspan", "2"], [1, "td-icon"], ["scope", "row"], ["mat-icon-button", "", 3, "disabled", "matTooltip"], ["mat-icon-button", "", 3, "click", "disabled", "matTooltip"]], template: function WorkersTableComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, WorkersTableComponent_ng_container_0_Template, 14, 4, "ng-container", 0);
      }
    }, dependencies: [TranslocoDirective, MatIcon, MatTooltip, MatIconButton], styles: ["\n\n.table-workers[_ngcontent-%COMP%] {\n  width: 100%;\n}\n.table-workers[_ngcontent-%COMP%]   th[_ngcontent-%COMP%], \n.table-workers[_ngcontent-%COMP%]   td[_ngcontent-%COMP%] {\n  padding-right: 20px;\n}\n.table-workers[_ngcontent-%COMP%]   th[_ngcontent-%COMP%] {\n  text-align: left;\n}\n.table-workers[_ngcontent-%COMP%]   thead[_ngcontent-%COMP%]   th[_ngcontent-%COMP%] {\n  padding-bottom: 10px;\n}\n/*# sourceMappingURL=workers-table.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(WorkersTableComponent, { className: "WorkersTableComponent", filePath: "src/app/workers/workers-table.component.ts", lineNumber: 13 });
})();

// src/app/workers/workers-card.component.ts
function WorkersCardComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-card", 1)(2, "mat-card-content", 2)(3, "mat-card")(4, "mat-card-content");
    \u0275\u0275element(5, "app-workers-table", 3);
    \u0275\u0275elementEnd()()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    \u0275\u0275advance(5);
    \u0275\u0275property("actions", true);
  }
}
var WorkersCardComponent = class _WorkersCardComponent {
  constructor() {
    this.workers = inject(WorkersService);
  }
  static {
    this.\u0275fac = function WorkersCardComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _WorkersCardComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _WorkersCardComponent, selectors: [["app-workers-card"]], decls: 1, vars: 0, consts: [[4, "transloco"], [1, "dashboard-card", "dashboard-card-health"], [1, "dashboard-card-content"], [3, "actions"]], template: function WorkersCardComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, WorkersCardComponent_ng_container_0_Template, 6, 1, "ng-container", 0);
      }
    }, dependencies: [TranslocoDirective, MatCard, MatCardContent, WorkersTableComponent] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(WorkersCardComponent, { className: "WorkersCardComponent", filePath: "src/app/workers/workers-card.component.ts", lineNumber: 10 });
})();

// src/app/workers/workers.module.ts
var WorkersModule = class _WorkersModule {
  static {
    this.\u0275fac = function WorkersModule_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _WorkersModule)();
    };
  }
  static {
    this.\u0275mod = /* @__PURE__ */ \u0275\u0275defineNgModule({ type: _WorkersModule });
  }
  static {
    this.\u0275inj = /* @__PURE__ */ \u0275\u0275defineInjector({ providers: [WorkersService], imports: [
      GraphQLModule,
      MatIcon,
      MatDialogModule,
      MatButton,
      MatIconButton,
      MatCard,
      MatCardHeader,
      MatGridTile,
      MatMenu,
      MatMenuItem
    ] });
  }
};

// src/app/dashboard/workers/dashboard-workers.component.ts
var _c02 = (a0) => [a0];
function DashboardWorkersComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275element(1, "app-document-title", 1);
    \u0275\u0275elementStart(2, "mat-card", 2)(3, "mat-card-header")(4, "mat-toolbar")(5, "h2")(6, "mat-icon");
    \u0275\u0275text(7, "manufacturing");
    \u0275\u0275elementEnd();
    \u0275\u0275text(8);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(9, "mat-card-content");
    \u0275\u0275element(10, "mat-divider");
    \u0275\u0275elementStart(11, "div", 3)(12, "mat-grid-list", 4)(13, "mat-grid-tile", 5);
    \u0275\u0275element(14, "app-workers-card");
    \u0275\u0275elementEnd()()()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r1 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275property("parts", \u0275\u0275pureFunction1(5, _c02, t_r1("routes.dashboard")));
    \u0275\u0275advance(7);
    \u0275\u0275textInterpolate(t_r1("routes.workers"));
    \u0275\u0275advance(4);
    \u0275\u0275property("cols", ctx_r1.breakpoints.sizeAtLeast("Medium") ? 2 : 1);
    \u0275\u0275advance();
    \u0275\u0275property("colspan", 1)("rowspan", 1);
  }
}
var DashboardWorkersComponent = class _DashboardWorkersComponent {
  constructor() {
    this.breakpoints = inject(BreakpointsService);
  }
  static {
    this.\u0275fac = function DashboardWorkersComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _DashboardWorkersComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _DashboardWorkersComponent, selectors: [["app-workers-dashboard"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], [3, "parts"], [1, "dashboard-card"], [1, "grid-container"], ["rowHeight", "600px", 3, "cols"], [3, "colspan", "rowspan"]], template: function DashboardWorkersComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, DashboardWorkersComponent_ng_container_0_Template, 15, 7, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatCard, MatCardContent, MatCardHeader, MatDivider, MatGridList, MatGridTile, MatIcon, MatToolbar, TranslocoDirective, DocumentTitleComponent, WorkersModule, WorkersCardComponent], styles: ["\n\n.grid-container[_ngcontent-%COMP%] {\n  margin: 20px;\n}\n.more-button[_ngcontent-%COMP%] {\n  position: absolute;\n  top: 5px;\n  right: 10px;\n}\napp-health-card[_ngcontent-%COMP%] {\n  width: 100%;\n  height: 100%;\n}\napp-health-card[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%] {\n  height: 100%;\n}\nmat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%] {\n  width: 100%;\n}\nmat-toolbar[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  position: relative;\n  top: 3px;\n  margin-right: 14px;\n  margin-left: 32px;\n}\n/*# sourceMappingURL=dashboard-workers.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(DashboardWorkersComponent, { className: "DashboardWorkersComponent", filePath: "src/app/dashboard/workers/dashboard-workers.component.ts", lineNumber: 14 });
})();
export {
  DashboardWorkersComponent
};
//# sourceMappingURL=chunk-6PZWTJGH.js.map
