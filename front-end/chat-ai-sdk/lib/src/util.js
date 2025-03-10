export function objectToQueryString(obj) {  
    if (!obj) return '';  
  
    const params = [];  
    for (let key in obj) {  
        if (obj.hasOwnProperty(key)) {  
            let value = obj[key];  
  
            // 如果值是数组或对象，我们需要将其转换为字符串  
            // 这里简单地使用JSON.stringify，但你可能需要更复杂的逻辑来处理嵌套对象或数组  
            if (Array.isArray(value) || typeof value === 'object' && value !== null) {  
                value = JSON.stringify(value);  
            }  
  
            // 确保值不是undefined（因为undefined在查询字符串中没有意义）  
            if (value !== undefined) {  
                params.push(encodeURIComponent(key) + '=' + encodeURIComponent(value));  
            }  
        }  
    }  
  
    return params.join('&');  
}

export function loadCss() {  
    const sdkEl = document.getElementById("ai_chat_js")
    const link = document.createElement("link");  
    const origin = new URL(sdkEl.src).origin
    
    link.type = "text/css";  
    link.rel = "stylesheet";  
    link.href = origin + '/sdk/style.css';
    document.getElementsByTagName("head")[0].appendChild(link);  
}  