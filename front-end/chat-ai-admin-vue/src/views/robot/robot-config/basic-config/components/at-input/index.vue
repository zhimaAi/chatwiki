<template>
  <div style="width: 100%; position: relative">
    <div></div>
    <Teleport :to="popupContainer">
      <div
        class="dropdown-list"
        :class="[dropdownPositionClass, { 'dropdown-list-semantic': semanticTags }]"
        ref="dropdownList"
        v-if="showDropdown"
        :style="positionStyle"
      >
        <slot name="dropdown" :options="showOptions" :select="onSelectOption" :cardOptions="cardOptions" :selectCard="selectCard" :close="hideDropdownMenus">
          <CascadePanel
            :options="showOptions"
            :checkAnyLevel="checkAnyLevel"
            @change="onSelectOption"
            @direction-change="handleDirectionChange"
          />
        </slot>
      </div>
    </Teleport>
    <div class="mention-input-warpper mention-promt-input" :class="{'no-border': noBorder, 'semantic-mention': semanticTags}" ref="JMentionContainer">
      <div
        ref="JMention"
        :style="[inputStyle]"
        class="j-mention"
        :class="[{ 'show-placeholder': localValue.length == 0 }, 'type-' + type]"
        :contenteditable="readonly ? 'false' : semanticTags ? 'true' : 'plaintext-only'"
        :aria-readonly="readonly"
        :placeholder="placeholder"
        @focus="onFocus"
        @input="debouncedUpdate"
        @click="clickJMention"
        @blur="handleBlur"
        @keydown="handleKeydown"
        @paste="handlePaste"
      ></div>
      <span class="placeholder" v-if="localValue.length == 0 || localValue == '\n'">{{ placeholder }}</span>
    </div>
  </div>
</template>

<script>
import { useI18n } from '@/hooks/web/useI18n'
import CascadePanel from './cascade-panel.vue'
import { createPromptToken, ensurePromptCaretSlots, getPromptSelection, readPromptText, renderPromptText, restorePromptSelection, tokenizePrompt } from '@/views/customer-service/basic-settings/prompt-tokens'

// 把树状的options抓换成编排的数组
const getTreeOptions = (options) => {
  let result = []
  options.forEach((opt) => {
    if (opt.children && opt.children.length > 0) {
      result.push(...getTreeOptions(opt.children))
    } else {
      result.push(opt)
    }
  })
  return result
}

