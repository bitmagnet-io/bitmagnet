import{d as M,e as y,f as w,g as f,h as g,j as S,k as D,l as I,m as X}from"./chunk-3PIWRAOL.js";function h(t,e){let r=+f(t)-+f(e);return r<0?-1:r>0?1:r}function T(t){return w(t,Date.now())}function _(t,e,r){let[s,o]=D(r?.in,t,e),c=s.getFullYear()-o.getFullYear(),i=s.getMonth()-o.getMonth();return c*12+i}function F(t){return e=>{let s=(t?Math[t]:Math.trunc)(e);return s===0?0:s}}function Y(t,e){return+f(t)-+f(e)}function z(t,e){let r=f(t,e?.in);return r.setHours(23,59,59,999),r}function L(t,e){let r=f(t,e?.in),s=r.getMonth();return r.setFullYear(r.getFullYear(),s+1,0),r.setHours(23,59,59,999),r}function N(t,e){let r=f(t,e?.in);return+z(r,e)==+L(r,e)}function b(t,e,r){let[s,o,c]=D(r?.in,t,t,e),i=h(o,c),n=Math.abs(_(o,c));if(n<1)return 0;o.getMonth()===1&&o.getDate()>27&&o.setDate(30),o.setMonth(o.getMonth()-i*n);let m=h(o,c)===-i;N(s)&&n===1&&h(s,c)===1&&(m=!1);let u=i*(n-+m);return u===0?0:u}function A(t,e,r){let s=Y(t,e)/1e3;return F(r?.roundingMethod)(s)}function v(t,e,r){let s=g(),o=r?.locale??s.locale??I,c=2520,i=h(t,e);if(isNaN(i))throw new RangeError("Invalid time value");let n=Object.assign({},r,{addSuffix:r?.addSuffix,comparison:i}),[m,u]=D(r?.in,...i>0?[e,t]:[t,e]),l=A(u,m),x=(S(u)-S(m))/1e3,a=Math.round((l-x)/60),p;if(a<2)return r?.includeSeconds?l<5?o.formatDistance("lessThanXSeconds",5,n):l<10?o.formatDistance("lessThanXSeconds",10,n):l<20?o.formatDistance("lessThanXSeconds",20,n):l<40?o.formatDistance("halfAMinute",0,n):l<60?o.formatDistance("lessThanXMinutes",1,n):o.formatDistance("xMinutes",1,n):a===0?o.formatDistance("lessThanXMinutes",1,n):o.formatDistance("xMinutes",a,n);if(a<45)return o.formatDistance("xMinutes",a,n);if(a<90)return o.formatDistance("aboutXHours",1,n);if(a<y){let d=Math.round(a/60);return o.formatDistance("aboutXHours",d,n)}else{if(a<c)return o.formatDistance("xDays",1,n);if(a<M){let d=Math.round(a/y);return o.formatDistance("xDays",d,n)}else if(a<M*2)return p=Math.round(a/M),o.formatDistance("aboutXMonths",p,n)}if(p=b(u,m),p<12){let d=Math.round(a/M);return o.formatDistance("xMonths",d,n)}else{let d=p%12,O=Math.trunc(p/12);return d<3?o.formatDistance("aboutXYears",O,n):d<9?o.formatDistance("overXYears",O,n):o.formatDistance("almostXYears",O+1,n)}}function H(t,e){return v(t,T(t),e)}var C=["years","months","weeks","days","hours","minutes","seconds"];function k(t,e){let r=g(),s=e?.locale??r.locale??I,o=e?.format??C,c=e?.zero??!1,i=e?.delimiter??" ";return s.formatDistance?o.reduce((m,u)=>{let l=`x${u.replace(/(^.)/,a=>a.toUpperCase())}`,x=t[u];return x!==void 0&&(c||t[u])?m.concat(s.formatDistance(l,x)):m},[]).join(i):""}var _t=(t,e)=>H(t,{addSuffix:!0,locale:X(e)}),Ft=(t,e)=>k(t,{locale:X(e)});export{_t as a,Ft as b};
