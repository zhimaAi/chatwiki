// prevent losing focus when clicking on the search box
let preventLosingFocus = false;
let isEditStatus = false;
function isPreview() {
  return getUrlParam('ispreview') === '1';
}

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

let isBtnHovering = false;
let isBoxHovering = false;
let actionType = "";
let message = {}

function setMessage(event, type) {
  let target = event.target
  let data = {
    doc_id: $('#doc_id').val()
  }

  if(type == 'title' || type == 'content'){
    data.title = $('#title').val()
    data.content = $('#content').val()
  }

  if(type == 'question'){
    data.question = $(target).data('question')
    data.id = $(target).data('id')
  }

  message = {...data}
}
function onShowFloatBtn(event, type, btns) {
  if(!isEditStatus){
    return;
  }

  setMessage(event, type)

  isBoxHovering = true;
  actionType = type;

  let floatBtnWarpper = document.getElementById("floatBtnWarpper");
  $(floatBtnWarpper).addClass("show");
  btns.forEach(function (btn) {
    floatBtnWarpper.querySelector(
      `.action-btn[data-type="${btn}"]`
    ).style.display = "flex";
  });

  floatBtnWarpper.style.display = "flex";

  const rect = event.target.getBoundingClientRect();

  const btnWidth = floatBtnWarpper.offsetWidth;
  const btnHeight = floatBtnWarpper.offsetHeight;

  const left = rect.left + rect.width / 2 - btnWidth / 2;
  const top = rect.top - btnHeight + 2;

  floatBtnWarpper.style.left = left + "px";
  floatBtnWarpper.style.top = top + "px";
}

function onHideFloatBtn (){
  isBoxHovering = false;

  hideFloatBtn();
}


let timer = null
function hideFloatBtn(event) {
  if(!isEditStatus){
    return;
  }
  
  if(timer){
    clearTimeout(timer)
  }

  timer = setTimeout(function () {
    if (isBtnHovering || isBoxHovering) {
      return;
    }

    isBtnHovering = false;
    actionType = "";

    let floatBtnWarpper = document.getElementById("floatBtnWarpper");

    $(floatBtnWarpper).removeClass("show");

    floatBtnWarpper.querySelectorAll(".action-btn").forEach(function (btn) {
      btn.style.display = "none";
    });

    floatBtnWarpper.style.display = "none";
    floatBtnWarpper.style.left = "-9999px";
    floatBtnWarpper.style.top = "-9999px";
  }, 50);
}

function onFloatBtnHover(event, status) {
  // 阻止事件冒泡和默认行为
  event.preventDefault();
  event.stopPropagation();

  if (status) {
    isBtnHovering = true;
  } else {
    isBtnHovering = false;
    hideFloatBtn();
  }
}

function clickQuestionGuide(e){
  let libraryKey = $("#library_key").val();
  let keyword = $(e.target).data('question');

  window.location.href = `/open/search/html/${libraryKey}?v=${encodeURIComponent(keyword)}`;
}

const postMessage = (actionKey, actionType) => {
  if (!window.parent) {
    return;
  }

  const data = {
    action: actionType,
    key: actionKey,
    data: message,
  };

  window.parent.postMessage(data, "*");
};

function handleActionBtn(actionKey) {
  postMessage(actionType, actionKey);
}

function handleAddQuestionBtn(){
  message = {
    id: '',
    question: '',
  }
  postMessage('question', 'add');
}

function handleDeleteQuestionBtn(){
  message = {
    id: '',
    question: '',
  }
  
  postMessage('question', 'delete');
}

function initMessage(){
  let data = {
    doc_id: $('#doc_id').val(),
    title: $('#title').val(),
    content: $('#content').val()
  }

  message = {...data}

  postMessage('init', 'init');
}

function setEditStatus(event){
  isEditStatus = event.data.type;
  if(isEditStatus){
    let questionGuideItems = document.querySelectorAll('#keyword-list .question-guide-item') || []

    if(questionGuideItems.length < 10){
      $('#keyword-list .add-btn-item').show();
    }
    let library_desc = $('#library-desc').text()

    if(!library_desc || library_desc.length === 0){
      $('#library-desc').text('请输入描述');
    }
  }else{
    isBtnHovering = false;
    isBoxHovering = false;
    actionType = "";

    let floatBtnWarpper = document.getElementById("floatBtnWarpper");

    $(floatBtnWarpper).removeClass("show");

    floatBtnWarpper.querySelectorAll(".action-btn").forEach(function (btn) {
      btn.style.display = "none";
    });

    floatBtnWarpper.style.display = "none";
    floatBtnWarpper.style.left = "-9999px";
    floatBtnWarpper.style.top = "-9999px";

    $('#keyword-list .add-btn-item').hide();
  }
}

$(function () {
  initMessage();
  initToggleCatalog();

  // 监听父页面发送的setEditStatus事件
  window.addEventListener('message', function(event) {
    if (event.data && event.data.action === 'setEditStatus') {
      setEditStatus(event.data);
    }
  });
});
