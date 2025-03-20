window.onload = function () {
  init();
};

// 动态加载高亮脚本
function loadHighlightScript(language) {
  return new Promise((resolve, reject) => {
    const languageMap = {
      "js": "javascript",
      "ts": "typescript",
    }
    const script = document.createElement("script");

    if(languageMap[language]){
      language = languageMap[language]
    }

    script.src = `/open/static/libs/highlight/languages/${language}.min.js`;
    script.onload = resolve;
    script.onerror = reject;
    document.head.appendChild(script);
  });
}

function init() {
  initToggleCatalog();

  document.querySelectorAll("pre code").forEach((el) => {
    let language = el.className.match(/language-(\w+)/)?.[1];

    // if(!language){
    //   language = 'markdown'
    // }

    if (language) {
      // 动态加载对应的高亮脚本
      loadHighlightScript(language)
        .then(() => {
          // 高亮脚本加载完成后，对块代码进行高亮处理
          hljs.highlightElement(el);
        })
        .catch((err) => {
          console.error(
            "Failed to load highlight script for language:",
            language,
            err
          );
        });
    }else{
      hljs.highlightElement(el);
    }
  });
}

