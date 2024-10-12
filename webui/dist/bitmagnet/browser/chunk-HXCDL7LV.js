import{a as Rt,b as he}from"./chunk-OS3AV5XL.js";import{a as Z}from"./chunk-AMUSMSN6.js";import{a as Nt}from"./chunk-BRPPGP57.js";import{Aa as me,B as Vt,Ba as _e,C as Lt,Ca as de,D as Ht,Da as ue,E as Gt,Ea as Ce,Fa as fe,G as Yt,I as Kt,L as qt,M as Jt,N as Ut,O as Xt,P as Qt,Pa as xe,Q as Zt,R as Wt,Sa as ge,Ua as V,_ as te,a as R,b as j,d as Pt,e as At,f as Ft,fa as W,g as Bt,i as Dt,ka as ee,n as Ot,na as ne,o as zt,sa as ie,ta as oe,ua as re,va as ae,wa as le,xa as se,ya as ce,z as jt,za as pe}from"./chunk-EH7PJXWD.js";import{m as kt}from"./chunk-FKMTSCBK.js";import{$b as c,$c as Q,B as at,Cb as C,Ea as M,Fa as $,Hb as g,Mb as d,N as Y,Na as St,Nb as K,Ob as k,Pb as P,Q as bt,Qb as s,Rb as r,Sb as h,Tb as T,Ub as b,Wb as B,Zb as A,_c as Mt,a as U,b as Tt,ha as X,hb as Et,hd as $t,jc as st,kb as o,kc as l,l as S,lb as lt,lc as x,mc as m,nc as wt,o as G,pc as N,qa as y,rc as It,sc as yt,ua as z,uc as E,vc as w,xa as vt}from"./chunk-3DR3CJRN.js";var De="array",Oe="bit",Te="bits",ze="byte",be="bytes",L="",Ne="exponent",Re="function",ve="iec",je="Invalid number",Ve="Invalid rounding method",ct="jedec",Le="object",Se=".",He="round",Ge="s",Ye="si",Ke="kbit",qe="kB",Je=" ",Ue="string",Xe="0",pt={symbol:{iec:{bits:["bit","Kibit","Mibit","Gibit","Tibit","Pibit","Eibit","Zibit","Yibit"],bytes:["B","KiB","MiB","GiB","TiB","PiB","EiB","ZiB","YiB"]},jedec:{bits:["bit","Kbit","Mbit","Gbit","Tbit","Pbit","Ebit","Zbit","Ybit"],bytes:["B","KB","MB","GB","TB","PB","EB","ZB","YB"]}},fullform:{iec:["","kibi","mebi","gibi","tebi","pebi","exbi","zebi","yobi"],jedec:["","kilo","mega","giga","tera","peta","exa","zetta","yotta"]}};function Ee(t,{bits:n=!1,pad:e=!1,base:i=-1,round:a=2,locale:p=L,localeOptions:_={},separator:D=L,spacer:v=Je,symbols:ke={},standard:I=L,output:ot=Ue,fullform:Pe=!1,fullforms:_t=[],exponent:dt=-1,roundingMethod:Ae=He,precision:q=0}={}){let f=dt,O=Number(t),u=[],F=0,rt=L;I===Ye?(i=10,I=ct):I===ve||I===ct?i=2:i===2?I=ve:(i=10,I=ct);let J=i===10?1e3:1024,Fe=Pe===!0,ut=O<0,Ct=Math[Ae];if(typeof t!="bigint"&&isNaN(t))throw new TypeError(je);if(typeof Ct!==Re)throw new TypeError(Ve);if(ut&&(O=-O),(f===-1||isNaN(f))&&(f=Math.floor(Math.log(O)/Math.log(J)),f<0&&(f=0)),f>8&&(q>0&&(q+=8-f),f=8),ot===Ne)return f;if(O===0)u[0]=0,rt=u[1]=pt.symbol[I][n?Te:be][f];else{F=O/(i===2?Math.pow(2,f*10):Math.pow(1e3,f)),n&&(F=F*8,F>=J&&f<8&&(F=F/J,f++));let H=Math.pow(10,f>0?a:0);u[0]=Ct(F*H)/H,u[0]===J&&f<8&&dt===-1&&(u[0]=1,f++),rt=u[1]=i===10&&f===1?n?Ke:qe:pt.symbol[I][n?Te:be][f]}if(ut&&(u[0]=-u[0]),q>0&&(u[0]=u[0].toPrecision(q)),u[1]=ke[u[1]]||u[1],p===!0?u[0]=u[0].toLocaleString():p.length>0?u[0]=u[0].toLocaleString(p,_):D.length>0&&(u[0]=u[0].toString().replace(Se,D)),e&&a>0){let H=u[0].toString(),ft=D||(H.match(/(\D)/g)||[]).pop()||Se,xt=H.toString().split(ft),gt=xt[1]||L,ht=gt.length,Be=a-ht;u[0]=`${xt[0]}${ft}${gt.padEnd(ht+Be,Xe)}`}return Fe&&(u[1]=_t[f]?_t[f]:pt.fullform[I][f]+(n?Oe:ze)+(u[0]===1?L:Ge)),ot===De?u:ot===Le?{value:u[0],symbol:u[1],exponent:f,unit:rt}:u.join(v)}var tt=(()=>{class t{constructor(){this.transloco=y(R)}transform(e){return Ee(e,{locale:this.transloco.getActiveLang(),base:2})}static{this.\u0275fac=function(i){return new(i||t)}}static{this.\u0275pipe=vt({name:"filesize",type:t,pure:!1,standalone:!0})}}return t})();var Qe=t=>t.toLowerCase().replaceAll(/[^a-z0-9\-]/g,"-").replace(/^-+/,"").replaceAll(/-+/g,"-"),we=Qe;var Ze={items:[],hasNextPage:!1,totalCount:0,aggregations:{queue:[],status:[]}},et=class{constructor(n,e,i){this.apollo=n,this.errorsService=e,this.currentRequest=new S(0),this.loadingSubject=new S(!1),this.loading$=this.loadingSubject.asObservable(),this.result=Ze,this.resultSubject=new S(this.result),this.result$=this.resultSubject.asObservable(),this.items$=this.resultSubject.pipe(at(a=>a.items)),i.subscribe(a=>{this.loadResult(a)}),this.resultSubject.subscribe(a=>{this.result=a})}connect({}){return this.items$}disconnect(){this.resultSubject.complete()}loadResult(n){this.currentSubscription&&(this.currentSubscription.unsubscribe(),this.currentSubscription=void 0),this.loadingSubject.next(!0);let e=this.currentRequest.getValue()+1;this.currentRequest.next(e);let i=this.apollo.query({query:xe,variables:n,fetchPolicy:"no-cache"}).pipe(at(a=>a.data.torrent.files)).pipe(Y(a=>(this.errorsService.addError(`Error loading item results: ${a.message}`),G)));this.currentSubscription=i.subscribe(a=>{e===this.currentRequest.getValue()&&(this.loadingSubject.next(!1),this.resultSubject.next(a))})}},nt=class{constructor(n){this.torrent=n,this.loading$=new S(!1).asObservable(),this.file={infoHash:n.infoHash,index:0,path:n.name,size:n.size,fileType:n.fileType,extension:n.extension,createdAt:n.createdAt,updatedAt:n.updatedAt},this.result={hasNextPage:!1,items:[this.file],totalCount:1},this.result$=new S(this.result).asObservable(),this.items$=new S([this.file]).asObservable()}connect({}){return this.items$}disconnect(){}};var it=class{constructor(n){let e={infoHash:n,limit:10,page:1};this.controlsSubject=new S(e),this.controls$=this.controlsSubject.asObservable(),this.controls$.pipe(bt(100)).subscribe(i=>{let a=this.variablesSubject.getValue(),p=Ie(i);JSON.stringify(a)!==JSON.stringify(p)&&this.variablesSubject.next(p)}),this.variablesSubject=new S(Ie(e)),this.variables$=this.variablesSubject.asObservable()}update(n){let e=this.controlsSubject.getValue(),i=n(e);JSON.stringify(e)!==JSON.stringify(i)&&this.controlsSubject.next(i)}handlePageEvent(n){this.update(e=>Tt(U({},e),{limit:n.pageSize,page:n.page}))}},Ie=t=>({input:{infoHashes:[t.infoHash],limit:t.limit,page:t.page,totalCount:!0,hasNextPage:!1}});var tn=(t,n)=>({x:t,y:n});function en(t,n){if(t&1&&(s(0,"p"),l(1),E(2,"number"),E(3,"number"),r()),t&2){let e=c().$implicit,i=c();o(),m(" ",e("torrents.showing_x_of_y_files",yt(5,tn,w(2,1,i.dataSource.result.totalCount),i.torrent.filesCount==null?"?":w(3,3,i.torrent.filesCount)))," ")}}function nn(t,n){if(t&1&&(s(0,"th",13),l(1),r()),t&2){let e=c().$implicit;o(),x(e("torrents.file_index"))}}function on(t,n){if(t&1&&(s(0,"td",14),l(1),r()),t&2){let e=n.$implicit,i=c(2);o(),m(" ",i.item(e).index," ")}}function rn(t,n){if(t&1&&(s(0,"th",13),l(1),r()),t&2){let e=c().$implicit;o(),x(e("torrents.file_path"))}}function an(t,n){if(t&1&&(s(0,"td",14),l(1),r()),t&2){let e=n.$implicit,i=c(2);o(),m(" ",i.item(e).path," ")}}function ln(t,n){if(t&1&&(s(0,"th",13),l(1),r()),t&2){let e=c().$implicit;o(),x(e("torrents.file_type"))}}function sn(t,n){if(t&1&&(s(0,"td",14),l(1),r()),t&2){let e,i=n.$implicit,a=c().$implicit,p=c();o(),m(" ",a("file_types."+((e=p.item(i).fileType)!==null&&e!==void 0?e:"unknown"))," ")}}function cn(t,n){if(t&1&&(s(0,"th",13),l(1),r()),t&2){let e=c().$implicit;o(),x(e("torrents.file_size"))}}function pn(t,n){if(t&1&&(s(0,"td",14),l(1),E(2,"filesize"),r()),t&2){let e=n.$implicit,i=c(2);o(),m(" ",w(2,1,i.item(e).size)," ")}}function mn(t,n){t&1&&h(0,"tr",15)}function _n(t,n){t&1&&h(0,"tr",16)}function dn(t,n){if(t&1){let e=B();T(0),s(1,"div",1),h(2,"mat-progress-bar",2),E(3,"async"),r(),C(4,en,4,8,"p"),s(5,"table",3),T(6,4),C(7,nn,2,1,"th",5)(8,on,2,1,"td",6),b(),T(9,7),C(10,rn,2,1,"th",5)(11,an,2,1,"td",6),b(),T(12,8),C(13,ln,2,1,"th",5)(14,sn,2,1,"td",6),b(),T(15,9),C(16,cn,2,1,"th",5)(17,pn,3,3,"td",6),b(),C(18,mn,1,0,"tr",10)(19,_n,1,0,"tr",11),r(),s(20,"app-paginator",12),A("paging",function(a){M(e);let p=c();return $(p.controller.handlePageEvent(a))}),r(),b()}if(t&2){let e=c();o(2),g("mode",w(3,13,e.dataSource.loading$)?"indeterminate":"determinate")("value",0),o(2),d(e.torrent.filesStatus==="over_threshold"?4:-1),o(),g("dataSource",e.dataSource)("multiTemplateDataRows",!0),o(13),g("matHeaderRowDef",e.displayedColumns),o(),g("matRowDefColumns",e.displayedColumns),o(),g("page",e.controls.page)("pageSize",e.controls.limit)("pageLength",e.dataSource.result.items.length)("totalLength",e.dataSource.result.totalCount)("totalIsEstimate",!1)("showLastPage",!0)}}var Me=(()=>{class t{constructor(){this.apollo=y(Dt),this.errorsService=y(Z),this.transloco=y(R),this.displayedColumns=["index","path","type","size"]}ngOnInit(){this.controller=new it(this.torrent.infoHash),this.dataSource=this.torrent.filesStatus==="single"?new nt(this.torrent):new et(this.apollo,this.errorsService,this.controller.variables$),this.controller.controls$.subscribe(e=>{this.controls=e})}item(e){return e}static{this.\u0275fac=function(i){return new(i||t)}}static{this.\u0275cmp=z({type:t,selectors:[["app-torrent-files-table"]],inputs:{torrent:"torrent"},standalone:!0,features:[N],decls:1,vars:0,consts:[[4,"transloco"],[1,"progress-bar-container"],[3,"mode","value"],["mat-table","",1,"table-results",3,"dataSource","multiTemplateDataRows"],["matColumnDef","index"],["mat-header-cell","",4,"matHeaderCellDef"],["mat-cell","",4,"matCellDef"],["matColumnDef","path"],["matColumnDef","type"],["matColumnDef","size"],["mat-header-row","",4,"matHeaderRowDef"],["mat-row","",4,"matRowDef","matRowDefColumns"],[3,"paging","page","pageSize","pageLength","totalLength","totalIsEstimate","showLastPage"],["mat-header-cell",""],["mat-cell",""],["mat-header-row",""],["mat-row",""]],template:function(i,a){i&1&&C(0,dn,21,15,"ng-container",0)},dependencies:[V,ee,ie,re,ce,ae,oe,pe,le,se,me,_e,j,Mt,Q,tt,he]})}}return t})();var fn=(t,n)=>n.key,xn=(t,n)=>n.id,gn=(t,n)=>n.metadataSource.key,hn=t=>({count:t});function Tn(t,n){if(t&1&&h(0,"img",3),t&2){let e=c().$implicit,i=c();g("ngSrc","https://image.tmdb.org/t/p/w300/"+n)("alt",e("torrents.poster"))("width",i.breakpoints.sizeAtLeast("Medium")?300:150)("height",i.breakpoints.sizeAtLeast("Medium")?450:225)}}function bn(t,n){if(t&1&&(s(0,"h2")(1,"a",14),l(2),r()()),t&2){let e=c().$implicit,i=c();o(),g("routerLink","permalink/"+i.torrentContent.infoHash)("matTooltip",e("torrents.permalink")),o(),x(i.torrentContent.torrent.name)}}function vn(t,n){if(t&1&&(s(0,"p",4)(1,"strong"),l(2),r(),l(3),E(4,"filesize"),r()),t&2){let e=c().$implicit,i=c();o(2),m("",e("torrents.size"),":"),o(),m(" ",w(4,2,i.torrentContent.torrent.size)," ")}}function Sn(t,n){if(t&1&&(s(0,"p",5)(1,"strong"),l(2),r(),l(3),E(4,"timeAgo"),r()),t&2){let e=c().$implicit,i=c();o(2),x(e("torrents.published")),o(),m(" ",w(4,2,i.torrentContent.publishedAt)," ")}}function En(t,n){if(t&1&&(s(0,"p",6)(1,"strong"),l(2),r(),l(3),r()),t&2){let e,i=c().$implicit,a=c();o(2),m("",i("torrents.s_l"),":"),o(),wt(" ",(e=a.torrentContent.seeders)!==null&&e!==void 0?e:"?"," / ",(e=a.torrentContent.leechers)!==null&&e!==void 0?e:"?"," ")}}function wn(t,n){if(t&1&&(s(0,"span"),l(1),r()),t&2){let e=n.$implicit,i=n.$index;o(),x((i>0?", ":"")+e.name)}}function In(t,n){if(t&1&&(s(0,"p")(1,"strong"),l(2),r(),l(3),r()),t&2){let e=c().$implicit,i=c();o(2),m("",e("torrents.title"),":"),o(),m(" ",i.torrentContent.content.title," ")}}function yn(t,n){if(t&1&&l(0),t&2){let e=n.$implicit,i=n.$index,a=c(3);m(" ",(i>0?", ":"")+e.name+(e.id===(a.torrentContent.content==null||a.torrentContent.content.originalLanguage==null?null:a.torrentContent.content.originalLanguage.id)?" (original)":"")," ")}}function Mn(t,n){if(t&1&&(s(0,"p")(1,"strong"),l(2),r(),l(3,"\xA0 "),k(4,yn,1,1,null,null,xn),r()),t&2){let e=c().$implicit,i=c();o(2),m("",e("torrents.languages"),":"),o(2),P(i.torrentContent.languages)}}function $n(t,n){if(t&1&&(s(0,"p")(1,"strong"),l(2),r(),l(3),r()),t&2){let e,i=c().$implicit,a=c();o(2),m("",i("torrents.original_release_date"),":"),o(),m(" ",(e=a.torrentContent.content==null?null:a.torrentContent.content.releaseDate)!==null&&e!==void 0?e:a.torrentContent.content==null?null:a.torrentContent.content.releaseYear," ")}}function kn(t,n){if(t&1&&(s(0,"p")(1,"strong"),l(2),r(),l(3),r()),t&2){let e=c().$implicit,i=c();o(2),m("",e("torrents.episodes"),":"),o(),m(" ",i.torrentContent.episodes.label," ")}}function Pn(t,n){if(t&1&&(s(0,"p"),l(1),r()),t&2){let e=c(2);o(),m(" ",e.torrentContent.content.overview," ")}}function An(t,n){if(t&1&&(T(0),s(1,"p")(2,"strong"),l(3),r(),l(4),r(),b()),t&2){let e=c().$implicit;o(3),m("",e("torrents.genres"),":"),o(),m(" ",n.join(", ")," ")}}function Fn(t,n){if(t&1&&(T(0),l(1),E(2,"number"),b()),t&2){let e=c(2).$implicit,i=c();o(),m("(",e("torrents.votes_count_n",It(3,hn,w(2,1,i.torrentContent.content==null?null:i.torrentContent.content.voteCount))),")")}}function Bn(t,n){if(t&1&&(s(0,"p")(1,"strong"),l(2),r(),l(3),C(4,Fn,3,5,"ng-container"),r()),t&2){let e=c().$implicit,i=c();o(2),m("",e("torrents.rating"),":"),o(),m(" ",i.torrentContent.content==null?null:i.torrentContent.content.voteAverage," / 10 "),o(),d((i.torrentContent.content==null?null:i.torrentContent.content.voteCount)!=null?4:-1)}}function Dn(t,n){if(t&1&&(l(0),s(1,"a",15),l(2),r()),t&2){let e=n.$implicit,i=n.$index;m(" ",i>0?", ":"",""),o(),g("href",e.url,Et),o(),x(e.metadataSource.name)}}function On(t,n){if(t&1&&(s(0,"p")(1,"strong"),l(2),r(),l(3,"\xA0 "),k(4,Dn,3,3,"a",15,gn),r()),t&2){let e=c().$implicit;o(2),m("",e("torrents.external_links"),":"),o(2),P(n)}}function zn(t,n){if(t&1&&(s(0,"span",16),l(1),r()),t&2){let e=c(2).$implicit;o(),x(e("torrents.files"))}}function Nn(t,n){t&1&&(s(0,"span",17),l(1),E(2,"number"),r()),t&2&&(o(),m("(",w(2,1,n),")"))}function Rn(t,n){if(t&1&&(s(0,"mat-icon"),l(1,"file_present"),r(),C(2,zn,2,1,"span",16)(3,Nn,3,3,"span",17)),t&2){let e,i=c(2);o(2),d(i.breakpoints.sizeAtLeast("Medium")?2:-1),o(),d((e=i.filesCount())?3:-1,e)}}function jn(t,n){if(t&1&&(s(0,"p"),l(1),r()),t&2){let e=c(2).$implicit;o(),x(e("torrents.files_no_info"))}}function Vn(t,n){if(t&1&&(s(0,"mat-card",18),C(1,jn,2,1,"p"),h(2,"app-torrent-files-table",19),r()),t&2){let e=c(2);o(),d(e.torrentContent.torrent.filesStatus==="no_info"?1:-1),o(),g("torrent",e.torrentContent.torrent)}}function Ln(t,n){if(t&1&&(s(0,"span",16),l(1),r()),t&2){let e=c(2).$implicit;o(),x(e("torrents.edit_tags"))}}function Hn(t,n){if(t&1&&(s(0,"mat-icon"),l(1,"sell"),r(),C(2,Ln,2,1,"span",16)),t&2){let e=c(2);o(2),d(e.breakpoints.sizeAtLeast("Medium")?2:-1)}}function Gn(t,n){if(t&1){let e=B();s(0,"mat-chip-row",25),A("edited",function(a){let p=M(e).$implicit,_=c(3);return $(_.renameTag(p,a.value))})("removed",function(){let a=M(e).$implicit,p=c(3);return $(p.deleteTag(a))}),l(1),s(2,"mat-icon",26),l(3,"cancel"),r()()}if(t&2){let e=n.$implicit;g("editable",!0),o(),m(" ",e," ")}}function Yn(t,n){if(t&1&&(s(0,"mat-option",24),l(1),r()),t&2){let e=n.$implicit;g("value",e),o(),x(e)}}function Kn(t,n){if(t&1){let e=B();s(0,"mat-card")(1,"mat-form-field",20)(2,"mat-chip-grid",null,0),k(4,Gn,4,2,"mat-chip-row",21,K),r(),s(6,"input",22),A("matChipInputTokenEnd",function(a){M(e);let p=c(2);return $(a.value&&p.addTag(a.value))}),r(),s(7,"mat-autocomplete",23,1),A("optionSelected",function(a){M(e);let p=c(2);return $(p.addTag(a.option.viewValue))}),k(9,Yn,2,2,"mat-option",24,K),r()()()}if(t&2){let e=st(3),i=st(8),a=c().$implicit,p=c();o(4),P(p.torrentContent.torrent.tagNames),o(2),g("placeholder",a("torrents.new_tag"))("formControl",p.newTagCtrl)("matAutocomplete",i)("matChipInputFor",e)("matChipInputSeparatorKeyCodes",p.separatorKeysCodes)("value",p.newTagCtrl.value),o(3),P(p.suggestedTags)}}function qn(t,n){if(t&1&&(s(0,"span",16),l(1),r()),t&2){let e=c(2).$implicit;o(),x(e("torrents.delete"))}}function Jn(t,n){if(t&1&&(s(0,"mat-icon"),l(1,"delete_forever"),r(),C(2,qn,2,1,"span",16)),t&2){let e=c(2);o(2),d(e.breakpoints.sizeAtLeast("Medium")?2:-1)}}function Un(t,n){if(t&1){let e=B();s(0,"mat-card")(1,"mat-card-content",27)(2,"p")(3,"strong"),l(4),r(),h(5,"br"),l(6),r()(),s(7,"mat-card-actions",28)(8,"button",29),A("click",function(){M(e);let a=c(2);return $(a.delete())}),s(9,"mat-icon"),l(10,"delete_forever"),r(),l(11),r()()()}if(t&2){let e=c().$implicit;o(4),x(e("torrents.delete_are_you_sure")),o(2),m("",e("torrents.delete_action_cannot_be_undone")," "),o(5),m("",e("torrents.delete")," ")}}function Xn(t,n){t&1&&(s(0,"mat-icon",30),l(1,"close"),r())}function Qn(t,n){t&1&&(s(0,"mat-tab"),C(1,Xn,2,0,"ng-template",12),r())}function Zn(t,n){if(t&1){let e=B();T(0),C(1,Tn,1,4,"img",3)(2,bn,3,3,"h2")(3,vn,5,4,"p",4)(4,Sn,5,4,"p",5)(5,En,4,3,"p",6),s(6,"p",7)(7,"strong"),l(8),r(),s(9,"span",8),l(10),r()(),s(11,"p")(12,"strong"),l(13),r(),l(14,"\xA0 "),k(15,wn,2,1,"span",null,fn),r(),C(17,In,4,2,"p")(18,Mn,6,1,"p")(19,$n,4,2,"p")(20,kn,4,2,"p")(21,Pn,2,1,"p")(22,An,5,2,"ng-container")(23,Bn,5,3,"p")(24,On,6,1,"p"),h(25,"mat-divider",9),s(26,"mat-tab-group",10),A("focusChange",function(a){M(e);let p=c();return $(p.selectTab(a.index==4?0:a.index))}),h(27,"mat-tab",11),s(28,"mat-tab"),C(29,Rn,4,2,"ng-template",12)(30,Vn,3,2,"ng-template",13),r(),s(31,"mat-tab"),C(32,Hn,3,1,"ng-template",12)(33,Kn,11,6,"ng-template",13),r(),s(34,"mat-tab"),C(35,Jn,3,1,"ng-template",12)(36,Un,12,3,"ng-template",13),r(),C(37,Qn,2,0,"mat-tab"),r(),b()}if(t&2){let e,i,a,p=n.$implicit,_=c();o(),d((e=_.getAttribute("poster_path","tmdb"))?1:-1,e),o(),d(_.heading?2:-1),o(),d(_.size?3:-1),o(),d(_.published?4:-1),o(),d(_.peers?5:-1),o(3),m("",p("torrents.info_hash"),":"),o(),g("matTooltip",p("torrents.copy_to_clipboard"))("cdkCopyToClipboard",_.torrentContent.infoHash),o(),x(_.torrentContent.infoHash),o(3),m("",p("torrents.source"),":"),o(2),P(_.torrentContent.torrent.sources),o(2),d(_.torrentContent.content?17:-1),o(),d(_.torrentContent.languages!=null&&_.torrentContent.languages.length?18:-1),o(),d(_.torrentContent.content!=null&&_.torrentContent.content.releaseYear?19:-1),o(),d(_.torrentContent.episodes?20:-1),o(),d(_.torrentContent.content!=null&&_.torrentContent.content.overview?21:-1),o(),d((i=_.getCollections("genre"))?22:-1,i),o(),d((_.torrentContent.content==null?null:_.torrentContent.content.voteAverage)!=null?23:-1),o(),d((a=_.torrentContent.content==null?null:_.torrentContent.content.externalLinks)?24:-1,a),o(2),g("selectedIndex",_.selectedTabIndex)("mat-stretch-tabs",!1),o(11),d(_.selectedTabIndex>0?37:-1)}}var xo=(()=>{class t{constructor(e,i){this.graphQLService=e,this.errorsService=i,this.breakpoints=y(Nt),this.heading=!0,this.size=!0,this.peers=!0,this.published=!0,this.updated=new St,this.newTagCtrl=new Ft(""),this.editedTags=Array(),this.suggestedTags=Array(),this.selectedTabIndex=0,this.separatorKeysCodes=[13,188],this.transloco=y(R),this.newTagCtrl.valueChanges.subscribe(a=>(a&&(a=we(a),this.newTagCtrl.setValue(a,{emitEvent:!1})),e.torrentSuggestTags({input:{prefix:a,exclusions:this.torrentContent.torrent.tagNames}}).pipe(X(p=>{this.suggestedTags.splice(0,this.suggestedTags.length,...p.suggestions.map(_=>_.name))})).subscribe()))}selectTab(e){this.selectedTabIndex=e}addTag(e){this.editTags(i=>[...i,e]),this.saveTags()}renameTag(e,i){this.editTags(a=>a.map(p=>p===e?i:p)),this.saveTags()}deleteTag(e){this.editTags(i=>i.filter(a=>a!==e)),this.saveTags()}editTags(e){this.editedTags=e(this.editedTags),this.newTagCtrl.reset()}saveTags(){this.graphQLService.torrentSetTags({infoHashes:[this.torrentContent.infoHash],tagNames:this.editedTags}).pipe(Y(e=>(this.errorsService.addError(`Error saving tags: ${e.message}`),G))).pipe(X(()=>{this.editedTags=[],this.updated.emit(null)})).subscribe()}delete(){this.graphQLService.torrentDelete({infoHashes:[this.torrentContent.infoHash]}).pipe(Y(e=>(this.errorsService.addError(`Error deleting torrent: ${e.message}`),G))).pipe(X(()=>{this.updated.emit(null)})).subscribe()}getAttribute(e,i){return this.torrentContent.content?.attributes?.find(a=>a.key===e&&(i===void 0||a.source===i))?.value}getCollections(e){let i=this.torrentContent.content?.collections?.filter(a=>a.type===e).map(a=>a.name);return i?.length?i.sort():void 0}filesCount(){return this.torrentContent.torrent.filesStatus==="single"?1:this.torrentContent.torrent.filesCount??void 0}static{this.\u0275fac=function(i){return new(i||t)(lt(ge),lt(Z))}}static{this.\u0275cmp=z({type:t,selectors:[["app-torrent-content"]],inputs:{torrentContent:"torrentContent",heading:"heading",size:"size",peers:"peers",published:"published"},outputs:{updated:"updated"},standalone:!0,features:[N],decls:1,vars:0,consts:[["chipGrid",""],["auto","matAutocomplete"],[4,"transloco"],[1,"poster",3,"ngSrc","alt","width","height"],[1,"size"],[1,"published"],[1,"peers"],[1,"info-hash"],[3,"matTooltip","cdkCopyToClipboard"],[2,"clear","both"],["animationDuration","0",3,"focusChange","selectedIndex","mat-stretch-tabs"],["aria-labelledby","hidden"],["mat-tab-label",""],["matTabContent",""],[3,"routerLink","matTooltip"],["target","_blank",3,"href"],[1,"label"],[1,"files-count"],[1,"torrent-files"],[3,"torrent"],["subscriptSizing","dynamic",1,"form-edit-tags"],[3,"editable"],["autocapitalize","none",3,"matChipInputTokenEnd","placeholder","formControl","matAutocomplete","matChipInputFor","matChipInputSeparatorKeyCodes","value"],[3,"optionSelected"],[3,"value"],[3,"edited","removed","editable"],["matChipRemove",""],[2,"margin-top","10px"],[1,"button-row"],["mat-stroked-button","","color","warn",3,"click"],[2,"margin-right","0"]],template:function(i,a){i&1&&C(0,Zn,38,21,"ng-container",2)},dependencies:[V,jt,Lt,Ot,Ht,zt,Gt,Kt,Yt,Zt,Wt,Jt,Xt,te,Vt,W,de,ue,Ce,fe,ne,Pt,At,Bt,kt,j,Q,tt,$t,Rt,Me],styles:["h2[_ngcontent-%COMP%]{margin-top:10px;max-width:900px;white-space:pre-wrap;word-break:break-all;overflow-wrap:break-word}.poster[_ngcontent-%COMP%]{float:right;margin:10px;border:1px solid currentColor}.info-hash[_ngcontent-%COMP%]{white-space:pre-wrap;word-break:break-all;overflow-wrap:break-word}.info-hash[_ngcontent-%COMP%]   span[_ngcontent-%COMP%]{padding-left:5px;cursor:crosshair;text-decoration:underline;text-decoration-style:dotted}.torrent-files[_ngcontent-%COMP%]{padding-top:10px;max-height:800px;overflow:scroll}.torrent-files[_ngcontent-%COMP%]   table[_ngcontent-%COMP%]{margin-bottom:10px;width:800px}.torrent-files[_ngcontent-%COMP%]   td[_ngcontent-%COMP%]{padding-right:20px;border-bottom:1px solid rgba(0,0,0,.12)}.torrent-files[_ngcontent-%COMP%]   tr[_ngcontent-%COMP%]:hover   td[_ngcontent-%COMP%]{background-color:#f5f5f5}.form-edit-tags[_ngcontent-%COMP%]     .mat-mdc-form-field-subscript-wrapper{display:none}.files-count[_ngcontent-%COMP%]{margin-left:4px}.mat-mdc-card-content[_ngcontent-%COMP%]   p[_ngcontent-%COMP%]{margin-top:0}  .mdc-tab[aria-labelledby=hidden]{display:none}  .mdc-tab[role=tab]{padding-left:15px;padding-right:15px}  .mdc-tab .label,   .mdc-tab .files-count{margin-left:8px}"]})}}return t})();var $e={movie:{singular:"Movie",plural:"Movies",icon:"movie"},tv_show:{singular:"TV Show",plural:"TV Shows",icon:"live_tv"},music:{singular:"Music",plural:"Music",icon:"music_note"},ebook:{singular:"E-Book",plural:"E-Books",icon:"auto_stories"},comic:{singular:"Comic",plural:"Comics",icon:"comic_bubble"},audiobook:{singular:"Audiobook",plural:"Audiobooks",icon:"mic"},software:{singular:"Software",plural:"Software",icon:"desktop_windows"},game:{singular:"Game",plural:"Games",icon:"sports_esports"},xxx:{singular:"XXX",plural:"XXX",icon:"18_up_rating"},null:{singular:"Unknown",plural:"Unknown",icon:"question_mark"}},ho=Object.entries($e).map(([t,n])=>U({key:t},n)),To=t=>t?$e[t]:void 0;var Wn=(t,n)=>n.id;function ti(t,n){if(t&1&&(s(0,"mat-chip",1)(1,"mat-icon",2),l(2,"sell"),r(),l(3),r()),t&2){let e=n.$implicit;o(3),m(" ",e," ")}}function ei(t,n){if(t&1&&l(0),t&2){let e=n.$implicit,i=n.$index,a=c(),p=c().$implicit;m(" ",p("languages."+e.id)+(i<a.length-1?", ":"")," ")}}function ni(t,n){t&1&&(s(0,"mat-chip")(1,"mat-icon",2),l(2,"translate"),r(),k(3,ei,1,1,null,null,Wn),r()),t&2&&(o(3),P(n))}function ii(t,n){t&1&&(s(0,"mat-chip"),l(1),r()),t&2&&(o(),x(n))}function oi(t,n){t&1&&(s(0,"mat-chip")(1,"mat-icon",2),l(2,"aspect_ratio"),r(),l(3),r()),t&2&&(o(3),x(n))}function ri(t,n){t&1&&(s(0,"mat-chip")(1,"mat-icon",2),l(2,"album"),r(),l(3),r()),t&2&&(o(3),x(n))}function ai(t,n){t&1&&(s(0,"mat-chip"),h(1,"mat-icon",3),l(2),r()),t&2&&(o(2),x(n))}function li(t,n){t&1&&(s(0,"mat-chip"),l(1),r()),t&2&&(o(),x(n))}function si(t,n){if(t&1&&(T(0),s(1,"mat-chip-set"),k(2,ti,4,1,"mat-chip",1,K),C(4,ni,5,0,"mat-chip")(5,ii,2,1,"mat-chip")(6,oi,4,1,"mat-chip")(7,ri,4,1,"mat-chip")(8,ai,3,1,"mat-chip")(9,li,2,1,"mat-chip"),r(),b()),t&2){let e,i,a,p,_,D,v=c();o(2),P(v.torrentContent.torrent.tagNames),o(2),d((e=v.torrentContent.languages)?4:-1,e),o(),d((i=v.torrentContent.video3d==null?null:v.torrentContent.video3d.slice(1))?5:-1,i),o(),d((a=v.torrentContent.videoResolution==null?null:v.torrentContent.videoResolution.slice(1))?6:-1,a),o(),d((p=v.torrentContent.videoSource)?7:-1,p),o(),d((_=v.torrentContent.videoCodec)?8:-1,_),o(),d((D=v.torrentContent.videoModifier)?9:-1,D)}}var Eo=(()=>{class t{static{this.\u0275fac=function(i){return new(i||t)}}static{this.\u0275cmp=z({type:t,selectors:[["app-torrent-chips"]],inputs:{torrentContent:"torrentContent"},standalone:!0,features:[N],decls:1,vars:0,consts:[[4,"transloco"],[1,"chip-primary"],["matChipAvatar",""],["matChipAvatar","","svgIcon","binary"]],template:function(i,a){i&1&&C(0,si,10,6,"ng-container",0)},dependencies:[V,Ut,qt,Qt,W,j],styles:["mat-chip-set[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%]{position:relative;left:4px}"]})}}return t})();export{tt as a,xo as b,$e as c,ho as d,To as e,Eo as f};