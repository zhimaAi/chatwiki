// useMathJax.js
import { nextTick } from 'vue';

const MATHJAX_SRC = "./libs/MathJax/es5/tex-mml-chtml.js";

export function useMathJax() {

  // 动态加载脚本的私有函数
  const initMathJax = () => {
    return new Promise((resolve, reject) => {
      if (window.MathJax && window.MathJax.typesetPromise) {
        resolve();
        return;
      }
      const script = document.createElement('script');
      script.src = MATHJAX_SRC;
      script.async = true;
      script.onload = () => resolve();
      script.onerror = () => reject(new Error('MathJax 加载失败'));
      document.head.appendChild(script);
    });
  };

  /**
   * 渲染函数
   * @param {HTMLElement|string} target - 指定元素或选择器
   */
  const renderMath = async (target = null) => {
    try {
      await initMathJax(); // 确保脚本已加载
      await nextTick();

      if (window.MathJax.typesetPromise) {
        const el = (target && target.value) ? target.value : target;

        if (el instanceof HTMLElement) {
          await window.MathJax.typesetPromise([el]);
        } else {
          // 首次进入或未传参时进行全局渲染
          await window.MathJax.typesetPromise();
        }
      }
    } catch (err) {
      console.error('MathJax 渲染出错:', err);
    }
  };

  return { renderMath };
}
