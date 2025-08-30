import {
  BreakpointObserver,
  Breakpoints
} from "./chunk-WWRDQTKJ.js";
import {
  inject,
  map,
  shareReplay,
  ɵɵdefineInjectable
} from "./chunk-DMMUMX3A.js";

// src/app/layout/breakpoints.service.ts
var sizes = ["XSmall", "Small", "Medium", "Large", "XLarge"];
var BreakpointsService = class _BreakpointsService {
  constructor() {
    this.breakpointObserver = inject(BreakpointObserver);
    this.state = this.breakpointObserver.observe([
      Breakpoints.XSmall,
      Breakpoints.Small,
      Breakpoints.Medium,
      Breakpoints.Large,
      Breakpoints.XLarge
    ]).pipe(map((result) => result.breakpoints), shareReplay());
    this.size$ = this.state.pipe(map((st) => sizes.find((s) => st[Breakpoints[s]]) ?? "Medium"));
    this.size = "Medium";
    this.size$.subscribe((s) => {
      this.size = s;
    });
  }
  sizeAtLeast(size) {
    return sizes.indexOf(size) <= sizes.indexOf(this.size);
  }
  static {
    this.\u0275fac = function BreakpointsService_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _BreakpointsService)();
    };
  }
  static {
    this.\u0275prov = /* @__PURE__ */ \u0275\u0275defineInjectable({ token: _BreakpointsService, factory: _BreakpointsService.\u0275fac, providedIn: "root" });
  }
};

export {
  BreakpointsService
};
//# sourceMappingURL=chunk-NQ6E5D5R.js.map
