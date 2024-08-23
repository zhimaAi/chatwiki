import{c as R}from"./vue-chunks-DrvJJrR0.js";var U={exports:{}},V;function tt(){return V||(V=1,function(B,K){(function(C,k){B.exports=k()})(R,function(){var C=1e3,k=6e4,F=36e5,A="millisecond",S="second",w="minute",O="hour",M="day",T="week",m="month",J="quarter",v="year",_="date",Z="Invalid Date",E=/^(\d{4})[-/]?(\d{1,2})?[-/]?(\d{0,2})[Tt\s]*(\d{1,2})?:?(\d{1,2})?:?(\d{1,2})?[.:]?(\d+)?$/,G=/\[([^\]]+)]|Y{1,4}|M{1,4}|D{1,2}|d{1,4}|H{1,2}|h{1,2}|a|A|m{1,2}|s{1,2}|Z{1,2}|SSS/g,P={name:"en",weekdays:"Sunday_Monday_Tuesday_Wednesday_Thursday_Friday_Saturday".split("_"),months:"January_February_March_April_May_June_July_August_September_October_November_December".split("_"),ordinal:function(i){var n=["th","st","nd","rd"],t=i%100;return"["+i+(n[(t-20)%10]||n[t]||n[0])+"]"}},I=function(i,n,t){var r=String(i);return!r||r.length>=n?i:""+Array(n+1-r.length).join(t)+i},Q={s:I,z:function(i){var n=-i.utcOffset(),t=Math.abs(n),r=Math.floor(t/60),e=t%60;return(n<=0?"+":"-")+I(r,2,"0")+":"+I(e,2,"0")},m:function i(n,t){if(n.date()<t.date())return-i(t,n);var r=12*(t.year()-n.year())+(t.month()-n.month()),e=n.clone().add(r,m),s=t-e<0,u=n.clone().add(r+(s?-1:1),m);return+(-(r+(t-e)/(s?e-u:u-e))||0)},a:function(i){return i<0?Math.ceil(i)||0:Math.floor(i)},p:function(i){return{M:m,y:v,w:T,d:M,D:_,h:O,m:w,s:S,ms:A,Q:J}[i]||String(i||"").toLowerCase().replace(/s$/,"")},u:function(i){return i===void 0}},x="en",D={};D[x]=P;var q="$isDayjsObject",N=function(i){return i instanceof L||!(!i||!i[q])},j=function i(n,t,r){var e;if(!n)return x;if(typeof n=="string"){var s=n.toLowerCase();D[s]&&(e=s),t&&(D[s]=t,e=s);var u=n.split("-");if(!e&&u.length>1)return i(u[0])}else{var o=n.name;D[o]=n,e=o}return!r&&e&&(x=e),e||!r&&x},f=function(i,n){if(N(i))return i.clone();var t=typeof n=="object"?n:{};return t.date=i,t.args=arguments,new L(t)},a=Q;a.l=j,a.i=N,a.w=function(i,n){return f(i,{locale:n.$L,utc:n.$u,x:n.$x,$offset:n.$offset})};var L=function(){function i(t){this.$L=j(t.locale,null,!0),this.parse(t),this.$x=this.$x||t.x||{},this[q]=!0}var n=i.prototype;return n.parse=function(t){this.$d=function(r){var e=r.date,s=r.utc;if(e===null)return new Date(NaN);if(a.u(e))return new Date;if(e instanceof Date)return new Date(e);if(typeof e=="string"&&!/Z$/i.test(e)){var u=e.match(E);if(u){var o=u[2]-1||0,c=(u[7]||"0").substring(0,3);return s?new Date(Date.UTC(u[1],o,u[3]||1,u[4]||0,u[5]||0,u[6]||0,c)):new Date(u[1],o,u[3]||1,u[4]||0,u[5]||0,u[6]||0,c)}}return new Date(e)}(t),this.init()},n.init=function(){var t=this.$d;this.$y=t.getFullYear(),this.$M=t.getMonth(),this.$D=t.getDate(),this.$W=t.getDay(),this.$H=t.getHours(),this.$m=t.getMinutes(),this.$s=t.getSeconds(),this.$ms=t.getMilliseconds()},n.$utils=function(){return a},n.isValid=function(){return this.$d.toString()!==Z},n.isSame=function(t,r){var e=f(t);return this.startOf(r)<=e&&e<=this.endOf(r)},n.isAfter=function(t,r){return f(t)<this.startOf(r)},n.isBefore=function(t,r){return this.endOf(r)<f(t)},n.$g=function(t,r,e){return a.u(t)?this[r]:this.set(e,t)},n.unix=function(){return Math.floor(this.valueOf()/1e3)},n.valueOf=function(){return this.$d.getTime()},n.startOf=function(t,r){var e=this,s=!!a.u(r)||r,u=a.p(t),o=function(p,$){var y=a.w(e.$u?Date.UTC(e.$y,$,p):new Date(e.$y,$,p),e);return s?y:y.endOf(M)},c=function(p,$){return a.w(e.toDate()[p].apply(e.toDate("s"),(s?[0,0,0,0]:[23,59,59,999]).slice($)),e)},h=this.$W,d=this.$M,l=this.$D,b="set"+(this.$u?"UTC":"");switch(u){case v:return s?o(1,0):o(31,11);case m:return s?o(1,d):o(0,d+1);case T:var g=this.$locale().weekStart||0,Y=(h<g?h+7:h)-g;return o(s?l-Y:l+(6-Y),d);case M:case _:return c(b+"Hours",0);case O:return c(b+"Minutes",1);case w:return c(b+"Seconds",2);case S:return c(b+"Milliseconds",3);default:return this.clone()}},n.endOf=function(t){return this.startOf(t,!1)},n.$set=function(t,r){var e,s=a.p(t),u="set"+(this.$u?"UTC":""),o=(e={},e[M]=u+"Date",e[_]=u+"Date",e[m]=u+"Month",e[v]=u+"FullYear",e[O]=u+"Hours",e[w]=u+"Minutes",e[S]=u+"Seconds",e[A]=u+"Milliseconds",e)[s],c=s===M?this.$D+(r-this.$W):r;if(s===m||s===v){var h=this.clone().set(_,1);h.$d[o](c),h.init(),this.$d=h.set(_,Math.min(this.$D,h.daysInMonth())).$d}else o&&this.$d[o](c);return this.init(),this},n.set=function(t,r){return this.clone().$set(t,r)},n.get=function(t){return this[a.p(t)]()},n.add=function(t,r){var e,s=this;t=Number(t);var u=a.p(r),o=function(d){var l=f(s);return a.w(l.date(l.date()+Math.round(d*t)),s)};if(u===m)return this.set(m,this.$M+t);if(u===v)return this.set(v,this.$y+t);if(u===M)return o(1);if(u===T)return o(7);var c=(e={},e[w]=k,e[O]=F,e[S]=C,e)[u]||1,h=this.$d.getTime()+t*c;return a.w(h,this)},n.subtract=function(t,r){return this.add(-1*t,r)},n.format=function(t){var r=this,e=this.$locale();if(!this.isValid())return e.invalidDate||Z;var s=t||"YYYY-MM-DDTHH:mm:ssZ",u=a.z(this),o=this.$H,c=this.$m,h=this.$M,d=e.weekdays,l=e.months,b=e.meridiem,g=function($,y,H,W){return $&&($[y]||$(r,s))||H[y].slice(0,W)},Y=function($){return a.s(o%12||12,$,"0")},p=b||function($,y,H){var W=$<12?"AM":"PM";return H?W.toLowerCase():W};return s.replace(G,function($,y){return y||function(H){switch(H){case"YY":return String(r.$y).slice(-2);case"YYYY":return a.s(r.$y,4,"0");case"M":return h+1;case"MM":return a.s(h+1,2,"0");case"MMM":return g(e.monthsShort,h,l,3);case"MMMM":return g(l,h);case"D":return r.$D;case"DD":return a.s(r.$D,2,"0");case"d":return String(r.$W);case"dd":return g(e.weekdaysMin,r.$W,d,2);case"ddd":return g(e.weekdaysShort,r.$W,d,3);case"dddd":return d[r.$W];case"H":return String(o);case"HH":return a.s(o,2,"0");case"h":return Y(1);case"hh":return Y(2);case"a":return p(o,c,!0);case"A":return p(o,c,!1);case"m":return String(c);case"mm":return a.s(c,2,"0");case"s":return String(r.$s);case"ss":return a.s(r.$s,2,"0");case"SSS":return a.s(r.$ms,3,"0");case"Z":return u}return null}($)||u.replace(":","")})},n.utcOffset=function(){return 15*-Math.round(this.$d.getTimezoneOffset()/15)},n.diff=function(t,r,e){var s,u=this,o=a.p(r),c=f(t),h=(c.utcOffset()-this.utcOffset())*k,d=this-c,l=function(){return a.m(u,c)};switch(o){case v:s=l()/12;break;case m:s=l();break;case J:s=l()/3;break;case T:s=(d-h)/6048e5;break;case M:s=(d-h)/864e5;break;case O:s=d/F;break;case w:s=d/k;break;case S:s=d/C;break;default:s=d}return e?s:a.a(s)},n.daysInMonth=function(){return this.endOf(m).$D},n.$locale=function(){return D[this.$L]},n.locale=function(t,r){if(!t)return this.$L;var e=this.clone(),s=j(t,r,!0);return s&&(e.$L=s),e},n.clone=function(){return a.w(this.$d,this)},n.toDate=function(){return new Date(this.valueOf())},n.toJSON=function(){return this.isValid()?this.toISOString():null},n.toISOString=function(){return this.$d.toISOString()},n.toString=function(){return this.$d.toUTCString()},i}(),z=L.prototype;return f.prototype=z,[["$ms",A],["$s",S],["$m",w],["$H",O],["$W",M],["$M",m],["$y",v],["$D",_]].forEach(function(i){z[i[1]]=function(n){return this.$g(n,i[0],i[1])}}),f.extend=function(i,n){return i.$i||(i(n,L,f),i.$i=!0),f},f.locale=j,f.isDayjs=N,f.unix=function(i){return f(1e3*i)},f.en=D[x],f.Ls=D,f.p={},f})}(U)),U.exports}export{tt as r};