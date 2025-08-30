import {
  formatTimeAgo
} from "./chunk-ORIQXXAG.js";
import {
  AppModule,
  MatFormField,
  MatIcon,
  MatIconButton,
  MatLabel,
  MatOption,
  MatSelect,
  MatTooltip,
  TranslocoDirective,
  TranslocoService
} from "./chunk-WWRDQTKJ.js";
import {
  DecimalPipe,
  EventEmitter,
  inject,
  numberAttribute,
  ɵsetClassDebugInfo,
  ɵɵInputTransformsFeature,
  ɵɵStandaloneFeature,
  ɵɵadvance,
  ɵɵconditional,
  ɵɵdefineComponent,
  ɵɵdefinePipe,
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
  ɵɵpureFunction2,
  ɵɵpureFunction3,
  ɵɵrepeater,
  ɵɵrepeaterCreate,
  ɵɵrepeaterTrackByIdentity,
  ɵɵresetView,
  ɵɵrestoreView,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate1
} from "./chunk-DMMUMX3A.js";

// src/app/pipes/int-estimate.pipe.ts
var IntEstimatePipe = class _IntEstimatePipe {
  constructor() {
    this.transloco = inject(TranslocoService);
  }
  transform(n, isEstimate = true, sigFigs = 2) {
    if (isEstimate && n > 0 && sigFigs > 0) {
      const magnitude = Math.floor(Math.log10(Math.abs(n)));
      const scale = Math.pow(10, magnitude - (sigFigs - 1));
      n = Math.round(n / scale) * scale;
    }
    const str = Intl.NumberFormat(this.transloco.getActiveLang()).format(n);
    return isEstimate ? `~${str}` : str;
  }
  static {
    this.\u0275fac = function IntEstimatePipe_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _IntEstimatePipe)();
    };
  }
  static {
    this.\u0275pipe = /* @__PURE__ */ \u0275\u0275definePipe({ name: "intEstimate", type: _IntEstimatePipe, pure: false, standalone: true });
  }
};

