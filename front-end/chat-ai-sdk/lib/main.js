import { loadCss } from "./src/util";
import { init } from "./src/index";

// 加载css
loadCss();

// 初始化
window.AiChatSDK = init();