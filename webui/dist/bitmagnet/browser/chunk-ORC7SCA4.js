import {
  PaginatorComponent,
  TimeAgoPipe
} from "./chunk-43HRGFU3.js";
import "./chunk-ORIQXXAG.js";
import "./chunk-3D6CEWET.js";
import {
  ErrorsService
} from "./chunk-75G4HS47.js";
import {
  DocumentTitleComponent
} from "./chunk-OOXMQI6S.js";
import {
  BreakpointsService
} from "./chunk-NQ6E5D5R.js";
import {
  Apollo,
  AppModule,
  CdkCopyToClipboard,
  MatCard,
  MatCardContent,
  MatCardHeader,
  MatCardTitle,
  MatCell,
  MatCellDef,
  MatCheckbox,
  MatColumnDef,
  MatDrawer,
  MatDrawerContainer,
  MatDrawerContent,
  MatFormField,
  MatHeaderCell,
  MatHeaderCellDef,
  MatHeaderRow,
  MatHeaderRowDef,
  MatIcon,
  MatIconButton,
  MatLabel,
  MatMiniFabButton,
  MatOption,
  MatProgressBar,
  MatRow,
  MatRowDef,
  MatSelect,
  MatTable,
  MatTooltip,
  QueueJobsDocument,
  SelectionModel,
  TranslocoDirective,
  TranslocoService
} from "./chunk-WWRDQTKJ.js";
import {
  animate,
  state,
  style,
  transition,
  trigger
} from "./chunk-VSVMRYN2.js";
import "./chunk-Y2ZC5Z2X.js";
import {
  AsyncPipe,
  BehaviorSubject,
  DecimalPipe,
  EMPTY,
  EventEmitter,
  SlicePipe,
  __spreadProps,
  __spreadValues,
  catchError,
  combineLatestWith,
  debounceTime,
  inject,
  map,
  ɵsetClassDebugInfo,
  ɵɵStandaloneFeature,
  ɵɵadvance,
  ɵɵattribute,
  ɵɵclassMap,
  ɵɵconditional,
  ɵɵdefineComponent,
  ɵɵelement,
  ɵɵelementContainerEnd,
  ɵɵelementContainerStart,
  ɵɵelementEnd,
  ɵɵelementStart,
  ɵɵgetCurrentView,
  ɵɵlistener,
  ɵɵnextContext,
  ɵɵpipe,
  ɵɵpipeBind1,
  ɵɵpipeBind3,
  ɵɵproperty,
  ɵɵpureFunction0,
  ɵɵpureFunction3,
  ɵɵreference,
  ɵɵrepeater,
  ɵɵrepeaterCreate,
  ɵɵresetView,
  ɵɵrestoreView,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate,
  ɵɵtextInterpolate1
} from "./chunk-DMMUMX3A.js";

