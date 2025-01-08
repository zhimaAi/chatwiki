import{a2 as j,a6 as f,$ as w,a0 as b,aF as o,aS as _,aR as p,aE as G,T as F,ca as O,a1 as E,a3 as I,al as P}from"../../index-P5YTNo6V.js";import{e as T,b as A,j as N}from"./vue-chunks-DrvJJrR0.js";const d=["wrap","nowrap","wrap-reverse"],m=["flex-start","flex-end","start","end","center","space-between","space-around","space-evenly","stretch","normal","left","right"],y=["center","start","end","flex-start","flex-end","self-start","self-end","baseline","normal","stretch"],V=(e,t)=>{const a={};return d.forEach(n=>{a[`${e}-wrap-${n}`]=t.wrap===n}),a},W=(e,t)=>{const a={};return y.forEach(n=>{a[`${e}-align-${n}`]=t.align===n}),a[`${e}-align-stretch`]=!t.align&&!!t.vertical,a},L=(e,t)=>{const a={};return m.forEach(n=>{a[`${e}-justify-${n}`]=t.justify===n}),a};function D(e,t){return j(f(f(f({},V(e,t)),W(e,t)),L(e,t)))}const J=e=>{const{componentCls:t}=e;return{[t]:{display:"flex","&-vertical":{flexDirection:"column"},"&-rtl":{direction:"rtl"},"&:empty":{display:"none"}}}},M=e=>{const{componentCls:t}=e;return{[t]:{"&-gap-small":{gap:e.flexGapSM},"&-gap-middle":{gap:e.flexGap},"&-gap-large":{gap:e.flexGapLG}}}},R=e=>{const{componentCls:t}=e,a={};return d.forEach(n=>{a[`${t}-wrap-${n}`]={flexWrap:n}}),a},z=e=>{const{componentCls:t}=e,a={};return y.forEach(n=>{a[`${t}-align-${n}`]={alignItems:n}}),a},H=e=>{const{componentCls:t}=e,a={};return m.forEach(n=>{a[`${t}-justify-${n}`]={justifyContent:n}}),a},X=w("Flex",e=>{const t=b(e,{flexGapSM:e.paddingXS,flexGap:e.padding,flexGapLG:e.paddingLG});return[J(t),M(t),R(t),z(t),H(t)]});function g(e){return["small","middle","large"].includes(e)}const q=()=>({prefixCls:o(),vertical:_(),wrap:o(),justify:o(),align:o(),flex:p([Number,String]),gap:p([Number,String]),component:G()});var B=function(e,t){var a={};for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&t.indexOf(n)<0&&(a[n]=e[n]);if(e!=null&&typeof Object.getOwnPropertySymbols=="function")for(var l=0,n=Object.getOwnPropertySymbols(e);l<n.length;l++)t.indexOf(n[l])<0&&Object.prototype.propertyIsEnumerable.call(e,n[l])&&(a[n[l]]=e[n[l]]);return a};const Q=T({name:"AFlex",inheritAttrs:!1,props:q(),setup(e,t){let{slots:a,attrs:n}=t;const{flex:l,direction:x}=O(),{prefixCls:s}=E("flex",e),[v,S]=X(s),C=A(()=>{var r;return[s.value,S.value,D(s.value,e),{[`${s.value}-rtl`]:x.value==="rtl",[`${s.value}-gap-${e.gap}`]:g(e.gap),[`${s.value}-vertical`]:(r=e.vertical)!==null&&r!==void 0?r:l==null?void 0:l.value.vertical}]});return()=>{var r;const{flex:u,gap:c,component:$="div"}=e,h=B(e,["flex","gap","component"]),i={};return u&&(i.flex=u),c&&!g(c)&&(i.gap=`${c}px`),v(N($,I({class:[n.class,C.value],style:[n.style,i]},P(h,["justify","wrap","align","vertical"])),{default:()=>[(r=a.default)===null||r===void 0?void 0:r.call(a)]}))}}}),Z=F(Q);export{Z as _};