import {
  Title
} from "./chunk-Y2ZC5Z2X.js";
import {
  inject,
  ɵsetClassDebugInfo,
  ɵɵNgOnChangesFeature,
  ɵɵStandaloneFeature,
  ɵɵdefineComponent,
  ɵɵelementContainer
} from "./chunk-DMMUMX3A.js";

// src/app/layout/document-title.component.ts
var DocumentTitleComponent = class _DocumentTitleComponent {
  constructor() {
    this.title = inject(Title);
    this.parts = [];
  }
  ngOnInit() {
    this.updateTitle();
  }
  ngOnChanges() {
    this.updateTitle();
  }
  updateTitle() {
    this.title.setTitle([...this.parts.filter(Boolean), "bitmagnet"].join(" - "));
  }
  static {
    this.\u0275fac = function DocumentTitleComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _DocumentTitleComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _DocumentTitleComponent, selectors: [["app-document-title"]], inputs: { parts: "parts" }, standalone: true, features: [\u0275\u0275NgOnChangesFeature, \u0275\u0275StandaloneFeature], decls: 1, vars: 0, template: function DocumentTitleComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275elementContainer(0);
      }
    }, encapsulation: 2 });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(DocumentTitleComponent, { className: "DocumentTitleComponent", filePath: "src/app/layout/document-title.component.ts", lineNumber: 9 });
})();

export {
  DocumentTitleComponent
};
//# sourceMappingURL=chunk-OOXMQI6S.js.map