// src/app/dashboard/queue/queue-jobs-table.component.ts
var _c0 = () => ["expandedDetail"];
function QueueJobsTableComponent_ng_container_0_th_6_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 18);
    \u0275\u0275text(1, "ID");
    \u0275\u0275elementEnd();
  }
}
function QueueJobsTableComponent_ng_container_0_td_7_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 19);
    \u0275\u0275listener("click", function QueueJobsTableComponent_ng_container_0_td_7_Template_td_click_0_listener($event) {
      const i_r2 = \u0275\u0275restoreView(_r1).$implicit;
      const ctx_r2 = \u0275\u0275nextContext(2);
      ctx_r2.toggleQueueJobId(i_r2.id);
      return \u0275\u0275resetView($event.stopPropagation());
    });
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const i_r2 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", ctx_r2.item(i_r2).id, " ");
  }
}
function QueueJobsTableComponent_ng_container_0_th_9_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 18);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r4("dashboard.queues.queue"), " ");
  }
}
function QueueJobsTableComponent_ng_container_0_td_10_Template(rf, ctx) {
  if (rf & 1) {
    const _r5 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 19);
    \u0275\u0275listener("click", function QueueJobsTableComponent_ng_container_0_td_10_Template_td_click_0_listener($event) {
      const i_r6 = \u0275\u0275restoreView(_r5).$implicit;
      const ctx_r2 = \u0275\u0275nextContext(2);
      ctx_r2.toggleQueueJobId(i_r6.id);
      return \u0275\u0275resetView($event.stopPropagation());
    });
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const i_r6 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", ctx_r2.item(i_r6).queue, " ");
  }
}
function QueueJobsTableComponent_ng_container_0_th_12_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 18);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r4("dashboard.queues.priority"), " ");
  }
}
function QueueJobsTableComponent_ng_container_0_td_13_Template(rf, ctx) {
  if (rf & 1) {
    const _r7 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 19);
    \u0275\u0275listener("click", function QueueJobsTableComponent_ng_container_0_td_13_Template_td_click_0_listener($event) {
      const i_r8 = \u0275\u0275restoreView(_r7).$implicit;
      const ctx_r2 = \u0275\u0275nextContext(2);
      ctx_r2.toggleQueueJobId(i_r8.id);
      return \u0275\u0275resetView($event.stopPropagation());
    });
    \u0275\u0275text(1);
    \u0275\u0275pipe(2, "number");
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const i_r8 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", \u0275\u0275pipeBind1(2, 1, ctx_r2.item(i_r8).priority), " ");
  }
}
function QueueJobsTableComponent_ng_container_0_th_15_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 18);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r4("general.status"), " ");
  }
}
function QueueJobsTableComponent_ng_container_0_td_16_Template(rf, ctx) {
  if (rf & 1) {
    const _r9 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 19);
    \u0275\u0275listener("click", function QueueJobsTableComponent_ng_container_0_td_16_Template_td_click_0_listener($event) {
      const i_r10 = \u0275\u0275restoreView(_r9).$implicit;
      const ctx_r2 = \u0275\u0275nextContext(2);
      ctx_r2.toggleQueueJobId(i_r10.id);
      return \u0275\u0275resetView($event.stopPropagation());
    });
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const i_r10 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", ctx_r2.item(i_r10).status, " ");
  }
}
function QueueJobsTableComponent_ng_container_0_th_18_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 18);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r4("general.error"));
  }
}
function QueueJobsTableComponent_ng_container_0_td_19_Conditional_1_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275text(0);
    \u0275\u0275pipe(1, "slice");
  }
  if (rf & 2) {
    const i_r12 = \u0275\u0275nextContext().$implicit;
    const ctx_r2 = \u0275\u0275nextContext(2);
    \u0275\u0275textInterpolate1(" ", \u0275\u0275pipeBind3(1, 1, ctx_r2.item(i_r12).error, 0, 20) + "...", " ");
  }
}
function QueueJobsTableComponent_ng_container_0_td_19_Template(rf, ctx) {
  if (rf & 1) {
    const _r11 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 19);
    \u0275\u0275listener("click", function QueueJobsTableComponent_ng_container_0_td_19_Template_td_click_0_listener($event) {
      const i_r12 = \u0275\u0275restoreView(_r11).$implicit;
      const ctx_r2 = \u0275\u0275nextContext(2);
      ctx_r2.toggleQueueJobId(i_r12.id);
      return \u0275\u0275resetView($event.stopPropagation());
    });
    \u0275\u0275template(1, QueueJobsTableComponent_ng_container_0_td_19_Conditional_1_Template, 2, 5);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const i_r12 = ctx.$implicit;
    \u0275\u0275advance();
    \u0275\u0275conditional(i_r12.error ? 1 : -1);
  }
}
function QueueJobsTableComponent_ng_container_0_th_21_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 20);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r4("dashboard.queues.created_at"), " ");
  }
}
function QueueJobsTableComponent_ng_container_0_td_22_Template(rf, ctx) {
  if (rf & 1) {
    const _r13 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 19);
    \u0275\u0275listener("click", function QueueJobsTableComponent_ng_container_0_td_22_Template_td_click_0_listener($event) {
      const i_r14 = \u0275\u0275restoreView(_r13).$implicit;
      const ctx_r2 = \u0275\u0275nextContext(2);
      ctx_r2.toggleQueueJobId(i_r14.id);
      return \u0275\u0275resetView($event.stopPropagation());
    });
    \u0275\u0275text(1);
    \u0275\u0275pipe(2, "timeAgo");
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const i_r14 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", \u0275\u0275pipeBind1(2, 1, ctx_r2.item(i_r14).createdAt), " ");
  }
}
function QueueJobsTableComponent_ng_container_0_th_24_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 20);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r4("dashboard.queues.ran_at"), " ");
  }
}
function QueueJobsTableComponent_ng_container_0_td_25_Conditional_1_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275text(0);
    \u0275\u0275pipe(1, "timeAgo");
  }
  if (rf & 2) {
    \u0275\u0275textInterpolate1(" ", \u0275\u0275pipeBind1(1, 1, ctx), " ");
  }
}
function QueueJobsTableComponent_ng_container_0_td_25_Template(rf, ctx) {
  if (rf & 1) {
    const _r15 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 19);
    \u0275\u0275listener("click", function QueueJobsTableComponent_ng_container_0_td_25_Template_td_click_0_listener($event) {
      const i_r16 = \u0275\u0275restoreView(_r15).$implicit;
      const ctx_r2 = \u0275\u0275nextContext(2);
      ctx_r2.toggleQueueJobId(i_r16.id);
      return \u0275\u0275resetView($event.stopPropagation());
    });
    \u0275\u0275template(1, QueueJobsTableComponent_ng_container_0_td_25_Conditional_1_Template, 2, 3);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    let tmp_4_0;
    const i_r16 = ctx.$implicit;
    \u0275\u0275advance();
    \u0275\u0275conditional((tmp_4_0 = i_r16.ranAt) ? 1 : -1, tmp_4_0);
  }
}
function QueueJobsTableComponent_ng_container_0_td_27_Conditional_13_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "h5")(1, "span", 23);
    \u0275\u0275text(2);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(3, "pre", 24);
    \u0275\u0275text(4);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const error_r17 = ctx;
    const t_r4 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275property("matTooltip", t_r4("torrents.copy_to_clipboard"))("cdkCopyToClipboard", error_r17);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1("", t_r4("general.error"), ":");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(error_r17);
  }
}
function QueueJobsTableComponent_ng_container_0_td_27_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "td", 21)(1, "div", 22)(2, "p")(3, "strong");
    \u0275\u0275text(4, "ID:");
    \u0275\u0275elementEnd();
    \u0275\u0275text(5, "\xA0");
    \u0275\u0275elementStart(6, "span", 23);
    \u0275\u0275text(7);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(8, "h5")(9, "span", 23);
    \u0275\u0275text(10);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(11, "pre", 24);
    \u0275\u0275text(12);
    \u0275\u0275elementEnd();
    \u0275\u0275template(13, QueueJobsTableComponent_ng_container_0_td_27_Conditional_13_Template, 5, 4);
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    let tmp_13_0;
    const i_r18 = ctx.$implicit;
    const t_r4 = \u0275\u0275nextContext().$implicit;
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275attribute("colspan", ctx_r2.displayedColumns.length);
    \u0275\u0275advance();
    \u0275\u0275property("@detailExpand", ctx_r2.expandedId.getValue() === i_r18.id ? "expanded" : "collapsed");
    \u0275\u0275advance(5);
    \u0275\u0275property("matTooltip", t_r4("torrents.copy_to_clipboard"))("cdkCopyToClipboard", ctx_r2.item(i_r18).id);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(ctx_r2.item(i_r18).id);
    \u0275\u0275advance(2);
    \u0275\u0275property("matTooltip", t_r4("torrents.copy_to_clipboard"))("cdkCopyToClipboard", ctx_r2.item(i_r18).payload);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1("", t_r4("dashboard.queues.payload"), ":");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(ctx_r2.beautifyPayload(ctx_r2.item(i_r18).payload));
    \u0275\u0275advance();
    \u0275\u0275conditional((tmp_13_0 = ctx_r2.item(i_r18).error) ? 13 : -1, tmp_13_0);
  }
}
function QueueJobsTableComponent_ng_container_0_tr_28_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275element(0, "tr", 25);
  }
}
function QueueJobsTableComponent_ng_container_0_tr_29_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275element(0, "tr", 26);
  }
  if (rf & 2) {
    const i_r19 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext(2);
    \u0275\u0275classMap("summary-row " + (i_r19.id === ctx_r2.expandedId.getValue() ? "expanded" : "collapsed"));
  }
}
function QueueJobsTableComponent_ng_container_0_tr_30_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275element(0, "tr", 26);
  }
  if (rf & 2) {
    const i_r20 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext(2);
    \u0275\u0275classMap("expanded-detail-row " + (i_r20.id === ctx_r2.expandedId.getValue() ? "expanded" : "collapsed"));
  }
}
function QueueJobsTableComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "div", 1);
    \u0275\u0275element(2, "mat-progress-bar", 2);
    \u0275\u0275pipe(3, "async");
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(4, "table", 3);
    \u0275\u0275elementContainerStart(5, 4);
    \u0275\u0275template(6, QueueJobsTableComponent_ng_container_0_th_6_Template, 2, 0, "th", 5)(7, QueueJobsTableComponent_ng_container_0_td_7_Template, 2, 1, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(8, 7);
    \u0275\u0275template(9, QueueJobsTableComponent_ng_container_0_th_9_Template, 2, 1, "th", 5)(10, QueueJobsTableComponent_ng_container_0_td_10_Template, 2, 1, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(11, 8);
    \u0275\u0275template(12, QueueJobsTableComponent_ng_container_0_th_12_Template, 2, 1, "th", 5)(13, QueueJobsTableComponent_ng_container_0_td_13_Template, 3, 3, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(14, 9);
    \u0275\u0275template(15, QueueJobsTableComponent_ng_container_0_th_15_Template, 2, 1, "th", 5)(16, QueueJobsTableComponent_ng_container_0_td_16_Template, 2, 1, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(17, 10);
    \u0275\u0275template(18, QueueJobsTableComponent_ng_container_0_th_18_Template, 2, 1, "th", 5)(19, QueueJobsTableComponent_ng_container_0_td_19_Template, 2, 1, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(20, 11);
    \u0275\u0275template(21, QueueJobsTableComponent_ng_container_0_th_21_Template, 2, 1, "th", 12)(22, QueueJobsTableComponent_ng_container_0_td_22_Template, 3, 3, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(23, 13);
    \u0275\u0275template(24, QueueJobsTableComponent_ng_container_0_th_24_Template, 2, 1, "th", 12)(25, QueueJobsTableComponent_ng_container_0_td_25_Template, 2, 1, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(26, 14);
    \u0275\u0275template(27, QueueJobsTableComponent_ng_container_0_td_27_Template, 14, 10, "td", 15);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275template(28, QueueJobsTableComponent_ng_container_0_tr_28_Template, 1, 0, "tr", 16)(29, QueueJobsTableComponent_ng_container_0_tr_29_Template, 1, 2, "tr", 17)(30, QueueJobsTableComponent_ng_container_0_tr_30_Template, 1, 2, "tr", 17);
    \u0275\u0275elementEnd();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275advance(2);
    \u0275\u0275property("mode", \u0275\u0275pipeBind1(3, 7, ctx_r2.dataSource.loading$) ? "indeterminate" : "determinate")("value", 0);
    \u0275\u0275advance(2);
    \u0275\u0275property("dataSource", ctx_r2.dataSource)("multiTemplateDataRows", true);
    \u0275\u0275advance(24);
    \u0275\u0275property("matHeaderRowDef", ctx_r2.displayedColumns);
    \u0275\u0275advance();
    \u0275\u0275property("matRowDefColumns", ctx_r2.displayedColumns);
    \u0275\u0275advance();
    \u0275\u0275property("matRowDefColumns", \u0275\u0275pureFunction0(9, _c0));
  }
}
var QueueJobsTableComponent = class _QueueJobsTableComponent {
  constructor() {
    this.transloco = inject(TranslocoService);
    this.displayedColumns = allColumns;
    this.updated = new EventEmitter();
    this.expandedId = new BehaviorSubject(null);
    this.items = Array();
  }
  ngOnInit() {
    this.dataSource.items$.subscribe((items) => {
      this.items = items;
      if (items.length) {
        const expandedId = this.expandedId.getValue();
        if (expandedId && !items.some(({ id }) => id === expandedId)) {
          this.expandedId.next(null);
        }
      }
    });
  }
  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    return this.items.every((i) => this.selection.isSelected(i.id));
  }
  /** Selects all rows if they are not all selected; otherwise clear selection. */
  toggleAllRows() {
    if (this.isAllSelected()) {
      this.selection.clear();
      return;
    }
    this.selection.select(...this.items.map((i) => i.id));
  }
  toggleQueueJobId(id) {
    if (this.expandedId.getValue() === id) {
      this.expandedId.next(null);
    } else {
      this.expandedId.next(id);
    }
  }
  /**
   * Workaround for untyped table cell definitions
   */
  item(item) {
    return item;
  }
  beautifyPayload(payload) {
    try {
      return JSON.stringify(JSON.parse(payload), null, 2);
    } catch (e) {
      return payload;
    }
  }
  static {
    this.\u0275fac = function QueueJobsTableComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _QueueJobsTableComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _QueueJobsTableComponent, selectors: [["app-queue-jobs-table"]], inputs: { dataSource: "dataSource", selection: "selection", displayedColumns: "displayedColumns" }, outputs: { updated: "updated" }, standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], [1, "progress-bar-container"], [3, "mode", "value"], ["mat-table", "", 1, "table-results", 3, "dataSource", "multiTemplateDataRows"], ["matColumnDef", "id"], ["mat-header-cell", "", 4, "matHeaderCellDef"], ["mat-cell", "", 3, "click", 4, "matCellDef"], ["matColumnDef", "queue"], ["matColumnDef", "priority"], ["matColumnDef", "status"], ["matColumnDef", "error"], ["matColumnDef", "createdAt"], ["mat-header-cell", "", "style", "text-align: center", 4, "matHeaderCellDef"], ["matColumnDef", "ranAt"], ["matColumnDef", "expandedDetail"], ["mat-cell", "", 4, "matCellDef"], ["mat-header-row", "", 4, "matHeaderRowDef"], ["mat-row", "", 3, "class", 4, "matRowDef", "matRowDefColumns"], ["mat-header-cell", ""], ["mat-cell", "", 3, "click"], ["mat-header-cell", "", 2, "text-align", "center"], ["mat-cell", ""], [1, "item-detail"], [1, "copy", 3, "matTooltip", "cdkCopyToClipboard"], [1, "payload"], ["mat-header-row", ""], ["mat-row", ""]], template: function QueueJobsTableComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, QueueJobsTableComponent_ng_container_0_Template, 31, 10, "ng-container", 0);
      }
    }, dependencies: [AppModule, CdkCopyToClipboard, MatProgressBar, MatTable, MatHeaderCellDef, MatHeaderRowDef, MatColumnDef, MatCellDef, MatRowDef, MatHeaderCell, MatCell, MatHeaderRow, MatRow, MatTooltip, TranslocoDirective, AsyncPipe, SlicePipe, DecimalPipe, TimeAgoPipe], styles: ["\n\n.item-detail[_ngcontent-%COMP%]    > [_ngcontent-%COMP%]:first-child {\n  padding-top: 20px;\n}\n.item-detail[_ngcontent-%COMP%]    > [_ngcontent-%COMP%]:last-child {\n  margin-bottom: 20px;\n}\ntr[_ngcontent-%COMP%]:not(.expanded-detail-row)   td[_ngcontent-%COMP%] {\n  cursor: pointer;\n}\ntr.expanded-detail-row[_ngcontent-%COMP%] {\n  height: 0;\n}\ntr.expanded-detail-row[_ngcontent-%COMP%]   h5[_ngcontent-%COMP%] {\n  margin: 0;\n  padding-top: 8px;\n}\ntr.expanded-detail-row[_ngcontent-%COMP%]   p[_ngcontent-%COMP%] {\n  margin: 0;\n  padding-top: 8px;\n  padding-bottom: 4px;\n}\ntr.expanded-detail-row[_ngcontent-%COMP%]   span.copy[_ngcontent-%COMP%] {\n  cursor: crosshair;\n  text-decoration: underline;\n  text-decoration-style: dotted;\n}\npre[_ngcontent-%COMP%] {\n  opacity: 0;\n  max-height: 200px;\n  max-width: 100px;\n  overflow: scroll;\n  background: rgba(119, 119, 119, 0.2);\n  padding: 10px;\n}\n.expanded-detail-row.expanded[_ngcontent-%COMP%]   pre[_ngcontent-%COMP%] {\n  opacity: 1;\n  max-width: 1200px;\n}\n/*# sourceMappingURL=queue-jobs-table.component.css.map */"], data: { animation: [
      trigger("detailExpand", [
        state("collapsed,void", style({ height: "0px", minHeight: "0" })),
        state("expanded", style({ height: "*" })),
        transition("expanded <=> collapsed", animate("225ms cubic-bezier(0.4, 0.0, 0.2, 1)"))
      ])
    ] } });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(QueueJobsTableComponent, { className: "QueueJobsTableComponent", filePath: "src/app/dashboard/queue/queue-jobs-table.component.ts", lineNumber: 41 });
})();
var allColumns = [
  "id",
  "queue",
  "priority",
  "status",
  "error",
  "createdAt",
  "ranAt"
];

