import{a as je,b as $e,c as j}from"./chunk-TJLO75RV.js";import{e as Re}from"./chunk-Y7K23DTG.js";import{l as Ae}from"./chunk-42PJPEMD.js";import{a as ve}from"./chunk-ASLGZ7DJ.js";import{a as Ne}from"./chunk-MSAOOVCY.js";import{a as De}from"./chunk-DSEDLZDW.js";import{G as ke,J as Te,M as Se,P as J,Q as ye,R as K,U as q,Ua as Y,a as Ce,b as G,i as ge,n as xe,na as we,oa as Oe,qa as Q,r as Me,ra as Pe,va as Ee,xa as Be,ya as Ve,za as Ie}from"./chunk-VAEZNV34.js";import"./chunk-6XXA7HXI.js";import"./chunk-CMNWCZJM.js";import{$b as p,B as ue,Cb as F,Ea as _,Fa as f,Hb as l,Kb as N,N as pe,Nb as D,Ob as E,Pb as B,Q as te,Qb as i,Rb as c,Sb as P,Tb as z,Ub as H,Wb as re,Zb as h,a as v,ad as be,b as k,ka as de,kb as a,kc as u,l as I,lc as T,mc as _e,o as me,pc as U,qa as O,qc as fe,sc as he,ua as L,vc as ne,wc as ie}from"./chunk-Z3WUIYN5.js";var R=(n,s)=>{let{bucketParams:e,earliestBucket:r}=We(n,s),t=Object.entries(s.torrent.metrics.buckets.reduce((b,x)=>{if(x.source!==(n.source??x.source))return b;let g=y(x.bucket,e);if(r&&r.index>g.index&&(g=void 0),!g)return b;let d=b[x.source]??[];return k(v({},b),{[x.source]:{created:x.updated?d.created:k(v({},d.created),{[g.key]:{count:x.count+(d.created?.[g.key]?.count??0),startTime:g.start}}),updated:x.updated?k(v({},d.updated),{[g.key]:{count:x.count+(d.updated?.[g.key]?.count??0),startTime:g.start}}):d.updated}})},{})).map(([b,x])=>{let g;if(Object.keys(x).length){let d=Array(),M=Le(Array("created","updated").flatMap(S=>{let A=Le(Object.entries(x[S]??{}).filter(([,ee])=>ee?.count).sort(([ee],[Ye])=>parseInt(ee)<parseInt(Ye)?1:-1)),w=Object.keys(A);if(!w.length)return[];let ce=parseInt(w[0]),le=parseInt(w[w.length-1]);return d.push(ce,le),[[S,{earliestBucket:ce,latestBucket:le,entries:A}]]}));d.sort(),g={bucketDuration:e.duration,earliestBucket:d[0],latestBucket:d[d.length-1],eventBuckets:M}}return{source:b,events:g,isEmpty:!g?.eventBuckets}}),m,C=t.flatMap(b=>b.events?[b.events.earliestBucket]:[]).sort()[0],o=t.flatMap(b=>b.events?[b.events.latestBucket]:[]).sort().reverse()[0];return C&&o&&(m={earliestBucket:C,latestBucket:o}),{params:k(v({},n),{buckets:e}),sourceSummaries:t,bucketSpan:m,availableSources:s.torrent.listSources.sources.map(b=>({key:b.key,name:b.name}))}},Le=n=>Object.fromEntries(n),We=(n,s)=>{let e=n.buckets.duration==="AUTO"?"hour":n.buckets.duration,r=n.buckets.multiplier==="AUTO"?1:n.buckets.multiplier,t=n.buckets.timeframe,m=new Date,C=y(m,{duration:e,multiplier:r}),o=y(m.getTime()-1e3*V[t],{duration:e,multiplier:r}),b=[o,...s.torrent.metrics.buckets.flatMap(d=>[y(d.bucket,{duration:e,multiplier:r})]),C].filter(d=>d.index>=o.index).sort((d,M)=>d.index-M.index),x=b[0],g=b[b.length-1];if(n.buckets.multiplier==="AUTO"){let M=g.index-x.index;r=Math.min(60,Math.max(Math.floor(M/(20*5))*5,1))}return{bucketParams:{duration:e,multiplier:r,timeframe:t},earliestBucket:y(m.getTime()-1e3*V[t],{duration:e,multiplier:r}),latestBucket:y(Math.max(m.getTime(),g.start.getTime()),{duration:e,multiplier:r})}},y=(n,s)=>{let e=new Date(n),r=1e3*W[s.duration]*s.multiplier,t=Math.floor(e.getTime()/r);return{key:`${t}`,index:t,start:new Date(t*r)}};var oe={duration:"minute",multiplier:1,timeframe:"hours_1"},Fe=["day","hour","minute"],W={minute:60,hour:60*60,day:60*60*24},ae={buckets:oe,autoRefresh:"off"},se={torrent:{metrics:{buckets:[]},listSources:{sources:[{key:"dht",name:"DHT"}]}}},X=["created","updated"],ze=["minutes_15","minutes_30","hours_1","hours_6","hours_12","days_1","weeks_1"],V={minutes_15:60*15,minutes_30:60*30,hours_1:60*60,hours_6:60*60*6,hours_12:60*60*12,days_1:60*60*24,weeks_1:60*60*24*7},He=["off","seconds_10","seconds_30","minutes_1","minutes_5"],Ue={off:null,seconds_10:10,seconds_30:30,minutes_1:60,minutes_5:60*5},Ge=R(ae,se);var Z=class{constructor(s,e=ae,r){this.apollo=s,this.errorsService=r,this.rawResultSubject=new I(se),this.resultSubject=new I(Ge),this.result$=this.resultSubject.asObservable(),this.loadingSubject=new I(!1),this.paramsSubject=new I(e),this.params$=this.paramsSubject.asObservable(),this.variablesSubject=new I(Je(e)),this.paramsSubject.pipe(te(50)).subscribe(t=>{let m=this.variablesSubject.getValue(),C=Je(t);JSON.stringify(m)!==JSON.stringify(C)?this.variablesSubject.next(C):this.resultSubject.next(R(t,this.rawResultSubject.getValue()))}),this.variablesSubject.pipe(te(50)).subscribe(t=>this.request(t)),this.rawResultSubject.subscribe(t=>{let m=this.paramsSubject.getValue();this.resultSubject.next(R(m,t)),this.setInterval(m.autoRefresh)})}setInterval(s){clearTimeout(this.refreshTimeout);let e=Ue[s??this.params.autoRefresh];e&&(this.refreshTimeout=setTimeout(()=>{this.refresh()},e*1e3))}get params(){return this.paramsSubject.getValue()}get bucketDuration(){let s=this.params.buckets.duration;return s==="AUTO"?"hour":s}get bucketMultiplier(){return this.resultSubject.getValue().params.buckets.multiplier??this.params.buckets.multiplier}get loading(){return this.loadingSubject.getValue()}setTimeframe(s){this.updateParams(e=>k(v({},e),{buckets:k(v({},e.buckets),{timeframe:s})}))}setSource(s){this.updateParams(e=>k(v({},e),{source:s??void 0}))}setBucketDuration(s,e){this.updateParams(r=>k(v({},r),{buckets:k(v({},r.buckets),{duration:s,multiplier:e??"AUTO"})}))}setBucketMultiplier(s){this.updateParams(e=>k(v({},e),{buckets:k(v({},e.buckets),{multiplier:s})}))}setEvent(s){this.updateParams(e=>k(v({},e),{event:s??void 0}))}setAutoRefreshInterval(s){this.updateParams(e=>k(v({},e),{autoRefresh:s}))}updateParams(s){this.paramsSubject.next(s(this.params))}refresh(){this.variablesSubject.next(this.variablesSubject.getValue())}request(s){return clearTimeout(this.refreshTimeout),this.loadingSubject.next(!0),this.apollo.query({query:ke,variables:s,fetchPolicy:"no-cache"}).pipe(ue(e=>{e&&(this.loadingSubject.next(!1),this.rawResultSubject.next(e.data))})).pipe(pe(e=>(this.errorsService.addError(`Failed to load torrent metrics: ${e.message}`),this.loadingSubject.next(!1),this.setInterval(),me))).subscribe()}},Je=n=>({input:{bucketDuration:n.buckets.duration==="AUTO"?"hour":n.buckets.duration,sources:n.source?[n.source]:void 0,startTime:new Date(new Date().getTime()-1e3*V[n.buckets.timeframe]).toISOString()}});var $={created:"primary",updated:"secondary"},qe=(()=>{class n{constructor(){this.themeInfo=O(Re),this.transloco=O(Ce)}create(e){let{colors:r}=this.themeInfo.info,t=Array(),m=[];if(e){let C=e.sourceSummaries.filter(d=>!d.isEmpty),o=Array.from(new Set(C.flatMap(d=>d.events?[d.events.earliestBucket,d.events.latestBucket]:[]))).sort(),b=new Date,x=Math.min(o[0],y(b.getTime()-1e3*V[e.params.buckets.timeframe],e.params.buckets).index),g=Math.max(o[o.length-1],y(b,e.params.buckets).index);if(o.length){for(let M=x;M<=g;M++)t.push(this.formatBucketKey(e.params.buckets,M));let d=X.filter(M=>(e.params.event??M)===M);for(let M of C)for(let S of d){let A=Array();for(let w=x;w<=g;w++)A.push(M.events?.eventBuckets?.[S]?.entries?.[`${w}`]?.count??0);m.push({yAxisID:"yCount",label:[M.source,this.transloco.translate("dashboard.queues."+S)].join("/"),data:A,borderColor:r[j($[S],50)],pointBackgroundColor:r[j($[S],20)],pointBorderColor:r[j($[S],80)],pointHoverBackgroundColor:r[j($[S],40)],pointHoverBorderColor:r[j($[S],60)]})}}}return{type:"line",options:{animation:!1,elements:{line:{tension:.5}},scales:{yCount:{position:"left",ticks:{callback:C=>parseInt(C).toLocaleString(this.transloco.getActiveLang())}}},plugins:{legend:{display:!0},decimation:{enabled:!0}}},data:{labels:t,datasets:m}}}formatBucketKey(e,r){let t;switch(e.duration){case"day":t="d LLL";break;case"hour":t="d LLL H:00";break;case"minute":t="H:mm";break}return je(1e3*W[e.duration]*e.multiplier*r,t,{locale:Ae(this.transloco.getActiveLang())})}static{this.\u0275fac=function(r){return new(r||n)}}static{this.\u0275prov=de({token:n,factory:n.\u0275fac,providedIn:"root"})}}return n})();var tt=(n,s)=>s.key,rt=()=>["dht"];function nt(n,s){if(n&1&&(i(0,"mat-option",6),u(1),c()),n&2){let e=s.$implicit,r=p().$implicit;l("value",e),a(),T(r("dashboard.interval."+e))}}function it(n,s){if(n&1&&(i(0,"mat-option",6),u(1),c()),n&2){let e=s.$implicit,r=p().$implicit;l("value",e),a(),T(r("dashboard.interval."+e+"s"))}}function ot(n,s){if(n&1&&(i(0,"mat-option",6),u(1),c()),n&2){let e=s.$implicit;l("value",e.key),a(),T(e.name)}}function at(n,s){if(n&1){let e=re();i(0,"button",16),h("click",function(){let t=_(e).$implicit,m=p(2);return f(m.torrentMetricsController.params.source===t||m.torrentMetricsController.setSource(t))}),i(1,"mat-icon"),u(2),c()()}if(n&2){let e=s.$implicit,r=p(2);N(r.torrentMetricsController.params.source===e?"selected":"deselected"),l("matTooltip",e),a(2),T(r.torrentMetricsController.params.source===e?"radio_button_checked":"radio_button_unchecked")}}function st(n,s){if(n&1&&(i(0,"mat-option",6),u(1),c()),n&2){let e=s.$implicit,r=p().$implicit;l("value",e),a(),T(r("dashboard.event."+e))}}function ct(n,s){if(n&1&&(i(0,"mat-option",6),u(1),c()),n&2){let e=s.$implicit,r=p().$implicit;l("value",e),a(),T(r("dashboard.interval."+e))}}function lt(n,s){if(n&1){let e=re();z(0),i(1,"mat-card")(2,"mat-card-content")(3,"mat-grid-list",1)(4,"mat-grid-tile",2)(5,"mat-card",3)(6,"mat-card-header")(7,"mat-card-title")(8,"h4"),u(9),c()()(),i(10,"mat-card-content")(11,"mat-form-field",4)(12,"mat-select",5),h("valueChange",function(t){_(e);let m=p();return f(m.torrentMetricsController.setTimeframe(t))}),E(13,nt,2,2,"mat-option",6,D),c()(),i(15,"div",7)(16,"button",8),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setTimeframe(t.timeframeNames[0]))}),i(17,"mat-icon"),u(18,"first_page"),c()(),i(19,"button",8),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setTimeframe(t.timeframeNames[t.timeframeNames.indexOf(t.torrentMetricsController.params.buckets.timeframe)-1]))}),i(20,"mat-icon"),u(21,"navigate_before"),c()(),i(22,"button",8),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setTimeframe(t.timeframeNames[t.timeframeNames.indexOf(t.torrentMetricsController.params.buckets.timeframe)+1]))}),i(23,"mat-icon"),u(24,"navigate_next"),c()(),i(25,"button",8),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setTimeframe(t.timeframeNames[t.timeframeNames.length-1]))}),i(26,"mat-icon"),u(27,"last_page"),c()()()()()(),i(28,"mat-grid-tile",2)(29,"mat-card",9)(30,"mat-card-header")(31,"mat-card-title")(32,"h4"),u(33),c()()(),i(34,"mat-card-content")(35,"mat-form-field",10)(36,"input",11),ne(37,"async"),h("change",function(t){_(e);let m=p();return f(m.handleMultiplierEvent(t))}),c()(),i(38,"mat-form-field",12)(39,"mat-select",5),h("valueChange",function(t){_(e);let m=p();return f(m.torrentMetricsController.setBucketDuration(t))}),E(40,it,2,2,"mat-option",6,D),c()(),i(42,"div",7)(43,"button",8),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setBucketMultiplier(t.torrentMetricsController.bucketMultiplier-1))}),i(44,"mat-icon"),u(45,"remove"),c()(),i(46,"button",13),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setBucketMultiplier(t.torrentMetricsController.bucketMultiplier+1))}),i(47,"mat-icon"),u(48,"add"),c()(),i(49,"button",8),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setBucketDuration(t.resolutionNames[0]))}),i(50,"mat-icon"),u(51,"first_page"),c()(),i(52,"button",8),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setBucketDuration(t.resolutionNames[t.resolutionNames.indexOf(t.torrentMetricsController.bucketDuration)-1]))}),i(53,"mat-icon"),u(54,"navigate_before"),c()(),i(55,"button",8),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setBucketDuration(t.resolutionNames[t.resolutionNames.indexOf(t.torrentMetricsController.bucketDuration)+1]))}),i(56,"mat-icon"),u(57,"navigate_next"),c()(),i(58,"button",8),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setBucketDuration(t.resolutionNames[t.resolutionNames.length-1]))}),i(59,"mat-icon"),u(60,"last_page"),c()()()()()(),i(61,"mat-grid-tile",2)(62,"mat-card")(63,"mat-card-header")(64,"mat-card-title")(65,"h4"),u(66),c()()(),i(67,"mat-card-content")(68,"mat-form-field",4)(69,"mat-select",5),h("valueChange",function(t){_(e);let m=p();return f(m.torrentMetricsController.setSource(t==="_all"?null:t))}),i(70,"mat-option",14),u(71,"All"),c(),E(72,ot,2,2,"mat-option",6,tt),ne(74,"async"),c()(),i(75,"div",15)(76,"button",16),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setSource(null))}),i(77,"mat-icon",17),u(78,"workspaces"),c()(),E(79,at,3,4,"button",18,D),c()()()(),i(81,"mat-grid-tile",2)(82,"mat-card")(83,"mat-card-header")(84,"mat-card-title")(85,"h4"),u(86),c()()(),i(87,"mat-card-content")(88,"mat-form-field",4)(89,"mat-select",5),h("valueChange",function(t){_(e);let m=p();return f(m.torrentMetricsController.setEvent(t==="_all"?null:t))}),i(90,"mat-option",14),u(91,"All"),c(),E(92,st,2,2,"mat-option",6,D),c()(),i(94,"div",15)(95,"button",16),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.setEvent(null))}),i(96,"mat-icon",17),u(97,"radio_button_checked"),c()(),i(98,"button",16),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.params.event==="created"||t.torrentMetricsController.setEvent("created"))}),i(99,"mat-icon"),u(100,"add_circle"),c()(),i(101,"button",16),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.params.event==="updated"||t.torrentMetricsController.setEvent("updated"))}),i(102,"mat-icon"),u(103,"check_circle"),c()()()()()(),i(104,"mat-grid-tile",2)(105,"mat-card",19)(106,"mat-card-header")(107,"mat-card-title")(108,"h4"),u(109),c()()(),i(110,"mat-card-content")(111,"mat-form-field",4)(112,"mat-select",5),h("valueChange",function(t){_(e);let m=p();return f(m.torrentMetricsController.setAutoRefreshInterval(t))}),E(113,ct,2,2,"mat-option",6,D),c()(),i(115,"div",15)(116,"button",16),h("click",function(){_(e);let t=p();return f(t.torrentMetricsController.refresh())}),i(117,"mat-icon"),u(118,"sync"),c()()()()()()(),i(119,"div",20),P(120,"mat-progress-bar",21),c(),i(121,"mat-grid-list",1)(122,"mat-grid-tile",2),P(123,"app-chart",22),c(),P(124,"mat-grid-tile",2),c()()(),H()}if(n&2){let e,r,t,m,C=s.$implicit,o=p();a(3),l("cols",o.breakpoints.sizeAtLeast("Large")?5:o.breakpoints.sizeAtLeast("Medium")?3:o.breakpoints.sizeAtLeast("Small")?2:1),a(),l("colspan",1)("rowspan",2),a(5),T(C("dashboard.metrics.timeframe")),a(3),l("value",o.torrentMetricsController.params.buckets.timeframe),a(),B(o.timeframeNames),a(3),l("disabled",o.timeframeNames.indexOf(o.torrentMetricsController.params.buckets.timeframe)<=0),a(3),l("disabled",o.timeframeNames.indexOf(o.torrentMetricsController.params.buckets.timeframe)<=0),a(3),l("disabled",o.timeframeNames.indexOf(o.torrentMetricsController.params.buckets.timeframe)>=o.timeframeNames.length-1),a(3),l("disabled",o.timeframeNames.indexOf(o.torrentMetricsController.params.buckets.timeframe)>=o.timeframeNames.length-1),a(3),l("colspan",1)("rowspan",2),a(5),_e(" ",C("dashboard.metrics.resolution")," "),a(3),l("placeholder",(e=(e=ie(37,57,o.torrentMetricsController.result$))==null||e.params==null||e.params.buckets==null||e.params.buckets.multiplier==null?null:e.params.buckets.multiplier.toString())!==null&&e!==void 0?e:"")("value",o.torrentMetricsController.params.buckets.multiplier),a(3),l("value",o.torrentMetricsController.bucketDuration),a(),B(o.resolutionNames),a(3),l("disabled",o.torrentMetricsController.bucketMultiplier===1),a(6),l("disabled",o.resolutionNames.indexOf(o.torrentMetricsController.bucketDuration)<=0),a(3),l("disabled",o.resolutionNames.indexOf(o.torrentMetricsController.bucketDuration)<=0),a(3),l("disabled",o.resolutionNames.indexOf(o.torrentMetricsController.bucketDuration)>=o.resolutionNames.length-1),a(3),l("disabled",o.resolutionNames.indexOf(o.torrentMetricsController.bucketDuration)>=o.resolutionNames.length-1),a(3),l("colspan",1)("rowspan",2),a(5),T(C("torrents.source")),a(3),l("value",(r=o.torrentMetricsController.params.source)!==null&&r!==void 0?r:"_all"),a(3),B((t=ie(74,59,o.torrentMetricsController.result$))==null?null:t.availableSources),a(4),N(o.torrentMetricsController.params.source?"deselected":"selected"),l("matTooltip","all"),a(3),B(fe(61,rt)),a(2),l("colspan",1)("rowspan",2),a(5),T(C("dashboard.metrics.event")),a(3),l("value",(m=o.torrentMetricsController.params.event)!==null&&m!==void 0?m:"_all"),a(3),B(o.eventNames),a(3),N(o.torrentMetricsController.params.event?"deselected":"selected"),l("matTooltip","all"),a(3),N(o.torrentMetricsController.params.event==="created"?"selected":"deselected"),l("matTooltip","created"),a(3),N(o.torrentMetricsController.params.event==="updated"?"selected":"deselected"),l("matTooltip","updated"),a(3),l("colspan",1)("rowspan",2),a(5),T(C("general.refresh")),a(3),l("value",o.torrentMetricsController.params.autoRefresh),a(),B(o.autoRefreshIntervalNames),a(3),l("matTooltip","Refresh"),a(4),l("mode",o.torrentMetricsController.loading?"indeterminate":"determinate")("value",0),a(),l("cols",o.breakpoints.sizeAtLeast("Large")?2:1),a(),l("colspan",1)("rowspan",6),a(),l("adapter",o.timeline)("$data",o.torrentMetricsController.result$)("height",400)("width",550)("title",C("dashboard.metrics.throughput")),a(),l("colspan",1)("rowspan",5)}}var Qe=(()=>{class n{constructor(){this.breakpoints=O(Ne),this.apollo=O(ge),this.torrentMetricsController=new Z(this.apollo,{buckets:oe,autoRefresh:"seconds_30"},O(ve)),this.timeline=O(qe),this.resolutionNames=Fe,this.timeframeNames=ze,this.autoRefreshIntervalNames=He,this.eventNames=X}ngOnDestroy(){this.torrentMetricsController.setAutoRefreshInterval("off")}handleMultiplierEvent(e){let r=e.currentTarget.value;this.torrentMetricsController.setBucketMultiplier(/^\d+$/.test(r)?parseInt(r):"AUTO")}static{this.\u0275fac=function(r){return new(r||n)}}static{this.\u0275cmp=L({type:n,selectors:[["app-torrent-metrics"]],standalone:!0,features:[U],decls:1,vars:0,consts:[[4,"transloco"],["rowHeight","100px",3,"cols"],[3,"colspan","rowspan"],[1,"form-timeframe"],["subscriptSizing","dynamic"],[3,"valueChange","value"],[3,"value"],[1,"paginator","actions"],["mat-icon-button","",3,"click","disabled"],[1,"form-resolution"],["subscriptSizing","dynamic",1,"form-input-multiplier"],["type","number","matInput","","min","1","step","1",3,"change","placeholder","value"],["subscriptSizing","dynamic",1,"form-select-duration"],["mat-icon-button","",3,"click"],["value","_all"],[1,"actions"],["mat-icon-button","",3,"click","matTooltip"],["fontSet","material-icons"],["mat-icon-button","",3,"class","matTooltip"],[1,"form-refresh"],[1,"progress-bar-container"],[3,"mode","value"],[3,"adapter","$data","height","width","title"]],template:function(r,t){r&1&&F(0,lt,125,62,"ng-container",0)},dependencies:[Y,xe,Me,J,K,q,ye,Se,Oe,we,Q,Pe,Ee,Ie,Ve,G,be,$e,Te],styles:[".actions[_ngcontent-%COMP%]{width:210px;padding-top:12px;--mdc-icon-button-state-layer-size: 32px}.actions[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%]{font-size:22px}.actions[_ngcontent-%COMP%]   button[_ngcontent-%COMP%]{margin-right:0}.progress-bar-container[_ngcontent-%COMP%]{width:100%;height:10px}mat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]{width:100%}mat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]   mat-card-content[_ngcontent-%COMP%]{min-width:190px}mat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]   h4[_ngcontent-%COMP%]{margin-bottom:16px;font-size:18px}mat-form-field[_ngcontent-%COMP%]{width:186px}.form-resolution[_ngcontent-%COMP%]   .actions[_ngcontent-%COMP%]{margin-left:-2px}.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]{width:60px;margin-right:10px}.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[_ngcontent-%COMP%]::-webkit-outer-spin-button, .form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[_ngcontent-%COMP%]::-webkit-inner-spin-button{-webkit-appearance:none;margin:0}.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[type=number][_ngcontent-%COMP%]{-moz-appearance:textfield}.form-resolution[_ngcontent-%COMP%]   .form-select-duration[_ngcontent-%COMP%]{width:116px}"]})}}return n})();var mt=(n,s)=>[n,s];function ut(n,s){if(n&1&&(z(0),P(1,"app-document-title",1),i(2,"mat-card",2)(3,"mat-card-header")(4,"mat-toolbar")(5,"h2"),P(6,"mat-icon",3),u(7),c()()(),i(8,"mat-card-content"),P(9,"app-torrent-metrics"),c()(),H()),n&2){let e=s.$implicit;a(),l("parts",he(2,mt,e("routes.torrents"),e("routes.dashboard"))),a(6),T(e("routes.torrents"))}}var or=(()=>{class n{static{this.\u0275fac=function(r){return new(r||n)}}static{this.\u0275cmp=L({type:n,selectors:[["app-torrents"]],standalone:!0,features:[U],decls:1,vars:0,consts:[[4,"transloco"],[3,"parts"],[1,"dashboard-card"],["svgIcon","magnet"]],template:function(r,t){r&1&&F(0,ut,10,5,"ng-container",0)},dependencies:[Y,J,K,q,Q,Be,G,Qe,De],styles:["mat-card-header[_ngcontent-%COMP%]{flex-wrap:wrap}mat-card-header[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%]{font-size:18px;margin:0 60px 0 48px;height:48px;line-height:48px}mat-card-header[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%]{position:relative;top:6px;margin-right:14px;line-height:1.25rem}mat-card-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]{flex:0 0 100%}mat-card-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   a[_ngcontent-%COMP%]{margin-top:2px}mat-card-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   a[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%]{margin-right:12px}"]})}}return n})();export{or as TorrentsDashboardComponent};