// src/app/paginator/paginator.component.ts
var _c0 = (a0, a1, a2) => ({ x: a0, y: a1, z: a2 });
var _c1 = (a0, a1) => ({ x: a0, y: a1 });
var _c2 = (a0) => [null, a0];
function PaginatorComponent_ng_container_0_For_7_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 4);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const size_r3 = ctx.$implicit;
    \u0275\u0275property("value", size_r3);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", size_r3, " ");
  }
}
function PaginatorComponent_ng_container_0_Conditional_9_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275text(0);
    \u0275\u0275pipe(1, "number");
    \u0275\u0275pipe(2, "number");
    \u0275\u0275pipe(3, "intEstimate");
  }
  if (rf & 2) {
    let tmp_3_0;
    const t_r4 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275textInterpolate1(" ", t_r4("paginator.x_to_y_of_z", \u0275\u0275pureFunction3(8, _c0, \u0275\u0275pipeBind1(1, 1, ctx_r1.firstItemIndex), \u0275\u0275pipeBind1(2, 3, ctx_r1.lastItemIndex), \u0275\u0275pipeBind2(3, 5, (tmp_3_0 = ctx_r1.totalLength) !== null && tmp_3_0 !== void 0 ? tmp_3_0 : 0, ctx_r1.totalIsEstimate))), " ");
  }
}
function PaginatorComponent_ng_container_0_Conditional_10_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275text(0);
    \u0275\u0275pipe(1, "number");
    \u0275\u0275pipe(2, "number");
  }
  if (rf & 2) {
    const t_r4 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275textInterpolate1(" ", t_r4("paginator.x_to_y", \u0275\u0275pureFunction2(5, _c1, \u0275\u0275pipeBind1(1, 1, ctx_r1.firstItemIndex), \u0275\u0275pipeBind1(2, 3, ctx_r1.lastItemIndex))), " ");
  }
}
function PaginatorComponent_ng_container_0_Conditional_21_Template(rf, ctx) {
  if (rf & 1) {
    const _r5 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "button", 7);
    \u0275\u0275listener("click", function PaginatorComponent_ng_container_0_Conditional_21_Template_button_click_0_listener() {
      let tmp_4_0;
      \u0275\u0275restoreView(_r5);
      const ctx_r1 = \u0275\u0275nextContext(2);
      ctx_r1.page = (tmp_4_0 = ctx_r1.pageCount) !== null && tmp_4_0 !== void 0 ? tmp_4_0 : 1;
      return \u0275\u0275resetView(ctx_r1.emitChange());
    });
    \u0275\u0275elementStart(1, "mat-icon");
    \u0275\u0275text(2, "last_page");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const t_r4 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275property("disabled", \u0275\u0275pureFunction1(2, _c2, ctx_r1.page).includes(ctx_r1.pageCount))("matTooltip", t_r4("paginator.last_page"));
  }
}
function PaginatorComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "div", 1)(2, "mat-form-field", 2)(3, "mat-label");
    \u0275\u0275text(4, "Items per page");
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(5, "mat-select", 3);
    \u0275\u0275listener("valueChange", function PaginatorComponent_ng_container_0_Template_mat_select_valueChange_5_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      ctx_r1.pageSize = $event;
      ctx_r1.page = 1;
      return \u0275\u0275resetView(ctx_r1.emitChange());
    });
    \u0275\u0275repeaterCreate(6, PaginatorComponent_ng_container_0_For_7_Template, 2, 2, "mat-option", 4, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(8, "p", 5);
    \u0275\u0275template(9, PaginatorComponent_ng_container_0_Conditional_9_Template, 4, 12)(10, PaginatorComponent_ng_container_0_Conditional_10_Template, 3, 8);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(11, "div", 6)(12, "button", 7);
    \u0275\u0275listener("click", function PaginatorComponent_ng_container_0_Template_button_click_12_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      ctx_r1.page = 1;
      return \u0275\u0275resetView(ctx_r1.emitChange());
    });
    \u0275\u0275elementStart(13, "mat-icon");
    \u0275\u0275text(14, "first_page");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(15, "button", 7);
    \u0275\u0275listener("click", function PaginatorComponent_ng_container_0_Template_button_click_15_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      ctx_r1.page = ctx_r1.page - 1;
      return \u0275\u0275resetView(ctx_r1.emitChange());
    });
    \u0275\u0275elementStart(16, "mat-icon");
    \u0275\u0275text(17, "navigate_before");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(18, "button", 7);
    \u0275\u0275listener("click", function PaginatorComponent_ng_container_0_Template_button_click_18_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      ctx_r1.page = ctx_r1.page + 1;
      return \u0275\u0275resetView(ctx_r1.emitChange());
    });
    \u0275\u0275elementStart(19, "mat-icon");
    \u0275\u0275text(20, "navigate_next");
    \u0275\u0275elementEnd()();
    \u0275\u0275template(21, PaginatorComponent_ng_container_0_Conditional_21_Template, 3, 4, "button", 8);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r4 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance(5);
    \u0275\u0275property("value", ctx_r1.pageSize);
    \u0275\u0275advance();
    \u0275\u0275repeater(ctx_r1.pageSizes);
    \u0275\u0275advance(3);
    \u0275\u0275conditional(ctx_r1.hasTotalLength ? 9 : 10);
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", !ctx_r1.hasPreviousPage)("matTooltip", t_r4("paginator.first_page"));
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", !ctx_r1.hasPreviousPage)("matTooltip", t_r4("paginator.previous_page"));
    \u0275\u0275advance(3);
    \u0275\u0275property("disabled", !ctx_r1.actuallyHasNextPage)("matTooltip", t_r4("paginator.next_page"));
    \u0275\u0275advance(3);
    \u0275\u0275conditional(ctx_r1.showLastPage ? 21 : -1);
  }
}
var PaginatorComponent = class _PaginatorComponent {
  constructor() {
    this.page = 1;
    this.pageSize = 10;
    this.pageSizes = [10, 20, 50, 100];
    this.pageLength = 0;
    this.totalLength = null;
    this.totalIsEstimate = false;
    this.hasNextPage = null;
    this.showLastPage = false;
    this.paging = new EventEmitter();
  }
  get firstItemIndex() {
    return (this.page - 1) * this.pageSize + 1;
  }
  get lastItemIndex() {
    return (this.page - 1) * this.pageSize + this.pageLength;
  }
  get hasTotalLength() {
    return typeof this.totalLength === "number";
  }
  get hasPreviousPage() {
    return this.page > 1;
  }
  get pageCount() {
    if (typeof this.totalLength !== "number") {
      return null;
    }
    return Math.ceil(this.totalLength / this.pageSize);
  }
  get actuallyHasNextPage() {
    if (typeof this.hasNextPage === "boolean") {
      return this.hasNextPage;
    }
    if (typeof this.totalLength !== "number") {
      return false;
    }
    return this.page * this.pageSize < this.totalLength;
  }
  emitChange() {
    this.paging.emit({
      page: this.page,
      pageSize: this.pageSize
    });
  }
  static {
    this.\u0275fac = function PaginatorComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _PaginatorComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _PaginatorComponent, selectors: [["app-paginator"]], inputs: { page: [2, "page", "page", numberAttribute], pageSize: [2, "pageSize", "pageSize", numberAttribute], pageSizes: "pageSizes", pageLength: [2, "pageLength", "pageLength", numberAttribute], totalLength: "totalLength", totalIsEstimate: "totalIsEstimate", hasNextPage: "hasNextPage", showLastPage: "showLastPage" }, outputs: { paging: "paging" }, standalone: true, features: [\u0275\u0275InputTransformsFeature, \u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], [1, "paginator"], ["subscriptSizing", "dynamic", 1, "field-items-per-page"], [3, "valueChange", "value"], [3, "value"], [1, "paginator-description"], [1, "paginator-navigation"], ["mat-icon-button", "", 3, "click", "disabled", "matTooltip"], ["mat-icon-button", "", 3, "disabled", "matTooltip"]], template: function PaginatorComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, PaginatorComponent_ng_container_0_Template, 22, 9, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatOption, MatIconButton, MatFormField, MatLabel, MatIcon, MatSelect, MatTooltip, TranslocoDirective, DecimalPipe, IntEstimatePipe], styles: ["\n\n.paginator[_ngcontent-%COMP%]    > *[_ngcontent-%COMP%] {\n  display: inline-block;\n  vertical-align: middle;\n}\n.paginator[_ngcontent-%COMP%]   p[_ngcontent-%COMP%] {\n  margin: 0 20px;\n}\n.paginator[_ngcontent-%COMP%]   .field-items-per-page[_ngcontent-%COMP%] {\n  width: 140px;\n}\n/*# sourceMappingURL=paginator.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(PaginatorComponent, { className: "PaginatorComponent", filePath: "src/app/paginator/paginator.component.ts", lineNumber: 19 });
})();

// src/app/pipes/time-ago.pipe.ts
var TimeAgoPipe = class _TimeAgoPipe {
  constructor() {
    this.transloco = inject(TranslocoService);
  }
  transform(value) {
    return formatTimeAgo(value, this.transloco.getActiveLang());
  }
  static {
    this.\u0275fac = function TimeAgoPipe_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _TimeAgoPipe)();
    };
  }
  static {
    this.\u0275pipe = /* @__PURE__ */ \u0275\u0275definePipe({ name: "timeAgo", type: _TimeAgoPipe, pure: false, standalone: true });
  }
};

export {
  IntEstimatePipe,
  PaginatorComponent,
  TimeAgoPipe
};
//# sourceMappingURL=chunk-43HRGFU3.js.map
