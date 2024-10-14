import { inject, Injectable } from "@angular/core";
import { BreakpointObserver, Breakpoints } from "@angular/cdk/layout";
import { Observable } from "rxjs";
import { map, shareReplay } from "rxjs/operators";

const sizes = ["XSmall", "Small", "Medium", "Large", "XLarge"] as const;

type Size = (typeof sizes)[number];

@Injectable({ providedIn: "root" })
export class BreakpointsService {
  private breakpointObserver = inject(BreakpointObserver);

  private state: Observable<Record<string, boolean>> = this.breakpointObserver
    .observe([
      Breakpoints.XSmall,
      Breakpoints.Small,
      Breakpoints.Medium,
      Breakpoints.Large,
      Breakpoints.XLarge,
    ])
    .pipe(
      map((result) => result.breakpoints),
      shareReplay(),
    );

  size$: Observable<Size> = this.state.pipe(
    map((st) => sizes.find((s) => st[Breakpoints[s]]) ?? "Medium"),
  );

  size: Size = "Medium";

  constructor() {
    this.size$.subscribe((s) => {
      this.size = s;
    });
  }

  sizeAtLeast(size: Size) {
    return sizes.indexOf(size) <= sizes.indexOf(this.size);
  }
}
