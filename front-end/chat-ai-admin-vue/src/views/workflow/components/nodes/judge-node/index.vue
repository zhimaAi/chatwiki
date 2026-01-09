<style lang="less" scoped>
@import '../form-block.less';

.judge-node {
  position: relative;

  .condition-box {
    position: relative;
    display: flex;
    align-items: center;
    .condition-left-box {
      margin-right: 4px;
    }
    .connection-text {
      display: inline-block;
      height: 18px;
      line-height: 18px;
      padding: 0 8px;
      font-size: 12px;
      font-weight: 400;
      border-radius: 4px;
      color: #595959;
      background: #e4e6eb;
    }
    .condition-line {
      height: 100%;
      width: 24px;
      margin-right: 4px;
      &::after {
        content: '';
        position: absolute;
        width: 24px;
        top: 12px;
        bottom: 12px;
        border: 1px solid #d9d9d9;
        border-right: 0;
        border-radius: 2px;
        border-top-right-radius: 0;
        border-bottom-right-radius: 0;
      }
    }

    .condition-body{
      flex: 1;
      overflow: hidden;
    }

    .field-items {
      .field-item {
        display: flex;
        align-items: center;
        min-height: 24px;
        line-height: 16px;
        padding: 2px 4px;
        margin-bottom: 2px;
        border-radius: 4px;
        border: 1px solid #d9d9d9;
        background: #fff;

        &:last-child {
          margin-bottom: 0;
        }

        .field-name,
        .field-value {
          flex: 1;
          line-height: 16px;
          font-size: 12px;
          font-weight: 400;
          font-size: 12px;
          color: #595959;
          overflow: hidden;
          white-space: nowrap;
          text-overflow: ellipsis;
        }
        .field-rule {
          height: 18px;
          line-height: 18px;
          padding: 0 8px;
          margin: 0 4px;
          font-size: 12px;
          font-weight: 400;
          border-radius: 4px;
          color: #595959;
          background: #e4e6eb;
        }
      }
    }
  }

  .if-box {
    position: relative;
    display: flex;
    align-items: center;
    margin-bottom: 8px;

    .if-box-label {
      width: 60px;
      line-height: 22px;
      margin-right: 8px;
      text-align: right;
      color: #262626;
    }

    .if-box-content {
      flex: 1;
      padding: 8px;
      border-radius: 6px;
      background: #f2f4f7;
    }
  }

  .else-box {
    display: flex;
    align-items: center;
    margin-bottom: 8px;

    .else-box-label {
      width: 60px;
      line-height: 22px;
      margin-right: 8px;
      text-align: right;
      color: #262626;
    }

    .else-box-content {
      line-height: 16px;
      padding: 3px 4px;
      font-size: 12px;
      border-radius: 4px;
      border: 1px solid #d9d9d9;
      background: #fff;
    }
  }
}
</style>

