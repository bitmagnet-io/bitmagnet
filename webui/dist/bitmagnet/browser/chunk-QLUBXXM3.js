import {
  FilesizePipe,
  TorrentChipsComponent,
  TorrentContentComponent,
  TorrentReprocessComponent
} from "./chunk-LUZJBAO3.js";
import {
  TimeAgoPipe
} from "./chunk-43HRGFU3.js";
import {
  contentTypeInfo
} from "./chunk-UGVUNZOV.js";
import {
  ErrorsService
} from "./chunk-75G4HS47.js";
import {
  BreakpointsService
} from "./chunk-NQ6E5D5R.js";
import {
  AppModule,
  COMMA,
  CdkCopyToClipboard,
  DefaultValueAccessor,
  ENTER,
  FormControl,
  FormControlDirective,
  GraphQLService,
  MatAutocomplete,
  MatAutocompleteTrigger,
  MatButton,
  MatCard,
  MatCardActions,
  MatCardContent,
  MatCell,
  MatCellDef,
  MatCheckbox,
  MatChipGrid,
  MatChipInput,
  MatChipRemove,
  MatChipRow,
  MatColumnDef,
  MatFormField,
  MatHeaderCell,
  MatHeaderCellDef,
  MatHeaderRow,
  MatHeaderRowDef,
  MatIcon,
  MatOption,
  MatProgressBar,
  MatRow,
  MatRowDef,
  MatTab,
  MatTabContent,
  MatTabGroup,
  MatTabLabel,
  MatTable,
  MatTooltip,
  NgControlStatus,
  TranslocoDirective
} from "./chunk-WWRDQTKJ.js";
import {
  animate,
  state,
  style,
  transition,
  trigger
} from "./chunk-VSVMRYN2.js";
import {
  ActivatedRoute,
  Router
} from "./chunk-Y2ZC5Z2X.js";
import {
  AsyncPipe,
  EMPTY,
  EventEmitter,
  Observable,
  __spreadProps,
  __spreadValues,
  catchError,
  inject,
  tap,
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
  ɵɵpipeBind2,
  ɵɵproperty,
  ɵɵpropertyInterpolate,
  ɵɵpureFunction0,
  ɵɵreference,
  ɵɵrepeater,
  ɵɵrepeaterCreate,
  ɵɵrepeaterTrackByIdentity,
  ɵɵresetView,
  ɵɵrestoreView,
  ɵɵsanitizeUrl,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate,
  ɵɵtextInterpolate1,
  ɵɵtextInterpolate2
} from "./chunk-DMMUMX3A.js";

