import{a as U}from"./chunk-RQ2LQKI2.js";import{Ea as K,Fa as Q,R as q,S as G,_a as R,a as b,b as D,s as j,w as H,wa as J}from"./chunk-6SW7H33Q.js";import{$b as o,Cb as C,Ea as g,Fa as u,Hb as m,Lc as v,Mb as S,Na as L,Nb as M,Ob as z,Pb as E,Qb as a,Rb as s,Tb as V,Ub as k,Wb as y,Zb as _,bd as $,kb as l,kc as c,mc as P,pc as F,qa as f,rc as N,sc as O,tc as A,ua as T,vc as h,wc as d,xa as x,xc as B,zb as w}from"./chunk-Z3WUIYN5.js";var X=(()=>{class t{constructor(){this.transloco=f(b)}transform(e,i=!0,n=2){if(i&&e>0&&n>0){let Y=Math.floor(Math.log10(Math.abs(e))),I=Math.pow(10,Y-(n-1));e=Math.round(e/I)*I}let p=Intl.NumberFormat(this.transloco.getActiveLang()).format(e);return i?`~${p}`:p}static{this.\u0275fac=function(i){return new(i||t)}}static{this.\u0275pipe=x({name:"intEstimate",type:t,pure:!1,standalone:!0})}}return t})();var Z=(t,r,e)=>({x:t,y:r,z:e}),tt=(t,r)=>({x:t,y:r}),et=t=>[null,t];function it(t,r){if(t&1&&(a(0,"mat-option",4),c(1),s()),t&2){let e=r.$implicit;m("value",e),l(),P(" ",e," ")}}function nt(t,r){if(t&1&&(c(0),h(1,"number"),h(2,"number"),h(3,"intEstimate")),t&2){let e,i=o().$implicit,n=o();P(" ",i("paginator.x_to_y_of_z",A(8,Z,d(1,1,n.firstItemIndex),d(2,3,n.lastItemIndex),B(3,5,(e=n.totalLength)!==null&&e!==void 0?e:0,n.totalIsEstimate)))," ")}}function at(t,r){if(t&1&&(c(0),h(1,"number"),h(2,"number")),t&2){let e=o().$implicit,i=o();P(" ",e("paginator.x_to_y",O(5,tt,d(1,1,i.firstItemIndex),d(2,3,i.lastItemIndex)))," ")}}function ot(t,r){if(t&1){let e=y();a(0,"button",7),_("click",function(){let n;g(e);let p=o(2);return p.page=(n=p.pageCount)!==null&&n!==void 0?n:1,u(p.emitChange())}),a(1,"mat-icon"),c(2,"last_page"),s()()}if(t&2){let e=o().$implicit,i=o();m("disabled",N(2,et,i.page).includes(i.pageCount))("matTooltip",e("paginator.last_page"))}}function rt(t,r){if(t&1){let e=y();V(0),a(1,"div",1)(2,"mat-form-field",2)(3,"mat-label"),c(4,"Items per page"),s(),a(5,"mat-select",3),_("valueChange",function(n){g(e);let p=o();return p.pageSize=n,p.page=1,u(p.emitChange())}),z(6,it,2,2,"mat-option",4,M),s()(),a(8,"p",5),C(9,nt,4,12)(10,at,3,8),s(),a(11,"div",6)(12,"button",7),_("click",function(){g(e);let n=o();return n.page=1,u(n.emitChange())}),a(13,"mat-icon"),c(14,"first_page"),s()(),a(15,"button",7),_("click",function(){g(e);let n=o();return n.page=n.page-1,u(n.emitChange())}),a(16,"mat-icon"),c(17,"navigate_before"),s()(),a(18,"button",7),_("click",function(){g(e);let n=o();return n.page=n.page+1,u(n.emitChange())}),a(19,"mat-icon"),c(20,"navigate_next"),s()(),C(21,ot,3,4,"button",8),s()(),k()}if(t&2){let e=r.$implicit,i=o();l(5),m("value",i.pageSize),l(),E(i.pageSizes),l(3),S(i.hasTotalLength?9:10),l(3),m("disabled",!i.hasPreviousPage)("matTooltip",e("paginator.first_page")),l(3),m("disabled",!i.hasPreviousPage)("matTooltip",e("paginator.previous_page")),l(3),m("disabled",!i.actuallyHasNextPage)("matTooltip",e("paginator.next_page")),l(3),S(i.showLastPage?21:-1)}}var It=(()=>{class t{constructor(){this.page=1,this.pageSize=10,this.pageSizes=[10,20,50,100],this.pageLength=0,this.totalLength=null,this.totalIsEstimate=!1,this.hasNextPage=null,this.showLastPage=!1,this.paging=new L}get firstItemIndex(){return(this.page-1)*this.pageSize+1}get lastItemIndex(){return(this.page-1)*this.pageSize+this.pageLength}get hasTotalLength(){return typeof this.totalLength=="number"}get hasPreviousPage(){return this.page>1}get pageCount(){return typeof this.totalLength!="number"?null:Math.ceil(this.totalLength/this.pageSize)}get actuallyHasNextPage(){return typeof this.hasNextPage=="boolean"?this.hasNextPage:typeof this.totalLength!="number"?!1:this.page*this.pageSize<this.totalLength}emitChange(){this.paging.emit({page:this.page,pageSize:this.pageSize})}static{this.\u0275fac=function(i){return new(i||t)}}static{this.\u0275cmp=T({type:t,selectors:[["app-paginator"]],inputs:{page:[2,"page","page",v],pageSize:[2,"pageSize","pageSize",v],pageSizes:"pageSizes",pageLength:[2,"pageLength","pageLength",v],totalLength:"totalLength",totalIsEstimate:"totalIsEstimate",hasNextPage:"hasNextPage",showLastPage:"showLastPage"},outputs:{paging:"paging"},standalone:!0,features:[w,F],decls:1,vars:0,consts:[[4,"transloco"],[1,"paginator"],["subscriptSizing","dynamic",1,"field-items-per-page"],[3,"valueChange","value"],[3,"value"],[1,"paginator-description"],[1,"paginator-navigation"],["mat-icon-button","",3,"click","disabled","matTooltip"],["mat-icon-button","",3,"disabled","matTooltip"]],template:function(i,n){i&1&&C(0,rt,22,9,"ng-container",0)},dependencies:[R,j,H,G,q,J,Q,K,D,$,X],styles:[".paginator[_ngcontent-%COMP%] > *[_ngcontent-%COMP%]{display:inline-block;vertical-align:middle}.paginator[_ngcontent-%COMP%]   p[_ngcontent-%COMP%]{margin:0 20px}.paginator[_ngcontent-%COMP%]   .field-items-per-page[_ngcontent-%COMP%]{width:140px}"]})}}return t})();var Et=(()=>{class t{constructor(){this.transloco=f(b)}transform(e){return U(e,this.transloco.getActiveLang())}static{this.\u0275fac=function(i){return new(i||t)}}static{this.\u0275pipe=x({name:"timeAgo",type:t,pure:!1,standalone:!0})}}return t})();export{X as a,It as b,Et as c};