export default {
  name: 'JMention',
  components: {
    CascadePanel
  },
  inject: [],
  emits: ['input', 'change', 'update:selectedList', 'focus', 'open'],
  props: {
    semanticTags: { type: Boolean, default: false },
    readonly: { type: Boolean, default: false },
    type: {
      type: String,
      default: 'input', // input, textarea
      validator(value) {
        return ['input', 'textarea'].includes(value)
      }
    },
    placeholder: {
      type: String,
      default: ''
    },
    defaultValue: {
      type: [String, Number],
      default: ''
    },
    options: {
      type: Array,
      default: () => []
    },
    defaultSelectedList: {
      type: Array,
      default: () => []
    },
    canRepeat: {
      type: Boolean,
      default: true
    },
    isContain: {
      type: Boolean,
      default: false
    },
    inputStyle: {
      type: String,
      default: ''
    },
    checkAnyLevel: {
      type: Boolean,
      default: false
    },
    noBorder: {
      type: Boolean,
      default: false
    },
    cardOptions: {
      type: Array,
      default: () => []
    },
    enableCardInsert: {
      type: Boolean,
      default: false
    },
    getPopupContainer: {
      type: Function,
      default: null
    }
  },
  setup() {
    const { t } = useI18n('views.robot.robot-config.basic-config.components.at-input.index')
    return {
      t
    }
  },
  data() {
    return {
      timer: null,
      showDropdown: false,
      positionStyle: '',
      selectedList: [],
      selectedIdSet: new Set(),
      showOptions: [],
      atText: '',
      localValue: ''
    }
  },
  computed: {
    popupContainer() {
      return this.resolvePopupContainer()
    },
    isBodyPopupContainer() {
      return this.popupContainer === document.body
    },
    dropdownPositionClass() {
      return this.isBodyPopupContainer ? 'dropdown-list-fixed' : 'dropdown-list-absolute'
    }
  },
  watch: {
    readonly(value) {
      if (value) this.hideDropdownMenus()
    },
    defaultValue(val) {
      if (val.trim() != this.localValue.trim()) {
        this.refresh()
      }
    },
    showDropdown(val) {
      this.onOpen(val)
      if (this.semanticTags) {
        this.popupResizeObserver?.disconnect()
        if (val) this.$nextTick(() => {
          const popup = this.$refs.dropdownList
          if (popup) this.popupResizeObserver?.observe(popup.firstElementChild || popup)
          this.positionSemanticDropdown()
        })
      }
    },
    options: { deep: true, handler() { this.initShowOptionList(); this.refreshSemanticLabels() } },
    cardOptions: { deep: true, handler() { this.refreshSemanticLabels() } }
  },
  mounted() {
    this.initShowOptionList()

    document.addEventListener('dragstart', this.handleDragstart)
    document.addEventListener('click', this.handleClick)
    document.addEventListener('scroll', this.onScroll, true)

    this.initData()
    if (this.semanticTags) {
      window.addEventListener('resize', this.positionSemanticDropdown)
      if (typeof ResizeObserver !== 'undefined') this.popupResizeObserver = new ResizeObserver(this.positionSemanticDropdown)
    }
  },
  beforeUnmount() {
    if (this.semanticTags) clearTimeout(this.timer)
    this.popupResizeObserver?.disconnect()
    window.removeEventListener('resize', this.positionSemanticDropdown)
    document.removeEventListener('dragstart', this.handleDragstart)
    document.removeEventListener('click', this.handleClick)
    document.removeEventListener('scroll', this.onScroll, true)
  },
  methods: {
    refreshSemanticLabels() {
      if (this.semanticTags && this.$refs.JMention) this.renderSemanticValue(readPromptText(this.$refs.JMention))
    },
    renderSemanticValue(value) {
      const options = [...getTreeOptions(this.options), ...getTreeOptions(this.cardOptions)]
      renderPromptText(this.$refs.JMention, value, options)
      this.selectedList = tokenizePrompt(value, options).filter(part => part.token).map(part => part.token)
      this.selectedIdSet = new Set(this.selectedList.map(option => option.value))
      this.localValue = value
    },
    handlePaste(event) {
      if (!this.semanticTags) return
      event.preventDefault()
      const selection = window.getSelection()
      if (!selection.rangeCount || !this.$refs.JMention.contains(selection.anchorNode)) return
      const range = selection.getRangeAt(0)
      const text = document.createTextNode(event.clipboardData.getData('text/plain'))
      range.deleteContents()
      range.insertNode(text)
      range.setStartAfter(text)
      range.collapse(true)
      selection.removeAllRanges()
      selection.addRange(range)
      this.updateValue()
    },
    resolvePopupContainer() {
      const container = this.getPopupContainer ? this.getPopupContainer() : null
      return container instanceof HTMLElement ? container : document.body
    },
    handleClick(e) {
      const target = e.target.closest('.dropdown-list')
      if (!target) {
        this.hideDropdownMenus()
      }
    },
    handleDragstart(e) {
      const target = e.target.closest('.dropdown-list')
      if (target) {
        e.preventDefault()
      }
    },
    onScroll(event) {
      if (this.semanticTags) {
        if (!this.$refs.dropdownList?.contains(event.target)) this.positionSemanticDropdown()
        return
      }
      if (!this.isBodyPopupContainer && !this.semanticTags) {
        return
      }
      if (this.$refs.dropdownList?.contains(event.target)) {
        return
      }
      this.hideDropdownMenus()
    },
    hideDropdownMenus() {
      this.showDropdown = false
    },
    positionSemanticDropdown() {
      const popup = this.$refs.dropdownList
      const range = this.dropdownAnchorRange
      if (!this.semanticTags || !this.showDropdown || !popup || !range || !this.$refs.JMention?.contains(range.startContainer)) return
      const raw = range.getBoundingClientRect()
      const rect = raw.bottom > raw.top ? raw : this.$refs.JMention.getBoundingClientRect()
      const margin = 8
      const gap = 4
      const width = Math.min(popup.offsetWidth, window.innerWidth - margin * 2)
      const height = popup.scrollHeight || popup.offsetHeight
      const below = Math.max(0, window.innerHeight - margin - rect.bottom - gap)
      const above = Math.max(0, rect.top - margin - gap)
      const useBelow = height <= below || below >= above
      const maxHeight = Math.min(useBelow ? below : above, window.innerHeight - margin * 2)
      const fittedHeight = Math.min(height, maxHeight)
      const top = useBelow ? rect.bottom + gap : rect.top - gap - fittedHeight
      this.positionStyle = {
        left: `${Math.max(margin, Math.min(rect.left, window.innerWidth - width - margin))}px`,
        top: `${Math.max(margin, Math.min(top, window.innerHeight - fittedHeight - margin))}px`,
        maxHeight: `${maxHeight}px`,
        maxWidth: `calc(100vw - ${margin * 2}px)`
      }
    },
    onFocus() {
      this.$emit('focus')
    },
    onOpen(val) {
      this.$emit('open', val)
    },
    getText(htmlStr) {
      const div = document.createElement('div')
      div.innerHTML = htmlStr
      if (this.semanticTags) return readPromptText(div)
      let text = ''
      // 遍历子节点
      Array.from(div.childNodes).forEach((node) => {
        let textContent = ''
        if (this.type === 'textarea') {
          // textarea
          if (node.nodeName === 'BR') {
            // 如果是BR节点,直接添加换行
            textContent = '\n'
          } else {
            if (node.classList && node.classList.contains('j-mention-at')) {
              textContent = node.dataset.value
            } else {
              textContent = node.textContent
            }
          }
        } else {
          // input
          if (node.classList && node.classList.contains('j-mention-at')) {
            textContent = node.dataset.value
          } else if (node.nodeName === 'BR') {
            textContent = ''
          } else {
            textContent = node.textContent.replace(/\n/g, '')
          }
        }

        // 去除只有&nbsp;字符的节点
        if (textContent === '&nbsp;' || textContent === ' ') {
          textContent = ''
        }

        // 拼接文本内容
        text += textContent
      })

      return text
    },
    onChange() {
      const jMention = this.$refs.JMention

      if (!jMention) return

      let htmlStr = jMention.innerHTML
      let text = this.getText(htmlStr)

      this.localValue = text

      let selectedListMap = {}

      this.selectedList.forEach((item) => {
        if (!selectedListMap[item.value] && text.indexOf(item.value) !== -1) {
          selectedListMap[item.value] = item
        }
      })

      let mySelectedList = Object.values(selectedListMap)

      this.$emit('change', text.trim(), mySelectedList)
    },
    refresh() {
      this.selectedList = []
      this.selectedIdSet = new Set()

      this.initShowOptionList()
      this.initData()
    },
    initData() {
      const JMention = this.$refs.JMention
      if (this.semanticTags) {
        this.renderSemanticValue(String(this.defaultValue ?? ''))
        return
      }
      // let html = this.defaultValue;
      let html = this.defaultValue.toString().replace(/</g, '&lt;').replace(/>/g, '&gt;')
      let defaultSelectedList = this.defaultSelectedList || []
      let treeOptions = getTreeOptions(this.options)
      let cardTreeOptions = getTreeOptions(this.cardOptions)
      let selectedList = [...defaultSelectedList, ...treeOptions, ...cardTreeOptions].filter(
        (item) => item && item.value && item.value != '' && item.typ !== 'node'
      )

      // 记录已经替换过的的value
      const replacedValues = new Set()

      selectedList.forEach((opt) => {
        // 使用正则匹配完整的value,避免部分匹配
        const regex = new RegExp(`(${opt.value})`, 'g')
        if (regex.test(this.defaultValue) && !replacedValues.has(opt.value)) {
          // 判断是否为小程序卡片类型
          const isCard = opt.value && opt.value.startsWith('【mini_card:')
          const tagClass = isCard ? 'j-mention-at mini-card-tag' : 'j-mention-at'
          // 小程序卡片复显时需包含图标
          const iconHtml = isCard
            ? '<span class="svg-action" style="font-size: 12px;"><svg aria-hidden="true" class="svg-icon" fill="currentColor" width="1em" height="1em"><use xlink:href="#applet" fill="currentColor"></use></svg></span>'
            : ''

          html = html.replace(
            regex,
            `<span class="${tagClass}" data-id="${opt.id}" contentEditable="false" data-value="${opt.value}" data-type="${isCard ? 'card' : ''}">${iconHtml}${opt.label}</span>`
          )

          replacedValues.add(opt.value)

          this.selectOption(opt, true)
        }
      })

      JMention.innerHTML = html
      this.localValue = this.defaultValue
    },
    initShowOptionList() {
      this.showOptions = this.options.filter((opt) => {
        if (this.atText) {
          return this.semanticTags ? `${opt.label} ${opt.description || ''}`.includes(this.atText) : opt.label.startsWith(this.atText)
        }
        if (this.canRepeat) return true
        return !this.selectedIdSet.has(opt.id + '')
      })
    },
    getLastChar(len = 1) {
      const selection = window.getSelection()
      if (selection.rangeCount <= 0) return
      const range = selection.getRangeAt(0)
      const startOffset = range.startOffset
      const startContainer = range.startContainer
      let textBeforeCursor = ''
      if (startContainer.nodeType === Node.TEXT_NODE) {
        textBeforeCursor = startContainer.textContent.slice(0, startOffset)
      } else if (startContainer.nodeType === Node.ELEMENT_NODE) {
        // 若光标在元素节点内，获取其内部文本
        const childNodes = startContainer.childNodes
        for (let i = 0; i < childNodes.length; i++) {
          const node = childNodes[i]
          if (node.nodeType === Node.TEXT_NODE) {
            if (i === childNodes.length - 1) {
              textBeforeCursor += node.textContent.slice(0, startOffset)
            } else {
              textBeforeCursor += node.textContent
            }
          }
        }
      }
      return textBeforeCursor.slice(-len)
    },
    getLastAtText() {
      const selection = window.getSelection()
      if (selection.rangeCount <= 0) return
      const range = selection.getRangeAt(0)
      const startOffset = range.startOffset
      const startContainer = range.startContainer
      let textBeforeCursor = ''
      if (startContainer.nodeType === Node.TEXT_NODE) {
        textBeforeCursor = startContainer.textContent.slice(0, startOffset)
      } else if (startContainer.nodeType === Node.ELEMENT_NODE) {
        // 若光标在元素节点内，获取其内部文本
        const childNodes = startContainer.childNodes
        for (let i = 0; i < childNodes.length; i++) {
          const node = childNodes[i]
          if (node.nodeType === Node.TEXT_NODE) {
            if (i === childNodes.length - 1) {
              textBeforeCursor += node.textContent.slice(0, startOffset)
            } else {
              textBeforeCursor += node.textContent
            }
          }
        }
      }
      const lastAtIndex = textBeforeCursor.lastIndexOf('/')
      if (lastAtIndex === -1) return ''
      return textBeforeCursor.slice(lastAtIndex + 1)
    },
    async getCaretPosition() {
      // const JMention = this.$refs.JMention;
      const selection = window.getSelection()
      if (selection.rangeCount <= 0) return { top: 0, left: 0 }
      const range = selection.getRangeAt(0)
      const rect = range.getBoundingClientRect()
      if (this.semanticTags) this.dropdownSelection = getPromptSelection(this.$refs.JMention)
      const dropdownListRect = await this.getDropdownListSize()
      const viewportWidth = window.innerWidth
      const viewportHeight = window.innerHeight

      let top = rect.bottom
      let left = rect.left

      // 自动垂直定位：如果下方空间不足，则移动到上方
      if (top + dropdownListRect.height > viewportHeight) {
        top = rect.top - dropdownListRect.height
      }

      // 自动水平定位：如果右侧空间不足，则向左对齐
      if (left + dropdownListRect.width > viewportWidth) {
        left = rect.right - dropdownListRect.width
      }

      // 确保不会超出视口顶部和左部
      if (top < 0) {
        top = 0
      }
      if (left < 0) {
        left = 0
      }

      if (!this.isBodyPopupContainer) {
        const containerRect = this.popupContainer.getBoundingClientRect()
        top = top - containerRect.top + this.popupContainer.scrollTop
        left = left - containerRect.left + this.popupContainer.scrollLeft

        if (top < 0) {
          top = 0
        }
        if (left < 0) {
          left = 0
        }
      }

      return {
        top,
        left
      }
    },
    async getDropdownListSize() {
      this.showDropdown = true // 临时显示元素
      await this.$nextTick() // 等待 DOM 更新
      const width = this.$refs.dropdownList.offsetWidth
      const height = this.$refs.dropdownList.offsetHeight
      this.showDropdown = false // 恢复隐藏
      return { width, height }
    },
    updateSelectedList() {
      const JMention = this.$refs.JMention
      const nodeList = JMention.querySelectorAll('.j-mention-at')
      this.selectedList = []
      this.selectedIdSet = new Set()

      for (const node of nodeList) {
        // const { value } = node.dataset;
        let value = node.dataset.value
        const item = this.options.find((opt) => {
          return opt.value == value
        })

        this.selectedList.push(item)
        this.selectedIdSet.add(value)
      }
      this.$emit('update:selectedList', this.selectedList)
      if (nodeList.length === this.selectedList.length && !this.atText && !this.showDropdown) {
        return
      }

      this.initShowOptionList()
    },
    checkStartWith(text) {
      if (this.semanticTags) return [...this.options, ...this.cardOptions].some(option => `${option.label} ${option.description || ''}`.includes(text))
      const list = this.options.filter((opt) => {
        return opt.label.startsWith(text)
      })
      return list.length > 0
    },
    updateValue() {
      this.onChange()
    },
    clickJMention(event) {
      if (this.semanticTags) {
        const editor = this.$refs.JMention
        const token = event.target.closest('.j-mention-at')
        const selection = window.getSelection()
        let slot
        if (token && editor.contains(token)) {
          ensurePromptCaretSlots(editor)
          const rect = token.getBoundingClientRect()
          if (event.clientX >= rect.left + rect.width / 2) slot = token.nextSibling
        } else if (selection.anchorNode === editor && selection.anchorOffset === editor.childNodes.length) slot = editor.lastChild
        if (slot?.dataset?.promptCaret) {
          editor.focus()
          const range = document.createRange()
          range.setStart(slot.firstChild, 1)
          range.collapse(true)
          selection.removeAllRanges()
          selection.addRange(range)
        }
      }
      if (this.showDropdown) {
        this.showDropdown = false
      }
    },
    debouncedUpdate(e) {
      // 1. 优先处理下拉菜单的UI逻辑，这需要实时响应
      this.handleDropdownLogic(e)

      // 2. 将数据更新操作放入防抖函数中
      // 清除上一个计时器
      if (this.timer) {
        clearTimeout(this.timer)
      }
      // 设置一个新的计时器，在用户停止输入200ms后执行 updateValue
      this.timer = setTimeout(() => {
        this.updateValue()
      }, 350)
    },
    handleBlur() {
      // 当输入框失焦时，为了确保数据最终被同步，立即执行一次 updateValue
      // 同时清除可能存在的防抖计时器，避免重复执行
      if (this.timer) {
        clearTimeout(this.timer)
      }
      this.updateValue()
      if (this.semanticTags && !this.showDropdown) this.refreshSemanticLabels()
    },
    async handleDropdownLogic(e) {
      // 获取光标前的最后一个字符。此操作依赖于当前的DOM状态，应在任何DOM修改前执行。
      const lastChar = this.getLastChar()

      // --- 场景1: 用户执行删除操作 ---
      if (e.inputType === 'deleteContentBackward' && !this.atText) {
        if (lastChar != '/') {
          this.showDropdown = false
        }
      } else {
        // --- 场景2: 用户输入'/'，触发提及功能 ---
        if (lastChar === '/') {
          this.showDropdown = true
          this.atText = ''
        } else {
          // --- 场景3: 用户在'/'后输入文本，进行筛选 ---
          const text = this.getLastAtText()
          this.atText = text
          const isExistStartWith = this.checkStartWith(text)
          if (text) {
            this.showDropdown = isExistStartWith
          }
        }
      }

      // 如果不显示下拉菜单，则提前退出定位逻辑
      if (!this.showDropdown) {
        return
      }

      if (this.semanticTags) {
        const selection = window.getSelection()
        if (!selection.rangeCount || !this.$refs.JMention?.contains(selection.anchorNode)) return
        this.dropdownSelection = getPromptSelection(this.$refs.JMention)
        this.dropdownAnchorRange = selection.getRangeAt(0).cloneRange()
        await this.$nextTick()
        this.positionSemanticDropdown()
        return
      }

      // --- 定位并最终显示下拉菜单 ---
      const { top, left } = await this.getCaretPosition()
      this.showDropdown = true
      this.positionStyle = `left:${left}px; top:${top}px;`
    },
    insertAtCaret(text, dataSet) {
      if (this.semanticTags) {
        const editor = this.$refs.JMention
        const value = readPromptText(editor)
        const offsets = this.dropdownSelection || getPromptSelection(editor)
        if (!offsets) return
        const trigger = '/' + this.atText
        if (value.slice(0, offsets.start).endsWith(trigger)) offsets.start -= trigger.length
        // 浮层搜索框获得焦点后，仍在原编辑位置插入，并保持编码用于保存。
        renderPromptText(editor, value, [...this.options, ...this.cardOptions])
        restorePromptSelection(editor, offsets)
        const selection = window.getSelection()
        const range = selection.getRangeAt(0)
        const caretSlot = range.startContainer.parentElement
        if (caretSlot?.dataset?.promptCaret) range.setStartAfter(caretSlot)
        range.deleteContents()
        const token = createPromptToken({ ...dataSet, label: text }, document)
        range.insertNode(token)
        ensurePromptCaretSlots(editor)
        editor.focus()
        range.setStart(token.nextSibling.firstChild, 1)
        range.collapse(true)
        selection.removeAllRanges()
        selection.addRange(range)
        this.dropdownSelection = null
        return
      }
      let len = this.atText.length + 1
      const lastChar = this.getLastChar(len)
      if (lastChar !== '/' + this.atText) len = 0
      const selection = window.getSelection()
      if (selection.rangeCount <= 0) return
      const range = selection.getRangeAt(0)
      const startOffset = range.startOffset
      const startContainer = range.startContainer

      let prevCharOffset = startOffset - len
      let prevCharNode = startContainer
      if (prevCharOffset < 0) {
        // 如果当前节点没有前一个字符，查找前一个兄弟节点
        let prevSibling = startContainer.previousSibling
        while (prevSibling) {
          if (prevSibling.nodeType === Node.TEXT_NODE) {
            prevCharOffset = prevSibling.textContent.length - len
            prevCharNode = prevSibling
            break
          }
          prevSibling = prevSibling.previousSibling
        }
      }

      if (prevCharNode.nodeType === Node.TEXT_NODE) {
        // 将触发词替换为空字符串，而不是空格，以避免在span后留下多余的空格
        prevCharNode.replaceData(prevCharOffset, len, '')

        // 创建 span 标签
        const span = document.createElement('span')
        span.contentEditable = 'false'
        // 根据类型区分样式：小程序卡片用蓝色 Tag，普通变量用白底 Tag
        if (dataSet.type === 'card') {
          span.classList.add('j-mention-at', 'mini-card-tag')
          // 插入小程序图标 + 文字
          const iconWrap = document.createElement('span')
          iconWrap.className = 'svg-action'
          iconWrap.style.cssText = 'font-size: 12px;'
          iconWrap.innerHTML = '<svg aria-hidden="true" class="svg-icon" fill="currentColor" width="1em" height="1em"><use xlink:href="#applet" fill="currentColor"></use></svg>'
          span.appendChild(iconWrap)
          span.appendChild(document.createTextNode(text))
        } else {
          span.classList.add('j-mention-at')
          span.style = this.atTextStyle
          span.textContent = `${text}`
        }

        // const tmp = document.createElement("span");
        // tmp.contentEditable = true;
        // tmp.innerHTML += "&nbsp;";
        // tmp.classList.add("at-space");

        for (const key in dataSet) {
          span.dataset[key] = dataSet[key]
        }

        // 在原字符位置插入 span 标签
        const newRange = document.createRange()
        newRange.setStart(prevCharNode, prevCharOffset)
        newRange.insertNode(span)
        // newRange.insertNode(tmp);

        // 将光标移动到 span 元素外部
        // 将光标移动到 span 元素之后
        range.setStartAfter(span)
        range.setEndAfter(span)

        selection.removeAllRanges()
        selection.addRange(range)
      }
    },
    onSelectOption(value, selectedValuePath, selectedPath) {
      // 兼容 SlashPopupPanel 直接传 opt 的调用方式
      if (!selectedPath) {
        this.selectOption(value)
        return
      }
      let item = selectedPath[selectedPath.length - 1]
      delete item['children']
      this.selectOption(item)
    },
    selectOption(opt, isInit = false) {
      if (!opt) return
      // 处理用户选择逻辑
      this.showDropdown = false

      // let text = opt.node_name + '.' + opt.text;

      let text = opt.label

      const dataSet = {
        id: opt.id,
        label: opt.label,
        value: opt.value,
        text: opt.text,
        index: this.selectedList.length
      }

      this.selectedList.push(opt)
      this.selectedIdSet.add(opt.value)

      if (!isInit) {
        this.insertAtCaret(text, dataSet)
        this.initShowOptionList()
        // 在插入提及内容后，立即更新值以确保数据同步
        this.updateValue()
      }
    },
    // 选择小程序卡片，插入蓝色 Tag
    selectCard(opt, isInit = false) {
      if (!opt) return
      this.showDropdown = false

      let text = opt.label

      const dataSet = {
        id: opt.id,
        label: opt.label,
        value: opt.value,
        text: opt.text,
        index: this.selectedList.length,
        type: 'card'
      }

      this.selectedList.push(opt)
      this.selectedIdSet.add(opt.value)

      if (!isInit) {
        this.insertAtCaret(text, dataSet)
        this.initShowOptionList()
        this.updateValue()
      }
    },

    handleDirectionChange(direction, width) {
      const leftMatch = this.positionStyle.match(/left:\s*(-?\d+\.?\d*)px/)
      const topMatch = this.positionStyle.match(/top:\s*(-?\d+\.?\d*)px/)

      if (!leftMatch || !topMatch) return

      const currentLeft = parseFloat(leftMatch[1])
      const currentTop = parseFloat(topMatch[1])
      let newLeft

      if (direction === 'left') {
        newLeft = currentLeft - width
      } else {
        newLeft = currentLeft + width
      }

      this.positionStyle = `left:${newLeft}px; top:${currentTop}px;`
    },
    handleKeydown(event) {
      if (this.semanticTags && ['Backspace', 'Delete'].includes(event.key)) {
        const editor = this.$refs.JMention
        const selection = window.getSelection()
        const slot = selection.anchorNode?.parentElement
        if (selection.isCollapsed && slot?.dataset?.promptCaret && selection.anchorOffset <= 1) {
          const token = event.key === 'Backspace' ? slot.previousSibling : slot.nextSibling
          if (token?.classList?.contains('j-mention-at')) {
            event.preventDefault()
            const offsets = getPromptSelection(editor)
            const position = offsets.start - (event.key === 'Backspace' ? token.dataset.value.length : 0)
            token.remove()
            this.renderSemanticValue(readPromptText(editor))
            restorePromptSelection(editor, { start: position, end: position })
            this.updateValue()
            return
          }
        }
      }
      if (this.semanticTags && this.showDropdown && event.key === 'Escape') {
        event.preventDefault()
        event.stopPropagation()
        this.hideDropdownMenus()
        return
      }
      // 当type为input时，阻止回车键换行
      if (this.type === 'input' && event.key === 'Enter') {
        event.preventDefault();
        return false;
      }
    },
  }
}
</script>
<style lang="less">
.mention-input-warpper.mention-promt-input .j-mention-at {
  border-radius: 6px;
  padding: 1px 8px;
  border: 1px solid #00000026;
  background: #fff;
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
  margin-right: 6px;
}
.mention-input-warpper.mention-promt-input .j-mention-at.mini-card-tag {
  border: none;
  background: #E5EFFF;
  color: #164799;
  font-size: 12px;
  line-height: 20px;
  padding: 1px 8px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.mention-input-warpper.mention-promt-input .j-mention-at.mini-card-tag .svg-action {
  display: inline-flex;
  align-items: center;
  color: inherit;
}
.mention-input-warpper.mention-promt-input .j-mention-at.mini-card-tag .svg-icon {
  display: inline-block;
  fill: currentColor;
}
</style>
<style lang="less" scoped>
.mention-input-warpper {
  position: relative;
  overflow: hidden;
  border-radius: 6px;
  padding: 4px 8px;
  border: 1px solid #d9d9d9;
  background: #fff;
  &.no-border{
    border: none;
    padding: 0;
    .placeholder{
      padding: 0;
    }
    &:hover {
      border: none;
    }
    &:focus {
      border: none;
    }
  }

  &:hover {
    border: 1px solid #409eff;
  }

  &:focus {
    border-color: #4d94ff;
    box-shadow: 0 0 0 2px rgba(5, 138, 255, 0.06);
    border-inline-end-width: 1px;
    outline: 0;
  }
}
.j-mention {
  position: relative;
  overflow: scroll;
  width: 100%;
  line-height: 24px;
  min-height: 24px;
  outline: none;
  background: #fff;

  &.type-input {
    overflow-y: hidden;
    overflow-x: scroll;
    white-space: nowrap;
  }

  &.type-textarea {
    white-space: pre-wrap;
    word-wrap: break-word;
    min-height: 64px;
  }

  &::-webkit-scrollbar {
    width: 0px;
    height: 0px;
  }

  /*定义滚动条轨道 内阴影+圆角*/
  &::-webkit-scrollbar-track {
    border-radius: 0px;
    background-color: #fafafa;
  }

  /*定义滑块 内阴影+圆角*/
  &::-webkit-scrollbar-thumb {
    border-radius: 0px;
    background: rgb(191, 191, 191);
  }
}
.mention-input-warpper.semantic-mention {
  .j-mention { padding-top: 2px; padding-bottom: 2px; }
  :deep(.j-mention-at) { display: inline-flex; align-items: center; min-height: 24px; box-sizing: border-box; padding: 0 8px; line-height: 22px; vertical-align: baseline; cursor: text; }
  :deep(.j-mention-at.mini-card-tag) { min-height: 22px; line-height: 20px; }
}
.placeholder {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  line-height: 24px;
  overflow: hidden;
  padding: 4px 8px;
  color: #bfbfbf;
  pointer-events: none;
  user-select: none;
}
.dropdown-list {
  padding: 0;
  margin: 0;
  background-color: #fff;
  z-index: 999999;
  min-width: 180px;
  overflow: hidden;
  user-select: none;
  border-radius: 8px;
  box-shadow:
    0 6px 16px 0 rgba(0, 0, 0, 0.08),
    0 3px 6px -4px rgba(0, 0, 0, 0.12),
    0 9px 28px 8px rgba(0, 0, 0, 0.05);

  .dropdown-list-nodata {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    font-size: 13px;
    padding: 8px;
  }

  .dropdown-list-item {
    padding: 8px;
    cursor: pointer;
    background: #fff;

    &:hover {
      background-color: rgba(0, 0, 0, 0.04) !important;
    }

    &.active {
      background-color: rgba(0, 0, 0, 0.09) !important;
    }
  }

  &::-webkit-scrollbar {
    width: 0px;
    height: 0px;
  }

  /*定义滚动条轨道 内阴影+圆角*/
  &::-webkit-scrollbar-track {
    border-radius: 0px;
    background-color: #fafafa;
  }

  /*定义滑块 内阴影+圆角*/
  &::-webkit-scrollbar-thumb {
    border-radius: 0px;
    background: rgb(191, 191, 191);
  }
}

.dropdown-list-fixed {
  position: fixed;
}
.dropdown-list-semantic {
  overflow-y: auto;
  overscroll-behavior: contain;
}

.dropdown-list-absolute {
  position: absolute;
}
</style>