// src/app/dashboard/queue/queue-jobs.datasource.ts
var emptyResult = {
  items: [],
  hasNextPage: false,
  totalCount: 0,
  aggregations: {
    queue: [],
    status: []
  }
};
var QueueJobsDatasource = class {
  constructor(apollo, errorsService, queryVariables) {
    this.apollo = apollo;
    this.errorsService = errorsService;
    this.currentRequest = new BehaviorSubject(0);
    this.loadingSubject = new BehaviorSubject(false);
    this.loading$ = this.loadingSubject.asObservable();
    this.result = emptyResult;
    this.resultSubject = new BehaviorSubject(this.result);
    this.result$ = this.resultSubject.asObservable();
    this.items$ = this.resultSubject.pipe(map((result) => result.items));
    queryVariables.subscribe((variables) => {
      this.variables = variables;
      this.loadResult(variables);
    });
    this.resultSubject.subscribe((result) => {
      this.result = result;
    });
  }
  connect({}) {
    return this.items$;
  }
  disconnect() {
    this.resultSubject.complete();
  }
  refresh() {
    if (this.variables) {
      this.loadResult(this.variables);
    }
  }
  loadResult(variables) {
    if (this.currentSubscription) {
      this.currentSubscription.unsubscribe();
      this.currentSubscription = void 0;
    }
    this.loadingSubject.next(true);
    const currentRequest = this.currentRequest.getValue() + 1;
    this.currentRequest.next(currentRequest);
    const result = this.apollo.query({
      query: QueueJobsDocument,
      variables,
      fetchPolicy: "no-cache"
    }).pipe(map((r) => r.data.queue.jobs)).pipe(catchError((err) => {
      this.errorsService.addError(`Error loading item results: ${err.message}`);
      return EMPTY;
    }));
    this.currentSubscription = result.subscribe((r) => {
      if (currentRequest === this.currentRequest.getValue()) {
        this.loadingSubject.next(false);
        this.resultSubject.next(r);
      }
    });
  }
};

