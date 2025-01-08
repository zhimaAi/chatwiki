import{a0 as E,b as h,w as z,a5 as R,L as o,M as c,V as a,F as C,a1 as N,a2 as q,j as t,Q as u,X as P,Y as g,P as x,r as I,v as H,B as K,u as w,a6 as Q,a7 as W,a8 as T,A as F,U,a9 as X,W as j}from"./vue-chunks-DrvJJrR0.js";import{L as Y,a as G}from"./locale-dropdown-BaQkvZHD.js";import{_ as k,d as J,a as $,b as Z,R as ee,c as te,r as oe,e as se,M as ae}from"../../index-P5YTNo6V.js";import{_ as re,B as ne}from"./index-BbkuOfbn.js";import{_ as le,M as ie}from"./index-CH9v5L6a.js";import"./index-DpfcpfQ0.js";import{D as ce}from"./dropdown-7DbQfkaS.js";import{v as ue}from"./validate-8MtiUsxW.js";import{F as D,_ as _e}from"./index-BUTxM_ZP.js";import{_ as de}from"./index-Dh-ZljIv.js";import{_ as pe}from"./Password-DoqQ5Cp-.js";import"./index-DbkTLXZJ.js";import"./dayjs-Bv4VEw_C.js";import"./axios-Cm0UX6qg.js";import"./qs-QntUJHOZ.js";import"./crypto-js-BpV8n4kc.js";import"./DownOutlined-BHsv0Dlx.js";import"./shallowequal-CfnxU2uU.js";import"./colors-CF1NjnDv.js";import"./ResizeObserver.es-B1PUzC5B.js";import"./collapseMotion-DfZt1qrJ.js";import"./slide-BQHNwZxO.js";import"./index-CCuYH0LT.js";import"./index-GBuEB8v_.js";import"./Dropdown-ZPlzQgPg.js";import"./RightOutlined-DIZF9Vab.js";import"./move-DKHu4ARt.js";import"./Input-Dgd5v23J.js";import"./FormItemContext-C-6BWMej.js";import"./index-DuYANmfA.js";import"./inputProps-CVwJ66mj.js";import"./responsiveObserve-BH6_flfv.js";import"./QuestionCircleOutlined-Dtch7f6L.js";import"./Search-DbjjEMkL.js";import"./SearchOutlined-BsZdEpth.js";import"./TextArea-2ujELiW_.js";const me={class:"navbar-wrapper"},ve={class:"navbar"},fe={__name:"layout-navbar",setup(b){const e=E(),i=h(()=>e.meta.activeMenu||""),_=h(()=>e.path.split("/")[1]),d=J(),m=h(()=>{const r=["ModelManage","TokenManage","TeamManage","AccountManage","CompanyManage","ClientSideManage"];let n=!1,{role_permission:v}=d,s=[];for(let l=0;l<v.length;l++){const f=v[l];if(f==="RobotManage"&&s.push({id:1,key:"robot",label:"robot",title:"机器人管理",path:"/robot/list"}),f==="LibraryManage"&&s.push({id:2,key:"library",label:"library",title:"知识库管理",path:"/library/list"}),f==="FormManage"&&s.push({id:3,key:"database",label:"database",title:"数据库",path:"/database/list"}),!n){for(let y=0;y<r.length;y++)if(r[y]===f){s.push({id:4,key:"user",label:"user",title:"系统管理",path:"/user/model"}),n=!0;break}}}return s.sort((l,f)=>l.id-f.id)});return z(()=>e.path,()=>{},{immediate:!0}),(r,n)=>{const v=R("router-link");return o(),c("div",me,[a("div",ve,[(o(!0),c(C,null,N(m.value,s=>(o(),c("div",{class:q(["nav-menu",{active:s.key===_.value||s.key===i.value}]),key:s.key},[t(v,{to:s.path,class:"nav-menu-name"},{default:u(()=>[P(g(s.title),1)]),_:2},1032,["to"])],2))),128))])])}}},he=k(fe,[["__scopeId","data-v-b3d5716e"]]),ye={class:"layout-breadcrumb"},be={key:1,class:"page-title"},ge={__name:"layout-breadcrumb",props:{title:{type:[String,Boolean],default:""},items:{type:[Array,Boolean],default:()=>[]}},setup(b){const e=b;return(i,_)=>{const d=R("router-link"),m=re,r=ne;return o(),c("div",ye,[e.items.length>0?(o(),x(r,{key:0},{default:u(()=>[(o(!0),c(C,null,N(e.items,(n,v)=>(o(),c(C,{key:n.title},[v!==e.items.length-1?(o(),x(m,{key:0},{default:u(()=>[t(d,{to:n.path},{default:u(()=>[P(g(n.title),1)]),_:2},1032,["to"])]),_:2},1024)):(o(),x(m,{key:1},{default:u(()=>[P(g(n.title),1)]),_:2},1024))],64))),128))]),_:1})):(o(),c("div",be,g(e.title),1))])}}},we=k(ge,[["__scopeId","data-v-c69a461d"]]),ke={class:"layout-footer"},Me={class:"copyright-text"},xe={__name:"layout-footer",setup(b){const e=$(),{t:i}=Z(),_=I(null);return H(()=>{_.value=setInterval(()=>{setTimeout(()=>{e.refreshToken()},0)},ee)}),K(()=>{clearInterval(_.value),_.value=null}),(d,m)=>(o(),c("div",ke,[a("div",Me,g(w(i)("common.copyright")),1)]))}},Se=k(xe,[["__scopeId","data-v-ab337acb"]]),$e={key:0,class:"user-dropdown"},Ce=["src"],Pe={class:"user-name"},Ie={__name:"user-dropdown",setup(b){const e=$(),{userInfo:i,avatar:_,user_name:d}=Q(e),m=()=>{e.logoutConfirm(!0)};return(r,n)=>{const v=te,s=le,l=ie,f=ce;return w(i)?(o(),c("div",$e,[t(f,null,{overlay:u(()=>[t(l,null,{default:u(()=>[t(s,null,{default:u(()=>[a("a",{href:"javascript:;",onClick:m},"退出登录")]),_:1})]),_:1})]),default:u(()=>[a("div",{class:"user-dropdown-link",onClick:n[0]||(n[0]=W(()=>{},["prevent"]))},[a("img",{class:"user-avatar",src:w(_),alt:""},null,8,Ce),a("span",Pe,g(w(d)),1),t(v,{name:"arrow-down",style:{"font-size":"16px",color:"#8c8c8c"}})])]),_:1})])):T("",!0)}}},Te=k(Ie,[["__scopeId","data-v-a80b4e62"]]),Re={class:"form-box"},Le={__name:"reset-password",setup(b,{expose:e}){const i=$(),_=D.useForm,d=I(!0),m=I("重置登录密码"),r=F({password:"",check_password:"",id:i.userInfo.user_id}),n=F({password:[{message:"请输入登录密码",required:!0},{validator:async(S,p)=>!ue(p)&&p?Promise.reject("密码必须包含字母、数字或者字符中的两种，6-32位"):Promise.resolve()}],check_password:[{message:"请输入确认密码",required:!0},{validator:async(S,p)=>p!=r.password&&p?Promise.reject("两次输入的密码不一致"):Promise.resolve()}]}),{resetFields:v,validate:s,validateInfos:l}=_(r,n),f=()=>{s().then(()=>{oe({...X(r)}).then(S=>{se.success("修改成功"),d.value=!1,i.reset(!0)})})},y=()=>{i.setResetPassModal()};return e({open}),(S,p)=>{const V=de,L=pe,B=_e,A=D,O=ae;return o(),c("div",null,[t(O,{open:d.value,"onUpdate:open":p[2]||(p[2]=M=>d.value=M),title:m.value,onOk:f,width:"520px",onCancel:y},{default:u(()=>[t(V,{class:"alert-box",message:"您还未修改初始密码，为了保障您的数据安全，请您尽快重置密码。",type:"info","show-icon":""}),a("div",Re,[t(A,{layout:"vertical"},{default:u(()=>[t(B,U({label:"登录密码"},w(l).password),{default:u(()=>[t(L,{value:r.password,"onUpdate:value":p[0]||(p[0]=M=>r.password=M),placeholder:"密码必须包含字母、数字或者字符中的两种，6-32位"},null,8,["value"])]),_:1},16),t(B,U({label:"确认密码"},w(l).check_password),{default:u(()=>[t(L,{value:r.check_password,"onUpdate:value":p[1]||(p[1]=M=>r.check_password=M),placeholder:"请重新输入密码"},null,8,["value"])]),_:1},16)]),_:1})])]),_:1},8,["open","title"])])}}},Be=k(Le,[["__scopeId","data-v-40aed5c2"]]),Fe={class:"admin-layout"},Ue={class:"layout-header-wrapper"},je={class:"layout-header"},De={class:"header-left"},Ee={class:"header-body"},Ne={class:"header-right"},Ve={class:"item-box"},Ae={class:"item-box"},Oe={class:"layout-body"},ze={key:0,class:"layout-breadcrumb-wrapper"},qe={class:"layout-footer-wrapper"},He={__name:"index",setup(b){const e=E(),i=$(),_=h(()=>e.meta.isCustomPage||!1),d=h(()=>e.meta.pageStyle||{}),m=h(()=>e.meta.bgColor||"#ffffff"),r=h(()=>e.meta.breadcrumb||[]),n=h(()=>e.meta.hideTitle||!1),v=h(()=>e.meta.title||!1),s=h(()=>{var l;return((l=i.userInfo)==null?void 0:l.d_pass)==1&&i.isShowResetPassModal});return(l,f)=>{const y=R("router-view");return o(),c("div",Fe,[a("div",Ue,[a("div",je,[a("div",De,[t(Y)]),a("div",Ee,[t(he)]),a("div",Ne,[a("div",Ve,[t(G)]),a("div",Ae,[t(Te)])])])]),a("div",Oe,[_.value?(o(),x(y,{key:0})):(o(),c("div",{key:1,class:"page-wrapper",style:j({"background-color":m.value})},[n.value?T("",!0):(o(),c("div",ze,[t(we,{items:r.value,title:v.value},null,8,["items","title"])])),a("div",{class:"page-container",style:j({...d.value})},[t(y)],4)],4))]),a("div",qe,[t(Se)]),s.value?(o(),x(Be,{key:0})):T("",!0)])}}},Pt=k(He,[["__scopeId","data-v-1c78c62b"]]);export{Pt as default};