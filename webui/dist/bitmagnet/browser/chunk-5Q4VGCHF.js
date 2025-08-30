import {
  DocumentTitleComponent
} from "./chunk-OOXMQI6S.js";
import {
  AppModule,
  MatCard,
  MatCardHeader,
  MatCardTitle,
  TranslocoDirective
} from "./chunk-WWRDQTKJ.js";
import "./chunk-VSVMRYN2.js";
import "./chunk-Y2ZC5Z2X.js";
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
  ɵɵpureFunction1,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate
} from "./chunk-DMMUMX3A.js";

// src/app/not-found/not-found.component.ts
var _c0 = (a0) => [a0];
function NotFoundComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275element(1, "app-document-title", 1);
    \u0275\u0275elementStart(2, "mat-card", 2)(3, "mat-card-header")(4, "mat-card-title")(5, "h2");
    \u0275\u0275text(6);
    \u0275\u0275elementEnd()()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r1 = ctx.$implicit;
    \u0275\u0275advance();
    \u0275\u0275property("parts", \u0275\u0275pureFunction1(2, _c0, t_r1("general.page_not_found")));
    \u0275\u0275advance(5);
    \u0275\u0275textInterpolate(t_r1("general.page_not_found"));
  }
}
var NotFoundComponent = class _NotFoundComponent {
  static {
    this.\u0275fac = function NotFoundComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _NotFoundComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _NotFoundComponent, selectors: [["app-not-found"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], [3, "parts"], [1, "card-not-found"]], template: function NotFoundComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, NotFoundComponent_ng_container_0_Template, 7, 4, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatCard, MatCardHeader, MatCardTitle, TranslocoDirective, DocumentTitleComponent], styles: ["\n\n.card-not-found[_ngcontent-%COMP%] {\n  max-width: 960px;\n  margin: 20px auto;\n}\n.card-not-found[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%] {\n  margin-top: 10px;\n}\n/*# sourceMappingURL=not-found.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(NotFoundComponent, { className: "NotFoundComponent", filePath: "src/app/not-found/not-found.component.ts", lineNumber: 12 });
})();
export {
  NotFoundComponent
};
//# sourceMappingURL=chunk-5Q4VGCHF.js.map
