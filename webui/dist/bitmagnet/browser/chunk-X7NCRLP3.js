import{a as Ue,b as Ge,c as te,d as Ke,e as Ye,f as ie,g as Je,h as We,i as j,j as Xe,k as Ze,l as et}from"./chunk-6OKZFKY7.js";import{a as Fe,b as ee,c as I}from"./chunk-2GOIBOCD.js";import{e as Z}from"./chunk-Y7K23DTG.js";import{b as He}from"./chunk-RQ2LQKI2.js";import{l as Re}from"./chunk-42PJPEMD.js";import{a as Te}from"./chunk-ZOC7YF5Y.js";import{a as Le}from"./chunk-B5NIMPHN.js";import{a as $e}from"./chunk-DSEDLZDW.js";import{Ba as De,Ea as Pe,Fa as Qe,J as qe,P as K,S as Ve,V as Y,W as Be,X as Oe,_ as J,_a as je,a as U,b as xe,m as ye,oa as Ee,s as Se,t as we,ta as W,ua as Ie,w as G,wa as X,xa as ze,ya as Ne,za as Ae}from"./chunk-TFL6A5NY.js";import"./chunk-6XXA7HXI.js";import"./chunk-CMNWCZJM.js";import{$b as d,B as le,Cb as fe,Ea as b,Fa as g,Hb as l,Kb as N,N as ue,Nb as A,Ob as D,Pb as P,Q as re,Qb as s,Rb as u,Sb as $,Tb as _e,Ub as he,Wb as ae,Zb as C,a as y,ad as Me,b as S,ka as H,kb as c,kc as f,l as Q,la as me,lc as V,mc as be,o as ce,pc as ge,qa as q,tc as Ce,ua as pe,va as de,vc as ke,wc as ve}from"./chunk-Z3WUIYN5.js";var tt=(()=>{class n{static{this.\u0275fac=function(i){return new(i||n)}}static{this.\u0275mod=de({type:n})}static{this.\u0275inj=me({imports:[K,X,Ee,we,G,Y,J,W,Ae,Ne,ee]})}}return n})();var mt={pending:"primary",processed:"success",failed:"error",retry:"caution"},it=(()=>{class n{constructor(){this.themeInfo=q(Z),this.transloco=q(U)}create(e,i){let{colors:t}=this.themeInfo.info,a=Array(),r=[];if(e&&Array.from(new Set(e.queues.flatMap(v=>v.events?[v.events.earliestBucket,v.events.latestBucket]:[]))).sort().length){let v=e.queues.filter(m=>!m.isEmpty);a.push(...v.map(m=>m.queue));let h=Array();switch(e.params.event){case"created":h.push("pending");break;case"processed":h.push("processed");break;case"failed":h.push("retry","failed");break;default:h.push(...Je);break}r.push(...h.map(m=>({label:this.transloco.translate("dashboard.queues."+m),data:v.map(_=>_.statusCounts[m]),backgroundColor:t[I(mt[m],50)]})))}return{type:"bar",options:{animation:!1,responsive:!0,scales:{x:{ticks:{callback:k=>parseInt(k).toLocaleString(this.transloco.getActiveLang())}},y:{}},indexAxis:"y",plugins:{legend:{display:i.legend}}},data:{labels:a,datasets:r}}}static{this.\u0275fac=function(i){return new(i||n)}}static{this.\u0275prov=H({token:n,factory:n.\u0275fac,providedIn:"root"})}}return n})();var L="\\d+",se="".concat(L,"(?:[\\.,]").concat(L,")?"),pt="(".concat(L,"Y)?(").concat(L,"M)?(").concat(L,"W)?(").concat(L,"D)?"),dt="T(".concat(se,"H)?(").concat(se,"M)?(").concat(se,"S)?"),ft="P(?:".concat(pt,"(?:").concat(dt,")?)"),_t=["years","months","weeks","days","hours","minutes","seconds"],nt={years:0,months:0,weeks:0,days:0,hours:0,minutes:0,seconds:0},ht=new RegExp(ft),rt=function(n){let o=n.replace(/,/g,".").match(ht);if(!o)throw new RangeError("invalid duration: ".concat(n));let e=o.slice(1);if(e.filter(function(i){return i!=null}).length===0)throw new RangeError("invalid duration: ".concat(n));if(e.filter(function(i){return/\./.test(i||"")}).length>1)throw new RangeError("only the smallest unit can be fractional");return e.reduce(function(i,t,a){return Object.assign(i,{[_t[a]]:parseFloat(t||"0")||0}),i},{})},bt=function(n,o){o||(o=new Date);let e=Object.assign({},nt,n),i=o.getTime(),t=new Date(i);t.setFullYear(t.getFullYear()+e.years),t.setMonth(t.getMonth()+e.months),t.setDate(t.getDate()+e.days);let a=e.hours*3600*1e3,r=e.minutes*60*1e3;return t.setMilliseconds(t.getMilliseconds()+e.seconds*1e3+a+r),t.setDate(t.getDate()+e.weeks*7),t},at=function(n,o){o||(o=new Date);let e=Object.assign({},nt,n),i=o.getTime(),t=new Date(i),a=bt(e,t),r=o.getTimezoneOffset(),k=a.getTimezoneOffset(),v=(r-k)*60;return(a.getTime()-t.getTime())/1e3+v};var ne=class{constructor(o,e=Ke,i){this.apollo=o,this.errorsService=i,this.rawResultSubject=new Q({queue:{metrics:{buckets:[]}}}),this.resultSubject=new Q(Ye),this.result$=this.resultSubject.asObservable(),this.loadingSubject=new Q(!1),this.paramsSubject=new Q(e),this.params$=this.paramsSubject.asObservable(),this.variablesSubject=new Q(ot(e)),this.paramsSubject.pipe(re(50)).subscribe(t=>{let a=this.variablesSubject.getValue(),r=ot(t);JSON.stringify(a)!==JSON.stringify(r)?this.variablesSubject.next(r):this.resultSubject.next(ct(t,this.rawResultSubject.getValue()))}),this.variablesSubject.pipe(re(50)).subscribe(t=>this.request(t)),this.rawResultSubject.subscribe(t=>{let a=this.paramsSubject.getValue();this.resultSubject.next(ct(a,t)),this.setInterval(a.autoRefresh)})}setInterval(o){clearTimeout(this.refreshTimeout);let e=et[o??this.params.autoRefresh];e&&(this.refreshTimeout=setTimeout(()=>{this.refresh()},e*1e3))}get params(){return this.paramsSubject.getValue()}get bucketDuration(){let o=this.params.buckets.duration;return o==="AUTO"?"hour":o}get bucketMultiplier(){return this.resultSubject.getValue().params.buckets.multiplier??this.params.buckets.multiplier}get loading(){return this.loadingSubject.getValue()}setTimeframe(o){this.updateParams(e=>S(y({},e),{buckets:S(y({},e.buckets),{timeframe:o})}))}setQueue(o){this.updateParams(e=>S(y({},e),{queue:o??void 0}))}setBucketDuration(o,e){this.updateParams(i=>S(y({},i),{buckets:S(y({},i.buckets),{duration:o,multiplier:e??"AUTO"})}))}setBucketMultiplier(o){this.updateParams(e=>S(y({},e),{buckets:S(y({},e.buckets),{multiplier:o})}))}setEvent(o){this.updateParams(e=>S(y({},e),{event:o??void 0}))}setAutoRefreshInterval(o){this.updateParams(e=>S(y({},e),{autoRefresh:o}))}updateParams(o){this.paramsSubject.next(o(this.params))}refresh(){this.variablesSubject.next(this.variablesSubject.getValue())}request(o){return clearTimeout(this.refreshTimeout),this.loadingSubject.next(!0),this.apollo.query({query:qe,variables:o,fetchPolicy:"no-cache"}).pipe(le(e=>{e&&(this.loadingSubject.next(!1),this.rawResultSubject.next(e.data))})).pipe(ue(e=>(this.errorsService.addError(`Failed to load queue metrics: ${e.message}`),this.loadingSubject.next(!1),this.setInterval(),ce))).subscribe()}},ot=n=>({input:{bucketDuration:n.buckets.duration==="AUTO"?"hour":n.buckets.duration,queues:n.queue?[n.queue]:void 0,startTime:n.buckets.timeframe==="all"?void 0:new Date(new Date().getTime()-1e3*j[n.buckets.timeframe]).toISOString()}}),st=n=>Object.fromEntries(n),ct=(n,o)=>{let{bucketParams:e,earliestBucket:i,latestBucket:t}=gt(n,o),a=Object.entries(o.queue.metrics.buckets.reduce((h,m)=>{if(m.queue!==(n.queue??m.queue))return h;let _,p;if((n.event??!0)&&(_=B(m.createdAtBucket,e),i&&i.index>_.index&&(_=void 0)),m.ranAtBucket&&n.event!=="created"&&(p=B(m.ranAtBucket,e),p&&(t.index<p.index||i&&i.index>p.index)&&(p=void 0)),m.queue!==n.queue&&!_&&(!p||m.status==="pending"))return h;let[M,w]=h[m.queue]??[Ue,[]],x=m.latency?at(rt(m.latency)):void 0;return S(y({},h),{[m.queue]:[(m.status==="pending"?_:p)?S(y({},M),{[m.status]:m.count+M[m.status]}):M,{created:_?S(y({},w.created),{[_.key]:{count:m.count+(w.created?.[_.key]?.count??0),latency:0,startTime:_.start}}):w.created,processed:p&&m.status==="processed"&&(n.event??!0)?S(y({},w.processed),{[p.key]:{count:m.count+(w.processed?.[p.key]?.count??0),latency:(w.processed?.[p.key]?.latency??0)+(x??0),startTime:p.start}}):w.processed,failed:p&&m.status==="failed"&&(n.event??!0)?S(y({},w.failed),{[p.key]:{count:m.count+(w.failed?.[p.key]?.count??0),latency:(w.failed?.[p.key]?.latency??0)+(x??0),startTime:p.start}}):w.failed}]})},{})).map(([h,[m,_]])=>{let p;if(Object.keys(_).length){let M=Array(),w=st(Array("created","processed","failed").flatMap(x=>{let O=st(Object.entries(_[x]??{}).filter(([,z])=>z?.count).sort(([z],[ut])=>parseInt(z)<parseInt(ut)?1:-1)),T=Object.keys(O);if(!T.length)return[];let E=parseInt(T[0]),R=parseInt(T[T.length-1]);return M.push(E,R),[[x,{earliestBucket:E,latestBucket:R,entries:O}]]}));M.sort(),p={bucketDuration:e.duration,earliestBucket:M[0],latestBucket:M[M.length-1],eventBuckets:w}}return{queue:h,statusCounts:m,events:p,isEmpty:!p?.eventBuckets}}),r,k=a.flatMap(h=>h.events?[h.events.earliestBucket]:[]).sort()[0],v=a.flatMap(h=>h.events?[h.events.latestBucket]:[]).sort().reverse()[0];return k&&v&&(r={earliestBucket:k,latestBucket:v}),{params:S(y({},n),{buckets:e}),queues:a,bucketSpan:r}},gt=(n,o)=>{let e=n.buckets.duration==="AUTO"?"hour":n.buckets.duration,i=n.buckets.multiplier==="AUTO"?1:n.buckets.multiplier,t=n.buckets.timeframe,a=new Date,r=B(a,{duration:e,multiplier:i}),k=t==="all"?void 0:B(a.getTime()-1e3*j[t],{duration:e,multiplier:i}),v=[...k?[k]:[],...o.queue.metrics.buckets.flatMap(_=>[B(_.createdAtBucket,{duration:e,multiplier:i}),..._.ranAtBucket?[B(_.ranAtBucket,{duration:e,multiplier:i})]:[]]),r].filter(_=>!k||_.index>=k.index).sort((_,p)=>_.index-p.index),h=v[0],m=v[v.length-1];if(n.buckets.multiplier==="AUTO"){let p=m.index-h.index;i=Math.min(60,Math.max(Math.floor(p/(20*5))*5,1))}return{bucketParams:{duration:e,multiplier:i,timeframe:t},earliestBucket:t==="all"?void 0:B(a.getTime()-1e3*j[t],{duration:e,multiplier:i}),latestBucket:B(Math.max(a.getTime(),m.start.getTime()),{duration:e,multiplier:i})}},B=(n,o)=>{let e=new Date(n),i=1e3*te[o.duration]*o.multiplier,t=Math.floor(e.getTime()/i);return{key:`${t}`,index:t,start:new Date(t*i)}};var F={created:"primary",processed:"success",failed:"error"},lt=(()=>{class n{constructor(){this.themeInfo=q(Z),this.transloco=q(U)}create(e,i){let{colors:t}=this.themeInfo.info,a=Array(),r=[];if(e){let k=e.queues.filter(p=>!p.isEmpty),v=Array.from(new Set(k.flatMap(p=>p.events?[p.events.earliestBucket,p.events.latestBucket]:[]))).sort(),h=new Date,m=e.params.buckets.timeframe==="all"?v[0]:Math.min(v[0],B(h.getTime()-1e3*j[e.params.buckets.timeframe],e.params.buckets).index),_=Math.max(v[v.length-1],B(h,e.params.buckets).index);if(v.length){for(let M=m;M<=_;M++)a.push(this.formatBucketKey(e.params.buckets,M));let p=ie.filter(M=>(e.params.event??M)===M);for(let M of k){for(let x of p){let O=Array();for(let T=m;T<=_;T++)O.push(M.events?.eventBuckets?.[x]?.entries?.[`${T}`]?.count??0);r.push({yAxisID:"yCount",label:M.queue+": "+this.transloco.translate("dashboard.queues."+x),data:O,borderColor:t[I(F[x],50)],pointBackgroundColor:t[I(F[x],20)],pointBorderColor:t[I(F[x],80)],pointHoverBackgroundColor:t[I(F[x],40)],pointHoverBorderColor:t[I(F[x],60)]})}if(["processed","failed"].filter(x=>p.includes(x)).length){let x=Array();for(let O=m;O<=_;O++){let T=["processed","failed"].filter(E=>p.includes(E)).reduce((E,R)=>{let z=M.events?.eventBuckets?.[R]?.entries?.[`${O}`];return z?.count?[(E?.[0]??0)+z.latency,(E?.[1]??0)+z.count]:E},null);x.push(T?T[0]/T[1]:null)}r.push({yAxisID:"yLatency",label:M.queue+": "+this.transloco.translate("dashboard.queues.latency"),data:x,borderColor:t["tertiary-50"],pointHoverBackgroundColor:t["tertiary-80"],pointHoverBorderColor:t["tertiary-20"]})}}}}return{type:"line",options:{animation:!1,responsive:!0,elements:{line:{tension:.5}},scales:{yCount:{position:"left",ticks:{callback:k=>parseInt(k).toLocaleString(this.transloco.getActiveLang())}},yLatency:{position:"right",ticks:{callback:this.formatDuration.bind(this)}}},plugins:{legend:{display:i.legend},decimation:{enabled:!0},tooltip:{callbacks:{label:k=>k.dataset.yAxisID==="yCount"?k.formattedValue:this.formatDuration(k.parsed.y)}}}},data:{labels:a,datasets:r}}}formatBucketKey(e,i){let t;switch(e.duration){case"day":t="d LLL";break;case"hour":t="d LLL H:00";break;case"minute":t="H:mm";break}return Fe(1e3*te[e.duration]*e.multiplier*i,t,{locale:Re(this.transloco.getActiveLang())})}formatDuration(e){if(typeof e=="string"&&(e=parseInt(e)),e===0)return"0";let i=e,t=0,a=0,r=0;return i>=60&&(t=Math.floor(i/60),i=i%60,t>=5&&(i=0,t>=60&&(a=Math.floor(t/60),t=t%60,a>=5&&(t=0,a>=24&&(r=Math.floor(a/24),a=a%24))))),He({days:r,hours:a,minutes:t,seconds:i},this.transloco.getActiveLang())}static{this.\u0275fac=function(i){return new(i||n)}}static{this.\u0275prov=H({token:n,factory:n.\u0275fac,providedIn:"root"})}}return n})();var Ct=(n,o,e)=>[n,o,e];function kt(n,o){if(n&1&&(s(0,"mat-option",7),f(1),u()),n&2){let e=o.$implicit,i=d().$implicit;l("value",e),c(),V(i("dashboard.interval."+e))}}function vt(n,o){if(n&1&&(s(0,"mat-option",7),f(1),u()),n&2){let e=o.$implicit,i=d().$implicit;l("value",e),c(),V(i("dashboard.interval."+e+"s"))}}function Mt(n,o){if(n&1&&(s(0,"mat-option",7),f(1),u()),n&2){let e=o.$implicit;l("value",e),c(),V(e)}}function xt(n,o){if(n&1){let e=ae();s(0,"button",17),C("click",function(){let t=b(e).$implicit,a=d(2);return g(a.queueMetricsController.params.queue===t||a.queueMetricsController.setQueue(t))}),s(1,"mat-icon"),f(2),u()()}if(n&2){let e=o.$implicit,i=d(2);N(i.queueMetricsController.params.queue===e?"selected":"deselected"),l("matTooltip",e),c(2),V(i.queueMetricsController.params.queue===e?"radio_button_checked":"radio_button_unchecked")}}function yt(n,o){if(n&1&&(s(0,"mat-option",7),f(1),u()),n&2){let e=o.$implicit,i=d().$implicit;l("value",e),c(),V(i("dashboard.event."+e))}}function St(n,o){if(n&1&&(s(0,"mat-option",7),f(1),u()),n&2){let e=o.$implicit,i=d().$implicit;l("value",e),c(),V(i("dashboard.interval."+e))}}function wt(n,o){if(n&1){let e=ae();_e(0),$(1,"app-document-title",1),s(2,"mat-card")(3,"mat-card-content")(4,"mat-grid-list",2)(5,"mat-grid-tile",3)(6,"mat-card",4)(7,"mat-card-header")(8,"mat-card-title")(9,"h4"),f(10),u()()(),s(11,"mat-card-content")(12,"mat-form-field",5)(13,"mat-select",6),C("valueChange",function(t){b(e);let a=d();return g(a.queueMetricsController.setTimeframe(t))}),D(14,kt,2,2,"mat-option",7,A),u()(),s(16,"div",8)(17,"button",9),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setTimeframe(t.timeframeNames[0]))}),s(18,"mat-icon"),f(19,"first_page"),u()(),s(20,"button",9),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setTimeframe(t.timeframeNames[t.timeframeNames.indexOf(t.queueMetricsController.params.buckets.timeframe)-1]))}),s(21,"mat-icon"),f(22,"navigate_before"),u()(),s(23,"button",9),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setTimeframe(t.timeframeNames[t.timeframeNames.indexOf(t.queueMetricsController.params.buckets.timeframe)+1]))}),s(24,"mat-icon"),f(25,"navigate_next"),u()(),s(26,"button",9),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setTimeframe(t.timeframeNames[t.timeframeNames.length-1]))}),s(27,"mat-icon"),f(28,"last_page"),u()()()()()(),s(29,"mat-grid-tile",3)(30,"mat-card",10)(31,"mat-card-header")(32,"mat-card-title")(33,"h4"),f(34),u()()(),s(35,"mat-card-content")(36,"mat-form-field",11)(37,"input",12),ke(38,"async"),C("change",function(t){b(e);let a=d();return g(a.handleMultiplierEvent(t))}),u()(),s(39,"mat-form-field",13)(40,"mat-select",6),C("valueChange",function(t){b(e);let a=d();return g(a.queueMetricsController.setBucketDuration(t))}),D(41,vt,2,2,"mat-option",7,A),u()(),s(43,"div",8)(44,"button",9),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setBucketMultiplier(t.queueMetricsController.bucketMultiplier-1))}),s(45,"mat-icon"),f(46,"remove"),u()(),s(47,"button",14),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setBucketMultiplier(t.queueMetricsController.bucketMultiplier+1))}),s(48,"mat-icon"),f(49,"add"),u()(),s(50,"button",9),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setBucketDuration(t.resolutionNames[0]))}),s(51,"mat-icon"),f(52,"first_page"),u()(),s(53,"button",9),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setBucketDuration(t.resolutionNames[t.resolutionNames.indexOf(t.queueMetricsController.bucketDuration)-1]))}),s(54,"mat-icon"),f(55,"navigate_before"),u()(),s(56,"button",9),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setBucketDuration(t.resolutionNames[t.resolutionNames.indexOf(t.queueMetricsController.bucketDuration)+1]))}),s(57,"mat-icon"),f(58,"navigate_next"),u()(),s(59,"button",9),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setBucketDuration(t.resolutionNames[t.resolutionNames.length-1]))}),s(60,"mat-icon"),f(61,"last_page"),u()()()()()(),s(62,"mat-grid-tile",3)(63,"mat-card")(64,"mat-card-header")(65,"mat-card-title")(66,"h4"),f(67),u()()(),s(68,"mat-card-content")(69,"mat-form-field",5)(70,"mat-select",6),C("valueChange",function(t){b(e);let a=d();return g(a.queueMetricsController.setQueue(t==="_all"?null:t))}),s(71,"mat-option",15),f(72),u(),D(73,Mt,2,2,"mat-option",7,A),u()(),s(75,"div",16)(76,"button",17),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setQueue(null))}),s(77,"mat-icon",18),f(78,"workspaces"),u()(),D(79,xt,3,4,"button",19,A),u()()()(),s(81,"mat-grid-tile",3)(82,"mat-card")(83,"mat-card-header")(84,"mat-card-title")(85,"h4"),f(86),u()()(),s(87,"mat-card-content")(88,"mat-form-field",5)(89,"mat-select",6),C("valueChange",function(t){b(e);let a=d();return g(a.queueMetricsController.setEvent(t==="_all"?null:t))}),s(90,"mat-option",15),f(91,"All"),u(),D(92,yt,2,2,"mat-option",7,A),u()(),s(94,"div",16)(95,"button",17),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.setEvent(null))}),s(96,"mat-icon",18),f(97,"radio_button_checked"),u()(),s(98,"button",17),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.params.event==="created"||t.queueMetricsController.setEvent("created"))}),s(99,"mat-icon"),f(100,"add_circle"),u()(),s(101,"button",17),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.params.event==="processed"||t.queueMetricsController.setEvent("processed"))}),s(102,"mat-icon"),f(103,"check_circle"),u()(),s(104,"button",17),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.params.event==="failed"||t.queueMetricsController.setEvent("failed"))}),s(105,"mat-icon"),f(106,"error"),u()()()()()(),s(107,"mat-grid-tile",3)(108,"mat-card",20)(109,"mat-card-header")(110,"mat-card-title")(111,"h4"),f(112),u()()(),s(113,"mat-card-content")(114,"mat-form-field",5)(115,"mat-select",6),C("valueChange",function(t){b(e);let a=d();return g(a.queueMetricsController.setAutoRefreshInterval(t))}),D(116,St,2,2,"mat-option",7,A),u()(),s(118,"div",16)(119,"button",17),C("click",function(){b(e);let t=d();return g(t.queueMetricsController.refresh())}),s(120,"mat-icon"),f(121,"sync"),u()()()()()()(),s(122,"div",21),$(123,"mat-progress-bar",22),u(),s(124,"mat-grid-list",2)(125,"mat-grid-tile",3),$(126,"app-chart",23),u(),s(127,"mat-grid-tile",3),$(128,"app-chart",23),u()()()(),he()}if(n&2){let e,i,t,a=o.$implicit,r=d();c(),l("parts",Ce(69,Ct,a("routes.visualize"),a("routes.queues"),a("routes.dashboard"))),c(3),l("cols",r.breakpoints.sizeAtLeast("Large")?5:r.breakpoints.sizeAtLeast("Medium")?3:r.breakpoints.sizeAtLeast("Small")?2:1),c(),l("colspan",1)("rowspan",2),c(5),V(a("dashboard.metrics.timeframe")),c(3),l("value",r.queueMetricsController.params.buckets.timeframe),c(),P(r.timeframeNames),c(3),l("disabled",r.timeframeNames.indexOf(r.queueMetricsController.params.buckets.timeframe)<=0),c(3),l("disabled",r.timeframeNames.indexOf(r.queueMetricsController.params.buckets.timeframe)<=0),c(3),l("disabled",r.timeframeNames.indexOf(r.queueMetricsController.params.buckets.timeframe)>=r.timeframeNames.length-1),c(3),l("disabled",r.timeframeNames.indexOf(r.queueMetricsController.params.buckets.timeframe)>=r.timeframeNames.length-1),c(3),l("colspan",1)("rowspan",2),c(5),be(" ",a("dashboard.metrics.resolution")," "),c(3),l("placeholder",(e=(e=ve(38,67,r.queueMetricsController.result$))==null||e.params==null||e.params.buckets==null||e.params.buckets.multiplier==null?null:e.params.buckets.multiplier.toString())!==null&&e!==void 0?e:"")("value",r.queueMetricsController.params.buckets.multiplier),c(3),l("value",r.queueMetricsController.bucketDuration),c(),P(r.resolutionNames),c(3),l("disabled",r.queueMetricsController.bucketMultiplier===1),c(6),l("disabled",r.resolutionNames.indexOf(r.queueMetricsController.bucketDuration)<=0),c(3),l("disabled",r.resolutionNames.indexOf(r.queueMetricsController.bucketDuration)<=0),c(3),l("disabled",r.resolutionNames.indexOf(r.queueMetricsController.bucketDuration)>=r.resolutionNames.length-1),c(3),l("disabled",r.resolutionNames.indexOf(r.queueMetricsController.bucketDuration)>=r.resolutionNames.length-1),c(3),l("colspan",1)("rowspan",2),c(5),V(a("dashboard.queues.queue")),c(3),l("value",(i=r.queueMetricsController.params.queue)!==null&&i!==void 0?i:"_all"),c(2),V(a("general.all")),c(),P(r.availableQueueNames),c(3),N(r.queueMetricsController.params.queue?"deselected":"selected"),l("matTooltip",a("general.all")),c(3),P(r.availableQueueNames),c(2),l("colspan",1)("rowspan",2),c(5),V(a("dashboard.metrics.event")),c(3),l("value",(t=r.queueMetricsController.params.event)!==null&&t!==void 0?t:"_all"),c(3),P(r.eventNames),c(3),N(r.queueMetricsController.params.event?"deselected":"selected"),l("matTooltip",a("general.all")),c(3),N(r.queueMetricsController.params.event==="created"?"selected":"deselected"),l("matTooltip",a("dashboard.queues.created")),c(3),N(r.queueMetricsController.params.event==="processed"?"selected":"deselected"),l("matTooltip",a("dashboard.queues.processed")),c(3),N(r.queueMetricsController.params.event==="failed"?"selected":"deselected"),l("matTooltip",a("dashboard.queues.failed")),c(3),l("colspan",1)("rowspan",2),c(5),V(a("general.refresh")),c(3),l("value",r.queueMetricsController.params.autoRefresh),c(),P(r.autoRefreshIntervalNames),c(3),l("matTooltip",a("general.refresh")),c(4),l("mode",r.queueMetricsController.loading?"indeterminate":"determinate")("value",0),c(),l("cols",r.breakpoints.sizeAtLeast("Large")?2:1),c(),l("colspan",1)("rowspan",5),c(),l("title",a("dashboard.queues.total_counts_by_status"))("adapter",r.totals)("$data",r.queueMetricsController.result$)("height",400)("width",550),c(),l("colspan",1)("rowspan",5),c(),l("title",a("dashboard.metrics.throughput"))("adapter",r.timeline)("$data",r.queueMetricsController.result$)("height",400)("width",550)}}var Di=(()=>{class n{constructor(){this.breakpoints=q(Le),this.apollo=q(ye),this.queueMetricsController=new ne(this.apollo,{buckets:{duration:"AUTO",multiplier:"AUTO",timeframe:"all"},autoRefresh:"seconds_30"},q(Te)),this.timeline=q(lt),this.totals=q(it),this.resolutionNames=Ge,this.timeframeNames=We,this.availableQueueNames=Xe,this.autoRefreshIntervalNames=Ze,this.eventNames=ie}ngOnInit(){this.queueMetricsController.result$.subscribe(e=>{if(this.queueMetricsController.params.buckets.timeframe==="all"&&this.queueMetricsController.params.buckets.duration==="AUTO"&&e.params.buckets.duration==="hour"){let i=e.bucketSpan;i&&i.latestBucket-i.earliestBucket<12&&this.queueMetricsController.setBucketDuration("minute")}})}ngOnDestroy(){this.queueMetricsController.setAutoRefreshInterval("off")}handleMultiplierEvent(e){let i=e.currentTarget.value;this.queueMetricsController.setBucketMultiplier(/^\d+$/.test(i)?parseInt(i):"AUTO")}static{this.\u0275fac=function(i){return new(i||n)}}static{this.\u0275cmp=pe({type:n,selectors:[["app-queue-visualize"]],standalone:!0,features:[ge],decls:1,vars:0,consts:[[4,"transloco"],[3,"parts"],["rowHeight","100px",3,"cols"],[3,"colspan","rowspan"],[1,"form-timeframe"],["subscriptSizing","dynamic"],[3,"valueChange","value"],[3,"value"],[1,"paginator","actions"],["mat-icon-button","",3,"click","disabled"],[1,"form-resolution"],["subscriptSizing","dynamic",1,"form-input-multiplier"],["type","number","matInput","","min","1","step","1",3,"change","placeholder","value"],["subscriptSizing","dynamic",1,"form-select-duration"],["mat-icon-button","",3,"click"],["value","_all"],[1,"actions"],["mat-icon-button","",3,"click","matTooltip"],["fontSet","material-icons"],["mat-icon-button","",3,"class","matTooltip"],[1,"form-refresh"],[1,"progress-bar-container"],[3,"mode","value"],[3,"title","adapter","$data","height","width"]],template:function(i,t){i&1&&fe(0,wt,129,73,"ng-container",0)},dependencies:[je,Se,G,Y,Oe,J,Be,Ve,Ie,W,X,ze,De,Qe,Pe,xe,Me,ee,K,tt,$e],styles:[".actions[_ngcontent-%COMP%]{width:210px;padding-top:12px;--mdc-icon-button-state-layer-size: 32px}.actions[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%]{font-size:22px}.actions[_ngcontent-%COMP%]   button[_ngcontent-%COMP%]{margin-right:0}.progress-bar-container[_ngcontent-%COMP%]{width:100%;height:10px}mat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]{width:100%}mat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]   mat-card-content[_ngcontent-%COMP%]{min-width:190px}mat-grid-tile[_ngcontent-%COMP%]   mat-card[_ngcontent-%COMP%]   h4[_ngcontent-%COMP%]{margin-bottom:16px;font-size:18px}mat-form-field[_ngcontent-%COMP%]{width:186px}.form-resolution[_ngcontent-%COMP%]   .actions[_ngcontent-%COMP%]{margin-left:-2px}.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]{width:60px;margin-right:10px}.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[_ngcontent-%COMP%]::-webkit-outer-spin-button, .form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[_ngcontent-%COMP%]::-webkit-inner-spin-button{-webkit-appearance:none;margin:0}.form-resolution[_ngcontent-%COMP%]   .form-input-multiplier[_ngcontent-%COMP%]   input[type=number][_ngcontent-%COMP%]{-moz-appearance:textfield}.form-resolution[_ngcontent-%COMP%]   .form-select-duration[_ngcontent-%COMP%]{width:116px}"]})}}return n})();export{Di as QueueVisualizeComponent};
