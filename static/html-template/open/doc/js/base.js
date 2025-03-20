function btoaUTF8UsingTextEncoder(str) {
  const encoder = new TextEncoder();
  const data = encoder.encode(str);
  let binaryString = '';
  for (let i = 0; i < data.length; i++) {
      binaryString += String.fromCharCode(data[i]);
  }
  return btoa(binaryString);
}

// 获取URL参数
const getUrlParam = function(name) {
  const searchParams = new URLSearchParams(window.location.search);
  return searchParams.get(name);
}

// 检查是不是在预览
let is_preview = false
let token = ''


const setPreview = function(data){
  is_preview = true
  let homeLink = $('#home-link').attr('href')
  // /open/home替换为/manage/libDocHome
  if(homeLink){
    let newLink = homeLink.replace('/open/home', '/manage/libDocHome')
    $('#home-link').attr('href', `${newLink}?ispreview=1&token=${data.token}`)
  }

  let menus = $('#directory-list .doc-body');
  
  menus.each(function() {
    let currentHref = $(this).attr('href');
    if(currentHref) {
      $(this).attr('href', `${currentHref}?ispreview=1&token=${data.token}`);
    }
  });
}

window.parent.postMessage({
  action: 'check_preview'
}, '*');

window.addEventListener('message', function(event) {
  const data = event.data

  if(data.action === 'setPreview'){
    setPreview(data)
  }
});


// 初始化目录相关功能
const closeAllCatalog = function(){
  let menus = document.querySelectorAll('#directory-list .directory-item')

  for(let menu of menus){
    let $menu = $(menu);
    let isExpanded = $(menu).hasClass("is-expanded");

    if (isExpanded) {
      let menuId = $menu.data("id");
       $(menu).removeClass("is-expanded");
      $('.directory-item[data-pid="' + menuId + '"]').hide();
    }
  }
}

const initToggleCatalog = function () {
  let docKey = $('#doc-key').val() || '';
  $('.directory-item[data-doc-key="'+docKey+'"]').addClass('active')

  $(".directory-list").on("click", ".directory-expanded-btn", function () {
    if ($(this).hasClass("directory-expanded-noop")) {
      return;
    }

    let $directoryItem = $(this).closest(".directory-item");
    // let pid = $directoryItem.data('pid')
    let id = $directoryItem.data("id");
    let isExpanded = $directoryItem.hasClass("is-expanded");

    // 根据data-pid来控制展开收起
    if (isExpanded) {
      $directoryItem.removeClass("is-expanded");
      $('.directory-item[data-pid="' + id + '"]').hide();
    } else {
      $directoryItem.addClass("is-expanded");
      $(".directory-item[data-pid=" + id + "]").show();
    }
  });

  // 收起所有目录

};

// 侧边栏菜单展开收起
function toggleSidebar() {
  const wikiSidebarWrapper = $('#wikiSidebarWrapper')

  if (wikiSidebarWrapper.hasClass('wiki-sidebar-open')) {
    wikiSidebarWrapper.removeClass('wiki-sidebar-open')
    setTimeout(() => {
      $('#wikiSidebarWrapper .wiki-sidebar-mask').hide()
    }, 200)
  }else {
    $('#wikiSidebarWrapper .wiki-sidebar-mask').show()
    wikiSidebarWrapper.addClass('wiki-sidebar-open')
  }
}

// 侧边栏搜索
function toSearchPage(keyword){
  let libraryKey = $("#library_key").val();
  window.location.href = `/open/search/html/${libraryKey}?v=${encodeURIComponent(keyword)}`;
}
function onSidebarSearch() {
  let keyword = $("#sidebar-search-input").val();
  if (!keyword) {
    return;
  }
  toSearchPage(keyword)
}

function handleEnterSidebarSearch(event){
  if (event.keyCode === 13) {
    onSidebarSearch();
  }
}

function onSearch(){
  var keyword = $("#search-input").val();
  if (!keyword) {
    return;
  }
  
  toSearchPage(keyword)
}

function handleEnterSearch(event){
  if (event.keyCode === 13) {
    onSearch();
  }
}

$(function(){

})