// src/app/torrents/torrents-bulk-actions.component.ts
function TorrentsBulkActionsComponent_ng_container_0_ng_template_4_Conditional_2_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "span", 7);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("torrents.copy"));
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_4_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-icon");
    \u0275\u0275text(1, "content_copy");
    \u0275\u0275elementEnd();
    \u0275\u0275template(2, TorrentsBulkActionsComponent_ng_container_0_ng_template_4_Conditional_2_Template, 2, 1, "span", 7);
  }
  if (rf & 2) {
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance(2);
    \u0275\u0275conditional(ctx_r1.breakpoints.sizeAtLeast("Medium") ? 2 : -1);
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_5_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-card")(1, "mat-card-actions", 8)(2, "button", 9);
    \u0275\u0275element(3, "mat-icon", 10);
    \u0275\u0275text(4);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(5, "button", 9)(6, "mat-icon");
    \u0275\u0275text(7, "tag");
    \u0275\u0275elementEnd();
    \u0275\u0275text(8);
    \u0275\u0275elementEnd()()();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance(2);
    \u0275\u0275property("disabled", !ctx_r1.selectedItems.length)("matTooltip", t_r3("torrents.copy_to_clipboard"))("cdkCopyToClipboard", ctx_r1.getSelectedMagnetLinks());
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate1("", t_r3("torrents.magnet_links"), " ");
    \u0275\u0275advance();
    \u0275\u0275property("disabled", !ctx_r1.selectedItems.length)("matTooltip", t_r3("torrents.copy_to_clipboard"))("cdkCopyToClipboard", ctx_r1.getSelectedInfoHashesLines());
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate1("", t_r3("torrents.info_hashes"), " ");
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_7_Conditional_2_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "span", 7);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("torrents.edit_tags"));
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_7_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-icon");
    \u0275\u0275text(1, "sell");
    \u0275\u0275elementEnd();
    \u0275\u0275template(2, TorrentsBulkActionsComponent_ng_container_0_ng_template_7_Conditional_2_Template, 2, 1, "span", 7);
  }
  if (rf & 2) {
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance(2);
    \u0275\u0275conditional(ctx_r1.breakpoints.sizeAtLeast("Medium") ? 2 : -1);
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_For_5_Template(rf, ctx) {
  if (rf & 1) {
    const _r5 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-chip-row", 18);
    \u0275\u0275listener("edited", function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_For_5_Template_mat_chip_row_edited_0_listener($event) {
      const tagName_r6 = \u0275\u0275restoreView(_r5).$implicit;
      const ctx_r1 = \u0275\u0275nextContext(3);
      return \u0275\u0275resetView(ctx_r1.renameTag(tagName_r6, $event.value));
    })("removed", function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_For_5_Template_mat_chip_row_removed_0_listener() {
      const tagName_r6 = \u0275\u0275restoreView(_r5).$implicit;
      const ctx_r1 = \u0275\u0275nextContext(3);
      return \u0275\u0275resetView(ctx_r1.deleteTag(tagName_r6));
    });
    \u0275\u0275text(1);
    \u0275\u0275elementStart(2, "mat-icon", 19);
    \u0275\u0275text(3, "cancel");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const tagName_r6 = ctx.$implicit;
    \u0275\u0275property("editable", true)("aria-description", "press enter to edit");
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", tagName_r6, " ");
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_For_10_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-option", 16);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const tagName_r7 = ctx.$implicit;
    \u0275\u0275property("value", tagName_r7);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(tagName_r7);
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_Template(rf, ctx) {
  if (rf & 1) {
    const _r4 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-card")(1, "mat-form-field", 11)(2, "mat-chip-grid", 12, 0);
    \u0275\u0275repeaterCreate(4, TorrentsBulkActionsComponent_ng_container_0_ng_template_8_For_5_Template, 4, 3, "mat-chip-row", 13, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(6, "input", 14);
    \u0275\u0275listener("matChipInputTokenEnd", function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_Template_input_matChipInputTokenEnd_6_listener($event) {
      \u0275\u0275restoreView(_r4);
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView($event.value && ctx_r1.addTag($event.value));
    });
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(7, "mat-autocomplete", 15, 1);
    \u0275\u0275listener("optionSelected", function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_Template_mat_autocomplete_optionSelected_7_listener($event) {
      \u0275\u0275restoreView(_r4);
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.addTag($event.option.viewValue));
    });
    \u0275\u0275repeaterCreate(9, TorrentsBulkActionsComponent_ng_container_0_ng_template_8_For_10_Template, 2, 2, "mat-option", 16, \u0275\u0275repeaterTrackByIdentity);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(11, "mat-card-actions", 8)(12, "button", 17);
    \u0275\u0275listener("click", function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_Template_button_click_12_listener() {
      \u0275\u0275restoreView(_r4);
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.setTags());
    });
    \u0275\u0275text(13);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(14, "button", 17);
    \u0275\u0275listener("click", function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_Template_button_click_14_listener() {
      \u0275\u0275restoreView(_r4);
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.putTags());
    });
    \u0275\u0275text(15);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(16, "button", 17);
    \u0275\u0275listener("click", function TorrentsBulkActionsComponent_ng_container_0_ng_template_8_Template_button_click_16_listener() {
      \u0275\u0275restoreView(_r4);
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.deleteTags());
    });
    \u0275\u0275text(17);
    \u0275\u0275elementEnd()()();
  }
  if (rf & 2) {
    const chipGrid_r8 = \u0275\u0275reference(3);
    const auto_r9 = \u0275\u0275reference(8);
    const t_r3 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance(4);
    \u0275\u0275repeater(ctx_r1.editedTags);
    \u0275\u0275advance(2);
    \u0275\u0275propertyInterpolate("placeholder", t_r3("torrents.tags.placeholder"));
    \u0275\u0275property("formControl", ctx_r1.newTagCtrl)("matAutocomplete", auto_r9)("matChipInputFor", chipGrid_r8)("matChipInputSeparatorKeyCodes", ctx_r1.separatorKeysCodes)("value", ctx_r1.newTagCtrl.value);
    \u0275\u0275advance(3);
    \u0275\u0275repeater(ctx_r1.suggestedTags);
    \u0275\u0275advance(3);
    \u0275\u0275propertyInterpolate("matTooltip", t_r3("torrents.tags.set_tip"));
    \u0275\u0275property("disabled", !ctx_r1.selectedItems.length);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r3("torrents.tags.set"), " ");
    \u0275\u0275advance();
    \u0275\u0275propertyInterpolate("matTooltip", t_r3("torrents.tags.put_tip"));
    \u0275\u0275property("disabled", !ctx_r1.selectedItems.length || !ctx_r1.editedTags.length && !ctx_r1.newTagCtrl.value);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r3("torrents.tags.put"), " ");
    \u0275\u0275advance();
    \u0275\u0275propertyInterpolate("matTooltip", t_r3("torrents.tags.delete_tip"));
    \u0275\u0275property("disabled", !ctx_r1.selectedItems.length || !ctx_r1.editedTags.length && !ctx_r1.newTagCtrl.value);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r3("torrents.tags.delete"), " ");
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_10_Conditional_2_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "span", 7);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("torrents.classification"));
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_10_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-icon");
    \u0275\u0275text(1, "category");
    \u0275\u0275elementEnd();
    \u0275\u0275template(2, TorrentsBulkActionsComponent_ng_container_0_ng_template_10_Conditional_2_Template, 2, 1, "span", 7);
  }
  if (rf & 2) {
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance(2);
    \u0275\u0275conditional(ctx_r1.breakpoints.sizeAtLeast("Medium") ? 2 : -1);
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_11_Template(rf, ctx) {
  if (rf & 1) {
    const _r10 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "app-torrent-reprocess", 20);
    \u0275\u0275listener("updated", function TorrentsBulkActionsComponent_ng_container_0_ng_template_11_Template_app_torrent_reprocess_updated_0_listener() {
      \u0275\u0275restoreView(_r10);
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.updated.emit(null));
    });
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275property("infoHashes", ctx_r1.selectedInfoHashes);
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_13_Conditional_2_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "span", 7);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext(2).$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("torrents.delete"));
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_13_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-icon");
    \u0275\u0275text(1, "delete_forever");
    \u0275\u0275elementEnd();
    \u0275\u0275template(2, TorrentsBulkActionsComponent_ng_container_0_ng_template_13_Conditional_2_Template, 2, 1, "span", 7);
  }
  if (rf & 2) {
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance(2);
    \u0275\u0275conditional(ctx_r1.breakpoints.sizeAtLeast("Medium") ? 2 : -1);
  }
}
function TorrentsBulkActionsComponent_ng_container_0_ng_template_14_Template(rf, ctx) {
  if (rf & 1) {
    const _r11 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "mat-card")(1, "mat-card-content")(2, "p")(3, "strong");
    \u0275\u0275text(4);
    \u0275\u0275elementEnd();
    \u0275\u0275element(5, "br");
    \u0275\u0275text(6);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(7, "mat-card-actions", 8)(8, "button", 21);
    \u0275\u0275listener("click", function TorrentsBulkActionsComponent_ng_container_0_ng_template_14_Template_button_click_8_listener() {
      \u0275\u0275restoreView(_r11);
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.deleteTorrents());
    });
    \u0275\u0275elementStart(9, "mat-icon");
    \u0275\u0275text(10, "delete_forever");
    \u0275\u0275elementEnd();
    \u0275\u0275text(11);
    \u0275\u0275elementEnd()()();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate(t_r3("torrents.delete_are_you_sure"));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate1("", t_r3("torrents.delete_action_cannot_be_undone"), ". ");
    \u0275\u0275advance(2);
    \u0275\u0275property("disabled", !ctx_r1.selectedItems.length);
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate1("", t_r3("torrents.delete"), " ");
  }
}
function TorrentsBulkActionsComponent_ng_container_0_Conditional_15_ng_template_1_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-icon", 22);
    \u0275\u0275text(1, "close");
    \u0275\u0275elementEnd();
  }
}
function TorrentsBulkActionsComponent_ng_container_0_Conditional_15_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "mat-tab");
    \u0275\u0275template(1, TorrentsBulkActionsComponent_ng_container_0_Conditional_15_ng_template_1_Template, 2, 0, "ng-template", 5);
    \u0275\u0275elementEnd();
  }
}
function TorrentsBulkActionsComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-tab-group", 3);
    \u0275\u0275listener("focusChange", function TorrentsBulkActionsComponent_ng_container_0_Template_mat_tab_group_focusChange_1_listener($event) {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext();
      return \u0275\u0275resetView(ctx_r1.selectTab($event.index == 5 ? 0 : $event.index));
    });
    \u0275\u0275element(2, "mat-tab", 4);
    \u0275\u0275elementStart(3, "mat-tab");
    \u0275\u0275template(4, TorrentsBulkActionsComponent_ng_container_0_ng_template_4_Template, 3, 1, "ng-template", 5)(5, TorrentsBulkActionsComponent_ng_container_0_ng_template_5_Template, 9, 8, "ng-template", 6);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(6, "mat-tab");
    \u0275\u0275template(7, TorrentsBulkActionsComponent_ng_container_0_ng_template_7_Template, 3, 1, "ng-template", 5)(8, TorrentsBulkActionsComponent_ng_container_0_ng_template_8_Template, 18, 15, "ng-template", 6);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(9, "mat-tab");
    \u0275\u0275template(10, TorrentsBulkActionsComponent_ng_container_0_ng_template_10_Template, 3, 1, "ng-template", 5)(11, TorrentsBulkActionsComponent_ng_container_0_ng_template_11_Template, 1, 1, "ng-template", 6);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(12, "mat-tab");
    \u0275\u0275template(13, TorrentsBulkActionsComponent_ng_container_0_ng_template_13_Template, 3, 1, "ng-template", 5)(14, TorrentsBulkActionsComponent_ng_container_0_ng_template_14_Template, 12, 4, "ng-template", 6);
    \u0275\u0275elementEnd();
    \u0275\u0275template(15, TorrentsBulkActionsComponent_ng_container_0_Conditional_15_Template, 2, 0, "mat-tab");
    \u0275\u0275elementEnd();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275property("selectedIndex", ctx_r1.selectedTabIndex)("mat-stretch-tabs", false);
    \u0275\u0275advance();
    \u0275\u0275property("aria-labelledby", "hidden");
    \u0275\u0275advance(13);
    \u0275\u0275conditional(ctx_r1.selectedTabIndex > 0 ? 15 : -1);
  }
}
var TorrentsBulkActionsComponent = class _TorrentsBulkActionsComponent {
  constructor() {
    this.graphQLService = inject(GraphQLService);
    this.errorsService = inject(ErrorsService);
    this.breakpoints = inject(BreakpointsService);
    this.selectedItems$ = new Observable();
    this.updated = new EventEmitter();
    this.separatorKeysCodes = [ENTER, COMMA];
    this.selectedTabIndex = 0;
    this.newTagCtrl = new FormControl("");
    this.editedTags = Array();
    this.suggestedTags = Array();
    this.selectedItems = new Array();
    this.selectedInfoHashes = new Array();
  }
  ngOnInit() {
    this.selectedItems$.subscribe((items) => {
      this.selectedItems = items;
      this.selectedInfoHashes = items.map((i) => i.infoHash);
    });
    this.newTagCtrl.reset();
  }
  selectTab(index) {
    this.selectedTabIndex = index;
  }
  getSelectedMagnetLinks() {
    return this.selectedItems.map((i) => i.torrent.magnetUri).join("\n");
  }
  getSelectedInfoHashesLines() {
    return this.selectedInfoHashes.join("\n");
  }
  addTag(tagName) {
    if (!this.editedTags.includes(tagName)) {
      this.editedTags.push(tagName);
    }
    this.newTagCtrl.reset();
    this.updateSuggestedTags();
  }
  deleteTag(tagName) {
    this.editedTags = this.editedTags.filter((t) => t !== tagName);
    this.updateSuggestedTags();
  }
  renameTag(fromTagName, toTagName) {
    this.editedTags = this.editedTags.map((t) => t === fromTagName ? toTagName : t);
    this.updateSuggestedTags();
  }
  putTags() {
    const infoHashes = this.selectedItems.map(({ infoHash }) => infoHash);
    if (!infoHashes.length) {
      return;
    }
    if (this.newTagCtrl.value) {
      this.addTag(this.newTagCtrl.value);
    }
    return this.graphQLService.torrentPutTags({
      infoHashes,
      tagNames: this.editedTags
    }).pipe(catchError((err) => {
      this.errorsService.addError(`Error putting tags: ${err.message}`);
      return EMPTY;
    })).pipe(tap(() => {
      this.updated.emit();
    })).subscribe();
  }
  setTags() {
    const infoHashes = this.selectedItems.map(({ infoHash }) => infoHash);
    if (!infoHashes.length) {
      return;
    }
    if (this.newTagCtrl.value) {
      this.addTag(this.newTagCtrl.value);
    }
    return this.graphQLService.torrentSetTags({
      infoHashes,
      tagNames: this.editedTags
    }).pipe(catchError((err) => {
      this.errorsService.addError(`Error setting tags: ${err.message}`);
      return EMPTY;
    })).pipe(tap(() => {
      this.updated.emit();
    })).subscribe();
  }
  deleteTags() {
    const infoHashes = this.selectedItems.map(({ infoHash }) => infoHash);
    if (!infoHashes.length) {
      return;
    }
    if (this.newTagCtrl.value) {
      this.addTag(this.newTagCtrl.value);
    }
    return this.graphQLService.torrentDeleteTags({
      infoHashes,
      tagNames: this.editedTags
    }).pipe(catchError((err) => {
      this.errorsService.addError(`Error deleting tags: ${err.message}`);
      return EMPTY;
    })).pipe(tap(() => {
      this.updated.emit();
    })).subscribe();
  }
  updateSuggestedTags() {
    return this.graphQLService.torrentSuggestTags({
      input: {
        prefix: this.newTagCtrl.value,
        exclusions: this.editedTags
      }
    }).pipe(tap((result) => {
      this.suggestedTags.splice(0, this.suggestedTags.length, ...result.suggestions.map((t) => t.name));
    })).subscribe();
  }
  deleteTorrents() {
    const infoHashes = this.selectedItems.map(({ infoHash }) => infoHash);
    this.graphQLService.torrentDelete({ infoHashes }).pipe(catchError((err) => {
      this.errorsService.addError(`Error deleting torrents: ${err.message}`);
      return EMPTY;
    })).pipe(tap(() => {
      this.updated.emit();
    })).subscribe();
  }
  static {
    this.\u0275fac = function TorrentsBulkActionsComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _TorrentsBulkActionsComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _TorrentsBulkActionsComponent, selectors: [["app-torrents-bulk-actions"]], inputs: { selectedItems$: "selectedItems$" }, outputs: { updated: "updated" }, standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [["chipGrid", ""], ["auto", "matAutocomplete"], [4, "transloco"], ["animationDuration", "0", 1, "tab-group-bulk-actions", 3, "focusChange", "selectedIndex", "mat-stretch-tabs"], [1, "bulk-tab-placeholder", 3, "aria-labelledby"], ["mat-tab-label", ""], ["matTabContent", ""], [1, "label"], [1, "button-row"], ["mat-stroked-button", "", 3, "disabled", "matTooltip", "cdkCopyToClipboard"], ["svgIcon", "magnet"], ["subscriptSizing", "dynamic", 1, "form-edit-tags"], ["aria-label", "Enter tags"], [3, "editable", "aria-description"], [3, "matChipInputTokenEnd", "placeholder", "formControl", "matAutocomplete", "matChipInputFor", "matChipInputSeparatorKeyCodes", "value"], [3, "optionSelected"], [3, "value"], ["mat-stroked-button", "", "color", "primary", 3, "click", "disabled", "matTooltip"], [3, "edited", "removed", "editable", "aria-description"], ["matChipRemove", ""], [3, "updated", "infoHashes"], ["mat-stroked-button", "", "color", "warning", 3, "click", "disabled"], [2, "margin-right", "0"]], template: function TorrentsBulkActionsComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, TorrentsBulkActionsComponent_ng_container_0_Template, 16, 4, "ng-container", 2);
      }
    }, dependencies: [AppModule, CdkCopyToClipboard, MatAutocomplete, MatOption, MatAutocompleteTrigger, MatButton, MatCard, MatCardActions, MatCardContent, MatChipGrid, MatChipInput, MatChipRemove, MatChipRow, MatFormField, MatIcon, MatTabContent, MatTabLabel, MatTab, MatTabGroup, MatTooltip, DefaultValueAccessor, NgControlStatus, FormControlDirective, TranslocoDirective, TorrentReprocessComponent], styles: ["\n\nmat-tab-group[_ngcontent-%COMP%] {\n  padding-left: 10px;\n}\n.mat-mdc-card[_ngcontent-%COMP%] {\n  margin-bottom: 10px;\n}\nbutton[_ngcontent-%COMP%] {\n  margin-right: 10px;\n}\np[_ngcontent-%COMP%] {\n  margin-top: 0;\n}\n  .mdc-tab[aria-labelledby=hidden] {\n  display: none;\n}\n/*# sourceMappingURL=torrents-bulk-actions.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(TorrentsBulkActionsComponent, { className: "TorrentsBulkActionsComponent", filePath: "src/app/torrents/torrents-bulk-actions.component.ts", lineNumber: 26 });
})();