// src/app/dashboard/queue/queue-jobs.controller.ts
var QueueJobsController = class {
  constructor(ctrl = initialControls) {
    this.controlsSubject = new BehaviorSubject(ctrl);
    this.controls$ = this.controlsSubject.asObservable();
    this.variablesSubject = new BehaviorSubject(controlsToQueryVariables(ctrl));
    this.variables$ = this.variablesSubject.asObservable();
    this.controls$.pipe(debounceTime(100)).subscribe((ctrl2) => {
      const currentParams = this.variablesSubject.getValue();
      const nextParams = controlsToQueryVariables(ctrl2);
      if (JSON.stringify(currentParams) !== JSON.stringify(nextParams)) {
        this.variablesSubject.next(nextParams);
      }
    });
  }
  update(fn) {
    const ctrl = this.controlsSubject.getValue();
    const next = fn(ctrl);
    if (JSON.stringify(ctrl) !== JSON.stringify(next)) {
      this.controlsSubject.next(next);
    }
  }
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  activateFilter(def, filter) {
    this.update((ctrl) => {
      const input = def.extractInput(ctrl.facets);
      return __spreadProps(__spreadValues({}, ctrl), {
        page: 1,
        facets: def.patchInput(ctrl.facets, __spreadProps(__spreadValues({}, input), {
          filter: Array.from(/* @__PURE__ */ new Set([...input.filter ?? [], filter])).sort()
        }))
      });
    });
  }
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  deactivateFilter(def, filter) {
    this.update((ctrl) => {
      const input = def.extractInput(ctrl.facets);
      const nextFilter = input.filter?.filter((value) => value !== filter);
      return __spreadProps(__spreadValues({}, ctrl), {
        page: 1,
        facets: def.patchInput(ctrl.facets, __spreadProps(__spreadValues({}, input), {
          filter: nextFilter?.length ? nextFilter : void 0
        }))
      });
    });
  }
  selectOrderBy(field) {
    const orderBy = {
      field,
      descending: orderByOptions.find((option) => option.field === field)?.descending ?? false
    };
    this.update((ctrl) => __spreadProps(__spreadValues({}, ctrl), {
      orderBy,
      page: 1
    }));
  }
  toggleOrderByDirection() {
    this.update((ctrl) => __spreadProps(__spreadValues({}, ctrl), {
      orderBy: __spreadProps(__spreadValues({}, ctrl.orderBy), {
        descending: !ctrl.orderBy.descending
      }),
      page: 1
    }));
  }
  handlePageEvent(event) {
    this.update((ctrl) => __spreadProps(__spreadValues({}, ctrl), {
      limit: event.pageSize,
      page: event.page
    }));
  }
};
var controlsToQueryVariables = (ctrl) => ({
  input: {
    limit: ctrl.limit,
    page: ctrl.page,
    totalCount: true,
    orderBy: [
      ctrl.orderBy,
      ...ctrl.orderBy.field !== "created_at" ? [
        {
          field: "created_at",
          descending: ctrl.orderBy.descending
        }
      ] : []
    ],
    queues: ctrl.queues,
    statuses: ctrl.statuses,
    facets: {
      queue: {
        aggregate: true,
        filter: ctrl.facets.queue.filter
      },
      status: {
        aggregate: true,
        filter: ctrl.facets.status.filter
      }
    }
  }
});
var orderByOptions = [
  {
    field: "created_at",
    descending: true
  },
  {
    field: "ran_at",
    descending: true
  },
  {
    field: "priority",
    descending: false
  }
];
var initialControls = {
  limit: 20,
  page: 1,
  orderBy: {
    field: "ran_at",
    descending: true
  },
  facets: {
    queue: {},
    status: {}
  }
};
var queueFacet = {
  key: "queue",
  extractInput: (f) => f.queue,
  patchInput: (f, i) => __spreadProps(__spreadValues({}, f), {
    queue: i
  }),
  extractAggregations: (aggs) => aggs.queue ?? [],
  resolveLabel: (agg) => agg.label
};
var statusFacet = {
  key: "status",
  extractInput: (f) => f.status,
  patchInput: (f, i) => __spreadProps(__spreadValues({}, f), {
    status: i
  }),
  extractAggregations: (aggs) => aggs.status ?? [],
  resolveLabel: (agg, t) => t.translate("dashboard.queues." + agg.label)
};
var facets = [queueFacet, statusFacet];

