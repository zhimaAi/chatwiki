import{aw as a}from"../../index-P5YTNo6V.js";const s=({file:t,category:r})=>a.post({headers:{"Content-Type":"multipart/form-data"},url:"/manage/upload",data:{file:t,category:r}}),o=({openid:t})=>a.get({url:"/chat/getWsUrl",params:{openid:t,debug:0}});export{o as g,s as u};