// src/app/torrents/torrents-table.component.ts
var _c0 = () => ["expandedDetail"];
function TorrentsTableComponent_ng_container_0_th_6_Template(rf, ctx) {
  if (rf & 1) {
    const _r1 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "th", 19)(1, "mat-checkbox", 20);
    \u0275\u0275listener("change", function TorrentsTableComponent_ng_container_0_th_6_Template_mat_checkbox_change_1_listener() {
      \u0275\u0275restoreView(_r1);
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.toggleAllRows());
    });
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275property("checked", ctx_r1.multiSelection.hasValue() && ctx_r1.isAllSelected())("indeterminate", ctx_r1.multiSelection.hasValue() && !ctx_r1.isAllSelected())("matTooltip", ctx_r1.isAllSelected() ? t_r3("torrents.deselect_all") : t_r3("torrents.select_all"));
  }
}
function TorrentsTableComponent_ng_container_0_td_7_Template(rf, ctx) {
  if (rf & 1) {
    const _r4 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 21)(1, "mat-checkbox", 22);
    \u0275\u0275listener("click", function TorrentsTableComponent_ng_container_0_td_7_Template_mat_checkbox_click_1_listener($event) {
      \u0275\u0275restoreView(_r4);
      return \u0275\u0275resetView($event.stopPropagation());
    })("change", function TorrentsTableComponent_ng_container_0_td_7_Template_mat_checkbox_change_1_listener($event) {
      const i_r5 = \u0275\u0275restoreView(_r4).$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView($event ? ctx_r1.multiSelection.toggle(ctx_r1.item(i_r5).infoHash) : null);
    });
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const i_r5 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275property("checked", ctx_r1.multiSelection.isSelected(ctx_r1.item(i_r5).infoHash));
  }
}
function TorrentsTableComponent_ng_container_0_th_9_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 19);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("torrents.summary"));
  }
}
function TorrentsTableComponent_ng_container_0_td_10_Conditional_5_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "p", 26);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const i_r7 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(ctx_r1.item(i_r7).torrent.name);
  }
}
function TorrentsTableComponent_ng_container_0_td_10_Template(rf, ctx) {
  if (rf & 1) {
    const _r6 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 23);
    \u0275\u0275listener("click", function TorrentsTableComponent_ng_container_0_td_10_Template_td_click_0_listener($event) {
      const i_r7 = \u0275\u0275restoreView(_r6).$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      ctx_r1.toggleSelectedTorrent(i_r7.infoHash);
      return \u0275\u0275resetView($event.stopPropagation());
    });
    \u0275\u0275elementStart(1, "mat-icon", 24);
    \u0275\u0275text(2);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(3, "span", 25);
    \u0275\u0275text(4);
    \u0275\u0275elementEnd();
    \u0275\u0275template(5, TorrentsTableComponent_ng_container_0_td_10_Conditional_5_Template, 2, 1, "p", 26);
    \u0275\u0275element(6, "app-torrent-chips", 27);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    let tmp_4_0;
    let tmp_5_0;
    const i_r7 = ctx.$implicit;
    const t_r3 = \u0275\u0275nextContext().$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275property("matTooltip", t_r3("content_types.singular." + ((tmp_4_0 = ctx_r1.item(i_r7).contentType) !== null && tmp_4_0 !== void 0 ? tmp_4_0 : "null")));
    \u0275\u0275advance();
    \u0275\u0275textInterpolate((tmp_5_0 = (tmp_5_0 = ctx_r1.contentTypeInfo(ctx_r1.item(i_r7).contentType)) == null ? null : tmp_5_0.icon) !== null && tmp_5_0 !== void 0 ? tmp_5_0 : "question_mark");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(ctx_r1.item(i_r7).title);
    \u0275\u0275advance();
    \u0275\u0275conditional(ctx_r1.item(i_r7).title !== ctx_r1.item(i_r7).torrent.name ? 5 : -1);
    \u0275\u0275advance();
    \u0275\u0275property("torrentContent", i_r7);
  }
}
function TorrentsTableComponent_ng_container_0_th_12_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 19);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("torrents.size"));
  }
}
function TorrentsTableComponent_ng_container_0_td_13_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "td", 21)(1, "span", 28);
    \u0275\u0275pipe(2, "filesize");
    \u0275\u0275text(3);
    \u0275\u0275pipe(4, "filesize");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const i_r8 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275property("matTooltip", \u0275\u0275pipeBind2(2, 2, ctx_r1.item(i_r8).torrent.size, 10));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(\u0275\u0275pipeBind1(4, 5, ctx_r1.item(i_r8).torrent.size));
  }
}
function TorrentsTableComponent_ng_container_0_th_15_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 19);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("torrents.published"));
  }
}
function TorrentsTableComponent_ng_container_0_td_16_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "td", 29)(1, "abbr", 30);
    \u0275\u0275text(2);
    \u0275\u0275pipe(3, "timeAgo");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const i_r9 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275propertyInterpolate("matTooltip", ctx_r1.item(i_r9).publishedAt);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", \u0275\u0275pipeBind1(3, 2, ctx_r1.item(i_r9).publishedAt), " ");
  }
}
function TorrentsTableComponent_ng_container_0_th_18_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 19)(1, "abbr", 24);
    \u0275\u0275text(2);
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275property("matTooltip", t_r3("torrents.seeders") + " / " + t_r3("torrents.leechers"));
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(t_r3("torrents.s_l"));
  }
}
function TorrentsTableComponent_ng_container_0_td_19_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "td", 21);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    let tmp_4_0;
    const i_r10 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275textInterpolate2(" ", (tmp_4_0 = ctx_r1.item(i_r10).seeders) !== null && tmp_4_0 !== void 0 ? tmp_4_0 : "?", " / ", (tmp_4_0 = ctx_r1.item(i_r10).leechers) !== null && tmp_4_0 !== void 0 ? tmp_4_0 : "?", " ");
  }
}
function TorrentsTableComponent_ng_container_0_th_21_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "th", 31);
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance();
    \u0275\u0275textInterpolate1(" ", t_r3("torrents.magnet"), " ");
  }
}
function TorrentsTableComponent_ng_container_0_td_22_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "td", 21)(1, "a", 32);
    \u0275\u0275element(2, "mat-icon", 33);
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const i_r11 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275advance();
    \u0275\u0275propertyInterpolate("href", ctx_r1.item(i_r11).torrent.magnetUri, \u0275\u0275sanitizeUrl);
  }
}
function TorrentsTableComponent_ng_container_0_td_24_Template(rf, ctx) {
  if (rf & 1) {
    const _r12 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "td", 21)(1, "div", 34);
    \u0275\u0275pipe(2, "async");
    \u0275\u0275elementStart(3, "mat-card", 35)(4, "mat-card-content")(5, "app-torrent-content", 36);
    \u0275\u0275pipe(6, "async");
    \u0275\u0275listener("updated", function TorrentsTableComponent_ng_container_0_td_24_Template_app_torrent_content_updated_5_listener() {
      const i_r13 = \u0275\u0275restoreView(_r12).$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.updated.emit(ctx_r1.item(i_r13).infoHash));
    })("tabSelected", function TorrentsTableComponent_ng_container_0_td_24_Template_app_torrent_content_tabSelected_5_listener($event) {
      const i_r13 = \u0275\u0275restoreView(_r12).$implicit;
      const ctx_r1 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r1.controller.selectTorrent(i_r13.infoHash, $event ? $event : null));
    });
    \u0275\u0275elementEnd()()()()();
  }
  if (rf & 2) {
    let tmp_5_0;
    let tmp_10_0;
    const i_r13 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275attribute("colspan", ctx_r1.displayedColumns.length);
    \u0275\u0275advance();
    \u0275\u0275property("@detailExpand", ((tmp_5_0 = \u0275\u0275pipeBind1(2, 7, ctx_r1.controller.selection$)) == null ? null : tmp_5_0.infoHash) == i_r13.infoHash ? "expanded" : "collapsed");
    \u0275\u0275advance(4);
    \u0275\u0275property("torrentContent", i_r13)("size", false)("published", ctx_r1.breakpoints.sizeAtLeast("Medium"))("peers", ctx_r1.breakpoints.sizeAtLeast("Medium"))("selectedTab", (tmp_10_0 = \u0275\u0275pipeBind1(6, 9, ctx_r1.controller.selection$)) == null ? null : tmp_10_0.tab);
  }
}
function TorrentsTableComponent_ng_container_0_tr_25_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275element(0, "tr", 37);
  }
}
function TorrentsTableComponent_ng_container_0_tr_26_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275element(0, "tr", 38);
    \u0275\u0275pipe(1, "async");
  }
  if (rf & 2) {
    let tmp_4_0;
    const i_r14 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext(2);
    \u0275\u0275classMap("summary-row " + ((tmp_4_0 = \u0275\u0275pipeBind1(1, 2, ctx_r1.controller.selection$)) == null ? null : tmp_4_0.infoHash) == i_r14.infoHash ? "expanded" : "collapsed");
  }
}
function TorrentsTableComponent_ng_container_0_tr_27_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275element(0, "tr", 39);
  }
}
function TorrentsTableComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "div", 1);
    \u0275\u0275element(2, "mat-progress-bar", 2);
    \u0275\u0275pipe(3, "async");
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(4, "table", 3);
    \u0275\u0275elementContainerStart(5, 4);
    \u0275\u0275template(6, TorrentsTableComponent_ng_container_0_th_6_Template, 2, 3, "th", 5)(7, TorrentsTableComponent_ng_container_0_td_7_Template, 2, 1, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(8, 7);
    \u0275\u0275template(9, TorrentsTableComponent_ng_container_0_th_9_Template, 2, 1, "th", 5)(10, TorrentsTableComponent_ng_container_0_td_10_Template, 7, 5, "td", 8);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(11, 9);
    \u0275\u0275template(12, TorrentsTableComponent_ng_container_0_th_12_Template, 2, 1, "th", 5)(13, TorrentsTableComponent_ng_container_0_td_13_Template, 5, 7, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(14, 10);
    \u0275\u0275template(15, TorrentsTableComponent_ng_container_0_th_15_Template, 2, 1, "th", 5)(16, TorrentsTableComponent_ng_container_0_td_16_Template, 4, 4, "td", 11);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(17, 12);
    \u0275\u0275template(18, TorrentsTableComponent_ng_container_0_th_18_Template, 3, 2, "th", 5)(19, TorrentsTableComponent_ng_container_0_td_19_Template, 2, 2, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(20, 13);
    \u0275\u0275template(21, TorrentsTableComponent_ng_container_0_th_21_Template, 2, 1, "th", 14)(22, TorrentsTableComponent_ng_container_0_td_22_Template, 3, 1, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275elementContainerStart(23, 15);
    \u0275\u0275template(24, TorrentsTableComponent_ng_container_0_td_24_Template, 7, 11, "td", 6);
    \u0275\u0275elementContainerEnd();
    \u0275\u0275template(25, TorrentsTableComponent_ng_container_0_tr_25_Template, 1, 0, "tr", 16)(26, TorrentsTableComponent_ng_container_0_tr_26_Template, 2, 4, "tr", 17)(27, TorrentsTableComponent_ng_container_0_tr_27_Template, 1, 0, "tr", 18);
    \u0275\u0275elementEnd();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance(2);
    \u0275\u0275property("mode", \u0275\u0275pipeBind1(3, 7, ctx_r1.dataSource.loading$) ? "indeterminate" : "determinate")("value", 0);
    \u0275\u0275advance(2);
    \u0275\u0275property("dataSource", ctx_r1.dataSource)("multiTemplateDataRows", true);
    \u0275\u0275advance(21);
    \u0275\u0275property("matHeaderRowDef", ctx_r1.displayedColumns);
    \u0275\u0275advance();
    \u0275\u0275property("matRowDefColumns", ctx_r1.displayedColumns);
    \u0275\u0275advance();
    \u0275\u0275property("matRowDefColumns", \u0275\u0275pureFunction0(9, _c0));
  }
}
var TorrentsTableComponent = class _TorrentsTableComponent {
  constructor() {
    this.route = inject(ActivatedRoute);
    this.router = inject(Router);
    this.breakpoints = inject(BreakpointsService);
    this.contentTypeInfo = contentTypeInfo;
    this.displayedColumns = allColumns;
    this.updated = new EventEmitter();
    this.items = Array();
  }
  ngOnInit() {
    this.dataSource.items$.subscribe((items) => {
      this.items = items;
    });
  }
  /** Whether the number of selected elements matches the total number of rows. */
  isAllSelected() {
    return this.items.every((i) => this.multiSelection.isSelected(i.infoHash));
  }
  /** Selects all rows if they are not all selected; otherwise clear selection. */
  toggleAllRows() {
    if (this.isAllSelected()) {
      this.multiSelection.clear();
      return;
    }
    this.multiSelection.select(...this.items.map((i) => i.infoHash));
  }
  toggleSelectedTorrent(infoHash) {
    this.controller.update((ctrl) => __spreadProps(__spreadValues({}, ctrl), {
      selectedTorrent: ctrl.selectedTorrent?.infoHash === infoHash ? void 0 : {
        infoHash,
        tab: ctrl.selectedTorrent?.tab
      }
    }));
  }
  /**
   * Workaround for untyped table cell definitions
   */
  item(item) {
    return item;
  }
  static {
    this.\u0275fac = function TorrentsTableComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _TorrentsTableComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _TorrentsTableComponent, selectors: [["app-torrents-table"]], inputs: { dataSource: "dataSource", controller: "controller", multiSelection: "multiSelection", displayedColumns: "displayedColumns" }, outputs: { updated: "updated" }, standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], [1, "progress-bar-container"], [3, "mode", "value"], ["mat-table", "", 1, "table-torrents", 3, "dataSource", "multiTemplateDataRows"], ["matColumnDef", "select"], ["mat-header-cell", "", 4, "matHeaderCellDef"], ["mat-cell", "", 4, "matCellDef"], ["matColumnDef", "summary"], ["mat-cell", "", 3, "click", 4, "matCellDef"], ["matColumnDef", "size"], ["matColumnDef", "publishedAt"], ["class", "td-published-at", "mat-cell", "", 4, "matCellDef"], ["matColumnDef", "peers"], ["matColumnDef", "magnet"], ["mat-header-cell", "", "style", "text-align: center", 4, "matHeaderCellDef"], ["matColumnDef", "expandedDetail"], ["mat-header-row", "", 4, "matHeaderRowDef"], ["mat-row", "", 3, "class", 4, "matRowDef", "matRowDefColumns"], ["mat-row", "", "class", "expanded-detail-row", 4, "matRowDef", "matRowDefColumns"], ["mat-header-cell", ""], [3, "change", "checked", "indeterminate", "matTooltip"], ["mat-cell", ""], [3, "click", "change", "checked"], ["mat-cell", "", 3, "click"], [3, "matTooltip"], [1, "title"], [1, "original-name"], [3, "torrentContent"], [1, "filesize", 3, "matTooltip"], ["mat-cell", "", 1, "td-published-at"], ["matTooltipClass", "tooltip-published-at", 3, "matTooltip"], ["mat-header-cell", "", 2, "text-align", "center"], [3, "href"], ["svgIcon", "magnet"], [1, "item-detail"], [1, "torrent-permalink"], [3, "updated", "tabSelected", "torrentContent", "size", "published", "peers", "selectedTab"], ["mat-header-row", ""], ["mat-row", ""], ["mat-row", "", 1, "expanded-detail-row"]], template: function TorrentsTableComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, TorrentsTableComponent_ng_container_0_Template, 28, 10, "ng-container", 0);
      }
    }, dependencies: [
      AppModule,
      MatCard,
      MatCardContent,
      MatCheckbox,
      MatIcon,
      MatProgressBar,
      MatTable,
      MatHeaderCellDef,
      MatHeaderRowDef,
      MatColumnDef,
      MatCellDef,
      MatRowDef,
      MatHeaderCell,
      MatCell,
      MatHeaderRow,
      MatRow,
      MatTooltip,
      TranslocoDirective,
      AsyncPipe,
      FilesizePipe,
      TimeAgoPipe,
      TorrentChipsComponent,
      TorrentContentComponent
    ], styles: ["\n\n.progress-bar-container[_ngcontent-%COMP%] {\n  height: 10px;\n}\nth.cdk-column-select[_ngcontent-%COMP%], \ntd.cdk-column-select[_ngcontent-%COMP%] {\n  padding-right: 0;\n}\ntd.mat-column-summary[_ngcontent-%COMP%] {\n  vertical-align: middle;\n  cursor: pointer;\n  white-space: pre-wrap;\n  padding-top: 8px;\n  padding-bottom: 8px;\n}\ntd.mat-column-summary[_ngcontent-%COMP%]   .title[_ngcontent-%COMP%] {\n  line-height: 30px;\n  overflow: hidden;\n  margin-right: 20px;\n  font-weight: bold;\n}\ntd.mat-column-summary[_ngcontent-%COMP%]   .original-name[_ngcontent-%COMP%] {\n  margin: 2px 0 8px 34px;\n}\ntd.mat-column-summary[_ngcontent-%COMP%]   .title[_ngcontent-%COMP%], \ntd.mat-column-summary[_ngcontent-%COMP%]   .original-name[_ngcontent-%COMP%] {\n  white-space: pre-wrap;\n  word-break: break-word;\n  overflow-wrap: break-word;\n}\ntd.mat-column-summary[_ngcontent-%COMP%]    > .mat-icon[_ngcontent-%COMP%] {\n  display: inline-block;\n  position: relative;\n  top: 6px;\n  margin-right: 10px;\n}\ntd.mat-column-summary[_ngcontent-%COMP%]   mat-chip-set[_ngcontent-%COMP%] {\n  display: inline-block;\n  margin-left: 10px;\n}\ntd.mat-column-summary[_ngcontent-%COMP%]   mat-chip-set[_ngcontent-%COMP%]   mat-chip[_ngcontent-%COMP%] {\n  margin: 2px 10px 2px 0;\n}\ntr.expanded-detail-row[_ngcontent-%COMP%] {\n  height: 0;\n}\ntr.mat-mdc-row.expanded[_ngcontent-%COMP%]   td[_ngcontent-%COMP%] {\n  border-bottom: 0;\n}\napp-torrent-content[_ngcontent-%COMP%] {\n  padding-top: 20px;\n  padding-bottom: 20px;\n}\n.mat-column-magnet[_ngcontent-%COMP%] {\n  text-align: center;\n}\n.mat-column-magnet[_ngcontent-%COMP%]   .mat-icon[_ngcontent-%COMP%] {\n  position: relative;\n  top: 3px;\n}\n.item-detail[_ngcontent-%COMP%] {\n  width: 100%;\n  overflow: hidden;\n}\n.td-published-at[_ngcontent-%COMP%]   abbr[_ngcontent-%COMP%] {\n  cursor: default;\n  text-decoration: underline;\n  text-decoration-style: dotted;\n}\n.cdk-column-peers[_ngcontent-%COMP%] {\n  white-space: nowrap;\n}\nspan.filesize[_ngcontent-%COMP%] {\n  text-decoration: underline;\n  text-decoration-style: dotted;\n  cursor: default;\n}\n/*# sourceMappingURL=torrents-table.component.css.map */"], data: { animation: [
      trigger("detailExpand", [
        state("collapsed,void", style({ height: "0px", minHeight: "0" })),
        state("expanded", style({ height: "*" })),
        transition("expanded <=> collapsed", animate("225ms cubic-bezier(0.4, 0.0, 0.2, 1)"))
      ])
    ] } });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(TorrentsTableComponent, { className: "TorrentsTableComponent", filePath: "src/app/torrents/torrents-table.component.ts", lineNumber: 52 });
})();
var allColumns = [
  "select",
  "summary",
  "size",
  "publishedAt",
  "peers",
  "magnet"
];
var compactColumns = ["select", "summary", "size", "magnet"];

export {
  TorrentsBulkActionsComponent,
  TorrentsTableComponent,
  allColumns,
  compactColumns
};
//# sourceMappingURL=chunk-QLUBXXM3.js.map
