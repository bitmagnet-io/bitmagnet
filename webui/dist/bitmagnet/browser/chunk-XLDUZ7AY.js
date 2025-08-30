import "./chunk-QLUBXXM3.js";
import "./chunk-LUZJBAO3.js";
import "./chunk-43HRGFU3.js";
import "./chunk-ORIQXXAG.js";
import "./chunk-UGVUNZOV.js";
import "./chunk-3D6CEWET.js";
import "./chunk-75G4HS47.js";
import {
  BreakpointsService
} from "./chunk-NQ6E5D5R.js";
import {
  AppModule,
  MatAnchor,
  MatDrawer,
  MatDrawerContainer,
  MatDrawerContent,
  MatIcon,
  MatIconButton,
  MatTooltip,
  TranslocoDirective
} from "./chunk-WWRDQTKJ.js";
import "./chunk-VSVMRYN2.js";
import {
  RouterLink,
  RouterLinkActive,
  RouterOutlet
} from "./chunk-Y2ZC5Z2X.js";
import {
  inject,
  ɵsetClassDebugInfo,
  ɵɵStandaloneFeature,
  ɵɵadvance,
  ɵɵattribute,
  ɵɵclassMap,
  ɵɵdefineComponent,
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
  ɵɵreference,
  ɵɵresetView,
  ɵɵrestoreView,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate
} from "./chunk-DMMUMX3A.js";

// src/app/dashboard/dashboard.component.ts
var _c0 = () => ({ exact: true });
function DashboardComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-drawer-container", 7)(2, "mat-drawer", 8, 0)(4, "nav")(5, "ul")(6, "li")(7, "a", 9, 1)(9, "mat-icon");
    \u0275\u0275text(10, "dashboard");
    \u0275\u0275elementEnd();
    \u0275\u0275text(11);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(12, "li")(13, "a", 10, 2)(15, "mat-icon");
    \u0275\u0275text(16, "manufacturing");
    \u0275\u0275elementEnd();
    \u0275\u0275text(17);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(18, "li")(19, "a", 11, 3);
    \u0275\u0275element(21, "mat-icon", 12);
    \u0275\u0275text(22);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(23, "li")(24, "a", 13, 4);
    \u0275\u0275element(26, "mat-icon", 14);
    \u0275\u0275text(27);
    \u0275\u0275elementEnd()()()()();
    \u0275\u0275elementStart(28, "mat-drawer-content")(29, "div", 15)(30, "button", 16);
    \u0275\u0275listener("click", function DashboardComponent_ng_container_0_Template_button_click_30_listener() {
      \u0275\u0275restoreView(_r1);
      const drawer_r2 = \u0275\u0275reference(3);
      return \u0275\u0275resetView(drawer_r2.toggle());
    });
    \u0275\u0275elementStart(31, "mat-icon", 17);
    \u0275\u0275text(32);
    \u0275\u0275elementEnd()()();
    \u0275\u0275element(33, "router-outlet", null, 5);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r3 = ctx.$implicit;
    const drawer_r2 = \u0275\u0275reference(3);
    const linkHome_r4 = \u0275\u0275reference(8);
    const linkWorkers_r5 = \u0275\u0275reference(14);
    const linkQueues_r6 = \u0275\u0275reference(20);
    const linkTorrents_r7 = \u0275\u0275reference(25);
    const ctx_r7 = \u0275\u0275nextContext();
    \u0275\u0275advance(2);
    \u0275\u0275property("mode", ctx_r7.breakpoints.sizeAtLeast("Medium") ? "side" : "over")("opened", ctx_r7.breakpoints.sizeAtLeast("Medium"));
    \u0275\u0275attribute("role", ctx_r7.breakpoints.sizeAtLeast("Medium") ? "navigation" : "dialog");
    \u0275\u0275advance(5);
    \u0275\u0275classMap(linkHome_r4.isActive ? "active" : "");
    \u0275\u0275property("routerLinkActiveOptions", \u0275\u0275pureFunction0(18, _c0));
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(t_r3("routes.home"));
    \u0275\u0275advance(2);
    \u0275\u0275classMap(linkWorkers_r5.isActive ? "active" : "");
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(t_r3("routes.workers"));
    \u0275\u0275advance(2);
    \u0275\u0275classMap(linkQueues_r6.isActive ? "active" : "");
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate(t_r3("routes.queues"));
    \u0275\u0275advance(2);
    \u0275\u0275classMap(linkTorrents_r7.isActive ? "active" : "");
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate(t_r3("routes.torrents"));
    \u0275\u0275advance(3);
    \u0275\u0275property("matTooltip", t_r3("torrents.toggle_drawer"));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(drawer_r2.opened ? "arrow_circle_left" : "arrow_circle_right");
  }
}
var DashboardComponent = class _DashboardComponent {
  constructor() {
    this.breakpoints = inject(BreakpointsService);
  }
  static {
    this.\u0275fac = function DashboardComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _DashboardComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _DashboardComponent, selectors: [["app-dashboard"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [["drawer", ""], ["linkHome", "routerLinkActive"], ["linkWorkers", "routerLinkActive"], ["linkQueues", "routerLinkActive"], ["linkTorrents", "routerLinkActive"], ["outlet", ""], [4, "transloco"], [1, "drawer-container"], [1, "drawer", 3, "mode", "opened"], ["mat-button", "", "routerLink", "/dashboard", "routerLinkActive", "", 3, "routerLinkActiveOptions"], ["mat-button", "", "routerLink", "workers", "routerLinkActive", ""], ["mat-button", "", "routerLink", "queues", "routerLinkActive", ""], ["svgIcon", "queue"], ["mat-button", "", "routerLink", "torrents", "routerLinkActive", ""], ["svgIcon", "magnet"], [1, "form-field-container", "button-container", "button-container-toggle-drawer"], ["type", "button", "mat-icon-button", "", 1, "button-toggle-drawer", 3, "click", "matTooltip"], ["aria-label", "Side nav toggle icon", "fontSet", "material-icons"]], template: function DashboardComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, DashboardComponent_ng_container_0_Template, 35, 19, "ng-container", 6);
      }
    }, dependencies: [AppModule, MatAnchor, MatIconButton, MatIcon, MatDrawer, MatDrawerContainer, MatDrawerContent, MatTooltip, RouterOutlet, RouterLink, RouterLinkActive, TranslocoDirective], styles: ["\n\nmat-drawer[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%] {\n  padding-top: 12px;\n  --mat-text-button-icon-spacing: 14px;\n}\nmat-drawer[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   ul[_ngcontent-%COMP%] {\n  list-style-type: none;\n  padding-left: 0;\n}\nmat-drawer[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   ul[_ngcontent-%COMP%]   a[_ngcontent-%COMP%] {\n  width: 100%;\n  font-size: var(--mat-expansion-container-text-size);\n  justify-content: flex-start;\n  padding-left: 20px;\n}\nmat-drawer[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   ul[_ngcontent-%COMP%]   li[_ngcontent-%COMP%] {\n  margin-bottom: 6px;\n}\nmat-drawer-content[_ngcontent-%COMP%]   .button-container-toggle-drawer[_ngcontent-%COMP%] {\n  position: absolute;\n  left: 20px;\n  top: 28px;\n  z-index: 100;\n}\n.drawer[_ngcontent-%COMP%] {\n  width: 220px;\n}\n/*# sourceMappingURL=dashboard.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(DashboardComponent, { className: "DashboardComponent", filePath: "src/app/dashboard/dashboard.component.ts", lineNumber: 20 });
})();
export {
  DashboardComponent
};
//# sourceMappingURL=chunk-XLDUZ7AY.js.map
