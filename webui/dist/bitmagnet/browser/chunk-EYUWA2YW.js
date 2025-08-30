import {
  availableQueueNames,
  statusNames
} from "./chunk-GSQBVGUV.js";
import {
  contentTypeList
} from "./chunk-UGVUNZOV.js";
import {
  ErrorsService
} from "./chunk-75G4HS47.js";
import {
  DocumentTitleComponent
} from "./chunk-OOXMQI6S.js";
import {
  Apollo,
  AppModule,
  GraphQLModule,
  MatAnchor,
  MatButton,
  MatCard,
  MatCardContent,
  MatCheckbox,
  MatDialog,
  MatDialogActions,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle,
  MatFormField,
  MatLabel,
  MatOption,
  MatProgressSpinner,
  MatSelect,
  QueueEnqueueReprocessTorrentsBatchDocument,
  QueuePurgeJobsDocument,
  TranslocoDirective
} from "./chunk-WWRDQTKJ.js";
import "./chunk-VSVMRYN2.js";
import "./chunk-Y2ZC5Z2X.js";
import {
  EMPTY,
  catchError,
  inject,
  map,
  ɵsetClassDebugInfo,
  ɵɵStandaloneFeature,
  ɵɵadvance,
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
  ɵɵproperty,
  ɵɵpureFunction3,
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

// src/app/dashboard/queue/queue-enqueue-reprocess-torrents-batch-dialog.component.ts
var _forTrack0 = ($index, $item) => $item.key;
function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_For_23_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 8);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const contentType_r4 = ctx.$implicit;
    const t_r5 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275property("value", contentType_r4.key);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r5("content_types.plural." + contentType_r4.key), " ");
  }
}
function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_Template(rf, ctx) {
  if (rf & 1) {
    const _r2 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "section")(1, "mat-checkbox", 4);
    \u0275\u0275listener("change", function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_Template_mat_checkbox_change_1_listener($event) {
      \u0275\u0275restoreView(_r2);
      const ctx_r2 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r2.purge = $event.checked);
    });
    \u0275\u0275text(2);
    \u0275\u0275elementEnd();
    \u0275\u0275element(3, "br");
    \u0275\u0275elementStart(4, "mat-checkbox", 4);
    \u0275\u0275listener("change", function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_Template_mat_checkbox_change_4_listener($event) {
      \u0275\u0275restoreView(_r2);
      const ctx_r2 = \u0275\u0275nextContext(2);
      ctx_r2.localSearchDisabled = !$event.checked;
      return \u0275\u0275resetView(ctx_r2.apisDisabled = !$event.checked ? true : ctx_r2.apisDisabled);
    });
    \u0275\u0275text(5);
    \u0275\u0275elementEnd();
    \u0275\u0275element(6, "br");
    \u0275\u0275elementStart(7, "mat-checkbox", 4);
    \u0275\u0275listener("change", function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_Template_mat_checkbox_change_7_listener($event) {
      \u0275\u0275restoreView(_r2);
      const ctx_r2 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r2.apisDisabled = !$event.checked);
    });
    \u0275\u0275text(8);
    \u0275\u0275elementEnd();
    \u0275\u0275element(9, "br");
    \u0275\u0275elementStart(10, "mat-checkbox", 4);
    \u0275\u0275listener("change", function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_Template_mat_checkbox_change_10_listener($event) {
      \u0275\u0275restoreView(_r2);
      const ctx_r2 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r2.classifierRematch = $event.checked);
    });
    \u0275\u0275text(11);
    \u0275\u0275elementEnd();
    \u0275\u0275element(12, "br");
    \u0275\u0275elementStart(13, "mat-checkbox", 4);
    \u0275\u0275listener("change", function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_Template_mat_checkbox_change_13_listener($event) {
      \u0275\u0275restoreView(_r2);
      const ctx_r2 = \u0275\u0275nextContext(2);
      ctx_r2.orphans = $event.checked;
      return \u0275\u0275resetView(ctx_r2.contentTypes = $event.checked ? ["all"] : ctx_r2.contentTypes);
    });
    \u0275\u0275text(14);
    \u0275\u0275elementEnd();
    \u0275\u0275element(15, "br");
    \u0275\u0275elementStart(16, "mat-form-field", 5)(17, "mat-label");
    \u0275\u0275text(18);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(19, "mat-select", 6);
    \u0275\u0275listener("selectionChange", function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_Template_mat_select_selectionChange_19_listener($event) {
      \u0275\u0275restoreView(_r2);
      const ctx_r2 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r2.onContentTypeSelectionChange($event));
    });
    \u0275\u0275elementStart(20, "mat-option", 7);
    \u0275\u0275text(21);
    \u0275\u0275elementEnd();
    \u0275\u0275repeaterCreate(22, QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_For_23_Template, 2, 2, "mat-option", 8, _forTrack0);
    \u0275\u0275elementEnd()()();
  }
  if (rf & 2) {
    const t_r5 = \u0275\u0275nextContext().$implicit;
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275property("checked", ctx_r2.purge);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r5("dashboard.queues.purge_queue_jobs"));
    \u0275\u0275advance(2);
    \u0275\u0275property("checked", !ctx_r2.localSearchDisabled);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r5("torrents.reprocess.match_content_by_local_search"));
    \u0275\u0275advance(2);
    \u0275\u0275property("checked", !ctx_r2.apisDisabled);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r5("torrents.reprocess.match_content_by_external_api_search"));
    \u0275\u0275advance(2);
    \u0275\u0275property("checked", ctx_r2.classifierRematch);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r5("torrents.reprocess.force_rematch"));
    \u0275\u0275advance(2);
    \u0275\u0275property("checked", ctx_r2.orphans);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r5("dashboard.queues.process_orphaned_torrents_only"));
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(t_r5("facets.content_type"));
    \u0275\u0275advance();
    \u0275\u0275property("value", ctx_r2.contentTypes);
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(t_r5("general.all"));
    \u0275\u0275advance();
    \u0275\u0275repeater(ctx_r2.allContentTypes);
  }
}
function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_6_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275element(0, "mat-spinner");
  }
}
function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_7_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "p");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r5 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r5("dashboard.queues.jobs_enqueued"));
  }
}
function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_9_Template(rf, ctx) {
  if (rf & 1) {
    const _r6 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "button", 9);
    \u0275\u0275listener("click", function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_9_Template_button_click_0_listener() {
      \u0275\u0275restoreView(_r6);
      const ctx_r2 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r2.handleEnqueue());
    });
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r5 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r5("dashboard.queues.enqueue_jobs"), " ");
  }
}
function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-card")(2, "h2", 1);
    \u0275\u0275text(3);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(4, "mat-dialog-content");
    \u0275\u0275template(5, QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_5_Template, 24, 13, "section")(6, QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_6_Template, 1, 0, "mat-spinner")(7, QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_7_Template, 2, 1, "p");
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(8, "mat-dialog-actions");
    \u0275\u0275template(9, QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Conditional_9_Template, 2, 1, "button", 2);
    \u0275\u0275elementStart(10, "button", 3);
    \u0275\u0275listener("click", function QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Template_button_click_10_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r2 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r2.dialogRef.close());
    });
    \u0275\u0275text(11);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r5 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate1(" ", t_r5("dashboard.queues.enqueue_torrent_processing_batch"), " ");
    \u0275\u0275advance(2);
    \u0275\u0275conditional(ctx_r2.stage === "PENDING" ? 5 : ctx_r2.stage === "REQUESTING" ? 6 : ctx_r2.stage === "DONE" ? 7 : -1);
    \u0275\u0275advance(4);
    \u0275\u0275conditional(ctx_r2.stage === "PENDING" ? 9 : -1);
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate1(" ", t_r5("general.dismiss"), " ");
  }
}
var QueueEnqueueReprocessTorrentsBatchDialogComponent = class _QueueEnqueueReprocessTorrentsBatchDialogComponent {
  constructor() {
    this.apollo = inject(Apollo);
    this.dialogRef = inject(MatDialogRef);
    this.errorsService = inject(ErrorsService);
    this.allContentTypes = contentTypeList;
    this.stage = "PENDING";
    this.purge = true;
    this.apisDisabled = true;
    this.localSearchDisabled = true;
    this.classifierRematch = false;
    this.contentTypes = ["all"];
    this.orphans = false;
  }
  handleEnqueue() {
    if (this.stage !== "PENDING") {
      return;
    }
    this.stage = "REQUESTING";
    this.apollo.mutate({
      mutation: QueueEnqueueReprocessTorrentsBatchDocument,
      variables: {
        input: {
          purge: this.purge,
          apisDisabled: this.apisDisabled,
          localSearchDisabled: this.localSearchDisabled,
          classifierRematch: this.classifierRematch,
          contentTypes: this.contentTypes.includes("all") ? void 0 : this.contentTypes.map((ct) => ct === "null" ? null : ct),
          orphans: this.orphans ? true : void 0
        }
      }
    }).pipe(catchError((error) => {
      this.errorsService.addError(error.message);
      this.dialogRef.close();
      return EMPTY;
    })).subscribe(() => {
      this.stage = "DONE";
      this.data.onEnqueued?.();
    });
  }
  onContentTypeSelectionChange(change) {
    if (!Array.isArray(change.value) || !change.value.length || change.value.includes("all") && (!this.contentTypes.includes("all") || change.value.length === 1)) {
      this.contentTypes = ["all"];
    } else {
      this.orphans = false;
      this.contentTypes = this.allContentTypes.map((ct) => ct.key).filter((ct) => change.value.includes(ct));
    }
  }
  static {
    this.\u0275fac = function QueueEnqueueReprocessTorrentsBatchDialogComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _QueueEnqueueReprocessTorrentsBatchDialogComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _QueueEnqueueReprocessTorrentsBatchDialogComponent, selectors: [["app-queue-enqueue-reprocess-torrents-batch-dialog"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], ["mat-dialog-title", ""], ["mat-stroked-button", "", "color", "warning"], ["mat-stroked-button", "", 3, "click"], [3, "change", "checked"], [1, "select-content-types"], ["multiple", "", 3, "selectionChange", "value"], ["value", "all"], [3, "value"], ["mat-stroked-button", "", "color", "warning", 3, "click"]], template: function QueueEnqueueReprocessTorrentsBatchDialogComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, QueueEnqueueReprocessTorrentsBatchDialogComponent_ng_container_0_Template, 12, 4, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatOption, MatButton, MatCard, MatCheckbox, MatDialogTitle, MatDialogActions, MatDialogContent, MatFormField, MatLabel, MatProgressSpinner, MatSelect, TranslocoDirective], styles: ["\n\nmat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%] {\n  width: 100%;\n  height: 100%;\n}\n.select-content-types[_ngcontent-%COMP%] {\n  margin-top: 10px;\n}\n/*# sourceMappingURL=queue-enqueue-reprocess-torrents-batch-dialog.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(QueueEnqueueReprocessTorrentsBatchDialogComponent, { className: "QueueEnqueueReprocessTorrentsBatchDialogComponent", filePath: "src/app/dashboard/queue/queue-enqueue-reprocess-torrents-batch-dialog.component.ts", lineNumber: 19 });
})();

