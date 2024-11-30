import{a as Nt,b as jt,c as Gt}from"./chunk-GHFSK4SQ.js";import{b as Qt}from"./chunk-ENCA42F6.js";import{c as Lt}from"./chunk-PX3TOMGC.js";import{a as dt}from"./chunk-ASLGZ7DJ.js";import{a as H}from"./chunk-MSAOOVCY.js";import{$ as wt,Da as Mt,Ea as kt,Fa as Et,Ga as At,Ha as Pt,I as _t,Ia as Ot,Ja as Vt,K as ut,Ka as Dt,La as $t,M as gt,Ma as Bt,N as xt,Na as Rt,O as Ct,Oa as Ht,P as ht,Pa as zt,Qa as Ft,R as Tt,T as ft,Ua as R,X as bt,Z as vt,aa as St,b as D,ba as yt,d as rt,e as lt,f as st,g as ct,n as mt,o as pt,qa as $,va as It,ya as B}from"./chunk-VAEZNV34.js";import{c as it,d as ot,g as Q,h as j,i as at}from"./chunk-6XXA7HXI.js";import{i as et,k as nt}from"./chunk-CMNWCZJM.js";import{$b as l,Cb as d,Ea as _,Fa as u,Gb as Y,Hb as p,Kb as J,Mb as y,N as k,Na as P,Nb as z,Ob as F,Pb as L,Qb as s,Rb as r,Sb as T,Tb as f,Ub as b,Wb as v,Zb as x,ad as tt,cc as S,h as U,ha as I,hb as K,jc as N,kb as o,kc as c,l as q,lc as g,mc as h,nc as W,o as M,pc as O,qa as w,qc as X,ua as A,vc as E,wc as V,xc as Z}from"./chunk-Z3WUIYN5.js";var Ae=(e,a)=>{let n=G(e,a)?.split(",").map(i=>i.trim()).filter(Boolean);return n?.length?Array.from(new Set(n)).sort():void 0},G=(e,a)=>typeof e[a]=="string"&&decodeURIComponent(e[a])||void 0,Pe=(e,a)=>{if(e&&e[a]&&/^\d+$/.test(e[a]))return parseInt(e[a])};function Jt(e,a){if(e&1&&(s(0,"span",7),c(1),r()),e&2){let t=l(2).$implicit;o(),g(t("torrents.copy"))}}function Wt(e,a){if(e&1&&(s(0,"mat-icon"),c(1,"content_copy"),r(),d(2,Jt,2,1,"span",7)),e&2){let t=l(2);o(2),y(t.breakpoints.sizeAtLeast("Medium")?2:-1)}}function Xt(e,a){if(e&1&&(s(0,"mat-card")(1,"mat-card-actions",8)(2,"button",9),T(3,"mat-icon",10),c(4),r(),s(5,"button",9)(6,"mat-icon"),c(7,"tag"),r(),c(8),r()()()),e&2){let t=l().$implicit,n=l();o(2),p("disabled",!n.selectedItems.length)("matTooltip",t("torrents.copy_to_clipboard"))("cdkCopyToClipboard",n.getSelectedMagnetLinks()),o(2),h("",t("torrents.magnet_links")," "),o(),p("disabled",!n.selectedItems.length)("matTooltip",t("torrents.copy_to_clipboard"))("cdkCopyToClipboard",n.getSelectedInfoHashes()),o(3),h("",t("torrents.info_hashes")," ")}}function Zt(e,a){if(e&1&&(s(0,"span",7),c(1),r()),e&2){let t=l(2).$implicit;o(),g(t("torrents.edit_tags"))}}function te(e,a){if(e&1&&(s(0,"mat-icon"),c(1,"sell"),r(),d(2,Zt,2,1,"span",7)),e&2){let t=l(2);o(2),y(t.breakpoints.sizeAtLeast("Medium")?2:-1)}}function ee(e,a){if(e&1){let t=v();s(0,"mat-chip-row",18),x("edited",function(i){let m=_(t).$implicit,C=l(3);return u(C.renameTag(m,i.value))})("removed",function(){let i=_(t).$implicit,m=l(3);return u(m.deleteTag(i))}),c(1),s(2,"mat-icon",19),c(3,"cancel"),r()()}if(e&2){let t=a.$implicit;p("editable",!0)("aria-description","press enter to edit"),o(),h(" ",t," ")}}function ne(e,a){if(e&1&&(s(0,"mat-option",16),c(1),r()),e&2){let t=a.$implicit;p("value",t),o(),g(t)}}function ie(e,a){if(e&1){let t=v();s(0,"mat-card")(1,"mat-form-field",11)(2,"mat-chip-grid",12,0),F(4,ee,4,3,"mat-chip-row",13,z),r(),s(6,"input",14),x("matChipInputTokenEnd",function(i){_(t);let m=l(2);return u(i.value&&m.addTag(i.value))}),r(),s(7,"mat-autocomplete",15,1),x("optionSelected",function(i){_(t);let m=l(2);return u(m.addTag(i.option.viewValue))}),F(9,ne,2,2,"mat-option",16,z),r()(),s(11,"mat-card-actions",8)(12,"button",17),x("click",function(){_(t);let i=l(2);return u(i.setTags())}),c(13),r(),s(14,"button",17),x("click",function(){_(t);let i=l(2);return u(i.putTags())}),c(15),r(),s(16,"button",17),x("click",function(){_(t);let i=l(2);return u(i.deleteTags())}),c(17),r()()()}if(e&2){let t=N(3),n=N(8),i=l().$implicit,m=l();o(4),L(m.editedTags),o(2),S("placeholder",i("torrents.tags.placeholder")),p("formControl",m.newTagCtrl)("matAutocomplete",n)("matChipInputFor",t)("matChipInputSeparatorKeyCodes",m.separatorKeysCodes)("value",m.newTagCtrl.value),o(3),L(m.suggestedTags),o(3),S("matTooltip",i("torrents.tags.set_tip")),p("disabled",!m.selectedItems.length),o(),h(" ",i("torrents.tags.set")," "),o(),S("matTooltip",i("torrents.tags.put_tip")),p("disabled",!m.selectedItems.length||!m.editedTags.length&&!m.newTagCtrl.value),o(),h(" ",i("torrents.tags.put")," "),o(),S("matTooltip",i("torrents.tags.delete_tip")),p("disabled",!m.selectedItems.length||!m.editedTags.length&&!m.newTagCtrl.value),o(),h(" ",i("torrents.tags.delete")," ")}}function oe(e,a){if(e&1&&(s(0,"span",7),c(1),r()),e&2){let t=l(2).$implicit;o(),g(t("torrents.delete"))}}function ae(e,a){if(e&1&&(s(0,"mat-icon"),c(1,"delete_forever"),r(),d(2,oe,2,1,"span",7)),e&2){let t=l(2);o(2),y(t.breakpoints.sizeAtLeast("Medium")?2:-1)}}function re(e,a){if(e&1){let t=v();s(0,"mat-card")(1,"mat-card-content")(2,"p")(3,"strong"),c(4),r(),T(5,"br"),c(6),r()(),s(7,"mat-card-actions",8)(8,"button",20),x("click",function(){_(t);let i=l(2);return u(i.deleteTorrents())}),s(9,"mat-icon"),c(10,"delete_forever"),r(),c(11),r()()()}if(e&2){let t=l().$implicit,n=l();o(4),g(t("torrents.delete_are_you_sure")),o(2),h("",t("torrents.delete_action_cannot_be_undone"),". "),o(2),p("disabled",!n.selectedItems.length),o(3),h("",t("torrents.delete")," ")}}function le(e,a){e&1&&(s(0,"mat-icon",21),c(1,"close"),r())}function se(e,a){e&1&&(s(0,"mat-tab"),d(1,le,2,0,"ng-template",5),r())}function ce(e,a){if(e&1){let t=v();f(0),s(1,"mat-tab-group",3),x("focusChange",function(i){_(t);let m=l();return u(m.selectTab(i.index==4?0:i.index))}),T(2,"mat-tab",4),s(3,"mat-tab"),d(4,Wt,3,1,"ng-template",5)(5,Xt,9,8,"ng-template",6),r(),s(6,"mat-tab"),d(7,te,3,1,"ng-template",5)(8,ie,18,15,"ng-template",6),r(),s(9,"mat-tab"),d(10,ae,3,1,"ng-template",5)(11,re,12,4,"ng-template",6),r(),d(12,se,2,0,"mat-tab"),r(),b()}if(e&2){let t=l();o(),p("selectedIndex",t.selectedTabIndex)("mat-stretch-tabs",!1),o(),p("aria-labelledby","hidden"),o(10),y(t.selectedTabIndex>0?12:-1)}}var tn=(()=>{class e{constructor(){this.graphQLService=w(_t),this.errorsService=w(dt),this.breakpoints=w(H),this.selectedItems$=new U,this.updated=new P,this.separatorKeysCodes=[13,188],this.selectedTabIndex=0,this.newTagCtrl=new st(""),this.editedTags=Array(),this.suggestedTags=Array(),this.selectedItems=new Array}ngOnInit(){this.selectedItems$.subscribe(t=>{this.selectedItems=t})}selectTab(t){this.selectedTabIndex=t}getSelectedMagnetLinks(){return this.selectedItems.map(t=>t.torrent.magnetUri).join(`
`)}getSelectedInfoHashes(){return this.selectedItems.map(t=>t.infoHash).join(`
`)}addTag(t){this.editedTags.includes(t)||this.editedTags.push(t),this.newTagCtrl.reset(),this.updateSuggestedTags()}deleteTag(t){this.editedTags=this.editedTags.filter(n=>n!==t),this.updateSuggestedTags()}renameTag(t,n){this.editedTags=this.editedTags.map(i=>i===t?n:i),this.updateSuggestedTags()}putTags(){let t=this.selectedItems.map(({infoHash:n})=>n);if(t.length)return this.newTagCtrl.value&&this.addTag(this.newTagCtrl.value),this.graphQLService.torrentPutTags({infoHashes:t,tagNames:this.editedTags}).pipe(k(n=>(this.errorsService.addError(`Error putting tags: ${n.message}`),M))).pipe(I(()=>{this.updated.emit()})).subscribe()}setTags(){let t=this.selectedItems.map(({infoHash:n})=>n);if(t.length)return this.newTagCtrl.value&&this.addTag(this.newTagCtrl.value),this.graphQLService.torrentSetTags({infoHashes:t,tagNames:this.editedTags}).pipe(k(n=>(this.errorsService.addError(`Error setting tags: ${n.message}`),M))).pipe(I(()=>{this.updated.emit()})).subscribe()}deleteTags(){let t=this.selectedItems.map(({infoHash:n})=>n);if(t.length)return this.newTagCtrl.value&&this.addTag(this.newTagCtrl.value),this.graphQLService.torrentDeleteTags({infoHashes:t,tagNames:this.editedTags}).pipe(k(n=>(this.errorsService.addError(`Error deleting tags: ${n.message}`),M))).pipe(I(()=>{this.updated.emit()})).subscribe()}updateSuggestedTags(){return this.graphQLService.torrentSuggestTags({input:{prefix:this.newTagCtrl.value,exclusions:this.editedTags}}).pipe(I(t=>{this.suggestedTags.splice(0,this.suggestedTags.length,...t.suggestions.map(n=>n.name))})).subscribe()}deleteTorrents(){let t=this.selectedItems.map(({infoHash:n})=>n);this.graphQLService.torrentDelete({infoHashes:t}).pipe(k(n=>(this.errorsService.addError(`Error deleting torrents: ${n.message}`),M))).pipe(I(()=>{this.updated.emit()})).subscribe()}static{this.\u0275fac=function(n){return new(n||e)}}static{this.\u0275cmp=A({type:e,selectors:[["app-torrents-bulk-actions"]],inputs:{selectedItems$:"selectedItems$"},outputs:{updated:"updated"},standalone:!0,features:[O],decls:1,vars:0,consts:[["chipGrid",""],["auto","matAutocomplete"],[4,"transloco"],["animationDuration","0",1,"tab-group-bulk-actions",3,"focusChange","selectedIndex","mat-stretch-tabs"],[1,"bulk-tab-placeholder",3,"aria-labelledby"],["mat-tab-label",""],["matTabContent",""],[1,"label"],[1,"button-row"],["mat-stroked-button","",3,"disabled","matTooltip","cdkCopyToClipboard"],["svgIcon","magnet"],["subscriptSizing","dynamic",1,"form-edit-tags"],["aria-label","Enter tags"],[3,"editable","aria-description"],[3,"matChipInputTokenEnd","placeholder","formControl","matAutocomplete","matChipInputFor","matChipInputSeparatorKeyCodes","value"],[3,"optionSelected"],[3,"value"],["mat-stroked-button","","color","primary",3,"click","disabled","matTooltip"],[3,"edited","removed","editable","aria-description"],["matChipRemove",""],["mat-stroked-button","","color","warning",3,"click","disabled"],[2,"margin-right","0"]],template:function(n,i){n&1&&d(0,ce,13,4,"ng-container",2)},dependencies:[R,ut,xt,mt,Ct,pt,ht,ft,Tt,wt,St,bt,vt,gt,$,Rt,Ht,zt,Ft,B,rt,lt,ct,D],styles:["mat-tab-group[_ngcontent-%COMP%]{padding-left:10px}.mat-mdc-card[_ngcontent-%COMP%]{margin-bottom:10px}button[_ngcontent-%COMP%]{margin-right:10px}p[_ngcontent-%COMP%]{margin-top:0}  .mdc-tab[aria-labelledby=hidden]{display:none}"]})}}return e})();var me=()=>["expandedDetail"];function pe(e,a){if(e&1){let t=v();s(0,"th",19)(1,"mat-checkbox",20),x("change",function(){_(t);let i=l(2);return u(i.toggleAllRows())}),r()()}if(e&2){let t=l().$implicit,n=l();o(),p("checked",n.selection.hasValue()&&n.isAllSelected())("indeterminate",n.selection.hasValue()&&!n.isAllSelected())("matTooltip",n.isAllSelected()?t("torrents.deselect_all"):t("torrents.select_all"))}}function de(e,a){if(e&1){let t=v();s(0,"td",21)(1,"mat-checkbox",22),x("click",function(i){return _(t),u(i.stopPropagation())})("change",function(i){let m=_(t).$implicit,C=l(2);return u(i?C.selection.toggle(C.item(m).infoHash):null)}),r()()}if(e&2){let t=a.$implicit,n=l(2);o(),p("checked",n.selection.isSelected(n.item(t).infoHash))}}function _e(e,a){if(e&1&&(s(0,"th",19),c(1),r()),e&2){let t=l().$implicit;o(),g(t("torrents.summary"))}}function ue(e,a){if(e&1&&(s(0,"p",26),c(1),r()),e&2){let t=l().$implicit,n=l(2);o(),g(n.item(t).torrent.name)}}function ge(e,a){if(e&1){let t=v();s(0,"td",23),x("click",function(i){let m=_(t).$implicit;return l(2).toggleTorrentContentId(m.id),u(i.stopPropagation())}),s(1,"mat-icon",24),c(2),r(),s(3,"span",25),c(4),r(),d(5,ue,2,1,"p",26),T(6,"app-torrent-chips",27),r()}if(e&2){let t,n,i=a.$implicit,m=l().$implicit,C=l();o(),p("matTooltip",m("content_types.singular."+((t=C.item(i).contentType)!==null&&t!==void 0?t:"null"))),o(),g((n=(n=C.contentTypeInfo(C.item(i).contentType))==null?null:n.icon)!==null&&n!==void 0?n:"question_mark"),o(2),g(C.item(i).title),o(),y(C.item(i).title!==C.item(i).torrent.name?5:-1),o(),p("torrentContent",i)}}function xe(e,a){if(e&1&&(s(0,"th",19),c(1),r()),e&2){let t=l().$implicit;o(),g(t("torrents.size"))}}function Ce(e,a){if(e&1&&(s(0,"td",21)(1,"span",28),E(2,"filesize"),c(3),E(4,"filesize"),r()()),e&2){let t=a.$implicit,n=l(2);o(),p("matTooltip",Z(2,2,n.item(t).torrent.size,10)),o(2),g(V(4,5,n.item(t).torrent.size))}}function he(e,a){if(e&1&&(s(0,"th",19),c(1),r()),e&2){let t=l().$implicit;o(),g(t("torrents.published"))}}function Te(e,a){if(e&1&&(s(0,"td",29)(1,"abbr",30),c(2),E(3,"timeAgo"),r()()),e&2){let t=a.$implicit,n=l(2);o(),S("matTooltip",n.item(t).publishedAt),o(),h(" ",V(3,2,n.item(t).publishedAt)," ")}}function fe(e,a){if(e&1&&(s(0,"th",19)(1,"abbr",24),c(2),r()()),e&2){let t=l().$implicit;o(),p("matTooltip",t("torrents.seeders")+" / "+t("torrents.leechers")),o(),g(t("torrents.s_l"))}}function be(e,a){if(e&1&&(s(0,"td",21),c(1),r()),e&2){let t,n=a.$implicit,i=l(2);o(),W(" ",(t=i.item(n).seeders)!==null&&t!==void 0?t:"?"," / ",(t=i.item(n).leechers)!==null&&t!==void 0?t:"?"," ")}}function ve(e,a){if(e&1&&(s(0,"th",31),c(1),r()),e&2){let t=l().$implicit;o(),h(" ",t("torrents.magnet")," ")}}function we(e,a){if(e&1&&(s(0,"td",21)(1,"a",32),T(2,"mat-icon",33),r()()),e&2){let t=a.$implicit,n=l(2);o(),S("href",n.item(t).torrent.magnetUri,K)}}function Se(e,a){if(e&1){let t=v();s(0,"td",21)(1,"div",34)(2,"app-torrent-content",35),x("updated",function(){let i=_(t).$implicit,m=l(2);return u(m.updated.emit(m.item(i).infoHash))}),r()()()}if(e&2){let t=a.$implicit,n=l(2);Y("colspan",n.displayedColumns.length),o(),p("@detailExpand",n.expandedId.getValue()===t.id?"expanded":"collapsed"),o(),p("torrentContent",t)("size",!1)("published",n.breakpoints.sizeAtLeast("Medium"))("peers",n.breakpoints.sizeAtLeast("Medium"))}}function ye(e,a){e&1&&T(0,"tr",36)}function Ie(e,a){if(e&1&&T(0,"tr",37),e&2){let t=a.$implicit,n=l(2);J("summary-row "+(t.id===n.expandedId.getValue()?"expanded":"collapsed"))}}function Me(e,a){e&1&&T(0,"tr",38)}function ke(e,a){if(e&1&&(f(0),s(1,"div",1),T(2,"mat-progress-bar",2),E(3,"async"),r(),s(4,"table",3),f(5,4),d(6,pe,2,3,"th",5)(7,de,2,1,"td",6),b(),f(8,7),d(9,_e,2,1,"th",5)(10,ge,7,5,"td",8),b(),f(11,9),d(12,xe,2,1,"th",5)(13,Ce,5,7,"td",6),b(),f(14,10),d(15,he,2,1,"th",5)(16,Te,4,4,"td",11),b(),f(17,12),d(18,fe,3,2,"th",5)(19,be,2,2,"td",6),b(),f(20,13),d(21,ve,2,1,"th",14)(22,we,3,1,"td",6),b(),f(23,15),d(24,Se,3,6,"td",6),b(),d(25,ye,1,0,"tr",16)(26,Ie,1,2,"tr",17)(27,Me,1,0,"tr",18),r(),b()),e&2){let t=l();o(2),p("mode",V(3,7,t.dataSource.loading$)?"indeterminate":"determinate")("value",0),o(2),p("dataSource",t.dataSource)("multiTemplateDataRows",!0),o(21),p("matHeaderRowDef",t.displayedColumns),o(),p("matRowDefColumns",t.displayedColumns),o(),p("matRowDefColumns",X(9,me))}}var vn=(()=>{class e{constructor(){this.route=w(et),this.router=w(nt),this.breakpoints=w(H),this.contentTypeInfo=Lt,this.displayedColumns=Ee,this.updated=new P,this.items=Array(),this.expandedId=new q(null)}ngOnInit(){this.dataSource.items$.subscribe(t=>{this.items=t}),this.route.queryParams.subscribe(t=>{let n=this.expandedId.getValue()??void 0,i=G(t,"expanded");n!==i&&this.expandedId.next(i??null)}),this.expandedId.subscribe(t=>{this.router.navigate([],{relativeTo:this.route,queryParams:{expanded:t?encodeURIComponent(t):void 0},queryParamsHandling:"merge"})})}isAllSelected(){return this.items.every(t=>this.selection.isSelected(t.infoHash))}toggleAllRows(){if(this.isAllSelected()){this.selection.clear();return}this.selection.select(...this.items.map(t=>t.infoHash))}toggleTorrentContentId(t){this.expandedId.getValue()===t?this.expandedId.next(null):this.expandedId.next(t)}item(t){return t}static{this.\u0275fac=function(n){return new(n||e)}}static{this.\u0275cmp=A({type:e,selectors:[["app-torrents-table"]],inputs:{dataSource:"dataSource",selection:"selection",displayedColumns:"displayedColumns"},outputs:{updated:"updated"},standalone:!0,features:[O],decls:1,vars:0,consts:[[4,"transloco"],[1,"progress-bar-container"],[3,"mode","value"],["mat-table","",1,"table-torrents",3,"dataSource","multiTemplateDataRows"],["matColumnDef","select"],["mat-header-cell","",4,"matHeaderCellDef"],["mat-cell","",4,"matCellDef"],["matColumnDef","summary"],["mat-cell","",3,"click",4,"matCellDef"],["matColumnDef","size"],["matColumnDef","publishedAt"],["class","td-published-at","mat-cell","",4,"matCellDef"],["matColumnDef","peers"],["matColumnDef","magnet"],["mat-header-cell","","style","text-align: center",4,"matHeaderCellDef"],["matColumnDef","expandedDetail"],["mat-header-row","",4,"matHeaderRowDef"],["mat-row","",3,"class",4,"matRowDef","matRowDefColumns"],["mat-row","","class","expanded-detail-row",4,"matRowDef","matRowDefColumns"],["mat-header-cell",""],[3,"change","checked","indeterminate","matTooltip"],["mat-cell",""],[3,"click","change","checked"],["mat-cell","",3,"click"],[3,"matTooltip"],[1,"title"],[1,"original-name"],[3,"torrentContent"],[1,"filesize",3,"matTooltip"],["mat-cell","",1,"td-published-at"],["matTooltipClass","tooltip-published-at",3,"matTooltip"],["mat-header-cell","",2,"text-align","center"],[3,"href"],["svgIcon","magnet"],[1,"item-detail"],[3,"updated","torrentContent","size","published","peers"],["mat-header-row",""],["mat-row",""],["mat-row","",1,"expanded-detail-row"]],template:function(n,i){n&1&&d(0,ke,28,10,"ng-container",0)},dependencies:[R,yt,$,It,Mt,Et,Vt,At,kt,Dt,Pt,Ot,$t,Bt,B,D,tt,Nt,Qt,jt,Gt],styles:[".progress-bar-container[_ngcontent-%COMP%]{height:10px}tr.expanded-detail-row[_ngcontent-%COMP%]   td[_ngcontent-%COMP%]{border-bottom-width:0}tr.expanded[_ngcontent-%COMP%] + tr.expanded-detail-row[_ngcontent-%COMP%]   td[_ngcontent-%COMP%]{border-bottom-width:1px}th.cdk-column-select[_ngcontent-%COMP%], td.cdk-column-select[_ngcontent-%COMP%]{padding-right:0}td.mat-column-summary[_ngcontent-%COMP%]{vertical-align:middle;cursor:pointer;white-space:pre-wrap;padding-top:8px;padding-bottom:8px}td.mat-column-summary[_ngcontent-%COMP%]   .title[_ngcontent-%COMP%]{line-height:30px;overflow:hidden;margin-right:20px;font-weight:700}td.mat-column-summary[_ngcontent-%COMP%]   .original-name[_ngcontent-%COMP%]{margin:2px 0 8px 34px}td.mat-column-summary[_ngcontent-%COMP%]   .title[_ngcontent-%COMP%], td.mat-column-summary[_ngcontent-%COMP%]   .original-name[_ngcontent-%COMP%]{white-space:pre-wrap;word-break:break-word;overflow-wrap:break-word}td.mat-column-summary[_ngcontent-%COMP%] > .mat-icon[_ngcontent-%COMP%]{display:inline-block;position:relative;top:6px;margin-right:10px}td.mat-column-summary[_ngcontent-%COMP%]   mat-chip-set[_ngcontent-%COMP%]{display:inline-block;margin-left:10px}td.mat-column-summary[_ngcontent-%COMP%]   mat-chip-set[_ngcontent-%COMP%]   mat-chip[_ngcontent-%COMP%]{margin:2px 10px 2px 0}tr.expanded-detail-row[_ngcontent-%COMP%]{height:0}tr.mat-mdc-row.expanded[_ngcontent-%COMP%]   td[_ngcontent-%COMP%]{border-bottom:0}tr.mat-mdc-row.expanded[_ngcontent-%COMP%] + .expanded-detail-row[_ngcontent-%COMP%] > td[_ngcontent-%COMP%]{padding-bottom:10px}.mat-column-magnet[_ngcontent-%COMP%]{text-align:center}.mat-column-magnet[_ngcontent-%COMP%]   .mat-icon[_ngcontent-%COMP%]{position:relative;top:3px}.item-detail[_ngcontent-%COMP%]{width:100%;overflow:hidden}.td-published-at[_ngcontent-%COMP%]   abbr[_ngcontent-%COMP%]{cursor:default;text-decoration:underline;text-decoration-style:dotted}.cdk-column-peers[_ngcontent-%COMP%]{white-space:nowrap}span.filesize[_ngcontent-%COMP%]{text-decoration:underline;text-decoration-style:dotted}"],data:{animation:[it("detailExpand",[j("collapsed,void",Q({height:"0px",minHeight:"0"})),j("expanded",Q({height:"*"})),at("expanded <=> collapsed",ot("225ms cubic-bezier(0.4, 0.0, 0.2, 1)"))])]}})}}return e})(),Ee=["select","summary","size","publishedAt","peers","magnet"],wn=["select","summary","size","magnet"];export{Ae as a,G as b,Pe as c,tn as d,vn as e,Ee as f,wn as g};
