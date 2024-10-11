import{a as D,b as N,c as z,d as De,e as Ne,f as ze,g as Qe}from"./chunk-Y6LBPD3T.js";import{c as Ae,d as Re}from"./chunk-HXCDL7LV.js";import{b as Le}from"./chunk-OS3AV5XL.js";import{a as Ce}from"./chunk-AMUSMSN6.js";import{a as ye}from"./chunk-BRPPGP57.js";import{$ as we,A as ve,B as Se,Oa as Be,S as be,Ta as je,Ua as qe,_ as Te,a as se,aa as Me,b as le,ba as Oe,d as pe,e as ue,f as me,fa as Pe,g as de,ga as ke,h as _e,i as ge,n as fe,na as Ie,oa as Ee,pa as $e,q as he,qa as Fe,r as xe,ra as Ve}from"./chunk-EH7PJXWD.js";import"./chunk-6XXA7HXI.js";import{j as ae,l as ce}from"./chunk-FKMTSCBK.js";import{$ as X,$b as l,$c as re,B as M,Cb as C,Ea as _,Fa as g,Gb as ee,Hb as f,Kb as L,Mb as b,N as W,O as Y,Ob as k,Pb as I,Q as K,Qb as p,Rb as u,Sb as A,Tb as te,Ub as ne,Wb as P,Zb as h,_c as oe,a,b as s,jc as R,kb as c,kc as m,l as y,lc as E,mc as T,nc as $,o as G,pc as ie,qa as O,ua as Z,uc as v,vc as S}from"./chunk-3DR3CJRN.js";var Q={items:[],totalCount:0,totalCountIsEstimate:!1,aggregations:{}},V=class{constructor(n,e,o){this.apollo=n,this.errorsService=e,this.currentRequest=new y(0),this.loadingSubject=new y(!1),this.loading$=this.loadingSubject.asObservable(),this.result=Q,this.resultSubject=new y(this.result),this.result$=this.resultSubject.asObservable(),this.items$=this.resultSubject.pipe(M(i=>i.items)),this.overallTotalCount$=this.resultSubject.pipe(M(i=>{let r=0,d=!1;for(let w of i.aggregations.contentType??[])r+=w.count,d=d||w.isEstimate;return{count:r,isEstimate:d}})),this.availableContentTypes$=this.resultSubject.pipe(X((i,r)=>Array.from(new Set([...i,...(r.aggregations.contentType??[]).flatMap(d=>d.value?[d.value]:[])])),[])),this.contentTypeCounts$=this.resultSubject.pipe(M(i=>Object.fromEntries((i.aggregations.contentType??[]).map(r=>[r.value,{count:r.count,isEstimate:r.isEstimate}])))),o.subscribe(i=>{this.input=i.input,this.loadResult({input:s(a({},i.input),{cached:!0})})}),this.resultSubject.subscribe(i=>{this.result=i})}connect({}){return this.items$}disconnect(){this.resultSubject.complete()}refresh(){this.loadResult({input:s(a({},this.input),{cached:!1})})}loadResult(n){this.currentSubscription&&(this.currentSubscription.unsubscribe(),this.currentSubscription=void 0),this.loadingSubject.next(!0);let e=this.currentRequest.getValue()+1;this.currentRequest.next(e);let o=this.apollo.query({query:Be,variables:n,fetchPolicy:"no-cache"}).pipe(M(i=>i.data.torrentContent.search)).pipe(W(i=>(this.errorsService.addError(`Error loading item results: ${i.message}`),G)));this.currentSubscription=o.subscribe(i=>{e===this.currentRequest.getValue()&&(this.loadingSubject.next(!1),this.resultSubject.next(i))})}};var He=t=>({input:{queryString:t.queryString,limit:t.limit,page:t.page,totalCount:!0,hasNextPage:!0,orderBy:[t.orderBy],facets:{contentType:{aggregate:!0,filter:t.contentType?[t.contentType==="null"?null:t.contentType]:void 0},genre:t.facets.genre.active?{aggregate:!0,filter:t.facets.genre.filter}:void 0,language:t.facets.language.active?{aggregate:t.facets.language.active,filter:t.facets.language.filter}:void 0,torrentFileType:t.facets.fileType.active?{aggregate:!0,filter:t.facets.fileType.filter}:void 0,torrentSource:t.facets.torrentSource.active?{aggregate:!0,filter:t.facets.torrentSource.filter}:void 0,torrentTag:t.facets.torrentTag.active?{aggregate:!0,filter:t.facets.torrentTag.filter}:void 0,videoResolution:t.facets.videoResolution.active?{aggregate:!0,filter:t.facets.videoResolution.filter}:void 0,videoSource:t.facets.videoSource.active?{aggregate:!0,filter:t.facets.videoSource.filter}:void 0}}}),x={active:!1},B=class{constructor(n){this.controlsSubject=new y(n),this.controls$=this.controlsSubject.asObservable(),this.paramsSubject=new y(He(n)),this.params$=this.paramsSubject.asObservable(),this.controls$.pipe(K(100)).subscribe(e=>{let o=this.paramsSubject.getValue(),i=He(e);JSON.stringify(o)!==JSON.stringify(i)&&this.paramsSubject.next(i)})}update(n){let e=this.controlsSubject.getValue(),o=n(e);JSON.stringify(e)!==JSON.stringify(o)&&this.controlsSubject.next(o)}selectLanguage(n){this.update(e=>s(a({},e),{language:n}))}selectContentType(n){this.update(e=>s(a({},e),{contentType:n,page:1,facets:s(a({},e.facets),{genre:H(n,Je.contentTypes)?e.facets.genre:x,videoResolution:H(n,Ue.contentTypes)?e.facets.videoResolution:x,videoSource:H(n,Ge.contentTypes)?e.facets.videoSource:x})}))}activateFacet(n){this.update(e=>s(a({},e),{facets:n.patchInput(e.facets,s(a({},n.extractInput(e.facets)),{active:!0}))}))}deactivateFacet(n){this.update(e=>{let o=n.extractInput(e.facets);return s(a({},e),{page:o.filter?1:e.page,facets:n.patchInput(e.facets,s(a({},o),{active:!1,filter:void 0}))})})}activateFilter(n,e){this.update(o=>{let i=n.extractInput(o.facets);return s(a({},o),{page:1,facets:n.patchInput(o.facets,s(a({},i),{filter:Array.from(new Set([...i.filter??[],e])).sort()}))})})}deactivateFilter(n,e){this.update(o=>{let i=n.extractInput(o.facets),r=i.filter?.filter(d=>d!==e);return s(a({},o),{page:1,facets:n.patchInput(o.facets,s(a({},i),{filter:r?.length?r:void 0}))})})}setQueryString(n){n=n||void 0,this.update(e=>{let o=e.orderBy;return n?n!==e.queryString&&(o=U):o.field==="relevance"&&(o=F),s(a({},e),{queryString:n,orderBy:o,page:n===e.queryString?e.page:1})})}selectOrderBy(n){let e={field:n,descending:J.find(o=>o.field===n)?.descending??!1};this.update(o=>s(a({},o),{orderBy:e.field!=="relevance"||o.queryString?e:F,page:1}))}toggleOrderByDirection(){this.update(n=>s(a({},n),{orderBy:s(a({},n.orderBy),{descending:!n.orderBy.descending}),page:1}))}handlePageEvent(n){this.update(e=>s(a({},e),{limit:n.pageSize,page:n.page}))}},et={key:"torrent_source",icon:"mediation",allowNull:!1,extractInput:t=>t.torrentSource,patchInput:(t,n)=>s(a({},t),{torrentSource:n}),extractAggregations:t=>t.torrentSource??[],resolveLabel:t=>t.label},tt={key:"torrent_tag",icon:"sell",allowNull:!1,extractInput:t=>t.torrentTag,patchInput:(t,n)=>s(a({},t),{torrentTag:n}),extractAggregations:t=>t.torrentTag??[],resolveLabel:t=>t.value},nt={key:"file_type",icon:"file_present",allowNull:!1,extractInput:t=>t.fileType,patchInput:(t,n)=>s(a({},t),{fileType:n}),extractAggregations:t=>t.torrentFileType??[],resolveLabel:(t,n)=>n.translate(`file_types.${t.value}`)},it={key:"language",icon:"translate",allowNull:!1,extractInput:t=>t.language,patchInput:(t,n)=>s(a({},t),{language:n}),extractAggregations:t=>t.language??[],resolveLabel:(t,n)=>n.translate(`languages.${t.value}`)},Je={key:"genre",icon:"theater_comedy",allowNull:!1,contentTypes:["movie","tv_show"],extractInput:t=>t.genre,patchInput:(t,n)=>s(a({},t),{genre:n}),extractAggregations:t=>t.genre??[],resolveLabel:t=>t.label},Ue={key:"video_resolution",icon:"aspect_ratio",allowNull:!0,contentTypes:["movie","tv_show","xxx"],extractInput:t=>t.videoResolution,patchInput:(t,n)=>s(a({},t),{videoResolution:n}),extractAggregations:t=>(t.videoResolution??[]).map(n=>s(a({},n),{value:n.value??null})),resolveLabel:t=>t.value?.slice(1)??"?"},Ge={key:"video_source",icon:"album",allowNull:!0,contentTypes:["movie","tv_show","xxx"],extractInput:t=>t.videoSource,patchInput:(t,n)=>s(a({},t),{videoSource:n}),extractAggregations:t=>(t.videoSource??[]).map(n=>s(a({},n),{value:n.value??null})),resolveLabel:t=>t.value??"?"},j=[et,tt,nt,it,Je,Ue,Ge],J=[{field:"relevance",descending:!0},{field:"published_at",descending:!0},{field:"updated_at",descending:!0},{field:"size",descending:!0},{field:"files_count",descending:!0},{field:"seeders",descending:!0},{field:"leechers",descending:!0},{field:"name",descending:!1}],F={field:"published_at",descending:!0},U={field:"relevance",descending:!0},H=(t,n)=>!n||t&&n.includes(t);var We=(t,n)=>n.key,ot=(t,n)=>n.field,Ye=(t,n)=>n.value;function rt(t,n){if(t&1&&(p(0,"small"),m(1),v(2,"number"),u()),t&2){let e=n;c(),$("",e.isEstimate?"~":"","",S(2,2,e.count),"")}}function at(t,n){if(t&1&&(p(0,"small"),m(1),v(2,"number"),u()),t&2){let e=n;c(),$("",e.isEstimate?"~":"","",S(2,2,e.count),"")}}function ct(t,n){t&1&&(p(0,"small"),m(1,"0"),u())}function st(t,n){if(t&1){let e=P();p(0,"li",5),h("click",function(){_(e);let i=l().$implicit,r=l(2);return g(r.controller.selectContentType(i.key))}),p(1,"mat-icon"),m(2),u(),m(3),C(4,at,3,4,"small"),v(5,"async"),C(6,ct,2,0,"small"),u()}if(t&2){let e,o=l().$implicit,i=l().$implicit,r=l();L(r.controls.contentType===o.key?"active":""),c(2),E(o.icon),c(),T(" ",i("content_types.plural."+o.key)," "),c(),b((e=(e=S(5,5,r.dataSource.contentTypeCounts$))==null?null:e[o.key])?4:6,e)}}function lt(t,n){if(t&1&&(C(0,st,7,7,"li",23),v(1,"async")),t&2){let e,o=n.$implicit,i=l(2);b(o.key==="null"||(e=S(1,1,i.dataSource.availableContentTypes$))!=null&&e.includes(o.key)?0:-1)}}function pt(t,n){if(t&1){let e=P();p(0,"mat-checkbox",29),h("change",function(i){let r=_(e).$implicit,d=l(3).$implicit,w=l(2);return g(i.checked?w.controller.activateFilter(d,r.value):w.controller.deactivateFilter(d,r.value))}),m(1),p(2,"small"),m(3),v(4,"number"),u()()}if(t&2){let e=n.$implicit,o=l(3).$implicit;f("checked",o.filter==null?null:o.filter.includes(e.value)),c(),T(" ",e.label," "),c(2),$("",e.isEstimate?"~":"","",S(4,4,e.count),"")}}function ut(t,n){if(t&1&&(p(0,"section",26),k(1,pt,5,6,"mat-checkbox",28,Ye),u()),t&2){let e=l(2).$implicit;c(),I(e.aggregations)}}function mt(t,n){if(t&1){let e=P();p(0,"mat-checkbox",31),h("change",function(){let i=_(e).$implicit,r=l(3).$implicit,d=l(2);return g(d.controller.activateFilter(r,i.value))}),m(1),p(2,"small"),m(3),v(4,"number"),u()()}if(t&2){let e=n.$implicit;c(),T(" ",e.label," "),c(2),$("",e.isEstimate?"~":"","",S(4,3,e.count),"")}}function dt(t,n){if(t&1&&(p(0,"section",27),k(1,mt,5,5,"mat-checkbox",30,Ye),u()),t&2){let e=l(2).$implicit;c(),I(e.aggregations)}}function _t(t,n){if(t&1){let e=P();p(0,"mat-expansion-panel",25),h("opened",function(){_(e);let i=l().$implicit,r=l(2);return g(r.controller.activateFacet(i))})("closed",function(){_(e);let i=l().$implicit,r=l(2);return g(r.controller.deactivateFacet(i))}),p(1,"mat-expansion-panel-header")(2,"mat-panel-title")(3,"mat-icon"),m(4),u(),m(5),u()(),C(6,ut,3,0,"section",26)(7,dt,3,0,"section",27),u()}if(t&2){let e=l().$implicit,o=l().$implicit;f("expanded",e.active),c(4),E(e.icon),c(),T(" ",o("facets."+e.key)," "),c(),b(e.filter!=null&&e.filter.length?6:7)}}function gt(t,n){if(t&1&&C(0,_t,8,4,"mat-expansion-panel",24),t&2){let e=n.$implicit;b(e.relevant?0:-1)}}function ft(t,n){if(t&1){let e=P();p(0,"button",17),h("click",function(){_(e);let i=l(2);return i.queryString.reset(),g(i.controller.setQueryString(null))}),p(1,"mat-icon"),m(2,"close"),u()()}if(t&2){let e=l().$implicit;f("matTooltip",e("torrents.clear_search"))}}function ht(t,n){if(t&1&&(p(0,"mat-option",32),m(1),u()),t&2){let e=l().$implicit,o=l().$implicit;f("value",e.field),c(),T(" ",o("torrents.ordering."+e.field)," ")}}function xt(t,n){if(t&1&&C(0,ht,2,2,"mat-option",32),t&2){let e=n.$implicit,o=l(2);b(e.field!="relevance"||o.queryString.value?0:-1)}}function Ct(t,n){if(t&1){let e=P();te(0),p(1,"mat-drawer-container",2)(2,"mat-drawer",3,0)(4,"mat-expansion-panel",4)(5,"mat-expansion-panel-header")(6,"mat-panel-title")(7,"mat-icon"),m(8,"interests"),u(),m(9),u()(),p(10,"section")(11,"nav")(12,"ul")(13,"li",5),h("click",function(){_(e);let i=l();return g(i.controller.selectContentType(null))}),p(14,"mat-icon",6),m(15,"emergency"),u(),m(16),C(17,rt,3,4,"small"),v(18,"async"),u(),k(19,lt,2,3,null,null,We),u()()()(),k(21,gt,1,1,null,null,We),v(23,"async"),u(),p(24,"mat-drawer-content")(25,"div",7)(26,"div",8)(27,"button",9),h("click",function(){_(e);let i=R(3);return g(i.toggle())}),p(28,"mat-icon",10),m(29),u()()(),p(30,"div",11)(31,"mat-form-field",12)(32,"input",13),h("keyup.enter",function(){_(e);let i=l();return g(i.controller.setQueryString(i.queryString.value))}),u(),C(33,ft,3,1,"button",14),u()(),p(34,"div",15)(35,"mat-form-field",12)(36,"mat-label"),m(37),u(),p(38,"mat-select",16),h("valueChange",function(i){_(e);let r=l();return g(r.controller.selectOrderBy(i))}),k(39,xt,1,1,null,null,ot),u()(),p(41,"button",17),h("click",function(){_(e);let i=l();return g(i.controller.toggleOrderByDirection())}),p(42,"mat-icon"),m(43),u()()(),p(44,"div",18)(45,"button",19),h("click",function(){_(e);let i=l();return g(i.dataSource.refresh())}),p(46,"mat-icon"),m(47,"sync"),u()()()(),A(48,"mat-divider"),p(49,"app-torrents-bulk-actions",20),h("updated",function(){_(e);let i=l();return g(i.dataSource.refresh())}),u(),A(50,"mat-divider"),p(51,"app-torrents-table",21),h("updated",function(){_(e);let i=l();return g(i.dataSource.refresh())}),u(),p(52,"app-paginator",22),h("paging",function(i){_(e);let r=l();return g(r.controller.handlePageEvent(i))}),u()()(),ne()}if(t&2){let e,o=n.$implicit,i=R(3),r=l();c(2),f("mode",r.breakpoints.sizeAtLeast("Medium")?"side":"over")("opened",r.breakpoints.sizeAtLeast("Medium")),ee("role",r.breakpoints.sizeAtLeast("Medium")?"navigation":"dialog"),c(2),f("expanded",r.breakpoints.sizeAtLeast("Medium")),c(5),T(" ",o("facets.content_type")," "),c(4),L(r.controls.contentType===null?"active":""),c(3),T("",o("content_types.plural.all")," "),c(),b((e=S(18,29,r.dataSource.overallTotalCount$))?17:-1,e),c(2),I(r.contentTypes),c(2),I(S(23,31,r.facets$)),c(6),f("matTooltip",o("torrents.toggle_drawer")),c(2),E(i.opened?"arrow_circle_left":"arrow_circle_right"),c(3),f("placeholder",o("torrents.search"))("formControl",r.queryString),c(),b(r.queryString.value?33:-1),c(4),E(o("torrents.order_by")),c(),f("value",r.controls.orderBy.field),c(),I(r.orderByOptions),c(2),f("matTooltip",o("torrents.order_direction_toggle")),c(2),E(r.controls.orderBy.descending?"arrow_downward":"arrow_upward"),c(2),f("matTooltip",o("torrents.refresh")),c(4),f("selectedItems$",r.selectedItems$),c(2),f("dataSource",r.dataSource)("displayedColumns",r.breakpoints.sizeAtLeast("Medium")?r.allColumns:r.compactColumns)("selection",r.selection),c(),f("page",r.controls.page)("pageSize",r.controls.limit)("pageLength",r.dataSource.result.items.length)("totalLength",r.dataSource.result.totalCount)("totalIsEstimate",r.dataSource.result.totalCountIsEstimate)("hasNextPage",r.dataSource.result.hasNextPage)}}var _n=(()=>{class t{constructor(){this.route=O(ae),this.router=O(ce),this.apollo=O(ge),this.errorsService=O(Ce),this.transloco=O(se),this.breakpoints=O(ye),this.contentTypes=Re,this.orderByOptions=J,this.allColumns=Ne,this.compactColumns=ze,this.queryString=new me(""),this.result=Q,this.selection=new _e(!0,[]),this.selectedItemsSubject=new y([]),this.selectedItems$=this.selectedItemsSubject.asObservable(),this.subscriptions=Array(),this.controls=s(a({},yt),{language:this.transloco.getActiveLang()}),this.controller=new B(this.controls),this.dataSource=new V(this.apollo,this.errorsService,this.controller.params$),this.subscriptions.push(this.controller.controls$.subscribe(e=>{this.controls=e})),this.facets$=this.controller.controls$.pipe(Y(this.dataSource.result$),M(([e,o])=>j.map(i=>s(a(a({},i),i.extractInput(e.facets)),{relevant:!i.contentTypes||!!(e.contentType&&e.contentType!=="null"&&i.contentTypes.includes(e.contentType)),aggregations:i.extractAggregations(o.aggregations).map(r=>s(a({},r),{label:i.resolveLabel(r,this.transloco)}))})))),this.subscriptions.push(this.dataSource.result$.subscribe(e=>{this.result=e;let o=new Set(e.items.map(({infoHash:i})=>i));this.selection.deselect(...this.selection.selected.filter(i=>!o.has(i)))}))}ngOnInit(){this.subscriptions.push(this.route.queryParams.subscribe(e=>{let o=N(e,"query");this.queryString.setValue(o??null),this.controller.update(i=>{let r=D(e,"facets"),d=i.orderBy;return o?o!==i.queryString&&(d=U):d.field==="relevance"&&(d=F),s(a({},i),{queryString:o,orderBy:d,contentType:vt(e,"content_type"),limit:z(e,"limit")??i.limit,page:z(e,"page")??i.page,facets:j.reduce((w,q)=>{let Xe=r?.includes(q.key)??!1,Ze=D(e,q.key);return q.patchInput(w,{active:Xe,filter:Ze})},i.facets)})})}),this.controller.controls$.subscribe(e=>{let o=e.page,i=e.limit;o===1&&(o=void 0),i===Ke&&(i=void 0),this.router.navigate([],{relativeTo:this.route,queryParams:a({query:e.queryString?encodeURIComponent(e.queryString):void 0,page:o,limit:i,content_type:e.contentType},St(e.facets)),queryParamsHandling:"merge"})}),this.selection.changed.subscribe(e=>{let o=new Set(e.source.selected);this.selectedItemsSubject.next(this.result.items.filter(i=>o.has(i.infoHash)))}))}ngOnDestroy(){this.subscriptions.forEach(e=>e.unsubscribe()),this.subscriptions=new Array}static{this.\u0275fac=function(o){return new(o||t)}}static{this.\u0275cmp=Z({type:t,selectors:[["app-torrents-search"]],standalone:!0,features:[ie],decls:1,vars:0,consts:[["drawer",""],[4,"transloco"],[1,"drawer-container"],[1,"drawer",3,"mode","opened"],[1,"panel-content-type",3,"expanded"],[3,"click"],["fontSet","material-icons"],[1,"search-form"],[1,"form-field-container","button-container","button-container-toggle-drawer"],["type","button","mat-icon-button","",1,"button-toggle-drawer",3,"click","matTooltip"],["aria-label","Side nav toggle icon","fontSet","material-icons"],[1,"form-field-container","form-field-container-search-query"],["subscriptSizing","dynamic"],["matInput","","autocapitalize","none",3,"keyup.enter","placeholder","formControl"],["mat-icon-button","",3,"matTooltip"],[1,"form-field-container","form-field-container-order-by"],[3,"valueChange","value"],["mat-icon-button","",3,"click","matTooltip"],[1,"form-field-container","button-container","button-container-refresh"],["mat-mini-fab","","color","primary",3,"click","matTooltip"],[3,"updated","selectedItems$"],[3,"updated","dataSource","displayedColumns","selection"],[3,"paging","page","pageSize","pageLength","totalLength","totalIsEstimate","hasNextPage"],[3,"class"],[3,"expanded"],[3,"opened","closed","expanded"],[1,"filtered"],[1,"unfiltered"],[3,"checked"],[3,"change","checked"],["checked","true"],["checked","true",3,"change"],[3,"value"]],template:function(o,i){o&1&&C(0,Ct,53,33,"ng-container",1)},dependencies:[qe,fe,xe,he,be,Te,we,Me,Oe,Se,ve,Pe,ke,Ee,Fe,Ve,$e,Ie,pe,ue,de,le,oe,re,je,Le,Qe,De],styles:[".mat-expansion-panel[_ngcontent-%COMP%]{margin-top:14px;margin-right:14px}.mat-expansion-panel[_ngcontent-%COMP%]   section[_ngcontent-%COMP%]{margin-left:-10px}.mat-expansion-panel.panel-content-type[_ngcontent-%COMP%]{margin-top:20px}.mat-expansion-panel.panel-content-type[_ngcontent-%COMP%]   section[_ngcontent-%COMP%]{margin-left:0}.mat-expansion-panel[_ngcontent-%COMP%]   ul[_ngcontent-%COMP%]{list-style:none;padding-left:0;margin:0}.mat-expansion-panel[_ngcontent-%COMP%]   mat-panel-title[_ngcontent-%COMP%], .mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%]{position:relative;line-height:40px;padding-left:40px}.mat-expansion-panel[_ngcontent-%COMP%]   mat-panel-title[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%], .mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%]{position:absolute;left:0;top:8px}.mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%]{cursor:pointer}.mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%]{top:6px}.mat-expansion-panel[_ngcontent-%COMP%]   li[_ngcontent-%COMP%]   small[_ngcontent-%COMP%]{float:right;font-size:.8rem}.mat-expansion-panel[_ngcontent-%COMP%]   mat-checkbox[_ngcontent-%COMP%]{display:block}.mat-expansion-panel[_ngcontent-%COMP%]   mat-checkbox[_ngcontent-%COMP%]     label{min-width:220px}.mat-expansion-panel[_ngcontent-%COMP%]   mat-checkbox[_ngcontent-%COMP%]   small[_ngcontent-%COMP%]{margin-left:10px;position:absolute;right:0}.search-form[_ngcontent-%COMP%]{padding-top:20px;padding-bottom:10px;position:relative;clear:both;display:flex;flex-wrap:wrap}.search-form[_ngcontent-%COMP%]   .form-field-container[_ngcontent-%COMP%]{display:inline-flex;flex-direction:column;position:relative;margin-left:20px;padding-bottom:20px}.search-form[_ngcontent-%COMP%]   .form-field-container[_ngcontent-%COMP%]   button[_ngcontent-%COMP%]{top:8px}.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-order-by[_ngcontent-%COMP%]{padding-right:40px}.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-order-by[_ngcontent-%COMP%]   button[_ngcontent-%COMP%]{position:absolute;right:0}.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-search-query[_ngcontent-%COMP%]{width:300px}.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-search-query[_ngcontent-%COMP%]   button[_ngcontent-%COMP%]{position:absolute;right:0}.search-form[_ngcontent-%COMP%]   .form-field-container.form-field-container-search-query[_ngcontent-%COMP%]     .mat-mdc-form-field-infix{padding-right:50px}.search-form[_ngcontent-%COMP%]   .button-container-toggle-direction[_ngcontent-%COMP%]{margin-left:4px}app-paginator[_ngcontent-%COMP%]{float:right;padding-top:14px;padding-bottom:20px}"],changeDetection:0})}}return t})(),Ke=20,yt={language:"en",page:1,limit:Ke,contentType:null,orderBy:F,facets:{genre:x,language:x,fileType:x,torrentSource:x,torrentTag:x,videoResolution:x,videoSource:x}},vt=(t,n)=>{let e=N(t,n);return e&&e in Ae?e:null},St=t=>{let[n,e]=j.reduce((o,i)=>{let r=i.extractInput(t);return r.active?[[...o[0],i.key],r.filter?s(a({},o[1]),{[i.key]:r.filter}):o[1]]:o},[[],{}]);return a({facets:n.length?n.join(","):void 0},Object.fromEntries(Object.entries(e).map(([o,i])=>[o,encodeURIComponent(i.join(","))])))};export{_n as TorrentsSearchComponent};
