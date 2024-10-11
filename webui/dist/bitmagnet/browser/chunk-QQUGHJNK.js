import{a as Fe,b as He,c as Z,d as Ue,e as Ge,f as ee,g as Ke,h as Ye,i as j,j as Je,k as We,l as Xe}from"./chunk-T2EFAQWC.js";import{a as Te,b as Re,c as O}from"./chunk-WTRXJNHX.js";import{e as X}from"./chunk-AV47NZN7.js";import{a as Se,k as qe,m as Ve}from"./chunk-AMUSMSN6.js";import{a as we}from"./chunk-BRPPGP57.js";import{B as Be,E as G,F as Ee,G as Oe,J as K,Na as Le,Ta as W,Ua as $e,Z as Ie,a as H,b as ve,ca as Y,da as ze,fa as J,ga as Ne,ha as Pe,i as xe,ia as Qe,ka as Ae,n as Me,na as je,o as ye,oa as De,r as U}from"./chunk-EH7PJXWD.js";import"./chunk-6XXA7HXI.js";import"./chunk-FKMTSCBK.js";import{$b as _,B as ce,Cb as de,Ea as g,Fa as C,Hb as u,Kb as I,N as le,Nb as z,Ob as N,Pb as P,Q as ne,Qb as o,Rb as l,Sb as F,Tb as fe,Ub as _e,Wb as re,Zb as k,_c as ke,a as M,b as y,ka as R,kb as c,kc as p,l as A,la as ue,lc as S,mc as he,o as se,pc as be,qa as T,ua as me,uc as ge,va as pe,vc as Ce}from"./chunk-3DR3CJRN.js";var Ze=(()=>{class n{static{this.\u0275fac=function(i){return new(i||n)}}static{this.\u0275mod=pe({type:n})}static{this.\u0275inj=ue({imports:[W,J,Ie,ye,U,G,K,Y,Qe,Pe]})}}return n})();var lt={pending:"primary",processed:"success",failed:"error",retry:"caution"},et=(()=>{class n{constructor(){this.themeInfo=T(X),this.transloco=T(H)}create(e){let{colors:i}=this.themeInfo.info,t=Array(),a=[];if(e&&Array.from(new Set(e.queues.flatMap(v=>v.events?[v.events.earliestBucket,v.events.latestBucket]:[]))).sort().length){let v=e.queues.filter(h=>!h.isEmpty);t.push(...v.map(h=>h.queue));let x=Array();switch(e.params.event){case"created":x.push("pending");break;case"processed":x.push("processed");break;case"failed":x.push("retry","failed");break;default:x.push(...Ke);break}a.push(...x.map(h=>({label:h,data:v.map(d=>d.statusCounts[h]),backgroundColor:i[O(lt[h],50)]})))}return{type:"bar",options:{animation:!1,scales:{x:{ticks:{callback:r=>parseInt(r).toLocaleString(this.transloco.getActiveLang())}},y:{}},indexAxis:"y",plugins:{legend:{display:!0}}},data:{labels:t,datasets:a}}}static{this.\u0275fac=function(i){return new(i||n)}}static{this.\u0275prov=R({token:n,factory:n.\u0275fac,providedIn:"root"})}}return n})();var D="\\d+",oe="".concat(D,"(?:[\\.,]").concat(D,")?"),ut="(".concat(D,"Y)?(").concat(D,"M)?(").concat(D,"W)?(").concat(D,"D)?"),mt="T(".concat(oe,"H)?(").concat(oe,"M)?(").concat(oe,"S)?"),pt="P(?:".concat(ut,"(?:").concat(mt,")?)"),dt=["years","months","weeks","days","hours","minutes","seconds"],tt={years:0,months:0,weeks:0,days:0,hours:0,minutes:0,seconds:0},ft=new RegExp(pt),it=function(n){let s=n.replace(/,/g,".").match(ft);if(!s)throw new RangeError("invalid duration: ".concat(n));let e=s.slice(1);if(e.filter(function(i){return i!=null}).length===0)throw new RangeError("invalid duration: ".concat(n));if(e.filter(function(i){return/\./.test(i||"")}).length>1)throw new RangeError("only the smallest unit can be fractional");return e.reduce(function(i,t,a){return Object.assign(i,{[dt[a]]:parseFloat(t||"0")||0}),i},{})},_t=function(n,s){s||(s=new Date);let e=Object.assign({},tt,n),i=s.getTime(),t=new Date(i);t.setFullYear(t.getFullYear()+e.years),t.setMonth(t.getMonth()+e.months),t.setDate(t.getDate()+e.days);let a=e.hours*3600*1e3,r=e.minutes*60*1e3;return t.setMilliseconds(t.getMilliseconds()+e.seconds*1e3+a+r),t.setDate(t.getDate()+e.weeks*7),t},nt=function(n,s){s||(s=new Date);let e=Object.assign({},tt,n),i=s.getTime(),t=new Date(i),a=_t(e,t),r=s.getTimezoneOffset(),v=a.getTimezoneOffset(),x=(r-v)*60;return(a.getTime()-t.getTime())/1e3+x};var te=class{constructor(s,e=Ue,i){this.apollo=s,this.errorsService=i,this.rawResultSubject=new A({queue:{metrics:{buckets:[]}}}),this.resultSubject=new A(Ge),this.result$=this.resultSubject.asObservable(),this.loadingSubject=new A(!1),this.paramsSubject=new A(e),this.params$=this.paramsSubject.asObservable(),this.variablesSubject=new A(rt(e)),this.paramsSubject.pipe(ne(50)).subscribe(t=>{let a=this.variablesSubject.getValue(),r=rt(t);JSON.stringify(a)!==JSON.stringify(r)?this.variablesSubject.next(r):this.resultSubject.next(ot(t,this.rawResultSubject.getValue()))}),this.variablesSubject.pipe(ne(50)).subscribe(t=>this.request(t)),this.rawResultSubject.subscribe(t=>{let a=this.paramsSubject.getValue();this.resultSubject.next(ot(a,t)),this.setInterval(a.autoRefresh)})}setInterval(s){clearTimeout(this.refreshTimeout);let e=Xe[s??this.params.autoRefresh];e&&(this.refreshTimeout=setTimeout(()=>{this.refresh()},e*1e3))}get params(){return this.paramsSubject.getValue()}get bucketDuration(){let s=this.params.buckets.duration;return s==="AUTO"?"hour":s}get bucketMultiplier(){return this.resultSubject.getValue().params.buckets.multiplier??this.params.buckets.multiplier}get loading(){return this.loadingSubject.getValue()}setTimeframe(s){this.updateParams(e=>y(M({},e),{buckets:y(M({},e.buckets),{timeframe:s})}))}setQueue(s){this.updateParams(e=>y(M({},e),{queue:s??void 0}))}setBucketDuration(s,e){this.updateParams(i=>y(M({},i),{buckets:y(M({},i.buckets),{duration:s,multiplier:e??"AUTO"})}))}setBucketMultiplier(s){this.updateParams(e=>y(M({},e),{buckets:y(M({},e.buckets),{multiplier:s})}))}setEvent(s){this.updateParams(e=>y(M({},e),{event:s??void 0}))}setAutoRefreshInterval(s){this.updateParams(e=>y(M({},e),{autoRefresh:s}))}updateParams(s){this.paramsSubject.next(s(this.params))}refresh(){this.variablesSubject.next(this.variablesSubject.getValue())}request(s){return clearTimeout(this.refreshTimeout),this.loadingSubject.next(!0),this.apollo.query({query:Le,variables:s,fetchPolicy:"no-cache"}).pipe(ce(e=>{e&&(this.loadingSubject.next(!1),this.rawResultSubject.next(e.data))})).pipe(le(e=>(this.errorsService.addError(`Failed to load queue metrics: ${e.message}`),this.loadingSubject.next(!1),this.setInterval(),se))).subscribe()}},rt=n=>({input:{bucketDuration:n.buckets.duration==="AUTO"?"hour":n.buckets.duration,queues:n.queue?[n.queue]:void 0,startTime:n.buckets.timeframe==="all"?void 0:new Date(new Date().getTime()-1e3*j[n.buckets.timeframe]).toISOString()}}),at=n=>Object.fromEntries(n),ot=(n,s)=>{let{bucketParams:e,earliestBucket:i,latestBucket:t}=ht(n,s),a=Object.entries(s.queue.metrics.buckets.reduce((h,d)=>{if(d.queue!==(n.queue??d.queue))return h;let f,m;if((n.event??!0)&&(f=q(d.createdAtBucket,e),i&&i.index>f.index&&(f=void 0)),d.ranAtBucket&&n.event!=="created"&&(m=q(d.ranAtBucket,e),m&&(t.index<m.index||i&&i.index>m.index)&&(m=void 0)),d.queue!==n.queue&&!f&&(!m||d.status==="pending"))return h;let[V,b]=h[d.queue]??[Fe,[]],w=d.latency?nt(it(d.latency)):void 0;return y(M({},h),{[d.queue]:[(d.status==="pending"?f:m)?y(M({},V),{[d.status]:d.count+V[d.status]}):V,{created:f?y(M({},b.created),{[f.key]:{count:d.count+(b.created?.[f.key]?.count??0),latency:0,startTime:f.start}}):b.created,processed:m&&d.status==="processed"&&(n.event??!0)?y(M({},b.processed),{[m.key]:{count:d.count+(b.processed?.[m.key]?.count??0),latency:(b.processed?.[m.key]?.latency??0)+(w??0),startTime:m.start}}):b.processed,failed:m&&d.status==="failed"&&(n.event??!0)?y(M({},b.failed),{[m.key]:{count:d.count+(b.failed?.[m.key]?.count??0),latency:(b.failed?.[m.key]?.latency??0)+(w??0),startTime:m.start}}):b.failed}]})},{})).map(([h,[d,f]])=>{let m;if(Object.keys(f).length){let V=Array(),b=at(Array("created","processed","failed").flatMap(w=>{let B=at(Object.entries(f[w]??{}).filter(([,ie])=>ie?.count).sort(([ie],[ct])=>parseInt(ie)<parseInt(ct)?1:-1)),E=Object.keys(B);if(!E.length)return[];let $=parseInt(E[0]),Q=parseInt(E[E.length-1]);return V.push($,Q),[[w,{earliestBucket:$,latestBucket:Q,entries:B}]]}));V.sort(),m={bucketDuration:e.duration,earliestBucket:V[0],latestBucket:V[V.length-1],eventBuckets:b}}return{queue:h,statusCounts:d,events:m,isEmpty:!m?.eventBuckets}}),r,v=a.flatMap(h=>h.events?[h.events.earliestBucket]:[]).sort()[0],x=a.flatMap(h=>h.events?[h.events.latestBucket]:[]).sort().reverse()[0];return v&&x&&(r={earliestBucket:v,latestBucket:x}),{params:y(M({},n),{buckets:e}),queues:a,bucketSpan:r}},ht=(n,s)=>{let e=n.buckets.duration==="AUTO"?"hour":n.buckets.duration,i=n.buckets.multiplier==="AUTO"?1:n.buckets.multiplier,t=n.buckets.timeframe,a=new Date,r=q(a,{duration:e,multiplier:i}),v=t==="all"?void 0:q(a.getTime()-1e3*j[t],{duration:e,multiplier:i}),x=[...v?[v]:[],...s.queue.metrics.buckets.flatMap(f=>[q(f.createdAtBucket,{duration:e,multiplier:i}),...f.ranAtBucket?[q(f.ranAtBucket,{duration:e,multiplier:i})]:[]]),r].filter(f=>!v||f.index>=v.index).sort((f,m)=>f.index-m.index),h=x[0],d=x[x.length-1];if(n.buckets.multiplier==="AUTO"){let m=d.index-h.index;i=Math.min(60,Math.max(Math.floor(m/(20*5))*5,1))}return{bucketParams:{duration:e,multiplier:i,timeframe:t},earliestBucket:t==="all"?void 0:q(a.getTime()-1e3*j[t],{duration:e,multiplier:i}),latestBucket:q(Math.max(a.getTime(),d.start.getTime()),{duration:e,multiplier:i})}},q=(n,s)=>{let e=new Date(n),i=1e3*Z[s.duration]*s.multiplier,t=Math.floor(e.getTime()/i);return{key:`${t}`,index:t,start:new Date(t*i)}};var L={created:"primary",processed:"success",failed:"error"},st=(()=>{class n{constructor(){this.themeInfo=T(X),this.transloco=T(H)}create(e){let{colors:i}=this.themeInfo.info,t=Array(),a=[];if(e){let r=e.queues.filter(f=>!f.isEmpty),v=Array.from(new Set(r.flatMap(f=>f.events?[f.events.earliestBucket,f.events.latestBucket]:[]))).sort(),x=new Date,h=e.params.buckets.timeframe==="all"?v[0]:Math.min(v[0],q(x.getTime()-1e3*j[e.params.buckets.timeframe],e.params.buckets).index),d=Math.max(v[v.length-1],q(x,e.params.buckets).index);if(v.length){for(let m=h;m<=d;m++)t.push(this.formatBucketKey(e.params.buckets,m));let f=ee.filter(m=>(e.params.event??m)===m);for(let m of r){for(let b of f){let w=Array();for(let B=h;B<=d;B++)w.push(m.events?.eventBuckets?.[b]?.entries?.[`${B}`]?.count??0);a.push({yAxisID:"yCount",label:[m.queue,b].join("/"),data:w,borderColor:i[O(L[b],50)],pointBackgroundColor:i[O(L[b],20)],pointBorderColor:i[O(L[b],80)],pointHoverBackgroundColor:i[O(L[b],40)],pointHoverBorderColor:i[O(L[b],60)]})}if(["processed","failed"].filter(b=>f.includes(b)).length){let b=Array();for(let w=h;w<=d;w++){let B=["processed","failed"].filter(E=>f.includes(E)).reduce((E,$)=>{let Q=m.events?.eventBuckets?.[$]?.entries?.[`${w}`];return Q?.count?[(E?.[0]??0)+Q.latency,(E?.[1]??0)+Q.count]:E},null);b.push(B?B[0]/B[1]:null)}a.push({yAxisID:"yLatency",label:[m.queue,"latency"].join("/"),data:b,borderColor:i["tertiary-50"],pointHoverBackgroundColor:i["tertiary-80"],pointHoverBorderColor:i["tertiary-20"]})}}}}return{type:"line",options:{animation:!1,elements:{line:{tension:.5}},scales:{yCount:{position:"left",ticks:{callback:r=>parseInt(r).toLocaleString(this.transloco.getActiveLang())}},yLatency:{position:"right",ticks:{callback:this.formatDuration.bind(this)}}},plugins:{legend:{display:!0},decimation:{enabled:!0}}},data:{labels:t,datasets:a}}}formatBucketKey(e,i){let t;switch(e.duration){case"day":t="d LLL";break;case"hour":t="d LLL H:00";break;case"minute":t="H:mm";break}return Te(1e3*Z[e.duration]*e.multiplier*i,t,{locale:qe(this.transloco.getActiveLang())})}formatDuration(e){if(typeof e=="string"&&(e=parseInt(e)),e===0)return"0";let i=e,t=0,a=0,r=0;return i>=60&&(t=Math.floor(i/60),i=i%60,t>=5&&(i=0,t>=60&&(a=Math.floor(t/60),t=t%60,a>=5&&(t=0,a>=24&&(r=Math.floor(a/24),a=a%24))))),Ve({days:r,hours:a,minutes:t,seconds:i},this.transloco.getActiveLang())}static{this.\u0275fac=function(i){return new(i||n)}}static{this.\u0275prov=R({token:n,factory:n.\u0275fac,providedIn:"root"})}}return n})();function bt(n,s){if(n&1&&(o(0,"mat-option",6),p(1),l()),n&2){let e=s.$implicit,i=_().$implicit;u("value",e),c(),S(i("dashboard.interval."+e))}}function gt(n,s){if(n&1&&(o(0,"mat-option",6),p(1),l()),n&2){let e=s.$implicit,i=_().$implicit;u("value",e),c(),S(i("dashboard.interval."+e+"s"))}}function Ct(n,s){if(n&1&&(o(0,"mat-option",6),p(1),l()),n&2){let e=s.$implicit;u("value",e),c(),S(e)}}function kt(n,s){if(n&1){let e=re();o(0,"button",17),k("click",function(){let t=g(e).$implicit,a=_(2);return C(a.queueMetricsController.params.queue===t||a.queueMetricsController.setQueue(t))}),o(1,"mat-icon"),p(2),l()()}if(n&2){let e=s.$implicit,i=_(2);I(i.queueMetricsController.params.queue===e?"selected":"deselected"),u("matTooltip",e),c(2),S(i.queueMetricsController.params.queue===e?"radio_button_checked":"radio_button_unchecked")}}function vt(n,s){if(n&1&&(o(0,"mat-option",6),p(1),l()),n&2){let e=s.$implicit,i=_().$implicit;u("value",e),c(),S(i("dashboard.event."+e))}}function xt(n,s){if(n&1&&(o(0,"mat-option",6),p(1),l()),n&2){let e=s.$implicit,i=_().$implicit;u("value",e),c(),S(i("dashboard.interval."+e))}}function Mt(n,s){if(n&1){let e=re();fe(0),o(1,"mat-card")(2,"mat-card-content")(3,"mat-grid-list",1)(4,"mat-grid-tile",2)(5,"mat-card",3)(6,"mat-card-header")(7,"mat-card-title")(8,"h4"),p(9),l()()(),o(10,"mat-card-content")(11,"mat-form-field",4)(12,"mat-select",5),k("valueChange",function(t){g(e);let a=_();return C(a.queueMetricsController.setTimeframe(t))}),N(13,bt,2,2,"mat-option",6,z),l()(),o(15,"div",7)(16,"button",8),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setTimeframe(t.timeframeNames[0]))}),o(17,"mat-icon"),p(18,"first_page"),l()(),o(19,"button",8),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setTimeframe(t.timeframeNames[t.timeframeNames.indexOf(t.queueMetricsController.params.buckets.timeframe)-1]))}),o(20,"mat-icon"),p(21,"navigate_before"),l()(),o(22,"button",8),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setTimeframe(t.timeframeNames[t.timeframeNames.indexOf(t.queueMetricsController.params.buckets.timeframe)+1]))}),o(23,"mat-icon"),p(24,"navigate_next"),l()(),o(25,"button",8),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setTimeframe(t.timeframeNames[t.timeframeNames.length-1]))}),o(26,"mat-icon"),p(27,"last_page"),l()()()()()(),o(28,"mat-grid-tile",2)(29,"mat-card",9)(30,"mat-card-header")(31,"mat-card-title")(32,"h4"),p(33),l()()(),o(34,"mat-card-content")(35,"mat-form-field",10)(36,"input",11),ge(37,"async"),k("change",function(t){g(e);let a=_();return C(a.handleMultiplierEvent(t))}),l()(),o(38,"mat-form-field",12)(39,"mat-select",5),k("valueChange",function(t){g(e);let a=_();return C(a.queueMetricsController.setBucketDuration(t))}),N(40,gt,2,2,"mat-option",6,z),l()(),o(42,"div",7)(43,"button",8),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setBucketMultiplier(t.queueMetricsController.bucketMultiplier-1))}),o(44,"mat-icon"),p(45,"remove"),l()(),o(46,"button",13),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setBucketMultiplier(t.queueMetricsController.bucketMultiplier+1))}),o(47,"mat-icon"),p(48,"add"),l()(),o(49,"button",8),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setBucketDuration(t.resolutionNames[0]))}),o(50,"mat-icon"),p(51,"first_page"),l()(),o(52,"button",8),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setBucketDuration(t.resolutionNames[t.resolutionNames.indexOf(t.queueMetricsController.bucketDuration)-1]))}),o(53,"mat-icon"),p(54,"navigate_before"),l()(),o(55,"button",8),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setBucketDuration(t.resolutionNames[t.resolutionNames.indexOf(t.queueMetricsController.bucketDuration)+1]))}),o(56,"mat-icon"),p(57,"navigate_next"),l()(),o(58,"button",8),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setBucketDuration(t.resolutionNames[t.resolutionNames.length-1]))}),o(59,"mat-icon"),p(60,"last_page"),l()()()()()(),o(61,"mat-grid-tile",2)(62,"mat-card",14)(63,"mat-card-header")(64,"mat-card-title")(65,"h4"),p(66),l()()(),o(67,"mat-card-content")(68,"mat-form-field",4)(69,"mat-select",5),k("valueChange",function(t){g(e);let a=_();return C(a.queueMetricsController.setQueue(t==="_all"?null:t))}),o(70,"mat-option",15),p(71),l(),N(72,Ct,2,2,"mat-option",6,z),l()(),o(74,"div",16)(75,"button",17),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setQueue(null))}),o(76,"mat-icon",18),p(77,"workspaces"),l()(),N(78,kt,3,4,"button",19,z),l()()()(),o(80,"mat-grid-tile",2)(81,"mat-card")(82,"mat-card-header")(83,"mat-card-title")(84,"h4"),p(85),l()()(),o(86,"mat-card-content")(87,"mat-form-field",4)(88,"mat-select",5),k("valueChange",function(t){g(e);let a=_();return C(a.queueMetricsController.setEvent(t==="_all"?null:t))}),o(89,"mat-option",15),p(90,"All"),l(),N(91,vt,2,2,"mat-option",6,z),l()(),o(93,"div",16)(94,"button",17),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.setEvent(null))}),o(95,"mat-icon",18),p(96,"radio_button_checked"),l()(),o(97,"button",17),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.params.event==="created"||t.queueMetricsController.setEvent("created"))}),o(98,"mat-icon"),p(99,"add_circle"),l()(),o(100,"button",17),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.params.event==="processed"||t.queueMetricsController.setEvent("processed"))}),o(101,"mat-icon"),p(102,"check_circle"),l()(),o(103,"button",17),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.params.event==="failed"||t.queueMetricsController.setEvent("failed"))}),o(104,"mat-icon"),p(105,"error"),l()()()()()(),o(106,"mat-grid-tile",2)(107,"mat-card",20)(108,"mat-card-header")(109,"mat-card-title")(110,"h4"),p(111),l()()(),o(112,"mat-card-content")(113,"mat-form-field",4)(114,"mat-select",5),k("valueChange",function(t){g(e);let a=_();return C(a.queueMetricsController.setAutoRefreshInterval(t))}),N(115,xt,2,2,"mat-option",6,z),l()(),o(117,"div",16)(118,"button",17),k("click",function(){g(e);let t=_();return C(t.queueMetricsController.refresh())}),o(119,"mat-icon"),p(120,"sync"),l()()()()()()(),o(121,"div",21),F(122,"mat-progress-bar",22),l(),o(123,"mat-grid-list",1)(124,"mat-grid-tile",2)(125,"mat-card")(126,"mat-card-header")(127,"mat-card-title")(128,"h4"),p(129),l()()(),o(130,"mat-card-content"),F(131,"app-chart",23),l()()(),o(132,"mat-grid-tile",2)(133,"mat-card")(134,"mat-card-header")(135,"mat-card-title")(136,"h4"),p(137),l()()(),o(138,"mat-card-content"),F(139,"app-chart",23),l()()()()()(),_e()}if(n&2){let e,i,t,a=s.$implicit,r=_();c(3),u("cols",r.breakpoints.sizeAtLeast("Large")?5:r.breakpoints.sizeAtLeast("Medium")?3:r.breakpoints.sizeAtLeast("Small")?2:1),c(),u("colspan",1)("rowspan",2),c(5),S(a("dashboard.metrics.timeframe")),c(3),u("value",r.queueMetricsController.params.buckets.timeframe),c(),P(r.timeframeNames),c(3),u("disabled",r.timeframeNames.indexOf(r.queueMetricsController.params.buckets.timeframe)<=0),c(3),u("disabled",r.timeframeNames.indexOf(r.queueMetricsController.params.buckets.timeframe)<=0),c(3),u("disabled",r.timeframeNames.indexOf(r.queueMetricsController.params.buckets.timeframe)>=r.timeframeNames.length-1),c(3),u("disabled",r.timeframeNames.indexOf(r.queueMetricsController.params.buckets.timeframe)>=r.timeframeNames.length-1),c(3),u("colspan",1)("rowspan",2),c(5),he(" ",a("dashboard.metrics.resolution")," "),c(3),u("placeholder",(e=(e=Ce(37,66,r.queueMetricsController.result$))==null||e.params==null||e.params.buckets==null||e.params.buckets.multiplier==null?null:e.params.buckets.multiplier.toString())!==null&&e!==void 0?e:"")("value",r.queueMetricsController.params.buckets.multiplier),c(3),u("value",r.queueMetricsController.bucketDuration),c(),P(r.resolutionNames),c(3),u("disabled",r.queueMetricsController.bucketMultiplier===1),c(6),u("disabled",r.resolutionNames.indexOf(r.queueMetricsController.bucketDuration)<=0),c(3),u("disabled",r.resolutionNames.indexOf(r.queueMetricsController.bucketDuration)<=0),c(3),u("disabled",r.resolutionNames.indexOf(r.queueMetricsController.bucketDuration)>=r.resolutionNames.length-1),c(3),u("disabled",r.resolutionNames.indexOf(r.queueMetricsController.bucketDuration)>=r.resolutionNames.length-1),c(3),u("colspan",1)("rowspan",2),c(5),S(a("dashboard.queues.queue")),c(3),u("value",(i=r.queueMetricsController.params.queue)!==null&&i!==void 0?i:"_all"),c(2),S(a("general.all")),c(),P(r.availableQueueNames),c(3),I(r.queueMetricsController.params.queue?"deselected":"selected"),u("matTooltip",a("general.all")),c(3),P(r.availableQueueNames),c(2),u("colspan",1)("rowspan",2),c(5),S(a("dashboard.metrics.event")),c(3),u("value",(t=r.queueMetricsController.params.event)!==null&&t!==void 0?t:"_all"),c(3),P(r.eventNames),c(3),I(r.queueMetricsController.params.event?"deselected":"selected"),u("matTooltip",a("general.all")),c(3),I(r.queueMetricsController.params.event==="created"?"selected":"deselected"),u("matTooltip",a("dashboard.queues.created")),c(3),I(r.queueMetricsController.params.event==="processed"?"selected":"deselected"),u("matTooltip",a("dashboard.queues.processed")),c(3),I(r.queueMetricsController.params.event==="failed"?"selected":"deselected"),u("matTooltip",a("dashboard.queues.failed")),c(3),u("colspan",1)("rowspan",2),c(5),S(a("general.refresh")),c(3),u("value",r.queueMetricsController.params.autoRefresh),c(),P(r.autoRefreshIntervalNames),c(3),u("matTooltip",a("general.refresh")),c(4),u("mode",r.queueMetricsController.loading?"indeterminate":"determinate")("value",0),c(),u("cols",r.breakpoints.sizeAtLeast("Large")?2:1),c(),u("colspan",1)("rowspan",5),c(5),S(a("dashboard.queues.total_counts_by_status")),c(2),u("adapter",r.totals)("$data",r.queueMetricsController.result$)("height",400)("width",550),c(),u("colspan",1)("rowspan",5),c(5),S(a("dashboard.metrics.throughput")),c(2),u("adapter",r.timeline)("$data",r.queueMetricsController.result$)("height",400)("width",550)}}var zi=(()=>{class n{constructor(){this.breakpoints=T(we),this.apollo=T(xe),this.queueMetricsController=new te(this.apollo,{buckets:{duration:"AUTO",multiplier:"AUTO",timeframe:"all"},autoRefresh:"seconds_30"},T(Se)),this.timeline=T(st),this.totals=T(et),this.resolutionNames=He,this.timeframeNames=Ye,this.availableQueueNames=Je,this.autoRefreshIntervalNames=We,this.eventNames=ee}ngOnInit(){this.queueMetricsController.result$.subscribe(e=>{if(this.queueMetricsController.params.buckets.timeframe==="all"&&this.queueMetricsController.params.buckets.duration==="AUTO"&&e.params.buckets.duration==="hour"){let i=e.bucketSpan;i&&i.latestBucket-i.earliestBucket<12&&this.queueMetricsController.setBucketDuration("minute")}})}ngOnDestroy(){this.queueMetricsController.setAutoRefreshInterval("off")}handleMultiplierEvent(e){let i=e.currentTarget.value;this.queueMetricsController.setBucketMultiplier(/^\d+$/.test(i)?parseInt(i):"AUTO")}static{this.\u0275fac=function(i){return new(i||n)}}static{this.\u0275cmp=me({type:n,selectors:[["app-queue-visualize"]],standalone:!0,features:[be],decls:1,vars:0,consts:[[4,"transloco"],["rowHeight","100px",3,"cols"],[3,"colspan","rowspan"],[1,"form-timeframe"],["subscriptSizing","dynamic"],[3,"valueChange","value"],[3,"value"],[1,"paginator","actions"],["mat-icon-button","",3,"click","disabled"],[1,"form-resolution"],["subscriptSizing","dynamic",1,"form-input-multiplier"],["type","number","matInput","","min","1","step","1",3,"change","placeholder","value"],["subscriptSizing","dynamic",1,"form-select-duration"],["mat-icon-button","",3,"click"],[1,"form-queues"],["value","_all"],[1,"actions"],["mat-icon-button","",3,"click","matTooltip"],["fontSet","material-icons"],["mat-icon-button","",3,"class","matTooltip"],[1,"form-refresh"],[1,"progress-bar-container"],[3,"mode","value"],[3,"adapter","$data","height","width"]],template:function(i,t){i&1&&de(0,Mt,140,68,"ng-container",0)},dependencies:[$e,Me,U,G,Oe,K,Ee,Be,ze,Y,J,Ne,Ae,De,je,ve,ke,Re,W,Ze],styles:[".actions[_ngcontent-%COMP%]{width:210px;padding-top:12px;--mdc-icon-button-state-layer-size: 32px}.actions[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%]{font-size:22px}.actions[_ngcontent-%COMP%]   button[_ngcontent-%COMP%]{margin-right:0}.progress-bar-container[_ngcontent-%COMP%]{width:100%;height:10px}mat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]{width:100%}mat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]   mat-card-content[_ngcontent-%COMP%]{min-width:190px}mat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]   h4[_ngcontent-%COMP%]{margin-bottom:16px;font-size:18px}mat-form-field[_ngcontent-%COMP%]{width:186px}.form-resolution[_ngcontent-%COMP%]   .actions[_ngcontent-%COMP%]{margin-left:-2px}.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]{width:60px;margin-right:10px}.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[_ngcontent-%COMP%]::-webkit-outer-spin-button, .form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[_ngcontent-%COMP%]::-webkit-inner-spin-button{-webkit-appearance:none;margin:0}.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[type=number][_ngcontent-%COMP%]{-moz-appearance:textfield}.form-resolution[_ngcontent-%COMP%]   .form-select-duration[_ngcontent-%COMP%]{width:116px}"]})}}return n})();export{zi as QueueVisualizeComponent};
