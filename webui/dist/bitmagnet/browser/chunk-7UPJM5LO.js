import{b as j,c as B}from"./chunk-FZR3LVFA.js";import"./chunk-ENCA42F6.js";import"./chunk-RQ2LQKI2.js";import{c as U}from"./chunk-PX3TOMGC.js";import"./chunk-42PJPEMD.js";import"./chunk-ASLGZ7DJ.js";import"./chunk-MSAOOVCY.js";import{a as R}from"./chunk-DSEDLZDW.js";import{E as I,J as S,P as w,Q as E,R as F,S as A,U as D,Ua as L,V as z,b,i as O,qa as H,va as q,ya as $}from"./chunk-VAEZNV34.js";import"./chunk-6XXA7HXI.js";import{i as y,k as P}from"./chunk-CMNWCZJM.js";import{$b as d,Cb as C,Hb as l,Mb as h,Qb as i,Rb as a,Sb as m,Tb as x,Ub as M,cc as v,hb as _,kb as r,kc as u,lc as f,pc as T,qa as s,sc as k,ua as g}from"./chunk-Z3WUIYN5.js";var G=(e,p)=>[e,p];function Q(e,p){e&1&&m(0,"mat-progress-bar",2)}function J(e,p){if(e&1&&(i(0,"mat-card",3)(1,"mat-card-header")(2,"mat-icon",4),u(3),a(),i(4,"mat-card-title")(5,"h2"),u(6),a(),i(7,"a",5),m(8,"mat-icon",6),a()(),i(9,"mat-card-subtitle"),m(10,"app-torrent-chips",7),a()(),i(11,"mat-card-content"),m(12,"app-torrent-content",8),a()()),e&2){let o,t,c=d().$implicit,n=d();r(2),l("matTooltip",c("content_types.singular."+((o=n.torrentContent.contentType)!==null&&o!==void 0?o:"null"))),r(),f((t=(t=n.contentTypeInfo(n.torrentContent.contentType))==null?null:t.icon)!==null&&t!==void 0?t:"question_mark"),r(3),f(n.torrentContent.torrent.name),r(),v("href",n.torrentContent.torrent.magnetUri,_),r(3),l("torrentContent",n.torrentContent),r(2),l("torrentContent",n.torrentContent)("heading",!1)}}function K(e,p){if(e&1&&(x(0),m(1,"app-document-title",1),C(2,Q,1,0,"mat-progress-bar",2)(3,J,13,7,"mat-card",3),M()),e&2){let o=p.$implicit,t=d();r(),l("parts",k(2,G,t.torrentContent==null?null:t.torrentContent.title,o("torrents.permalink"))),r(),h(t.torrentContent?3:2)}}var st=(()=>{class e{constructor(){this.route=s(y),this.router=s(P),this.apollo=s(O),this.contentTypeInfo=U}ngOnInit(){this.route.paramMap.subscribe(o=>{let t=o.get("infoHash");if(typeof t!="string"||!/^[0-9a-f]{40}$/.test(t))return this.notFound();this.apollo.query({query:I,variables:{input:{infoHashes:[t]}},fetchPolicy:"no-cache"}).subscribe(c=>{let n=c.data.torrentContent.search.items;if(n.length===0)return this.notFound();this.torrentContent=n[0]})})}notFound(){this.router.navigate(["/not-found"],{skipLocationChange:!0})}static{this.\u0275fac=function(t){return new(t||e)}}static{this.\u0275cmp=g({type:e,selectors:[["app-torrent-permalink"]],standalone:!0,features:[T],decls:1,vars:0,consts:[[4,"transloco"],[3,"parts"],["mode","indeterminate"],[1,"torrent-permalink"],["matCardAvatar","",3,"matTooltip"],[1,"magnet-link",3,"href"],["svgIcon","magnet"],[3,"torrentContent"],[3,"torrentContent","heading"]],template:function(t,c){t&1&&C(0,K,4,5,"ng-container",0)},dependencies:[L,w,z,F,D,A,E,H,q,$,b,S,B,j,R],styles:[".torrent-permalink[_ngcontent-%COMP%]{max-width:900px;margin:20px auto}.torrent-permalink[_ngcontent-%COMP%]   mat-card-title[_ngcontent-%COMP%]   h2[_ngcontent-%COMP%]{margin:0;font-size:24px;word-break:break-word;overflow-wrap:break-word;padding-right:80px}.torrent-permalink[_ngcontent-%COMP%]   mat-card-title[_ngcontent-%COMP%]   .magnet-link[_ngcontent-%COMP%]{position:absolute;right:20px;top:20px}.torrent-permalink[_ngcontent-%COMP%]   .mat-mdc-card-avatar[_ngcontent-%COMP%]{font-size:44px;margin-top:-10px;border-radius:0;overflow:visible}.torrent-permalink[_ngcontent-%COMP%]   mat-card-subtitle[_ngcontent-%COMP%]{margin:16px 0 14px -56px;font-size:6px}"]})}}return e})();export{st as TorrentPermalinkComponent};
