function handleFocus() {
  $("#search-box").addClass("focus");
  handleInput();
}
function handleBlur() {
  var value = $("#search-input").val();
  if (value) {
    return;
  }

  $("#search-box").removeClass("focus");
  handleInput();
}

function handleInput() {
  var value = $("#search-input").val();

  if (value.length > 0) {
    $("#search-box").addClass("has-value");
  } else {
    $("#search-box").removeClass("has-value");
  }
}

function handleClear() {
  preventLosingFocus = true;
  $("#search-input").val("");
  $("#search-input").focus();
  handleInput();
}

// 生成搜索结果列表
function generateSearchResultList(list) {
  let html = "";
  let keyword = $("#search-input").val() || "";

  for (let i = 0; i < list.length; i++) {
    let item = list[i];
    let title = item.title;
    let content = item.content;

    if (keyword.length) {
      title = item.title.replace(
        new RegExp(keyword, "gi"),
        (match) => `<span class="keyword-text">${match}</span>`
      );
      content = item.content.replace(
        new RegExp(keyword, "gi"),
        (match) => `<span class="keyword-text">${match}</span>`
      );
    }

    html += `<li class="search-result-item-wrap">
              <a class="search-result-item" target="_blank" href="/open/doc/${item.doc_key}">
                <div class="doc-title-box">
                  <i class="doc-icon"></i>
                  <h3 class="doc-title">${title}</h3>
                </div>
                <p class="doc-content">
                ${content}
                </p>
              </a>
            </li>`;
  }
  return html;
}

function handleSearch() {
  let libraryKey = $("#library_key").val();
  let keyword = $("#search-input").val() || "";

  if (!keyword) {
    return;
  }

  $("#search-box").addClass("has-value focus");

  $.ajax({
    url: `/open/search/query/${libraryKey}?v=${encodeURIComponent(keyword)}`,
    type: "get",
    dataType: "json",
    success: function (res) {
      if (res.res != 0) {
        return;
      }
      let list = res.data || [];
      let resultListHtml = generateSearchResultList(list);
      $("#search-result-list").html(resultListHtml);
    },
  });
}
const aiResult = {
  text: "",
  documents: [],
};

const updateAiResultText = function () {
  $(".ai-result-box").show();
  $("#ai-result-text").text(aiResult.text);
};

const updateAiResultDocs = function () {
  let html = "";
  aiResult.documents.forEach((doc) => {
    html += `<li class="ai-result-doc-item-wrap">
                  <a class="ai-result-doc-item" target="_blank" href="/open/doc/${doc.doc_key}">
                    <div class="doc-title-box">
                      <img
                        class="doc-icon"
                        src="/open/static/img/doc_icon2.svg?v=1"
                        alt=""
                      />
                      <h3 class="doc-title">${doc.file_name}</h3>
                    </div>
                    <p class="doc-content" style="display: none">
                    ${doc.content || ''}
                    </p>
                  </a>
                </li>`;
  });
  $('.ai-result-doc-box .label-text').text('相关文档 ('+aiResult.documents.length+')');
  $('#ai-result-docs').html(html);
  $('.ai-result-doc-box').show();
};

function getAiResult() {
  let libraryKey = $("#library_key").val();
  let keyword = $("#search-input").val() || "";

  if (!keyword) {
    return;
  }

  aiResult.text = "";
  aiResult.documents = [];

  const ctrl = new AbortController();
  FetchEventSource.fetchEventSource(
    `/open/summary/${libraryKey}?v=${encodeURIComponent(keyword)}`,
    {
      method: "GET",
      headers: {
        // "Content-Type": "application/json",
      },
      // 允许在页面隐藏时继续接收消息(开启后不再触发自动重连的问题)
      openWhenHidden: true,
      signal: ctrl.signal,
      async onopen(response) {
        if (
          response.ok &&
          response.headers.get("content-type") === "text/event-stream"
        ) {
          return; // everything's good
        } else if (
          response.status >= 400 &&
          response.status < 500 &&
          response.status !== 429
        ) {
          // client-side errors are usually non-retriable:
          throw new Error("连接出错");
        } else {
          throw new Error("连接出错");
        }
      },
      onmessage(res) {
        if (res.event == "sending") {
          aiResult.text += res.data;
          updateAiResultText();
        } else if (res.event == "quote_file") {
          aiResult.documents = JSON.parse(res.data);
          updateAiResultDocs();
        } else if (res.event == "done") {
        }
      },
      onclose() {
        ctrl.abort();
      },
      onerror(err) {
        console.log(err);
        throw err;
      },
    }
  );
}

// 展开收起相关文档
function toggleAiResultDocs() {
  if ($(".ai-result-doc-box").hasClass("expanded")) {
    $(".ai-result-doc-box").removeClass("expanded");
  } else {
    $(".ai-result-doc-box").addClass("expanded");
  }
}

$(function () {
  initToggleCatalog();
  handleSearch();
  getAiResult();
});