// src/app/dashboard/queue/queue-jobs.component.ts
var _forTrack0 = ($index, $item) => $item.key;
var _forTrack1 = ($index, $item) => $item.field;
var _forTrack2 = ($index, $item) => $item.value;
var _c02 = (a0, a1, a2) => [a0, a1, a2];
function QueueJobsComponent_ng_container_0_For_6_Conditional_4_For_2_Template(rf, ctx) {
  if (rf & 1) {
    const _r2 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-checkbox", 21);
    \u0275\u0275listener("change", function QueueJobsComponent_ng_container_0_For_6_Conditional_4_For_2_Template_mat_checkbox_change_0_listener($event) {
      const agg_r3 = \u0275\u0275restoreView(_r2).$implicit;
      const facet_r4 = \u0275\u0275nextContext(2).$implicit;
      const ctx_r4 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView($event.checked ? ctx_r4.controller.activateFilter(facet_r4, agg_r3.value) : ctx_r4.controller.deactivateFilter(facet_r4, agg_r3.value));
    });
    \u0275\u0275text(1);
    \u0275\u0275elementStart(2, "small");
    \u0275\u0275text(3);
    \u0275\u0275pipe(4, "number");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const agg_r3 = ctx.$implicit;
    const facet_r4 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275property("checked", facet_r4.filter == null ? null : facet_r4.filter.includes(agg_r3.value));
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", agg_r3.label, " ");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(\u0275\u0275pipeBind1(4, 3, agg_r3.count));
  }
}
function QueueJobsComponent_ng_container_0_For_6_Conditional_4_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-card-content", 18);
    \u0275\u0275repeaterCreate(1, QueueJobsComponent_ng_container_0_For_6_Conditional_4_For_2_Template, 5, 5, "mat-checkbox", 20, _forTrack2);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const facet_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275repeater(facet_r4.aggregations);
  }
}
function QueueJobsComponent_ng_container_0_For_6_Conditional_5_For_2_Template(rf, ctx) {
  if (rf & 1) {
    const _r6 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-checkbox", 23);
    \u0275\u0275listener("change", function QueueJobsComponent_ng_container_0_For_6_Conditional_5_For_2_Template_mat_checkbox_change_0_listener() {
      const agg_r7 = \u0275\u0275restoreView(_r6).$implicit;
      const facet_r4 = \u0275\u0275nextContext(2).$implicit;
      const ctx_r4 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r4.controller.activateFilter(facet_r4, agg_r7.value));
    });
    \u0275\u0275text(1);
    \u0275\u0275elementStart(2, "small");
    \u0275\u0275text(3);
    \u0275\u0275pipe(4, "number");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const agg_r7 = ctx.$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", agg_r7.label, " ");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(\u0275\u0275pipeBind1(4, 2, agg_r7.count));
  }
}
function QueueJobsComponent_ng_container_0_For_6_Conditional_5_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-card-content", 19);
    \u0275\u0275repeaterCreate(1, QueueJobsComponent_ng_container_0_For_6_Conditional_5_For_2_Template, 5, 4, "mat-checkbox", 22, _forTrack2);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const facet_r4 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275repeater(facet_r4.aggregations);
  }
}
function QueueJobsComponent_ng_container_0_For_6_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-card")(1, "mat-card-header")(2, "mat-card-title");
    \u0275\u0275text(3);
    \u0275\u0275elementEnd()();
    \u0275\u0275template(4, QueueJobsComponent_ng_container_0_For_6_Conditional_4_Template, 3, 0, "mat-card-content", 18)(5, QueueJobsComponent_ng_container_0_For_6_Conditional_5_Template, 3, 0, "mat-card-content", 19);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const facet_r4 = ctx.$implicit;
    const t_r8 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate1(" ", t_r8("facets." + facet_r4.key), " ");
    \u0275\u0275advance();
    \u0275\u0275conditional((facet_r4.filter == null ? null : facet_r4.filter.length) ? 4 : 5);
  }
}
function QueueJobsComponent_ng_container_0_For_20_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 12);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const option_r10 = ctx.$implicit;
    const t_r8 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("value", option_r10.field);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r8("dashboard.queues." + option_r10.field), " ");
  }
}
function QueueJobsComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275element(1, "app-document-title", 2);
    \u0275\u0275elementStart(2, "mat-drawer-container", 3)(3, "mat-drawer", 4, 0);
    \u0275\u0275repeaterCreate(5, QueueJobsComponent_ng_container_0_For_6_Template, 6, 2, "mat-card", null, _forTrack0);
    \u0275\u0275pipe(7, "async");
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(8, "mat-drawer-content")(9, "div", 5)(10, "div", 6)(11, "button", 7);
    \u0275\u0275listener("click", function QueueJobsComponent_ng_container_0_Template_button_click_11_listener() {
      \u0275\u0275restoreView(_r1);
      const drawer_r9 = \u0275\u0275reference(4);
      return \u0275\u0275resetView(drawer_r9.toggle());
    });
    \u0275\u0275elementStart(12, "mat-icon", 8);
    \u0275\u0275text(13);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(14, "div", 9)(15, "mat-form-field", 10)(16, "mat-label");
    \u0275\u0275text(17);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(18, "mat-select", 11);
    \u0275\u0275listener("valueChange", function QueueJobsComponent_ng_container_0_Template_mat_select_valueChange_18_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r4 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r4.controller.selectOrderBy($event));
    });
    \u0275\u0275repeaterCreate(19, QueueJobsComponent_ng_container_0_For_20_Template, 2, 2, "mat-option", 12, _forTrack1);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(21, "button", 13);
    \u0275\u0275listener("click", function QueueJobsComponent_ng_container_0_Template_button_click_21_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r4 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r4.controller.toggleOrderByDirection());
    });
    \u0275\u0275elementStart(22, "mat-icon");
    \u0275\u0275text(23);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementStart(24, "div", 14)(25, "button", 15);
    \u0275\u0275listener("click", function QueueJobsComponent_ng_container_0_Template_button_click_25_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r4 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r4.dataSource.refresh());
    });
    \u0275\u0275elementStart(26, "mat-icon");
    \u0275\u0275text(27, "sync");
    \u0275\u0275elementEnd()()()();
    \u0275\u0275element(28, "app-queue-jobs-table", 16);
    \u0275\u0275elementStart(29, "app-paginator", 17);
    \u0275\u0275listener("paging", function QueueJobsComponent_ng_container_0_Template_app_paginator_paging_29_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r4 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r4.controller.handlePageEvent($event));
    });
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r8 = ctx.$implicit;
    const drawer_r9 = \u0275\u0275reference(4);
    const ctx_r4 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275property("parts", \u0275\u0275pureFunction3(21, _c02, t_r8("routes.jobs"), t_r8("routes.queues"), t_r8("routes.dashboard")));
    \u0275\u0275advance(2);
    \u0275\u0275property("mode", ctx_r4.breakpoints.sizeAtLeast("Medium") ? "side" : "over")("opened", ctx_r4.breakpoints.sizeAtLeast("Medium"));
    \u0275\u0275attribute("role", ctx_r4.breakpoints.sizeAtLeast("Medium") ? "navigation" : "dialog");
    \u0275\u0275advance(2);
    \u0275\u0275repeater(\u0275\u0275pipeBind1(7, 19, ctx_r4.facets$));
    \u0275\u0275advance(6);
    \u0275\u0275property("matTooltip", t_r8("torrents.toggle_drawer"));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(drawer_r9.opened ? "arrow_circle_left" : "arrow_circle_right");
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(t_r8("torrents.order_by"));
    \u0275\u0275advance();
    \u0275\u0275property("value", ctx_r4.controls.orderBy.field);
    \u0275\u0275advance();
    \u0275\u0275repeater(ctx_r4.orderByOptions);
    \u0275\u0275advance(2);
    \u0275\u0275property("matTooltip", t_r8("torrents.order_direction_toggle"));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(ctx_r4.controls.orderBy.descending ? "arrow_downward" : "arrow_upward");
    \u0275\u0275advance(2);
    \u0275\u0275property("matTooltip", t_r8("torrents.refresh"));
    \u0275\u0275advance(3);
    \u0275\u0275property("dataSource", ctx_r4.dataSource)("selection", ctx_r4.selection);
    \u0275\u0275advance();
    \u0275\u0275property("page", ctx_r4.controls.page)("pageSize", ctx_r4.controls.limit)("pageLength", ctx_r4.dataSource.result.items.length)("totalLength", ctx_r4.dataSource.result.totalCount)("totalIsEstimate", false)("showLastPage", true);
  }
}
var QueueJobsComponent = class _QueueJobsComponent {
  constructor() {
    this.apollo = inject(Apollo);
    this.errorsService = inject(ErrorsService);
    this.breakpoints = inject(BreakpointsService);
    this.transloco = inject(TranslocoService);
    this.controller = new QueueJobsController();
    this.dataSource = new QueueJobsDatasource(this.apollo, this.errorsService, this.controller.variables$);
    this.selection = new SelectionModel();
    this.orderByOptions = orderByOptions;
    this.facets$ = this.controller.controls$.pipe(combineLatestWith(this.dataSource.result$), map(([controls, result]) => facets.map((f) => __spreadProps(__spreadValues(__spreadValues({}, f), f.extractInput(controls.facets)), {
      aggregations: f.extractAggregations(result.aggregations).map((agg) => __spreadProps(__spreadValues({}, agg), {
        label: f.resolveLabel(agg, this.transloco)
      }))
    }))));
    this.controller.controls$.subscribe((ctrl) => {
      this.controls = ctrl;
    });
  }
  static {
    this.\u0275fac = function QueueJobsComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _QueueJobsComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _QueueJobsComponent, selectors: [["app-queue-jobs"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [["drawer", ""], [4, "transloco"], [3, "parts"], [1, "drawer-container"], [1, "drawer", 3, "mode", "opened"], [1, "query-form"], [1, "form-field-container", "button-container", "button-container-toggle-drawer"], ["type", "button", "mat-icon-button", "", 1, "button-toggle-drawer", 3, "click", "matTooltip"], ["fontSet", "material-icons"], [1, "form-field-container", "form-field-container-order-by"], ["subscriptSizing", "dynamic"], [3, "valueChange", "value"], [3, "value"], ["mat-icon-button", "", 3, "click", "matTooltip"], [1, "form-field-container", "button-container", "button-container-refresh"], ["mat-mini-fab", "", "color", "primary", 3, "click", "matTooltip"], [3, "dataSource", "selection"], [3, "paging", "page", "pageSize", "pageLength", "totalLength", "totalIsEstimate", "showLastPage"], [1, "filtered"], [1, "unfiltered"], [3, "checked"], [3, "change", "checked"], ["checked", "true"], ["checked", "true", 3, "change"]], template: function QueueJobsComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, QueueJobsComponent_ng_container_0_Template, 30, 25, "ng-container", 1);
      }
    }, dependencies: [
      AppModule,
      MatOption,
      MatIconButton,
      MatMiniFabButton,
      MatCard,
      MatCardContent,
      MatCardHeader,
      MatCardTitle,
      MatCheckbox,
      MatFormField,
      MatLabel,
      MatIcon,
      MatSelect,
      MatDrawer,
      MatDrawerContainer,
      MatDrawerContent,
      MatTooltip,
      TranslocoDirective,
      AsyncPipe,
      DecimalPipe,
      PaginatorComponent,
      QueueJobsTableComponent,
      DocumentTitleComponent
    ], styles: ["\n\n.drawer[_ngcontent-%COMP%] {\n  width: 220px;\n}\n.query-form[_ngcontent-%COMP%] {\n  padding-top: 20px;\n  padding-bottom: 10px;\n  position: relative;\n  clear: both;\n  display: flex;\n  flex-wrap: wrap;\n}\n.query-form[_ngcontent-%COMP%]   .form-field-container[_ngcontent-%COMP%] {\n  display: inline-flex;\n  flex-direction: column;\n  position: relative;\n  margin-left: 20px;\n  padding-bottom: 20px;\n}\n.query-form[_ngcontent-%COMP%]   .form-field-container[_ngcontent-%COMP%]   button[_ngcontent-%COMP%] {\n  top: 8px;\n}\n.query-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-order-by[_ngcontent-%COMP%] {\n  padding-right: 40px;\n}\n.query-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-order-by[_ngcontent-%COMP%]   button[_ngcontent-%COMP%] {\n  position: absolute;\n  right: 0;\n}\n.query-form[_ngcontent-%COMP%]   .form-field-container.button-container-toggle-drawer[_ngcontent-%COMP%] {\n  margin-left: 5px;\n}\n.query-form[_ngcontent-%COMP%]   .button-container-toggle-direction[_ngcontent-%COMP%] {\n  margin-left: 4px;\n}\napp-paginator[_ngcontent-%COMP%] {\n  float: right;\n  padding-top: 14px;\n  padding-bottom: 20px;\n}\n/*# sourceMappingURL=queue-jobs.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(QueueJobsComponent, { className: "QueueJobsComponent", filePath: "src/app/dashboard/queue/queue-jobs.component.ts", lineNumber: 34 });
})();
export {
  QueueJobsComponent
};
//# sourceMappingURL=chunk-ORC7SCA4.js.map
