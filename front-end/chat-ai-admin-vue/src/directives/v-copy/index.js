import { message } from 'ant-design-vue'

export const copy = {
  mounted: function (el, binding) {
    const { text, onSuccess } = binding.value

    el.$value = text

    // el控件定义 onclick 事件
    el.handler = () => {
      // if (!el.$value) {
      //   // 值为空的时候，给出提示。可根据项目UI仔细设计
      //   console.log('无复制内容')
      //   return
      // }
      // 动态创建 textarea 标签
      const textarea = document.createElement('textarea')
      // 将该 textarea 设为 readonly 防止 iOS 下自动唤起键盘，同时将 textarea 移出可视区域
      textarea.readOnly = 'readonly'
      textarea.style.position = 'fixed'
      textarea.style.left = '-9999px'
      // 将要 copy 的值赋给 textarea 标签的 value 属性
      textarea.value = el.$value
      // 将 textarea 插入到 body 中
      document.body.appendChild(textarea)
      // 选中值并复制
      textarea.select()
      const result = document.execCommand('Copy')
      if (result) {
        if (typeof onSuccess === 'function') {
          onSuccess()
        } else {
          message.success('复制成功')
        }
      }
      document.body.removeChild(textarea)
    }
    // 绑定点击事件，就是所谓的一键 copy 啦
    el.addEventListener('click', el.handler)
  },
  // 当传进来的值更新的时候触发
  beforeUpdate(el, { value }) {
    el.$value = value.text
  },
  // 指令与元素解绑的时候，移除事件绑定
  unmounted(el) {
    el.removeEventListener('click', el.handler)
  }
}

export default copy
