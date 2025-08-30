import {
  HealthCardComponent,
  HealthModule
} from "./chunk-WNZLSZUT.js";
import {
  DocumentTitleComponent
} from "./chunk-OOXMQI6S.js";
import {
  BreakpointsService
} from "./chunk-NQ6E5D5R.js";
import {
  AppModule,
  MatCard,
  MatCardContent,
  MatCardHeader,
  MatDivider,
  MatGridList,
  MatGridTile,
  MatIcon,
  MatToolbar,
  TranslocoDirective
} from "./chunk-WWRDQTKJ.js";
import "./chunk-VSVMRYN2.js";
import "./chunk-Y2ZC5Z2X.js";
import {
  inject,
  ɵsetClassDebugInfo,
  ɵɵStandaloneFeature,
  ɵɵadvance,
  ɵɵdefineComponent,
  ɵɵelement,
  ɵɵelementContainerEnd,
  ɵɵelementContainerStart,
  ɵɵelementEnd,
  ɵɵelementStart,
  ɵɵnextContext,
  ɵɵproperty,
  ɵɵpureFunction1,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate
} from "./chunk-DMMUMX3A.js";

// src/app/dashboard/dashboard-home.component.ts
var _c0 = (a0) => [a0];
function DashboardHomeComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275element(1, "app-document-title", 1);
    \u0275\u0275elementStart(2, "mat-card", 2)(3, "mat-card-header")(4, "mat-toolbar")(5, "h2")(6, "mat-icon");
    \u0275\u0275text(7, "dashboard");
    \u0275\u0275elementEnd();
    \u0275\u0275text(8);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(9, "mat-card-content");
    \u0275\u0275element(10, "mat-divider");
    \u0275\u0275elementStart(11, "div", 3)(12, "mat-grid-list", 4)(13, "mat-grid-tile", 5);
    \u0275\u0275element(14, "app-health-card");
    \u0275\u0275elementEnd()()()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r1 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275property("parts", \u0275\u0275pureFunction1(5, _c0, t_r1("routes.dashboard")));
    \u0275\u0275advance(7);
    \u0275\u0275textInterpolate(t_r1("routes.dashboard"));
    \u0275\u0275advance(4);
    \u0275\u0275property("cols", ctx_r1.breakpoints.sizeAtLeast("Medium") ? 2 : 1);
    \u0275\u0275advance();
    \u0275\u0275property("colspan", 1)("rowspan", 1);
  }
}
var DashboardHomeComponent = class _DashboardHomeComponent {
  constructor() {
    this.breakpoints = inject(BreakpointsService);
  }
  static {
    this.\u0275fac = function DashboardHomeComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _DashboardHomeComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _DashboardHomeComponent, selectors: [["app-dashboard"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], [3, "parts"], [1, "dashboard-card"], [1, "grid-container"], ["rowHeight", "500px", 3, "cols"], [3, "colspan", "rowspan"]], template: function DashboardHomeComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, DashboardHomeComponent_ng_container_0_Template, 15, 7, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatCard, MatCardContent, MatCardHeader, MatDivider, MatGridList, MatGridTile, MatIcon, MatToolbar, TranslocoDirective, HealthModule, HealthCardComponent, DocumentTitleComponent], styles: ["\n\n.grid-container[_ngcontent-%COMP%] {\n  margin: 20px;\n}\n.more-button[_ngcontent-%COMP%] {\n  position: absolute;\n  top: 5px;\n  right: 10px;\n}\napp-health-card[_ngcontent-%COMP%] {\n  width: 100%;\n  height: 100%;\n}\napp-health-card[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%] {\n  height: 100%;\n}\nmat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%] {\n  width: 100%;\n}\nmat-toolbar[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  position: relative;\n  top: 3px;\n  margin-right: 14px;\n  margin-left: 32px;\n}\n/*# sourceMappingURL=dashboard-home.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(DashboardHomeComponent, { className: "DashboardHomeComponent", filePath: "src/app/dashboard/dashboard-home.component.ts", lineNumber: 14 });
})();
export {
  DashboardHomeComponent
};
//# sourceMappingURL=chunk-BI2TVXFB.js.map
