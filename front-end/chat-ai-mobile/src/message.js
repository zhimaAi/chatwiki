import { useChatStore } from '@/stores/modules/chat'
window.addEventListener("message", async (e) => {
  var { data, source } = e;
  // console.log(e,'e==')
  if (source == window.parent) {
    const chatStore = useChatStore()
    const { upDataUiStyle, updataQuickComand } = chatStore
    var type = data.type;
    switch (type) {
      case "onPreview": {
        upDataUiStyle(data.data)
        break;
      }
      case "updataQuickComand": {
        updataQuickComand(data.data)
        break;
      }
    }
  }
});