<template>
  <node-common
    :properties="properties"
    :title="properties.node_name"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    style="width: 420px;"
  >
    <div class="judge-node">
      <div class="if-box" :ref="(el) => ifBoxRefs[index] = el" v-for="(item, index) in formState.term" :key="index">
        <div class="if-box-label">{{ index == 0 ? 'if' : 'else if' }}</div>
        <div class="if-box-content">
          <div class="condition-box">
            <div class="condition-left-box" v-if="item.terms.length > 1">
              <span class="connection-text">{{ item.is_or == 1 ? 'or' : 'and' }}</span>
            </div>
            <div class="condition-line" v-if="item.terms.length > 1"></div>
            <div class="condition-body">
              <div class="field-items">
                <div class="field-item" v-for="(subItem, index) in item.terms" :key="index">
                  <span class="field-name">
                    <user-question-text :value="subItem.variable" />
                  </span>
                  <span class="field-rule">{{ getRuleLabel(subItem) }}</span>
                  <span class="field-value">
                    {{ subItem.value }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="else-box">
        <div class="else-box-label">else</div>
        <div class="else-box-content">
          <div class="else-desc">不符合上述所有分支的条件时，走默认分支</div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { ref, reactive, watch, inject, onMounted, nextTick } from 'vue'
import NodeCommon from '../base-node.vue'
import userQuestionText  from '../user-question-text.vue'

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

const resetSize = inject('resetSize')
const getNode = inject('getNode')
const updateAnchorList = inject('updateAnchorList')

// 特殊节点列表
let specialNodeList = [
  'special.lib_paragraph_list',
  'special.llm_reply_content',
  'specify-reply-node'
]

const formState = reactive({
  term: []
})

const variableOptions = ref([])

// 用于收集所有if-box的ref
const ifBoxRefs = ref([])

function getOptions() {
  let list = getNode().getAllParentVariable()

  variableOptions.value = handleOptions(list)
}

// 递归处理Options
function handleOptions(options) {
  options.forEach((item) => {
    if (item.typ == 'node') {
      if (item.node_type == 1) {
        item.value = 'global'
      } else {
        item.value = item.node_id
      }
    } else {
      item.value = item.key
    }

    if (item.children && item.children.length > 0) {
      item.children = handleOptions(item.children)
    }
  })

  return options
}

let baseTypeOptions = [
  {
    label: '等于',
    value: 1
  },
  {
    label: '不等于',
    value: 2
  },
  {
    label: '包含',
    value: 3
  },
  {
    label: '不包含',
    value: 4
  },
  {
    label: '为空',
    value: 5
  },
  {
    label: '不为空',
    value: 6
  }
]

let baseTypeOptions2 = [
  {
    label: '包含其中一项',
    value: 3
  },
  {
    label: '不包含其中一项',
    value: 4
  },
  {
    label: '为空',
    value: 5
  },
  {
    label: '不为空',
    value: 6
  }
]

function getTypeByVariable(data) {
  let typ = ''
  if (data.variable && data.variable.length > 0) {
    let slectItem = variableOptions.value.filter((item) => item.key == data.variable[0])
    if (slectItem && slectItem.length) {
      typ = slectItem[0].typ
    }
    if (typ == 'object') {
      if (slectItem[0] && slectItem[0].children) {
        slectItem = slectItem[0].subs.filter((item) => item.key == data.variable[1])
      }
      if (slectItem && slectItem.length) {
        typ = slectItem[0].typ
      }

      if (typ == 'object') {
        if (slectItem[0] && slectItem[0].children) {
          slectItem = slectItem[0].subs.filter((item) => item.key == data.variable[2])
        }
        if (slectItem && slectItem.length) {
          typ = slectItem[0].typ
        }
      }
    }
  }
  return typ
}

function getTypeOptions(data) {
  if (!data.variable) {
    return []
  }
  let typ = getTypeByVariable(data)

  if (typ == '') {
    if (data.is_mult) {
      return baseTypeOptions2
    } else {
      return baseTypeOptions
    }
  }

  if (typ == 'boole') {
    return [
      {
        label: '等于',
        value: 1
      },
      {
        label: '不等于',
        value: 2
      }
    ]
  }

  if (typ.includes('array')) {
    return baseTypeOptions2
  }

  return baseTypeOptions
}

const getRuleLabel =  (item) => {
  let options = getTypeOptions(item)
  let index = options.findIndex((it) => it.value == item.type)

  return options[index].label
}

const init = () => {
  try {
    let term = JSON.parse(props.properties.node_params).term || []

    term = term.map((item) => {
      let terms = item.terms.map((it) => {
        // 判断是不是特殊节点
        let specialKey = ''

        for (let i = 0; i < specialNodeList.length; i++) {
          if (it.variable.indexOf(specialNodeList[i]) > -1) {
            specialKey = specialNodeList[i]
            break
          }
        }

        if (specialKey != '') {
          let arr = it.variable.split('.')
          it.variable = [arr[0], specialKey]
        } else {
          it.variable = it.variable.split('.')
        }
        return {
          ...it,
          type: it.type > 0 ? it.type : 1,
          key: Math.random() * 10000
        }
      })

      return {
        ...item,
        is_or: item.is_or ? 1 : 0,
        terms,
        key: Math.random() * 10000
      }
    })

    formState.term = term

    nextTick(() => {
      resetSize()
      updateAnchor()
    })
  } catch (error) {
    console.log(error)
  }
}

const updateAnchor = () => {
  const  nodeSortKey = props.properties.nodeSortKey
  let items =  []
  let ifBoxCountOffsetTop = 56;
  for(let i =0; i < formState.term.length; i++){
    let item  = formState.term[i]
    if (item) {
      let height = 0;
      let top = ifBoxCountOffsetTop;
      let padding = 8;
      let margin  = 2;
      let subItemHeight = 24;
      let subItemNum = item.terms.length;

      height =  subItemNum  * subItemHeight + padding * 2 + margin * (subItemNum - 1);
      
      if(i > 0 &&  i < formState.term.length){
        top += 8
      }

      ifBoxCountOffsetTop = top + height;

      items.push({
        id: nodeSortKey + '-anchor_' + i,
        offsetHeight: height,
        offsetTop: top
      })
    }
  }
 
  // 插入else-box
  if(items.length > 0){
    items.push({
       id: nodeSortKey + '-anchor_right',
      offsetHeight: 24,
      offsetTop: ifBoxCountOffsetTop + 8,
    })
  }

  updateAnchorList(items)
}

watch(() => props.properties.node_params, (newVal, oldVal) => {
  if (newVal != oldVal) {
    init()
  }
}, {
  deep: true
})

onMounted(() => {
  getOptions()

  init()
})
</script>
