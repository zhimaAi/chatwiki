<template>
  <div style="width: 100%;position: relative;">
    <div></div>
    <Teleport to="body">
      <div class="dropdown-list" ref="dropdownList" v-if="showDropdown" :style="positionStyle">
        <CascadePanel
          :options="showOptions"
          @change="onSelectOption"
          @direction-change="handleDirectionChange"
        />
      </div>
    </Teleport>
    <div class="mention-input-warpper" ref="JMentionContainer" >
      <div
      ref="JMention"
      :style="[inputStyle]"
      class="j-mention"
      :class="[{'show-placeholder': localValue.length == 0, }, 'type-'+type]"
      contenteditable="plaintext-only"
      :placeholder="placeholder"
      @wheel.stop=""
      @focus="onFocus"
      @input="divInput"
      @keydown="dropDownKeydown"
      @click="clickJMention"></div>
        <span class="placeholder" v-if="localValue.length == 0">{{ placeholder }}</span>
    </div>
  </div>
</template>

<script>
import CascadePanel from './cascade-panel.vue'

// 把树状的options抓换成编排的数组
const getTreeOptions = (options) => {
  let result = [];
  options.forEach((opt) => {
    if (opt.children && opt.children.length > 0) {
      result.push(...getTreeOptions(opt.children));
    } else {
      result.push(opt);
    }
  });
  return result;
}

