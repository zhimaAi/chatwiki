import{aH as l,W as a,a9 as o,aD as c,aj as f,aS as s}from"../../index-QYI4dmfl.js";const u=e=>e!=null&&(Array.isArray(e)?l(e).length:!0);function S(e){return u(e.prefix)||u(e.suffix)||u(e.allowClear)}function x(e){return u(e.addonBefore)||u(e.addonAfter)}function B(e){return typeof e>"u"||e===null?"":String(e)}function C(e,i,t,r){if(!t)return;const n=i;if(i.type==="click"){Object.defineProperty(n,"target",{writable:!0}),Object.defineProperty(n,"currentTarget",{writable:!0});const d=e.cloneNode(!0);n.target=d,n.currentTarget=d,d.value="",t(n);return}if(r!==void 0){Object.defineProperty(n,"target",{writable:!0}),Object.defineProperty(n,"currentTarget",{writable:!0}),n.target=e,n.currentTarget=e,e.value=r,t(n);return}t(n)}function F(e,i){if(!e)return;e.focus(i);const{cursor:t}=i||{};if(t){const r=e.value.length;switch(t){case"start":e.setSelectionRange(0,0);break;case"end":e.setSelectionRange(r,r);break;default:e.setSelectionRange(0,r)}}}const p=()=>({addonBefore:o.any,addonAfter:o.any,prefix:o.any,suffix:o.any,clearIcon:o.any,affixWrapperClassName:String,groupClassName:String,wrapperClassName:String,inputClassName:String,allowClear:{type:Boolean,default:void 0}}),y=()=>a(a({},p()),{value:{type:[String,Number,Symbol],default:void 0},defaultValue:{type:[String,Number,Symbol],default:void 0},inputElement:o.any,prefixCls:String,disabled:{type:Boolean,default:void 0},focused:{type:Boolean,default:void 0},triggerFocus:Function,readonly:{type:Boolean,default:void 0},handleReset:Function,hidden:{type:Boolean,default:void 0}}),g=()=>a(a({},y()),{id:String,placeholder:{type:[String,Number]},autocomplete:String,type:c("text"),name:String,size:{type:String},autofocus:{type:Boolean,default:void 0},lazy:{type:Boolean,default:!0},maxlength:Number,loading:{type:Boolean,default:void 0},bordered:{type:Boolean,default:void 0},showCount:{type:[Boolean,Object]},htmlSize:Number,onPressEnter:Function,onKeydown:Function,onKeyup:Function,onFocus:Function,onBlur:Function,onChange:Function,onInput:Function,"onUpdate:value":Function,onCompositionstart:Function,onCompositionend:Function,valueModifiers:Object,hidden:{type:Boolean,default:void 0},status:String}),b=()=>f(g(),["wrapperClassName","groupClassName","inputClassName","affixWrapperClassName"]),N=()=>a(a({},f(b(),["prefix","addonBefore","addonAfter","suffix"])),{rows:Number,autosize:{type:[Boolean,Object],default:void 0},autoSize:{type:[Boolean,Object],default:void 0},onResize:{type:Function},onCompositionstart:s(),onCompositionend:s(),valueModifiers:Object});export{x as a,y as b,b as c,N as d,B as f,S as h,g as i,C as r,F as t};