import {
  AppModule,
  MatCard,
  MatCardContent,
  MatCardHeader,
  MatIcon,
  MatTabLink,
  MatTabNav,
  MatTabNavPanel,
  MatToolbar,
  TranslocoDirective
} from "./chunk-WWRDQTKJ.js";
import "./chunk-VSVMRYN2.js";
import {
  RouterLink,
  RouterLinkActive,
  RouterOutlet
} from "./chunk-Y2ZC5Z2X.js";
import {
  ɵsetClassDebugInfo,
  ɵɵStandaloneFeature,
  ɵɵadvance,
  ɵɵdefineComponent,
  ɵɵelement,
  ɵɵelementContainerEnd,
  ɵɵelementContainerStart,
  ɵɵelementEnd,
  ɵɵelementStart,
  ɵɵproperty,
  ɵɵreference,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate
} from "./chunk-DMMUMX3A.js";

// src/app/dashboard/queue/queue-dashboard.component.ts
function QueueDashboardComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-card", 6)(2, "mat-card-header")(3, "mat-toolbar")(4, "nav", 7)(5, "h2");
    \u0275\u0275element(6, "mat-icon", 8);
    \u0275\u0275text(7);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(8, "a", 9, 0)(10, "mat-icon");
    \u0275\u0275text(11, "monitoring");
    \u0275\u0275elementEnd();
    \u0275\u0275text(12);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(13, "a", 10, 1)(15, "mat-icon");
    \u0275\u0275text(16, "toc");
    \u0275\u0275elementEnd();
    \u0275\u0275text(17);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(18, "a", 11, 2)(20, "mat-icon");
    \u0275\u0275text(21, "construction");
    \u0275\u0275elementEnd();
    \u0275\u0275text(22);
    \u0275\u0275elementEnd()();
    \u0275\u0275element(23, "mat-tab-nav-panel", null, 3);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(25, "mat-card-content");
    \u0275\u0275element(26, "router-outlet", null, 4);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r1 = ctx.$implicit;
    const linkVisualize_r2 = \u0275\u0275reference(9);
    const linkJobs_r3 = \u0275\u0275reference(14);
    const linkAdmin_r4 = \u0275\u0275reference(19);
    const tabPanel_r5 = \u0275\u0275reference(24);
    \u0275\u0275advance(4);
    \u0275\u0275property("tabPanel", tabPanel_r5);
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate(t_r1("routes.queues"));
    \u0275\u0275advance();
    \u0275\u0275property("active", linkVisualize_r2.isActive);
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(t_r1("routes.visualize"));
    \u0275\u0275advance();
    \u0275\u0275property("active", linkJobs_r3.isActive);
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(t_r1("routes.jobs"));
    \u0275\u0275advance();
    \u0275\u0275property("active", linkAdmin_r4.isActive);
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(t_r1("routes.admin"));
  }
}
var QueueDashboardComponent = class _QueueDashboardComponent {
  static {
    this.\u0275fac = function QueueDashboardComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _QueueDashboardComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _QueueDashboardComponent, selectors: [["app-queue-dashboard"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [["linkVisualize", "routerLinkActive"], ["linkJobs", "routerLinkActive"], ["linkAdmin", "routerLinkActive"], ["tabPanel", ""], ["outlet", ""], [4, "transloco"], [1, "dashboard-card"], ["mat-tab-nav-bar", "", "mat-stretch-tabs", "false", "mat-align-tabs", "start", 3, "tabPanel"], ["svgIcon", "queue"], ["mat-tab-link", "", "routerLink", "visualize", "routerLinkActive", "", 3, "active"], ["mat-tab-link", "", "routerLink", "jobs", "routerLinkActive", "", 3, "active"], ["mat-tab-link", "", "routerLink", "admin", "routerLinkActive", "", 3, "active"]], template: function QueueDashboardComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, QueueDashboardComponent_ng_container_0_Template, 28, 8, "ng-container", 5);
      }
    }, dependencies: [AppModule, MatCard, MatCardContent, MatCardHeader, MatIcon, MatTabNav, MatTabNavPanel, MatTabLink, MatToolbar, RouterOutlet, RouterLink, RouterLinkActive, TranslocoDirective], styles: ["\n\nmat-card-header[_ngcontent-%COMP%] {\n  flex-wrap: wrap;\n}\nmat-card-header[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%] {\n  font-size: 18px;\n  margin: 0 60px 0 48px;\n  height: 48px;\n  line-height: 48px;\n}\nmat-card-header[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  position: relative;\n  top: 6px;\n  margin-right: 14px;\n  line-height: 1.25rem;\n}\nmat-card-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%] {\n  flex: 0 0 100%;\n}\nmat-card-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   a[_ngcontent-%COMP%] {\n  margin-top: 2px;\n}\nmat-card-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   a[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  margin-right: 12px;\n}\n/*# sourceMappingURL=queue-dashboard.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(QueueDashboardComponent, { className: "QueueDashboardComponent", filePath: "src/app/dashboard/queue/queue-dashboard.component.ts", lineNumber: 11 });
})();
export {
  QueueDashboardComponent
};
//# sourceMappingURL=chunk-2IVFA22B.js.map