export default {
  name: "JMention",
  components: {
    CascadePanel,
  },
  inject: [],
  emits: ["input", 'change', "update:selectedList", "focus", 'open'],
  props: {
    type: {
      type: String,
      default: "input", // input, textarea
      validator(value) {
        return ["input", "textarea"].includes(value);
      },
    },
    placeholder: {
      type: String,
      default: "请输入",
    },
    defaultValue: {
      type: [String, Number],
      default: "",
    },
    options: {
      type: Array,
      default: () => [],
    },
    defaultSelectedList: {
      type: Array,
      default: () => [],
    },
    canRepeat: {
      type: Boolean,
      default: true,
    },
    isContain: {
      type: Boolean,
      default: false,
    },
    inputStyle: {
      type: String,
      default: "",
    }
  },
  data() {
    return {
      timer: null,
      showDropdown: false,
      showAtList: false,
      positionStyle: "",
      chooseIndex: 0,
      selectedList: [],
      selectedIdSet: new Set(),
      showOptions: [],
      atText: "",
      localValue: "",
    };
  },
  watch: {
    defaultValue(val){
      if(val != this.localValue){
        this.refresh()
      }
    },
    chooseIndex() {},
    showAtList(val){
      this.onOpen(val);
    },
    options(){
      this.initShowOptionList();
    },
  },
  mounted() {
    this.initShowOptionList();

    document.addEventListener("dragstart", this.handleDragstart);
    document.addEventListener("click", this.handleClick);
    document.addEventListener("scroll", this.onScroll);

    this.initData();
  },
  beforeUnmount() {
    document.removeEventListener("dragstart", this.handleDragstart);
    document.removeEventListener("click", this.handleClick);
    document.removeEventListener("scroll", this.onScroll);
  },
  methods: {
    handleClick(e){
      const target = e.target.closest(".dropdown-list");
      if (!target) {
        this.hideDropdownMenus()
      }
    },
    handleDragstart(e){
      const target = e.target.closest(".dropdown-list");
      if (target) {
        e.preventDefault();
      }
    },
    onScroll(){
      this.hideDropdownMenus()
    },
    hideDropdownMenus() {
      this.showDropdown = false;
      this.showAtList = false;
    },
    onFocus(){
      this.$emit("focus");
    },
    onOpen(val){
      this.$emit("open", val);
    },
    getText(htmlStr){
      const div = document.createElement('div');
      div.innerHTML = htmlStr;
      let text = '';
        // 遍历子节点
        Array.from(div.childNodes).forEach(node => {
        let textContent = '';
        if(this.type === 'textarea'){
          // textarea
          if(node.nodeName === 'BR'){ // 如果是BR节点,直接添加换行
            textContent = '\n';
          }else{
            if (node.classList && node.classList.contains('j-mention-at')) {
              textContent = node.dataset.value;
            } else{
              textContent = node.textContent
            }
          }
        }else{
          // input
          if (node.classList && node.classList.contains('j-mention-at')) {
            textContent = node.dataset.value;
          }else if(node.nodeName === 'BR'){
              textContent = '';
          }else{
            textContent = node.textContent.replace(/\n/g, '');
          }
        }

        // 去除只有&nbsp;字符的节点
        if (textContent === '&nbsp;' || textContent === ' ') {
          textContent = '';
        }

        // 拼接文本内容
        text += textContent;
      });

      return text;
    },
    onChange(){
      const jMention = this.$refs.JMention;
      let htmlStr = jMention.innerHTML;
      let text = this.getText(htmlStr);

      this.localValue = text

      let selectedListMap = {};

      this.selectedList.forEach(item => {
        if(!selectedListMap[item.value] && text.indexOf(item.value) !== -1){
          selectedListMap[item.value] = item;
        }
      });

      let mySelectedList = Object.values(selectedListMap);

      this.$emit("change", text.trim(), mySelectedList);
    },
    refresh(){
      this.selectedList = [];
      this.selectedIdSet = new Set();

      this.initShowOptionList();
      this.initData();
    },
    initData() {
      const JMention = this.$refs.JMention;
      // let html = this.defaultValue;
      let html = this.defaultValue.toString()
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
      let defaultSelectedList = this.defaultSelectedList || [];
      let treeOptions = getTreeOptions(this.options);
      let selectedList = [...defaultSelectedList, ...treeOptions].filter(item => item && item.value && item.value != '' && item.typ !== 'node');

      // 记录已经替换过的的value
      const replacedValues = new Set();

      selectedList.forEach((opt) => {
        // 使用正则匹配完整的value,避免部分匹配
        const regex = new RegExp(`(${opt.value})`, 'g');
        if(regex.test(this.defaultValue) && !replacedValues.has(opt.value)){
          let text = opt.node_name + '.' + opt.text;

          text = text.replace(/\./g, '/')

          html = html.replace(regex,
            `<span class="j-mention-at" data-id="${opt.node_id}" contentEditable="false" data-value="${opt.value}">${text}</span> `
          );

          replacedValues.add(opt.value);

          this.selectOption(opt, true);
        }
      });

      JMention.innerHTML = html;
      this.localValue = this.defaultValue;
    },
    initShowOptionList() {
      this.showOptions = this.options.filter((opt) => {
        if (this.atText) {
          return opt.label.startsWith(this.atText);
        }
        if (this.canRepeat) return true;
        return !this.selectedIdSet.has(opt.id + "");
      });
    },
    dropDownKeydown(event) {
      const { keyCode } = event;
      if (!this.showDropdown) {
        // 阻止enter键的默认行为
        if (keyCode === 13 && this.type === 'input') {
          event.preventDefault();
          event.stopPropagation();
          return;
        }
        return
      };

      const keyCodeList = [13, 38, 40];
      if (!keyCodeList.includes(keyCode)) return;
      event.preventDefault();
      event.stopPropagation();

      if (keyCode === 13) {
        return;
      }

      const map = {
        38: -1,
        40: 1,
      };
      this.chooseIndex += map[keyCode] || 0;
      this.chooseIndex = Math.max(
        0,
        Math.min(this.chooseIndex, this.showOptions.length - 1)
      );
    },
    getLastChar(len = 1) {
      const selection = window.getSelection();
      if (selection.rangeCount <= 0) return;
      const range = selection.getRangeAt(0);
      const startOffset = range.startOffset;
      const startContainer = range.startContainer;
      let textBeforeCursor = "";
      if (startContainer.nodeType === Node.TEXT_NODE) {
        textBeforeCursor = startContainer.textContent.slice(
          0,
          startOffset
        );
      } else if (startContainer.nodeType === Node.ELEMENT_NODE) {
        // 若光标在元素节点内，获取其内部文本
        const childNodes = startContainer.childNodes;
        for (let i = 0; i < childNodes.length; i++) {
          const node = childNodes[i];
          if (node.nodeType === Node.TEXT_NODE) {
            if (i === childNodes.length - 1) {
              textBeforeCursor += node.textContent.slice(
                0,
                startOffset
              );
            } else {
              textBeforeCursor += node.textContent;
            }
          }
        }
      }
      return textBeforeCursor.slice(-len);
    },
    getLastAtText() {
      const selection = window.getSelection();
      if (selection.rangeCount <= 0) return;
      const range = selection.getRangeAt(0);
      const startOffset = range.startOffset;
      const startContainer = range.startContainer;
      let textBeforeCursor = "";
      if (startContainer.nodeType === Node.TEXT_NODE) {
        textBeforeCursor = startContainer.textContent.slice(
          0,
          startOffset
        );
      } else if (startContainer.nodeType === Node.ELEMENT_NODE) {
        // 若光标在元素节点内，获取其内部文本
        const childNodes = startContainer.childNodes;
        for (let i = 0; i < childNodes.length; i++) {
          const node = childNodes[i];
          if (node.nodeType === Node.TEXT_NODE) {
            if (i === childNodes.length - 1) {
              textBeforeCursor += node.textContent.slice(
                0,
                startOffset
              );
            } else {
              textBeforeCursor += node.textContent;
            }
          }
        }
      }
      const lastAtIndex = textBeforeCursor.lastIndexOf("/");
      if (lastAtIndex === -1) return "";
      return textBeforeCursor.slice(lastAtIndex + 1);
    },
    async getCaretPosition() {
      // const JMention = this.$refs.JMention;
      const selection = window.getSelection();
      if (selection.rangeCount <= 0) return { top: 0, left: 0 };
      const range = selection.getRangeAt(0);
      const rect = range.getBoundingClientRect();
      const dropdownListRect = await this.getDropdownListSize();
      const viewportWidth = window.innerWidth;
      const viewportHeight = window.innerHeight;

      let top = rect.bottom;
      let left = rect.left;

      // 自动垂直定位：如果下方空间不足，则移动到上方
      if (top + dropdownListRect.height > viewportHeight) {
        top = rect.top - dropdownListRect.height;
      }

      // 自动水平定位：如果右侧空间不足，则向左对齐
      if (left + dropdownListRect.width > viewportWidth) {
        left = rect.right - dropdownListRect.width;
      }

      // 确保不会超出视口顶部和左部
      if (top < 0) {
        top = 0;
      }
      if (left < 0) {
        left = 0;
      }

      return {
        top,
        left,
      };
    },
    async getDropdownListSize() {
      this.showDropdown = true; // 临时显示元素
      await this.$nextTick(); // 等待 DOM 更新
      const width = this.$refs.dropdownList.offsetWidth;
      const height = this.$refs.dropdownList.offsetHeight;
      this.showDropdown = false; // 恢复隐藏
      return { width, height };
    },
    updateSelectedList() {
      const JMention = this.$refs.JMention;
      const nodeList = JMention.querySelectorAll(".j-mention-at");
      this.selectedList = [];
      this.selectedIdSet = new Set();

      for (const node of nodeList) {
        // const { value } = node.dataset;
        let value = node.dataset.value;
        const item = this.options.find((opt) => {
          return opt.value == value;
        });

        this.selectedList.push(item);
        this.selectedIdSet.add(value);
      }
      this.$emit("update:selectedList", this.selectedList);
      if (
        nodeList.length === this.selectedList.length &&
        !this.atText &&
        !this.showDropdown
      ){
        return;
      }

      this.initShowOptionList();
    },
    checkStartWith(text) {
      const list = this.options.filter((opt) => {
        return opt.label.startsWith(text);
      });
      return list.length > 0;
    },
    updateValue() {
      if(this.timer){
        clearTimeout(this.timer);
        this.timer = null;
      }

      this.timer = setTimeout(() => {
        this.onChange();
      }, 100);
    },
    clickJMention() {
      this.showDropdown = false;
      this.showAtList = false;
    },
    async divInput(e) {
      this.updateValue();
      const lastChar = this.getLastChar();
      if (e.inputType === "deleteContentBackward" && !this.atText) {
        this.showDropdown = false;
        this.showAtList = false;
      }
      if (lastChar === "/") {
        this.showDropdown = true;
        this.atText = "";
      } else {
        const text = this.getLastAtText();
        this.atText = text;
        const isExistStartWith = this.checkStartWith(text);
        if (text) {
          this.showDropdown = isExistStartWith;
        }
      }
      // this.updateSelectedList();
      if (!this.showDropdown) {
        return;
      }
      const { top, left } = await this.getCaretPosition();
      this.showDropdown = true;
      this.positionStyle = `left:${left}px; top:${top}px;`;
      this.chooseIndex = 0;
      this.showAtList = true;
    },
    insertAtCaret(text, dataSet) {
      let len = this.atText.length + 1;
      const lastChar = this.getLastChar(len);
      if (lastChar !== "/" + this.atText) len = 0;
      const selection = window.getSelection();
      if (selection.rangeCount <= 0) return;
      const range = selection.getRangeAt(0);
      const startOffset = range.startOffset;
      const startContainer = range.startContainer;

      let prevCharOffset = startOffset - len;
      let prevCharNode = startContainer;
      if (prevCharOffset < 0) {
        // 如果当前节点没有前一个字符，查找前一个兄弟节点
        let prevSibling = startContainer.previousSibling;
        while (prevSibling) {
          if (prevSibling.nodeType === Node.TEXT_NODE) {
            prevCharOffset = prevSibling.textContent.length - len;
            prevCharNode = prevSibling;
            break;
          }
          prevSibling = prevSibling.previousSibling;
        }
      }

      if (prevCharNode.nodeType === Node.TEXT_NODE) {
        prevCharNode.replaceData(prevCharOffset, len, " ");

        // 创建 span 标签
        const span = document.createElement("span");
        span.contentEditable = "false";
        span.classList.add("j-mention-at");
        span.style = this.atTextStyle;
        span.textContent = `${text}`;

        const tmp = document.createElement("span");
        tmp.contentEditable = true;
        tmp.innerHTML += "&nbsp;";
        tmp.classList.add("at-space");

        for (const key in dataSet) {
          span.dataset[key] = dataSet[key];
        }

        // 在原字符位置插入 span 标签
        const newRange = document.createRange();
        newRange.setStart(prevCharNode, prevCharOffset);
        newRange.insertNode(span);
        // newRange.insertNode(tmp);

        // 将光标移动到 span 元素外部
        let nextNode = span.nextSibling;
        if (!nextNode) {
          // 如果 span 元素后面没有兄弟节点，创建一个新的文本节点
          nextNode = document.createTextNode("");
          tmp.parentNode.appendChild(nextNode);
        }
        range.setStart(nextNode, 1);
        range.setEnd(nextNode, 1);
        selection.removeAllRanges();
        selection.addRange(range);
      }
    },
    onSelectOption(value, selectedValuePath, selectedPath) {
      let item = selectedPath[selectedPath.length - 1];
      delete item['children'];
      this.selectOption(item)
    },
    selectOption(opt, isInit = false) {
      if (!opt) return;
      // 处理用户选择逻辑
      this.showDropdown = false;
      this.showAtList = false;

      let text = opt.node_name + '.' + opt.text;

      text = text.replace(/\./g, '/')

      const dataSet = {
        id: opt.id ,
        label: opt.label,
        value: opt.value,
        text: opt.text,
        index: this.selectedList.length,
      };

      this.selectedList.push(opt);
      this.selectedIdSet.add(opt.value);

      if (!isInit) {
        this.insertAtCaret(text, dataSet);
        this.initShowOptionList();
        this.updateValue();
      }
    },

    handleDirectionChange(direction, width) {
      const leftMatch = this.positionStyle.match(/left:\s*(-?\d+\.?\d*)px/);
      const topMatch = this.positionStyle.match(/top:\s*(-?\d+\.?\d*)px/);

      if (!leftMatch || !topMatch) return;

      const currentLeft = parseFloat(leftMatch[1]);
      const currentTop = parseFloat(topMatch[1]);
      let newLeft;

      if (direction === 'left') {
        newLeft = currentLeft - width;
      } else {
        newLeft = currentLeft + width;
      }

      this.positionStyle = `left:${newLeft}px; top:${currentTop}px;`;
    },
  },
};
</script>
<style lang="less">
.mention-input-warpper .j-mention-at {
  border-radius: 4px;
  padding: 1px 4px;
  background: #f2f4f5;
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
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

  &.type-input{
    overflow-y: hidden;
    overflow-x: scroll;
    white-space: nowrap;
  }

  &.type-textarea{
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
.placeholder{
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
    white-space: nowrap;
  }
.dropdown-list {
  position: fixed;
  padding: 0;
  margin: 0;
  background-color: #fff;
  z-index: 999999;
  min-width: 180px;
  overflow: hidden;
  user-select: none;
  border-radius: 8px;
  box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.08), 0 3px 6px -4px rgba(0, 0, 0, 0.12), 0 9px 28px 8px rgba(0, 0, 0, 0.05);

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
</style>