// src/app/dashboard/queue/queue-purge-jobs-dialog.component.ts
function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_For_6_Template(rf, ctx) {
  if (rf & 1) {
    const _r4 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-checkbox", 6);
    \u0275\u0275listener("change", function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_For_6_Template_mat_checkbox_change_0_listener($event) {
      \u0275\u0275restoreView(_r4);
      const ctx_r2 = \u0275\u0275nextContext(3);
      return \u0275\u0275resetView(ctx_r2.handleQueueEvent($event));
    });
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const queue_r5 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext(3);
    \u0275\u0275property("value", queue_r5)("checked", ctx_r2.queues == null ? null : ctx_r2.queues.includes(queue_r5));
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(queue_r5);
  }
}
function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_For_13_Template(rf, ctx) {
  if (rf & 1) {
    const _r6 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-checkbox", 6);
    \u0275\u0275listener("change", function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_For_13_Template_mat_checkbox_change_0_listener($event) {
      \u0275\u0275restoreView(_r6);
      const ctx_r2 = \u0275\u0275nextContext(3);
      return \u0275\u0275resetView(ctx_r2.handleStatusEvent($event));
    });
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const status_r7 = ctx.$implicit;
    const t_r8 = \u0275\u0275nextContext(2).$implicit;
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275property("value", status_r7)("checked", ctx_r2.statuses == null ? null : ctx_r2.statuses.includes(status_r7));
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r8("dashboard.queues." + status_r7));
  }
}
function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_Template(rf, ctx) {
  if (rf & 1) {
    const _r2 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "section")(1, "h4");
    \u0275\u0275text(2);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(3, "mat-checkbox", 4);
    \u0275\u0275listener("change", function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_Template_mat_checkbox_change_3_listener($event) {
      \u0275\u0275restoreView(_r2);
      const ctx_r2 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r2.handleQueueEvent($event));
    });
    \u0275\u0275text(4);
    \u0275\u0275elementEnd();
    \u0275\u0275repeaterCreate(5, QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_For_6_Template, 2, 3, "mat-checkbox", 5, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(7, "section")(8, "h4");
    \u0275\u0275text(9);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(10, "mat-checkbox", 4);
    \u0275\u0275listener("change", function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_Template_mat_checkbox_change_10_listener($event) {
      \u0275\u0275restoreView(_r2);
      const ctx_r2 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r2.handleStatusEvent($event));
    });
    \u0275\u0275text(11);
    \u0275\u0275elementEnd();
    \u0275\u0275repeaterCreate(12, QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_For_13_Template, 2, 3, "mat-checkbox", 5, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r8 = \u0275\u0275nextContext().$implicit;
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate1("", t_r8("dashboard.queues.queues"), ":");
    \u0275\u0275advance();
    \u0275\u0275property("checked", ctx_r2.queues === void 0);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r8("general.all"));
    \u0275\u0275advance();
    \u0275\u0275repeater(ctx_r2.availableQueueNames);
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate1("", t_r8("general.status"), ":");
    \u0275\u0275advance();
    \u0275\u0275property("checked", ctx_r2.statuses === void 0);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r8("general.all"));
    \u0275\u0275advance();
    \u0275\u0275repeater(ctx_r2.statusNames);
  }
}
function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_6_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275element(0, "mat-spinner");
  }
}
function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_7_Conditional_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "p");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r8 = \u0275\u0275nextContext(2).$implicit;
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275textInterpolate2("", t_r8("general.error"), ": ", ctx_r2.error.message, "");
  }
}
function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_7_Conditional_1_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "p");
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r8 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r8("dashboard.queues.queue_purged"));
  }
}
function QueuePurgeJobsDialogComponent_ng_container_0_Conditional_7_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275template(0, QueuePurgeJobsDialogComponent_ng_container_0_Conditional_7_Conditional_0_Template, 2, 2, "p")(1, QueuePurgeJobsDialogComponent_ng_container_0_Conditional_7_Conditional_1_Template, 2, 1, "p");
  }
  if (rf & 2) {
    const ctx_r2 = \u0275\u0275nextContext(2);
    \u0275\u0275conditional(ctx_r2.error ? 0 : 1);
  }
}
function QueuePurgeJobsDialogComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-card")(2, "h2", 1);
    \u0275\u0275text(3);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(4, "mat-dialog-content");
    \u0275\u0275template(5, QueuePurgeJobsDialogComponent_ng_container_0_Conditional_5_Template, 14, 6)(6, QueuePurgeJobsDialogComponent_ng_container_0_Conditional_6_Template, 1, 0, "mat-spinner")(7, QueuePurgeJobsDialogComponent_ng_container_0_Conditional_7_Template, 2, 1);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(8, "mat-dialog-actions")(9, "button", 2);
    \u0275\u0275listener("click", function QueuePurgeJobsDialogComponent_ng_container_0_Template_button_click_9_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r2 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r2.handlePurgeJobs());
    });
    \u0275\u0275text(10);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(11, "button", 3);
    \u0275\u0275listener("click", function QueuePurgeJobsDialogComponent_ng_container_0_Template_button_click_11_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r2 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r2.dialogRef.close());
    });
    \u0275\u0275text(12);
    \u0275\u0275elementEnd()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r8 = ctx.$implicit;
    const ctx_r2 = \u0275\u0275nextContext();
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate(t_r8("dashboard.queues.purge_queue_jobs"));
    \u0275\u0275advance(2);
    \u0275\u0275conditional(ctx_r2.stage === "PENDING" ? 5 : ctx_r2.stage === "REQUESTING" ? 6 : ctx_r2.stage === "DONE" ? 7 : -1);
    \u0275\u0275advance(4);
    \u0275\u0275property("disabled", ctx_r2.stage !== "PENDING");
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r8("dashboard.queues.purge_jobs"), " ");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate1(" ", t_r8("general.dismiss"), " ");
  }
}
var QueuePurgeJobsDialogComponent = class _QueuePurgeJobsDialogComponent {
  constructor() {
    this.apollo = inject(Apollo);
    this.dialogRef = inject(MatDialogRef);
    this.availableQueueNames = availableQueueNames;
    this.statusNames = statusNames;
    this.stage = "PENDING";
  }
  handleQueueEvent(event) {
    if (event.source.value === "_all") {
      this.queues = void 0;
      return;
    }
    if (event.checked) {
      let queues = this.queues ?? [];
      if (!queues.includes(event.source.value)) {
        queues = [...queues, event.source.value];
      }
      if (queues.length === this.availableQueueNames.length) {
        event.source.checked = false;
        this.queues = void 0;
      } else {
        this.queues = queues;
      }
    } else {
      const queues = this.queues?.filter((q) => q !== event.source.value);
      if (!queues?.length) {
        this.queues = void 0;
      } else {
        this.queues = queues;
      }
    }
  }
  handleStatusEvent(event) {
    if (event.source.value === "_all") {
      this.statuses = void 0;
      return;
    }
    if (event.checked) {
      let statuses = this.statuses ?? [];
      if (!statuses.includes(event.source.value)) {
        statuses = [
          ...statuses,
          event.source.value
        ];
      }
      if (statuses.length === this.statusNames.length) {
        event.source.checked = false;
        this.statuses = void 0;
      } else {
        this.statuses = statuses;
      }
    } else {
      const statuses = this.statuses?.filter((s) => s !== event.source.value);
      if (!statuses?.length) {
        this.statuses = void 0;
      } else {
        this.statuses = statuses;
      }
    }
  }
  handlePurgeJobs() {
    if (this.stage !== "PENDING") {
      return;
    }
    this.stage = "REQUESTING";
    this.apollo.mutate({
      mutation: QueuePurgeJobsDocument,
      variables: {
        input: {
          queues: this.queues,
          statuses: this.statuses
        }
      }
    }).pipe(catchError((err) => {
      this.stage = "DONE";
      this.error = err;
      return EMPTY;
    }), map(() => {
      this.stage = "DONE";
      this.data?.onPurged?.();
    })).subscribe();
  }
  static {
    this.\u0275fac = function QueuePurgeJobsDialogComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _QueuePurgeJobsDialogComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _QueuePurgeJobsDialogComponent, selectors: [["app-queue-purge-jobs-dialog"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], ["mat-dialog-title", ""], ["mat-stroked-button", "", "color", "warning", 3, "click", "disabled"], ["mat-stroked-button", "", 3, "click"], ["value", "_all", 3, "change", "checked"], [3, "value", "checked"], [3, "change", "value", "checked"]], template: function QueuePurgeJobsDialogComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, QueuePurgeJobsDialogComponent_ng_container_0_Template, 13, 5, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatButton, MatCard, MatCheckbox, MatDialogTitle, MatDialogActions, MatDialogContent, MatProgressSpinner, TranslocoDirective, GraphQLModule], styles: ["\n\nmat-dialog-content[_ngcontent-%COMP%] {\n  min-height: 240px;\n  overflow: visible;\n}\nmat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%] {\n  width: 100%;\n  height: 100%;\n}\n/*# sourceMappingURL=queue-purge-jobs-dialog.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(QueuePurgeJobsDialogComponent, { className: "QueuePurgeJobsDialogComponent", filePath: "src/app/dashboard/queue/queue-purge-jobs-dialog.component.ts", lineNumber: 19 });
})();

