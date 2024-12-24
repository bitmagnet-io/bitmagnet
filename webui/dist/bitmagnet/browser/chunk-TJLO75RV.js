import{c as _,e as lt}from"./chunk-Y7K23DTG.js";import{a as C,b as ft,e as f,f as s,g as x,h as l,i as F,j as mt,k as ht}from"./chunk-42PJPEMD.js";import{a as dt}from"./chunk-MSAOOVCY.js";import{P as nt,Q as at,R as ot,U as it,Ua as ut,a as tt,b as et,qa as st,r as rt,ya as ct}from"./chunk-VAEZNV34.js";import{$b as v,Cb as $,Ea as G,Fa as B,Hb as W,Kb as P,Qb as y,Rb as D,Sb as X,Tb as R,Ub as V,Wb as j,Zb as A,a as L,b as N,ec as Z,fc as J,gc as K,h as H,kb as O,kc as E,mc as U,pc as z,qa as k,ua as Q}from"./chunk-Z3WUIYN5.js";function Et(r,e){if(r&1){let t=j();R(0),y(1,"mat-card")(2,"mat-card-header")(3,"mat-card-title")(4,"button",1),A("click",function(){G(t);let a=v();return B(a.toggleLegend())}),y(5,"mat-icon",2),E(6,"legend_toggle"),D()(),y(7,"h4"),E(8),D()()(),y(9,"mat-card-content")(10,"div"),X(11,"canvas",3),D()()(),V()}if(r&2){let t=e.$implicit,n=v();O(4),P(n.chartShowLegend?"selected":"deselected"),W("matTooltip",t("dashboard.metrics.toggle_chart_legend")),O(4),U(" ",n.title," "),O(2),P(n.breakpoints.sizeAtLeast("Large")?"app-chart":"app-chart-small"),O(),W("data",n.chartConfig.data)("options",n.chartConfig.options)("type",n.chartConfig.type)("height",n.height)("width",n.width)}}var ie=(()=>{class r{constructor(){this.themeInfo=k(lt),this.transloco=k(tt),this.breakpoints=k(dt),this.hiddenDatasets=new Map,this.$data=new H,this.width=500,this.height=500,this.title="chart title",this.chartShowLegend=!0}ngOnInit(){this.updateChart(),this.$data.subscribe(t=>{this.data=t,this.updateChart()}),this.themeInfo.info$.subscribe(()=>{this.updateChart()}),this.transloco.langChanges$.subscribe(()=>{this.updateChart()})}toggleLegend(){this.chartShowLegend=!this.chartShowLegend,this.legendConfig()}setDatasetVisible(){for(let t=0;t<this.basechart.chart.data.datasets.length;t++){let n=this.basechart.chart.getDatasetMeta(t);this.hiddenDatasets.has(n.label)&&(n.hidden=this.hiddenDatasets.get(n.label))}}legendOnClick(t,n){let a=this.basechart.chart.getDatasetMeta(n.datasetIndex);this.hiddenDatasets.set(a.label,!a.hidden),this.setDatasetVisible(),this.basechart.chart.update()}legendConfig(){setTimeout(()=>{this.basechart.chart.legend.options.display=this.chartShowLegend,this.basechart.chart.options.plugins.legend.onClick=this.legendOnClick.bind(this),this.setDatasetVisible(),this.basechart.chart.update()},5)}updateChart(){this.chartConfig=this.adapter.create(this.data),this.legendConfig()}static{this.\u0275fac=function(n){return new(n||r)}}static{this.\u0275cmp=Q({type:r,selectors:[["app-chart"]],viewQuery:function(n,a){if(n&1&&Z(_,5),n&2){let o;J(o=K())&&(a.basechart=o.first)}},inputs:{$data:"$data",adapter:"adapter",width:"width",height:"height",title:"title"},standalone:!0,features:[z],decls:1,vars:0,consts:[[4,"transloco"],["type","button","mat-icon-button","",3,"click","matTooltip"],["aria-label","Toggle legend","fontSet","material-icons",1,"icon-header"],["baseChart","",3,"data","options","type","height","width"]],template:function(n,a){n&1&&$(0,Et,12,11,"ng-container",0)},dependencies:[ut,rt,nt,ot,it,at,st,ct,et,_],styles:[".app-chart[_ngcontent-%COMP%]{position:relative;min-height:30vh;width:35vw}.app-chart-small[_ngcontent-%COMP%]{position:relative;width:calc(40vw + 100px)}button[_ngcontent-%COMP%]{vertical-align:middle}h4[_ngcontent-%COMP%]{margin:0;display:inline;vertical-align:middle;font-size:18px}"]})}}return r})();function I(r,e){let t=s(r,e?.in);return t.setHours(0,0,0,0),t}function gt(r,e,t){let[n,a]=mt(t?.in,r,e),o=I(n),c=I(a),m=+o-F(o),h=+c-F(c);return Math.round((m-h)/ft)}function pt(r,e){let t=s(r,e?.in);return t.setFullYear(t.getFullYear(),0,1),t.setHours(0,0,0,0),t}function wt(r,e){let t=s(r,e?.in);return gt(t,pt(t))+1}function p(r,e){return l(r,N(L({},e),{weekStartsOn:1}))}function M(r,e){let t=s(r,e?.in),n=t.getFullYear(),a=f(t,0);a.setFullYear(n+1,0,4),a.setHours(0,0,0,0);let o=p(a),c=f(t,0);c.setFullYear(n,0,4),c.setHours(0,0,0,0);let m=p(c);return t.getTime()>=o.getTime()?n+1:t.getTime()>=m.getTime()?n:n-1}function xt(r,e){let t=M(r,e),n=f(e?.in||r,0);return n.setFullYear(t,0,4),n.setHours(0,0,0,0),p(n)}function bt(r,e){let t=s(r,e?.in),n=+p(t)-+xt(t);return Math.round(n/C)+1}function Y(r,e){let t=s(r,e?.in),n=t.getFullYear(),a=x(),o=e?.firstWeekContainsDate??e?.locale?.options?.firstWeekContainsDate??a.firstWeekContainsDate??a.locale?.options?.firstWeekContainsDate??1,c=f(e?.in||r,0);c.setFullYear(n+1,0,o),c.setHours(0,0,0,0);let m=l(c,e),h=f(e?.in||r,0);h.setFullYear(n,0,o),h.setHours(0,0,0,0);let T=l(h,e);return+t>=+m?n+1:+t>=+T?n:n-1}function Ot(r,e){let t=x(),n=e?.firstWeekContainsDate??e?.locale?.options?.firstWeekContainsDate??t.firstWeekContainsDate??t.locale?.options?.firstWeekContainsDate??1,a=Y(r,e),o=f(e?.in||r,0);return o.setFullYear(a,0,n),o.setHours(0,0,0,0),l(o,e)}function yt(r,e){let t=s(r,e?.in),n=+l(t,e)-+Ot(t,e);return Math.round(n/C)+1}function i(r,e){let t=r<0?"-":"",n=Math.abs(r).toString().padStart(e,"0");return t+n}var g={y(r,e){let t=r.getFullYear(),n=t>0?t:1-t;return i(e==="yy"?n%100:n,e.length)},M(r,e){let t=r.getMonth();return e==="M"?String(t+1):i(t+1,2)},d(r,e){return i(r.getDate(),e.length)},a(r,e){let t=r.getHours()/12>=1?"pm":"am";switch(e){case"a":case"aa":return t.toUpperCase();case"aaa":return t;case"aaaaa":return t[0];case"aaaa":default:return t==="am"?"a.m.":"p.m."}},h(r,e){return i(r.getHours()%12||12,e.length)},H(r,e){return i(r.getHours(),e.length)},m(r,e){return i(r.getMinutes(),e.length)},s(r,e){return i(r.getSeconds(),e.length)},S(r,e){let t=e.length,n=r.getMilliseconds(),a=Math.trunc(n*Math.pow(10,t-3));return i(a,e.length)}};var b={am:"am",pm:"pm",midnight:"midnight",noon:"noon",morning:"morning",afternoon:"afternoon",evening:"evening",night:"night"},q={G:function(r,e,t){let n=r.getFullYear()>0?1:0;switch(e){case"G":case"GG":case"GGG":return t.era(n,{width:"abbreviated"});case"GGGGG":return t.era(n,{width:"narrow"});case"GGGG":default:return t.era(n,{width:"wide"})}},y:function(r,e,t){if(e==="yo"){let n=r.getFullYear(),a=n>0?n:1-n;return t.ordinalNumber(a,{unit:"year"})}return g.y(r,e)},Y:function(r,e,t,n){let a=Y(r,n),o=a>0?a:1-a;if(e==="YY"){let c=o%100;return i(c,2)}return e==="Yo"?t.ordinalNumber(o,{unit:"year"}):i(o,e.length)},R:function(r,e){let t=M(r);return i(t,e.length)},u:function(r,e){let t=r.getFullYear();return i(t,e.length)},Q:function(r,e,t){let n=Math.ceil((r.getMonth()+1)/3);switch(e){case"Q":return String(n);case"QQ":return i(n,2);case"Qo":return t.ordinalNumber(n,{unit:"quarter"});case"QQQ":return t.quarter(n,{width:"abbreviated",context:"formatting"});case"QQQQQ":return t.quarter(n,{width:"narrow",context:"formatting"});case"QQQQ":default:return t.quarter(n,{width:"wide",context:"formatting"})}},q:function(r,e,t){let n=Math.ceil((r.getMonth()+1)/3);switch(e){case"q":return String(n);case"qq":return i(n,2);case"qo":return t.ordinalNumber(n,{unit:"quarter"});case"qqq":return t.quarter(n,{width:"abbreviated",context:"standalone"});case"qqqqq":return t.quarter(n,{width:"narrow",context:"standalone"});case"qqqq":default:return t.quarter(n,{width:"wide",context:"standalone"})}},M:function(r,e,t){let n=r.getMonth();switch(e){case"M":case"MM":return g.M(r,e);case"Mo":return t.ordinalNumber(n+1,{unit:"month"});case"MMM":return t.month(n,{width:"abbreviated",context:"formatting"});case"MMMMM":return t.month(n,{width:"narrow",context:"formatting"});case"MMMM":default:return t.month(n,{width:"wide",context:"formatting"})}},L:function(r,e,t){let n=r.getMonth();switch(e){case"L":return String(n+1);case"LL":return i(n+1,2);case"Lo":return t.ordinalNumber(n+1,{unit:"month"});case"LLL":return t.month(n,{width:"abbreviated",context:"standalone"});case"LLLLL":return t.month(n,{width:"narrow",context:"standalone"});case"LLLL":default:return t.month(n,{width:"wide",context:"standalone"})}},w:function(r,e,t,n){let a=yt(r,n);return e==="wo"?t.ordinalNumber(a,{unit:"week"}):i(a,e.length)},I:function(r,e,t){let n=bt(r);return e==="Io"?t.ordinalNumber(n,{unit:"week"}):i(n,e.length)},d:function(r,e,t){return e==="do"?t.ordinalNumber(r.getDate(),{unit:"date"}):g.d(r,e)},D:function(r,e,t){let n=wt(r);return e==="Do"?t.ordinalNumber(n,{unit:"dayOfYear"}):i(n,e.length)},E:function(r,e,t){let n=r.getDay();switch(e){case"E":case"EE":case"EEE":return t.day(n,{width:"abbreviated",context:"formatting"});case"EEEEE":return t.day(n,{width:"narrow",context:"formatting"});case"EEEEEE":return t.day(n,{width:"short",context:"formatting"});case"EEEE":default:return t.day(n,{width:"wide",context:"formatting"})}},e:function(r,e,t,n){let a=r.getDay(),o=(a-n.weekStartsOn+8)%7||7;switch(e){case"e":return String(o);case"ee":return i(o,2);case"eo":return t.ordinalNumber(o,{unit:"day"});case"eee":return t.day(a,{width:"abbreviated",context:"formatting"});case"eeeee":return t.day(a,{width:"narrow",context:"formatting"});case"eeeeee":return t.day(a,{width:"short",context:"formatting"});case"eeee":default:return t.day(a,{width:"wide",context:"formatting"})}},c:function(r,e,t,n){let a=r.getDay(),o=(a-n.weekStartsOn+8)%7||7;switch(e){case"c":return String(o);case"cc":return i(o,e.length);case"co":return t.ordinalNumber(o,{unit:"day"});case"ccc":return t.day(a,{width:"abbreviated",context:"standalone"});case"ccccc":return t.day(a,{width:"narrow",context:"standalone"});case"cccccc":return t.day(a,{width:"short",context:"standalone"});case"cccc":default:return t.day(a,{width:"wide",context:"standalone"})}},i:function(r,e,t){let n=r.getDay(),a=n===0?7:n;switch(e){case"i":return String(a);case"ii":return i(a,e.length);case"io":return t.ordinalNumber(a,{unit:"day"});case"iii":return t.day(n,{width:"abbreviated",context:"formatting"});case"iiiii":return t.day(n,{width:"narrow",context:"formatting"});case"iiiiii":return t.day(n,{width:"short",context:"formatting"});case"iiii":default:return t.day(n,{width:"wide",context:"formatting"})}},a:function(r,e,t){let a=r.getHours()/12>=1?"pm":"am";switch(e){case"a":case"aa":return t.dayPeriod(a,{width:"abbreviated",context:"formatting"});case"aaa":return t.dayPeriod(a,{width:"abbreviated",context:"formatting"}).toLowerCase();case"aaaaa":return t.dayPeriod(a,{width:"narrow",context:"formatting"});case"aaaa":default:return t.dayPeriod(a,{width:"wide",context:"formatting"})}},b:function(r,e,t){let n=r.getHours(),a;switch(n===12?a=b.noon:n===0?a=b.midnight:a=n/12>=1?"pm":"am",e){case"b":case"bb":return t.dayPeriod(a,{width:"abbreviated",context:"formatting"});case"bbb":return t.dayPeriod(a,{width:"abbreviated",context:"formatting"}).toLowerCase();case"bbbbb":return t.dayPeriod(a,{width:"narrow",context:"formatting"});case"bbbb":default:return t.dayPeriod(a,{width:"wide",context:"formatting"})}},B:function(r,e,t){let n=r.getHours(),a;switch(n>=17?a=b.evening:n>=12?a=b.afternoon:n>=4?a=b.morning:a=b.night,e){case"B":case"BB":case"BBB":return t.dayPeriod(a,{width:"abbreviated",context:"formatting"});case"BBBBB":return t.dayPeriod(a,{width:"narrow",context:"formatting"});case"BBBB":default:return t.dayPeriod(a,{width:"wide",context:"formatting"})}},h:function(r,e,t){if(e==="ho"){let n=r.getHours()%12;return n===0&&(n=12),t.ordinalNumber(n,{unit:"hour"})}return g.h(r,e)},H:function(r,e,t){return e==="Ho"?t.ordinalNumber(r.getHours(),{unit:"hour"}):g.H(r,e)},K:function(r,e,t){let n=r.getHours()%12;return e==="Ko"?t.ordinalNumber(n,{unit:"hour"}):i(n,e.length)},k:function(r,e,t){let n=r.getHours();return n===0&&(n=24),e==="ko"?t.ordinalNumber(n,{unit:"hour"}):i(n,e.length)},m:function(r,e,t){return e==="mo"?t.ordinalNumber(r.getMinutes(),{unit:"minute"}):g.m(r,e)},s:function(r,e,t){return e==="so"?t.ordinalNumber(r.getSeconds(),{unit:"second"}):g.s(r,e)},S:function(r,e){return g.S(r,e)},X:function(r,e,t){let n=r.getTimezoneOffset();if(n===0)return"Z";switch(e){case"X":return Dt(n);case"XXXX":case"XX":return w(n);case"XXXXX":case"XXX":default:return w(n,":")}},x:function(r,e,t){let n=r.getTimezoneOffset();switch(e){case"x":return Dt(n);case"xxxx":case"xx":return w(n);case"xxxxx":case"xxx":default:return w(n,":")}},O:function(r,e,t){let n=r.getTimezoneOffset();switch(e){case"O":case"OO":case"OOO":return"GMT"+kt(n,":");case"OOOO":default:return"GMT"+w(n,":")}},z:function(r,e,t){let n=r.getTimezoneOffset();switch(e){case"z":case"zz":case"zzz":return"GMT"+kt(n,":");case"zzzz":default:return"GMT"+w(n,":")}},t:function(r,e,t){let n=Math.trunc(+r/1e3);return i(n,e.length)},T:function(r,e,t){return i(+r,e.length)}};function kt(r,e=""){let t=r>0?"-":"+",n=Math.abs(r),a=Math.trunc(n/60),o=n%60;return o===0?t+String(a):t+String(a)+e+i(o,2)}function Dt(r,e){return r%60===0?(r>0?"-":"+")+i(Math.abs(r)/60,2):w(r,e)}function w(r,e=""){let t=r>0?"-":"+",n=Math.abs(r),a=i(Math.trunc(n/60),2),o=i(n%60,2);return t+a+e+o}var Ct=(r,e)=>{switch(r){case"P":return e.date({width:"short"});case"PP":return e.date({width:"medium"});case"PPP":return e.date({width:"long"});case"PPPP":default:return e.date({width:"full"})}},Mt=(r,e)=>{switch(r){case"p":return e.time({width:"short"});case"pp":return e.time({width:"medium"});case"ppp":return e.time({width:"long"});case"pppp":default:return e.time({width:"full"})}},_t=(r,e)=>{let t=r.match(/(P+)(p+)?/)||[],n=t[1],a=t[2];if(!a)return Ct(r,e);let o;switch(n){case"P":o=e.dateTime({width:"short"});break;case"PP":o=e.dateTime({width:"medium"});break;case"PPP":o=e.dateTime({width:"long"});break;case"PPPP":default:o=e.dateTime({width:"full"});break}return o.replace("{{date}}",Ct(n,e)).replace("{{time}}",Mt(a,e))},Yt={p:Mt,P:_t};var Ft=/^D+$/,It=/^Y+$/,qt=["D","DD","YY","YYYY"];function Tt(r){return Ft.test(r)}function St(r){return It.test(r)}function Wt(r,e,t){let n=Lt(r,e,t);if(console.warn(n),qt.includes(r))throw new RangeError(n)}function Lt(r,e,t){let n=r[0]==="Y"?"years":"days of the month";return`Use \`${r.toLowerCase()}\` instead of \`${r}\` (in \`${e}\`) for formatting ${n} to the input \`${t}\`; see: https://github.com/date-fns/date-fns/blob/master/docs/unicodeTokens.md`}function Pt(r){return r instanceof Date||typeof r=="object"&&Object.prototype.toString.call(r)==="[object Date]"}function vt(r){return!(!Pt(r)&&typeof r!="number"||isNaN(+s(r)))}var Nt=/[yYQqMLwIdDecihHKkms]o|(\w)\1*|''|'(''|[^'])+('|$)|./g,Ht=/P+p+|P+|p+|''|'(''|[^'])+('|$)|./g,Qt=/^'([^]*?)'?$/,Gt=/''/g,Bt=/[a-zA-Z]/;function kr(r,e,t){let n=x(),a=t?.locale??n.locale??ht,o=t?.firstWeekContainsDate??t?.locale?.options?.firstWeekContainsDate??n.firstWeekContainsDate??n.locale?.options?.firstWeekContainsDate??1,c=t?.weekStartsOn??t?.locale?.options?.weekStartsOn??n.weekStartsOn??n.locale?.options?.weekStartsOn??0,m=s(r,t?.in);if(!vt(m))throw new RangeError("Invalid time value");let h=e.match(Ht).map(d=>{let u=d[0];if(u==="p"||u==="P"){let S=Yt[u];return S(d,a.formatLong)}return d}).join("").match(Nt).map(d=>{if(d==="''")return{isToken:!1,value:"'"};let u=d[0];if(u==="'")return{isToken:!1,value:$t(d)};if(q[u])return{isToken:!0,value:d};if(u.match(Bt))throw new RangeError("Format string contains an unescaped latin alphabet character `"+u+"`");return{isToken:!1,value:d}});a.localize.preprocessor&&(h=a.localize.preprocessor(m,h));let T={firstWeekContainsDate:o,weekStartsOn:c,locale:a};return h.map(d=>{if(!d.isToken)return d.value;let u=d.value;(!t?.useAdditionalWeekYearTokens&&St(u)||!t?.useAdditionalDayOfYearTokens&&Tt(u))&&Wt(u,e,String(r));let S=q[u[0]];return S(m,u,a.localize,T)}).join("")}function $t(r){let e=r.match(Qt);return e?e[1].replace(Gt,"'"):r}var Cr=(r,e)=>`${r}-${e}`;export{kr as a,ie as b,Cr as c};
