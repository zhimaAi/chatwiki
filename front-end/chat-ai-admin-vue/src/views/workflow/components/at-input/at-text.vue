<template>
  <span ref="JMention"  class="j-mention"  :class="['type-'+type]" ></span>
</template>

<script>
import { useI18n } from '@/hooks/web/useI18n'

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
  inject: [],
  emits: ["input", 'change', "update:selectedList", "focus", 'open', 'resize'],
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
      default: "ph_input",
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
  setup() {
    const { t } = useI18n('views.workflow.components.at-input.at-text')
    return {
      t
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
  computed: {
    computedPlaceholder() {
      return this.placeholder === 'ph_input' ? this.t('ph_input') : this.placeholder
    }
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
    this.initData();
  },
  beforeUnmount() {
  },
  methods: {
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

      this.localValue = text;

      let selectedListMap = {};

      this.selectedList.forEach(item => {
        if(!selectedListMap[item.value] && text.indexOf(item.value) !== -1){
          selectedListMap[item.value] = item;
        }
      });

      let mySelectedList = Object.values(selectedListMap);

      this.$emit("change", text, mySelectedList);
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
            `<span class="j-mention-at" data-id="${opt.node_id}" contentEditable="false" data-value="${opt.value}">${text}</span>`
          );

          replacedValues.add(opt.value);

          this.selectOption(opt, true);
        }
      });

      JMention.innerHTML = html;
      this.localValue = this.defaultValue;

      this.$nextTick(() => {
        this.$emit('resize');
      })
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

      // 获取视口位置
      let top = rect.bottom;
      let left = rect.left;

      // 确保下拉列表不会超出视口
      const viewportWidth = window.innerWidth;
      const viewportHeight = window.innerHeight;

      // 检查右边界
      if (left + dropdownListRect.width > viewportWidth) {
        left = viewportWidth - dropdownListRect.width;
      }

      // 检查底部边界
      if (top + dropdownListRect.height > viewportHeight) {
        // 如果下方空间不足，将下拉列表显示在输入框上方
        top = rect.top - dropdownListRect.height;
      }

      return {
        top,
        left,
        x: left,
        y: top
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
      this.onChange();
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

        // const tmp = document.createElement("span");
        // tmp.contentEditable = true;
        // tmp.innerHTML += "&nbsp;";
        // tmp.classList.add("at-space");

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
          // nextNode = document.createTextNode("");
          // tmp.parentNode.appendChild(nextNode);
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
  },
};
</script>
<style lang="less" scoped>

</style>