// src/app/dashboard/queue/queue-admin.component.ts
var _c0 = (a0, a1, a2) => [a0, a1, a2];
function QueueAdminComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275element(1, "app-document-title", 1);
    \u0275\u0275elementStart(2, "mat-card")(3, "mat-card-content")(4, "ul")(5, "li")(6, "a", 2);
    \u0275\u0275listener("click", function QueueAdminComponent_ng_container_0_Template_a_click_6_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.openDialogPurgeJobs());
    });
    \u0275\u0275text(7);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(8, "li")(9, "a", 2);
    \u0275\u0275listener("click", function QueueAdminComponent_ng_container_0_Template_a_click_9_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.openDialogEnqueueReprocessTorrentsBatch());
    });
    \u0275\u0275text(10);
    \u0275\u0275elementEnd()()()()();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r3 = ctx.$implicit;
    \u0275\u0275advance();
    \u0275\u0275property("parts", \u0275\u0275pureFunction3(3, _c0, t_r3("routes.admin"), t_r3("routes.queues"), t_r3("routes.dashboard")));
    \u0275\u0275advance(6);
    \u0275\u0275textInterpolate(t_r3("dashboard.queues.purge_queue_jobs"));
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate(t_r3("dashboard.queues.enqueue_torrent_processing_batch"));
  }
}
var QueueAdminComponent = class _QueueAdminComponent {
  constructor() {
    this.dialog = inject(MatDialog);
  }
  openDialogPurgeJobs() {
    this.dialog.open(QueuePurgeJobsDialogComponent);
  }
  openDialogEnqueueReprocessTorrentsBatch() {
    this.dialog.open(QueueEnqueueReprocessTorrentsBatchDialogComponent);
  }
  static {
    this.\u0275fac = function QueueAdminComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _QueueAdminComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _QueueAdminComponent, selectors: [["app-queue-admin"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], [3, "parts"], ["mat-button", "", 3, "click"]], template: function QueueAdminComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, QueueAdminComponent_ng_container_0_Template, 11, 7, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatAnchor, MatCard, MatCardContent, TranslocoDirective, DocumentTitleComponent], styles: ["\n\nul[_ngcontent-%COMP%] {\n  list-style-type: none;\n  padding-left: 0;\n}\nul[_ngcontent-%COMP%]   li[_ngcontent-%COMP%] {\n  margin-bottom: 6px;\n}\n/*# sourceMappingURL=queue-admin.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(QueueAdminComponent, { className: "QueueAdminComponent", filePath: "src/app/dashboard/queue/queue-admin.component.ts", lineNumber: 15 });
})();
export {
  QueueAdminComponent
};
//# sourceMappingURL=chunk-EYUWA2YW.js.map
